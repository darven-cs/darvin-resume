# Phase 1 执行计划: 项目骨架与数据层

**Phase**: 1 of 7 — 项目骨架与数据层
**Goal**: 项目可以启动运行，简历数据可以持久化存储和检索，为所有后续功能提供数据基础
**Created**: 2026-04-05

## 研究摘要

### 环境信息
- Go: 1.25.3 linux/amd64
- Node: v22.19.0, npm: 11.6.0
- **Wails CLI 未安装**（需在 01-01 中安装）

### 关键依赖版本
| 依赖 | 版本 | 用途 |
|------|------|------|
| Wails v2 | v2.12.0 | 桌面框架 |
| modernc.org/sqlite | v1.48.1 | 纯Go SQLite驱动（无需CGO） |
| goose | v3.27.0 | 数据库迁移工具 |
| Vue | 3.5.32 | 前端框架 |
| Vue Router | 5.0.4 | 前端路由 |
| Vite | 8.0.3 | 构建工具 |
| markdown-it | 14.1.1 | Markdown渲染引擎 |
| TypeScript | 6.0.5 | 类型系统 |

### 关键技术决策
1. **使用 modernc.org/sqlite 而非 mattn/go-sqlite3**: 纯Go实现，无需CGO，简化跨平台构建
2. **使用 goose 做数据库迁移**: 支持Go和SQL两种迁移格式，社区活跃
3. **Wails Bindings 自动生成**: Go公开方法自动暴露给前端TypeScript，无需手动绑定

---

## Plan 01-01: Wails v2 项目初始化与项目骨架

**目标**: Wails 应用可以正常启动并显示空白页面

### 前置条件
- [ ] 安装 Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- [ ] 验证 Wails 依赖: `wails doctor`

### 任务

#### T1: 初始化 Wails 项目
- 在项目根目录下使用 `wails init` 初始化（需先备份现有文件）
- 使用模板: `vue-ts`（Vue 3 + TypeScript）
- 项目名: `Darvin-Resume`
- 由于项目目录已存在（有 .planning/、CLAUDE.md 等），需要：
  1. 在临时目录创建 Wails 项目
  2. 将生成的文件合并到当前目录
  3. 保留已有的 .planning/、CLAUDE.md、.gitignore、docs/

#### T2: 配置 wails.json
```json5
{
  "name": "Darvin-Resume",
  "outputfilename": "Darvin-Resume",
  "frontend:dir": "frontend",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "wailsjsdir": "frontend/src/wailsjs",
  "assetdir": "frontend/dist",
  "info": {
    "companyName": "Darvin-Resume",
    "productName": "Darvin-Resume",
    "productVersion": "0.1.0",
    "comments": "Markdown原生、隐私优先的本地化简历工具"
  }
}
```

#### T3: Go 后端目录结构
```
.
├── main.go              # 入口，Wails配置
├── app.go               # App结构体，生命周期方法
├── internal/
│   ├── database/
│   │   ├── db.go        # 数据库连接管理
│   │   └── migrations/  # SQL迁移文件
│   ├── model/
│   │   └── resume.go    # 简历数据模型
│   ├── service/
│   │   └── resume.go    # 简历业务逻辑
│   └── converter/
│       └── json2md.go   # JSON→Markdown转换
├── wails.json
└── frontend/            # Vue 3前端
```

