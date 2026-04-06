import { ref, watch, onUnmounted } from 'vue'
import { EventsOff, EventsOn } from '../wailsjs/wailsjs/runtime/runtime'
import type { AIStreamChunk, AIError, AIErrorCode } from '../types/ai'
import { AICancelOperation } from '../wailsjs/wailsjs/go/main/App'

/**
 * Maps backend error codes to structured AIErrorCode.
 */
const backendErrorCodeMap: Record<string, AIErrorCode> = {
  network: 'NETWORK_ERROR',
  auth: 'AUTH_FAILED',
  rate_limit: 'RATE_LIMIT',
  api: 'UNKNOWN',
  timeout: 'TIMEOUT',
  cancelled: 'ABORTED',
  build_request: 'FORMAT_ERROR',
  parse_response: 'FORMAT_ERROR',
}

/**
 * User-friendly messages for each error code.
 */
const friendlyMessages: Record<AIErrorCode, string> = {
  NETWORK_ERROR: '网络连接失败，请检查网络后重试',
  TIMEOUT: '请求超时，请增加超时时间或减少内容后重试',
  AUTH_FAILED: 'API 密钥无效或已过期，请在设置中更新',
  RATE_LIMIT: '请求过于频繁，请稍后再试',
  TOKEN_LIMIT: '简历内容超出了模型处理的 Token 限制，请精简后重试',
  FORMAT_ERROR: 'AI 返回格式异常，已保留原始输出，可手动编辑',
  ABORTED: '操作已取消，已生成的内容已保留',
  UNKNOWN: '发生未知错误，请重试',
}

/**
 * Converts a backend error string to a structured AIError.
 */
function convertBackendError(errorStr: string): AIError {
  const lower = errorStr.toLowerCase()

  if (lower.includes('network') || lower.includes('connection') ||
    lower.includes('refused') || lower.includes('dial')) {
    return {
      code: 'NETWORK_ERROR',
      message: friendlyMessages.NETWORK_ERROR,
      recoverable: true,
      retryable: true,
    }
  }
  if (lower.includes('timeout') || lower.includes('deadline')) {
    return {
      code: 'TIMEOUT',
      message: friendlyMessages.TIMEOUT,
      recoverable: true,
      retryable: true,
    }
  }
  if (lower.includes('auth') || lower.includes('unauthorized') ||
    lower.includes('401') || lower.includes('403')) {
    return {
      code: 'AUTH_FAILED',
      message: friendlyMessages.AUTH_FAILED,
      recoverable: false,
      retryable: false,
    }
  }
  if (lower.includes('rate limit') || lower.includes('429')) {
    return {
      code: 'RATE_LIMIT',
      message: friendlyMessages.RATE_LIMIT,
      recoverable: true,
      retryable: true,
    }
  }
  if (lower.includes('token') || lower.includes('context_length')) {
    return {
      code: 'TOKEN_LIMIT',
      message: friendlyMessages.TOKEN_LIMIT,
      recoverable: true,
      retryable: false,
    }
  }
  if (lower.includes('cancel')) {
    return {
      code: 'ABORTED',
      message: friendlyMessages.ABORTED,
      recoverable: true,
      retryable: false,
    }
  }

  return {
    code: 'UNKNOWN',
    message: friendlyMessages.UNKNOWN,
    detail: errorStr,
    recoverable: true,
    retryable: false,
  }
}

/**
 * Composable for handling AI streaming responses via Wails Events.
 * Listens to ai:stream:{operationId} events and accumulates content chunks.
 * Integrates error handling (AIAI-10 to AIAI-13).
 *
 * @param operationIdRef - Ref containing the unique identifier for the streaming operation.
 *                           Pass a ref so this composable always listens to the current ID,
 *                           fixing the bug where a new operation ID was generated per message
 *                           but the listener was registered with the initial ID.
 * @returns Streaming state and control functions
 */
export function useAIStream(operationIdRef: { value: string }) {
  const content = ref('')
  const isStreaming = ref(false)
  const error = ref<string | null>(null)
  const aiError = ref<AIError | null>(null)

  let debounceTimer: ReturnType<typeof setTimeout> | null = null
  const DEBOUNCE_MS = 16 // ~60fps for smooth typewriter effect
  let currentEventName = ''

  // 【关键修复】handleEvent 必须在 watch 之前定义，因为 watch(immediate:true) 会立即执行
  const handleEvent = (data: AIStreamChunk) => {
    switch (data.type) {
      case 'content':
        if (data.content) {
          // Accumulate content (debounced DOM update for performance)
          if (debounceTimer) clearTimeout(debounceTimer)
          debounceTimer = setTimeout(() => {
            // Content is already accumulated via ref
          }, DEBOUNCE_MS)
          content.value += data.content ?? ''
        }
        break
      case 'error': {
        // Convert backend error string to structured AIError
        const errorStr = data.error ?? 'Unknown error'
        error.value = errorStr
        aiError.value = convertBackendError(errorStr)
        isStreaming.value = false
        // Log to console for debugging (AIAI-10)
        console.error('[AI Stream Error]', aiError.value, { operationId: operationIdRef.value })
        break
      }
      case 'done':
        isStreaming.value = false
        // Clear any pending debounce
        if (debounceTimer) {
          clearTimeout(debounceTimer)
          debounceTimer = null
        }
        break
    }
  }

  // 【关键修复】动态注册事件监听：每次 operationIdRef.value 变化时，
  // 取消旧监听器并注册新监听器。解决新消息使用新 ID 时前端仍监听旧 ID 的问题。
  watch(
    () => operationIdRef.value,
    (newId) => {
      const newEventName = `ai:stream:${newId}`

      // 取消旧监听器（如果有）
      if (currentEventName && currentEventName !== newEventName) {
        EventsOff(currentEventName)
        // 重置状态，为新操作做准备
        content.value = ''
        error.value = null
        aiError.value = null
      }

      // 注册新监听器
      currentEventName = newEventName
      EventsOn(currentEventName, handleEvent)
      isStreaming.value = true
    },
    { immediate: true }
  )

  /**
   * Abort the streaming operation.
   * Implements AIAI-13: user-initiated abort with content preservation.
   */
  const abort = async () => {
    try {
      await AICancelOperation(operationIdRef.value)
    } catch (err) {
      console.warn('[AI] Backend cancel failed:', err)
    }

    if (currentEventName) {
      EventsOff(currentEventName)
    }
    isStreaming.value = false

    if (debounceTimer) {
      clearTimeout(debounceTimer)
      debounceTimer = null
    }

    error.value = 'cancelled'
    aiError.value = {
      code: 'ABORTED',
      message: `操作已取消，已保留 ${content.value.length} 字内容`,
      recoverable: true,
      retryable: false,
    }

    console.log('[AI] Aborted, preserved content length:', content.value.length)
  }

  /**
   * Reset the streaming state for reuse.
   */
  const reset = () => {
    content.value = ''
    error.value = null
    aiError.value = null
    isStreaming.value = false
    if (debounceTimer) {
      clearTimeout(debounceTimer)
      debounceTimer = null
    }
  }

  // Cleanup on unmount
  onUnmounted(() => {
    if (currentEventName) {
      EventsOff(currentEventName)
    }
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
  })

  return {
    content,
    isStreaming,
    error,
    aiError,
    abort,
    reset,
  }
}
