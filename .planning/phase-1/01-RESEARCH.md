# Phase 1: 项目骨架与数据层 - Research

**Researched:** 2026-04-05
**Domain:** Wails v2 桌面应用框架 + Go/SQLite 数据层 + Vue 3 前端骨架
**Confidence:** HIGH

## Summary

Phase 1 是整个项目的基础设施搭建阶段，需要完成 Wails v2 项目初始化、SQLite 存储层实现、JSON-Markdown 正向同步三个核心任务。研究覆盖了 Wails v2 项目结构、Go-SQLite 库选型、数据库迁移工具、markdown-it Vue 3 集成方案等关键技术决策。

**核心推荐：** 使用 `wails init` 脚手架创建项目骨架，`modernc.org/sqlite` (纯 Go，无 CGo) 作为 SQLite 驱动，`golang-migrate` 管理数据库迁移，`markdown-it` 直接集成而非通过 Vue 包装库。

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| EDIT-04 | 编辑器预览与PDF导出采用完全一致的渲染引擎与CSS样式表 | markdown-it 统一渲染引擎集成方案，Go 端模板渲染管道 |
| TMPL-07 | JSON→Markdown正向同步自动执行 | JSON→Markdown 转换器设计，结构化数据到 Markdown 模板映射 |
| TMPL-08 | Markdown→JSON反向同步需用户手动触发（本阶段仅建立数据模型基础） | 简历 JSON Schema 设计预留反向同步扩展点 |
</phase_requirements>

## User Constraints (from CONTEXT.md)

### Locked Decisions
- **框架**: Wails v2 (Go + WebView)
- **后端**: Go (文件IO、SQLite事务、业务逻辑)
- **前端**: Vue 3 + TypeScript
- **数据库**: SQLite 3
- **渲染**: Markdown-it（统一渲染引擎）
- **架构分层**: Bridge层 → Domain业务层(ResumeManager, TemplateRender) → Infrastructure层(SQLiteStore) → 前端(Vue 3组件+路由)
- **项目约束**: 冷启动<2s，内存<200MB，安装包<50MB；全量数据本地存储；支持 Windows、macOS、Linux 三平台；预览与PDF导出100%渲染一致

### Claude's Discretion
（本阶段无明确 Claude 自由裁量项，技术实现细节由研究推荐确定）

### Deferred Ideas (OUT OF SCOPE)
- TMPL-08 反向同步完整实现（Phase 3 AI能力接入后）
- Monaco Editor 集成（Phase 2）
- AI 相关功能（Phase 3）
- 模板系统、PDF导出（Phase 5）

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Wails v2 | v2.12.0 | 桌面应用框架 | 项目确定选型，Go+WebView 轻量方案 [VERIFIED: go list -m] |
| modernc.org/sqlite | v1.48.1 | SQLite 驱动 (纯Go) | 无 CGo 依赖，跨平台交叉编译简单 [VERIFIED: go list -m] |
| Vue 3 | 3.5.32 | 前端框架 | 项目确定选型 [VERIFIED: npm view] |
| TypeScript | 6.0.5 | 类型安全 | 项目确定选型 [VERIFIED: npm view] |
| markdown-it | 14.1.1 | Markdown 渲染引擎 | 项目确定选型，统一预览与导出渲染 [VERIFIED: npm view] |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| vue-router | 5.0.4 | 前端路由 | 页面导航（简历列表、编辑器、设置） [VERIFIED: npm view] |
| golang-migrate | v4.19.1 | 数据库迁移 | SQLite Schema 版本管理 [VERIFIED: go list -m] |
| Vite | 8.0.3 | 前端构建工具 | Wails Vue 模板自带 [VERIFIED: npm view] |
| @vitejs/plugin-vue | 6.0.5 | Vue Vite 插件 | Wails Vue 模板自带 [VERIFIED: npm view] |
| goose | v3.27.0 | 数据库迁移 (备选) | 如果 golang-migrate 不满足需求 [VERIFIED: go list -m] |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| modernc.org/sqlite | mattn/go-sqlite3 | mattn 更快、功能更丰富，但依赖 CGo，交叉编译配置复杂，对 Wails 三平台构建增加额外复杂度 |
| golang-migrate | goose | goose 支持 Go 迁移脚本（更强编程能力），golang-migrate 社区更广泛、SQL 迁移文件更直观 |
| @f3ve/vue-markdown-it | 直接使用 markdown-it | vue-markdown-it 是维护中的 Vue 3 包装器，但直接使用 markdown-it + onMounted 更灵活，避免额外依赖 |