#### T4: 前端目录结构
```
frontend/
├── src/
│   ├── App.vue
│   ├── main.ts
│   ├── router/
│   │   └── index.ts     # Vue Router 配置
│   ├── views/
│   │   └── HomeView.vue # 主页
│   ├── components/       # 公共组件
│   ├── composables/      # 组合式函数
│   ├── stores/           # 状态管理
│   ├── types/            # TypeScript 类型定义
│   │   └── resume.ts    # 简历类型
│   └── wailsjs/          # Wails自动生成
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

#### T5: 基础 Vue Router 配置
- 路由: `/` → HomeView（空白首页）
- 路由: `/editor/:id` → EditorView（编辑器占位页）
- 创建路由实例并挂载到 App.vue

#### T6: 基础布局框架（侧边栏+主内容区）
在 App.vue 中实现双栏布局：
```vue
<template>
  <div class="app-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2>Darvin-Resume</h2>
      </div>
      <nav class="sidebar-nav">
        <RouterLink to="/">简历列表</RouterLink>
      </nav>
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>
```
- 侧边栏固定宽度 240px，深色背景
- 主内容区自适应剩余宽度
- 基础 CSS 样式（无第三方UI库）

#### T7: Markdown-it 渲染引擎初始化（EDIT-04 基础设施）
安装并配置 markdown-it 作为统一渲染引擎：
- 安装: `npm install markdown-it` (v14.1.1)
- 安装类型: `npm install -D @types/markdown-it`
- 创建 `frontend/src/utils/markdown.ts`:
```typescript
import MarkdownIt from 'markdown-it'

// 统一渲染引擎实例 — 预览和导出共用
export const md = new MarkdownIt({
  html: false,        // 安全：不解析HTML标签
  breaks: true,       // 换行符转 <br>
  linkify: true,      // 自动识别URL
  typographer: true,  // 优化排版
})

// 导出渲染方法供后续Phase使用
export function renderMarkdown(content: string): string {
  return md.render(content)
}
```
> **EDIT-04 关联**: 此渲染引擎实例将在 Phase 2 中同时用于预览区和 PDF 导出，
> 确保 100% 渲染一致性。本阶段仅完成初始化和基础配置。

#### T8: 验证构建
- `wails dev` 能正常启动开发服务器
- 窗口显示侧边栏+主内容区布局无报错
- `wails build` 能成功构建

### 验收标准
- [ ] `wails dev` 启动后显示应用窗口，可见侧边栏和主内容区
- [ ] 前端和后端通信正常（前端可调用后端 Greet 方法测试）
- [ ] markdown-it 已安装并可在前端 import 使用
- [ ] `wails build` 成功生成可执行文件

---

## Plan 01-02: SQLite 存储层与简历 JSON Schema

**目标**: SQLite 数据库正确初始化，简历结构化 JSON Schema 定义完整，CRUD 操作可用

### 依赖
- Plan 01-01 完成（项目骨架存在）

### 任务

#### T1: 数据库连接管理 (internal/database/db.go)
```go
// 功能：
// - 初始化 SQLite 连接（使用 modernc.org/sqlite）
// - 数据库文件路径：用户数据目录/Darvin-Resume/data.db
// - 启用 WAL 模式、外键约束
// - 连接池配置
// - 优雅关闭
```

#### T2: 数据库迁移 (internal/database/migrations/)
使用 goose 管理迁移：

**迁移 001: 创建 resumes 表**
```sql
CREATE TABLE resumes (
    id TEXT PRIMARY KEY,           -- UUID
    title TEXT NOT NULL DEFAULT '',
    json_data TEXT NOT NULL,       -- 简历结构化JSON
    markdown_content TEXT NOT NULL DEFAULT '', -- 生成的Markdown
    template_id TEXT NOT NULL DEFAULT 'default',
    custom_css TEXT NOT NULL DEFAULT '',
    module_order TEXT NOT NULL DEFAULT '[]',  -- 模块顺序JSON数组
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at DATETIME
);
```

**迁移 002: 创建 settings 表**
```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### T3: 简历结构化 JSON Schema 定义 (internal/model/resume.go)

