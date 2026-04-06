package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// Version constants
const (
	KeyVersion   = 1
	KeyLength    = 32 // AES-256
	SaltLength   = 16
	NonceLength  = 12 // GCM standard
	argon2Time   = 1
	argon2Memory = 64 * 1024 // 64MB
	argon2Threads = 4
)

// ErrInvalidCiphertext is returned when ciphertext is too short or corrupted.
var ErrInvalidCiphertext = errors.New("invalid ciphertext: too short or corrupted")

// ErrDecryptionFailed is returned when decryption fails (wrong key or corrupted data).
var ErrDecryptionFailed = errors.New("decryption failed: wrong passphrase or corrupted data")

// Encrypt encrypts plaintext using AES-256-GCM with Argon2id key derivation.
// The output format is: version(1) + salt(16) + nonce(12) + ciphertext+tag
func Encrypt(plaintext []byte, passphrase string) ([]byte, error) {
	// Generate random salt
	salt := make([]byte, SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	// Derive key using Argon2id
	key := argon2.IDKey([]byte(passphrase), salt, argon2Time, argon2Memory, argon2Threads, KeyLength)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate random nonce
	nonce := make([]byte, NonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal (encrypt + append auth tag)
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Assemble output: version(1) + salt(16) + nonce(12) + ciphertext+tag
	result := make([]byte, 1+SaltLength+NonceLength+len(ciphertext))
	result[0] = KeyVersion
	copy(result[1:1+SaltLength], salt)
	copy(result[1+SaltLength:1+SaltLength+NonceLength], nonce)
	copy(result[1+SaltLength+NonceLength:], ciphertext)

	// Clear key from memory (best effort in Go)
	for i := range key {
		key[i] = 0
	}

	return result, nil
}

// Decrypt decrypts ciphertext encrypted by Encrypt.
// passphrase must be the same as used during encryption.
func Decrypt(encrypted []byte, passphrase string) ([]byte, error) {
	minLen := 1 + SaltLength + NonceLength + 16 // version + salt + nonce + minimum ciphertext (tag)
	if len(encrypted) < minLen {
		return nil, ErrInvalidCiphertext
	}

	// Parse header
	version := encrypted[0]
	if version != KeyVersion {
		return nil, ErrInvalidCiphertext
	}
	salt := encrypted[1 : 1+SaltLength]
	nonce := encrypted[1+SaltLength : 1+SaltLength+NonceLength]
	ciphertext := encrypted[1+SaltLength+NonceLength:]

	// Derive key using Argon2id
	key := argon2.IDKey([]byte(passphrase), salt, argon2Time, argon2Memory, argon2Threads, KeyLength)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Open (decrypt + verify auth tag)
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	// Clear key from memory (best effort in Go)
	for i := range key {
		key[i] = 0
	}

	return plaintext, nil
}

// EncryptedSize returns the size of encrypted output for a given plaintext size.
// Useful for pre-allocating buffers.
func EncryptedSize(plaintextLen int) int {
	return 1 + SaltLength + NonceLength + plaintextLen + 16 // +16 for GCM tag
}

// MustEncrypt encrypts and panics on error. Use for testing only.
func MustEncrypt(plaintext []byte, passphrase string) []byte {
	out, err := Encrypt(plaintext, passphrase)
	if err != nil {
		panic(err)
	}
	return out
}

// MustDecrypt decrypts and panics on error. Use for testing only.
func MustDecrypt(encrypted []byte, passphrase string) []byte {
	out, err := Decrypt(encrypted, passphrase)
	if err != nil {
		panic(err)
	}
	return out
}

// DeriveKey derives a 32-byte key from passphrase and salt using Argon2id.
// Exported for use by other packages that need consistent key derivation.
func DeriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(passphrase), salt, argon2Time, argon2Memory, argon2Threads, KeyLength)
}

// SaltSize returns the salt length used for key derivation.
func SaltSize() int { return SaltLength }

// NonceSize returns the GCM nonce length.
func NonceSize() int { return NonceLength }

// KeySize returns the derived key length.
func KeySize() int { return KeyLength }

// Version returns the current encryption version.
func Version() int { return KeyVersion }

// ReadVersion extracts the version byte from encrypted data without full decryption.
func ReadVersion(encrypted []byte) (int, error) {
	if len(encrypted) < 1 {
		return 0, ErrInvalidCiphertext
	}
	return int(encrypted[0]), nil
}

