/**
 * Unified AI Error Handling Composable
 *
 * Provides centralized error display, retry logic, and data preservation
 * for all AI operations. Implements AIAI-10 to AIAI-13.
 *
 * - AIAI-10: Network/API errors with retry button
 * - AIAI-11: Token limit with content reduction suggestions
 * - AIAI-12: Format errors with auto-retry (backend) and raw output fallback
 * - AIAI-13: User abort with content preservation
 */

import { ref, computed } from 'vue'
import type { AIError, AIErrorCode } from '../types/ai'

/**
 * Global error toast state (singleton pattern for app-wide error display).
 */
const globalError = ref<AIError | null>(null)
const globalOperationId = ref<string>('')
const globalVisible = ref(false)

/**
 * Maps backend error codes to structured AIError.
 * Used by convertBackendError() to normalize errors from the backend.
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
 * Determines if an error can be retried automatically.
 */
function isRetryable(code: AIErrorCode): boolean {
  return [
    'NETWORK_ERROR',
    'TIMEOUT',
    'RATE_LIMIT',
  ].includes(code)
}

/**
 * Determines if an error is recoverable (data preserved).
 */
function isRecoverable(code: AIErrorCode): boolean {
  return code !== 'AUTH_FAILED'
}

/**
 * Converts a raw backend error string to a structured AIError.
 * Handles both structured errors and plain error messages.
 */
export function convertBackendError(error: unknown): AIError {
  if (!error) {
    return {
      code: 'UNKNOWN',
      message: friendlyMessages.UNKNOWN,
      recoverable: true,
      retryable: false,
    }
  }

  // Already a structured error
  if (typeof error === 'object' && error !== null && 'code' in error) {
    const e = error as Record<string, unknown>
    const code = (e.code as string) || 'UNKNOWN'

    // Check if it's a legacy backend code
    const mappedCode = backendErrorCodeMap[code] || code as AIErrorCode
    const validCodes: AIErrorCode[] = [
      'NETWORK_ERROR', 'TIMEOUT', 'AUTH_FAILED', 'RATE_LIMIT',
      'TOKEN_LIMIT', 'FORMAT_ERROR', 'ABORTED', 'UNKNOWN',
    ]

    if (!validCodes.includes(mappedCode)) {
      return {
        code: 'UNKNOWN',
        message: (e.message as string) || friendlyMessages.UNKNOWN,
        detail: String(error),
        recoverable: true,
        retryable: false,
      }
    }

    return {
      code: mappedCode,
      message: (e.message as string) || friendlyMessages[mappedCode],
      detail: String(error),
      recoverable: isRecoverable(mappedCode),
      retryable: isRetryable(mappedCode),
    }
  }

  // Plain string or unknown error
  const errorStr = String(error)
  const lowerError = errorStr.toLowerCase()

  // Try to detect error type from message
  if (lowerError.includes('network') || lowerError.includes('connection') ||
    lowerError.includes('refused') || lowerError.includes('dial')) {
    return {
      code: 'NETWORK_ERROR',
      message: friendlyMessages.NETWORK_ERROR,
      detail: errorStr,
      recoverable: true,
      retryable: true,
    }
  }
  if (lowerError.includes('timeout') || lowerError.includes('deadline')) {
    return {
      code: 'TIMEOUT',
      message: friendlyMessages.TIMEOUT,
      detail: errorStr,
      recoverable: true,
      retryable: true,
    }
  }
  if (lowerError.includes('auth') || lowerError.includes('unauthorized') ||
    lowerError.includes('401') || lowerError.includes('403')) {
    return {
      code: 'AUTH_FAILED',
      message: friendlyMessages.AUTH_FAILED,
      detail: errorStr,
      recoverable: false,
      retryable: false,
    }
  }
  if (lowerError.includes('rate limit') || lowerError.includes('429')) {
    return {
      code: 'RATE_LIMIT',
      message: friendlyMessages.RATE_LIMIT,
      detail: errorStr,
      recoverable: true,
      retryable: true,
    }
  }
  if (lowerError.includes('token') || lowerError.includes('context_length')) {
    return {
      code: 'TOKEN_LIMIT',
      message: friendlyMessages.TOKEN_LIMIT,
      detail: errorStr,
      recoverable: true,
      retryable: false,
    }
  }
  if (lowerError.includes('cancel')) {
    return {
      code: 'ABORTED',
      message: friendlyMessages.ABORTED,
      detail: errorStr,
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
 * useAIError composable for component-level error handling.
 * Provides reactive error state and handlers.
 */
export function useAIError() {
  const currentError = ref<AIError | null>(null)
  const currentOperationId = ref<string>('')
  const retryCallback = ref<(() => void) | null>(null)

  /**
   * Display an error as a toast notification.
   * Also logs to console for debugging.
   */
  function showError(error: unknown, operationId: string = '') {
    const aiError = convertBackendError(error)

    // Log to console for debugging
    console.error('[AI Error]', aiError, { operationId })

    // Update reactive state
    currentError.value = aiError
    currentOperationId.value = operationId
    globalError.value = aiError
    globalOperationId.value = operationId
    globalVisible.value = true
  }

  /**
   * Clear the current error state.
   */
  function clearError() {
    currentError.value = null
    currentOperationId.value = ''
    globalError.value = null
    globalOperationId.value = ''
    globalVisible.value = false
  }

  /**
   * Set a retry callback for the current error.
   * The callback will be invoked when the user clicks "重试".
   */
  function setRetryCallback(cb: (() => void) | null) {
    retryCallback.value = cb
  }

  /**
   * Invoke the retry callback if available.
   */
  function retry() {
    if (retryCallback.value) {
      retryCallback.value()
    }
  }

  /**
   * Returns true if the current error is serious (requires user action).
   * Serious errors (AUTH_FAILED) don't auto-dismiss.
   */
  const isSerious = computed(() => {
    if (!currentError.value) return false
    return currentError.value.code === 'AUTH_FAILED'
  })

  /**
   * Returns true if the current error allows retry.
   */
  const canRetry = computed(() => {
    if (!currentError.value) return false
    return currentError.value.retryable
  })

  return {
    // State
    currentError,
    currentOperationId,
    // Global (singleton)
    globalError,
    globalOperationId,
    globalVisible,
    // Computed
    isSerious,
    canRetry,
    // Methods
    showError,
    clearError,
    setRetryCallback,
    retry,
    convertBackendError,
  }
}

/**
 * Modal state for serious errors (e.g., TOKEN_LIMIT).
 */
const modalState = ref<{
  visible: boolean
  title: string
  content: string
  actions: Array<{ label: string; action: string }>
}>({
  visible: false,
  title: '',
  content: '',
  actions: [],
})

/**
 * Shows a modal for serious errors that need user action.
 * Used for TOKEN_LIMIT, AUTH_FAILED, and RATE_LIMIT.
 */
export function showErrorModal(
  title: string,
  content: string,
  actions: Array<{ label: string; action: string }> = []
) {
  modalState.value = {
    visible: true,
    title,
    content,
    actions,
  }
}

/**
 * Closes the error modal.
 */
export function closeErrorModal() {
  modalState.value.visible = false
}

export { modalState }

/**
 * Global error modal state (singleton).
 */
export const useGlobalErrorModal = () => modalState
