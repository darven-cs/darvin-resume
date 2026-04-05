# Phase 2: 核心编辑器 - Research

**Gathered:** 2026-04-05
**Status:** Ready for planning

## 技术方案研究

### 1. Monaco Editor 集成方案

#### 方案对比

| 方案 | 优点 | 缺点 | 推荐 |
|------|------|------|------|
| @monaco-editor/loader + npm | 离线可用、版本可控、打包优化 | 需要额外 webpack/vite 配置 | **推荐** |
| CDN (jsdelivr/unpkg) | 配置简单 | 无法离线、启动慢、违反隐私要求 | 不采用 |
| monaco-editor-core 直接引入 | 轻量 | 需要手动处理 workers | 备选 |

#### 推荐方案：@monaco-editor/loader + Vite

```bash
npm install @monaco-editor/loader monaco-editor
```

**Vite 配置 (vite.config.ts):**
```typescript
import monacoEditor from 'vite-plugin-monaco-editor';

export default defineConfig({
  plugins: [
    monacoEditor({
      customWorkers: [
        { label: 'editorWorkerService', entry: 'monaco-editor/esm/vs/editor/editor.worker' },
        { label: 'json', entry: 'monaco-editor/esm/vs/language/json/json.worker' },
      ],
    }),
  ],
});
```

**Vue 组件封装:**
```typescript
// src/components/MonacoEditor.vue
import loader from '@monaco-editor/loader';
import { ref, onMounted, onUnmounted, watch } from 'vue';

const editor = ref<any>(null);
let monacoInstance: any = null;

onMounted(async () => {
  monacoInstance = await loader.init();
  // 配置 markdown 语法高亮
  monacoInstance.languages.setMonarchTokensProvider('markdown', markdownTokens);
  editor.value = monacoInstance.editor.create(containerRef.value, {
    value: props.modelValue,
    language: 'markdown',
    theme: 'vs-dark',
    fontSize: 14,
    lineHeight: 1.6,
    wordWrap: 'on',
    minimap: { enabled: false },
    scrollBeyondLastLine: false,
    automaticLayout: true,
  });
});
```

#### Monaco Editor 版本选择
- **目标版本：** 0.45.0+ (与 Vue 3 + Vite 兼容稳定)
- **关键插件：** markdown-it (语法高亮)、custom theme

---

### 2. 双栏布局与拖拽调整

#### Vue 3 响应式布局实现

**Split Pane 组件设计:**
```vue
<!-- src/components/SplitPane.vue -->
<template>
  <div class="split-pane" :style="containerStyle">
    <div class="pane left-pane" :style="leftPaneStyle">
      <slot name="left" />
    </div>
    <div
      class="divider"
      @mousedown="startDrag"
      :class="{ dragging: isDragging }"
    >
      <div class="divider-handle" />
    </div>
    <div class="pane right-pane" :style="rightPaneStyle">
      <slot name="right" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue';

const props = defineProps<{
  defaultRatio?: number; // 默认 50
  minWidth?: number; // 每栏最小宽度 300px
}>();

const splitRatio = ref(props.defaultRatio ?? 50);
const isDragging = ref(false);

const containerStyle = computed(() => ({
  display: 'flex',
  height: '100%',
  width: '100%',
}));

const leftPaneStyle = computed(() => ({
  width: `${splitRatio.value}%`,
  minWidth: `${props.minWidth}px`,
}));

const rightPaneStyle = computed(() => ({
  width: `${100 - splitRatio.value}%`,
  minWidth: `${props.minWidth}px`,
}));

function startDrag(e: MouseEvent) {
  isDragging.value = true;
  const startX = e.clientX;
  const startRatio = splitRatio.value;

  const onMove = (e: MouseEvent) => {
    const container = document.querySelector('.split-pane') as HTMLElement;
    if (!container) return;
    const containerWidth = container.offsetWidth;
    const delta = e.clientX - startX;
    const deltaRatio = (delta / containerWidth) * 100;
    const newRatio = Math.max(20, Math.min(80, startRatio + deltaRatio));
    splitRatio.value = newRatio;
  };

  const onUp = () => {
    isDragging.value = false;
    document.removeEventListener('mousemove', onMove);
    document.removeEventListener('mouseup', onUp);
  };

  document.addEventListener('mousemove', onMove);
  document.addEventListener('mouseup', onUp);
}
</script>
```