```go
// Resume 结构体定义（对应 JSON Schema）
type Resume struct {
    ID           string       `json:"id"`
    Title        string       `json:"title"`
    BasicInfo    BasicInfo    `json:"basicInfo"`
    Modules      []Module     `json:"modules"`      // 按用户定义顺序
    TemplateID   string       `json:"templateId"`
    CustomCSS    string       `json:"customCss"`
    CreatedAt    time.Time    `json:"createdAt"`
    UpdatedAt    time.Time    `json:"updatedAt"`
}

type BasicInfo struct {
    Name        string `json:"name"`
    Phone       string `json:"phone"`
    Email       string `json:"email"`
    Avatar      string `json:"avatar"`       // 头像URL/base64
    Website     string `json:"website"`
    GitHub      string `json:"github"`
    Address     string `json:"address"`
    Summary     string `json:"summary"`      // 一句话介绍
}

// Module 定义通用模块结构
type Module struct {
    Type     string          `json:"type"`     // education/skills/projects/internship/campus/awards/certificates/evaluation
    Title    string          `json:"title"`    // 模块标题（如"教育经历"）
    Order    int             `json:"order"`    // 排序权重
    Items    json.RawMessage `json:"items"`    // 模块内容（不同类型结构不同）
    Visible  bool            `json:"visible"`  // 是否可见
}

// EducationItem 教育经历
type EducationItem struct {
    School      string `json:"school"`
    Degree      string `json:"degree"`       // 本科/硕士/博士/大专
    Major       string `json:"major"`
    StartDate   string `json:"startDate"`
    EndDate     string `json:"endDate"`
    GPA         string `json:"gpa"`
    Highlights  []string `json:"highlights"` // 核心3-4条
}

// SkillItem 专业技能
type SkillItem struct {
    Category    string   `json:"category"`    // 分类名称
    Skills      []string `json:"skills"`      // 技能列表
}

// ProjectItem 项目经历
type ProjectItem struct {
    Name        string   `json:"name"`
    Role        string   `json:"role"`
    StartDate   string   `json:"startDate"`
    EndDate     string   `json:"endDate"`
    Description string   `json:"description"`
    Highlights  []string `json:"highlights"`
    TechStack   []string `json:"techStack"`
}

// InternshipItem 实习经历
type InternshipItem struct {
    Company     string   `json:"company"`
    Position    string   `json:"position"`
    StartDate   string   `json:"startDate"`
    EndDate     string   `json:"endDate"`
    Description string   `json:"description"`
    Highlights  []string `json:"highlights"`
}

// CampusItem 校园经历
type CampusItem struct {
    Name        string   `json:"name"`
    Role        string   `json:"role"`
    StartDate   string   `json:"startDate"`
    EndDate     string   `json:"endDate"`
    Description string   `json:"description"`
    Highlights  []string `json:"highlights"`
}

// AwardItem 获奖经历
type AwardItem struct {
    Name        string `json:"name"`
    Level       string `json:"level"`       // 国家级/省级/校级
    Date        string `json:"date"`
    Description string `json:"description"`
}

// CertificateItem 证书
type CertificateItem struct {
    Name        string `json:"name"`
    Issuer      string `json:"issuer"`
    Date        string `json:"date"`
    Score       string `json:"score"`
}

// EvaluationItem 自我评价
type EvaluationItem struct {
    Content     string `json:"content"`
}
```

#### T4: CRUD 操作 (internal/service/resume.go)
```go
// ResumeService 接口
type ResumeService interface {
    Create(ctx context.Context, resume *Resume) error
    GetByID(ctx context.Context, id string) (*Resume, error)
    List(ctx context.Context) ([]*ResumeListItem, error)
    Update(ctx context.Context, resume *Resume) error
    Delete(ctx context.Context, id string) error  // 软删除
    UpdateJSON(ctx context.Context, id string, jsonData string) error
}
```

实现要点：
- Create: 生成UUID、初始化空JSON、触发JSON→Markdown同步
- GetByID: 查询未删除的简历，解析JSON
- List: 仅返回id/title/updatedAt（列表页用）
- Update: 更新JSON字段、自动触发JSON→Markdown同步
- Delete: 软删除（设置is_deleted=true、deleted_at）
- UpdateJSON: 单独更新JSON数据并重新生成Markdown

#### T5: 应用启动时初始化数据库
- 在 `app.startup()` 中调用 database.Init()
- 自动运行数据库迁移

#### T6: 单元测试 — CRUD 操作 (internal/service/resume_test.go)
测试文件: `internal/service/resume_test.go`
运行命令: `go test ./internal/service/ -v`

