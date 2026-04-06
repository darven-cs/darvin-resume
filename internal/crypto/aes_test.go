package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("sk-ant-api03-test-key-1234567890")
	passphrase := "test-passphrase-123"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypted != plaintext: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDifferentCiphertext(t *testing.T) {
	plaintext := []byte("same plaintext")
	passphrase := "same passphrase"

	enc1, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt 1 failed: %v", err)
	}

	enc2, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt 2 failed: %v", err)
	}

	// Same input should produce different ciphertext (random salt + nonce)
	if bytes.Equal(enc1, enc2) {
		t.Fatal("Expected different ciphertexts for same input, but they are identical")
	}

	// But both should decrypt to the same plaintext
	dec1, _ := Decrypt(enc1, passphrase)
	dec2, _ := Decrypt(enc2, passphrase)

	if !bytes.Equal(dec1, dec2) {
		t.Fatalf("Both ciphertexts should decrypt to same plaintext")
	}
	if !bytes.Equal(dec1, plaintext) {
		t.Fatalf("Decrypted != original plaintext")
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	plaintext := []byte("secret data")
	passphrase := "correct-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Try to decrypt with wrong passphrase
	_, err = Decrypt(encrypted, "wrong-passphrase")
	if err == nil {
		t.Fatal("Expected Decrypt to fail with wrong passphrase, but it succeeded")
	}
	if err != ErrDecryptionFailed {
		t.Fatalf("Expected ErrDecryptionFailed, got: %v", err)
	}
}

func TestDecryptCorruptedData(t *testing.T) {
	plaintext := []byte("original data")
	passphrase := "test-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Corrupt the ciphertext (change a byte in the middle)
	corrupted := make([]byte, len(encrypted))
	copy(corrupted, encrypted)
	if len(corrupted) > 20 {
		corrupted[20] ^= 0xFF
	}

	_, err = Decrypt(corrupted, passphrase)
	if err == nil {
		t.Fatal("Expected Decrypt to fail with corrupted data, but it succeeded")
	}
}

func TestEncryptEmptyPlaintext(t *testing.T) {
	plaintext := []byte("")
	passphrase := "passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt empty plaintext failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt empty ciphertext failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypted != empty plaintext")
	}
}

func TestEncryptLargeData(t *testing.T) {
	// 10KB of data
	plaintext := make([]byte, 10*1024)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}
	passphrase := "large-data-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt large data failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt large data failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypted large data != original (got %d bytes, want %d)", len(decrypted), len(plaintext))
	}
}

func TestEncryptOutputFormat(t *testing.T) {
	plaintext := []byte("test")
	passphrase := "pass"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Expected format: version(1) + salt(16) + nonce(12) + ciphertext+tag
	expectedMinLen := 1 + SaltLength + NonceLength + len(plaintext) + 16
	if len(encrypted) < expectedMinLen {
		t.Fatalf("Encrypted length too short: got %d, want >= %d", len(encrypted), expectedMinLen)
	}

	// Version byte should be 1
	if encrypted[0] != KeyVersion {
		t.Fatalf("Version byte: got %d, want %d", encrypted[0], KeyVersion)
	}
}

func TestEncryptedToHexDecryptFromHex(t *testing.T) {
	plaintext := []byte("api-key-test-12345")
	passphrase := "hex-test-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	hexStr := EncryptedToHex(encrypted)
	if hexStr[:4] != "enc:" {
		t.Fatalf("Hex string should start with 'enc:', got: %s", hexStr[:4])
	}

	decrypted, err := DecryptFromHex(hexStr, passphrase)
	if err != nil {
		t.Fatalf("DecryptFromHex failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypted != plaintext")
	}
}

