import { ref, computed } from 'vue'
import { getConfig, saveConfig, validateAPIKey } from '../services/ai'
import type { AIConfig } from '../types/ai'

/**
 * Composable for managing AI configuration state.
 * Provides loading, validation, and persistence of AI settings.
 */
export function useAIConfig() {
  const config = ref<AIConfig>({
    apiKey: '',
    baseURL: 'https://api.anthropic.com',
    defaultModel: 'claude-sonnet-4-20250514',
    maxTokens: 4096,
    timeoutSeconds: 60,
  })
  const isLoading = ref(false)
  const isSaving = ref(false)
  const isValidating = ref(false)
  const error = ref<string | null>(null)
  const validationError = ref<string | null>(null)

  /**
   * Loads AI configuration from the backend.
   */
  async function loadConfig() {
    isLoading.value = true
    error.value = null
    try {
      const loaded = await getConfig()
      config.value = loaded
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load AI config'
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Saves AI configuration to the backend.
   */
  async function persistConfig() {
    isSaving.value = true
    error.value = null
    try {
      await saveConfig(config.value)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to save AI config'
      throw err
    } finally {
      isSaving.value = false
    }
  }

  /**
   * Validates the API key by sending a test request.
   */
  async function validateKey() {
    isValidating.value = true
    validationError.value = null
    try {
      const isValid = await validateAPIKey(config.value.apiKey, config.value.baseURL)
      if (!isValid) {
        validationError.value = 'API key is invalid or expired'
      }
      return isValid
    } catch (err) {
      validationError.value = err instanceof Error ? err.message : 'Validation failed'
      return false
    } finally {
      isValidating.value = false
    }
  }

  /**
   * Form validation for AI config form.
   */
  const validationErrors = computed(() => {
    const errors: Record<string, string> = {}
    if (!config.value.apiKey.trim()) {
      errors.apiKey = 'API Key is required'
    }
    if (!config.value.baseURL.trim()) {
      errors.baseURL = 'Base URL is required'
    } else {
      try {
        new URL(config.value.baseURL)
      } catch {
        errors.baseURL = 'Invalid URL format'
      }
    }
    if (!config.value.defaultModel.trim()) {
      errors.defaultModel = 'Model name is required'
    }
    if (config.value.maxTokens <= 0) {
      errors.maxTokens = 'Max tokens must be positive'
    }
    if (config.value.maxTokens > 8192) {
      errors.maxTokens = 'Max tokens should be <= 8192'
    }
    if (config.value.timeoutSeconds <= 0) {
      errors.timeoutSeconds = 'Timeout must be positive'
    }
    if (config.value.timeoutSeconds > 300) {
      errors.timeoutSeconds = 'Timeout should be <= 300 seconds'
    }
    return errors
  })

  const isValid = computed(() => Object.keys(validationErrors.value).length === 0)

  /**
   * Checks if AI is configured (has API key).
   */
  const isConfigured = computed(() => config.value.apiKey.trim().length > 0)

  return {
    config,
    isLoading,
    isSaving,
    isValidating,
    error,
    validationError,
    validationErrors,
    isValid,
    isConfigured,
    loadConfig,
    persistConfig,
    validateKey,
  }
}
