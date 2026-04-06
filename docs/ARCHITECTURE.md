# Darvin-Resume 架构与编码指南

> 技术栈：Wails v2 + Go + Vue 3 + TypeScript + Monaco Editor + SQLite
> 最后更新：2026-04-06

---

## 1. 项目总览

Darvin-Resume 是一款本地化简历制作工具，核心特点：

- **Wails v2** 桌面应用框架 —— Go 后端 + WebView 前端
- **100% 本地存储** —— SQLite 数据库，无云端依赖
- **预览与 PDF 导出 100% 渲染一致** —— 共用同一份 CSS (`editor.css`)
- **AI 辅助** —— 用户自备 API Key，后端调用 Claude API

```
┌─────────────────────────────────────────────────┐
│                  Wails v2 框架                   │
│  ┌──────────────────┐  ┌──────────────────────┐  │
│  │   Go 后端         │  │   Vue 3 前端          │  │
│  │   app.go (桥接层) │◄►│   wailsjs/ (绑定)    │  │
│  │   internal/       │  │   components/        │  │
│  │     ai/           │  │   composables/       │  │
│  │     database/     │  │   services/          │  │
│  │     service/      │  │   views/             │  │
│  │     model/        │  │   styles/            │  │
│  │     export/       │  │   utils/             │  │
│  │     settings/     │  │                      │  │
│  └──────────────────┘  └──────────────────────┘  │
│          │                       │                │
│     SQLite DB              Monaco Editor         │
│     (本地文件)            (Markdown编辑)          │
└─────────────────────────────────────────────────┘
```

---

## 2. 目录结构

