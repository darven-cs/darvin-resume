# Phase 6 研究 — 数据安全与备份技术方案

> 研究日期: 2026-04-06
> 项目: Darvin-Resume (Wails v2 + Go + Vue 3 + SQLite)
> 状态: 研究完成，待编写执行计划

---

## 一、项目现状分析

### 1.1 当前数据存储

| 数据 | 存储位置 | 安全级别 | 表 |
|------|---------|---------|-----|
| 简历内容 (JSON + Markdown) | SQLite `resumes` | 高 (含个人隐私) | `resumes` |
| AI API Key | SQLite `settings` | **极高** (明文!) | `settings` |
| AI 对话历史 | SQLite `ai_messages` | 中 | `ai_messages` |
| 版本快照 | SQLite `snapshots` | 高 (含个人隐私) | `snapshots` |
| 应用设置 | SQLite `settings` | 低 | `settings` |

### 1.2 关键安全风险

- **AI API Key 明文存储** (config.go L24 注释已确认: "APIKey is stored in plaintext this phase")
- **简历个人信息无加密** (姓名、电话、邮箱、地址等 PII 数据明文存储)
- **无自动备份机制** (数据库损坏 = 数据丢失)
- **无备份导出功能** (用户无法手动导出/导入数据)

### 1.3 技术约束

- 使用 `modernc.org/sqlite` (纯 Go，无 CGO 依赖) — 不支持 SQLite 的 C 扩展
- `golang.org/x/crypto` 已在 go.mod 间接依赖中 (v0.48.0)
- 三平台支持: Windows / macOS / Linux
- **只用 Go 标准库 + golang.org/x/crypto** — 零新增外部依赖

---

## 二、AES-256-GCM 加密方案

### 2.1 技术选型

| 组件 | 选择 | 理由 |
|------|------|------|
| 加密算法 | AES-256-GCM | Go 标准库 `crypto/aes` + `crypto/cipher` 原生支持，认证加密(AEAD) |
| 密钥派生 | Argon2id | `golang.org/x/crypto/argon2`，抗 GPU/ASIC 暴力破解，2024 年密码哈希竞赛冠军 |
| 随机数 | `crypto/rand` | Go 标准库 CSPRNG |
| 密钥存储 | 平台原生 Keyring | 零依赖方案见下文 |

### 2.2 加密数据格式

```
┌──────────┬──────────┬──────────────┬──────────────┬──────────────┐
│ Version  │   Salt   │  Nonce(IV)   │  Ciphertext  │  Auth Tag    │
│ 1 byte   │ 16 bytes │ 12 bytes     │  N bytes     │  16 bytes    │
└──────────┴──────────┴──────────────┴──────────────┴──────────────┘
```

总开销 = 1 + 16 + 12 + 16 = **45 字节**固定开销（对简历数据可忽略不计）

### 2.3 核心加密代码设计

```go
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	KeyVersion    = 1
	KeyLength     = 32       // AES-256
	SaltLength    = 16
	NonceLength   = 12       // GCM 标准
	Argon2Time    = 1        // 迭代次数
	Argon2Memory  = 64 * 1024 // 64MB 内存
	Argon2Threads = 4        // 并行线程
)

// Encrypt 使用 AES-256-GCM 加密明文
// passphrase: 用户密码或设备密钥
// 返回格式: version(1) + salt(16) + nonce(12) + ciphertext + tag(16)
func Encrypt(plaintext []byte, passphrase string) ([]byte, error) {
	// 1. 生成随机 salt
	salt := make([]byte, SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	// 2. Argon2id 派生密钥
	key := argon2.IDKey(
		[]byte(passphrase), salt,
		Argon2Time, Argon2Memory, Argon2Threads, KeyLength,
	)

	// 3. 创建 AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	// 4. 生成随机 nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	// 5. 加密 + 认证 (Seal 会追加 16 字节 auth tag)
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// 6. 组装最终输出: version + salt + nonce + ciphertext(tag 内含)
	result := make([]byte, 0, 1+SaltLength+NonceLength+len(ciphertext))
	result = append(result, KeyVersion)
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	// 清除内存中的密钥
	for i := range key {
		key[i] = 0
	}

	return result, nil
}

// Decrypt 使用 AES-256-GCM 解密密文
func Decrypt(encrypted []byte, passphrase string) ([]byte, error) {
	// 1. 解析头部
	if len(encrypted) < 1+SaltLength+NonceLength+aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	version := encrypted[0]
	_ = version // 预留版本兼容
	salt := encrypted[1 : 1+SaltLength]
	nonce := encrypted[1+SaltLength : 1+SaltLength+NonceLength]
	ciphertext := encrypted[1+SaltLength+NonceLength:]

	// 2. Argon2id 派生同一密钥
	key := argon2.IDKey(
		[]byte(passphrase), salt,
		Argon2Time, Argon2Memory, Argon2Threads, KeyLength,
	)

	// 3. 解密 + 验证
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed (wrong password or corrupted data): %w", err)
	}

	// 清除密钥
	for i := range key {
		key[i] = 0
	}

	return plaintext, nil
}
```

