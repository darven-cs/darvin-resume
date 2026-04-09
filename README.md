<div align="center">

# Darvin-Resume

**Markdown-Native, Privacy-First, AI-Powered Resume Builder**

[![Wails](https://img.shields.io/badge/Wails-v2-blue?logo=go)](https://wails.io/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-4.x-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[中文文档](README_CN.md)

</div>

---

> **Editor preview matches PDF export exactly. What you see is what you get.**

Darvin-Resume is a **local-first** resume builder designed for CS graduates and junior developers (0-3 years experience). All data is stored locally in SQLite. Bring your own AI API key — the platform never intercepts your requests.

## Why Darvin-Resume?

| Pain Point | Solution |
|---|---|
| Layout breaks after export | Unified rendering engine (markdown-it) — preview = export |
| Hard to manage multiple versions | Local SQLite storage with structured JSON schema |
| AI tools upload your resume to the cloud | 100% local processing, BYOK (Bring Your Own Key) |
| Generic resume builders lack dev features | Monaco Editor, Markdown-native, VS Code-like experience |

## Features

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

**Resume Management**
- AI-guided creation wizard: step-by-step module selection, field filling, real-time AI polish
- Blank page mode: start from scratch with full AI support
- Resume list: card grid, search, sort, rename, duplicate, delete
- Recycle bin: recover or permanently delete within 30 days
- Auto-save: 30-second interval, AI completion trigger, page navigation trigger, Ctrl/Cmd+S manual save

**Templates & Export**
- 4 built-in resume templates with instant switching, content preserved
- Visual style tuning: color, font size, line height, margins, font family — real-time preview
- Custom CSS editor with property whitelist, one-click reset to defaults
- PDF export via system print or Chromedp, strict A4 pagination matching preview 100%
- Version snapshots: manual/auto creation, diff comparison, one-click rollback

**Data Security & Backup**
- AES-256-GCM encrypted API key storage with device-unique key
- One-click encrypted backup export, configurable auto-backup (daily/weekly/monthly)
- Backup import with automatic pre-restore snapshot, full data integrity

**UI & Experience**
- Design token system: light/dark/system theme with CSS custom properties
- Responsive layout: auto single-pane below 1200px, manual view switching
- Keyboard shortcuts: built-in (Ctrl/Cmd+R/T/D) + customizable with persistence
- Global toast notifications for all operations with error recovery
- Template demo previews on empty home page, loading states on all async operations

## Tech Stack

| Layer | Technology |
|---|---|
| Framework | [Wails v2](https://wails.io/) (Go + WebView) |
| Frontend | Vue 3 + TypeScript + Vite |
| Editor | [Monaco Editor](https://monaco-editor.github.io/) |
| Rendering | [markdown-it](https://github.com/markdown-it/markdown-it) |
| Database | SQLite (via [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)) |
| AI | Claude Messages API (SSE streaming) |

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) >= 1.25
- [Node.js](https://nodejs.org/) >= 16
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Development

```bash
git clone https://github.com/darven-cs/darvin-resume.git
cd darvin-resume
cd frontend && npm install && cd ..
wails dev
```

### Build

```bash
wails build
```

The binary will be in `build/bin/`.

## Project Structure

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
│       ├── components/       # Vue components (MonacoEditor, A4Page, AI* ...)
│       ├── composables/      # Vue composables (useAI*, useDebounce)
│       ├── styles/           # Shared CSS (preview + export)
│       └── views/            # Page views
├── docs/                     # Documentation & requirements
└── wails.json               # Wails project config
```

## Roadmap

| Phase | Description | Status |
|---|---|---|
| 1. Project Skeleton & Data Layer | Wails init, SQLite, JSON schema, routing | ✅ Done |
| 2. Core Editor | Monaco, split-pane, real-time preview, A4 boundary | ✅ Done |
| 3. AI Core | Claude API, streaming, floating toolbar, diff, chat, parsing | ✅ Done |
| 4. Resume Management | List, search, AI wizard, blank page, recycle bin, auto-save | ✅ Done |
| 5. Templates & Export | Built-in templates, style tuning, PDF export, version snapshots | ✅ Done |
| 6. Data Security | AES-256 encrypted API keys, backup & restore, auto-backup | ✅ Done |
| 7. UI Polish | Design tokens, light/dark theme, responsive, shortcuts, toast, error handling | ✅ Done |

## License

[MIT](LICENSE)
