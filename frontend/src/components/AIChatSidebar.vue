<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'
import type { ChatMessage } from '../types/ai'
import { useAIStream } from '../composables/useAIStream'
import {
  getChatHistory,
  saveChatMessage,
  clearChatHistory,
  sendChatMessage,
} from '../services/ai'

// Props
interface Props {
  visible: boolean
  resumeId: string
  jobTarget: string
  editorContent: string
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'insertText', text: string): void
}>()

// ============================================================
// State
// ============================================================

const messages = ref<ChatMessage[]>([])
const userInput = ref('')
const isStreaming = ref(false)
const streamingContent = ref('')
const quotedText = ref<string | null>(null)
const operationId = ref(crypto.randomUUID())
const isLoading = ref(false)
const isInitialized = ref(false)
const scrollContainer = ref<HTMLElement | null>(null)
const textareaRef = ref<HTMLTextAreaElement | null>(null)
const rafId = ref<number | null>(null)

// Use the AI stream composable with a ref (not a plain string)
// This fixes Bug 4: useAIStream now accepts a ref so it always listens to the current operationId
const { content: streamedContent, isStreaming: streamActive, error: streamError, abort, reset } = useAIStream(operationId)

// Welcome message
const welcomeMessage: ChatMessage = {
  id: 'welcome',
  resumeId: props.resumeId,
  role: 'assistant',
  content: '你好！我是简历优化助手。我可以帮你：\n\n- 优化简历内容的表达\n- 提供简历写作建议\n- 回答关于简历的任何问题\n\n请选择一段简历内容，点击「引用」按钮开始对话。',
  timestamp: Date.now(),
}

// ============================================================
// Computed
// ============================================================

const hasQuotedText = computed(() => quotedText.value !== null)

// ============================================================
// Initialization
// ============================================================

async function initialize() {
  if (isInitialized.value) return
  isLoading.value = true
  try {
    const history = await getChatHistory(props.resumeId)
    if (history.length > 0) {
      messages.value = history
    } else {
      // Show welcome message for new conversations
      messages.value = [{ ...welcomeMessage, resumeId: props.resumeId, timestamp: Date.now() }]
    }
  } catch (err) {
    console.error('Failed to load chat history:', err)
    messages.value = [{ ...welcomeMessage, resumeId: props.resumeId, timestamp: Date.now() }]
  } finally {
    isLoading.value = false
    isInitialized.value = true
  }
}

// Watch visible to initialize on open
watch(() => props.visible, (newVal) => {
  if (newVal) {
    initialize()
    // Reset operation ID for each new session
    operationId.value = crypto.randomUUID()
  }
}, { immediate: true })

// ============================================================
// Streaming content watcher (RAF-based typewriter)
// ============================================================

watch(streamedContent, (newContent) => {
  if (newContent) {
    if (rafId.value) cancelAnimationFrame(rafId.value)
    rafId.value = requestAnimationFrame(() => {
      if (messages.value.length > 0) {
        messages.value[messages.value.length - 1].content = newContent
      }
    })
  }
})

watch(isStreaming, (streaming) => {
  isStreaming.value = streaming
})

// ============================================================
// Auto-scroll
// ============================================================

watch(() => messages.value.length, () => {
  nextTick(() => {
    scrollToBottom()
  })
})

watch(streamedContent, () => {
  nextTick(() => {
    scrollToBottom()
  })
})

function scrollToBottom() {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTo({
      top: scrollContainer.value.scrollHeight,
      behavior: 'smooth',
    })
  }
}

// ============================================================
// Send message
// ============================================================

