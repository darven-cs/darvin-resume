---
phase: 03-ai-core
plan: 03-01
subsystem: ai
tags: [claude-api, sse, streaming, wails, vue3, typescript]

# Dependency graph
requires:
  - phase: 02-core-editor
    provides: EditorView, Monaco integration, split pane layout
provides:
  - Claude Messages API HTTP client with streaming (SSE) support
  - AI configuration persistence via SQLite settings table
  - Frontend SSE stream composable with Wails EventsOn
  - AI service layer and config composable
affects: [03-02, 03-03, 03-04] (AI operation integration)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - SSE streaming via Wails EventsEmit/EventsOn
    - AI config stored in settings table with ai.* keys
    - Backend-first API key (never exposed to frontend)
    - Debounced DOM updates for typewriter effect

key-files:
  created:
    - internal/ai/config.go - AIConfig struct, error codes, defaults
    - internal/ai/client.go - Claude API HTTP client, SSE streaming, message building
    - internal/ai/client_test.go - 14 unit tests covering all paths
    - internal/settings/settings.go - SQLite key-value CRUD for settings table
    - app.go - AI bridge methods (GetAIConfig, SaveAIConfig, ValidateAPIKey, AISendMessage, AISendMessageSync)
    - frontend/src/types/ai.ts - AIStreamChunk, AIConfig, error codes/messages
    - frontend/src/composables/useAIStream.ts - SSE handling via Wails EventsOn
    - frontend/src/composables/useAIConfig.ts - Config state, validation, persistence
    - frontend/src/services/ai.ts - Bridge call wrappers
  modified:
    - app.go - Added AI bridge methods and runtime import
    - frontend/vite.config.ts - Added @ alias
    - frontend/tsconfig.json - Added paths.@/* for @ imports

key-decisions:
  - "API key stored plaintext this phase (Phase 6 adds encryption)"
  - "SSE chunks emitted as Wails EventsOn events with ai:stream:{operationId} namespace"
  - "16ms debounce for DOM updates (~60fps typewriter effect)"
  - "Claude Messages API POST /v1/messages with stream:true for streaming"
  - "Backend sends system prompt as a user message (simpler than role distinction)"

patterns-established:
  - "AI bridge methods always load/save via settings service"
  - "StreamEvents() parses SSE line-by-line, handles both content_block_delta and backwards-compat delta formats"
  - "Frontend uses local wailsjs paths (not npm @wailsio/runtime)"

requirements-completed: [AIAI-01, AIAI-03, AIAI-07]

# Metrics
duration: 6.5min
completed: 2026-04-05
---

# Plan 03-01: Claude API 适配与 SSE 流式传输

**Claude Messages API 集成，支持自定义 BaseURL，流式 SSE 响应通过 Wails Events 传输到前端，打字机效果渲染**

## Performance

- **Duration:** 6.5 min (392s)
- **Started:** 2026-04-05T12:13:32Z
- **Completed:** 2026-04-05T12:20:04Z
- **Tasks:** 9 (6 commits + frontend build)
- **Files modified:** 13 files (5 created backend, 6 created frontend, 2 modified)

## Accomplishments

- Go AI 客户端支持 Claude Messages API 流式和非流式调用
- AI 配置通过 SQLite settings 表持久化（ai.* keys）
- 前端 SSE composable 通过 Wails EventsOn 监听流式响应
- 完整的错误处理体系（network/auth/rate_limit/api/timeout/cancelled）
- 14 个 Go 单元测试覆盖所有代码路径
- TypeScript 前端类型、服务层、composables 完整实现

## Task Commits

1. **AI 配置模型** - `f5d9722` (feat): AIConfig struct, error codes, defaults
2. **Settings 服务** - `b2ca208` (feat): SQLite key-value CRUD
3. **Claude API 客户端** - `32c6e51` (feat): HTTP client, SSE streaming, message building
4. **AI 桥接方法** - `936863d` (feat): GetAIConfig, SaveAIConfig, ValidateAPIKey, AISendMessage, AISendMessageSync
5. **AI 客户端单元测试** - `0c45bd5` (test): 14 tests (Chat success/error/timeout, ValidateAPIKey, StreamEvents, BuildMessages, Validate)
6. **前端 AI 基础设施** - `73e93eb` (feat): types, composables, service layer

**Plan metadata:** `73e93eb` (feat: complete plan 03-01)

## Files Created/Modified

### Backend (Go)

- `internal/ai/config.go` - AIConfig 结构体，错误码，网络/认证/限流/API 错误定义，默认值常量，Validate() 方法
- `internal/ai/client.go` - Claude Messages API HTTP 客户端，Chat()/ChatStream() 方法，StreamEvents() SSE 解析，ValidateAPIKey() 验证，BuildMessages() 消息构建，LoadConfig/SaveConfig 持久化，SystemPrompt 常量
- `internal/ai/client_test.go` - 14 个测试覆盖成功/失败/超时路径
- `internal/settings/settings.go` - Get/Set/Delete/GetWithDefault for settings 表
- `app.go` - 新增 6 个 AI 桥接方法 + 辅助函数，导入 wails runtime

### Frontend (TypeScript/Vue)

- `frontend/src/types/ai.ts` - AIStreamChunk, AIConfig 接口，AIErrorCodes/AIErrorMessages 常量
- `frontend/src/composables/useAIStream.ts` - useAIStream(operationId) composable，EventsOn 监听，16ms debounce，abort/reset
- `frontend/src/composables/useAIConfig.ts` - useAIConfig() composable，loadConfig/persistConfig/validateKey，表单验证
- `frontend/src/services/ai.ts` - getConfig/saveConfig/validateAPIKey/sendMessage/sendMessageSync 桥接调用
- `frontend/vite.config.ts` - 添加 @ -> src 路径别名
- `frontend/tsconfig.json` - 添加 baseUrl 和 paths.@/* 配置

## Decisions Made

- **API Key 明文存储**：本阶段明文，Phase 6 加密（AIAI-02）
- **SSE 事件命名**：`ai:stream:{operationId}` 命名空间，支持多操作隔离
- **Debounce 间隔**：16ms（约 60fps）确保打字机效果流畅
- **System Prompt 策略**：作为 user 消息发送（比 role 区分更简单）
- **前端 Wails 导入**：使用本地 `wailsjs/wailsjs/runtime/runtime` 而非 npm `@wailsio/runtime`

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

- **TypeScript @ 路径别名未配置**：tsconfig.json 缺少 paths 配置，vite.config.ts 缺少 resolve.alias。修复：添加了 tsconfig baseUrl + paths 和 vite resolve.alias
- **Wails 类型生成**：`wailsjs/wailsjs/go/main/App.d.ts` 已自动生成新方法类型，无需手动更新

## Next Phase Readiness

- AI 后端基础设施完成，可被 Phase 03-02 及后续阶段调用
- 前端 SSE composable 可直接用于 03-02 的 AI 操作集成
- AI 配置 UI 可在 03-02 中使用 useAIConfig composable
- 需要注意：Wails 每次运行需重新生成前端绑定（`wails generate` 或开发模式自动）

---
*Plan: 03-01*
*Completed: 2026-04-05*