```
darvin-resume/
├── main.go                    # Go 入口：Wails 应用配置、窗口初始化
├── app.go                     # Go 桥接层：所有前端可调用的方法
├── go.mod / go.sum            # Go 依赖管理
├── wails.json                 # Wails 项目配置
│
├── internal/                  # Go 后端核心（不对外暴露）
│   ├── ai/                    # Claude API 调用
│   │   ├── client.go          #   HTTP 客户端：Chat、ChatStream、StreamEvents
│   │   ├── chat.go            #   聊天历史 CRUD（ai_messages 表）
│   │   ├── config.go          #   AI 配置管理（APIKey、Model、BaseURL）
│   │   └── error.go           #   错误类型定义
│   ├── converter/
│   │   └── json2md.go         #   JSON 简历数据 → Markdown 转换
│   ├── database/
│   │   ├── db.go              #   SQLite 连接初始化 + 迁移执行
│   │   └── migrations/        #   SQL 迁移脚本（001~004）
│   ├── export/
│   │   └── chromedp.go        #   Chromedp 无头浏览器 PDF 导出
│   ├── model/
│   │   ├── resume.go          #   Resume 数据模型 + CRUD 请求/响应类型
│   │   └── snapshot.go        #   Snapshot 快照模型
│   ├── service/
│   │   ├── resume.go          #   Resume 业务逻辑（CRUD + JSON 更新）
│   │   └── snapshot.go        #   Snapshot 业务逻辑（创建/回滚/对比）
│   └── settings/
│       └── settings.go        #   键值对设置存储（settings 表）
│
├── frontend/
│   ├── index.html             # SPA 入口
│   ├── vite.config.ts         # Vite 配置
│   ├── package.json           # npm 依赖
│   └── src/
│       ├── main.ts            # Vue 应用入口
│       ├── App.vue            # 根组件（路由出口）
│       ├── style.css          # 全局基础样式
│       │
│       ├── views/             # 页面组件
│       │   ├── HomeView.vue   #   简历列表页（卡片网格 + 回收站）
│       │   └── EditorView.vue #   编辑器页（工具栏 + Monaco + 预览 + AI）
│       │
│       ├── components/        # UI 组件
│       │   ├── A4Page.vue           # A4 预览容器（Markdown → HTML 渲染）
│       │   ├── MonacoEditor.vue     # 代码编辑器（AI 工具栏 + 右键菜单）
│       │   ├── SplitPane.vue        # 分栏面板
│       │   ├── ExportDialog.vue     # PDF 导出对话框
│       │   ├── StyleEditor.vue      # 样式调整面板（字号/行高/边距/色彩）
│       │   ├── AIChatSidebar.vue    # AI 聊天侧边栏（多轮对话）
│       │   ├── AIFloatingToolbar.vue# AI 浮动工具栏（润色/翻译/缩写/重写）
│       │   ├── AIDiffView.vue       # AI Diff 对比视图
│       │   ├── AIConfigModal.vue    # AI 配置弹窗（API Key 设置）
│       │   ├── AIErrorToast.vue     # AI 错误提示
│       │   ├── ContextMenu.vue      # 右键上下文菜单
│       │   ├── CreateModeModal.vue  # 创建模式选择（空白/向导）
│       │   ├── ResumeCard.vue       # 简历卡片
│       │   ├── ResumeParserModal.vue# 简历解析弹窗
│       │   ├── ResumeWizardSidebar.vue# 简历向导侧边栏
│       │   ├── SnapshotSidebar.vue  # 版本快照侧边栏
│       │   ├── TemplateSelector.vue # 模板选择器
│       │   ├── JobTargetChip.vue    # 目标岗位标签
│       │   ├── RecycleBinSection.vue# 回收站折叠区
│       │   ├── SaveStatusIndicator.vue# 保存状态指示器
│       │   ├── TokenLimitModal.vue  # Token 限制弹窗
│       │   └── wizard/             # 向导子组件
│       │       ├── WizardStepForm.vue
│       │       ├── WizardModuleSelect.vue
│       │       └── WizardGenerate.vue
│       │
│       ├── composables/       # Vue 3 组合函数（可复用逻辑）
│       │   ├── useAutoSave.ts      # 自动保存（防抖 + 后端同步）
│       │   ├── useAIConfig.ts      # AI 配置加载/保存
│       │   ├── useAIError.ts       # AI 错误处理
│       │   ├── useAISelection.ts   # Monaco 选区管理 + AI 操作触发
│       │   ├── useAIStream.ts      # Wails 事件流式响应
│       │   ├── useDebounce.ts      # 防抖工具
│       │   ├── useResumeList.ts    # 简历列表状态管理
│       │   ├── useSnapshot.ts      # 快照管理
│       │   └── useTemplate.ts      # 模板和自定义 CSS 管理
│       │
│       ├── services/          # API 调用封装
│       │   └── ai.ts               # AI 服务（调用 Wails 绑定方法）
│       │
│       ├── types/             # TypeScript 类型定义
│       │   ├── ai.ts               # AI 相关类型（ChatMessage、AIStreamChunk 等）
│       │   ├── resume.ts           # Resume 数据结构
│       │   └── snapshot.ts         # Snapshot 数据结构
│       │
│       ├── utils/             # 工具函数
│       │   ├── markdown.ts         # markdown-it 渲染引擎配置
│       │   ├── resume.ts           # 简历数据工具函数
│       │   └── sanitizeCSS.ts      # CSS 白名单校验 + CSS 变量构建
│       │
│       ├── router/
│       │   └── index.ts            # 路由配置（/ → HomeView, /editor/:id → EditorView）
│       │
│       ├── styles/            # 样式文件
│       │   ├── editor.css          # ★ 预览和 PDF 共用样式表（核心）
│       │   └── templates/          # 简历模板 CSS
│       │       ├── template-minimal.css
│       │       ├── template-campus.css
│       │       ├── template-academic.css
│       │       └── template-dual-col.css
│       │
│       └── wailsjs/           # Wails 自动生成的绑定（勿手动修改）
│           ├── go/main/App.js      # Go 方法 JS 绑定
│           ├── go/main/App.d.ts    # Go 方法 TS 类型声明
│           ├── wailsjs/go/models.ts# Go 结构体 TS 映射
│           └── runtime/runtime.js  # Wails 运行时（EventsOn/EventsEmit）
│
├── wailsjs/                   # 项目根级 Wails 绑定（构建时生成）
├── docs/                      # 项目文档
├── build/                     # 构建输出（.gitignore）
├── darwin/                    # macOS 平台配置
└── windows/                   # Windows 平台配置
```

---

## 3. 架构分层

### 3.1 后端分层（Go）

```
┌──────────────────────────────────────┐
│  app.go — Wails 桥接层               │  ← 前端直接调用这一层
│  (AISendMessage, CreateResume, ...)  │
├──────────────────────────────────────┤
│  service/ — 业务逻辑层               │  ← 不感知 Wails
│  (ResumeService, SnapshotService)    │
├──────────────────────────────────────┤
│  model/ — 数据模型                   │  ← 纯数据结构
│  (Resume, Snapshot, Request/Response)│
├──────────────────────────────────────┤
│  database/ — 数据访问层              │  ← SQLite 连接 + 迁移
│  (db.go, migrations/)                │
├──────────────────────────────────────┤
│  settings/ — 键值存储                │  ← 通用配置（AI Key 等）
│  ai/ — AI 客户端                     │  ← Claude API 调用
│  export/ — PDF 导出                  │  ← Chromedp 无头浏览器
│  converter/ — 格式转换               │  ← JSON → Markdown
└──────────────────────────────────────┘
```

