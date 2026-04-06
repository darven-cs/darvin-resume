# Phase 5: 模板与导出 - Context

**Gathered:** 2026-04-06
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 5 实现模板系统（4套内置模板 + 个人模板）、样式可视化调整（滑块/选择器）、自定义CSS白名单安全机制、PDF导出（系统打印+Chromedp备选）、版本快照管理（创建/自动创建/历史/Diff/回滚）。

**Upstream dependency:** Phase 4 完成 — 简历列表、空白页创建、自动保存机制均已就绪。
**Downstream dependency:** Phase 6 的 API Key 加密、备份恢复依赖模板和快照系统。

**Out of scope:** 深色/浅色主题（Phase 7）、API Key 加密存储（Phase 6）、数据备份包导入导出（Phase 6）、ATS 优化等高级功能。
</domain>

<decisions>
## Implementation Decisions

### 模板系统架构（TMPL-01, TMPL-02, TMPL-03）

- **D-01:** 4套内置模板以 CSS 类名前缀区分，存储在 `frontend/src/styles/templates/` 目录下的独立 `.css` 文件中
  - `template-minimal.css` — 极简通用风，白底黑字单栏
  - `template-dual-col.css` — 双栏简约风，CSS `column-count: 2`
  - `template-academic.css` — 学术科研风，宽松行距，教育背景突出
  - `template-campus.css` — 大厂校招风，大标题活力配色

- **D-02:** 模板切换机制：在 `A4Page.vue` 中通过动态注入 `<style>` 标签实现模板 CSS 切换，模板 ID（如 `minimal`）作为类名添加到 `.page-content` 的父容器
- **D-03:** 模板切换仅更新 CSS，不重新执行 markdown-it 渲染，内容完全复用（TMPL-02）
- **D-04:** 个人模板保存：将当前 `custom_css` 内容存储到 `resumes.custom_css` 列，模板切换时保留
- **D-05:** 模板数据结构：使用已有的 `resumes.template_id` 和 `resumes.custom_css` 列，**不新建 `templates` 表**（Phase 5 仅实现模板切换逻辑，个人模板保存依赖 custom_css 列）
- **D-06:** 内置模板 ID：`minimal`（默认）、`dual-col`、`academic`、`campus`

### 样式可视化调整（TMPL-04）

- **D-07:** 通过 CSS 变量实现样式调整：`:root` 中定义 `--resume-*` 变量（主色调/字号/行高/边距/字体），滑块/选择器修改这些变量值
- **D-08:** 样式调整 UI 作为编辑器工具栏或侧边栏面板，通过 `useTemplate.ts` composable 管理当前模板 ID、样式参数、自定义 CSS
- **D-09:** 实时预览：CSS 变量变更通过 Vue reactive state 反映到 DOM，无需重新渲染 markdown-it

### 自定义CSS白名单（TMPL-05）

- **D-10:** 前端使用 DOMPurify + PostCSS 做 CSS 白名单校验，30+ 安全属性白名单（字体/颜色/边距/行高/列表样式）
- **D-11:** 危险属性黑名单：`position`、`z-index`、`url(`、`expression(`、`@import`、`animation`、`transition`
- **D-12:** 用户输入的 CSS 经过白名单校验后存储到 `resumes.custom_css` 列，渲染时注入 `<style>` 标签
- **D-13:** 新增依赖：`dompurify`、`postcss`、`postcss-selector-parser`（npm）

### PDF导出（EXPT-01, EXPT-02, EXPT-03, EXPT-04）