async function handleSend() {
  const text = userInput.value.trim()
  if (!text || isStreaming.value) return

  // Build the message content (with quoted text if any)
  let fullContent = text
  if (quotedText.value) {
    fullContent = `【引用内容】\n${quotedText.value}\n\n【我的问题】\n${text}`
  }

  // Create user message
  const userMsg: ChatMessage = {
    id: crypto.randomUUID(),
    resumeId: props.resumeId,
    role: 'user',
    content: text,
    timestamp: Date.now(),
    quotedText: quotedText.value || undefined,
  }

  // Add user message
  messages.value.push(userMsg)

  // Clear input and quoted text
  userInput.value = ''
  quotedText.value = null

  // Create assistant message placeholder
  const assistantMsg: ChatMessage = {
    id: crypto.randomUUID(),
    resumeId: props.resumeId,
    role: 'assistant',
    content: '',
    timestamp: Date.now(),
  }
  messages.value.push(assistantMsg)

  // Generate new operation ID for this stream
  const opId = crypto.randomUUID()
  operationId.value = opId

  // Reset the stream composable for the new operation
  reset()

  // Start streaming
  isStreaming.value = true
  try {
    // Get recent history for context (up to 10 messages)
    const historyForContext = messages.value
      .filter(m => m.id !== 'welcome' && m.id !== userMsg.id && m.id !== assistantMsg.id)
      .slice(-10)

    // Send the message
    const response = await sendChatMessage(
      opId,
      fullContent,
      props.jobTarget,
      historyForContext,
      props.editorContent,
    )

    // Update the assistant message with the response
    assistantMsg.content = response || streamedContent.value
  } catch (err) {
    console.error('Chat error:', err)
    assistantMsg.content = `抱歉，发生了错误：${err instanceof Error ? err.message : '未知错误'}`
  } finally {
    isStreaming.value = false

    // Save both messages to history
    try {
      await saveChatMessage(userMsg)
      await saveChatMessage(assistantMsg)
    } catch (err) {
      console.error('Failed to save messages:', err)
    }
  }
}

// ============================================================
// Quote selected text
// ============================================================

function handleQuote() {
  emit('insertText', '[引用内容]') // Placeholder - parent will handle actual selection
  // We can't access the Monaco editor selection from here directly
  // The parent component should listen to a 'quote' event
  quotedText.value = '[用户选中的简历内容]'
}

function handleQuoteFromSelection(selection: string) {
  if (selection.trim()) {
    quotedText.value = selection
  }
}

// Expose method for parent to set quoted text
defineExpose({
  setQuotedText: handleQuoteFromSelection,
  clearQuotedText: () => { quotedText.value = null },
})

// ============================================================
// Insert to editor
// ============================================================

function handleInsert(msgContent: string) {
  if (!msgContent.trim()) return
  emit('insertText', msgContent)
}

// ============================================================
// Clear history
// ============================================================

async function handleClear() {
  if (!confirm('确定要清空对话历史吗？')) return
  try {
    await clearChatHistory(props.resumeId)
    messages.value = [{ ...welcomeMessage, resumeId: props.resumeId, timestamp: Date.now() }]
  } catch (err) {
    console.error('Failed to clear chat history:', err)
  }
}

// ============================================================
// Textarea auto-resize
// ============================================================

function autoResizeTextarea() {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
    textareaRef.value.style.height = Math.min(textareaRef.value.scrollHeight, 120) + 'px'
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

// ============================================================
// Close
// ============================================================

function handleClose() {
  if (isStreaming.value) {
    abort()
    isStreaming.value = false
  }
  emit('close')
}
</script>

<template>
  <Transition name="sidebar">
    <div v-if="visible" class="chat-sidebar">
      <!-- Header -->
      <div class="sidebar-header">
        <span class="sidebar-title">AI 助手</span>
        <div class="header-actions">
          <button class="icon-btn" title="清空对话" @click="handleClear">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6" />
              <path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" />
              <path d="M10 11v6M14 11v6" />
              <path d="M9 6V4a1 1 0 011-1h4a1 1 0 011 1v2" />
            </svg>
          </button>
          <button class="icon-btn close-btn" title="关闭" @click="handleClose">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Messages -->
      <div ref="scrollContainer" class="messages-container">
        <div v-if="isLoading" class="loading-state">
          <div class="spinner"></div>
          <span>加载对话历史...</span>
        </div>

        <template v-else>
          <div
            v-for="msg in messages"
            :key="msg.id"
            :class="['message', msg.role]"
          >
            <div class="message-bubble">
              <pre class="message-content">{{ msg.content }}</pre>
              <button
                v-if="msg.role === 'assistant' && msg.content.trim()"
                class="insert-btn"
                title="插入到编辑区"
                @click="handleInsert(msg.content)"
              >
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M16 4h2a2 2 0 012 2v14a2 2 0 01-2 2H6a2 2 0 01-2-2V6a2 2 0 012-2h2" />
                  <rect x="8" y="2" width="8" height="4" rx="1" ry="1" />
                </svg>
                插入
              </button>
            </div>
          </div>

          <div v-if="isStreaming" class="message assistant streaming">
            <div class="message-bubble">
              <pre class="message-content">{{ streamedContent }}<span class="cursor">|</span></pre>
            </div>
          </div>
        </template>
      </div>

      <!-- Quoted text indicator -->
      <div v-if="quotedText" class="quoted-indicator">
        <div class="quoted-content">
          <span class="quoted-label">引用：</span>
          <span class="quoted-text">{{ quotedText.slice(0, 80) }}{{ quotedText.length > 80 ? '...' : '' }}</span>
        </div>
        <button class="clear-quote" @click="quotedText = null" title="取消引用">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>

      <!-- Input area -->
      <div class="input-area">
        <div class="input-wrapper">
          <textarea
            ref="textareaRef"
            v-model="userInput"
            class="message-input"
            placeholder="输入消息... (Shift+Enter 换行，Enter 发送)"
            :disabled="isStreaming"
            rows="1"
            @input="autoResizeTextarea"
            @keydown="handleKeydown"
          ></textarea>
          <button
            class="send-btn"
            :disabled="!userInput.trim() || isStreaming"
            @click="handleSend"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="22" y1="2" x2="11" y2="13" />
              <polygon points="22 2 15 22 11 13 2 9 22 2" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.chat-sidebar {
  position: fixed;
  top: 0;
  right: 0;
  width: 400px;
  height: 100vh;
  background: #ffffff;
  border-left: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  z-index: 100;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.08);
}

