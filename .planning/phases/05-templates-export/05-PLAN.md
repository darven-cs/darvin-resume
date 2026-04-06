# Plan 05: 模板与导出

**Phase:** 5 of 7
**Goal:** 用户可以选择模板、调整样式并导出排版精确的PDF，同时可以管理简历的版本历史
**Depends on:** Phase 4 complete
**Requirements:** TMPL-01, TMPL-02, TMPL-03, TMPL-04, TMPL-05, TMPL-06, EXPT-01, EXPT-02, EXPT-03, EXPT-04, EXPT-05, EXPT-06, EXPT-07, EXPT-08, EXPT-09

## Success Criteria (what must be TRUE)

1. 用户可以在4套内置模板（极简通用/大厂校招/学术科研/双栏简约）间自由切换，切换后简历内容完整保留，仅渲染样式变化
2. 用户可以通过滑块/颜色选择器/字体下拉调整主色调、字号、行高、边距、字体，调整后实时在预览区看到效果
3. 用户可以输入自定义 CSS（白名单安全机制过滤），可以一键重置默认样式，可以将当前样式保存为个人模板
4. 用户可以导出 PDF（默认系统打印，可选 Chromedp 高级模式），导出严格遵循 A4 尺寸和分页规则，与预览 100% 一致
5. 用户可以手动创建版本快照（标签+备注），PDF导出后自动创建快照，查看版本历史并对比两版本差异，一键回滚（回滚前自动快照当前）
6. 快照存储完整数据（JSON+Markdown+模板+CSS），快照上限50个自动清理最旧的

## Plans

| Plan | Name | Status | Dependencies |
|------|------|--------|-------------|
| 05-01 | 模板系统 — 4套内置模板 + 切换机制 + 个人模板保存 | Planned | — |
| 05-02 | 样式调整 — 可视化调整 + CSS白名单 + 一键重置 | Planned | — |
| 05-03 | PDF导出 — 系统打印 + Chromedp备选 + A4分页规则 + 导出参数 | Planned | 05-01 |
| 05-04 | 版本快照 — 创建/自动创建 + 历史列表 + Diff对比 + 回滚 | Planned | 05-03 |

## Decision Coverage Matrix

| Requirement | Plan | Coverage | Notes |
|-------------|------|----------|-------|
| TMPL-01 | 05-01 | Full | 4套内置模板 CSS |
| TMPL-02 | 05-01 | Full | 仅切换 CSS，数据复用 |
| TMPL-03 | 05-01 | Full | custom_css 列存储 |
| TMPL-04 | 05-02 | Full | CSS 变量 + 滑块/选择器 |
| TMPL-05 | 05-02 | Full | DOMPurify + PostCSS 白名单 |
| TMPL-06 | 05-02 | Full | 恢复默认样式逻辑 |
| EXPT-01 | 05-03 | Full | window.print() 系统打印 |
| EXPT-02 | 05-03 | Full | Chromedp 可选高级模式 |
| EXPT-03 | 05-03 | Full | @media print + break-inside: avoid |
| EXPT-04 | 05-03 | Full | 页码范围/分页线/DPI 参数 |
| EXPT-05 | 05-04 | Full | 手动创建快照 |
| EXPT-06 | 05-04 | Full | 导出后自动快照 hook |
| EXPT-07 | 05-04 | Full | 快照列表 + diff 库对比 |
| EXPT-08 | 05-04 | Full | 一键回滚（回滚前自动快照） |
| EXPT-09 | 05-04 | Full | 完整数据存储 + 50个上限 |

## New Dependencies

### npm (frontend)
```
dompurify: ^3.3.3       # CSS 白名单校验
postcss: ^8.5.8          # CSS AST 解析
postcss-selector-parser: ^7.1.3  # CSS 选择器校验
```

### go.mod (backend)
```
github.com/chromedp/chromedp  # 可选高级导出（默认不导入，仅条件编译）
```

## New Files

### Backend (Go)
```
internal/
├── model/
│   └── snapshot.go              # Snapshot model 定义
├── service/
│   └── snapshot.go              # SnapshotService 接口 + 实现
└── database/
    └── migrations/
        └── 004_create_snapshots_table.sql
app.go                             # 新增 Snapshot 相关 bridge 方法
```

### Frontend (Vue/TS)
```
frontend/src/
├── styles/templates/
│   ├── template-minimal.css
│   ├── template-dual-col.css
│   ├── template-academic.css
│   └── template-campus.css
├── components/
│   ├── TemplateSelector.vue
│   ├── StyleEditor.vue
│   ├── ExportDialog.vue
│   └── SnapshotSidebar.vue
├── composables/
│   ├── useTemplate.ts
│   └── useSnapshot.ts
├── utils/
│   └── sanitizeCSS.ts
└── styles/
    └── editor.css (追加 @media print)
```

## Execution Order

```
Wave 1 (并行):
  05-01: 模板系统
    - 4套内置模板 CSS 文件
    - TemplateSelector 组件
    - useTemplate composable
    - A4Page.vue 模板注入
    - 后端模板更新 bridge

Wave 2 (依赖 05-01):
  05-02: 样式调整
    - StyleEditor 组件（滑块/选择器）
    - sanitizeCSS 工具函数
    - 一键重置逻辑
    - useTemplate 扩展

Wave 3 (依赖 05-01):
  05-03: PDF导出
    - @media print CSS 规则
    - ExportDialog 组件
    - Chromedp 后端服务
    - 多页分页 JS 逻辑（保守版本）

Wave 4 (依赖 05-03):
  05-04: 版本快照
    - snapshots 表迁移
    - SnapshotService
    - SnapshotSidebar 组件
    - Diff 对比功能
    - 回滚逻辑 + 导出自动快照 hook
```

**并行策略：**
- 05-01 是所有其他 plan 的基础（提供模板选择、模板 CSS 注入、useTemplate composable），需最先执行
- 05-02 和 05-03 均可依赖 05-01，在 05-01 完成后并行执行
- 05-04 依赖 05-03（因为包含导出后自动快照的 hook），串行在 05-03 之后

## Wave Analysis

| Wave | Plans | Parallelizable | Key Dependencies |
|------|-------|---------------|-----------------|
| 1 | 05-01 | — | Phase 4 完成 |
| 2 | 05-02, 05-03 | Yes（无文件冲突） | 05-01 |
| 3 | 05-04 | — | 05-03 |
