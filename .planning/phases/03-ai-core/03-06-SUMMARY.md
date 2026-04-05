---
phase: 03-ai-core
plan: "03-06"
subsystem: ai
tags: [error-handling, toast, modal, wails, vue3, typescript, go]

# Dependency graph
requires:
  - phase: 03-01
    provides: Claude API adapter, SSE streaming, AI configuration infrastructure
provides:
  - Structured AI error type system (8 error codes)
  - Backend error classification with auto-retry for format errors
  - Frontend error toast component with severity-based UI
  - Token limit modal with actionable suggestions
  - AI config modal with validation
  - Unified error handling composable
  - User abort with content preservation via operation cancellation
affects: [03-02, 03-03, 03-04, 03-05] (all AI operation integration)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - VS Code-style error toast with severity icons and auto-dismiss
    - Backend error code classification via string matching
    - Operation cancellation via sync.Map (Go)
    - Singleton global error state for app-wide toast display
    - Format validation with 1 auto-retry before fallback to raw output

key-files:
  created:
    - internal/ai/error.go - AI error classification, Code constants, ToAIError()
    - frontend/src/types/ai.ts - AIErrorCode (8 codes), AIError interface, AIErrorMessages
    - frontend/src/components/AIErrorToast.vue - VS Code-style toast with retry + detail
    - frontend/src/components/TokenLimitModal.vue - Token limit modal with suggestions
    - frontend/src/components/AIConfigModal.vue - Full AI settings modal
    - frontend/src/composables/useAIError.ts - Unified error composable, modal state
    - frontend/src/composables/useAIStream.ts - Updated with structured error + abort
  modified:
    - internal/ai/client.go - Added ChatWithRetry, activeOperations sync.Map, CancelOperation
    - app.go - Added AICancelOperation bridge method
    - frontend/src/composables/useAIStream.ts - Added aiError ref, error conversion, abort()

key-decisions:
  - "Toast auto-dismisses for non-serious errors (3s), serious errors (AUTH_FAILED) require manual close"
  - "Format validation retries once before returning raw response as fallback"
  - "Content preserved on abort - streamingContent ref stays intact after cancel"
  - "Backend CancelOperation uses sync.Map for thread-safe operation tracking"
  - "AIConfigModal uses existing useAIConfig composable for all validation logic"

patterns-established:
  - "Error code normalization: backend string codes -> structured AIErrorCode on frontend"
  - "Singleton global error state via module-level ref for app-wide toast"
  - "Operation cancellation: frontend calls backend -> backend removes from sync.Map"

requirements-completed: [AIAI-10, AIAI-11, AIAI-12, AIAI-13]

# Metrics
duration: 10.3min
completed: 2026-04-05
---

# Plan 03-06: AI 异常兜底处理

**完整覆盖 AI 调用全异常场景：网络失败/Token 超限/格式异常/用户中断，提供清晰的错误提示和可执行的恢复方案，确保无数据丢失**

## Performance

- **Duration:** 10.3 min (618s)
- **Started:** 2026-04-05T12:25:02Z
- **Completed:** 2026-04-05T12:35:20Z
- **Tasks:** 6 commits (all 6 deliverables complete)
- **Files modified:** 8 files (5 created, 3 modified)

## Accomplishments

- 完整的 AI 错误类型体系：8 种错误码（NETWORK_ERROR/TIMEOUT/AUTH_FAILED/RATE_LIMIT/TOKEN_LIMIT/FORMAT_ERROR/ABORTED/UNKNOWN）
- 后端错误分类与自动重试：`ClassifyError()` + `ChatWithRetry()` 格式校验最多重试 1 次后回退原始输出
- VS Code 风格 Toast 通知组件：按严重程度着色（error/warning/info），自动消失（3s），重试按钮，开发者详情展开
- Token 超限 Modal：提供可操作的精简建议，支持一键聚焦编辑器或分批处理
- 统一错误处理 Hook：`useAIError()` 提供全局错误状态、错误转换、Modal 状态管理
- 用户中断内容保留：前端 `abort()` + 后端 `AICancelOperation()` 通过 `sync.Map` 追踪活跃操作

## Task Commits

1. **Task 1-2: 后端错误分类与类型扩展** - `b1f9ca8` (feat)
2. **Task 3-4: 前端错误处理集成** - `e767af3` (feat)
3. **Task 5: Vue 错误组件** - `0949b66` (feat)
4. **Task 6: useAIError composable** - `d52a6e5` (feat)
5. **Task 7: AIConfigModal** - `5f83436` (feat)

**Plan metadata:** `5f83436` (feat: complete plan 03-06)

## Files Created/Modified

### Backend (Go)

- `internal/ai/error.go` - 新增文件：AIErrorCode 常量、errorMapping 表、ClassifyError()、ToAIError()、AIError 结构体
- `internal/ai/client.go` - 新增：ChatWithRetry() 格式校验重试、activeOperations sync.Map、RegisterOperation/UnregisterOperation/CancelOperation
- `app.go` - 新增：AICancelOperation() 桥接方法

### Frontend (TypeScript/Vue)

- `frontend/src/types/ai.ts` - 新增：AIErrorCode 类型（8 种）、AIError 接口、AIErrorMessages、LegacyErrorCodeMap、ChatMessage/ChatConversation 接口
- `frontend/src/components/AIErrorToast.vue` - 新增文件：VS Code 风格 Toast，严重程度着色，自动消失，重试按钮，详情展开
- `frontend/src/components/TokenLimitModal.vue` - 新增文件：Token 超限 Modal，含精简建议列表和两个操作按钮
- `frontend/src/components/AIConfigModal.vue` - 新增文件：AI 设置 Modal，含全部 5 个配置项、验证按钮、保存逻辑
- `frontend/src/composables/useAIError.ts` - 新增文件：统一错误处理，错误转换，Modal 状态，全局单例
- `frontend/src/composables/useAIStream.ts` - 修改：新增 aiError ref、结构化错误转换、abort() 调用 AICancelOperation

## Decisions Made

- **Toast 自动消失策略**：非严重错误（NETWORK/TIMEOUT/RATE_LIMIT/FORMAT_ERROR/ABORTED）3 秒后自动消失；AUTH_FAILED 需要手动关闭
- **格式校验重试策略**：ChatWithRetry() 最多重试 1 次，第二次仍失败返回原始内容而非报错，确保数据不丢失
- **操作取消策略**：后端通过 sync.Map 追踪活跃操作，取消时移除映射；前端在 abort() 中保留已生成内容
- **错误码标准化**：后端保持原有字符串错误码，前端通过 convertBackendError() 统一转换为结构化 AIErrorCode

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

- **ChatMessage 类型冲突**：client.go 中定义的 ChatMessage 与 chat.go 中已存在的 ChatMessage 重复。修复：移除 client.go 中的重复定义，保留 chat.go 中的定义，并修正 app.go 中的类型转换

## Next Phase Readiness

- 完整的错误处理基础设施已就绪，可被 03-02（AI 操作集成）、03-03、03-04、03-05 使用
- 所有 AI 操作可通过 useAIError().showError() 统一处理错误
- AIConfigModal 可直接集成到设置页面或侧边栏
- 需注意：前端 Wails 绑定需在运行 `wails dev` 时重新生成（`AICancelOperation` 已在 App.d.ts 中自动生成）

---
*Plan: 03-06*
*Completed: 2026-04-05*

## Self-Check: PASSED

All 7 key files verified present. All 5 commit hashes verified in git log.