/* Slide animation */
.sidebar-enter-active,
.sidebar-leave-active {
  transition: transform 300ms ease-out;
}

.sidebar-enter-from,
.sidebar-leave-to {
  transform: translateX(100%);
}

/* Header */
.sidebar-header {
  height: 48px;
  min-height: 48px;
  padding: 0 12px;
  background: #f8f8f8;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
}

.header-actions {
  display: flex;
  gap: 4px;
}

.icon-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  transition: background-color 0.15s, color 0.15s;
}

.icon-btn:hover {
  background: #e8e8e8;
  color: #1a1a1a;
}

.close-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

/* Messages */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  height: 100%;
  color: #888;
  font-size: 13px;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #e0e0e0;
  border-top-color: #0078d4;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Message bubbles */
.message {
  display: flex;
  max-width: 85%;
}

.message.user {
  align-self: flex-end;
}

.message.assistant {
  align-self: flex-start;
}

.message-bubble {
  position: relative;
  padding: 10px 14px;
  border-radius: 12px;
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
}

.message.user .message-bubble {
  background: #0078d4;
  color: #ffffff;
  border-bottom-right-radius: 4px;
}

.message.assistant .message-bubble {
  background: #f1f3f4;
  color: #1a1a1a;
  border-bottom-left-radius: 4px;
}

.message-content {
  margin: 0;
  white-space: pre-wrap;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
}

/* Streaming cursor */
.cursor {
  display: inline-block;
  animation: blink 0.8s step-end infinite;
  color: inherit;
  opacity: 0.7;
}

@keyframes blink {
  0%, 100% { opacity: 0.7; }
  50% { opacity: 0; }
}

/* Insert button */
.insert-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.06);
  border: none;
  border-radius: 4px;
  font-size: 11px;
  color: inherit;
  cursor: pointer;
  opacity: 0.7;
  transition: opacity 0.15s, background-color 0.15s;
}

.message.user .insert-btn {
  background: rgba(255, 255, 255, 0.2);
}

.message.assistant .insert-btn:hover {
  opacity: 1;
  background: rgba(0, 0, 0, 0.1);
}

.message.user .insert-btn:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.3);
}

/* Quoted text indicator */
.quoted-indicator {
  padding: 8px 16px;
  background: #fff7ed;
  border-top: 1px solid #fed7aa;
  display: flex;
  align-items: center;
  gap: 8px;
}

.quoted-content {
  flex: 1;
  font-size: 12px;
  color: #9a3412;
  overflow: hidden;
  display: flex;
  gap: 4px;
}

.quoted-label {
  font-weight: 600;
  flex-shrink: 0;
}

.quoted-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.clear-quote {
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: #9a3412;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  flex-shrink: 0;
}

.clear-quote:hover {
  background: rgba(0, 0, 0, 0.08);
}

/* Input area */
.input-area {
  padding: 12px 16px;
  border-top: 1px solid #e0e0e0;
  background: #fafafa;
}

.input-wrapper {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}

.message-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #d0d0d0;
  border-radius: 8px;
  font-size: 13px;
  font-family: inherit;
  line-height: 1.4;
  resize: none;
  outline: none;
  background: #ffffff;
  max-height: 120px;
  min-height: 36px;
  transition: border-color 0.15s;
}

.message-input:focus {
  border-color: #0078d4;
}

.message-input:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.send-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #0078d4;
  color: #ffffff;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background-color 0.15s, opacity 0.15s;
}

.send-btn:hover:not(:disabled) {
  background: #0066b8;
}

.send-btn:disabled {
  background: #b0b0b0;
  cursor: not-allowed;
}
</style>