// MarshalHeader encodes version + salt + nonce into a byte slice.
// Useful for testing or custom encryption schemes.
func MarshalHeader(version byte, salt, nonce []byte) []byte {
	n := 1 + len(salt) + len(nonce)
	out := make([]byte, n)
	out[0] = version
	copy(out[1:], salt)
	copy(out[1+len(salt):], nonce)
	return out
}

// ParseHeader extracts version, salt, and nonce from encrypted data.
func ParseHeader(encrypted []byte) (version byte, salt, nonce []byte, err error) {
	minLen := 1 + SaltLength + NonceLength
	if len(encrypted) < minLen {
		return 0, nil, nil, ErrInvalidCiphertext
	}
	version = encrypted[0]
	salt = make([]byte, SaltLength)
	nonce = make([]byte, NonceLength)
	copy(salt, encrypted[1:1+SaltLength])
	copy(nonce, encrypted[1+SaltLength:1+SaltLength+NonceLength])
	return version, salt, nonce, nil
}

// GenerateSalt creates a random salt for key derivation.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	_, err := io.ReadFull(rand.Reader, salt)
	return salt, err
}

// GenerateNonce creates a random nonce for AES-GCM.
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, NonceLength)
	_, err := io.ReadFull(rand.Reader, nonce)
	return nonce, err
}

// VersionBytes encodes version as a single byte.
func VersionBytes() byte {
	return KeyVersion
}

// EncryptWithSalt encrypts with a provided salt (for deterministic testing).
// This is only exported for the backup package's use; do not use in production.
func EncryptWithSalt(plaintext []byte, passphrase string, salt []byte) ([]byte, error) {
	if len(salt) != SaltLength {
		return nil, errors.New("salt must be exactly 16 bytes")
	}

	// Derive key using Argon2id
	key := argon2.IDKey([]byte(passphrase), salt, argon2Time, argon2Memory, argon2Threads, KeyLength)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate random nonce
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, err
	}

	// Seal
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Assemble output
	result := make([]byte, 1+SaltLength+NonceLength+len(ciphertext))
	result[0] = KeyVersion
	copy(result[1:1+SaltLength], salt)
	copy(result[1+SaltLength:1+SaltLength+NonceLength], nonce)
	copy(result[1+SaltLength+NonceLength:], ciphertext)

	for i := range key {
		key[i] = 0
	}

	return result, nil
}

// EncryptedToHex converts encrypted bytes to hex string (for settings storage).
func EncryptedToHex(encrypted []byte) string {
	return "enc:" + hexEncode(encrypted)
}

// DecryptFromHex decrypts from hex string (removes "enc:" prefix automatically).
func DecryptFromHex(hexStr string, passphrase string) ([]byte, error) {
	if len(hexStr) < 4 {
		return nil, ErrInvalidCiphertext
	}
	// Remove "enc:" prefix if present
	if hexStr[:4] == "enc:" {
		hexStr = hexStr[4:]
	}
	encrypted, err := hexDecode(hexStr)
	if err != nil {
		return nil, err
	}
	return Decrypt(encrypted, passphrase)
}

// hexEncode encodes bytes to hex string (no external dependency).
func hexEncode(data []byte) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = hexChars[b>>4]
		result[i*2+1] = hexChars[b&0x0f]
	}
	return string(result)
}

// hexDecode decodes hex string to bytes.
func hexDecode(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, errors.New("hex: odd length")
	}
	result := make([]byte, len(s)/2)
	for i := 0; i < len(s)/2; i++ {
		var hi, lo byte
		switch c := s[i*2]; {
		case c >= '0' && c <= '9':
			hi = c - '0'
		case c >= 'a' && c <= 'f':
			hi = c - 'a' + 10
		case c >= 'A' && c <= 'F':
			hi = c - 'A' + 10
		default:
			return nil, errors.New("hex: invalid character")
		}
		switch c := s[i*2+1]; {
		case c >= '0' && c <= '9':
			lo = c - '0'
		case c >= 'a' && c <= 'f':
			lo = c - 'a' + 10
		case c >= 'A' && c <= 'F':
			lo = c - 'A' + 10
		default:
			return nil, errors.New("hex: invalid character")
		}
		result[i] = hi<<4 | lo
	}
	return result, nil
}
