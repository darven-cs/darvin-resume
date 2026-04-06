package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ErrDeviceKeyUnavailable is returned when the device key cannot be obtained.
type ErrDeviceKeyUnavailable struct {
	Reason string
}

func (e *ErrDeviceKeyUnavailable) Error() string {
	return "device key unavailable: " + e.Reason
}

// DeviceKey returns a unique, deterministic device identifier for the current platform.
// The key is derived from platform-specific identifiers and is never stored on disk.
// Format: SHA-256 hex of "Darvin-Resume-DeviceKey-V1:" + platform_seed
func DeviceKey() (string, error) {
	seed, err := getPlatformSeed()
	if err != nil {
		return "", &ErrDeviceKeyUnavailable{Reason: err.Error()}
	}

	// Prepend application identity so the same machine gets different keys for different apps
	prefix := "Darvin-Resume-DeviceKey-V1:"
	combined := prefix + seed

	// SHA-256 hash for consistent 64-char hex output
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:]), nil
}

// getPlatformSeed returns a platform-specific seed string.
func getPlatformSeed() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsMachineGuid()
	case "darwin":
		return getMacOSUUID()
	case "linux":
		return getLinuxMachineId()
	default:
		return "", &ErrDeviceKeyUnavailable{Reason: "unsupported platform: " + runtime.GOOS}
	}
}

// getWindowsMachineGuid reads the Windows MachineGuid from the registry.
func getWindowsMachineGuid() (string, error) {
	cmd := exec.Command("powershell", "-Command",
		"Get-ItemProperty 'HKLM:\\SOFTWARE\\Microsoft\\Cryptography' | Select-Object -ExpandProperty MachineGuid")
	output, err := cmd.Output()
	if err != nil {
		return "", &ErrDeviceKeyUnavailable{Reason: "powershell failed: " + err.Error()}
	}
	guid := strings.TrimSpace(string(output))
	if guid == "" {
		return "", &ErrDeviceKeyUnavailable{Reason: "MachineGuid is empty"}
	}
	return guid, nil
}

// getMacOSUUID reads the IOPlatformUUID from IORegistry.
func getMacOSUUID() (string, error) {
	cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	output, err := cmd.Output()
	if err != nil {
		return "", &ErrDeviceKeyUnavailable{Reason: "ioreg failed: " + err.Error()}
	}

	// Extract IOPlatformUUID value from output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IOPlatformUUID") {
			// Format: "           \"IOPlatformUUID\" = \"XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX\""
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				uuid := parts[3]
				if uuid != "" {
					return uuid, nil
				}
			}
		}
	}

	return "", &ErrDeviceKeyUnavailable{Reason: "IOPlatformUUID not found in ioreg output"}
}

// getLinuxMachineId reads the machine-id from standard Linux paths.
func getLinuxMachineId() (string, error) {
	// Try /etc/machine-id first (standard location)
	data, err := os.ReadFile("/etc/machine-id")
	if err == nil {
		id := strings.TrimSpace(string(data))
		if id != "" {
			return id, nil
		}
	}

	// Fallback to /var/lib/dbus/machine-id
	data, err = os.ReadFile("/var/lib/dbus/machine-id")
	if err == nil {
		id := strings.TrimSpace(string(data))
		if id != "" {
			return id, nil
		}
	}

	return "", &ErrDeviceKeyUnavailable{Reason: "no machine-id found in /etc/machine-id or /var/lib/dbus/machine-id"}
}

// IsErrDeviceKeyUnavailable checks if an error is of type ErrDeviceKeyUnavailable.
func IsErrDeviceKeyUnavailable(err error) bool {
	_, ok := err.(*ErrDeviceKeyUnavailable)
	return ok
}
