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
 * @param operationId - Unique identifier for the streaming operation
 * @returns Streaming state and control functions
 */
export function useAIStream(operationId: string) {
  const content = ref('')
  const isStreaming = ref(false)
  const error = ref<string | null>(null)
  const aiError = ref<AIError | null>(null)

  let debounceTimer: ReturnType<typeof setTimeout> | null = null
  const DEBOUNCE_MS = 16 // ~60fps for smooth typewriter effect

  const eventName = `ai:stream:${operationId}`

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
        console.error('[AI Stream Error]', aiError.value, { operationId })
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

  // Start listening to events
  EventsOn(eventName, handleEvent)
  isStreaming.value = true

  /**
   * Abort the streaming operation.
   * Implements AIAI-13: user-initiated abort with content preservation.
   *
   * 1. Notifies backend to cancel the operation
   * 2. Stops listening to events
   * 3. Preserves already-generated content (content.value stays intact)
   * 4. Logs the abort event
   */
  const abort = async () => {
    try {
      // Notify backend to cancel the operation
      await AICancelOperation(operationId)
    } catch (err) {
      // Backend cancellation failed - non-fatal, continue with frontend cleanup
      console.warn('[AI] Backend cancel failed:', err)
    }

    // Stop listening to events
    EventsOff(eventName)
    isStreaming.value = false

    if (debounceTimer) {
      clearTimeout(debounceTimer)
      debounceTimer = null
    }

    // Set abort error
    error.value = 'cancelled'
    aiError.value = {
      code: 'ABORTED',
      message: `操作已取消，已保留 ${content.value.length} 字内容`,
      recoverable: true,
      retryable: false,
    }

    // Log abort with content preservation info
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
    EventsOff(eventName)
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