**安装命令:**

```bash
# Wails CLI 安装（首次开发环境搭建）
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 项目初始化（自动生成前端模板）
wails init -n open-resume -t vue-ts

# Go 后端依赖
go get modernc.org/sqlite
go get github.com/golang-migrate/migrate/v4

# 前端依赖（在 frontend/ 目录下）
cd frontend
npm install vue-router markdown-it
```

**版本验证记录:**
- Wails v2: v2.12.0 (2026-04-05, `go list -m -versions`)
- modernc.org/sqlite: v1.48.1 (2026-04-03, `go list -m -json`)
- golang-migrate: v4.19.1 (2026-04-05, `go list -m -versions`)
- markdown-it: 14.1.1 (2026-04-05, `npm view`)
- vue-router: 5.0.4 (2026-04-05, `npm view`)
- Vue: 3.5.32 (2026-04-05, `npm view`)
- TypeScript: 6.0.5 (2026-04-05, `npm view`)

## Architecture Patterns

### 推荐项目结构

```
open-resume/
├── app.go                          # Wails App 结构体，生命周期方法
├── main.go                         # 入口，wails.Run 配置
├── wails.json                      # Wails 项目配置
├── go.mod / go.sum
├── build/                          # 构建资源（图标、平台配置）
│   ├── appicon.png
│   ├── darwin/
│   └── windows/
├── internal/                       # Go 内部包（不对外暴露）
│   ├── bridge/                     # Bridge 层 — Wails 绑定结构体
│   │   └── resume_bridge.go        # 简历相关前后端通信路由
│   ├── domain/                     # Domain 业务层
│   │   ├── resume_manager.go       # 简历管理（CRUD 逻辑）
│   │   ├── template_render.go      # JSON→Markdown 模板渲染
│   │   └── models.go               # 领域模型定义
│   └── infrastructure/             # Infrastructure 基础设施层
│       ├── sqlite_store.go         # SQLite 存储实现
│       └── db.go                   # 数据库初始化、迁移
├── migrations/                     # SQL 迁移文件
│   ├── 001_init_schema.up.sql
│   └── 001_init_schema.down.sql
└── frontend/                       # Vue 3 + TypeScript 前端
    ├── index.html
    ├── package.json
    ├── tsconfig.json
    ├── vite.config.ts
    └── src/
        ├── App.vue                 # 根组件（侧边栏+主内容区布局）
        ├── main.ts                 # 入口
        ├── router/
        │   └── index.ts            # Vue Router 路由配置
        ├── views/
        │   ├── ResumeListView.vue  # 简历列表页（占位）
        │   └── EditorView.vue      # 编辑器页（占位）
        ├── components/
        │   ├── AppSidebar.vue      # 侧边栏导航
        │   └── MarkdownPreview.vue # Markdown 预览组件
        ├── composables/            # Vue 组合式函数
        │   └── useWails.ts         # Wails 调用封装
        ├── lib/
        │   └── markdown.ts         # markdown-it 实例与配置
        └── styles/
            └── main.css            # 全局样式
```

### Pattern 1: Wails Binding 模式（Bridge 层）

**What:** Go 结构体的公共方法通过 Wails `Bind` 选项自动暴露给前端，Wails 生成对应的 TypeScript 声明。
**When to use:** 所有需要前后端通信的场景。

**Example:**

```go
// Source: Wails v2 官方文档 https://wails.io/docs/howdoesitwork
// internal/bridge/resume_bridge.go
package bridge

import (
    "context"
    "open-resume/internal/domain"
)

type ResumeBridge struct {
    ctx    context.Context
    manager *domain.ResumeManager
}

func NewResumeBridge(manager *domain.ResumeManager) *ResumeBridge {
    return &ResumeBridge{manager: manager}
}

func (b *ResumeBridge) Startup(ctx context.Context) {
    b.ctx = ctx
}

// Wails 自动将此方法暴露为前端可调用函数
// 返回值支持基本类型、struct、slice、map（需 json tag）
func (b *ResumeBridge) ListResumes() ([]domain.Resume, error) {
    return b.manager.ListResumes(b.ctx)
}

func (b *ResumeBridge) CreateResume(title string) (*domain.Resume, error) {
    return b.manager.CreateResume(b.ctx, title)
}

func (b *ResumeBridge) GetResume(id string) (*domain.Resume, error) {
    return b.manager.GetResume(b.ctx, id)
}

func (b *ResumeBridge) UpdateResumeContent(id string, jsonData string) error {
    return b.manager.UpdateContent(b.ctx, id, jsonData)
}
```

