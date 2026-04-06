# Phase 6 Context — 数据安全与备份

> 生成日期: 2026-04-06

## 现状摘要

### API Key 存储
- **明文存储**在 SQLite `settings` 表（`key="ai.apiKey"`）
- 代码注释已确认: `config.go:24` — "APIKey is stored in plaintext this phase; encryption added in Phase 6"
- 存储链路: 前端 `AIConfigModal.vue` → `useAIConfig.ts` → `ai.ts` → `app.go:SaveAIConfig()` → `ai/client.go:SaveConfig()` → `settings.Set(ctx, "ai.apiKey", cfg.APIKey)`

### 数据库
- SQLite 3（modernc.org/sqlite，纯 Go），4 张表：resumes/settings/ai_messages/snapshots
- WAL 模式，全局 `*sql.DB`（`database.DB`）
- 迁移使用 goose v3

### 架构模式
- Bridge: `app.go` → Service(Interface) → 全局 DB
- 无 repository/store 抽象层
- Service 接口 + 私有结构体 + 工厂函数

### 已有导出
- PDF 导出（系统打印 + Chromedp）
- 快照版本管理（snapshots 表）
- **无数据备份/导入机制**

## 关键约束
- 零新增外部依赖（仅 Go 标准库 + 已有 `golang.org/x/crypto`）
- 三平台支持（Windows/macOS/Linux）
- 设备密钥不存储文件，通过平台特征动态计算

## 相关文件

| 文件 | 角色 |
|------|------|
| `internal/ai/config.go` | AI 配置常量（SettingKeyAPIKey 等） |
| `internal/ai/client.go` | AI 客户端、配置读写（SaveConfig/LoadConfig） |
| `internal/settings/settings.go` | 通用 Key-Value 设置存储 |
| `internal/database/db.go` | 数据库初始化（getUserDataDir） |
| `internal/service/resume.go` | Resume 业务逻辑 |
| `internal/service/snapshot.go` | Snapshot 业务逻辑 |
| `app.go` | Wails Bridge 层 |
| `frontend/src/composables/useAIConfig.ts` | 前端 AI 配置管理 |
| `frontend/src/components/AIConfigModal.vue` | 前端 AI 配置弹窗 |
