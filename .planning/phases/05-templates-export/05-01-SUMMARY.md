---
phase: "05-templates-export"
plan: "01"
subsystem: "template-system"
tags: ["template", "css", "vue", "wails"]
dependency_graph:
  requires: []
  provides:
    - "frontend/src/styles/templates/*.css"
    - "frontend/src/composables/useTemplate.ts"
    - "frontend/src/components/TemplateSelector.vue"
    - "A4Page template injection"
  affects:
    - "EditorView.vue (future integration)"
tech_stack:
  added: []
  patterns: ["CSS class-based template switching", "Wails bridge for template persistence"]
key_files:
  created:
    - "frontend/src/styles/templates/template-minimal.css"
    - "frontend/src/styles/templates/template-dual-col.css"
    - "frontend/src/styles/templates/template-academic.css"
    - "frontend/src/styles/templates/template-campus.css"
    - "frontend/src/composables/useTemplate.ts"
    - "frontend/src/components/TemplateSelector.vue"
  modified:
    - "frontend/src/components/A4Page.vue"
    - "internal/service/resume.go"
    - "app.go"
    - "frontend/src/wailsjs/wailsjs/go/main/App.d.ts"
    - "frontend/src/wailsjs/wailsjs/go/main/App.js"
decisions:
  - "使用 CSS 类名前缀区分模板 (template-minimal/dual-col/academic/campus)"
  - "模板切换仅改变 CSS 类名，不重新渲染 markdown-it"
  - "自定义 CSS 通过动态 <style> 标签注入到 A4Page 容器"
  - "WailsJS 生成的 bindings 用于前端调用后端 bridge 方法"
metrics:
  duration: "469s (~7.8 min)"
  completed: "2026-04-06T02:33:49Z"
  tasks_completed: 5
  commits: 8
---

# Phase 05 Plan 01: 模板系统 Summary

## 执行状态

**Plan:** 05-01
**Status:** 完成
**Duration:** 469s (~7.8 分钟)
**Commits:** 8

## 一句话总结

实现了4套内置模板 CSS（极简通用/双栏简约/学术科研/大厂校招）、模板选择器组件、useTemplate 状态管理 composable、A4Page 模板 CSS 注入支持，以及后端 UpdateResumeTemplate/UpdateResumeCustomCSS bridge 方法。

## 完成的任务

| Task | 名称 | Commit | 状态 |
|------|------|--------|------|
| 1 | 创建4套内置模板 CSS | `5d7246f` | 完成 |
| 2 | 创建 useTemplate.ts composable | `dccbd16` -> `50677f7` | 完成 |
| 3 | 创建 TemplateSelector.vue 组件 | `8db3026` | 完成 |
| 4 | 修改 A4Page.vue 支持模板注入 | `406ed6b` | 完成 |
| 5 | 后端新增模板 bridge 方法 | `e5bc981` -> `5fe16f1` | 完成 |

## 详细交付物

### Task 1: 4套内置模板 CSS 文件

| 文件 | 行数 | 说明 |
|------|------|------|
| `template-minimal.css` | 136行 | 极简通用风，18pt 加粗下划线标题，flex space-between 联系人行 |
| `template-dual-col.css` | 148行 | 双栏简约风，`column-count: 2`，h1/h2 全宽独占 |
| `template-academic.css` | 150行 | 学术科研风，行高 1.8，居中标题，左侧蓝色色条教育背景 |
| `template-campus.css` | 158行 | 大厂校招风，700 加粗标题 + letter-spacing，左侧 3pt 黑色竖条 |

所有选择器使用 `.template-{name} .page-content` 前缀，复用 `:root` CSS 变量。

### Task 2: useTemplate.ts Composable

- `TemplateDef` 接口：`id`, `name`, `nameEn`, `category`, `tag`
- `BUILTIN_TEMPLATES` 常量数组：4套内置模板定义
- `useTemplate(resumeId)` 返回：`currentTemplateId`, `customCss`, `isLoading`, `loadTemplate`, `selectTemplate`, `saveAsTemplate`, `currentTemplateCSS`, `currentTemplateClass`, `templates`
- 使用 WailsJS 生成的 `GetResume` / `UpdateResumeTemplate` / `UpdateResumeCustomCSS` bindings

### Task 3: TemplateSelector.vue 组件

- 网格布局（大屏4列/中屏2列/小屏2列），卡片 140px 宽
- 4套模板的 CSS 微型预览（div 模拟简历外观）
- 当前选中模板：蓝色边框 + 浅蓝背景
- hover 效果：轻微放大 + 阴影
- props: `modelValue`（当前模板ID），emits: `update:modelValue`, `change`

### Task 4: A4Page.vue 模板 CSS 注入

- 新增 props: `templateId`, `customCss`
- 静态导入所有4套模板 CSS
- 动态 `templateClass` 计算属性：`template-minimal` / `template-dual-col` / `template-academic` / `template-campus`
- 动态注入自定义 CSS 到容器（`<style id="a4page-template-dynamic-css">`）
- `onUnmounted` 清理动态 style 标签
- 暴露 `getHTMLForExport()` 方法

### Task 5: 后端 Bridge 方法

- `ResumeService` 接口新增 `UpdateTemplateID` 和 `UpdateCustomCSS`
- `UpdateTemplateID`: 验证 templateId 为有效值（minimal/dual-col/academic/campus）
- `UpdateCustomCSS`: 直接存储 customCss 字符串
- `app.go` 新增 `UpdateResumeTemplate()` 和 `UpdateResumeCustomCSS()`
- Wails bindings 重新生成，包含新方法

## 验证结果

- TypeScript: `vue-tsc --noEmit` 无错误
- Go: `go build ./...` 编译通过
- 4套模板 CSS 文件均存在，每套 >25 行

## 偏差记录

**无偏差** - 所有任务按计划完成。

## 后续计划

- **05-02 (样式可视化调整):** 实现 StyleEditor.vue 组件，滑块调整字号/行高/主色调，实时预览效果
- **05-03 (PDF导出):** 实现 ExportDialog.vue，`window.print()` + `@media print` CSS 规则
- **05-04 (版本快照):** 实现快照系统（创建/列表/Diff/回滚）

## Self-Check

- [x] 4套模板 CSS 文件存在且内容完整
- [x] useTemplate composable 导出完整的状态和方法
- [x] TemplateSelector.vue 展示4套模板卡片，可点击切换
- [x] A4Page.vue 接收 templateId 和 customCss props
- [x] 后端 UpdateResumeTemplate 和 UpdateResumeCustomCSS 方法编译通过
- [x] 模板切换仅改变 CSS 样式（通过类名）
- [x] 个人模板（customCss）切换模板时保留
- [x] Wails bindings 重新生成并提交
