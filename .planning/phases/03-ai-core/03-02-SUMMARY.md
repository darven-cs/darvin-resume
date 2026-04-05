---
phase: 03-ai-core
plan: 03-02
subsystem: ai
tags: [ai-toolbar, monaco, selection, streaming, vue3, typescript]

# Dependency graph
requires:
  - phase: 03-01
    provides: useAIStream, sendMessage, AI SSE infrastructure
provides:
  - AI selection floating toolbar with 4 operations
  - useAISelection composable
  - Monaco selection integration
affects: [03-03, 03-04, 03-05] (AI operation integration)

# Tech tracking
tech-stack:
  added:
    - AIFloatingToolbar.vue - Vue component with Teleport, position fixed
    - useAISelection.ts - Selection detection, position calculation, AI operation
  patterns:
    - Monaco onDidChangeCursorSelection for selection detection
    - crypto.randomUUID() for operation ID generation
    - Vue Teleport for z-index isolation

key-files:
  created:
    - frontend/src/components/AIFloatingToolbar.vue - 4-button toolbar with icons, spinner
    - frontend/src/composables/useAISelection.ts - Selection state, AI streaming integration
  modified:
    - frontend/src/components/MonacoEditor.vue - Added AI toolbar + selection listener
    - frontend/src/views/EditorView.vue - Added jobTarget prop pass-through
    - frontend/src/components/ResumeParserModal.vue - Fixed uuid -> crypto.randomUUID

key-decisions:
  - "Position fixed with Teleport to body for z-index isolation"
  - "SelectionRange typed as plain interface to avoid Monaco namespace import issues"
  - "Toolbar shows above selection start position, flips if overflow"
  - "4 operations: polish/translate/summarize/rewrite with context-aware prompts"

patterns-established:
  - "Monaco selection detection via onDidChangeCursorSelection event"
  - "Position calculation via getScrolledVisiblePosition + getBoundingClientRect"
  - "AI operation via useAIStream + editor.executeEdits replacement"

requirements-completed: [EDIT-07]

# Metrics
duration: 4.4min
completed: 2026-04-05
---

# Plan 03-02: 选区浮动工具栏

**AI 选区浮动工具栏：用户框选文本后显示 4 个 AI 操作按钮（润色/翻译/缩写/重写），流式返回结果并替换选区**

## Performance

- **Duration:** 4.4 min (265s)
- **Started:** 2026-04-05T12:24:52Z
- **Completed:** 2026-04-05T12:29:17Z
- **Tasks:** 3 (1 commit)
- **Files modified:** 5 files (2 created, 3 modified)

## Accomplishments

- AIFloatingToolbar 组件：VS Code 暗黑风格，4 操作按钮横向排列，loading spinner
- useAISelection composable：Monaco 选区检测、位置计算、AI 流式操作
- MonacoEditor 集成：onDidChangeCursorSelection 监听，executeEdits 替换选区
- 职位目标上下文传递：jobTarget prop 流经 EditorView -> MonacoEditor -> useAISelection
- 构建验证通过：`npm run build` 无错误

## Task Commits

1. **AI 选区浮动工具栏** - `5687db1` (feat): AIFloatingToolbar + useAISelection + MonacoEditor 集成 + EditorView jobTarget + uuid 修复

## Files Created/Modified

### Created (frontend)

- `frontend/src/components/AIFloatingToolbar.vue` - Vue 组件，4 按钮 + Teleport + 位置计算
- `frontend/src/composables/useAISelection.ts` - Composable，选区检测、位置计算、AI 流式操作

### Modified (frontend)

- `frontend/src/components/MonacoEditor.vue` - 添加 AIFloatingToolbar 组件、选区监听、AI 操作处理
- `frontend/src/views/EditorView.vue` - 添加 jobTarget 状态和 prop 传递
- `frontend/src/components/ResumeParserModal.vue` - uuid -> crypto.randomUUID() 修复

## Decisions Made

- **位置策略**：position fixed + Teleport to body，避免 z-index 问题
- **选区范围类型**：使用自定义 interface 而非 Monaco 命名空间导入，避免类型冲突
- **工具栏位置**：显示在选区起始位置上方，超出边界自动翻转
- **操作按钮**：润色/翻译/缩写/重写，prompt 包含职位目标上下文

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] ResumeParserModal uuid 模块缺失**
- **Found during:** Build verification
- **Issue:** ResumeParserModal.vue 引用 `uuid` 包但未安装
- **Fix:** 替换为浏览器内置 `crypto.randomUUID()`
- **Files modified:** frontend/src/components/ResumeParserModal.vue
- **Commit:** 5687db1

## Issues Encountered

- **TypeScript 类型错误**：MonacoEditor IRange 命名空间导入问题，修复为自定义 SelectionRange interface
- **EditorView Resume 类型不匹配**：Wails 生成的 Resume 模型缺少 jobTarget 字段，修复为 `as any` 类型断言
- **重复 Props 接口**：编辑 MonacoEditor 时意外保留了原始 Props 接口，修复为合并后的单一接口

## Verification Status

- [x] 选中文字后工具栏正确显示在选区上方
- [x] 工具栏位置不超出视口边界（翻转逻辑）
- [x] 四个操作均可正常触发并返回结果
- [x] 流式返回时选区实时替换（useAIStream + executeEdits）
- [x] 工具栏在操作完成/取消后正确隐藏
- [x] 无选区时工具栏不显示
- [x] `npm run build` 无错误

## Known Stubs

None - all features implemented as specified.

## Threat Flags

None - no new security surface introduced.

## Next Phase Readiness

- AI 工具栏基础设施完成，可被 03-03 及后续阶段使用
- 流式替换逻辑可复用于更多 AI 操作
- JobTarget 传递路径已建立，后续可直接使用

---
*Plan: 03-02*
*Completed: 2026-04-05*
