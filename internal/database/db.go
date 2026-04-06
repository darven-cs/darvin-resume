package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
)

// DB is the global database connection
var DB *sql.DB

// Mutex for thread-safe access
var mu sync.Mutex

// Init initializes the database connection and runs migrations
func Init() error {
	mu.Lock()
	defer mu.Unlock()

	// Get user data directory
	userDataDir, err := getUserDataDir()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return err
	}

	// Database file path
	dbPath := filepath.Join(userDataDir, "data.db")

	// Open database connection
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Configure connection pool
	DB.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	DB.SetMaxIdleConns(1)
	DB.SetConnMaxLifetime(0) // Lifetime 0 means connections are reused forever

	// Enable WAL mode for better concurrency
	if _, err := DB.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return err
	}

	// Enable foreign keys
	if _, err := DB.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return err
	}

	// Enable synchronous=NORMAL for better performance without losing durability
	if _, err := DB.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		return err
	}

	// Enable cache_size = -64000 (64MB)
	if _, err := DB.Exec("PRAGMA cache_size=-64000"); err != nil {
		return err
	}

	// Run migrations
	if err := runMigrations(dbPath); err != nil {
		return err
	}

	log.Println("Database initialized at:", dbPath)
	return nil
}

// runMigrations runs database migrations using goose
func runMigrations(dbPath string) error {
	// Get the directory containing migration files
	_, currentFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(currentFile), "migrations")

	// Run migrations with options
	goose.SetDialect("sqlite")
	if err := goose.Up(DB, migrationsDir); err != nil {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}

// Close closes the database connection
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if DB != nil {
		err := DB.Close()
		DB = nil
		return err
	}
	return nil
}

// getUserDataDir returns the platform-specific user data directory
func getUserDataDir() (string, error) {
	// Try to use the standard user data directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Use XDG_DATA_HOME on Linux, ~/Library/Application Support on macOS,
	// or %APPDATA% on Windows
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		switch runtime.GOOS {
		case "darwin":
			dataHome = filepath.Join(homeDir, "Library", "Application Support")
		case "windows":
			dataHome = os.Getenv("APPDATA")
			if dataHome == "" {
				dataHome = filepath.Join(homeDir, "AppData", "Roaming")
			}
		default: // linux and others
			dataHome = filepath.Join(homeDir, ".local", "share")
		}
	}

	return filepath.Join(dataHome, "Darvin-Resume"), nil
}

// GetUserDataDir returns the platform-specific user data directory (exported for backup package).
func GetUserDataDir() (string, error) {
	return getUserDataDir()
}

// Reinit re-initializes the database connection after a restore operation.
// It closes the existing connection and reopens with the same path.
func Reinit() error {
	mu.Lock()
	defer mu.Unlock()

	// Close existing connection if any
	if DB != nil {
		if err := DB.Close(); err != nil {
			return err
		}
		DB = nil
	}

	// Re-initialize using the same path logic as Init
	userDataDir, err := getUserDataDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(userDataDir, "data.db")

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)
	DB.SetConnMaxLifetime(0)

	// Set pragmas
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=-64000",
	}
	for _, pragma := range pragmas {
		if _, err := DB.Exec(pragma); err != nil {
			return err
		}
	}

	// Run migrations
	_, currentFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(currentFile), "migrations")
	goose.SetDialect("sqlite")
	if err := goose.Up(DB, migrationsDir); err != nil {
		return err
	}

	log.Println("Database re-initialized at:", dbPath)
	return nil
}
