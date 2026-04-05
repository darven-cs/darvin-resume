// AI TypeScript type definitions

/**
 * Represents a streaming chunk from the AI service.
 * Emitted via Wails EventsEmit("ai:stream:{operationId}", chunk).
 */
export interface AIStreamChunk {
  type: 'content' | 'error' | 'done'
  content?: string
  error?: string
}

/**
 * AI service configuration stored in the backend settings table.
 */
export interface AIConfig {
  apiKey: string
  baseURL: string
  defaultModel: string
  maxTokens: number
  timeoutSeconds: number
}

/**
 * AI operation types for selection-based actions.
 */
export type AIOperationType = 'polish' | 'translate' | 'summarize' | 'rewrite'

/**
 * Error codes for AI operations (AIAI-10 to AIAI-13).
 * Covers all exception scenarios: network, timeout, auth, rate limit,
 * token limit, format error, abort, and unknown.
 */
export type AIErrorCode =
  | 'NETWORK_ERROR'      // 网络连接失败
  | 'TIMEOUT'            // 请求超时
  | 'AUTH_FAILED'        // API Key 无效/过期
  | 'RATE_LIMIT'         // 请求频率超限
  | 'TOKEN_LIMIT'        // Token 超限
  | 'FORMAT_ERROR'       // 返回格式异常
  | 'ABORTED'            // 用户主动中断
  | 'UNKNOWN'            // 未知错误

/**
 * AI error structure with user-friendly message and recovery metadata.
 */
export interface AIError {
  code: AIErrorCode
  message: string        // 用户友好的错误描述
  detail?: string        // 详细技术信息（仅开发调试显示）
  recoverable: boolean   // 是否可恢复
  retryable: boolean     // 是否可重试
}

/**
 * Maps error codes to user-friendly messages.
 */
export const AIErrorMessages: Record<AIErrorCode, string> = {
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
 * Legacy error codes (kept for backward compatibility).
 * @deprecated Use AIErrorCode instead
 */
export const AIErrorCodes = {
  Network: 'network',
  Auth: 'auth',
  RateLimit: 'rate_limit',
  API: 'api',
  Timeout: 'timeout',
  Cancelled: 'cancelled',
  BuildRequest: 'build_request',
  ParseResponse: 'parse_response',
} as const

/**
 * Maps legacy backend error codes to new error codes.
 * Used by the frontend to convert backend errors to structured AIError.
 */
export const LegacyErrorCodeMap: Record<string, AIErrorCode> = {
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
 * Chat message in the AI conversation sidebar.
 */
export interface ChatMessage {
  id: string
  resumeId: string
  role: 'user' | 'assistant'
  content: string
  timestamp: number
  quotedText?: string  // 引用的选中文本
}

/**
 * Chat conversation container.
 */
export interface ChatConversation {
  id: string
  resumeId: string
  messages: ChatMessage[]
  createdAt: number
}