### 3.2 前端分层（Vue 3）

```
┌──────────────────────────────────────┐
│  views/ — 页面                       │  ← 路由级别组件
│  (HomeView, EditorView)              │
├──────────────────────────────────────┤
│  components/ — UI 组件               │  ← 可复用的视图组件
│  (A4Page, MonacoEditor, ...)         │
├──────────────────────────────────────┤
│  composables/ — 组合函数             │  ← 可复用的业务逻辑
│  (useAutoSave, useAIStream, ...)     │
├──────────────────────────────────────┤
│  services/ — 服务层                  │  ← 调用 Wails 绑定
│  (ai.ts → AISendMessage, ...)        │
├──────────────────────────────────────┤
│  wailsjs/ — Wails 绑定               │  ← 自动生成，勿手动改
│  (App.js, App.d.ts, models.ts)       │
└──────────────────────────────────────┘
```

---

## 4. 核心数据流

### 4.1 编辑 → 预览 → 保存

```
MonacoEditor (用户编辑 Markdown)
    │ @input → emit('update:modelValue')
    ▼
EditorView (接收 markdown 文本)
    │ :content="markdownText" + :custom-css="customCss"
    ▼
A4Page (渲染预览)
    │ renderMarkdown(content) → HTML
    │ injectCustomCSS(css) → <style> 注入
    │ templateClass → 应用模板 CSS 类名
    ▼
预览区实时更新

    │ 同时触发 useAutoSave (防抖 2s)
    ▼
UpdateJSON(id, { markdownContent }) → Go → SQLite
```

### 4.2 AI 操作（润色/翻译/缩写/重写）

```
MonacoEditor (用户选中文字)
    │ onDidChangeCursorSelection → updateSelection()
    ▼
AIFloatingToolbar (显示浮动按钮)
    │ @click → emit('operate', 'polish')
    ▼
MonacoEditor.handleAIOperation(operation)
    │ performAIOperation(operation)
    │ → useAIStream(operationId) 注册事件监听
    │ → sendMessage(id, prompt, jobTarget)
    ▼
Wails Bridge → Go AISendMessage()
    │ client.ChatStream() → Claude API
    │ StreamEvents() → EventsEmit("ai:stream:{id}", chunk)
    ▼
useAIStream (接收流式事件)
    │ EventsOn → content.value += chunk
    ▼
AIDiffView (显示原文 vs 修改对比)
    │ @accept → editor.executeEdits() 替换选区
```

### 4.3 AI 聊天

```
AIChatSidebar (用户输入消息)
    │ sendChatMessage(id, prompt, jobTarget, history, resumeContent)
    ▼
Wails Bridge → Go AISendChatMessage()
    │ BuildChatMessages() (含历史消息 + 简历内容)
    │ client.ChatStream() → EventsEmit chunks
    ▼
useAIStream (接收流式事件)
    │ watch(streamedContent) → 更新消息气泡
    ▼
saveChatMessage() → 持久化到 ai_messages 表
```

### 4.4 PDF 导出

```
ExportDialog (用户选择导出模式)
    │
    ├── 系统打印: window.print() + @media print CSS
    │
    └── 高级导出: 克隆 .preview-container DOM
                   → 注入 CSS 变量实际值
                   → 构建完整 HTML 文档
                   → ExportPDFFromHTML(html, path)
                   → Go: Chromedp PrintToPDF()
                   → 写入 PDF 文件
```

---

## 5. 数据库设计

SQLite，迁移文件在 `internal/database/migrations/`：

| 表名 | 迁移文件 | 用途 |
|------|---------|------|
| `resumes` | `001_create_resumes_table.sql` | 简历主表（title, markdown_content, json_data, template_id, custom_css, ...） |
| `settings` | `002_create_settings_table.sql` | 键值对设置（ai_api_key, ai_base_url, ...） |
| `ai_messages` | `003_create_ai_messages_table.sql` | AI 聊天历史（resume_id, role, content, quoted_text） |
| `snapshots` | `004_create_snapshots_table.sql` | 版本快照（resume_id, markdown_content, trigger_type, ...） |

**数据库位置**：`~/.darvin-resume/darvin-resume.db`（由 Wails 确定用户数据目录）

---

