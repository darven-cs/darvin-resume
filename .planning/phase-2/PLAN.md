# Phase 2: 核心编辑器 — 执行计划

**Phase:** 2 of 7
**Goal:** 用户可以在双栏界面中流畅编写Markdown并实时看到与导出一致的预览效果
**Depends on:** Phase 1 (✅ 已完成)
**Status:** Ready for execution

## 验收标准 (Success Criteria)

> Goal-backward verification — 以下标准必须全部满足

1. **用户在Monaco编辑器中编写Markdown，编辑体验与VS Code一致**
   - 撤销/重做（Ctrl+Z / Ctrl+Y）
   - 多光标编辑（Ctrl+D）
   - 块选择、查找替换（Ctrl+F）

2. **预览区实时渲染Markdown内容，输入到渲染延迟<200ms**

3. **双栏布局支持拖拽调整宽度比例，预览区固定显示A4分页线**

4. **用户可以折叠/展开Markdown内容块，可以拖拽排序内容块，可以点击行首图标弹出快捷菜单**

## 需求覆盖

| ID | 需求 | 计划 | 状态 |
|----|------|------|------|
| EDIT-01 | Monaco Editor 核心 | 02-01 | ✅ |
| EDIT-02 | VS Code 兼容键位 | 02-01 | ✅ |
| EDIT-03 | 双栏布局 | 02-02, 02-03 | ✅ |
| EDIT-05 | A4 分页线 | 02-02 | ✅ |
| EDIT-06 | 实时预览 <200ms | 02-03 | ✅ |
| EDIT-08 | 折叠/展开图标 | 02-04 | ✅ |
| EDIT-09 | 拖拽排序 | 02-04 | ✅ |
| EDIT-10 | 快捷菜单 | 02-04 | ✅ |

## 执行计划

### Wave 1: Monaco Editor 基础集成 (02-01)

**目标:** Monaco Editor 集成、Markdown 语法高亮、VS Code 兼容键位配置

**文件变更:**
- `frontend/package.json` — 添加 @monaco-editor/loader, monaco-editor 依赖
- `frontend/vite.config.ts` — 配置 vite-plugin-monaco-editor
- `frontend/src/components/MonacoEditor.vue` — Monaco Editor Vue 3 封装组件

**子任务:**
1. 安装 npm 包并配置 Vite
2. 创建 MonacoEditor.vue 组件

**验收:** Monaco Editor 在 Vue 组件中正常加载，TypeScript 编译无错误

---

### Wave 2: 双栏布局 + A4 分页线 (02-02)

**目标:** SplitPane 双栏布局组件、拖拽调整栏宽、A4 页面边界线

**文件变更:**
- `frontend/src/components/SplitPane.vue` — 双栏分割组件
- `frontend/src/components/A4Page.vue` — A4 分页线组件
- `frontend/src/styles/editor.css` — 编辑器样式表

**子任务:**
1. 创建 SplitPane.vue 组件（拖拽分隔条、50:50 默认比例、300px 最小宽度）
2. 创建 A4Page.vue 组件（210mm×297mm A4 尺寸、页面边界线）

**验收:** 双栏布局可拖拽调整，A4 分页线正确显示

**依赖:** 02-01 (MonacoEditor.vue)

---

### Wave 3: 实时预览同步 (02-03)

**目标:** 实时预览同步(<200ms)、编辑器预览滚动同步、响应式单栏切换

**文件变更:**
- `frontend/src/views/EditorView.vue` — 编辑器主视图
- `frontend/src/composables/useDebounce.ts` — 防抖工具函数
- `frontend/src/utils/markdown.ts` — 渲染工具扩展

**子任务:**
1. 创建 useDebounce.ts 防抖工具
2. 整合 EditorView.vue（Monaco + SplitPane + A4Page）
3. 实现 150ms 防抖预览更新
4. 实现编辑器预览滚动同步
5. 实现 <1200px 自动单栏切换

**验收:** 输入到渲染延迟 <200ms，滚动同步正常，响应式切换正常

**依赖:** 02-01, 02-02

---

### Wave 4: 行级交互 (02-04)

**目标:** 折叠/展开、拖拽排序、快捷菜单

**文件变更:**
- `frontend/src/components/MonacoEditor.vue` — 扩展行级交互功能
- `frontend/src/components/ContextMenu.vue` — 上下文菜单组件

**子任务:**
1. 实现折叠/展开图标（Gutter 装饰、Monaco Decorations API）
2. 实现拖拽排序（HTML5 Drag API）和快捷菜单

**验收:** 可折叠/展开内容块、可拖拽排序、快捷菜单正常弹出

**依赖:** 02-01

---

## 决策覆盖

| D-ID | 决策 | 计划 |
|------|------|------|
| D-01 | Monaco via npm (不依赖 CDN) | 02-01 |
| D-02 | VS Code Dark 主题默认 | 02-01 |
| D-03 | 字体 14px，行高 1.6 | 02-01 |
| D-04 | 中文语言支持 | 02-01 |
| D-05 | VS Code 所有编辑功能 | 02-01 |
| D-06 | 50:50 默认比例 | 02-02 |
| D-07 | 每栏最小 300px | 02-02 |
| D-08 | 拖拽手柄视觉反馈 | 02-02 |
| D-09 | <1200px 自动单栏 | 02-02, 02-03 |
| D-10 | 150ms 防抖 | 02-03 |
| D-11 | 滚动同步 | 02-03 |
| D-12 | 渲染错误处理 | 02-03 |
| D-13 | 折叠图标 | 02-04 |
| D-14 | HTML5 拖拽 | 02-04 |
| D-15 | 快捷菜单 | 02-04 |
| D-16 | 10k 行性能 | 02-04 |
| D-17 | A4 边界线 | 02-02 |
| D-18 | 分页虚线 | 02-02 |

## 延期项目 (Deferred)

- **EDIT-07** (选区 AI 操作) → Phase 3
- **EDIT-11** (Diff 对比) → Phase 3
- **AI Rewrite 菜单项** → Phase 3 实现功能
- **多页分页逻辑** → Phase 5 PDF 导出时完善
- **深色/浅色主题切换** → Phase 7

## 风险评估

| 风险 | 概率 | 影响 | 缓解 |
|------|------|------|------|
| Monaco 冷启动慢 | 中 | 高 | 使用 @monaco-editor/loader + npm 离线包 |
| 大文档性能 | 低 | 中 | Monaco 内置虚拟化，Phase 7 进一步优化 |
| 渲染不一致 | 低 | 高 | 统一 markdown-it 实例 |

## 下一步

```bash
# 执行 Wave 1
/gsd-execute-phase 02 --wave 1

# 或执行完整 Phase 2
/gsd-execute-phase 02
```

---

*Generated: 2026-04-05*
*Plans: 4 plans in 4 waves*
*Status: Ready for execution (blockers fixed)*