- **D-14:** 默认方案：`window.print()` 通过 Wails WebView 触发系统原生打印对话框，用户选择「另存为 PDF」
- **D-15:** 高级可选方案：Chromedp 无头浏览器渲染，`go get github.com/chromedp/chromedp`（默认不导入，仅在用户开启高级模式时使用）
- **D-16:** `@media print` CSS 规则追加到 `editor.css`：隐藏 UI 元素、设置 A4 尺寸、启用 `break-inside: avoid`、强制背景色打印
- **D-17:** 导出参数：页码范围（默认全部）、隐藏/显示分页线（默认隐藏）、DPI（默认 96，高级模式可调）
- **D-18:** 多页分页：Phase 5 采用保守策略 — 按 DOM 高度计算分页，单页内容正常渲染，溢出时提示用户调整字号或内容（完整智能分页作为后续迭代）
- **D-19:** Chromedp PDF 生成：通过 `page.PrintToPDF()` 实现，A4 尺寸（210mm x 297mm），`print-color-adjust: exact`，`break-inside: avoid`

### 版本快照管理（EXPT-05, EXPT-06, EXPT-07, EXPT-08, EXPT-09）

- **D-20:** 快照存储到 `snapshots` 表（含 `id/resume_id/label/note/trigger_type/json_data/markdown_content/template_id/custom_css/created_at`）
- **D-21:** 手动创建快照（EXPT-05）：用户点击「创建快照」按钮，输入标签+备注，调用后端 `CreateSnapshot()`
- **D-22:** PDF导出后自动创建快照（EXPT-06）：导出成功时 hook 自动创建，`trigger_type = 'auto_pdf_export'`
- **D-23:** 快照列表（EXPT-07）：侧边栏面板展示历史快照列表，每条显示标签+触发类型+时间
- **D-24:** Diff对比（EXPT-07）：复用前端已有的 `diff` npm 包（v5.2.2），两版本 Markdown 内容对比，以 side-by-side 视图展示
- **D-25:** 一键回滚（EXPT-08）：回滚前自动创建当前快照（`trigger_type = 'rollback'`），然后恢复数据
- **D-26:** 快照上限50个（EXPT-09）：创建新快照时检查数量，超出自动删除最旧的；`DELETE FROM snapshots WHERE id IN (SELECT id FROM snapshots WHERE resume_id = ? ORDER BY created_at ASC LIMIT ?)`
- **D-27:** 快照存储完整数据（EXPT-09）：`json_data` + `markdown_content` + `template_id` + `custom_css`，快照恢复时还原所有字段

### Claude's Discretion

- 模板选择器的 UI 布局和交互细节（卡片展示还是下拉列表）
- 样式调整面板的具体滑块样式和颜色选择器实现
- 快照列表的 UI 布局和 Diff 对比视图的具体样式
- 导出参数界面的具体交互
- Chromedp 的具体浏览器启动策略（Chrome vs Chromium 回退）
- 多页分页的 DOM 高度计算具体算法（保守策略优先）
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning and implementing.**

### Project Specifications
- `.planning/REQUIREMENTS.md` — TMPL-01~06, EXPT-01~09
- `.planning/ROADMAP.md` — Phase 5 goal, requirements, success criteria
- `.planning/phase-5/05-RESEARCH.md` — 技术调研结果，包含模板系统架构、PDF导出方案、快照数据模型

### Prior Phase Context
- `.planning/phases/04-resume-mgmt/04-CONTEXT.md` — Phase 4 decisions（创建模式、列表管理、回收站、自动保存）
- `.planning/phases/02-core-editor/02-CONTEXT.md` — Editor layout and rendering decisions
- `.planning/phases/03-ai-core/03-CONTEXT.md` — AI infrastructure

### Code References
- `frontend/src/styles/editor.css` — 统一样式表（预览与PDF共用），需追加 `@media print` 规则
- `frontend/src/components/A4Page.vue` — A4预览容器，需修改为支持模板 CSS 注入
- `frontend/src/utils/markdown.ts` — markdown-it 渲染引擎（共用）
- `frontend/src/types/resume.ts` — Resume, Module, BasicInfo 类型定义
- `frontend/src/views/EditorView.vue` — 编辑器主视图，需集成模板选择器和导出功能
- `internal/database/db.go` — SQLite 连接和查询
- `internal/service/resume.go` — ResumeService 接口
- `internal/model/resume.go` — Resume Go model（含 TemplateID, CustomCSS 字段）
- `app.go` — Wails bridge 方法绑定
- `internal/database/migrations/001_create_resumes_table.sql` — template_id 和 custom_css 列已存在