### 2.4 设备唯一密钥方案（零外部依赖）

**方案核心: 基于设备特征生成确定性密钥，不存储密钥本身。**

```go
package crypto

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DeviceKey 基于设备硬件特征生成唯一密钥种子
// 不使用任何外部依赖，纯 Go 标准库
func DeviceKey() (string, error) {
	var seed string

	switch runtime.GOOS {
	case "windows":
		// Windows: 读取注册表 MachineGuid
		out, err := exec.Command("powershell",
			"-Command",
			"Get-ItemProperty 'HKLM:\\SOFTWARE\\Microsoft\\Cryptography' | Select-Object -ExpandProperty MachineGuid",
		).Output()
		if err != nil {
			return "", fmt.Errorf("read Windows MachineGuid: %w", err)
		}
		seed = strings.TrimSpace(string(out))

	case "darwin":
		// macOS: 读取 IOPlatformUUID
		out, err := exec.Command("ioreg",
			"-rd1", "-c", "IOPlatformExpertDevice",
		).Output()
		if err != nil {
			return "", fmt.Errorf("read macOS IOPlatformUUID: %w", err)
		}
		// 从输出中提取 UUID
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "IOPlatformUUID") {
				parts := strings.Split(line, `"`)
				if len(parts) >= 4 {
					seed = parts[3]
				}
			}
		}

	case "linux":
		// Linux: 读取 /etc/machine-id
		data, err := os.ReadFile("/etc/machine-id")
		if err != nil {
			// fallback: /var/lib/dbus/machine-id
			data, err = os.ReadFile("/var/lib/dbus/machine-id")
			if err != nil {
				return "", fmt.Errorf("read Linux machine-id: %w", err)
			}
		}
		seed = strings.TrimSpace(string(data))
	}

	if seed == "" {
		return "", fmt.Errorf("failed to generate device key on %s", runtime.GOOS)
	}

	// 用应用名 + 设备 ID 做 SHA-256 哈希，作为最终密钥种子
	// 这样同一设备不同应用会得到不同密钥
	appSeed := "Darvin-Resume-DeviceKey-V1:" + seed
	hash := sha256.Sum256([]byte(appSeed))
	return fmt.Sprintf("%x", hash), nil
}
```

### 2.5 加密策略分层

| 层级 | 保护对象 | 加密方式 | 触发时机 |
|------|---------|---------|---------|
| **L1 - API Key 加密** | `settings` 表中的 `ai.apiKey` | AES-256-GCM + 设备密钥 | 应用启动时解密到内存，保存时加密 |
| **L2 - 备份文件加密** | 导出的 .darvin-backup 文件 | AES-256-GCM + 用户密码 | 用户导出备份时 |
| **L3 - 全库加密** | (未来) 整个 SQLite 文件 | SQLCipher 或文件级加密 | 不在本期范围 |

### 2.6 密钥管理架构

```
┌──────────────────────────────────────────────────────┐
│                    Darvin-Resume                      │
├──────────────────────────────────────────────────────┤
│                                                      │
│  ┌─────────────────┐    ┌──────────────────────┐     │
│  │  DeviceKey()     │───>│  SHA-256 哈希种子     │     │
│  │  (设备特征)      │    │  (32字节密钥材料)     │     │
│  └─────────────────┘    └──────────┬───────────┘     │
│                                    │                  │
│                                    ▼                  │
│                         ┌─────────────────────┐      │
│                         │  Argon2id KDF        │      │
│                         │  + 随机 salt          │      │
│                         │  = 32 字节 AES 密钥   │      │
│                         └──────────┬──────────┘      │
│                                    │                  │
│                                    ▼                  │
│                         ┌─────────────────────┐      │
│                         │  AES-256-GCM         │      │
│                         │  加密/解密数据        │      │
│                         └─────────────────────┘      │
│                                                      │
│  ┌──────────────────────────────────────────────┐   │
│  │  备份文件加密: 用户密码 → Argon2id → AES-256  │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────┘
```

**关键设计决策:**

1. **设备密钥不存储**: 通过设备特征每次动态计算，无需存储密钥文件
2. **每次加密用新 salt + nonce**: 相同明文每次加密产出不同密文，防止重放攻击
3. **GCM 认证标签**: 确保数据完整性，篡改 = 解密失败
4. **内存中密钥用后清零**: `for i := range key { key[i] = 0 }`

---

## 三、SQLite 数据库备份方案

### 3.1 备份方式对比

| 方式 | 一致性 | 性能 | 适用场景 | 推荐度 |
|------|--------|------|---------|--------|
| **VACUUM INTO** | 完美 | 中 | SQLite 3.27+ (modernc 支持) | 推荐 |
| 文件复制 (热) | 可能不一致 | 快 | 不推荐 WAL 模式下 | 不推荐 |
| 文件复制 (冷) | 完美 | 快 | 关闭数据库后 | 可选 |
| SQLite Backup API | 完美 | 慢 | CGO 版本 | 不适用 (modernc) |

### 3.2 推荐方案: VACUUM INTO + archive/zip

**为什么选 VACUUM INTO:**
- modernc.org/sqlite 支持 VACUUM INTO（纯 Go，无 CGO）
- 在 WAL 模式下生成一致的数据库快照
- 不需要停止数据库连接
- 输出是一个完整的、经过优化的 SQLite 文件

### 3.3 备份文件格式设计

**文件扩展名: `.darvin-backup`**

```
darvin-resume-backup-20260406-143022.darvin-backup
│
├── [ZIP 结构]
│   ├── backup.json          ← 元数据清单
│   ├── data.db              ← SQLite 数据库文件 (VACUUM INTO 导出)
│   └── checksum.sha256      ← data.db 的 SHA-256 校验
│
└── [可选加密]
    整个 ZIP 内容用 AES-256-GCM 加密后存储
