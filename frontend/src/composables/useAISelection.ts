import { ref } from 'vue'
import { sendMessage } from '../services/ai'
import { useAIStream } from './useAIStream'
import type { AIOperationType } from '../types/ai'

/**
 * Operation configuration: labels and prompts.
 */
export const OPERATION_LABELS: Record<AIOperationType, string> = {
  polish: '润色',
  translate: '翻译',
  summarize: '缩写',
  rewrite: '重写',
}

type PromptBuilder = (text: string, job: string) => string

const OPERATION_PROMPTS: Record<AIOperationType, PromptBuilder> = {
  polish: (text, job) =>
    `请润色以下简历内容，使语言更专业、精炼${job ? `（目标岗位：${job}）` : ''}：\n${text}`,
  translate: (text, _job) =>
    `请将以下内容翻译成英文，保持专业语气：\n${text}`,
  summarize: (text, _job) =>
    `请将以下内容压缩到原长度的60%，保留关键信息：\n${text}`,
  rewrite: (text, job) =>
    `请重写以下简历内容，换一种表达方式${job ? `（目标岗位：${job}）` : ''}：\n${text}`,
}

interface SelectionRange {
  startLineNumber: number
  startColumn: number
  endLineNumber: number
  endColumn: number
}

/**
 * Composable for handling text selection in Monaco Editor and triggering AI operations.
 *
 * @param editorGetter - Getter function returning the Monaco editor instance
 * @param jobTarget - Current job target for AI context
 */
export function useAISelection(
  editorGetter: () => unknown,
  jobTarget: string = ''
) {
  // Selection state
  const hasSelection = ref(false)
  const selectedText = ref('')
  const selectionRange = ref<SelectionRange | null>(null)
  const toolbarPosition = ref({ top: 0, left: 0 })
  const toolbarVisible = ref(false)

  // AI operation state
  const isLoading = ref(false)
  const currentOperation = ref('')

  // Current operation ID for streaming
  let currentOperationId = ''

  /**
   * Updates selection state and toolbar visibility/position.
   */
  function updateSelection() {
    const editor = editorGetter() as {
      getSelection: () => SelectionRange | null
      getModel: () => {
        getValueInRange: (r: SelectionRange) => string
      } | null
      getScrolledVisiblePosition: (pos: { lineNumber: number; column: number }) => { top: number; left: number } | null
      getDomNode: () => HTMLElement | null
    } | null

    if (!editor) return

    const selection = editor.getSelection()
    if (!selection || (selection.startLineNumber === selection.endLineNumber && selection.startColumn === selection.endColumn)) {
      hasSelection.value = false
      selectedText.value = ''
      selectionRange.value = null
      toolbarVisible.value = false
      return
    }

    const model = editor.getModel()
    if (!model) return

    const text = model.getValueInRange(selection)
    if (!text || !text.trim()) {
      hasSelection.value = false
      selectedText.value = ''
      selectionRange.value = null
      toolbarVisible.value = false
      return
    }

    hasSelection.value = true
    selectedText.value = text
    selectionRange.value = selection

    // 计算工具栏位置
    const coords = editor.getScrolledVisiblePosition({
      lineNumber: selection.startLineNumber,
      column: selection.startColumn,
    })
    if (!coords) {
      toolbarVisible.value = false
      return
    }

    const domNode = editor.getDomNode()
    if (!domNode) {
      toolbarVisible.value = false
      return
    }

    const rect = domNode.getBoundingClientRect()
    toolbarPosition.value = {
      top: rect.top + coords.top + window.scrollY,
      left: rect.left + coords.left + window.scrollX,
    }

    toolbarVisible.value = true
  }

  /**
   * Performs an AI operation on the selected text.
   * Sets up streaming listener, calls the backend, and returns streamed content.
   *
   * @param operation - The operation type (polish/translate/summarize/rewrite)
   * @returns The streamed content after completion
   */
  async function performAIOperation(operation: AIOperationType): Promise<string> {
    if (!hasSelection.value || !selectedText.value || !selectionRange.value) {
      return ''
    }

    // Generate a unique operation ID
    currentOperationId = crypto.randomUUID()
    const prompt = OPERATION_PROMPTS[operation](selectedText.value, jobTarget)

    isLoading.value = true
    currentOperation.value = operation

    // Set up streaming listener
    const { content, isStreaming, abort } = useAIStream(currentOperationId)

    // Track streamed content for replacement
    let streamedResult = ''
    const unwatch = setInterval(() => {
      streamedResult = content.value
    }, 50) // ~20fps for streaming preview

    try {
      // Call the backend - this starts the SSE stream
      await sendMessage(currentOperationId, prompt, jobTarget, false)

      // Wait for streaming to complete
      await new Promise<void>((resolve) => {
        const checkInterval = setInterval(() => {
          const streaming = isStreaming as { value: boolean }
          if (!streaming.value) {
            clearInterval(checkInterval)
            clearInterval(unwatch)
            resolve()
          }
        }, 100)
      })

      return content.value
    } catch (err) {
      console.error('AI operation failed:', err)
      abort()
      clearInterval(unwatch)
      return ''
    } finally {
      isLoading.value = false
      currentOperation.value = ''
    }
  }

  /**
   * Hides the toolbar and resets selection state.
   */
  function hideToolbar() {
    toolbarVisible.value = false
  }

  return {
    // Selection state
    hasSelection,
    selectedText,
    selectionRange,
    toolbarPosition,
    toolbarVisible,
    // AI operation state
    isLoading,
    currentOperation,
    // Methods
    updateSelection,
    performAIOperation,
    hideToolbar,
  }
}
