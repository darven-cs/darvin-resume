<template>
  <div class="monaco-editor-wrapper">
    <div ref="containerRef" class="monaco-editor-container" />
    <!-- 上下文菜单 per EDIT-10 -->
    <ContextMenu
      :items="contextMenuItems"
      :x="contextMenuPosition.x"
      :y="contextMenuPosition.y"
      :visible="contextMenuVisible"
      @close="contextMenuVisible = false"
    />
    <!-- AI 浮动工具栏 per EDIT-07 -->
    <AIFloatingToolbar
      :visible="toolbarVisible"
      :position="toolbarPosition"
      :selected-text="selectedText"
      :loading="isLoading"
      :current-operation="currentOperation"
      @operate="handleAIOperation"
    />
    <!-- AI Diff 对比视图 per EDIT-11 -->
    <AIDiffView
      :visible="aiDiff.visible"
      :original-text="aiDiff.originalText"
      :modified-text="aiDiff.modifiedText"
      :operation-type="aiDiff.operationType"
      :streaming="aiDiff.streaming"
      :position="toolbarPosition"
      @accept="handleDiffAccept"
      @reject="handleDiffReject"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as monacoEditor from 'monaco-editor'
import ContextMenu from './ContextMenu.vue'
import AIFloatingToolbar from './AIFloatingToolbar.vue'
import AIDiffView from './AIDiffView.vue'
import type { MenuItem } from './ContextMenu.vue'
import { useAISelection } from '../composables/useAISelection'
import type { AIOperationType } from '../types/ai'

// Monaco Editor 实例类型
let editor: any = null
let monaco: any = null

const containerRef = ref<HTMLElement | null>(null)

// Emits 定义
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
  (e: 'load', editor: any, monaco: any): void
}>()

// ============================================================
// 行级交互功能 per EDIT-08, EDIT-09, EDIT-10
// ============================================================

// 折叠状态存储
const foldedRanges = new Map<number, { endLine: number; originalContent: string }>()
let currentDecorations: string[] = []

// 上下文菜单状态 per EDIT-10
const contextMenuItems = ref<MenuItem[]>([])
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })

// 拖拽排序状态 per EDIT-09
let dragSourceLine: number | null = null

// ============================================================
// AI 选区工具栏状态 per EDIT-07
// ============================================================

// Props: jobTarget 用于 AI 上下文
interface Props {
  modelValue?: string
  language?: string
  readOnly?: boolean
  options?: Record<string, any>
  jobTarget?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  language: 'markdown',
  readOnly: false,
  options: () => ({}),
  jobTarget: '',
})

// AI 选区 composable
const {
  hasSelection,
  selectedText,
  selectionRange,
  toolbarVisible,
  toolbarPosition,
  isLoading,
  currentOperation,
  updateSelection,
  performAIOperation,
} = useAISelection(() => editor, props.jobTarget)

// ============================================================
// AI Diff 对比视图状态 per EDIT-11
// ============================================================

interface AIDiffState {
  visible: boolean
  originalText: string
  modifiedText: string
  operationType: AIOperationType
  range: { startLineNumber: number; startColumn: number; endLineNumber: number; endColumn: number } | null
  streaming: boolean
}

const aiDiff = reactive<AIDiffState>({
  visible: false,
  originalText: '',
  modifiedText: '',
  operationType: 'polish',
  range: null,
  streaming: false,
})

// ============================================================
// 行首图标装饰 (Gutter Icons) per EDIT-08
// ============================================================

