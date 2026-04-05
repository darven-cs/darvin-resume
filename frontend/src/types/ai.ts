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
 * Error codes matching backend AI error types.
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

export type AIErrorCode = typeof AIErrorCodes[keyof typeof AIErrorCodes]

/**
 * Maps backend error codes to user-friendly messages.
 */
export const AIErrorMessages: Record<AIErrorCode, string> = {
  [AIErrorCodes.Network]: '网络连接失败，请检查网络设置',
  [AIErrorCodes.Auth]: 'API 密钥无效或已过期',
  [AIErrorCodes.RateLimit]: '请求过于频繁，请稍后再试',
  [AIErrorCodes.API]: 'AI 服务返回错误',
  [AIErrorCodes.Timeout]: '请求超时，请增加超时时间或减少内容',
  [AIErrorCodes.Cancelled]: '操作已取消',
  [AIErrorCodes.BuildRequest]: '请求构建失败',
  [AIErrorCodes.ParseResponse]: '响应解析失败',
}
