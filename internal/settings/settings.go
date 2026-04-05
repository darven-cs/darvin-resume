package settings

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Darvin-Resume/internal/database"
)

// ErrSettingNotFound is returned when a setting key does not exist.
var ErrSettingNotFound = errors.New("setting not found")

// Get retrieves a setting value by key.
// Returns ErrSettingNotFound if the key does not exist.
func Get(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM settings WHERE key = ?`
	var value string
	err := database.DB.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrSettingNotFound
		}
		return "", err
	}
	return value, nil
}

// Set creates or updates a setting value.
func Set(ctx context.Context, key string, value string) error {
	query := `
		INSERT INTO settings (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
	`
	_, err := database.DB.ExecContext(ctx, query, key, value, time.Now())
	return err
}

// Delete removes a setting by key.
func Delete(ctx context.Context, key string) error {
	query := `DELETE FROM settings WHERE key = ?`
	_, err := database.DB.ExecContext(ctx, query, key)
	return err
}

// GetWithDefault retrieves a setting, returning defaultValue if not found.
func GetWithDefault(ctx context.Context, key string, defaultValue string) string {
	val, err := Get(ctx, key)
	if err != nil {
		return defaultValue
	}
	return val
}
