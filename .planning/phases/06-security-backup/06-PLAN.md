# Plan 06: 数据安全与备份

**Phase:** 6 of 7
**Goal:** 用户数据受到加密保护，用户可以备份和恢复全部数据，不担心数据丢失
**Depends on:** Phase 5 complete
**Requirements:** AIAI-02, EXPT-10, EXPT-11, EXPT-12

## Success Criteria (what must be TRUE)

1. 用户输入的 API Key 经过 AES-256-GCM 加密后存储，应用重启后可以正确解密使用，旧版明文 API Key 自动迁移
2. 用户可以一键导出全量数据为加密压缩备份包（.darvin-backup），可以选择设置密码保护
3. 用户可以选择备份包一键导入恢复，恢复前自动为当前数据创建备份，恢复后所有数据完整
4. 用户可以设置自动备份周期（禁用/每日/每周/每月），自动备份使用设备密钥加密，最大保留10个备份文件

## Plans

| Plan | Name | Status | Dependencies |
|------|------|--------|-------------|
| 06-01 | 加密基础设施 + API Key 安全存储 | Planned | — |
| 06-02 | 手动备份与恢复 | Planned | 06-01 |
| 06-03 | 自动备份 + 备份管理 UI | Planned | 06-02 |

## Decision Coverage Matrix

| Requirement | Plan | Coverage | Notes |
|-------------|------|----------|-------|
| AIAI-02 | 06-01 | Full | AES-256-GCM + 设备密钥加密 API Key |
| EXPT-10 | 06-02 | Full | VACUUM INTO + ZIP + 可选密码加密 |
| EXPT-11 | 06-02 | Full | 恢复前 pre-restore 备份 + 校验 |
| EXPT-12 | 06-03 | Full | time.Ticker 定时 + 最大数量限制 |

## New Dependencies

### go.mod (backend)
```
无新增依赖 — 使用 Go 标准库 + 已有 golang.org/x/crypto
```

### npm (frontend)
```
无新增依赖
```

## New Files

### Backend (Go)
```
internal/
├── crypto/
│   ├── aes.go         # Encrypt() / Decrypt() — AES-256-GCM 加密解密
│   ├── aes_test.go    # 加密/解密单元测试
│   ├── device.go      # DeviceKey() — 平台设备唯一标识生成
│   └── device_test.go # 设备密钥单元测试
├── backup/
│   ├── backup.go      # CreateBackup() / RestoreBackup() / ListBackups()
│   └── scheduler.go   # Scheduler — time.Ticker 定时备份调度
└── ai/
    └── secure_config.go  # SaveSecureAPIKey() / LoadSecureAPIKey() + 明文迁移
```

### Modified Files
```
internal/ai/client.go    # 改用 secure_config 读写 API Key
internal/ai/config.go    # 更新注释
app.go                   # 新增 backup bridge 方法 + 集成调度器
main.go                  # 调度器生命周期管理
```

### Frontend (Vue/TS)
```
frontend/src/
├── components/
│   └── BackupManager.vue    # 备份管理面板（导出/导入/列表/设置）
└── composables/
    └── useBackup.ts         # 备份操作 composable
```

## Execution Order

```
Wave 1:
  06-01: 加密基础设施 + API Key 安全存储
    - internal/crypto/aes.go — AES-256-GCM 加密/解密
    - internal/crypto/device.go — 设备唯一密钥
    - internal/ai/secure_config.go — 安全存取 + 明文迁移
    - 修改 internal/ai/client.go 使用加密存储

Wave 2 (依赖 06-01 的 crypto 模块):
  06-02: 手动备份与恢复
    - internal/backup/backup.go — 备份创建/恢复/列表
    - app.go 新增 backup bridge 方法
    - 前端 BackupManager.vue + useBackup.ts

Wave 3 (依赖 06-02 的 backup 模块):
  06-03: 自动备份 + 备份管理 UI 完善
    - internal/backup/scheduler.go — 定时调度器
    - app.go/main.go 集成调度器
    - BackupManager.vue 添加自动备份设置
    - settings 表存储备份配置
```

**并行策略：**
- 06-01 是所有其他 plan 的基础（提供加密模块），需最先执行
- 06-02 依赖 06-01 的 crypto 模块做备份文件加密
- 06-03 依赖 06-02 的 backup 模块做定时备份
- 严格串行：06-01 → 06-02 → 06-03

## Wave Analysis

| Wave | Plans | Parallelizable | Key Dependencies |
|------|-------|---------------|-----------------|
| 1 | 06-01 | — | Phase 5 完成 |
| 2 | 06-02 | — | 06-01 (crypto 模块) |
| 3 | 06-03 | — | 06-02 (backup 模块) |