#### 响应式断点处理

```vue
<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';

const windowWidth = ref(window.innerWidth);
const isSinglePane = computed(() => windowWidth.value < 1200);
const activeView = ref<'editor' | 'preview'>('editor');

function handleResize() {
  windowWidth.value = window.innerWidth;
  if (windowWidth.value >= 1200) {
    // 切回双栏模式
  }
}

onMounted(() => window.addEventListener('resize', handleResize));
onUnmounted(() => window.removeEventListener('resize', handleResize));
</script>
```

---

### 3. 实时预览同步

#### Markdown-it 性能优化

```typescript
// src/utils/markdown.ts (已初始化，参考 Phase 1)
import MarkdownIt from 'markdown-it';

// 配置高性能渲染
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: false,
  highlight: (str: string, lang: string) => {
    // 简单的高亮，无重型依赖
    return `<pre><code>${md.utils.escapeHtml(str)}</code></pre>`;
  },
});

// 防抖函数
function debounce<T extends (...args: any[]) => void>(
  fn: T,
  delay: number
): T {
  let timer: ReturnType<typeof setTimeout> | null = null;
  return ((...args: any[]) => {
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  }) as T;
}

// 渲染函数
export function renderMarkdown(content: string): string {
  return md.render(content);
}

// 防抖版本（150ms，满足 <200ms 要求）
export const debouncedRender = debounce(renderMarkdown, 150);
```

#### 编辑器与预览滚动同步

```typescript
// 编辑区滚动 → 预览区同步
editor.onDidScrollChange((e: IContentSizeChangedEvent) => {
  const scrollTop = e.scrollTop;
  const scrollHeight = e.scrollHeight;
  const clientHeight = e.height;
  const ratio = scrollTop / (scrollHeight - clientHeight);

  const preview = document.querySelector('.preview-content');
  if (preview) {
    preview.scrollTop = ratio * (preview.scrollHeight - preview.clientHeight);
  }
});
```

---

### 4. 行级交互

#### Monaco Editor Decorations API

```typescript
// 行级装饰：折叠图标、拖拽手柄
const decorations: string[] = [];
const newDecorations = editor.deltaDecorations(decorations, [
  {
    range: new monaco.Range(lineNumber, 1, lineNumber, 1),
    options: {
      isWholeLine: true,
      glyphMarginClassName: 'line-gutter-icon',
      glyphMarginHoverMessage: { value: 'Click for options' },
      beforeContentClassName: 'drag-handle',
    },
  },
]);
decorations.length = 0;
decorations.push(...newDecorations);
```

#### 折叠/展开图标实现

```typescript
// 监听编辑器内容变化，更新装饰
editor.onDidChangeModelContent(() => {
  updateLineDecorations();
});

function updateLineDecorations() {
  const model = editor.getModel();
  if (!model) return;

  const lines = model.getLineCount();
  const decorations: any[] = [];

  for (let i = 1; i <= lines; i++) {
    const lineContent = model.getLineContent(i);
    // 检测可折叠的块（标题、列表等）
    if (isFoldable(lineContent)) {
      decorations.push({
        range: new monaco.Range(i, 1, i, 1),
        options: {
          isWholeLine: true,
          glyphMarginClassName: 'fold-gutter',
          glyphMarginHoverMessage: { value: 'Click to fold' },
        },
      });
    }
  }

  editor.deltaDecorations(currentDecorations, decorations);
  currentDecorations = decorations;
}

function isFoldable(line: string): boolean {
  const trimmed = line.trim();
  return /^#{1,6}\s/.test(trimmed) || // 标题
         /^[-*+]\s/.test(trimmed) ||   // 无序列表
         /^\d+\.\s/.test(trimmed);     // 有序列表
}
```

#### 快捷菜单实现

