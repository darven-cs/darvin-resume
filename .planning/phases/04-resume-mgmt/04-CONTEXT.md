# Phase 4: 简历创建与管理 - Context

**Gathered:** 2026-04-05
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 4 delivers resume creation workflows (AI-guided wizard + blank page), resume list management with card grid, recycle bin, and auto-save mechanism. This phase transforms the app from "editor-only" to a complete resume management tool.

**Out of scope:** Template switching (Phase 5), PDF export (Phase 5), dark/light themes (Phase 7), API Key encryption (Phase 6)
</domain>

<decisions>
## Implementation Decisions

### Resume Creation Entry Point (RESM-05)
- **D-01:** 点击 HomeView 的「新建简历」按钮后弹出选择弹窗，提供「AI引导创建」和「空白页创建」两个选项
- **D-02:** 空白页创建：直接进入编辑器，生成空白 Markdown 模板（标题 + 基础模块占位符），用户自由编辑
- **D-03:** AI引导创建：进入编辑器后自动打开侧边栏向导，引导完成后关闭向导进入正常编辑

### AI Guided Creation - Sidebar Wizard (RESM-01, RESM-02, RESM-03, RESM-04)
- **D-04:** 向导在编辑器右侧侧边栏（复用 AIChatSidebar 的布局模式）中执行，不离开编辑器页面
- **D-05:** 流程分三步：①勾选模块 → ②按勾选顺序逐步填写 → ③生成完整简历
- **D-06:** 第一步：模块勾选界面，展示可选模块（基础信息/专业技能/项目经历/自我评价默认勾选，校园经历/实习经历/获奖/证书可选），用户调整后点击「开始填写」
- **D-07:** 第二步：逐步填写，每个模块展示填写表单（含字段标签 + 输入框 + 示例提示），底部有「AI润色」按钮和「下一步」按钮
- **D-08:** AI润色为手动触发 — 用户填写完一个模块后手动点击「AI润色」，润色结果在表单下方以对比方式展示，用户确认后更新填写内容
- **D-09:** 每个模块支持：跳过（保留空模板）、重新润色、修改润色结果
- **D-10:** 全部模块填写完成后，点击「生成简历」，系统将所有模块数据组合生成完整 Markdown 内容填入编辑器
- **D-11:** 中途退出向导时弹出确认框，选择保存为草稿 — 已填写内容生成 Markdown 保存为简历，用户下次可在编辑器中继续完善
- **D-12:** 向导组件命名为 `ResumeWizardSidebar.vue`，通过 prop 控制显示/隐藏

### HomeView - Resume List (RESM-06)
- **D-13:** HomeView 重构为卡片网格布局，响应式排列（大屏3-4列，中屏2列，小屏1列）
- **D-14:** 每张简历卡片展示：标题、最近修改时间、Markdown 内容渲染预览缩略图
- **D-15:** 卡片右上角显示快捷操作按钮（重命名/复制/删除），hover 或点击更多图标时展开
- **D-16:** 点击卡片主体区域进入编辑器（路由跳转 `/editor/:id`）
- **D-17:** 列表页顶部包含：搜索框 + 新建按钮，支持按修改时间排序
- **D-18:** 空状态：无简历时展示引导文案 + 醒目新建按钮 + 简历模板Demo预览（Phase 5 完善，本阶段展示占位）

### Recycle Bin (RESM-07)
- **D-19:** 回收站作为 HomeView 底部折叠区，点击「回收站」标题展开/折叠
- **D-20:** 展开后显示已软删除的简历列表（紧凑列表形式，非卡片）
- **D-21:** 每条记录显示：标题 + 删除时间 + 「恢复」按钮 + 「永久删除」按钮
- **D-22:** 永久删除需二次确认弹窗
- **D-23:** 回收站无30天自动清理逻辑 — 用户手动管理（降低复杂度）
- **D-24:** 后端需新增方法：`RestoreResume(ctx, id)` 恢复软删除、`PermanentDeleteResume(ctx, id)` 永久删除、`ListDeletedResumes(ctx)` 列出已删除简历、`DuplicateResume(ctx, id)` 复制简历

### Auto-Save (RESM-08)
- **D-25:** 编辑器标题栏右侧显示保存状态文字：「已保存」(绿色) / 「保存中...」(灰色) / 「未保存」(橙色)
- **D-26:** 自动保存触发条件：内容变更后 30 秒间隔 + AI操作完成时 + 页面切换/离开时
- **D-27:** 将自动保存逻辑从 EditorView 内联抽取为独立 composable `useAutoSave.ts`
- **D-28:** 保存失败时显示错误状态 + 重试按钮，不阻塞用户继续编辑