## 6. Wails 桥接方法一览

`app.go` 中定义的方法会自动暴露给前端，通过 `wailsjs/go/main/App.js` 调用：

### 简历管理
| Go 方法 | 前端调用 | 功能 |
|---------|---------|------|
| `CreateResume(title)` | `CreateResume(title)` | 创建新简历 |
| `GetResume(id)` | `GetResume(id)` | 获取简历详情 |
| `UpdateResume(id, json)` | `UpdateResume(id, json)` | 更新简历 |
| `DeleteResume(id)` | `DeleteResume(id)` | 删除简历（移入回收站） |
| `ListResumes()` | `ListResumes()` | 列出所有简历 |
| `UpdateResumeTemplate(id, tid)` | `UpdateResumeTemplate(id, tid)` | 更新模板 ID |
| `UpdateResumeCustomCSS(id, css)` | `UpdateResumeCustomCSS(id, css)` | 保存自定义 CSS |

### AI 功能
| Go 方法 | 功能 |
|---------|------|
| `AISendMessage(id, prompt, job, ctx)` | 流式 AI 文本处理 |
| `AISendMessageSync(id, prompt, job, ctx)` | 非流式 AI 调用 |
| `AISendChatMessage(id, prompt, job, history, content)` | 带历史的聊天 |
| `AICancelOperation(id)` | 取消正在进行的 AI 操作 |
| `GetAIConfig()` / `SaveAIConfig(cfg)` | AI 配置管理 |
| `ValidateAPIKey(key, url)` | 验证 API Key |
| `GetChatHistory(resumeId)` / `SaveChatMessage(msg)` / `ClearChatHistory(id)` | 聊天历史 |

### 快照管理
| Go 方法 | 功能 |
|---------|------|
| `CreateSnapshot(resumeId, label, note, trigger)` | 创建版本快照 |
| `ListSnapshots(resumeId)` | 列出快照 |
| `GetSnapshot(id)` | 获取快照详情 |
| `DiffSnapshots(id1, id2)` | 对比两个快照 |
| `RollbackToSnapshot(resumeId, snapshotId)` | 回滚到快照 |
| `DeleteSnapshot(id)` | 删除快照 |

### PDF 导出
| Go 方法 | 功能 |
|---------|------|
| `ExportPDFFromHTML(html, path)` | Chromedp 无头浏览器导出 PDF |
| `ShowSaveDialog(options)` | 系统文件保存对话框 |

### 设置
| Go 方法 | 功能 |
|---------|------|
| `GetSetting(key)` / `SetSetting(key, val)` | 通用键值对读写 |

---

## 7. CSS 架构（预览一致性核心）

**这是本项目最核心的架构约束**：编辑器预览与 PDF 导出必须 100% 渲染一致。

```
editor.css（全局，预览 + PDF 共用）
├── :root CSS 变量定义
│   ├── --resume-font-size: 10.5pt
│   ├── --resume-line-height: 1.6
│   ├── --resume-padding: 20mm
│   ├── --resume-primary-color: #1a1a1a
│   └── --resume-font-family: ...
│
├── .page-content 样式（使用 CSS 变量）
│   ├── font-size: var(--resume-font-size)
│   ├── line-height: var(--resume-line-height)
│   └── color: var(--resume-primary-color)
│
├── .a4-page 样式
│   ├── padding: var(--resume-padding)
│   └── background-color: var(--resume-bg-color)
│
└── Markdown 渲染样式（h1~h4, ul, table, a, hr 等）

templates/*.css（各模板覆盖样式）
├── template-minimal.css    → .template-minimal .page-content { ... }
├── template-campus.css     → .template-campus .page-content { ... }
├── template-academic.css   → .template-academic .page-content { ... }
└── template-dual-col.css   → .template-dual-col .page-content { ... }
```

**CSS 变量修改流程**：
```
StyleEditor UI → applyCSSVarsToPage() → document.documentElement.style.setProperty()
                                                ↓
                                     CSS 变量级联到 .page-content
```

---

## 8. 编码介入指南

### 8.1 开发环境搭建

```bash
# 1. 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 2. 安装前端依赖
cd frontend && npm install

# 3. 启动开发模式（热重载）
cd .. && wails dev

# 4. 构建发布版本
wails build
```

开发模式下：
- Go 后端修改自动重编译
- 前端修改自动热重载
- 访问 Wails 内置 DevTools：`http://localhost:34115`

---

### 8.2 添加新的 Wails 桥接方法

**步骤**：

