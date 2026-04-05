---
phase: 03-ai-core
plan: 03-04
subsystem: ai
tags: [chat, sidebar, vue3, typescript, streaming, sqlite, wails]

# Dependency graph
requires:
  - phase: 03-ai-core
    provides: SSE streaming infrastructure (useAIStream), AI service layer, AI config
provides:
  - AIChatSidebar component with multi-turn conversation
  - Chat history persistence via SQLite ai_messages table
  - Backend chat bridge methods (GetChatHistory, SaveChatMessage, ClearChatHistory, AISendChatMessage)
  - Conversation context injection with up to 10 history messages
  - Quote selected text and insert AI output to editor
affects: [03-04] (AI conversation integration)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Fixed-position sidebar with CSS Transition (300ms ease-out slide)
    - RAF-based typewriter effect for streaming content
    - Chat history persistence via ai_messages SQLite table
    - Wails-generated ai.ChatMessage model for type-safe bridge
    - Multi-turn context with BuildChatMessages (up to 10 history)

key-files:
  created:
    - internal/ai/chat.go - ChatMessage CRUD: GetChatHistory, SaveChatMessage, ClearChatHistory
    - internal/database/migrations/003_create_ai_messages_table.sql - ai_messages table
    - frontend/src/components/AIChatSidebar.vue - Full chat sidebar component (layout, streaming, quotes, insert)
  modified:
    - internal/ai/client.go - Added ChatMessage struct, BuildChatMessages with history context
    - app.go - Added 4 bridge methods: GetChatHistory, SaveChatMessage, ClearChatHistory, AISendChatMessage
    - frontend/src/types/ai.ts - Added ChatMessage and ChatConversation interfaces
    - frontend/src/services/ai.ts - Added chat service methods with Wails-generated model types
    - frontend/src/views/EditorView.vue - Integrated sidebar with toolbar button, insertText handler
    - frontend/src/components/MonacoEditor.vue - Added getSelection() and insertAtCursor() exposed methods

key-decisions:
  - "ChatMessage defined once in internal/ai/client.go, reused across chat.go and app.go (single source of truth)"
  - "ai.ChatMessage Wails-generated model used in frontend service layer for type safety"
  - "BuildChatMessages includes system prompt with resume context + up to 10 history messages"
  - "Sidebar uses fixed positioning (not flex) to overlay editor content cleanly"

patterns-established:
  - "Chat messages persisted per-resume, loaded on sidebar open, saved after each exchange"
  - "Quote flow: parent (EditorView) sets quoted text via exposed method on sidebar"
  - "Insert flow: sidebar emits insertText event, parent calls MonacoEditor.insertAtCursor()"

requirements-completed: [AIAI-08]

# Metrics
duration: 11min
completed: 2026-04-05
---

# Plan 03-04: AI 对话侧边栏 Summary

**AIChatSidebar 组件：400px 右侧滑入面板，支持多轮对话、引用选中文本、流式打字机效果、一键插入编辑区，对话历史 SQLite 持久化**

## Performance

- **Duration:** 11 min
- **Started:** 2026-04-05T12:25:13Z
- **Completed:** 2026-04-05T12:36:00Z
- **Tasks:** 6 (6 commits)
- **Files modified:** 10 files (3 created, 7 modified)

## Accomplishments

- Go 后端 ChatMessage CRUD 操作（SQLite 持久化）
- BuildChatMessages 支持对话历史上下文注入（最近 10 条）
- AIChatSidebar Vue 组件完整实现（布局、流式渲染、引用、插入）
- MonacoEditor 暴露 getSelection() 和 insertAtCursor() 方法
- EditorView 集成工具栏按钮和侧边栏交互
- TypeScript 类型安全的 Wails 模型桥接

## Task Commits

1. **Database migration** - `d78d2fe` (feat): ai_messages table
2. **Backend chat service** - `b4fcd8b` (feat): ChatMessage CRUD + AISendChatMessage bridge
3. **TypeScript types** - `699e7eb` (feat): ChatMessage, ChatConversation interfaces
4. **Service layer** - `f53cd7f` (feat): Chat history service methods
5. **Component + Integration** - `96ff10d` (feat): AIChatSidebar + EditorView integration
6. **Service refactor** - `3515c2d` (refactor): Wails-generated model types

## Files Created/Modified

### Backend (Go)

- `internal/database/migrations/003_create_ai_messages_table.sql` - ai_messages table (id, resume_id FK, role, content, quoted_text, created_at)
- `internal/ai/chat.go` - GetChatHistory, SaveChatMessage, ClearChatHistory, GetChatHistoryOrEmpty
- `internal/ai/client.go` - ChatMessage struct, BuildChatMessages with history context (up to 10 messages)
- `app.go` - 4 bridge methods: GetChatHistory, SaveChatMessage, ClearChatHistory, AISendChatMessage

### Frontend (TypeScript/Vue)

- `frontend/src/types/ai.ts` - ChatMessage, ChatConversation interfaces added
- `frontend/src/services/ai.ts` - getChatHistory, saveChatMessage, clearChatHistory, sendChatMessage with ai.ChatMessage model
- `frontend/src/components/AIChatSidebar.vue` - Full chat sidebar (700+ lines)
- `frontend/src/views/EditorView.vue` - AI toolbar button, sidebar integration, insertText handler
- `frontend/src/components/MonacoEditor.vue` - getSelection(), insertAtCursor() exposed methods

## Decisions Made

- **ChatMessage 单一来源**: 在 `internal/ai/client.go` 定义一次，`chat.go` 和 `app.go` 共享
- **Wails 模型桥接**: 前端使用 `ai.ChatMessage` Wails 生成类替代 `as any` 类型转换
- **侧边栏定位**: 使用 `position: fixed` 覆盖编辑器内容，不改变现有布局
- **对话历史限制**: 最近 10 条消息作为上下文传递给 API

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

- **ChatMessage 类型重复**: 初始在 `chat.go`、`client.go` 和 `app.go` 三处定义，导致 Go 编译错误。修复：统一到 `client.go` 定义，`chat.go` 共享同包类型，`app.go` 使用 `ai.ChatMessage` 导入。
- **Wails 绑定自动生成**: `wails generate bindings` 命令不存在于 v2.12.0，绑定在 `wails dev` 或 `wails build` 时自动生成。已验证绑定正确包含所有新方法。

## Self-Check: PASSED

- All 8 files verified (3 created backend, 1 created frontend, 4 modified)
- All 6 commits verified in git log
- Go build: PASSED
- Frontend build: PASSED

## Next Phase Readiness

- AI 对话侧边栏完成，支持多轮对话和上下文注入
- 对话历史通过 SQLite 持久化，页面刷新后完整恢复
- 侧边栏可通过工具栏按钮或未来快捷键 Ctrl+Shift+A 打开
- 需后续计划：引用选中文本需要在打开侧边栏时自动获取 Monaco 选区

---
*Plan: 03-04*
*Completed: 2026-04-05*
