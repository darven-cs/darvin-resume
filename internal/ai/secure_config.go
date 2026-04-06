package ai

import (
	"context"
	"strings"

	"Darvin-Resume/internal/crypto"
	"Darvin-Resume/internal/settings"
)

// SaveSecureAPIKey encrypts and stores the API key securely.
// If apiKey is empty, stores empty string (clears the key).
func SaveSecureAPIKey(ctx context.Context, apiKey string) error {
	if apiKey == "" {
		// Store empty string (clears the key)
		return settings.Set(ctx, SettingKeyAPIKey, "")
	}

	// Get device key for encryption
	deviceKey, err := crypto.DeviceKey()
	if err != nil {
		return err
	}

	// Encrypt with device key
	encrypted, err := crypto.Encrypt([]byte(apiKey), deviceKey)
	if err != nil {
		return err
	}

	// Store with "enc:" prefix to indicate encrypted format
	encoded := crypto.EncryptedToHex(encrypted)
	return settings.Set(ctx, SettingKeyAPIKey, encoded)
}

// LoadSecureAPIKey loads and decrypts the API key from secure storage.
// If the stored value is plaintext (no "enc:" prefix), it automatically
// migrates to encrypted storage and returns the plaintext value.
func LoadSecureAPIKey(ctx context.Context) (string, error) {
	value, err := settings.Get(ctx, SettingKeyAPIKey)
	if err != nil {
		if err == settings.ErrSettingNotFound {
			return "", nil
		}
		return "", err
	}

	if value == "" {
		return "", nil
	}

	// Check if it's already encrypted
	if !strings.HasPrefix(value, "enc:") {
		// Plaintext value — auto-migrate to encrypted storage
		if err := SaveSecureAPIKey(ctx, value); err != nil {
			// Migration failed, but we still return the plaintext value for backwards compatibility
			return value, nil
		}
		return value, nil
	}

	// Encrypted value — decrypt it
	deviceKey, err := crypto.DeviceKey()
	if err != nil {
		return "", err
	}

	decrypted, err := crypto.DecryptFromHex(value, deviceKey)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// IsAPIKeyEncrypted checks whether the stored API key is in encrypted format.
func IsAPIKeyEncrypted(ctx context.Context) (bool, error) {
	value, err := settings.Get(ctx, SettingKeyAPIKey)
	if err != nil {
		if err == settings.ErrSettingNotFound {
			return false, nil
		}
		return false, err
	}
	return strings.HasPrefix(value, "enc:"), nil
}

// MigratePlaintextAPIKey explicitly migrates a plaintext API key to encrypted storage.
// Returns true if migration was performed, false if already encrypted or empty.
func MigratePlaintextAPIKey(ctx context.Context) (bool, error) {
	value, err := settings.Get(ctx, SettingKeyAPIKey)
	if err != nil {
		if err == settings.ErrSettingNotFound {
			return false, nil
		}
		return false, err
	}

	if value == "" {
		return false, nil
	}

	if strings.HasPrefix(value, "enc:") {
		return false, nil // Already encrypted
	}

	// Migrate
	if err := SaveSecureAPIKey(ctx, value); err != nil {
		return false, err
	}
	return true, nil
}