1. 在 `app.go` 中添加方法（必须是 `*App` 的方法）：

```go
// app.go
func (a *App) MyNewMethod(param1 string, param2 int) (string, error) {
    // 调用 service/model 层
    result, err := service.DoSomething(param1, param2)
    if err != nil {
        return "", err
    }
    return result, nil
}
```

2. 运行 `wails dev` —— Wails 自动生成前端绑定文件：
   - `wailsjs/go/main/App.js`（JS 绑定）
   - `wailsjs/go/main/App.d.ts`（TS 类型）

3. 在前端调用：

```typescript
import { MyNewMethod } from '../wailsjs/go/main/App'

const result = await MyNewMethod('hello', 42)
```

**注意**：只有 `*App` 结构体的公开方法（大写开头）才会被暴露给前端。

---

### 8.3 添加新的前端组件

**步骤**：

1. 在 `frontend/src/components/` 创建 `.vue` 文件

2. 如果是页面级组件，在 `views/` 创建并在 `router/index.ts` 注册路由：

```typescript
// router/index.ts
{
  path: '/my-page',
  name: 'my-page',
  component: () => import('../views/MyNewView.vue')
}
```

3. 如果是 UI 组件，在父组件中 import：

```vue
<script setup>
import MyComponent from '../components/MyComponent.vue'
</script>

<template>
  <MyComponent :prop1="value" @event1="handler" />
</template>
```

---

### 8.4 添加新的 Composable

**步骤**：

1. 在 `frontend/src/composables/` 创建文件

2. 使用 Vue 3 组合函数模式：

```typescript
// composables/useMyFeature.ts
import { ref, onMounted, onUnmounted } from 'vue'
import { MyWailsMethod } from '../wailsjs/go/main/App'

export function useMyFeature() {
  const data = ref<string>('')
  const loading = ref(false)

  async function fetchData() {
    loading.value = true
    try {
      data.value = await MyWailsMethod()
    } catch (err) {
      console.error('Failed:', err)
    } finally {
      loading.value = false
    }
  }

  return { data, loading, fetchData }
}
```

3. 在组件中使用：

```vue
<script setup>
import { useMyFeature } from '../composables/useMyFeature'

const { data, loading, fetchData } = useMyFeature()
</script>
```

---

### 8.5 添加新的数据库表

**步骤**：

1. 在 `internal/database/migrations/` 创建新的迁移文件（递增编号）：

```sql
-- 005_create_my_table.sql
CREATE TABLE IF NOT EXISTS my_table (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

2. 在 `internal/model/` 创建数据模型：

```go
// model/my_model.go
package model

type MyModel struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"createdAt"`
}
```

3. 在 `internal/service/` 创建业务逻辑：

```go
// service/my_service.go
package service

import (
    "context"
    "database/sql"
    "Darvin-Resume/internal/database"
    "Darvin-Resume/internal/model"
)

func GetMyItems(ctx context.Context) ([]model.MyModel, error) {
    rows, err := database.DB.QueryContext(ctx, "SELECT id, name, created_at FROM my_table")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    // ... scan rows ...
}
```

4. 在 `app.go` 添加桥接方法（参考 8.2）

---

### 8.6 添加新的简历模板

**步骤**：

1. 在 `frontend/src/styles/templates/` 创建 CSS 文件：

```css
/* template-mytemplate.css */
.template-mytemplate .page-content {
    font-family: 'Georgia', serif;
    column-count: 2;
    column-gap: 20px;
}

.template-mytemplate .page-content h1 {
    font-size: 20pt;
    text-align: center;
    border-bottom: 2px solid var(--resume-primary-color);
}
```

2. 在 `A4Page.vue` 导入 CSS：

```typescript
// A4Page.vue <script setup>
import '../styles/templates/template-mytemplate.css'
```

3. 在 `A4Page.vue` 的 `validIds` 数组中添加模板 ID：

```typescript
const validIds = ['minimal', 'dual-col', 'academic', 'campus', 'mytemplate']
```

4. 在 `TemplateSelector.vue` 添加模板选项卡

---

### 8.7 修改 AI 提示词

AI 提示词定义在两个位置：

**编辑器 AI 操作（润色/翻译/缩写/重写）**：

```typescript
// frontend/src/composables/useAISelection.ts
const OPERATION_PROMPTS: Record<AIOperationType, PromptBuilder> = {
  polish: (text, job) => `请润色以下简历内容...：\n${text}`,
  translate: (text, _job) => `请将以下内容翻译成英文...：\n${text}`,
  summarize: (text, _job) => `请将以下内容压缩到原长度的60%...：\n${text}`,
  rewrite: (text, job) => `请重写以下简历内容...：\n${text}`,
}
```

