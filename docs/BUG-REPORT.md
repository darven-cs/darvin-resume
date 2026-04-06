# Darvin-Resume Bug 分析报告

> 生成时间：2026-04-06

---

## Bug 1: PDF导出 — 系统打印空白 + 高级导出只一页

### 1-A: 系统打印（window.print()）空白

**文件**: `frontend/src/components/ExportDialog.vue`

**根因**: `@media print` 规则写在 `<style scoped>` 块里，被 Vue scoped CSS 处理机制破坏。

```css
/* ExportDialog.vue:449-767 —— <style scoped> 块 */
@media print {
  body > *:not(.print-container) { display: none !important; }  /* 不生效 */
  .a4-page { padding: 20mm !important; }
}
```

Vue scoped CSS 通过给选择器附加 `[data-v-xxxx]` 属性实现隔离，`@media print` 块内的选择器也会被附加。结果 `body > *:not(.print-container)` 变成了 `body > *:not(.print-container)[data-v-xxxx]`，而 `<body>` 本身和它的直接子元素都没有这个 data 属性，所以规则永远匹配不到任何元素。

**修复方案**: 将 `@media print` 块移到 `<style>` (非scoped)，或在 `editor.css` 中单独写一个全局 `@media print` 规则。

---

### 1-B: 高级导出（Chromedp）只导出一页

**文件**:
- `frontend/src/components/A4Page.vue:80-86`
- `internal/export/chromedp.go:100-120`

**根因**: 多页分页逻辑是 TODO，尚未实现。

```typescript
// A4Page.vue:80-86
const pages = computed(() => {
  if (!props.content) return []
  const html = renderMarkdown(props.content)
  // Phase 5: 多页分页逻辑将在 PDF 导出时完善（实际页面换行计算）
  return [html]  // ← 永远只返回1页！
})
```

所有简历内容被渲染成一个 HTML 字符串，放进一个 `.a4-page` div。Chromedp 的 `PrintToPDF` 只生成一个 PDF 页面，内容超出一页的部分被截断。

**修复方案**: 实现 A4 高度感知的分页逻辑，将 HTML 分割成多个 `.a4-page` div，每个加 `page-break-after: always` CSS。

---

## Bug 2: AI润色/缩写/重写/翻译 — 点击后无反应

**文件**: `frontend/src/components/MonacoEditor.vue:567-598`

**根因**: 点击浮动工具栏按钮时，Monaco编辑器失去焦点导致选区被清空，`handleAIOperation` 提前返回。

```typescript
// MonacoEditor.vue:567-568
async function handleAIOperation(operation: string) {
  if (!editor || !selectionRange.value) return  // ← 点击工具栏后 selectionRange 已是 null！
}
```

浮动工具栏使用 `<Teleport to="body">` 挂载在 DOM 树外部。当用户鼠标点击工具栏按钮时：

1. 点击事件 → Monaco 编辑器 `mousedown` 事件触发
2. 编辑器失去焦点，`onDidChangeCursorSelection` 回调执行
3. `updateSelection()` 检测无选中区 → `selectionRange.value = null`
4. 用户松手 → 工具栏 `@click` 触发 `handleAIOperation`
5. `selectionRange.value === null` → 函数直接 return

**修复方案**: 在 `handleAIOperation` 开头，将 `selectedText` 和 `selectionRange` 的值复制到局部变量（闭包捕获），而不是依赖可能被清空的响应式状态。

**需求补充**: AI聊天侧边栏（`AIChatSidebar.vue`）的引用功能是空实现占位符，需要实现选中简历内容并引用的完整流程。

---

## Bug 3: 样式调整 — 字号/行高/边距无效

**文件**:
- `frontend/src/components/A4Page.vue:123`
- `frontend/src/components/StyleEditor.vue:342-350`

**根因**: 双重根因。

**根因1**: `A4Page.vue` 的 scoped CSS 硬编码了样式值，覆盖了 CSS 变量的全局规则。