测试用例：
```go
// TestCreateResume — 创建简历，验证UUID非空、字段初始化正确
// TestGetResume — 创建后按ID查询，验证JSON解析完整
// TestListResumes — 创建3条简历后列表返回3条，按更新时间降序
// TestUpdateResume — 更新JSON数据，验证markdown_content同步更新
// TestDeleteResume — 软删除后GetByID返回error，List不包含已删除项
// TestUpdateJSON — 单独更新JSON，验证Markdown自动重新生成
```

#### T7: 单元测试 — JSON→Markdown 转换 (internal/converter/json2md_test.go)
测试文件: `internal/converter/json2md_test.go`
运行命令: `go test ./internal/converter/ -v`

测试用例（覆盖所有模块类型）：
```go
// TestConvertBasicInfo — 验证姓名作标题、联系方式一行、分隔线存在
// TestConvertEducation — 验证学校/专业/学位/时间格式、亮点列表
// TestConvertSkills — 验证分类加粗、技能顿号分隔
// TestConvertProjects — 验证项目名/角色/时间、亮点列表、技术栈标签
// TestConvertInternship — 验证公司/职位/时间格式
// TestConvertCampus — 验证校园经历格式
// TestConvertAwards — 验证奖项/级别/时间
// TestConvertCertificates — 验证证书/颁发机构/时间
// TestConvertEvaluation — 验证自我评价纯文本
// TestModuleOrdering — 验证模块按 order 字段排序输出
// TestEmptyModule — 验证 Visible=false 的模块不出现在输出中
// TestFullResume — 完整简历（含所有模块类型）的端到端转换
```

验证命令: `go test ./internal/... -v -count=1`

### 验收标准
- [ ] 应用启动时数据库文件自动创建
- [ ] migrations 自动执行，表结构正确
- [ ] `go test ./internal/service/ -v` 全部通过（6个测试用例）
- [ ] `go test ./internal/converter/ -v` 全部通过（12个测试用例）
- [ ] 软删除正确标记 deleted_at

---

## Plan 01-03: Bridge 层绑定、JSON-Markdown 正向同步、前端框架

**目标**: 前端可调用后端 CRUD 操作，JSON→Markdown 转换正确，前端页面框架就绪

### 依赖
- Plan 01-01 完成（项目骨架）
- Plan 01-02 完成（存储层）

### 任务

#### T1: Bridge 层绑定 (app.go)
```go
// App 结构体扩展，暴露给前端的方法
func (a *App) CreateResume(title string) (*model.Resume, error)
func (a *App) GetResume(id string) (*model.Resume, error)
func (a *App) ListResumes() ([]*model.ResumeListItem, error)
func (a *App) UpdateResume(id string, jsonData string) error
func (a *App) DeleteResume(id string) error
func (a *App) UpdateResumeModule(id string, moduleType string, moduleData string) error
```

Wails Bind 配置：
```go
// main.go
err := wails.Run(&options.App{
    Title:     "Darvin-Resume",
    Width:     1280,
    Height:    800,
    MinWidth:  1200,
    MinHeight: 700,
    AssetServer: &assetserver.Options{
        Assets: assets,
    },
    OnStartup:  app.startup,
    OnShutdown: app.shutdown,
    Bind: []interface{}{
        app,
    },
})
```

#### T2: JSON→Markdown 正向同步 (internal/converter/json2md.go)

转换规则：
```
简历结构 → Markdown 输出格式：

# {姓名}
{一句话介绍} | {电话} | {邮箱} | {GitHub}

---

## 教育经历
### {学校} — {专业} ({学位})
{开始日期} - {结束日期}
- {亮点1}
- {亮点2}

## 专业技能
- **{分类1}**: {技能1}、{技能2}、{技能3}
- **{分类2}**: {技能1}、{技能2}

## 项目经历
### {项目名} — {角色}
{开始日期} - {结束日期}
{描述}
- {亮点1}
- {亮点2}
**技术栈**: {技术1}、{技术2}

## 实习经历
### {公司} — {职位}
{开始日期} - {结束日期}
{描述}
- {亮点1}

## 自我评价
{内容}
```

