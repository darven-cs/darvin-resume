/**
 * AI Service Layer
 *
 * Bridges frontend calls to backend Wails methods.
 * All AI operations go through the backend to keep API keys secure.
 */

import { GetAIConfig, SaveAIConfig, ValidateAPIKey, AISendMessage, AISendMessageSync } from '../wailsjs/wailsjs/go/main/App'
import type { AIConfig } from '../types/ai'

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
