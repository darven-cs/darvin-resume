---
phase: 04-resume-mgmt
plan: "04-03"
subsystem: ui
tags: [vue3, typescript, auto-save, recycle-bin, wails]

# Dependency graph
requires:
  - phase: 04-resume-mgmt
    provides: 04-01 (resume list CRUD), 04-02 (wizard integration)
provides:
  - Blank page resume creation with template injection
  - Collapsible recycle bin with restore/permanent-delete
  - Auto-save composable with 30s interval and 4-state indicator
  - Editor navigation (back button, inline title edit, route leave guard)
affects:
  - Phase 05 (template switching, PDF export)
  - Phase 06 (API key encryption)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Auto-save composable pattern (useAutoSave.ts)
    - Collapsible section with Transition animation
    - Inline title editing with keyboard shortcuts
    - Route leave guard with dirty checking

key-files:
  created:
    - frontend/src/composables/useAutoSave.ts (auto-save composable)
    - frontend/src/components/SaveStatusIndicator.vue (4-state save indicator)
    - frontend/src/components/RecycleBinSection.vue (recycle bin section)
  modified:
    - frontend/src/views/EditorView.vue (extract save logic, add nav/title/guard)
    - frontend/src/views/HomeView.vue (integrate RecycleBinSection, blank template)

key-decisions:
  - "useAutoSave uses setInterval (not setTimeout) for 30s auto-save per D-26"
  - "RecycleBinSection manages its own data fetching, emits events to parent for list refresh"
  - "Title editing uses input vs button toggle (v-if) for clean UX"
  - "Route leave guard checks saveStatus === 'unsaved' (not 'saving' or 'error')"

requirements-completed: [RESM-05, RESM-07, RESM-08]

# Metrics
duration: 7min
completed: 2026-04-06
---

# Phase 04 Plan 03: 空白页创建、回收站、自动保存 Summary

**空白页创建、回收站折叠区、自动保存 composable + 状态指示器实现完成，EditorView 重构为使用自动保存 composable 并添加导航和标题编辑功能**

## Performance

- **Duration:** 7 min
- **Started:** 2026-04-06T16:43:15Z
- **Completed:** 2026-04-06T16:50:00Z
- **Tasks:** 5 (3 new files + 2 modified views)
- **Files modified:** 5

## Accomplishments

- 实现空白页创建：CreateResume + UpdateResume 注入空白 Markdown 模板，跳转编辑器
- 实现回收站折叠区 RecycleBinSection.vue：展开/折叠动画、恢复/永久删除、列表自动刷新
- 实现 useAutoSave.ts composable：30 秒 setInterval 定时器、isDirty 标记、四种保存状态
- 实现 SaveStatusIndicator.vue：四种状态（已保存/保存中/未保存/保存失败）+ 重试按钮
- EditorView.vue 重构：替换内联保存逻辑为 useAutoSave composable
- EditorView.vue 导航增强：返回按钮、inline 标题编辑（回车确认/ESC 取消）、路由离开守卫

## Task Commits

All tasks committed atomically:

1. **feat(04-03): implement blank page creation, recycle bin, and auto-save** - `ef631ff` (feat)
   - frontend/src/composables/useAutoSave.ts (new)
   - frontend/src/components/SaveStatusIndicator.vue (new)
   - frontend/src/components/RecycleBinSection.vue (new)
   - frontend/src/views/EditorView.vue (modified)
   - frontend/src/views/HomeView.vue (modified)

## Files Created/Modified

- `frontend/src/composables/useAutoSave.ts` - 自动保存 composable，管理保存状态、定时器、isDirty 标记
- `frontend/src/components/SaveStatusIndicator.vue` - 保存状态指示器，四种状态 + 过渡动画 + 重试按钮
- `frontend/src/components/RecycleBinSection.vue` - 回收站折叠区，支持恢复/永久删除、展开动画
- `frontend/src/views/EditorView.vue` - 重构：提取 useAutoSave、添加返回按钮、标题编辑、离开守卫
- `frontend/src/views/HomeView.vue` - 集成 RecycleBinSection、空白页创建注入模板内容

## Decisions Made

- **useAutoSave 使用 setInterval 而非 setTimeout**：30 秒间隔通过 setInterval 实现，isDirty 标记控制是否真正保存
- **RecycleBinSection 自管理数据获取**：展开时懒加载数据，通过 emit 事件通知父组件刷新主列表
- **标题编辑使用 v-if 切换 input/button**：点击标题变为 input，回车确认或 ESC 取消
- **路由守卫检查 saveStatus === 'unsaved'**：仅在有变更未保存时弹出确认框，saving/error 状态不阻塞离开

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## Self-Check

- [x] useAutoSave.ts created with correct interface
- [x] SaveStatusIndicator.vue created with 4 states
- [x] RecycleBinSection.vue created with expand/collapse animation
- [x] HomeView.vue imports and uses RecycleBinSection
- [x] HomeView.vue blank page creation injects template via UpdateResume
- [x] EditorView.vue imports useAutoSave and SaveStatusIndicator
- [x] EditorView.vue has back button with unsaved check
- [x] EditorView.vue has inline title editing with RenameResume
- [x] EditorView.vue has onBeforeRouteLeave guard
- [x] Build passes (vue-tsc --noEmit + vite build)

## Next Phase Readiness

- useAutoSave composable 可被 Phase 5 的 PDF 导出前自动保存复用
- RecycleBinSection 自管理数据加载，无需父组件传递数据
- EditorView 已具备完整的简历管理闭环（创建/编辑/保存/删除/恢复）

---
*Phase: 04-resume-mgmt / 04-03*
*Completed: 2026-04-06*