```

### 3.4 backup.json 元数据格式

```json
{
  "version": "1.0",
  "app": "Darvin-Resume",
  "createdAt": "2026-04-06T14:30:22+08:00",
  "encrypted": true,
  "compression": "zip",
  "tables": {
    "resumes": {"count": 5},
    "settings": {"count": 8},
    "ai_messages": {"count": 42},
    "snapshots": {"count": 15}
  },
  "checksum": "sha256:abcdef1234567890...",
  "dbSize": 245760
}
```

### 3.5 核心备份代码设计

```go
package backup

import (
	"archive/zip"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"Darvin-Resume/internal/database"
)

const (
	BackupVersion    = "1.0"
	MaxBackupFiles   = 10          // 最大保留备份数
	BackupDirName    = "backups"   // 备份子目录名
)

// BackupMeta 备份文件元数据
type BackupMeta struct {
	Version     string            `json:"version"`
	App         string            `json:"app"`
	CreatedAt   string            `json:"createdAt"`
	Encrypted   bool              `json:"encrypted"`
	Compression string            `json:"compression"`
	Tables      map[string]TableInfo `json:"tables"`
	Checksum    string            `json:"checksum"`
	DBSize      int64             `json:"dbSize"`
}

type TableInfo struct {
	Count int `json:"count"`
}

// getBackupDir 返回备份目录路径
func getBackupDir() (string, error) {
	userDataDir, err := getUserDataDir()
	if err != nil {
		return "", err
	}
	backupDir := filepath.Join(userDataDir, BackupDirName)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}
	return backupDir, nil
}