```go
// main.go — 注册绑定
func main() {
    store := infrastructure.NewSQLiteStore("open-resume.db")
    manager := domain.NewResumeManager(store)
    renderer := domain.NewTemplateRender()
    resumeBridge := bridge.NewResumeBridge(manager)

    app := NewApp()

    wails.Run(&options.App{
        Title:  "Open-Resume",
        Width:  1280,
        Height: 800,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        OnStartup:  func(ctx context.Context) {
            app.startup(ctx)
            resumeBridge.Startup(ctx)
        },
        Bind: []interface{}{
            app,
            resumeBridge, // 注册 Bridge
        },
    })
}
```

**前端调用:**

```typescript
// frontend/src/composables/useWails.ts
import { ListResumes, CreateResume, GetResume, UpdateResumeContent }
  from "../../wailsjs/go/bridge/ResumeBridge"

// Wails 自动生成类型声明，调用返回 Promise
export function useResume() {
  async function list() {
    return await ListResumes()
  }

  async function create(title: string) {
    return await CreateResume(title)
  }

  async function get(id: string) {
    return await GetResume(id)
  }

  async function updateContent(id: string, jsonData: string) {
    return await UpdateResumeContent(id, jsonData)
  }

  return { list, create, get, updateContent }
}
```

### Pattern 2: 简历 JSON Schema 设计

**What:** 简历结构化数据作为全系统唯一可信数据源，模块化组织，支持用户自定义模块顺序。
**When to use:** 所有简历数据的读写、JSON-Markdown 转换。

**Example:**

```go
// internal/domain/models.go
package domain

import "time"

// Resume 简历主记录
type Resume struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    TemplateID  string    `json:"templateId"`
    Content     ResumeContent `json:"content"`       // 结构化JSON内容
    Markdown    string    `json:"markdown"`         // 生成的Markdown
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

// ResumeContent 简历结构化数据（唯一数据源）
type ResumeContent struct {
    BasicInfo       *BasicInfoSection       `json:"basicInfo,omitempty"`
    Education       *EducationSection       `json:"education,omitempty"`
    Skills          *SkillsSection          `json:"skills,omitempty"`
    Projects        *ProjectsSection        `json:"projects,omitempty"`
    SelfEvaluation  *SelfEvaluationSection  `json:"selfEvaluation,omitempty"`
    Internship      *InternshipSection      `json:"internship,omitempty"`
    Campus          *CampusSection          `json:"campus,omitempty"`
    Awards          *AwardsSection          `json:"awards,omitempty"`
    Certificates    *CertificatesSection    `json:"certificates,omitempty"`
    ModuleOrder     []string                `json:"moduleOrder"` // 模块显示顺序
}

// BasicInfoSection 基础信息
type BasicInfoSection struct {
    Enabled  bool           `json:"enabled"`
    Name     string         `json:"name"`
    Phone    string         `json:"phone"`
    Email    string         `json:"email"`
    Github   string         `json:"github,omitempty"`
    Website  string         `json:"website,omitempty"`
    Summary  string         `json:"summary,omitempty"`
    Extra    []KeyValueItem `json:"extra,omitempty"` // 自定义键值对
}

// KeyValueItem 通用键值对（支持用户自定义字段）
type KeyValueItem struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}
```

### Pattern 3: JSON→Markdown 正向同步（TemplateRender）

**What:** 结构化 JSON 数据通过 Go 端模板渲染器自动生成 Markdown 文本。
**When to use:** 每次简历 JSON 内容变更后自动触发。

**Example:**

```go
// internal/domain/template_render.go
package domain

import (
    "bytes"
    "fmt"
    "strings"
    "text/template"
)

type TemplateRender struct {
    tpl *template.Template
}

func NewTemplateRender() *TemplateRender {
    tpl := template.Must(template.New("resume").Parse(resumeMarkdownTemplate))
    return &TemplateRender{tpl: tpl}
}

// RenderToMarkdown 将结构化 JSON 转换为 Markdown
func (r *TemplateRender) RenderToMarkdown(content *ResumeContent) (string, error) {
    var buf bytes.Buffer
    if err := r.tpl.Execute(&buf, content); err != nil {
        return "", fmt.Errorf("render markdown: %w", err)
    }
    return strings.TrimSpace(buf.String()), nil
}

const resumeMarkdownTemplate = `{{- if .BasicInfo}}{{- if .BasicInfo.Enabled}}
# {{.BasicInfo.Name}}

