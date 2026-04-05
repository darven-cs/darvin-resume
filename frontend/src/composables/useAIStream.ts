import { ref, onUnmounted } from 'vue'
import { EventsOff, EventsOn } from '../wailsjs/wailsjs/runtime/runtime'
import type { AIStreamChunk } from '../types/ai'

/**
 * Composable for handling AI streaming responses via Wails Events.
 * Listens to ai:stream:{operationId} events and accumulates content chunks.
 *
 * @param operationId - Unique identifier for the streaming operation
 * @returns Streaming state and control functions
 */
export function useAIStream(operationId: string) {
  const content = ref('')
  const isStreaming = ref(false)
  const error = ref<string | null>(null)

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
      case 'error':
        error.value = data.error ?? 'Unknown error'
        isStreaming.value = false
        break
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
   * Abort the streaming operation by stopping event listening.
   * The backend will continue until context cancellation is implemented.
   */
  const abort = () => {
    EventsOff(eventName)
    isStreaming.value = false
    if (debounceTimer) {
      clearTimeout(debounceTimer)
      debounceTimer = null
    }
  }

  /**
   * Reset the streaming state for reuse.
   */
  const reset = () => {
    content.value = ''
    error.value = null
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
    abort,
    reset,
  }
}
