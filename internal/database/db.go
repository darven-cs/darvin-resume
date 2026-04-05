package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// DB is the global database connection
var DB *sql.DB

// Init initializes the database connection
func Init() error {
	var err error
	DB, err = sql.Open("sqlite", "data.db")
	if err != nil {
		return err
	}

	// Enable WAL mode
	_, err = DB.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return err
	}

	// Enable foreign keys
	_, err = DB.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		return err
	}

	log.Println("Database initialized")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