function isFoldable(line: string): boolean {
  // 标题行
  if (/^#{1,6}\s/.test(line)) return true
  // 无序列表
  if (/^[-*+]\s/.test(line)) return true
  // 有序列表
  if (/^\d+\.\s/.test(line)) return true
  return false
}

function getIndentLevel(line: string): number {
  return line.search(/\S/)
}

function updateLineDecorations() {
  const model = editor.getModel()
  if (!model || !monaco) return

  const decorations: any[] = []

  for (let i = 1; i <= model.getLineCount(); i++) {
    const lineContent = model.getLineContent(i)
    const trimmed = lineContent.trim()

    // 检测可折叠的块 per EDIT-08
    if (isFoldable(trimmed)) {
      const isFolded = foldedRanges.has(i)
      decorations.push({
        range: new monaco.Range(i, 1, i, 1),
        options: {
          isWholeLine: true,
          // 行首图标 per EDIT-08
          glyphMarginClassName: isFolded ? 'fold-gutter fold-gutter-collapsed' : 'fold-gutter fold-gutter-expanded',
          glyphMarginHoverMessage: {
            value: isFolded ? '**展开内容**\n点击展开 | 右键菜单' : '**折叠内容**\n点击折叠 | 右键菜单'
          },
          // 折叠范围标记
          linesDecorationsClassName: 'fold-decoration',
        },
      })
    }
  }

  currentDecorations = editor.deltaDecorations(currentDecorations, decorations)
}

// ============================================================
// 折叠/展开功能 per EDIT-08
// ============================================================

function foldBlock(startLine: number) {
  const model = editor.getModel()
  if (!model || !monaco) return

  // 找到折叠范围（下一个同级标题或同级列表为止）
  const startContent = model.getLineContent(startLine)
  const startIndent = getIndentLevel(startContent)
  let endLine = startLine + 1

  while (endLine <= model.getLineCount()) {
    const lineContent = model.getLineContent(endLine)
    const trimmed = lineContent.trim()
    if (!trimmed) { endLine++; continue }

    const indent = getIndentLevel(lineContent)
    const isHeader = /^#{1,6}\s/.test(trimmed)
    const isList = /^[-*+]\s/.test(trimmed) || /^\d+\.\s/.test(trimmed)

    if (indent <= startIndent && (isHeader || isList)) break
    if (isHeader && /^#/.test(trimmed)) break
    endLine++
  }

  endLine--

  if (endLine > startLine) {
    const originalContent = model.getLinesContent().slice(startLine - 1, endLine).join('\n')
    foldedRanges.set(startLine, { endLine, originalContent })

    // 替换为折叠标记行
    editor.executeEdits('fold', [{
      range: new monaco.Range(startLine, 1, endLine, model.getLineMaxColumn(endLine)),
      text: `${startContent} <... ${endLine - startLine} lines hidden ...>`
    }])

    updateLineDecorations()
  }
}

function unfoldBlock(startLine: number) {
  const range = foldedRanges.get(startLine)
  if (!range || !monaco) return

  const model = editor.getModel()
  if (!model) return

  // 恢复完整内容
  editor.executeEdits('unfold', [{
    range: new monaco.Range(startLine, 1, startLine, model.getLineMaxColumn(startLine)),
    text: range.originalContent
  }])

  foldedRanges.delete(startLine)
  updateLineDecorations()
}

// ============================================================
// 上移/下移/删除功能 per EDIT-10
// ============================================================

function moveLine(lineNumber: number, direction: -1 | 1) {
  const model = editor.getModel()
  if (!model || !monaco) return

  const targetLine = lineNumber + direction
  if (targetLine < 1 || targetLine > model.getLineCount()) return

  const currentLine = model.getLineContent(lineNumber)
  const targetLineContent = model.getLineContent(targetLine)

  editor.executeEdits('moveLine', [
    {
      range: new monaco.Range(lineNumber, 1, lineNumber, model.getLineMaxColumn(lineNumber)),
      text: targetLineContent,
    },
    {
      range: new monaco.Range(targetLine, 1, targetLine, model.getLineMaxColumn(targetLine)),
      text: currentLine,
    },
  ])

  // 移动后光标位置跟随
  editor.setPosition({ lineNumber: targetLine, column: 1 })
}

function deleteLine(lineNumber: number) {
  const model = editor.getModel()
  if (!model || !monaco) return

  editor.executeEdits('deleteLine', [{
    range: new monaco.Range(lineNumber, 1, lineNumber, model.getLineMaxColumn(lineNumber)),
    text: '',
  }])
}

// ============================================================
// 拖拽排序 per EDIT-09
// ============================================================

function moveLineRange(sourceLine: number, targetLine: number) {
  const model = editor.getModel()
  if (!model || !monaco) return

  if (sourceLine === targetLine) return

  const lines: string[] = []
  let startLine = sourceLine
  let endLine = sourceLine

  // 找到要移动的连续块
  const sourceIndent = getIndentLevel(model.getLineContent(sourceLine))

  // 向上扩展
  while (startLine > 1) {
    const prevContent = model.getLineContent(startLine - 1)
    const prevIndent = getIndentLevel(prevContent)
    if (prevIndent < sourceIndent && /^[-*+]/.test(prevContent.trim())) break
    if (/^#{1,6}\s/.test(prevContent.trim()) && prevIndent <= sourceIndent) break
    startLine--
  }

  // 向下扩展
  while (endLine < model.getLineCount()) {
    const nextContent = model.getLineContent(endLine + 1)
    const nextIndent = getIndentLevel(nextContent)
    if (nextIndent < sourceIndent && /^[-*+]/.test(nextContent.trim())) break
    if (/^#{1,6}\s/.test(nextContent.trim()) && nextIndent <= sourceIndent) break
    endLine++
  }

  const blockLines = model.getLinesContent().slice(startLine - 1, endLine).join('\n')

  editor.executeEdits('moveBlock', [
    {
      range: new monaco.Range(startLine, 1, endLine, model.getLineMaxColumn(endLine)),
      text: '',
    },
  ])

  // 重新计算插入位置（可能在删除后发生变化）
  const insertLine = targetLine > endLine ? targetLine - (endLine - startLine) : targetLine

  editor.executeEdits('insertBlock', [{
    range: new monaco.Range(insertLine, 1, insertLine, 1),
    text: blockLines + '\n',
  }])
}

// ============================================================
// 编辑器初始化
// ============================================================

function initEditor() {
  if (!containerRef.value) return

  monaco = monacoEditor
  editor = monaco.editor.create(containerRef.value, {
    value: props.modelValue,
    language: props.language,
    theme: 'vs',
    readOnly: props.readOnly,
    fontSize: 14,
    lineHeight: 1.6,
    fontFamily: "'JetBrains Mono', 'Fira Code', Consolas, 'Courier New', monospace",
    fontLigatures: true,
    glyphMargin: true,
    minimap: { enabled: false },
    automaticLayout: true,
    wordWrap: 'on',
    lineNumbers: 'on',
    renderLineHighlight: 'line',
    scrollBeyondLastLine: false,
    padding: { top: 16, bottom: 16 },
    smoothScrolling: true,
    cursorBlinking: 'smooth',
    cursorSmoothCaretAnimation: 'on',
    multiCursorModifier: 'ctrlCmd',
    undoPerCell: false,
    // 快捷键配置 - VS Code 兼容
    keyboard: {
      // 使用默认的 VS Code 键位映射
    },
    // Find/Replace 配置
    find: {
      addExtraSpaceOnTop: false,
      autoFindInSelection: 'never',
      seedSearchStringFromSelection: 'always',
    },
    // 快捷查找
    quickSuggestions: {
      other: true,
      comments: false,
      strings: false,
    },
    suggestOnTriggerCharacters: true,
    acceptSuggestionOnEnter: 'on',
    tabCompletion: 'on',
    formatOnPaste: false,
    formatOnType: false,
    // Markdown 特定配置
    ...props.options,
  })

  // 监听内容变化
  editor.onDidChangeModelContent(() => {
    const value = editor.getValue()
    emit('update:modelValue', value)
    emit('change', value)
    // 内容变化时更新装饰
    updateLineDecorations()
  })

  // ============================================================
  // 选区变更监听 - AI 工具栏 per EDIT-07
  // ============================================================
  editor.onDidChangeCursorSelection((e: any) => {
    updateSelection()
  })

  // ============================================================
  // 鼠标点击事件（折叠/展开 + 快捷菜单）per EDIT-08, EDIT-10
  // ============================================================
  editor.onMouseDown((e: any) => {
    const target = e.target

    // 点击行首图标区域
    if (target.type === monaco.editor.MouseTargetType.GUTTER_GLYPH_MARGIN) {
      const lineNumber = target.position?.lineNumber
      if (!lineNumber) return

      const lineContent = editor.getModel()?.getLineContent(lineNumber) || ''
      const trimmed = lineContent.trim()

      if (!isFoldable(trimmed)) return

      // 左键点击：折叠/展开
      if (e.event?.leftButton) {
        if (foldedRanges.has(lineNumber)) {
          // 展开
          unfoldBlock(lineNumber)
        } else {
          // 折叠
          foldBlock(lineNumber)
        }
      }
      return
    }

    // 右键点击：显示快捷菜单 per EDIT-10
    if (e.event?.rightButton) {
      const position = editor.getPosition()
      if (!position) return

      const coords = editor.getScrolledVisiblePosition(position)
      if (!coords) return

      const editorDom = editor.getDomNode()
      if (!editorDom) return

      const rect = editorDom.getBoundingClientRect()

      contextMenuItems.value = [
        {
          id: 'move-up',
          label: '上移',
          icon: '↑',
          shortcut: 'Alt+↑',
          action: () => moveLine(position.lineNumber, -1)
        },
        {
          id: 'move-down',
          label: '下移',
          icon: '↓',
          shortcut: 'Alt+↓',
          action: () => moveLine(position.lineNumber, 1)
        },
        {
          id: 'ai-rewrite',
          label: 'AI 重写',
          icon: '✦',
          shortcut: '',
          disabled: true,  // Phase 3 占位
          action: () => {} // Placeholder
        },
        {
          id: 'delete',
          label: '删除',
          icon: '✕',
          shortcut: 'Ctrl+Shift+Delete',
          danger: true,
          action: () => deleteLine(position.lineNumber)
        },
      ]

      contextMenuPosition.value = {
        x: rect.left + coords.left + 40,
        y: rect.top + coords.top + 20,
      }
      contextMenuVisible.value = true
    }
  })

  // ============================================================
  // 拖拽排序 per EDIT-09
  // ============================================================
  const editorDom = editor.getDomNode()

  editorDom?.addEventListener('dragstart', (e: DragEvent) => {
    const position = editor.getPosition()
    if (position) {
      dragSourceLine = position.lineNumber
      if (e.dataTransfer) {
        e.dataTransfer.effectAllowed = 'move'
        e.dataTransfer.setData('text/plain', String(dragSourceLine))
      }
    }
  })

  editorDom?.addEventListener('dragover', (e: DragEvent) => {
    if (dragSourceLine === null) return
    e.preventDefault()
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'move'
    }
  })

  editorDom?.addEventListener('drop', (e: DragEvent) => {
    if (dragSourceLine === null) return
    e.preventDefault()

    const position = editor.getPosition()
    if (!position) {
      dragSourceLine = null
      return
    }

    const targetLine = position.lineNumber
    if (targetLine !== dragSourceLine && targetLine > 0) {
      moveLineRange(dragSourceLine, targetLine)
    }

    dragSourceLine = null
  })

  editorDom?.addEventListener('dragend', () => {
    dragSourceLine = null
  })

  // 初始更新行首装饰
  updateLineDecorations()

  emit('load', editor, monaco)
}

// 暴露方法给父组件
defineExpose({
  getEditor: () => editor,
  getMonaco: () => monaco,
  focus: () => editor?.focus(),
  layout: () => editor?.layout(),
  /**
   * Gets the currently selected text in the editor.
   */
  getSelection: () => editor?.getModel()?.getValueInRange(editor.getSelection()!),
  /**
   * Inserts text at the current cursor position.
   */
  insertAtCursor: (text: string) => {
    if (!editor) return
    const position = editor.getPosition()
    if (!position) return
    editor.executeEdits('insert', [{
      range: new monaco.Range(position.lineNumber, position.column, position.lineNumber, position.column),
      text: text,
    }])
    editor.focus()
  },
})

// ============================================================
// AI 操作处理器 per EDIT-07 + Diff 对比 per EDIT-11
// ============================================================

async function handleAIOperation(operation: string) {
  if (!editor) return

  // 【关键修复】用 Monaco API 直接获取选区，绕过 mousedown 清空 reactive 状态的问题
  // 事件顺序：mousedown → Monaco onDidChangeCursorSelection → selectionRange.value=null → click → handleAIOperation
  // 所以这里必须从 Monaco API 取原始选区，而不是依赖响应式状态
  const rawSelection = editor.getSelection()
  if (!rawSelection) return
  const model = editor.getModel()
  if (!model) return
  const rawText = model.getValueInRange(rawSelection)
  if (!rawText || !rawText.trim()) return

  // 闭包捕获捕获的原始选区，后续不再依赖响应式状态
  const originalText = rawText
  const savedRange = {
    startLineNumber: rawSelection.startLineNumber,
    startColumn: rawSelection.startColumn,
    endLineNumber: rawSelection.endLineNumber,
    endColumn: rawSelection.endColumn,
  }

  // 显示 diff 视图，进入 streaming 状态
  aiDiff.visible = true
  aiDiff.originalText = originalText
  aiDiff.modifiedText = ''
  aiDiff.operationType = operation as AIOperationType
  aiDiff.range = savedRange
  aiDiff.streaming = true

  // 隐藏工具栏（diff 视图替代显示）
  toolbarVisible.value = false

  // 执行 AI 操作（将原始选区文本作为参数传入，不再依赖响应式状态）
  const result = await performAIOperation(operation as AIOperationType, originalText)

  // 更新 diff 内容
  if (result) {
    aiDiff.modifiedText = result
  }
  aiDiff.streaming = false

  // 如果操作失败（空结果），关闭 diff 视图
  if (!result) {
    aiDiff.visible = false
  }
}

/**
 * 接受 AI 修改：使用 executeEdits 替换选区内容，合并到单个 undo 单元
 */
function handleDiffAccept() {
  if (!editor || !monaco || !aiDiff.range || !aiDiff.modifiedText) return

  const range = new monaco.Range(
    aiDiff.range.startLineNumber,
    aiDiff.range.startColumn,
    aiDiff.range.endLineNumber,
    aiDiff.range.endColumn,
  )

  // 使用 executeEdits 替换，支持 Ctrl+Z 一次性撤销
  editor.executeEdits('ai-accept', [{
    range,
    text: aiDiff.modifiedText,
    forceMoveMarkers: true,
  }])

  // 关闭 diff 视图
  closeDiffView()
}

/**
 * 拒绝 AI 修改：关闭 diff 视图，原文保持不变
 */
function handleDiffReject() {
  closeDiffView()
}

/**
 * 关闭 diff 视图并重置状态
 */
function closeDiffView() {
  aiDiff.visible = false
  aiDiff.originalText = ''
  aiDiff.modifiedText = ''
  aiDiff.range = null
  aiDiff.streaming = false
}

// 监听外部 modelValue 变化，同步到编辑器
watch(
  () => props.modelValue,
  (newValue) => {
    if (editor && newValue !== editor.getValue()) {
      const position = editor.getPosition()
      editor.setValue(newValue)
      if (position) {
        editor.setPosition(position)
      }
    }
  },
)

// 监听 language 变化
watch(
  () => props.language,
  (newLang) => {
    if (editor) {
      const model = editor.getModel()
      if (model) {
        monaco?.editor.setModelLanguage(model, newLang)
      }
    }
  },
)

// 监听 readOnly 变化
watch(
  () => props.readOnly,
  (readOnly) => {
    editor?.updateOptions({ readOnly })
  },
)

onMounted(() => {
  nextTick(() => {
    initEditor()
  })
})

onUnmounted(() => {
  editor?.dispose()
  editor = null
  monaco = null
  foldedRanges.clear()
  currentDecorations = []
})
</script>

<style scoped>
.monaco-editor-wrapper {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
}

.monaco-editor-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>

<style>
/* 行首折叠图标样式 (非 scoped，作用于 Monaco 编辑器内) per EDIT-08 */
.fold-gutter {
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  font-size: 10px;
  user-select: none;
}

.fold-gutter-expanded::before {
  content: '▼';
  color: #6e7681;
  font-size: 8px;
}

.fold-gutter-collapsed::before {
  content: '▶';
  color: #6e7681;
  font-size: 8px;
}

.fold-gutter:hover::before {
  color: #c6c6c6;
}

.fold-decoration {
  border-left: 2px solid #3b82f680;
  margin-left: 2px;
}
</style>