// CreateBackup 创建数据库备份
// passphrase: 空字符串 = 不加密，非空 = AES-256-GCM 加密
func CreateBackup(passphrase string) (string, error) {
	backupDir, err := getBackupDir()
	if err != nil {
		return "", fmt.Errorf("get backup dir: %w", err)
	}

	// 1. VACUUM INTO 临时文件 (生成一致的 SQLite 快照)
	tmpDB, err := os.CreateTemp(backupDir, "vacuum-*.tmp")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpDBPath := tmpDB.Name()
	tmpDB.Close()
	defer os.Remove(tmpDBPath)

	// 执行 VACUUM INTO
	_, err = database.DB.Exec(fmt.Sprintf("VACUUM INTO '%s'", tmpDBPath))
	if err != nil {
		return "", fmt.Errorf("vacuum into: %w", err)
	}

	// 2. 读取临时数据库文件
	dbData, err := os.ReadFile(tmpDBPath)
	if err != nil {
		return "", fmt.Errorf("read vacuumed db: %w", err)
	}

	// 3. 计算 SHA-256 校验和
	hash := sha256.Sum256(dbData)
	checksum := hex.EncodeToString(hash[:])

	// 4. 获取各表行数
	tables, err := getTableCounts()
	if err != nil {
		return "", fmt.Errorf("get table counts: %w", err)
	}

	// 5. 构建元数据
	meta := BackupMeta{
		Version:     BackupVersion,
		App:         "Darvin-Resume",
		CreatedAt:   time.Now().Format(time.RFC3339),
		Encrypted:   passphrase != "",
		Compression: "zip",
		Tables:      tables,
		Checksum:    "sha256:" + checksum,
		DBSize:      int64(len(dbData)),
	}

	metaJSON, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal meta: %w", err)
	}

	// 6. 生成备份文件名
	filename := fmt.Sprintf(
		"darvin-resume-backup-%s.darvin-backup",
		time.Now().Format("20060102-150405"),
	)
	backupPath := filepath.Join(backupDir, filename)

	// 7. 创建 ZIP
	if err := createBackupZip(backupPath, metaJSON, dbData); err != nil {
		return "", fmt.Errorf("create zip: %w", err)
	}

	// 8. 如果需要加密，对 ZIP 文件整体加密
	if passphrase != "" {
		if err := encryptBackupFile(backupPath, passphrase); err != nil {
			os.Remove(backupPath) // 清理失败的文件
			return "", fmt.Errorf("encrypt backup: %w", err)
		}
	}

	// 9. 清理旧备份
	if err := cleanupOldBackups(backupDir, MaxBackupFiles); err != nil {
		// 清理失败不影响主流程，只记录日志
		fmt.Printf("warning: cleanup old backups failed: %v\n", err)
	}

	return backupPath, nil
}

// createBackupZip 创建包含元数据和数据库的 ZIP 文件
func createBackupZip(path string, metaJSON []byte, dbData []byte) error {
	zipFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 写入 backup.json
	metaWriter, err := zipWriter.Create("backup.json")
	if err != nil {
		return err
	}
	if _, err := metaWriter.Write(metaJSON); err != nil {
		return err
	}

	// 写入 data.db
	dbWriter, err := zipWriter.Create("data.db")
	if err != nil {
		return err
	}
	if _, err := dbWriter.Write(dbData); err != nil {
		return err
	}

	return nil
}

// getTableCounts 获取各表行数
func getTableCounts() (map[string]TableInfo, error) {
	tables := []string{"resumes", "settings", "ai_messages", "snapshots"}
	result := make(map[string]TableInfo)

	for _, table := range tables {
		var count int
		err := database.DB.QueryRow(
			fmt.Sprintf("SELECT COUNT(*) FROM %s", table),
		).Scan(&count)
		if err != nil {
			return nil, err
		}
		result[table] = TableInfo{Count: count}
	}

	return result, nil
}