**系统提示词（所有 AI 调用共享）**：

```go
// internal/ai/client.go
const SystemPrompt = `你是一位专业的简历优化助手。请根据用户需求优化简历内容。
要求：只返回优化后的文本内容，不要添加解释。`
```

---

### 8.8 调试技巧

| 场景 | 方法 |
|------|------|
| Go 后端调试 | `fmt.Println()` 输出到终端（wails dev 控制台） |
| 前端调试 | Wails DevTools: `http://localhost:34115`，或 `console.log()` |
| SQL 调试 | 在 `database/db.go` 开启 `_journal_mode=WAL` + 日志 |
| CSS 调试 | DevTools Elements 面板检查 computed styles |
| Wails 事件调试 | `EventsOn('ai:stream:*', console.log)` 检查事件流 |
| PDF 导出调试 | 检查 Chromedp 临时文件：`/tmp/darvin-resume-export-*.html` |

---

## 9. 关键文件速查

| 我想改... | 找哪个文件 |
|-----------|-----------|
| 简历列表页布局 | `views/HomeView.vue` |
| 编辑器工具栏 | `views/EditorView.vue`（顶部按钮区） |
| Monaco 编辑器行为 | `components/MonacoEditor.vue` |
| AI 浮动工具栏按钮 | `components/AIFloatingToolbar.vue` |
| AI 聊天对话 | `components/AIChatSidebar.vue` |
| AI Diff 对比 | `components/AIDiffView.vue` |
| A4 预览渲染 | `components/A4Page.vue` |
| 简历卡片样式 | `components/ResumeCard.vue` |
| 样式调整面板 | `components/StyleEditor.vue` |
| PDF 导出逻辑 | `components/ExportDialog.vue` + `internal/export/chromedp.go` |
| AI 流式响应 | `composables/useAIStream.ts` |
| AI 选区操作 | `composables/useAISelection.ts` |
| 自动保存 | `composables/useAutoSave.ts` |
| 模板管理 | `composables/useTemplate.ts` |
| AI 服务调用 | `services/ai.ts` |
| CSS 变量/预览样式 | `styles/editor.css` |
| 简历模板 CSS | `styles/templates/*.css` |
| CSS 白名单校验 | `utils/sanitizeCSS.ts` |
| Markdown 渲染 | `utils/markdown.ts` |
| 数据库表结构 | `internal/database/migrations/*.sql` |
| 简历 CRUD | `internal/service/resume.go` |
| AI API 调用 | `internal/ai/client.go` |
| AI 配置管理 | `internal/ai/config.go` |
| 桥接方法定义 | `app.go` |
| 应用入口/窗口配置 | `main.go` |
| 路由配置 | `router/index.ts` |
| TS 类型定义 | `types/*.ts` |

---

## 10. 常见开发场景

### 场景 1：给 AI 聊天加"引用选中文本"功能

```
1. MonacoEditor 暴露 getSelection() 方法（已有）
2. EditorView 打开 AIChatSidebar 时获取选中文本
3. 调用 AIChatSidebar.setQuotedText(selectedText)（方法已暴露）
4. AIChatSidebar 的 handleQuote() 使用 quotedText（已有 ref）
5. sendChatMessage() 中带上 quotedText
```

涉及的文件：`EditorView.vue`、`AIChatSidebar.vue`

### 场景 2：添加新的 AI 操作类型

```
1. types/ai.ts — AIOperationType 联合类型加新值
2. composables/useAISelection.ts — OPERATION_PROMPTS 加新 prompt
3. components/AIFloatingToolbar.vue — operations 数组加新按钮
4. AIDiffView.vue — OPERATION_LABELS 加新标签（自动同步）
```

### 场景 3：实现多页 PDF 导出

```
1. A4Page.vue — pages computed 中实现 HTML 按高度分页
2. A4Page.vue — 每页生成独立 .a4-page div
3. ExportDialog.vue — Chromedp 模式已有 page-break-after: always
4. internal/export/chromedp.go — PrintToPDF 应自动处理多页
```

### 场景 4：添加新的设置项

```
1. internal/settings/settings.go — 添加新 SettingKey 常量
2. internal/ai/config.go — LoadConfig/SaveConfig 读写新配置
3. frontend AIConfigModal.vue — 添加 UI 控件
4. frontend composables/useAIConfig.ts — 添加响应式状态
```
