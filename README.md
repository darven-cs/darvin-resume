<div align="center">

# Open-Resume

**Markdown-Native, Privacy-First, AI-Powered Resume Builder**

Markdown 原生 | 隐私优先 | AI 深度协同 | 100% 本地运行

[![Wails](https://img.shields.io/badge/Wails-v2.12--blue?logo=go)](https://wails.io/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-4.x-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

> **Editor preview matches PDF export exactly. What you see is what you get.**
>
> **编辑器预览与 PDF 导出 100% 排版一致，所见即所得，零排版焦虑。**

---

## English

### What is Open-Resume?

Open-Resume is a **local-first** resume builder designed for CS graduates and junior developers (0-3 years experience). It solves four core pain points: layout chaos, multi-version management difficulty, low content polishing efficiency, and resume privacy leaks.

All data is stored locally in SQLite. Bring your own AI API key — the platform never intercepts your requests.

### Why Open-Resume?

| Pain Point | Solution |
|---|---|
| Layout breaks after export | Unified rendering engine (markdown-it) — preview = export |
| Hard to manage multiple versions | Local SQLite storage with structured JSON schema |
| AI tools upload your resume to the cloud | 100% local processing, BYOK (Bring Your Own Key) |
| Generic online resume builders lack developer features | Monaco Editor, Markdown-native, VS Code-like experience |

### Key Features

**Editor**
- Monaco Editor with Markdown syntax highlighting and VS Code keybindings
- Split-pane layout with real-time preview (< 200ms latency)
- A4 page boundary visualization
- Drag-to-resize pane width

**AI Integration (Claude API)**
- Floating toolbar: select text → polish / translate / abbreviate / rewrite
- Diff view: accept or reject AI changes with visual comparison
- Chat sidebar: multi-turn conversation with context, one-click insert
- Resume parser: paste old resume → auto-generate structured Markdown
- SSE streaming for real-time AI response rendering

**Privacy & Security**
- All data stored locally in SQLite — no cloud upload
- BYOK: bring your own API key, platform never intercepts requests
- Cross-platform: Windows, macOS, Linux

### Tech Stack

| Layer | Technology |
|---|---|
| Framework | [Wails v2](https://wails.io/) (Go + WebView) |
| Frontend | Vue 3 + TypeScript + Vite |
| Editor | [Monaco Editor](https://monaco-editor.github.io/) |
| Rendering | [markdown-it](https://github.com/markdown-it/markdown-it) |
| Database | SQLite (via [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)) |
| AI | Claude Messages API (SSE streaming) |

### Getting Started

#### Prerequisites

- [Go](https://go.dev/dl/) >= 1.25
- [Node.js](https://nodejs.org/) >= 16
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### Development

```bash
# Clone the repository
git clone https://github.com/darven-cs/darvin-resume.git
cd darvin-resume

# Install frontend dependencies
cd frontend && npm install && cd ..

# Run in development mode (hot reload)
wails dev
```

#### Build

```bash
# Build production binary
wails build
```

The binary will be in `build/bin/`.

### Project Structure

```
├── app.go                    # Wails app bindings (Go backend)
├── main.go                   # Application entry point
├── internal/
│   ├── ai/                   # Claude API client & SSE streaming
│   ├── converter/            # JSON ↔ Markdown conversion
│   ├── database/             # SQLite layer & migrations
│   ├── model/                # Data models
│   ├── service/              # Business logic services
│   └── settings/             # Key-value settings persistence
├── frontend/
│   └── src/
│       ├── components/       # Vue components
│       │   ├── MonacoEditor.vue      # Code editor
│       │   ├── A4Page.vue            # A4 preview container
│       │   ├── SplitPane.vue         # Resizable split layout
│       │   ├── AIFloatingToolbar.vue # AI selection toolbar
│       │   ├── AIDiffView.vue        # Diff comparison view
│       │   ├── AIChatSidebar.vue     # AI chat panel
│       │   └── ...
│       ├── composables/      # Vue composables (useAI*, useDebounce)
│       ├── styles/           # Shared CSS (preview + export)
│       └── views/            # Page views
├── docs/                     # Documentation & requirements
└── wails.json               # Wails project config
```

### Roadmap

| Phase | Description | Status |
|---|---|---|
| 1. Project Skeleton & Data Layer | Wails init, SQLite, JSON schema, routing | Done |
| 2. Core Editor | Monaco, split-pane, real-time preview, A4 boundary | Done |
| 3. AI Core | Claude API, streaming, floating toolbar, diff, chat, parsing | Done |
| 4. Resume Management | List, search, AI-guided creation, trash, auto-save | Planned |
| 5. Templates & Export | Built-in templates, style tuning, PDF export, versioning | Planned |
| 6. Data Security | Encrypted API keys, backup & restore | Planned |
| 7. UI Polish | Dark mode, responsive, shortcuts, error handling | Planned |

### License

[MIT](LICENSE)

---

## 中文

### Open-Resume 是什么？

Open-Resume 是一款面向计算机专业应届生及毕业 3 年内初级开发者的**本地化简历制作与管理工具**。基于 Markdown 原生编写，所有数据 100% 本地存储于 SQLite，用户自备 AI API Key，平台不截留任何请求数据。

### 核心优势

| 痛点 | 解决方案 |
|---|---|
| 导出后排版错乱 | 统一渲染引擎 (markdown-it)，预览 = 导出，所见即所得 |
| 多版本管理困难 | 本地 SQLite 结构化存储，JSON Schema 定义简历数据 |
| 在线 AI 工具会上传简历 | 100% 本地处理，自带 API Key，零隐私泄露 |
| 通用简历工具缺乏开发者特性 | Monaco Editor、Markdown 原生、VS Code 级编辑体验 |

### 功能特性

**编辑器**
- Monaco Editor：Markdown 语法高亮，VS Code 兼容键位
- 双栏布局：左侧编辑 + 右侧实时预览，延迟 < 200ms
- A4 分页线可视化，拖拽调整栏宽

**AI 协同（Claude API）**
- 浮动工具栏：选中文本 → 润色 / 翻译 / 缩写 / 重写
- Diff 对比视图：可视化比较 AI 修改，逐项接受/拒绝
- AI 对话侧边栏：多轮对话、引用文本、一键插入
- 简历解析：粘贴旧简历 → 自动生成结构化 Markdown
- SSE 流式传输，打字机效果实时渲染

**隐私与安全**
- 全量数据本地 SQLite 存储，无强制云端上传
- 用户自备 API Key，平台不截留请求数据
- 跨平台支持：Windows、macOS、Linux

### 技术栈

| 层级 | 技术 |
|---|---|
| 框架 | [Wails v2](https://wails.io/)（Go + WebView） |
| 前端 | Vue 3 + TypeScript + Vite |
| 编辑器 | [Monaco Editor](https://monaco-editor.github.io/) |
| 渲染 | [markdown-it](https://github.com/markdown-it/markdown-it) |
| 数据库 | SQLite（[modernc.org/sqlite](https://gitlab.com/cznic/sqlite)） |
| AI | Claude Messages API（SSE 流式传输） |

### 快速开始

#### 环境要求

- [Go](https://go.dev/dl/) >= 1.25
- [Node.js](https://nodejs.org/) >= 16
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### 开发模式

```bash
# 克隆仓库
git clone https://github.com/darven-cs/darvin-resume.git
cd darvin-resume

# 安装前端依赖
cd frontend && npm install && cd ..

# 启动开发模式（热重载）
wails dev
```

#### 构建

```bash
# 构建生产版本
wails build
```

构建产物位于 `build/bin/`。

### 项目结构

```
├── app.go                    # Wails 应用绑定（Go 后端）
├── main.go                   # 应用入口
├── internal/
│   ├── ai/                   # Claude API 客户端 & SSE 流式传输
│   ├── converter/            # JSON ↔ Markdown 转换
│   ├── database/             # SQLite 存储层 & 数据库迁移
│   ├── model/                # 数据模型
│   ├── service/              # 业务逻辑服务
│   └── settings/             # 键值对配置持久化
├── frontend/
│   └── src/
│       ├── components/       # Vue 组件
│       │   ├── MonacoEditor.vue      # 代码编辑器
│       │   ├── A4Page.vue            # A4 预览容器
│       │   ├── SplitPane.vue         # 可拖拽分栏布局
│       │   ├── AIFloatingToolbar.vue # AI 选区工具栏
│       │   ├── AIDiffView.vue        # Diff 对比视图
│       │   ├── AIChatSidebar.vue     # AI 对话面板
│       │   └── ...
│       ├── composables/      # Vue 组合式函数（useAI*, useDebounce）
│       ├── styles/           # 共享 CSS（预览 + 导出）
│       └── views/            # 页面视图
├── docs/                     # 文档 & 需求规格
└── wails.json               # Wails 项目配置
```

### 开发路线

| 阶段 | 描述 | 状态 |
|---|---|---|
| 1. 项目骨架与数据层 | Wails 初始化、SQLite、JSON Schema、路由 | 已完成 |
| 2. 核心编辑器 | Monaco、双栏布局、实时预览、A4 分页线 | 已完成 |
| 3. AI 核心能力 | Claude API、流式传输、浮动工具栏、Diff、对话、解析 | 已完成 |
| 4. 简历创建与管理 | 列表、搜索、AI 引导创建、回收站、自动保存 | 计划中 |
| 5. 模板与导出 | 内置模板、样式调整、PDF 导出、版本管理 | 计划中 |
| 6. 数据安全 | API Key 加密存储、备份与恢复 | 计划中 |
| 7. 界面打磨 | 深色模式、响应式、快捷键、异常兜底 | 计划中 |

### 许可证

[MIT](LICENSE)
