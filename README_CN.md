<div align="center">

# Darvin-Resume

**Markdown 原生 · 隐私优先 · AI 深度协同 · 100% 本地运行**

[![Wails](https://img.shields.io/badge/Wails-v2-blue?logo=go)](https://wails.io/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-4.x-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[English](README.md)

</div>

---

> **编辑器预览与 PDF 导出 100% 排版一致，所见即所得，零排版焦虑。**

Darvin-Resume 是一款面向计算机专业应届生及毕业 3 年内初级开发者的**本地化简历制作与管理工具**。所有数据 100% 本地存储于 SQLite，用户自备 AI API Key，平台不截留任何请求数据。

## 核心优势

| 痛点 | 解决方案 |
|---|---|
| 导出后排版错乱 | 统一渲染引擎 (markdown-it)，预览 = 导出，所见即所得 |
| 多版本管理困难 | 本地 SQLite 结构化存储，JSON Schema 定义简历数据 |
| 在线 AI 工具会上传简历 | 100% 本地处理，自带 API Key，零隐私泄露 |
| 通用简历工具缺乏开发者特性 | Monaco Editor、Markdown 原生、VS Code 级编辑体验 |

## 功能特性

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
- AES-256-GCM 加密存储 API Key，设备唯一密钥
- 一键加密备份导出，支持每日/每周/每月自动备份
- 跨平台支持：Windows、macOS、Linux

**模板与导出**
- 4 套内置简历模板，切换后内容完整保留，仅样式变化
- 可视化样式调整：主色调、字号、行高、边距、字体族，实时预览
- 自定义 CSS 编辑器（白名单属性），一键重置默认样式
- PDF 导出（系统打印或 Chromedp），严格遵循 A4 尺寸和分页，与预览 100% 一致
- 版本快照管理：手动/自动创建、历史列表、Diff 对比、一键回滚

**简历管理**
- AI 引导式创建：分步勾选模块、填写内容、AI 实时润色，完成生成完整 Markdown
- 空白页创建：从头开始，全流程 AI 支持
- 简历列表：卡片视图、搜索、排序、重命名、复制、删除
- 回收站：30 天内可恢复，彻底删除
- 自动保存：每 30 秒、AI 完成时、页面切换时自动保存，支持 Ctrl/Cmd+S 手动保存

**界面与体验**
- 设计令牌系统：浅色/深色/跟随系统主题，CSS 变量驱动
- 响应式布局：窗口宽度 < 1200px 自动切换单栏，支持手动切换视图
- 快捷键体系：内置快捷键（Ctrl/Cmd+R/T/D），支持自定义并持久化存储
- 全局 Toast 通知，覆盖所有操作与异常场景
- 首页模板 Demo 预览，所有异步操作加载状态

## 技术栈

| 层级 | 技术 |
|---|---|
| 框架 | [Wails v2](https://wails.io/)（Go + WebView） |
| 前端 | Vue 3 + TypeScript + Vite |
| 编辑器 | [Monaco Editor](https://monaco-editor.github.io/) |
| 渲染 | [markdown-it](https://github.com/markdown-it/markdown-it) |
| 数据库 | SQLite（[modernc.org/sqlite](https://gitlab.com/cznic/sqlite)） |
| AI | Claude Messages API（SSE 流式传输） |

## 快速开始

### 环境要求

- [Go](https://go.dev/dl/) >= 1.25
- [Node.js](https://nodejs.org/) >= 16
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发模式

```bash
git clone https://github.com/darven-cs/darvin-resume.git
cd darvin-resume
cd frontend && npm install && cd ..
wails dev
```

### 构建

```bash
wails build
```

构建产物位于 `build/bin/`。

## 项目结构

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
│       ├── components/       # Vue 组件（MonacoEditor, A4Page, AI* ...）
│       ├── composables/      # Vue 组合式函数（useAI*, useDebounce）
│       ├── styles/           # 共享 CSS（预览 + 导出）
│       └── views/            # 页面视图
├── docs/                     # 文档 & 需求规格
└── wails.json               # Wails 项目配置
```

## 开发路线

| 阶段 | 描述 | 状态 |
|---|---|---|
| 1. 项目骨架与数据层 | Wails 初始化、SQLite、JSON Schema、路由 | 已完成 |
| 2. 核心编辑器 | Monaco、双栏布局、实时预览、A4 分页线 | 已完成 |
| 3. AI 核心能力 | Claude API、流式传输、浮动工具栏、Diff、对话、解析 | 已完成 |
| 4. 简历创建与管理 | 列表、搜索、AI 引导创建、回收站、自动保存 | ✅ 已完成 |
| 5. 模板与导出 | 内置模板、样式调整、PDF 导出、版本快照 | ✅ 已完成 |
| 6. 数据安全 | AES-256 加密存储、备份与恢复、自动备份 | ✅ 已完成 |
| 7. 界面打磨 | 设计令牌、主题切换、响应式、快捷键、Toast、异常兜底 | ✅ 已完成 |

## 许可证

[MIT](LICENSE)
