package backup

import (
	"time"
)

// Auto backup settings keys (stored in settings table)
const (
	SettingKeyAutoBackupEnabled  = "backup.autoEnabled"  // "true" / "false"
	SettingKeyAutoBackupInterval = "backup.autoInterval" // "daily" / "weekly" / "monthly"
)

// Interval values
const (
	IntervalDaily   = "daily"
	IntervalWeekly  = "weekly"
	IntervalMonthly = "monthly"
)

// ParseInterval converts an interval string to time.Duration.
func ParseInterval(interval string) time.Duration {
	switch interval {
	case IntervalDaily:
		return 24 * time.Hour
	case IntervalWeekly:
		return 7 * 24 * time.Hour
	case IntervalMonthly:
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour // Default to daily
	}
}

// IsValidInterval checks if an interval string is valid.
func IsValidInterval(interval string) bool {
	return interval == IntervalDaily || interval == IntervalWeekly || interval == IntervalMonthly
}