```typescript
// 点击行首图标 → 显示上下文菜单
editor.onMouseDown((e: any) => {
  if (e.target.type === monaco.MouseTargetType.GUTTER_GLYPH_MARGIN) {
    const lineNumber = e.target.position?.lineNumber;
    showContextMenu(lineNumber);
  }
});

function showContextMenu(lineNumber: number) {
  // 使用 Vue Portal 渲染菜单到 body
  contextMenuVisible.value = true;
  contextMenuPosition.value = { lineNumber };

  // 菜单选项
  menuOptions.value = [
    { label: '上移', action: () => moveLine(lineNumber, -1) },
    { label: '下移', action: () => moveLine(lineNumber, 1) },
    { label: 'AI 重写', action: () => aiRewrite(lineNumber), disabled: true }, // Phase 3
    { label: '删除', action: () => deleteLine(lineNumber) },
  ];
}
```

#### 拖拽排序（HTML5 Drag API）

```typescript
// 拖拽开始
function onDragStart(e: DragEvent, lineNumber: number) {
  e.dataTransfer?.setData('text/plain', String(lineNumber));
  e.dataTransfer!.effectAllowed = 'move';
}

// 拖拽结束
function onDrop(e: DragEvent, targetLine: number) {
  const sourceLine = parseInt(e.dataTransfer?.getData('text/plain') || '0');
  if (sourceLine && sourceLine !== targetLine) {
    moveLineRange(sourceLine, targetLine);
  }
}
```

---

### 5. A4 分页线

#### CSS 实现方案

```css
/* A4 尺寸计算（96dpi） */
/* A4: 210mm × 297mm */
/* 96dpi 下: 210mm × 96/25.4 ≈ 794px, 297mm × 96/25.4 ≈ 1123px */

.preview-container {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: auto;
  background: #f5f5f5;
}

.a4-page {
  width: 210mm;
  min-height: 297mm;
  margin: 20px auto;
  padding: 20mm;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
}

/* A4 页面边界线 */
.a4-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border: 1px solid rgba(0, 120, 212, 0.3);
  pointer-events: none;
}

/* 分页符线 */
.page-break {
  position: relative;
  page-break-after: always;
  border-bottom: 1px dashed rgba(0, 120, 212, 0.4);
}

.page-break::after {
  content: attr(data-page);
  position: absolute;
  right: -40px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 10px;
  color: rgba(0, 120, 212, 0.5);
}
```

#### JavaScript 分页计算

```typescript
function calculatePageBreaks(content: string, pageHeight: number = 1123): number[] {
  // 估算每页高度（考虑字体大小、行高、边距）
  // 返回需要分页的行索引数组
  const lines = content.split('\n');
  const pageBreaks: number[] = [];
  let currentHeight = 0;
  const lineHeight = 24; // 约 14px 字体 + 1.6 行高

  lines.forEach((line, index) => {
    const lineHeight = estimateLineHeight(line);
    if (currentHeight + lineHeight > pageHeight) {
      pageBreaks.push(index);
      currentHeight = lineHeight;
    } else {
      currentHeight += lineHeight;
    }
  });

  return pageBreaks;
}
```

---

### 6. 依赖清单

```json
{
  "dependencies": {
    "@monaco-editor/loader": "^1.4.0",
    "monaco-editor": "^0.45.0",
    "markdown-it": "^14.0.0"
  },
  "devDependencies": {
    "vite-plugin-monaco-editor": "^1.1.0"
  }
}
```

---

### 7. 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| Monaco Editor 加载慢 | 中 | 高 | 使用 @monaco-editor/loader + 缓存 |
| 大文档性能问题 | 低 | 中 | 虚拟滚动、懒加载 |
| 渲染不一致 | 低 | 高 | 统一 markdown-it 实例 |
| Vite 配置兼容 | 低 | 低 | 使用标准配置 |

---

## 参考文档

- Monaco Editor API: https://microsoft.github.io/monaco-editor/api/
- @monaco-editor/loader: https://github.com/suren-atoyan/monaco-react
- Markdown-it: https://github.com/markdown-it/markdown-it
- Vue 3 Composition API: https://vuejs.org/guide/extras/composition-api-faq.html
