package backup

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"Darvin-Resume/internal/crypto"
	"Darvin-Resume/internal/database"
)

// Backup constants
const (
	BackupVersion  = "1.0"
	MaxBackupFiles = 10
	BackupDirName  = "backups"
	AppIdentifier  = "Darvin-Resume"
	BackupFileExt  = ".darvin-backup"
)

// BackupMeta contains metadata about a backup file.
type BackupMeta struct {
	Version     string          `json:"version"`
	App         string          `json:"app"`
	CreatedAt   string          `json:"createdAt"`
	Encrypted   bool            `json:"encrypted"`
	Compression string          `json:"compression"`
	Tables      map[string]int `json:"tables"`
	Checksum    string          `json:"checksum"`
	DBSize      int64           `json:"dbSize"`
	Filename    string          `json:"filename"`
}

// ErrBackupFailed is returned when backup creation fails.
var ErrBackupFailed = errors.New("backup failed")

// ErrRestoreFailed is returned when backup restoration fails.
var ErrRestoreFailed = errors.New("restore failed")

// ErrInvalidBackup is returned when the backup file is invalid.
var ErrInvalidBackup = errors.New("invalid backup file")

// getBackupDir returns the path to the backup directory, creating it if needed.
func getBackupDir() (string, error) {
	userDataDir, err := database.GetUserDataDir()
	if err != nil {
		return "", err
	}
	backupDir := filepath.Join(userDataDir, BackupDirName)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}
	return backupDir, nil
}

// getDBPath returns the path to the main database file.
func getDBPath() (string, error) {
	userDataDir, err := database.GetUserDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userDataDir, "data.db"), nil
}

// getTableCounts returns the row count for each table.
func getTableCounts() (map[string]int, error) {
	tables := []string{"resumes", "settings", "ai_messages", "snapshots"}
	counts := make(map[string]int)

	for _, table := range tables {
		var count int
		err := database.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			// Table might not exist, skip it
			continue
		}
		counts[table] = count
	}

	return counts, nil
}

// checksumFile computes SHA-256 checksum of a file.
func checksumFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// CreateBackup creates a full database backup.
// If passphrase is non-empty, the backup file is encrypted with AES-256-GCM.
// Returns the path to the created backup file.
func CreateBackup(passphrase string) (string, error) {
	backupDir, err := getBackupDir()
	if err != nil {
		return "", fmt.Errorf("get backup dir: %w", err)
	}

	dbPath, err := getDBPath()
	if err != nil {
		return "", fmt.Errorf("get db path: %w", err)
	}

	// Get table counts
	tables, err := getTableCounts()
	if err != nil {
		return "", fmt.Errorf("get table counts: %w", err)
	}

	// Get DB size
	dbInfo, err := os.Stat(dbPath)
	if err != nil {
		return "", fmt.Errorf("stat db: %w", err)
	}

	// Compute checksum of the database file
	checksum, err := checksumFile(dbPath)
	if err != nil {
		return "", fmt.Errorf("checksum db: %w", err)
	}

	// Create metadata
	meta := BackupMeta{
		Version:     BackupVersion,
		App:         AppIdentifier,
		CreatedAt:   time.Now().Format(time.RFC3339),
		Encrypted:   passphrase != "",
		Compression: "none",
		Tables:      tables,
		Checksum:    checksum,
		DBSize:      dbInfo.Size(),
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return "", fmt.Errorf("marshal meta: %w", err)
	}

	// Generate backup filename
	timestamp := time.Now().Format("20060102-150405")
	backupFilename := fmt.Sprintf("darvin-resume-backup-%s%s", timestamp, BackupFileExt)
	backupPath := filepath.Join(backupDir, backupFilename)

	// Create ZIP file
	f, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("create backup file: %w", err)
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	// Add backup.json to ZIP
	jsonEntry, err := zipWriter.Create("backup.json")
	if err != nil {
		return "", fmt.Errorf("create backup.json entry: %w", err)
	}
	if _, err := jsonEntry.Write(metaJSON); err != nil {
		return "", fmt.Errorf("write backup.json: %w", err)
	}

	// VACUUM INTO to a temp file for consistent snapshot
	tempVacuum, err := os.CreateTemp("", "darvin-vacuum-*.db")
	if err != nil {
		return "", fmt.Errorf("create temp vacuum file: %w", err)
	}
	tempVacuumPath := tempVacuum.Name()
	tempVacuum.Close()

	defer os.Remove(tempVacuumPath)

	// Execute VACUUM INTO
	_, err = database.DB.Exec(fmt.Sprintf("VACUUM INTO '%s'", tempVacuumPath))
	if err != nil {
		return "", fmt.Errorf("vacuum into: %w", err)
	}

	// Read vacuumed database
	vacuumData, err := os.ReadFile(tempVacuumPath)
	if err != nil {
		return "", fmt.Errorf("read vacuumed db: %w", err)
	}

	// Add data.db to ZIP
	dbEntry, err := zipWriter.Create("data.db")
	if err != nil {
		return "", fmt.Errorf("create data.db entry: %w", err)
	}
	if _, err := dbEntry.Write(vacuumData); err != nil {
		return "", fmt.Errorf("write data.db: %w", err)
	}

	// Close ZIP before encryption
	if err := zipWriter.Close(); err != nil {
		return "", fmt.Errorf("close zip: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close backup file: %w", err)
	}

	// Encrypt if passphrase provided
	if passphrase != "" {
		encryptedData, err := os.ReadFile(backupPath)
		if err != nil {
			return "", fmt.Errorf("read backup for encryption: %w", err)
		}

		encrypted, err := crypto.Encrypt(encryptedData, passphrase)
		if err != nil {
			return "", fmt.Errorf("encrypt backup: %w", err)
		}

		if err := os.WriteFile(backupPath, encrypted, 0644); err != nil {
			return "", fmt.Errorf("write encrypted backup: %w", err)
		}
	}

	// Cleanup old backups (keep MaxBackupFiles)
	if err := cleanupOldBackups(backupDir, MaxBackupFiles); err != nil {
		// Non-fatal: just log it
		fmt.Printf("warning: cleanup old backups failed: %v\n", err)
	}

	return backupPath, nil
}

