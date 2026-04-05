<div align="center">

# Open-Resume

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

Open-Resume is a **local-first** resume builder designed for CS graduates and junior developers (0-3 years experience). All data is stored locally in SQLite. Bring your own AI API key — the platform never intercepts your requests.

## Why Open-Resume?

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
| 1. Project Skeleton & Data Layer | Wails init, SQLite, JSON schema, routing | Done |
| 2. Core Editor | Monaco, split-pane, real-time preview, A4 boundary | Done |
| 3. AI Core | Claude API, streaming, floating toolbar, diff, chat, parsing | Done |
| 4. Resume Management | List, search, AI-guided creation, trash, auto-save | Planned |
| 5. Templates & Export | Built-in templates, style tuning, PDF export, versioning | Planned |
| 6. Data Security | Encrypted API keys, backup & restore | Planned |
| 7. UI Polish | Dark mode, responsive, shortcuts, error handling | Planned |

## License

[MIT](LICENSE)