// encryptBackupFile 对备份文件进行 AES-256-GCM 加密
func encryptBackupFile(path string, passphrase string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	encrypted, err := Encrypt(data, passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encrypted, 0644)
}
```

### 3.6 恢复 (导入) 备份设计

```go
package backup

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// RestoreBackup 从备份文件恢复数据库
// passphrase: 如果备份是加密的，提供密码；否则为空
func RestoreBackup(backupPath string, passphrase string) error {
	// 1. 读取备份文件
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("read backup file: %w", err)
	}

	// 2. 如果有密码，先解密
	if passphrase != "" {
		data, err = Decrypt(data, passphrase)
		if err != nil {
			return fmt.Errorf("decrypt backup: %w", err)
		}
	}

	// 3. 解析 ZIP
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("parse zip: %w", err)
	}

	var metaJSON []byte
	var dbData []byte

	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}

		switch file.Name {
		case "backup.json":
			metaJSON = content
		case "data.db":
			dbData = content
		}
	}

	if metaJSON == nil || dbData == nil {
		return fmt.Errorf("invalid backup: missing backup.json or data.db")
	}

	// 4. 验证元数据
	var meta BackupMeta
	if err := json.Unmarshal(metaJSON, &meta); err != nil {
		return fmt.Errorf("parse backup meta: %w", err)
	}
	if meta.App != "Darvin-Resume" {
		return fmt.Errorf("invalid backup: not a Darvin-Resume backup")
	}

	// 5. TODO: 验证校验和

	// 6. 将数据库文件写入目标位置
	// 关闭当前数据库连接
	database.Close()

	userDataDir, err := getUserDataDir()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(userDataDir, "data.db")

	// 备份当前数据库（以防恢复失败）
	currentBackup := dbPath + ".pre-restore"
	if err := os.Rename(dbPath, currentBackup); err != nil {
		return fmt.Errorf("backup current db: %w", err)
	}

	// 写入新数据库
	if err := os.WriteFile(dbPath, dbData, 0644); err != nil {
		// 恢复失败，回滚
		os.Rename(currentBackup, dbPath)
		return fmt.Errorf("write restored db: %w", err)
	}

	// 成功后删除恢复前的备份
	os.Remove(currentBackup)

	// 7. 重新打开数据库
	if err := database.Init(); err != nil {
		return fmt.Errorf("reinit database: %w", err)
	}

	return nil
}
```

---

## 四、定时备份机制

### 4.1 方案设计

**使用 Go 标准库 `time.Ticker` + `context.Context` 实现可取消的定时任务，无需外部依赖。**

### 4.2 核心设计

```go
package backup

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DefaultInterval = 30 * time.Minute  // 默认 30 分钟自动备份
)

// Scheduler 定时备份调度器
type Scheduler struct {
	mu       sync.Mutex
	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	running  bool
}

// NewScheduler 创建备份调度器
func NewScheduler(interval time.Duration) *Scheduler {
	if interval < 5*time.Minute {
		interval = 5 * time.Minute // 最小间隔 5 分钟
	}
	return &Scheduler{
		interval: interval,
	}
}

// Start 启动定时备份
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("scheduler already running")
	}

	schedCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.running = true

	s.wg.Add(1)
	go s.run(schedCtx)

	log.Printf("[BackupScheduler] started with interval %v", s.interval)
	return nil
}

// run 定时任务主循环
func (s *Scheduler) run(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[BackupScheduler] stopped")
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
			return

		case <-ticker.C:
			// 执行自动备份 (使用设备密钥，无需用户密码)
			deviceKey, err := DeviceKey()
			if err != nil {
				log.Printf("[BackupScheduler] get device key failed: %v", err)
				continue
			}

			path, err := CreateBackup(deviceKey)
			if err != nil {
				log.Printf("[BackupScheduler] backup failed: %v", err)
				continue
			}
			log.Printf("[BackupScheduler] auto backup created: %s", path)
		}
	}
}

// Stop 停止定时备份
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
	s.running = false
	log.Println("[BackupScheduler] stopped gracefully")
}