{{- if .BasicInfo.Phone}} {{.BasicInfo.Phone}} |{{- end}}
{{- if .BasicInfo.Email}} {{.BasicInfo.Email}} |{{- end}}
{{- if .BasicInfo.Github}} {{.BasicInfo.Github}} |{{- end}}

{{- if .BasicInfo.Summary}}

{{.BasicInfo.Summary}}
{{- end}}

{{- range .BasicInfo.Extra}}
- **{{.Key}}**: {{.Value}}
{{- end}}
{{end}}{{end}}

{{range $idx, $mod := .ModuleOrder}}
{{- if eq $mod "education"}}{{template "education" $.Education}}{{end}}
{{- if eq $mod "skills"}}{{template "skills" $.Skills}}{{end}}
{{- if eq $mod "projects"}}{{template "projects" $.Projects}}{{end}}
{{- if eq $mod "selfEvaluation"}}{{template "selfEval" $.SelfEvaluation}}{{end}}
{{- if eq $mod "internship"}}{{template "internship" $.Internship}}{{end}}
{{- if eq $mod "campus"}}{{template "campus" $.Campus}}{{end}}
{{- if eq $mod "awards"}}{{template "awards" $.Awards}}{{end}}
{{- if eq $mod "certificates"}}{{template "certificates" $.Certificates}}{{end}}
{{- end}}
`
```

### Pattern 4: SQLite 存储层

**What:** 使用 modernc.org/sqlite 纯 Go 驱动实现 SQLite 存储，数据库文件存储在用户数据目录。
**When to use:** 所有数据持久化场景。

**Example:**

```go
// internal/infrastructure/sqlite_store.go
package infrastructure

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"

    _ "modernc.org/sqlite"
)

type SQLiteStore struct {
    db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
    // 确保数据目录存在
    dir := filepath.Dir(dbPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("create db directory: %w", err)
    }

    db, err := sql.Open("sqlite", dbPath) // 注意：modernc.org/sqlite 驱动名是 "sqlite"
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }

    // SQLite 性能优化配置
    pragmas := []string{
        "PRAGMA journal_mode=WAL",       // Write-Ahead Logging，提高并发读性能
        "PRAGMA busy_timeout=5000",      // 忙等超时5秒
        "PRAGMA foreign_keys=ON",        // 启用外键约束
        "PRAGMA synchronous=NORMAL",     // 平衡安全性和性能
    }
    for _, p := range pragmas {
        if _, err := db.Exec(p); err != nil {
            return nil, fmt.Errorf("set pragma %s: %w", p, err)
        }
    }

    return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
    return s.db.Close()
}

func (s *SQLiteStore) DB() *sql.DB {
    return s.db
}
```

**数据库表设计:**

```sql
-- migrations/001_init_schema.up.sql

CREATE TABLE IF NOT EXISTS resumes (
    id          TEXT PRIMARY KEY,        -- UUID
    title       TEXT NOT NULL DEFAULT '',
    template_id TEXT NOT NULL DEFAULT 'minimal',
    content     TEXT NOT NULL '{}',      -- JSON string (ResumeContent)
    markdown    TEXT NOT NULL DEFAULT '',
    is_deleted  INTEGER NOT NULL DEFAULT 0,  -- 软删除标记
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_resumes_updated_at ON resumes(updated_at);
CREATE INDEX idx_resumes_is_deleted ON resumes(is_deleted);

-- 预留：版本快照表（Phase 5 启用）
-- CREATE TABLE IF NOT EXISTS snapshots (
--     id          TEXT PRIMARY KEY,
--     resume_id   TEXT NOT NULL REFERENCES resumes(id),
--     label       TEXT NOT NULL DEFAULT '',
--     content     TEXT NOT NULL,
--     markdown    TEXT NOT NULL,
--     template_id TEXT NOT NULL,
--     created_at  TEXT NOT NULL DEFAULT (datetime('now'))
-- );
```

### Pattern 5: 前端 Markdown 渲染组件

