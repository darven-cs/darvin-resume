---
phase: "05-templates-export"
plan: "02"
subsystem: ui
tags: [css-variables, postcss, vue3, whitelist, template]

# Dependency graph
requires:
  - phase: "05-01"
    provides: "useTemplate composable, TemplateSelector, A4Page with templateId/customCss props, 4 template CSS files"
provides:
  - "CSS 变量系统 --resume-* 在 editor.css 中定义"
  - "sanitizeCustomCSS() 白名单校验工具（postcss）"
  - "StyleEditor.vue 组件：6 种样式控件 + 实时预览 + 一键重置"
  - "StyleEditor 集成到 EditorView 工具栏"
affects: [05-03, 05-04, PDF 导出]

# Tech tracking
tech-stack:
  added: [dompurify, postcss, postcss-selector-parser, @types/dompurify]
  patterns:
    - "CSS 变量驱动样式系统（--resume-* 变量覆盖硬编码值）"
    - "PostCSS 白名单安全过滤（属性名白名单 + 危险值正则黑名单 + 选择器前缀限制）"
    - "组件内实时预览（document.documentElement.style.setProperty）"
    - "样式变更自动保存防抖（500ms debounce）"

key-files:
  created:
    - "frontend/src/utils/sanitizeCSS.ts - CSS 白名单校验工具"
    - "frontend/src/components/StyleEditor.vue - 样式调整面板组件"
  modified:
    - "frontend/src/styles/editor.css - 新增 --resume-* CSS 变量定义"
    - "frontend/src/views/EditorView.vue - 集成 StyleEditor 和 useTemplate"

key-decisions:
  - "CSS 变量追加到现有 :root 块后，不替换原有变量，避免冲突"
  - "自定义 CSS textarea 实时校验，用 removedCount 统计被过滤行数"
  - "StyleEditor 作为预览区侧边栏，不影响编辑器宽度"
  - "样式变更通过 document.documentElement.style.setProperty 实时应用，不依赖模板文件"
  - "重置弹窗使用 position:fixed overlay + confirm/cancel 双按钮"

patterns-established:
  - "CSS 变量作为样式调整的统一入口（与模板 CSS 解耦）"
  - "sanitizeCustomCSS() 作为安全护栏，所有用户自定义 CSS 必须经过此函数"

requirements-completed: [TMPL-04, TMPL-05, TMPL-06]

# Metrics
duration: 7min
completed: 2026-04-06
---

# Phase 05 Plan 02: 样式调整功能 Summary

**可视化样式调整面板（TMPL-04）+ CSS 白名单安全机制（TMPL-05）+ 一键重置（TMPL-06），通过 CSS 变量系统实现实时预览与安全保存**

## Performance

- **Duration:** 7 min
- **Started:** 2026-04-06T02:37:03Z
- **Completed:** 2026-04-06T02:43:51Z
- **Tasks:** 3/3
- **Files modified:** 4 (3 created, 1 modified)

## Accomplishments

- 在 editor.css 中定义完整的 --resume-* CSS 变量系统（主色调、标题色、链接色、背景色、字号、行高、边距、字体）
- 创建 sanitizeCustomCSS() 工具，使用 postcss 解析 + 属性名白名单（38 项）+ 危险值正则黑名单 + 选择器前缀限制
- StyleEditor.vue 组件包含 6 种控件：主色调选择器（含 6 色预设板）、字号滑块（8-14pt）、行高滑块（1.2-2.0）、页面边距滑块（15-30mm）、字体选择器（6 种字体栈）、自定义 CSS 文本区（2000 字符限制 + 实时校验反馈）
- StyleEditor 集成到 EditorView 工具栏，预览区侧边栏显示，样式变更实时反映到预览区
- 一键重置按钮带确认弹窗，恢复所有 DEFAULT_CSS_VARS 默认值

## Task Commits

Each task was committed atomically:

1. **Task 1: CSS 变量系统** - `662ad0e` (feat)
2. **Task 2: sanitizeCSS 工具** - `6f21768` (feat)
3. **Task 3: StyleEditor 组件** - `226a3c0` (feat)

**Plan metadata commit:** [见下方 final commit]

## Files Created/Modified

- `frontend/src/utils/sanitizeCSS.ts` - PostCSS 白名单 CSS 校验工具，含 DEFAULT_CSS_VARS、sanitizeCustomCSS()、buildCSSVars() 三个导出
- `frontend/src/components/StyleEditor.vue` - 样式调整面板组件，6 种控件 + 实时预览 + 重置确认弹窗，约 600 行
- `frontend/src/styles/editor.css` - 新增 50 行 CSS 变量定义，.page-content 改用 var(--resume-*) 引用
- `frontend/src/views/EditorView.vue` - 新增 StyleEditor 工具栏按钮、useTemplate 集成、preview-content-wrapper 布局、A4Page 传入 templateId/customCss 参数

## Decisions Made

- CSS 变量追加到 :root 块之后（不替换原有 --a4-* 变量），通过 cascade 覆盖基础样式
- 样式变更不立即写入模板 CSS 文件，而是写入 customCss 字段（存储完整 CSS 字符串，含变量块 + 用户自定义部分）
- 编辑器中预览通过 `document.documentElement.style.setProperty` 实现，无需重新渲染 A4Page 组件
- 自定义 CSS 校验使用行数差作为 removedCount 的粗略估算（而非精确声明计数）

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

- 编辑器 CSS 中 JSDoc 注释含 `/* ... */` 导致 vue-tsc 解析错误（第 272 行）—— 修复：将多行注释中的 `/*` 改为普通文字描述
- postcss-selector-parser 作为依赖已存在但未使用——按计划引入了 postcss（已有）和 dompurify + @types/dompurify
- ExportDialog.vue 有预存的 `window.go` 类型错误——判定为 pre-existing 问题，不在 05-02 范围内，未修复

## Threat Surface Scan

| Flag | File | Description |
|------|------|-------------|
| sanitize_whitelist | sanitizeCSS.ts | 所有用户 CSS 经 postcss 白名单过滤，属性值黑名单拒绝 url/expression/javascript:/@import/@keyframes/@media/data: |

## Next Phase Readiness

- CSS 变量系统就绪，05-03 PDF 导出可使用相同变量确保预览与导出 100% 一致
- sanitizeCustomCSS 工具可供任何需要用户输入 CSS 的场景复用
- StyleEditor 组件已集成到 EditorView，后续可扩展更多样式控件（如段间距、列表缩进等）

---
*Phase: 05-templates-export 02*
*Completed: 2026-04-06*
