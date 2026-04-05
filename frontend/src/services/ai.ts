/**
 * AI Service Layer
 *
 * Bridges frontend calls to backend Wails methods.
 * All AI operations go through the backend to keep API keys secure.
 */

import {
  GetAIConfig,
  SaveAIConfig,
  ValidateAPIKey,
  AISendMessage,
  AISendMessageSync,
  GetChatHistory,
  SaveChatMessage,
  ClearChatHistory,
  AISendChatMessage,
} from '../wailsjs/wailsjs/go/main/App'
import type { AIConfig, ChatMessage } from '../types/ai'
import { ai } from '../wailsjs/wailsjs/go/models'

/**
 * Retrieves the current AI configuration from the backend.
 */
export async function getConfig(): Promise<AIConfig> {
  const config = await GetAIConfig()
  return {
    apiKey: (config['apiKey'] as string) ?? '',
    baseURL: (config['baseURL'] as string) ?? 'https://api.anthropic.com',
    defaultModel: (config['defaultModel'] as string) ?? 'claude-sonnet-4-20250514',
    maxTokens: (config['maxTokens'] as number) ?? 4096,
    timeoutSeconds: (config['timeoutSeconds'] as number) ?? 60,
  }
}

/**
 * Saves AI configuration to the backend.
 */
export async function saveConfig(config: Partial<AIConfig>): Promise<void> {
  const current = await getConfig()
  const merged: Record<string, unknown> = {
    apiKey: config.apiKey ?? current.apiKey,
    baseURL: config.baseURL ?? current.baseURL,
    defaultModel: config.defaultModel ?? current.defaultModel,
    maxTokens: config.maxTokens ?? current.maxTokens,
    timeoutSeconds: config.timeoutSeconds ?? current.timeoutSeconds,
  }
  await SaveAIConfig(merged)
}

/**
 * Validates an API key by sending a test request to the backend.
 */
export async function validateAPIKey(apiKey: string, baseURL?: string): Promise<boolean> {
  const config = baseURL ? { baseURL } : await getConfig()
  return ValidateAPIKey(apiKey, config.baseURL)
}

/**
 * Sends a streaming chat message to the AI service.
 * Chunks are delivered via Wails EventsOn('ai:stream:{operationId}').
 *
 * @param operationId - Unique ID for tracking this operation (use uuid)
 * @param prompt - User input or operation content
 * @param jobTarget - Target job position for context
 * @param includeFullContext - Whether to include full resume as context
 * @returns Full response content (accumulated after streaming completes)
 */
export async function sendMessage(
  operationId: string,
  prompt: string,
  jobTarget: string = '',
  includeFullContext: boolean = false
): Promise<string> {
  return AISendMessage(operationId, prompt, jobTarget, includeFullContext)
}

/**
 * Sends a non-streaming chat message as a fallback.
 * Use when streaming is not needed or for quick operations.
 */
export async function sendMessageSync(
  operationId: string,
  prompt: string,
  jobTarget: string = '',
  includeFullContext: boolean = false
): Promise<string> {
  return AISendMessageSync(operationId, prompt, jobTarget, includeFullContext)
}

// ============================================================
// Chat History Service
// ============================================================

/**
 * Retrieves chat history for a given resume from the backend.
 * Returns an empty array if no history exists.
 */
export async function getChatHistory(resumeId: string): Promise<ChatMessage[]> {
  const messages = await GetChatHistory(resumeId)
  return (messages || []).map((msg) => {
    const m = msg instanceof ai.ChatMessage ? msg : new ai.ChatMessage(msg)
    return {
      id: m.id || '',
      resumeId: m.resumeId || resumeId,
      role: (m.role === 'user' || m.role === 'assistant') ? m.role : 'user',
      content: m.content || '',
      timestamp: m.createdAt || Date.now(),
      quotedText: m.quotedText || undefined,
    }
  })
}

/**
 * Saves a chat message to the backend for persistence.
 */
export async function saveChatMessage(message: ChatMessage): Promise<void> {
  await SaveChatMessage(new ai.ChatMessage({
    id: message.id,
    resumeId: message.resumeId,
    role: message.role,
    content: message.content,
    quotedText: message.quotedText || '',
    createdAt: message.timestamp,
  }))
}

/**
 * Clears all chat history for a given resume.
 */
export async function clearChatHistory(resumeId: string): Promise<void> {
  await ClearChatHistory(resumeId)
}

/**
 * Sends a chat message with conversation history context.
 * Streaming chunks are delivered via Wails EventsOn('ai:stream:{operationId}').
 *
 * @param operationId - Unique ID for this operation
 * @param prompt - User message content
 * @param jobTarget - Target job position for context
 * @param historyMessages - Recent conversation history (up to 10 messages)
 * @param resumeContent - Current resume markdown content for context
 * @returns Full AI response content
 */
export async function sendChatMessage(
  operationId: string,
  prompt: string,
  jobTarget: string = '',
  historyMessages: ChatMessage[] = [],
  resumeContent: string = ''
): Promise<string> {
  const history = historyMessages.map((msg) => new ai.ChatMessage({
    id: msg.id,
    resumeId: msg.resumeId,
    role: msg.role,
    content: msg.content,
    quotedText: msg.quotedText || '',
    createdAt: msg.timestamp,
  }))
  return AISendChatMessage(operationId, prompt, jobTarget, history, resumeContent)
}