// RestoreBackup restores data from a backup file.
// If passphrase is non-empty, the backup file is decrypted before restoration.
// Before restoring, a pre-restore backup is created automatically.
// If restoration fails, the pre-restore backup is used for rollback.
func RestoreBackup(backupPath string, passphrase string) error {
	// Read backup file data
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("read backup file: %w", err)
	}

	// Decrypt if encrypted
	isEncrypted := isEncryptedBackup(data)
	if isEncrypted {
		if passphrase == "" {
			return fmt.Errorf("%w: backup file is encrypted, password required", ErrInvalidBackup)
		}
		decrypted, err := crypto.Decrypt(data, passphrase)
		if err != nil {
			return fmt.Errorf("%w: decryption failed: %v", ErrRestoreFailed, err)
		}
		data = decrypted
	}

	// Open ZIP from memory
	reader, err := zip.NewReader(strings.NewReader(string(data)), int64(len(data)))
	if err != nil {
		// Maybe it's an unencrypted ZIP (no encryption)
		// Try opening as regular ZIP file
		zipFile, openErr := os.Open(backupPath)
		if openErr == nil {
			zipReader, openErr := zip.NewReader(strings.NewReader(string(data)), int64(len(data)))
			if openErr == nil {
				reader = zipReader
				err = nil
			} else {
				zipFile.Close()
			}
		}
		if err != nil {
			return fmt.Errorf("%w: open zip: %v", ErrInvalidBackup, err)
		}
	}

	// Extract and validate backup.json
	var meta BackupMeta
	var dbData []byte

	for _, file := range reader.File {
		switch file.Name {
		case "backup.json":
			rc, err := file.Open()
			if err != nil {
				return fmt.Errorf("%w: open backup.json: %v", ErrInvalidBackup, err)
			}
			jsonData, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return fmt.Errorf("%w: read backup.json: %v", ErrInvalidBackup, err)
			}
			if err := json.Unmarshal(jsonData, &meta); err != nil {
				return fmt.Errorf("%w: parse backup.json: %v", ErrInvalidBackup, err)
			}

			// Validate app identifier
			if meta.App != AppIdentifier {
				return fmt.Errorf("%w: wrong app (got %q, want %q)", ErrInvalidBackup, meta.App, AppIdentifier)
			}

		case "data.db":
			rc, err := file.Open()
			if err != nil {
				return fmt.Errorf("%w: open data.db: %v", ErrInvalidBackup, err)
			}
			dbData, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return fmt.Errorf("%w: read data.db: %v", ErrInvalidBackup, err)
			}
		}
	}

	if meta.App == "" {
		return fmt.Errorf("%w: missing backup.json", ErrInvalidBackup)
	}

	// Create pre-restore backup of current database
	dbPath, err := getDBPath()
	if err != nil {
		return fmt.Errorf("get db path: %w", err)
	}

	// Close database connection
	if err := database.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}

	// Backup current DB with pre-restore prefix
	preRestorePath := dbPath + ".pre-restore"
	if _, err := os.Stat(dbPath); err == nil {
		src, err := os.Open(dbPath)
		if err != nil {
			return fmt.Errorf("open current db for pre-restore backup: %w", err)
		}
		dst, err := os.Create(preRestorePath)
		if err != nil {
			src.Close()
			return fmt.Errorf("create pre-restore backup: %w", err)
		}
		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()
		if err != nil {
			return fmt.Errorf("copy pre-restore backup: %w", err)
		}
	}

	// Write restored database
	if err := os.WriteFile(dbPath, dbData, 0644); err != nil {
		// Rollback to pre-restore
		if _, statErr := os.Stat(preRestorePath); statErr == nil {
			os.Rename(preRestorePath, dbPath)
		}
		return fmt.Errorf("%w: write restored db: %v", ErrRestoreFailed, err)
	}

	// Re-initialize database connection
	if err := database.Reinit(); err != nil {
		// Rollback to pre-restore
		if _, statErr := os.Stat(preRestorePath); statErr == nil {
			os.Remove(dbPath)
			os.Rename(preRestorePath, dbPath)
			database.Reinit()
		}
		return fmt.Errorf("%w: reinitialize database: %v", ErrRestoreFailed, err)
	}

	// Remove pre-restore backup on success
	os.Remove(preRestorePath)

	return nil
}