**What:** 直接使用 markdown-it 库在前端渲染 Markdown 预览，与后端 PDF 导出共享同一份 CSS 样式表。
**When to use:** 编辑器预览区域。

**Example:**

```typescript
// frontend/src/lib/markdown.ts
import MarkdownIt from 'markdown-it'

// 创建统一的 markdown-it 实例
// 配置与后端渲染保持一致，确保预览与导出100%渲染一致
const md = new MarkdownIt({
  html: false,         // 不允许HTML标签（安全）
  linkify: true,       // 自动识别链接
  typographer: true,   // 排版优化
  breaks: true,        // 换行符转为 <br>
})

export default md
```

```vue
<!-- frontend/src/components/MarkdownPreview.vue -->
<script setup lang="ts">
import { ref, watch, onMounted, useTemplateRef } from 'vue'
import md from '../lib/markdown'

const props = defineProps<{
  markdown: string
}>()

const previewRef = useTemplateRef<HTMLDivElement>('preview')

watch(() => props.markdown, (newVal) => {
  if (previewRef.value) {
    previewRef.value.innerHTML = md.render(newVal)
  }
}, { immediate: true })
</script>

<template>
  <div ref="preview" class="markdown-preview"></div>
</template>

<style scoped>
.markdown-preview {
  /* 统一 CSS 样式 — 预览与 PDF 导出共用 */
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  font-size: 14px;
  line-height: 1.6;
  color: #333;
  padding: 20px;
}
</style>
```

### Anti-Patterns to Avoid

- **在 Bridge 层写业务逻辑**: Bridge 层应仅做参数传递和类型转换，业务逻辑全部在 Domain 层。违反会导致逻辑散落、难以测试。
- **直接在前端操作数据库**: 所有数据操作必须通过 Wails Binding 经 Go 后端完成。前端不应绕过后端直接访问 SQLite 文件。
- **硬编码 Markdown 模板字符串散落各处**: 模板应集中在 TemplateRender 中管理，支持后续多模板扩展。
- **使用 `mattn/go-sqlite3` 后发现交叉编译问题再切换**: 项目初期就应选择 `modernc.org/sqlite`，避免后期切换成本。
- **JSON Schema 设计不考虑模块顺序**: 简历模块顺序是核心需求（用户自定义），必须在 Schema 中预留 `moduleOrder` 字段。

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| 数据库迁移 | 手动执行 CREATE TABLE SQL | golang-migrate | 版本管理、回滚支持、多环境一致 |
| UUID 生成 | 自定义ID生成器 | `github.com/google/uuid` | RFC 4122 标准，碰撞概率极低 |
| JSON 序列化/反序列化 | 手动字符串拼接 | Go 标准库 `encoding/json` | 类型安全、性能优化、错误处理 |
| Markdown 解析 | 正则表达式或手写解析器 | markdown-it | 边界情况极多（嵌套列表、代码块、转义字符） |
| 跨平台数据目录 | 硬编码路径 | `os.UserConfigDir()` (Go 1.13+) | 各平台数据目录规范不同 |

**Key insight:** Wails 的 Binding 机制已经解决了前后端通信问题，不需要额外的 HTTP server 或 WebSocket 层。Go struct 方法自动变成前端 TypeScript 函数调用。

## Common Pitfalls

### Pitfall 1: Wails 未安装导致初始化失败

**What goes wrong:** 执行 `wails init` 命令时报 "command not found"。
**Why it happens:** Wails CLI 需要单独安装，不属于 Go 标准工具链。当前开发环境确认未安装。
**How to avoid:** 先执行 `go install github.com/wailsapp/wails/v2/cmd/wails@latest`，确认 `wails doctor` 检查通过。
**Warning signs:** `wails` 命令不存在；`wails doctor` 报缺少依赖。

### Pitfall 2: modernc.org/sqlite 驱动名不是 "sqlite3"

**What goes wrong:** `sql.Open("sqlite3", ...)` 会找不到驱动。
**Why it happens:** `mattn/go-sqlite3` 驱动名是 `"sqlite3"`，但 `modernc.org/sqlite` 驱动名是 `"sqlite"`。
**How to avoid:** 使用 `modernc.org/sqlite` 时，`sql.Open("sqlite", dbPath)`。
**Warning signs:** 运行时报 "unknown driver sqlite3" 错误。