```css
/* A4Page.vue:123-132 — <style scoped> */
.a4-page {
  padding: 20mm;  /* ← 硬编码！CSS 变量被覆盖 */
}
```

CSS 特异性比较：`.a4-page[data-v-xxxxx]` (scoped) > `.a4-page` (全局)，scoped 规则永远优先。

**根因2**: `StyleEditor.vue` 的 `saveToBackend()` 函数有 bug，没有保存 CSS 数据到后端。

```typescript
// StyleEditor.vue:342-350
function saveToBackend() {
  if (saveTimerHandle) clearTimeout(saveTimerHandle)
  saveTimer = Date.now()
  saveTimerHandle = setTimeout(() => {
    emit('update:modelValue', props.modelValue)  // ← 只 emit 面板可见性，没保存 CSS！
  }, 500)
}
```

它没有调用 `buildFullCSS()` 生成 CSS 字符串，也没有调用 `UpdateResumeCustomCSS` 持久化到后端。刷新页面后样式丢失。

**修复方案**:
1. 将 A4Page.vue scoped CSS 中的硬编码值改为 CSS 变量
2. 修复 `saveToBackend()` 调用 `buildFullCSS()` 并通过 `UpdateResumeCustomCSS` 保存

---

## Bug 4: AI助手 — 没有流式输出

**文件**: `frontend/src/components/AIChatSidebar.vue:37-45, 180-219`

**根因**: `useAIStream` 在组件 setup 时创建，监听的是初始 operationId，每次发消息都生成新的 operationId，导致事件监听错位。

```typescript
// AIChatSidebar.vue:37-45 — setup 阶段执行，只执行一次
const operationId = ref(crypto.randomUUID())
const { content: streamedContent, isStreaming: streamActive, ... } = useAIStream(operationId.value)

// AIChatSidebar.vue:180-219 — 每次发消息
async function handleSend() {
  const opId = crypto.randomUUID()      // ← 新的 operationId
  operationId.value = opId
  reset()                               // ← 只清空状态，不重新注册事件监听！
  await sendChatMessage(opId, ...)     // ← Go 后端向 ai:stream:{opId} 发事件
  // ← 前端还在监听 ai:stream:{初始Id}，新事件全部丢失！
}
```

**修复方案**: 每次发送消息时，创建新的 `useAIStream` 实例，或者让 `useAIStream` 支持动态更新监听的事件名。

---

## Bug 5: 设置按钮 — 应在简历列表最下方

**文件**: `frontend/src/views/HomeView.vue:41-76`

**根因**: `HomeView.vue` 的简历列表区域没有"设置"入口，设置按钮目前仅存在于 `EditorView.vue` 工具栏。

**修复方案**: 在 `HomeView.vue` 的 `resume-grid-area` 最下方（`RecycleBinSection` 之前）添加设置入口按钮。

---

## 汇总表

| # | 问题 | 严重度 | 根因文件 | 根因 |
|---|------|--------|----------|------|
| 1A | 系统打印空白 | 高 | `ExportDialog.vue` | `@media print` 在 scoped 块内，选择器被 data 属性破坏 |
| 1B | 高级导出只一页 | 高 | `A4Page.vue:80` | 多页分页未实现，永远 `return [html]` 单页 |
| 2 | AI操作无反应 | 高 | `MonacoEditor.vue:568` | 点击工具栏清空选区后 `handleAIOperation` guard 触发 |
| 3 | 样式调整无效 | 中 | `A4Page.vue:123` + `StyleEditor.vue:342` | scoped CSS 硬编码覆盖变量 + `saveToBackend()` 不保存 |
| 4 | AI无流式输出 | 中 | `AIChatSidebar.vue:45` | `useAIStream` 用初始 operationId，后续消息新 Id 事件丢失 |
| 5 | 设置入口位置 | 低 | `HomeView.vue` | 简历列表页缺少设置入口 |