// ListBackups returns metadata for all backup files in the backup directory.
func ListBackups() ([]BackupMeta, error) {
	backupDir, err := getBackupDir()
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupMeta{}, nil
		}
		return nil, err
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	var backups []BackupMeta
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), BackupFileExt) {
			continue
		}

		backupPath := filepath.Join(backupDir, entry.Name())
		meta, err := readBackupMeta(backupPath)
		if err != nil {
			// If we can't read the meta, still include the file with minimal info
			info, _ := entry.Info()
			meta = BackupMeta{
				Filename: entry.Name(),
				Encrypted: true, // Assume encrypted if we can't read meta
				App:      "",
			}
			if info != nil {
				meta.DBSize = info.Size()
				meta.CreatedAt = info.ModTime().Format(time.RFC3339)
			}
		}
		backups = append(backups, meta)
	}

	// Sort by creation time descending (newest first)
	sort.Slice(backups, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, backups[i].CreatedAt)
		t2, _ := time.Parse(time.RFC3339, backups[j].CreatedAt)
		return t1.After(t2)
	})

	return backups, nil
}

// readBackupMeta reads the metadata from a backup file without full decryption.
func readBackupMeta(backupPath string) (BackupMeta, error) {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return BackupMeta{}, err
	}

	// Check if encrypted
	if isEncryptedBackup(data) {
		// Can't read metadata without password, return minimal info
		info, _ := os.Stat(backupPath)
		meta := BackupMeta{
			Filename:  filepath.Base(backupPath),
			Encrypted: true,
			App:       "",
		}
		if info != nil {
			meta.DBSize = info.Size()
			meta.CreatedAt = info.ModTime().Format(time.RFC3339)
		}
		return meta, nil
	}

	// Try to open as ZIP
	reader, err := zip.NewReader(strings.NewReader(string(data)), int64(len(data)))
	if err != nil {
		// Not a valid ZIP
		info, _ := os.Stat(backupPath)
		meta := BackupMeta{
			Filename:  filepath.Base(backupPath),
			Encrypted: false,
			App:       "",
		}
		if info != nil {
			meta.DBSize = info.Size()
			meta.CreatedAt = info.ModTime().Format(time.RFC3339)
		}
		return meta, nil
	}

	for _, file := range reader.File {
		if file.Name == "backup.json" {
			rc, err := file.Open()
			if err != nil {
				return BackupMeta{}, err
			}
			jsonData, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return BackupMeta{}, err
			}

			var meta BackupMeta
			if err := json.Unmarshal(jsonData, &meta); err != nil {
				return BackupMeta{}, err
			}
			meta.Filename = filepath.Base(backupPath)
			return meta, nil
		}
	}

	// No backup.json found
	info, _ := os.Stat(backupPath)
	meta := BackupMeta{
		Filename:  filepath.Base(backupPath),
		Encrypted: false,
		App:       "",
	}
	if info != nil {
		meta.DBSize = info.Size()
		meta.CreatedAt = info.ModTime().Format(time.RFC3339)
	}
	return meta, nil
}

// isEncryptedBackup checks if the data appears to be encrypted (starts with version byte).
func isEncryptedBackup(data []byte) bool {
	// Encrypted data starts with version byte (1) + salt, etc.
	// ZIP files start with "PK" (0x50, 0x4B)
	if len(data) < 2 {
		return false
	}
	return data[0] != 0x50 || data[1] != 0x4B
}

// cleanupOldBackups removes the oldest backup files beyond maxBackups.
func cleanupOldBackups(backupDir string, maxBackups int) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	// Collect backup files with their modification times
	type backupEntry struct {
		name    string
		path    string
		modTime time.Time
	}
	var backups []backupEntry

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), BackupFileExt) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		backups = append(backups, backupEntry{
			name:    entry.Name(),
			path:    filepath.Join(backupDir, entry.Name()),
			modTime: info.ModTime(),
		})
	}

	// Sort by modification time oldest first
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].modTime.Before(backups[j].modTime)
	})

	// Remove oldest files beyond maxBackups
	if len(backups) > maxBackups {
		toDelete := backups[:len(backups)-maxBackups]
		for _, b := range toDelete {
			if err := os.Remove(b.path); err != nil {
				// Non-fatal, log and continue
				fmt.Printf("warning: failed to remove old backup %s: %v\n", b.path, err)
			}
		}
	}

	return nil
}

// GetBackupDir returns the backup directory path.
func GetBackupDir() (string, error) {
	return getBackupDir()
}