### Pitfall 3: Wails Binding 不支持 Go 指针类型返回

**What goes wrong:** Binding 方法返回 `*Resume` 但前端收不到数据或报错。
**Why it happens:** Wails Binding 只支持值类型和特定引用类型（slice、map），指针返回值需要确保结构体字段有 `json` tag。
**How to avoid:** 确保所有返回的 struct 字段都有 `json` tag，避免返回 interface{} 类型。
**Warning signs:** 前端调用返回 undefined 或空对象。

### Pitfall 4: SQLite WAL 模式未配置导致并发写入阻塞

**What goes wrong:** 写入操作频繁超时或报 "database is locked"。
**Why it happens:** SQLite 默认使用 DELETE journal mode，写操作会阻塞所有读操作。
**How to avoid:** 连接时执行 `PRAGMA journal_mode=WAL` 和 `PRAGMA busy_timeout=5000`。
**Warning signs:** 偶发的 "database is locked" 错误。

### Pitfall 5: 前端 Vite 开发服务器与 Wails dev 模式端口不匹配

**What goes wrong:** `wails dev` 启动后前端页面空白或404。
**Why it happens:** `wails.json` 中 `frontend:dev:serverUrl` 配置的端口与 Vite 实际启动端口不一致。
**How to avoid:** `wails init` 生成的默认配置通常是正确的，但修改 Vite 端口时需同步更新 `wails.json`。
**Warning signs:** `wails dev` 控制台显示连接失败，浏览器控制台显示资源加载错误。

### Pitfall 6: 嵌入资源路径不匹配

**What goes wrong:** `//go:embed all:frontend/dist` 找不到资源。
**Why it happens:** 前端未构建或 `frontend/dist` 目录不存在。
**How to avoid:** 开发时使用 `wails dev`（自动处理），构建时确保先执行 `wails build`（内部先构建前端）。
**Warning signs:** 编译时报 "no matching files found" 错误。

## Code Examples

### Wails 项目入口完整模板

```go
// Source: Wails v2 官方文档 https://wails.io/docs/howdoesitwork
// main.go
package main

import (
    "embed"
    "log"
    "os"
    "path/filepath"

    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "open-resume/internal/bridge"
    "open-resume/internal/domain"
    "open-resume/internal/infrastructure"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // 获取平台数据目录
    configDir, err := os.UserConfigDir()
    if err != nil {
        log.Fatal(err)
    }
    dbPath := filepath.Join(configDir, "open-resume", "open-resume.db")

    // 初始化存储层
    store, err := infrastructure.NewSQLiteStore(dbPath)
    if err != nil {
        log.Fatal(err)
    }
    defer store.Close()

    // 初始化业务层
    manager := domain.NewResumeManager(store)
    renderer := domain.NewTemplateRender()

    // 初始化 Bridge 层
    resumeBridge := bridge.NewResumeBridge(manager, renderer)

    // 创建应用
    app := NewApp()

    // 启动 Wails
    err = wails.Run(&options.App{
        Title:     "Open-Resume",
        Width:     1280,
        Height:    800,
        MinWidth:  1200,
        MinHeight: 700,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        OnStartup: func(ctx context.Context) {
            app.startup(ctx)
            resumeBridge.Startup(ctx)
        },
        OnShutdown: func(ctx context.Context) {
            store.Close()
        },
        Bind: []interface{}{
            app,
            resumeBridge,
        },
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

### wails.json 配置模板

```json
{
    "$schema": "https://wails.io/schemas/config.v2.json",
    "name": "open-resume",
    "outputfilename": "open-resume",
    "frontend:install": "npm install",
    "frontend:build": "npm run build",
    "frontend:dev:watcher": "npm run dev",
    "frontend:dev:serverUrl": "auto",
    "author": {
        "name": "Open-Resume"
    }
}
```

### 前端路由配置

```typescript
// frontend/src/router/index.ts
import { createRouter, createWebHashHistory } from 'vue-router'

// 注意：Wails 应用必须使用 Hash 模式
// 因为 WebView 没有真实的 HTTP 服务器处理 HTML5 History 路由
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'resumes',
      component: () => import('../views/ResumeListView.vue'),
    },
    {
      path: '/editor/:id',
      name: 'editor',
      component: () => import('../views/EditorView.vue'),
    },
  ],
})