// IsRunning 返回调度器是否在运行
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
```

### 4.3 旧备份清理机制

```go
// cleanupOldBackups 清理旧备份文件，保留最近 maxBackups 个
func cleanupOldBackups(backupDir string, maxBackups int) error {
	// 1. 列出所有 .darvin-backup 文件
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	var backups []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".darvin-backup") {
			backups = append(backups, entry)
		}
	}

	// 2. 如果不超过限制，无需清理
	if len(backups) <= maxBackups {
		return nil
	}

	// 3. 按修改时间排序 (旧 -> 新)
	sort.Slice(backups, func(i, j int) bool {
		info1, _ := backups[i].Info()
		info2, _ := backups[j].Info()
		return info1.ModTime().Before(info2.ModTime())
	})

	// 4. 删除最旧的文件
	deleteCount := len(backups) - maxBackups
	for i := 0; i < deleteCount; i++ {
		path := filepath.Join(backupDir, backups[i].Name())
		if err := os.Remove(path); err != nil {
			log.Printf("warning: failed to delete old backup %s: %v", path, err)
		} else {
			log.Printf("[BackupCleanup] deleted old backup: %s", path)
		}
	}

	return nil
}
```

### 4.4 Wails 应用集成方式

```go
// 在 app.go 中的集成方式

type App struct {
	ctx         context.Context
	svc         service.ResumeService
	backupSched *backup.Scheduler  // 新增
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// ... 现有数据库初始化 ...

	// 启动定时备份 (30 分钟间隔)
	a.backupSched = backup.NewScheduler(30 * time.Minute)
	if err := a.backupSched.Start(ctx); err != nil {
		log.Printf("warning: backup scheduler failed to start: %v", err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	// 停止备份调度器
	if a.backupSched != nil {
		a.backupSched.Stop()
	}

	// ... 现有数据库关闭 ...
}

// --- 前端可调用的 Bridge 方法 ---

// CreateManualBackup 手动创建备份 (用户可选密码)
func (a *App) CreateManualBackup(password string) (string, error) {
	return backup.CreateBackup(password)
}

// RestoreFromBackup 从备份文件恢复
func (a *App) RestoreFromBackup(filePath string, password string) error {
	return backup.RestoreBackup(filePath, password)
}

// ListBackups 列出所有本地备份
func (a *App) ListBackups() ([]map[string]interface{}, error) {
	return backup.ListBackups()
}

// ShowOpenBackupDialog 打开文件选择对话框选择备份文件
func (a *App) ShowOpenBackupDialog() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择备份文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "Darvin Resume 备份", Pattern: "*.darvin-backup"},
		},
	})
	return path, err
}
```

---

## 五、API Key 安全存储方案

### 5.1 方案: 设备密钥加密 + SQLite 存储

**不引入 keyring 外部依赖，而是使用设备密钥加密后存入现有 settings 表。**

```go
// internal/ai/secure_config.go

package ai

import (
	"context"
	"fmt"

	"Darvin-Resume/internal/crypto"
	"Darvin-Resume/internal/settings"
)

// SaveSecureAPIKey 安全保存 API Key (加密后存储)
func SaveSecureAPIKey(ctx context.Context, apiKey string) error {
	if apiKey == "" {
		return settings.Set(ctx, SettingKeyAPIKey, "")
	}

	deviceKey, err := crypto.DeviceKey()
	if err != nil {
		return fmt.Errorf("get device key: %w", err)
	}

	encrypted, err := crypto.Encrypt([]byte(apiKey), deviceKey)
	if err != nil {
		return fmt.Errorf("encrypt api key: %w", err)
	}

	// 存储为 hex 编码
	encoded := fmt.Sprintf("enc:%x", encrypted)
	return settings.Set(ctx, SettingKeyAPIKey, encoded)
}

// LoadSecureAPIKey 安全读取 API Key (解密)
func LoadSecureAPIKey(ctx context.Context) (string, error) {
	value, err := settings.Get(ctx, SettingKeyAPIKey)
	if err != nil {
		return "", err
	}

	// 空值或未加密的旧值直接返回
	if value == "" || !strings.HasPrefix(value, "enc:") {
		return value, nil
	}

	// 解密
	hexStr := strings.TrimPrefix(value, "enc:")
	encrypted, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("decode encrypted key: %w", err)
	}

	deviceKey, err := crypto.DeviceKey()
	if err != nil {
		return "", fmt.Errorf("get device key: %w", err)
	}

	decrypted, err := crypto.Decrypt(encrypted, deviceKey)
	if err != nil {
		return "", fmt.Errorf("decrypt api key: %w", err)
	}

	return string(decrypted), nil
}
```

### 5.2 迁移策略

```
启动时检测:
  settings["ai.apiKey"] 的值
  ├── 空字符串 → 无需处理
  ├── "enc:..." → 已加密，正常解密
  └── 其他 (明文 API Key)
      → 用 DeviceKey 加密
      → 更新 settings["ai.apiKey"] = "enc:..."
      → 日志: "API Key migrated to encrypted storage"