func TestDecryptFromHexPlaintext(t *testing.T) {
	// Simulate a plaintext value (no "enc:" prefix)
	hexStr := "enc:48656c6c6f576f726c64" // "HelloWorld" in hex

	decrypted, err := DecryptFromHex(hexStr, "wrong")
	if err == nil {
		t.Fatal("Expected DecryptFromHex to fail for non-encrypted hex data")
	}
	_ = decrypted

	// Now test with actual encrypted data
	plaintext := []byte("hello world")
	passphrase := "test"
	encrypted, _ := Encrypt(plaintext, passphrase)
	hexStr2 := EncryptedToHex(encrypted)

	decrypted2, err := DecryptFromHex(hexStr2, passphrase)
	if err != nil {
		t.Fatalf("DecryptFromHex failed: %v", err)
	}
	if !bytes.Equal(decrypted2, plaintext) {
		t.Fatalf("Decrypted != plaintext")
	}
}

func TestReadVersion(t *testing.T) {
	plaintext := []byte("test")
	passphrase := "pass"

	encrypted, _ := Encrypt(plaintext, passphrase)

	version, err := ReadVersion(encrypted)
	if err != nil {
		t.Fatalf("ReadVersion failed: %v", err)
	}
	if version != KeyVersion {
		t.Fatalf("Version: got %d, want %d", version, KeyVersion)
	}

	// Test too-short input
	_, err = ReadVersion([]byte{})
	if err == nil {
		t.Fatal("Expected error for empty input")
	}
}

func TestParseHeader(t *testing.T) {
	plaintext := []byte("test data")
	passphrase := "passphrase"

	encrypted, _ := Encrypt(plaintext, passphrase)

	version, salt, nonce, err := ParseHeader(encrypted)
	if err != nil {
		t.Fatalf("ParseHeader failed: %v", err)
	}

	if version != KeyVersion {
		t.Fatalf("Version: got %d, want %d", version, KeyVersion)
	}
	if len(salt) != SaltLength {
		t.Fatalf("Salt length: got %d, want %d", len(salt), SaltLength)
	}
	if len(nonce) != NonceLength {
		t.Fatalf("Nonce length: got %d, want %d", len(nonce), NonceLength)
	}

	// Too-short input
	_, _, _, err = ParseHeader([]byte{})
	if err == nil {
		t.Fatal("Expected error for short input")
	}
}

func TestGenerateSalt(t *testing.T) {
	salt1, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}
	if len(salt1) != SaltLength {
		t.Fatalf("Salt length: got %d, want %d", len(salt1), SaltLength)
	}

	salt2, _ := GenerateSalt()
	if bytes.Equal(salt1, salt2) {
		t.Fatal("Two generated salts should not be equal")
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce1, err := GenerateNonce()
	if err != nil {
		t.Fatalf("GenerateNonce failed: %v", err)
	}
	if len(nonce1) != NonceLength {
		t.Fatalf("Nonce length: got %d, want %d", len(nonce1), NonceLength)
	}

	nonce2, _ := GenerateNonce()
	if bytes.Equal(nonce1, nonce2) {
		t.Fatal("Two generated nonces should not be equal")
	}
}

func TestEncryptedSize(t *testing.T) {
	for _, size := range []int{0, 1, 100, 1024, 10000} {
		plaintext := make([]byte, size)
		encLen := EncryptedSize(size)
		encrypted, _ := Encrypt(plaintext, "pass")
		if encLen != len(encrypted) {
			t.Fatalf("EncryptedSize(%d): got %d, actual %d", size, encLen, len(encrypted))
		}
	}
}

func TestDeriveKey(t *testing.T) {
	salt := make([]byte, SaltLength)
	for i := range salt {
		salt[i] = byte(i)
	}

	key1 := DeriveKey("passphrase", salt)
	if len(key1) != KeyLength {
		t.Fatalf("Key length: got %d, want %d", len(key1), KeyLength)
	}

	// Same inputs should produce same key
	key2 := DeriveKey("passphrase", salt)
	if !bytes.Equal(key1, key2) {
		t.Fatal("Same inputs should produce same derived key")
	}

	// Different salt should produce different key
	salt2 := make([]byte, SaltLength)
	for i := range salt2 {
		salt2[i] = byte(i + 1)
	}
	key3 := DeriveKey("passphrase", salt2)
	if bytes.Equal(key1, key3) {
		t.Fatal("Different salts should produce different keys")
	}
}