### Navigation
- **D-29:** EditorView 顶部添加返回按钮，点击返回 HomeView（路由 `/`）
- **D-30:** EditorView 标题栏支持点击编辑简历标题（替代独立重命名流程）
- **D-31:** 路由离开守卫：编辑器有未保存内容时弹窗确认是否离开

### New Backend Methods
- **D-32:** `RenameResume(ctx, id, title)` — 修改简历标题（或复用 UpdateResume，只传 title 字段）
- **D-33:** `DuplicateResume(ctx, id)` — 复制简历（新记录 + "(副本)" 后缀）
- **D-34:** `RestoreResume(ctx, id)` — 恢复软删除（is_deleted = false, deleted_at = null）
- **D-35:** `PermanentDeleteResume(ctx, id)` — 物理删除记录
- **D-36:** `ListDeletedResumes(ctx)` — 查询 is_deleted = true 的记录

### Claude's Discretion
- 卡片网格具体 CSS 样式和间距
- 侧边栏向导的动画和过渡效果
- 搜索算法（前端过滤 vs 后端查询）
- 草稿保存的具体 Markdown 模板格式
- 预览缩略图的截取方式
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning and implementing.**

### Project Specifications
- `.planning/REQUIREMENTS.md` — RESM-01 through RESM-08
- `.planning/PROJECT.md` — Architecture decisions, tech stack constraints
- `.planning/ROADMAP.md` — Phase 4 success criteria and dependencies

### Prior Phase Context
- `.planning/phases/03-ai-core/03-CONTEXT.md` — AI infrastructure (streaming, error handling, config)
- `.planning/phases/02-core-editor/02-CONTEXT.md` — Editor layout and rendering decisions

### Code References
- `frontend/src/views/HomeView.vue` — Current skeleton, needs full rewrite
- `frontend/src/views/EditorView.vue` — Editor with inline auto-save, needs refactor + nav
- `frontend/src/components/ResumeParserModal.vue` — AI import flow, wizard should follow similar patterns
- `frontend/src/components/AIChatSidebar.vue` — Sidebar layout pattern to reuse for wizard
- `frontend/src/types/resume.ts` — Resume, ResumeListItem, Module types
- `internal/service/resume.go` — ResumeService interface, needs new methods
- `internal/database/db.go` — SQLite queries, soft-delete is_deleted column exists
- `app.go` — Bridge methods, needs new method bindings
- `frontend/src/composables/useDebounce.ts` — Reusable debounce for auto-save composable

### Technical Documentation
- Vue Router: https://router.vuejs.org/
- Wails Runtime: https://wails.io/docs/reference/runtime/
</canonical_refs>

<file_plan>
## Phase 4 File Plan

### Backend (Go)
```
internal/
├── service/
│   └── resume.go          # Add: RenameResume, DuplicateResume, RestoreResume, PermanentDeleteResume, ListDeletedResumes
└── database/
    └── db.go              # Add: SQL queries for new service methods
app.go                     # Add: Bridge bindings for new methods
```

### Frontend (Vue/TS)
```
frontend/src/
├── components/
│   ├── CreateModeModal.vue        # 创建模式选择弹窗（AI引导/空白页）
│   ├── ResumeWizardSidebar.vue    # AI引导式创建侧边栏向导
│   ├── ResumeCard.vue             # 简历卡片组件
│   ├── RecycleBinSection.vue      # 回收站折叠区
│   └── SaveStatusIndicator.vue    # 保存状态指示器
├── composables/
│   ├── useAutoSave.ts             # 自动保存 composable（从 EditorView 抽取）
│   └── useResumeList.ts           # 简历列表状态管理
├── views/
│   ├── HomeView.vue               # 重构：卡片网格 + 搜索 + 回收站
│   └── EditorView.vue             # 重构：抽取自动保存 + 添加导航 + 标题编辑
└── router/
    └── index.ts                   # 可能需要更新路由守卫
```
</file_plan>

<deferred>
## Deferred Ideas

**Deferred to Phase 5 (Templates & Export):**
- 简历模板 Demo 预览（空状态展示）
- 模板选择作为创建流程的一部分

**Deferred to Phase 6 (Security):**
- 自动备份与简历草稿关联

**Deferred to Phase 7 (UI Polish):**
- 卡片拖拽排序
- 快捷键支持（创建/删除）
- 简历缩略图缓存优化

**Deferred to Future Versions:**
- 30天自动清理回收站
- 简历标签/分类管理
- 批量操作（批量删除/导出）
</deferred>
