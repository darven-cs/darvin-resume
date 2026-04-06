package crypto

import (
	"testing"
)

func TestDeviceKey(t *testing.T) {
	key, err := DeviceKey()
	if err != nil {
		t.Fatalf("DeviceKey failed: %v", err)
	}

	// Should be a valid 64-character hex string (SHA-256)
	if len(key) != 64 {
		t.Fatalf("DeviceKey length: got %d, want 64", len(key))
	}

	// Should be all hex characters
	for i, c := range key {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Fatalf("DeviceKey char %d is not hex: %c", i, c)
		}
	}
}

func TestDeviceKeyConsistent(t *testing.T) {
	key1, err := DeviceKey()
	if err != nil {
		t.Fatalf("DeviceKey 1 failed: %v", err)
	}

	key2, err := DeviceKey()
	if err != nil {
		t.Fatalf("DeviceKey 2 failed: %v", err)
	}

	if key1 != key2 {
		t.Fatalf("DeviceKey is not consistent: got %q then %q", key1, key2)
	}
}

func TestIsErrDeviceKeyUnavailable(t *testing.T) {
	// Should return false for nil
	if IsErrDeviceKeyUnavailable(nil) {
		t.Fatal("IsErrDeviceKeyUnavailable(nil) should return false")
	}

	// Should return true for ErrDeviceKeyUnavailable
	err := &ErrDeviceKeyUnavailable{Reason: "test"}
	if !IsErrDeviceKeyUnavailable(err) {
		t.Fatal("IsErrDeviceKeyUnavailable should return true for ErrDeviceKeyUnavailable")
	}

	// Should return false for other errors
	err2 := ErrInvalidCiphertext
	if IsErrDeviceKeyUnavailable(err2) {
		t.Fatal("IsErrDeviceKeyUnavailable should return false for non-device errors")
	}
}

func TestGetPlatformSeed(t *testing.T) {
	seed, err := getPlatformSeed()
	if err != nil {
		t.Fatalf("getPlatformSeed failed: %v", err)
	}

	if seed == "" {
		t.Fatal("getPlatformSeed returned empty string")
	}

	// Seed should be consistent
	seed2, _ := getPlatformSeed()
	if seed != seed2 {
		t.Fatalf("getPlatformSeed not consistent: got %q then %q", seed, seed2)
	}
}