```

---

## 六、完整架构总览

```
┌─────────────────────────────────────────────────────────────────┐
│                      Darvin-Resume App                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  internal/crypto/                                        │  │
│  │  ├── aes.go       → Encrypt() / Decrypt() (AES-256-GCM) │  │
│  │  └── device.go    → DeviceKey() (设备唯一标识)            │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  internal/backup/                                        │  │
│  │  ├── backup.go     → CreateBackup() / RestoreBackup()    │  │
│  │  ├── scheduler.go  → Scheduler (time.Ticker 定时备份)    │  │
│  │  └── cleanup.go    → cleanupOldBackups() (旧文件清理)    │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  internal/ai/                                            │  │
│  │  ├── config.go          → 现有配置结构                    │  │
│  │  └── secure_config.go   → SaveSecureAPIKey/LoadSecure    │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  app.go (Bridge Methods)                                 │  │
│  │  ├── CreateManualBackup(password)                        │  │
│  │  ├── RestoreFromBackup(path, password)                   │  │
│  │  └── ListBackups()                                       │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  数据流:                                                        │
│  ┌──────────┐  VACUUM INTO  ┌──────────┐  ZIP+AES  ┌───────┐ │
│  │ data.db  │──────────────>│ temp.db  │──────────>│ .d-b  │ │
│  │ (SQLite) │               │ (快照)   │  加密压缩  │ 备份  │ │
│  └──────────┘               └──────────┘           └───────┘ │
└─────────────────────────────────────────────────────────────────┘
```

---

## 七、新增依赖评估

| 依赖 | 类型 | 大小 | 必要性 |
|------|------|------|--------|
| `crypto/aes` | Go 标准库 | 0 | 已内置 |
| `crypto/cipher` | Go 标准库 | 0 | 已内置 |
| `crypto/rand` | Go 标准库 | 0 | 已内置 |
| `crypto/sha256` | Go 标准库 | 0 | 已内置 |
| `archive/zip` | Go 标准库 | 0 | 已内置 |
| `encoding/hex` | Go 标准库 | 0 | 已内置 |
| `golang.org/x/crypto/argon2` | x 扩展库 | ~50KB | **已在 go.mod** |
| `time.Ticker` | Go 标准库 | 0 | 已内置 |

**结论: 零新增外部依赖。** 所有需要的功能要么是 Go 标准库，要么已在项目间接依赖中。

---

## 八、执行建议

### 8.1 推荐的子任务拆分

| 任务 | 优先级 | 预计复杂度 | 依赖 |
|------|--------|-----------|------|
| T1: `internal/crypto/` 加密模块 | P0 | 中 | 无 |
| T2: API Key 加密存储 + 明文迁移 | P0 | 低 | T1 |
| T3: `internal/backup/` 备份/恢复模块 | P1 | 中 | T1 |
| T4: 定时备份调度器 | P1 | 低 | T3 |
| T5: 前端备份管理 UI | P2 | 中 | T3, T4 |
| T6: 文件对话框 + 进度反馈 | P2 | 低 | T3 |

### 8.2 安全注意事项

1. **不要在日志中打印密钥或密文**
2. **加密后的 API Key 用 `enc:` 前缀标识**，方便迁移检测
3. **备份恢复时要先验证校验和**，防止损坏数据
4. **恢复前保留 pre-restore 备份**，防止恢复失败导致数据丢失
5. **DeviceKey 在 Linux 依赖 `/etc/machine-id`**，重装系统后密钥会变化（这是预期行为）

### 8.3 测试策略

- 加密/解密单元测试 (不同长度的明文、边界情况)
- 设备密钥获取测试 (三平台)
- VACUUM INTO 一致性测试
- 备份/恢复往返测试 (backup → restore → 数据一致)
- 旧备份清理测试
- API Key 加密迁移测试 (明文 → 加密 → 读取)
