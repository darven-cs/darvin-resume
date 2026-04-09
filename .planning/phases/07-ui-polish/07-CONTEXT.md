# Phase 7: 界面打磨与健壮性 — 上下文文档

**创建时间**: 2026-04-08
**Phase目标**: 应用界面专业美观，支持深色/浅色切换，窗口自适应，所有异常场景有兜底，交互流畅完整

## 需求映射

| 需求ID | 描述 | 当前状态 |
|--------|------|----------|
| UIUX-01 | 极简主义视觉风格，低饱和度专业色系 | ❌ 颜色硬编码，无设计令牌体系 |
| UIUX-02 | 深色/浅色双模式，系统主题自动切换 | ❌ 完全不存在 |
| UIUX-03 | 编辑区等宽字体，预览区无衬线字体 | ⚠️ 已有但硬编码在组件内 |
| UIUX-04 | 响应式适配，最小1200px | ⚠️ EditorView已有基础实现 |
| UIUX-05 | <1200px自动单栏，手动切换编辑/预览 | ⚠️ 已有isSinglePane计算逻辑 |
| UIUX-06 | AI快捷键(Ctrl+R/T/D) | ❌ 全局快捷键系统不存在 |
| UIUX-07 | 自定义快捷键，本地持久化 | ❌ 完全不存在 |
| UIUX-08 | 空状态设计，模板Demo预览 | ⚠️ HomeView有简单空状态 |
| UIUX-09 | 全异步操作加载提示 | ⚠️ SaveStatusIndicator仅覆盖保存 |
| UIUX-10 | 全操作即时反馈(toast/弹窗) | ❌ 仅AIErrorToast，无通用Toast |
| UIUX-11 | 全场景异常兜底 | ⚠️ AI错误完善，其他场景缺失 |

## 当前代码库状态

### 样式架构

**全局样式文件** (6个):
- `frontend/src/styles/editor.css` — 简历渲染CSS变量+A4纸张+Markdown元素+打印规则
- `frontend/src/styles/templates/template-{minimal,dual-col,academic,campus}.css` — 4套简历模板
- `frontend/src/style.css` — 极简全局基础样式

**UI样式现状**:
- 无CSS框架，全部手写原生CSS
- 颜色值大量硬编码（`#1e1e1e`, `#252526`, `#2d2d2d`, `#3c3c3c` 等在25个组件中重复）
- 简历CSS变量（`--resume-*`）仅控制简历内容渲染，不控制UI主题
- 无设计令牌(design tokens)体系

### 主题系统

**现状**: 完全不存在
- UI硬编码为暗色（`body { background-color: #1e1e1e }`）
- Monaco编辑器固定 `theme: 'vs'`（浅色）
- 无 `prefers-color-scheme` 媒体查询
- 后端无主题相关settings key或Bridge方法
- `main.go` 窗口背景色硬编码 `RGBA{R: 27, G: 38, B: 54, A: 1}`

### 响应式布局

**已实现**:
- `EditorView.vue`: `isSinglePane = computed(() => windowWidth < 1200)` — 双栏/单栏自动切换
- `HomeView.vue`: `grid-template-columns: repeat(auto-fill, minmax(280px, 1fr))` — 简历卡片自适应
- `TemplateSelector.vue`: `@media (max-width: 600px)` — 模板网格2列

**缺失**:
- Sidebar响应式（固定240px宽度，窄屏下可能挤压主内容）
- 没有手动切换编辑/预览视图的UI按钮（仅自动切换）
- 单栏模式下缺少Tab切换UI

### 快捷键

**现状**: 零散分布，无统一系统
- `Ctrl+S` 保存 — `EditorView.vue`
- Monaco编辑器内置VS Code键位
- 右键菜单硬编码快捷键标签（Alt+↑/↓, Ctrl+Shift+Delete）
- 无AI快捷键（Ctrl+R/T/D）
- 后端无快捷键存储
- 无自定义快捷键机制

### Toast/反馈系统

**已有**:
- `AIErrorToast.vue` — AI错误专用，3级严重度，支持重试
- `SaveStatusIndicator.vue` — 保存状态嵌入式指示器

**缺失**:
- 通用Toast通知组件（成功/警告/错误/信息）
- 导出成功/失败反馈（用console.log）
- 重命名/删除/复制操作反馈（用confirm()）
- 加载状态覆盖（多数异步操作无loading提示）

### 空状态

**已有**:
- HomeView: SVG图标 + "还没有简历" + 新建按钮
- HomeView: 搜索无结果提示
- A4Page: "开始编写Markdown"提示

**缺失**:
- 通用空状态组件（每个页面各自内联实现）
- 模板Demo预览卡片

### 后端接口

**需要新增**:
1. `ui.theme` settings key — `light`/`dark`/`system`
2. `ui.shortcuts` settings key — JSON序列化的快捷键映射
3. `GetTheme()` / `SetTheme()` Bridge方法
4. `GetShortcuts()` / `SetShortcuts()` Bridge方法

**已有的settings基础设施**: `internal/settings/settings.go` 提供通用key-value CRUD

### 关键约束

1. **简历渲染不可影响**: 主题切换仅影响UI外壳，简历预览始终保持白色A4纸张+用户选择模板样式
2. **Monaco主题同步**: 深色模式需同步切换Monaco编辑器主题(`vs` ↔ `vs-dark`)
3. **性能要求**: 主题切换即时生效，无闪烁；冷启动<2s不受影响
4. **CSS变量策略**: 新增`--ui-*`前缀的设计令牌，与现有`--resume-*`简历变量隔离
5. **Wails窗口背景**: 主题切换需同步更新Go端窗口背景色

## 依赖关系

- Phase 1-6 全部完成 ✅
- 无外部依赖阻塞