实现要点：
- 按 module_order 排序后依次转换
- 每种模块类型对应一个转换函数
- 基础信息部分单独处理（姓名作为标题，联系方式一行）
- 分隔线 `---` 仅在基础信息和第一个模块之间

#### T3: 前端类型定义 (frontend/src/types/resume.ts)
与 Go 后端模型对应的 TypeScript 类型定义，确保前后端类型一致。

#### T4: 前端视图页面
- **HomeView.vue**: 空白首页，展示"Darvin-Resume"标题，含"新建简历"按钮占位
- **EditorView.vue**: 空白编辑器页面占位（Phase 2 完善）
- 路由已在 01-01 T5 配置，此处仅创建视图组件

#### T5: 集成测试 — Bridge 通信验证 (frontend/src/__tests__/bridge.test.ts)
测试文件: `frontend/src/__tests__/bridge.test.ts`（在 dev 模式下手动执行）
运行方式: 在 `wails dev` 启动的浏览器控制台中执行验证脚本

但更可靠的方式是通过后端集成测试验证 Bridge 链路：
测试文件: `internal/service/resume_integration_test.go`
运行命令: `go test ./internal/service/ -v -run Integration`

```go
// TestIntegration_BridgeCRUD — 验证完整 CRUD 链路
// 1. CreateResume → 验证返回值含 UUID、title 匹配
// 2. GetResume(id) → 验证 JSON 解析、Markdown 非空
// 3. ListResumes → 验证列表含刚创建的简历
// 4. UpdateResume(id, newJSON) → 验证 Markdown 内容同步变更
// 5. DeleteResume(id) → 验证后续 GetByID 返回 err
// 6. 重启模拟（关闭DB重连）→ GetResume(id) 仍可查到数据
```

#### T6: 前端冒烟测试
验证流程（手动执行）：
1. `wails dev` 启动应用
2. 浏览器控制台执行: `window['go']['main']['App']['CreateResume']('测试')`
3. 验证返回结果非空、前端无报错
4. 点击侧边栏"简历列表"链接，确认路由跳转正常
5. 浏览器地址栏输入 `/editor/test-id`，确认编辑器页面渲染

### 验收标准
- [ ] `go test ./internal/service/ -v -run Integration` 全部通过
- [ ] 前端控制台可调用 Bridge 方法并获得返回值
- [ ] JSON→Markdown 转换输出格式正确（由 converter 测试保证，非手动对比）
- [ ] 前端路由 `/` 和 `/editor/:id` 均可正常访问
- [ ] 应用窗口标题显示 "Darvin-Resume"

---

## 执行顺序与依赖

```
01-01 ──→ 01-02 ──→ 01-03
(骨架)    (数据层)   (Bridge+同步+前端)
```

三个 Plan 严格串行，每个 Plan 完成后验证再进入下一个。

## Phase 验收标准（必须全部通过）

1. ✅ Wails 应用可以正常启动，显示基础页面框架（侧边栏+主内容区）
2. ✅ SQLite 数据库正确初始化，简历结构化 JSON Schema 完整定义并可读写
3. ✅ 新建一条简历记录后关闭应用重新打开，数据完整保留
4. ✅ JSON 到 Markdown 的正向同步可以正确执行，修改 JSON 字段后生成的 Markdown 内容准确

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| Wails CLI 安装失败（网络问题） | 阻塞 01-01 | 手动下载预编译二进制 / 使用 GOPROXY |
| modernc.org/sqlite 性能不足 | 影响大量数据操作 | Phase 1 数据量小，后续可切换 mattn |
| Vue 3 + Vite 8 兼容性 | 前端构建问题 | Wails 模板自带版本，不升级 |
| JSON Schema 设计不合理 | 影响后续所有Phase | 严格参照需求文档，预留扩展字段 |

---
*Plan created: 2026-04-05*