### Technical Documentation
- Wails v2 Runtime: https://wails.io/docs/reference/runtime/
- Chromedp GitHub: https://github.com/chromedp/chromedp
- DOMPurify: https://github.com/cure53/DOMPurify
- PostCSS: https://postcss.org/
- diff npm: https://www.npmjs.com/package/diff

### Existing npm Packages
- `diff: ^5.2.2` — 已有，用于快照 Diff
- `markdown-it: ^14.1.1` — 已有
- `dompurify: ^3.3.3` — 需新增
- `postcss: ^8.5.8` — 需新增
- `postcss-selector-parser` — 需新增

### Existing Go Packages
- `github.com/google/uuid v1.6.0` — 已有
- `github.com/chromedp/chromedp` — 需新增（可选）
</canonical_refs>

<file_plan>
## Phase 5 File Plan

### Backend (Go)
```
internal/
├── service/
│   ├── resume.go              # 修改：UpdateTemplateID, UpdateCustomCSS
│   └── snapshot.go            # 新建：SnapshotService 全套方法
├── model/
│   └── snapshot.go            # 新建：Snapshot model
├── database/
│   └── migrations/
│       └── 004_create_snapshots_table.sql  # 新建
app.go                           # 修改：新增 Snapshot 相关 bridge 方法
```

### Frontend (Vue/TS)
```
frontend/src/
├── styles/
│   └── templates/
│       ├── template-minimal.css    # 新建
│       ├── template-dual-col.css   # 新建
│       ├── template-academic.css   # 新建
│       └── template-campus.css     # 新建
├── components/
│   ├── TemplateSelector.vue         # 新建：模板选择器（TMPL-01, TMPL-02）
│   ├── StyleEditor.vue             # 新建：样式调整面板（TMPL-04, TMPL-05, TMPL-06）
│   ├── ExportDialog.vue            # 新建：PDF导出参数对话框（EXPT-01~04）
│   ├── SnapshotSidebar.vue        # 新建：版本快照侧边栏（EXPT-05~09）
│   └── AIDiffView.vue             # 修改：复用为快照 Diff 对比（已有 diff 库）
├── composables/
│   ├── useTemplate.ts              # 新建：模板状态管理
│   └── useSnapshot.ts             # 新建：快照状态管理
├── utils/
│   └── sanitizeCSS.ts              # 新建：CSS 白名单校验（TMPL-05）
├── views/
│   └── EditorView.vue             # 修改：集成模板选择器、样式调整、导出按钮、快照侧边栏
└── styles/
    └── editor.css                 # 修改：追加 @media print 规则
```

### 前端新增 npm 依赖
```
dompurify: ^3.3.3
postcss: ^8.5.8
postcss-selector-parser: ^7.1.3
```
</file_plan>

<deferred>
## Deferred Ideas

**Deferred to Phase 6 (数据安全与备份):**
- API Key AES-256-GCM 加密存储
- 全量数据备份包导出/导入（EXPT-10, EXPT-11, EXPT-12）
- 自动备份周期设置

**Deferred to Phase 7 (界面打磨与健壮性):**
- 深色/浅色主题跟随系统切换（UIUX-02）
- 响应式布局适配，窗口宽度 <1200px 单栏切换（UIUX-04, UIUX-05）
- 自定义快捷键体系（UIUX-06, UIUX-07）
- 空状态展示模板 Demo 预览（UIUX-08）

**Deferred to Future Versions:**
- 模板 `templates` 表（存储用户自定义模板定义，而非仅 custom_css）
- 完整智能多页分页算法
- 快照导出为独立文件
- 多语言模板支持
- 简历标签/分类管理
</deferred>