export default router
```

### 数据库迁移文件

```sql
-- migrations/001_init_schema.up.sql
CREATE TABLE IF NOT EXISTS resumes (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL DEFAULT '',
    template_id TEXT NOT NULL DEFAULT 'minimal',
    content     TEXT NOT NULL DEFAULT '{}',
    markdown    TEXT NOT NULL DEFAULT '',
    is_deleted  INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_resumes_updated_at ON resumes(updated_at);
CREATE INDEX idx_resumes_is_deleted ON resumes(is_deleted);
```

```sql
-- migrations/001_init_schema.down.sql
DROP TABLE IF EXISTS resumes;
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| mattn/go-sqlite3 (CGo) | modernc.org/sqlite (纯Go) | 成熟于2023-2024 | 交叉编译无需 CGo 工具链，构建复杂度大幅降低 |
| Wails v1 | Wails v2 | 2022 | 改进的 API、更好的开发体验、Vite 支持 |
| 手动 SQL 迁移 | golang-migrate / goose | 成熟方案 | 版本化 Schema 管理，支持回滚 |
| vue3-markdown-it (过时) | 直接使用 markdown-it + Vue 3 Composition API | 2024+ | vue3-markdown-it 4年未更新，直接使用更灵活可靠 |
| Electron | Wails (Go + WebView) | 2022+ | 安装包从 150MB+ 降至 <50MB，内存从 500MB+ 降至 <200MB |

**Deprecated/outdated:**
- `vue3-markdown-it` npm 包：约4年未更新，不推荐使用 [VERIFIED: npm registry]
- Wails v1：已被 v2 完全替代 [CITED: wails.io]
- `mattn/go-sqlite3` 在 Wails 场景下不推荐：CGo 交叉编译复杂度高 [ASSUMED]

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | 后端编译 | ✓ | 1.25.3 | — |
| Node.js | 前端构建 | ✓ | 22.19.0 | — |
| npm | 前端依赖管理 | ✓ | 11.6.0 | — |
| Wails CLI | 项目初始化/开发/构建 | ✗ | — | 安装: `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| GCC/CGo | (不要求) | — | — | 使用 modernc.org/sqlite 无需 CGo |

**Missing dependencies with no fallback:**
- Wails CLI — 阻塞性依赖，必须在开发前安装。安装命令: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`，安装后运行 `wails doctor` 检查环境完整性。

**Missing dependencies with fallback:**
- 无。除 Wails CLI 外所有依赖均已就绪。

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go testing (后端) + Vitest (前端) |
| Config file | Go: 内置; 前端: vitest.config.ts (Wave 0 创建) |
| Quick run command (Go) | `go test ./internal/... -v -short` |
| Quick run command (前端) | `cd frontend && npm run test` |
| Full suite command | `go test ./... -v && cd frontend && npm run test -- --run` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| EDIT-04 | markdown-it 渲染 HTML 输出一致性 | unit | `go test ./internal/domain/... -run TestRender -v` | ❌ Wave 0 |
| TMPL-07 | JSON→Markdown 转换正确性 | unit | `go test ./internal/domain/... -run TestTemplateRender -v` | ❌ Wave 0 |
| TMPL-07 | 前端 Markdown 预览渲染 | unit | `cd frontend && npm run test -- MarkdownPreview` | ❌ Wave 0 |
| TMPL-08 | 数据模型支持反向同步扩展 | unit | `go test ./internal/domain/... -run TestResumeContent -v` | ❌ Wave 0 |
| — | SQLite CRUD 操作 | unit | `go test ./internal/infrastructure/... -v` | ❌ Wave 0 |
| — | Wails Binding 调用链 | integration | `go test ./internal/bridge/... -v` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/... -v -short && cd frontend && npm run test`
- **Per wave merge:** `go test ./... -v && cd frontend && npm run test -- --run`
- **Phase gate:** 全部测试绿色通过后再执行 `/gsd-verify-work`

### Wave 0 Gaps
- [ ] `internal/domain/template_render_test.go` — 覆盖 TMPL-07 JSON→Markdown 转换
- [ ] `internal/domain/models_test.go` — 覆盖 TMPL-08 数据模型序列化/反序列化
- [ ] `internal/infrastructure/sqlite_store_test.go` — 覆盖 SQLite CRUD 操作
- [ ] `internal/bridge/resume_bridge_test.go` — 覆盖 Binding 层调用
- [ ] `frontend/vitest.config.ts` — 前端测试框架配置
- [ ] `frontend/src/components/__tests__/MarkdownPreview.spec.ts` — 覆盖 EDIT-04 渲染
- [ ] Vitest 安装: `cd frontend && npm install -D vitest @vue/test-utils` — 前端测试框架

## Security Domain

> 本阶段安全需求较低，主要是基础数据层搭建。完整安全措施在后续阶段实施。

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | no | 桌面应用，无用户认证 |
| V3 Session Management | no | 无网络会话 |
| V4 Access Control | no | 单用户本地应用 |
| V5 Input Validation | yes | Go 端 JSON 序列化验证，SQL 参数化查询 |
| V6 Cryptography | no | Phase 6 实现 API Key 加密 |

### Known Threat Patterns for Wails + SQLite

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| SQL 注入 | Tampering | 使用参数化查询 (`$1`, `$2`)，禁止字符串拼接 SQL |
| JSON 注入 | Tampering | Go 端 `encoding/json` 严格序列化/反序列化 |
| 路径遍历 | Tampering | 使用 `filepath.Join` 构建路径，不信任用户输入 |
| WebView XSS | Tampering | markdown-it 设置 `html: false`，禁止渲染原始 HTML |

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | `modernc.org/sqlite` 在 Wails 三平台交叉编译中无需额外 CGo 配置 | Standard Stack | 需回退到 mattn/go-sqlite3 并配置 CGo |
| A2 | Wails v2 v2.12.0 与 Go 1.25.3 兼容 | Environment Availability | 需降级 Go 版本或升级 Wails |
| A3 | `wails init -t vue-ts` 模板生成的 Vue 3 项目使用 Vite 8.x | Architecture Patterns | 可能需要手动升级 Vite 配置 |
| A4 | Vue Router 5.x 的 Hash 模式在 Wails WebView 中正常工作 | Code Examples | 需寻找替代路由方案 |

## Open Questions

1. **Wails v2 与 Go 1.25.3 兼容性**
   - What we know: Wails v2 v2.12.0 是最新稳定版，Go 1.25.3 已安装
   - What's unclear: Wails v2 v2.12.0 是否已适配 Go 1.25（Go 版本较新）
   - Recommendation: 安装 Wails CLI 后运行 `wails doctor` 验证，如果报版本不兼容则按提示处理

2. **数据库文件存储位置**
   - What we know: 应使用 `os.UserConfigDir()` 获取平台数据目录
   - What's unclear: 是否需要提供自定义存储路径选项
   - Recommendation: Phase 1 先使用默认路径（`~/.config/open-resume/` on Linux, `%APPDATA%/open-resume/` on Windows），后续 Phase 可增加自定义路径

3. **markdown-it 插件策略**
   - What we know: 基础 markdown-it 满足大部分 Markdown 渲染需求
   - What's unclear: 是否需要 markdown-it 插件（如脚注、数学公式等）
   - Recommendation: Phase 1 使用基础配置，按需在后续 Phase 添加插件

## Sources

### Primary (HIGH confidence)
- Wails v2 官方文档 https://wails.io/docs/introduction — 项目结构、Binding 机制、Options API
- `go list -m -versions github.com/wailsapp/wails/v2` — Wails v2.12.0 最新版本确认
- `go list -m -json modernc.org/sqlite@latest` — modernc.org/sqlite v1.48.1 版本确认
- `npm view markdown-it version` — markdown-it 14.1.1 版本确认
- `npm view vue-router version` — vue-router 5.0.4 版本确认

### Secondary (MEDIUM confidence)
- `npm view @f3ve/vue-markdown-it` — Vue 3 markdown-it 包装器信息
- golang-migrate GitHub — 迁移工具用法参考
- modernc.org/sqlite 与 mattn/go-sqlite3 对比 — 多来源交叉验证

### Tertiary (LOW confidence)
- Wails v2 与 Go 1.25.3 兼容性 — 未在实际环境中验证，需 `wails doctor` 确认 [ASSUMED]

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — 所有版本通过 registry 命令验证
- Architecture: HIGH — 基于 Wails 官方文档和 Go 标准实践
- Pitfalls: HIGH — 基于 Wails 文档、SQLite 官方建议、实际开发经验
- Security: MEDIUM — 本阶段安全需求简单，但 WebView 安全最佳实践需持续关注

**Research date:** 2026-04-05
**Valid until:** 2026-05-05 (30 days, 稳定技术栈)
