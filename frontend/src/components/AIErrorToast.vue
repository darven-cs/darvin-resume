<template>
  <Teleport to="body">
    <Transition name="toast">
      <div
        v-if="visible"
        class="ai-error-toast"
        :class="toastClass"
        role="alert"
        aria-live="polite"
      >
        <div class="toast-icon">
          <span v-if="severityIcon" v-html="severityIcon"></span>
        </div>
        <div class="toast-content">
          <div class="toast-title">{{ title }}</div>
          <div class="toast-message">{{ displayMessage }}</div>
          <div v-if="showDetail && error.detail" class="toast-detail">
            <code>{{ error.detail }}</code>
          </div>
        </div>
        <div class="toast-actions">
          <button
            v-if="error.retryable"
            class="toast-btn toast-btn-primary"
            @click="handleRetry"
          >
            重试
          </button>
          <button
            v-if="isDevMode"
            class="toast-btn toast-btn-secondary"
            @click="showDetail = !showDetail"
          >
            {{ showDetail ? '收起详情' : '查看详情' }}
          </button>
          <button
            class="toast-btn toast-btn-close"
            @click="handleClose"
            aria-label="关闭"
          >
            <span class="close-icon">&times;</span>
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import type { AIError, AIErrorCode } from '../types/ai'

interface Props {
  error: AIError
  operationId: string
  autoDismiss?: boolean  // 非严重错误自动消失
  dismissDelay?: number  // 自动消失延迟 (ms)
}

const props = withDefaults(defineProps<Props>(), {
  autoDismiss: true,
  dismissDelay: 3000,
})

const emit = defineEmits<{
  retry: [operationId: string]
  close: []
}>()

const visible = ref(false)
const showDetail = ref(false)

// 开发模式检测
const isDevMode = ref(false)
try {
  isDevMode.value = import.meta.env.DEV
} catch {
  isDevMode.value = false
}

// 根据错误类型计算严重程度
const severity = computed<'error' | 'warning' | 'info'>(() => {
  switch (props.error.code) {
    case 'AUTH_FAILED':
      return 'error'
    case 'TOKEN_LIMIT':
    case 'RATE_LIMIT':
      return 'warning'
    case 'NETWORK_ERROR':
    case 'TIMEOUT':
      return 'warning'
    case 'ABORTED':
      return 'info'
    default:
      return 'error'
  }
})

// CSS class for toast styling
const toastClass = computed(() => ({
  'toast-error': severity.value === 'error',
  'toast-warning': severity.value === 'warning',
  'toast-info': severity.value === 'info',
  'toast-retryable': props.error.retryable,
}))

// 标题（基于错误码）
const title = computed(() => {
  const titles: Record<AIErrorCode, string> = {
    NETWORK_ERROR: '网络连接失败',
    TIMEOUT: '请求超时',
    AUTH_FAILED: 'API 密钥无效',
    RATE_LIMIT: '请求过于频繁',
    TOKEN_LIMIT: '内容超出限制',
    FORMAT_ERROR: 'AI 返回格式异常',
    ABORTED: '操作已取消',
    UNKNOWN: '发生错误',
  }
  return titles[props.error.code] || '发生错误'
})

// 用户友好的错误描述（优先使用 error.message，否则使用预定义消息）
const displayMessage = computed(() => {
  if (props.error.message) {
    return props.error.message
  }
  return props.error.message
})

// 严重错误不自动消失
const shouldAutoDismiss = computed(() => {
  return props.autoDismiss && severity.value !== 'error'
})

// SVG 图标（基于严重程度）
const severityIcon = computed(() => {
  if (severity.value === 'error') {
    return `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="10" cy="10" r="9" stroke="currentColor" stroke-width="1.5"/>
      <path d="M10 6v5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      <circle cx="10" cy="13.5" r="0.75" fill="currentColor"/>
    </svg>`
  }
  if (severity.value === 'warning') {
    return `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M10 2L18 17H2L10 2Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
      <path d="M10 8v4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      <circle cx="10" cy="14.5" r="0.75" fill="currentColor"/>
    </svg>`
  }
  // info
  return `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="10" cy="10" r="9" stroke="currentColor" stroke-width="1.5"/>
    <path d="M10 9v5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
    <circle cx="10" cy="6.5" r="0.75" fill="currentColor"/>
  </svg>`
})

let dismissTimer: ReturnType<typeof setTimeout> | null = null

function startDismissTimer() {
  if (dismissTimer) {
    clearTimeout(dismissTimer)
    dismissTimer = null
  }
  if (shouldAutoDismiss.value) {
    dismissTimer = setTimeout(() => {
      handleClose()
    }, props.dismissDelay)
  }
}

function handleRetry() {
  emit('retry', props.operationId)
  handleClose()
}

function handleClose() {
  visible.value = false
  if (dismissTimer) {
    clearTimeout(dismissTimer)
    dismissTimer = null
  }
  emit('close')
}

// 显示时重置状态并启动计时器
watch(() => props.error, (newError) => {
  if (newError && newError.code) {
    showDetail.value = false
    visible.value = true
    startDismissTimer()
  }
}, { immediate: true })

onUnmounted(() => {
  if (dismissTimer) {
    clearTimeout(dismissTimer)
  }
})
</script>

<style scoped>
.ai-error-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 10000;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  min-width: 320px;
  max-width: 480px;
  border-radius: var(--ui-radius-md);
  box-shadow: var(--ui-shadow-md), 0 0 0 1px rgba(255, 255, 255, 0.1);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  font-size: 13px;
  line-height: 1.5;
}

.toast-error {
  background: var(--ui-bg-active);
  border: 1px solid var(--ui-danger);
  color: var(--ui-danger);
}

.toast-warning {
  background: var(--ui-bg-active);
  border: 1px solid var(--ui-warning);
  color: var(--ui-warning);
}

.toast-info {
  background: var(--ui-bg-active);
  border: 1px solid var(--ui-accent);
  color: var(--ui-accent);
}

.toast-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  margin-top: 1px;
}

.toast-error .toast-icon { color: var(--ui-danger); }
.toast-warning .toast-icon { color: var(--ui-warning); }
.toast-info .toast-icon { color: var(--ui-accent); }

.toast-content {
  flex: 1;
  min-width: 0;
}

.toast-title {
  font-weight: 600;
  margin-bottom: 2px;
  color: inherit;
}

.toast-message {
  color: inherit;
  opacity: 0.9;
}

.toast-detail {
  margin-top: 6px;
  padding: 6px 8px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 3px;
  font-size: 11px;
  overflow-x: auto;
}

.toast-detail code {
  color: var(--ui-text-tertiary);
  font-family: 'SF Mono', 'Consolas', monospace;
  white-space: pre-wrap;
  word-break: break-all;
}

.toast-actions {
  display: flex;
  align-items: center;
  gap: var(--ui-radius-sm);
  flex-shrink: 0;
}

.toast-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 4px 10px;
  border: none;
  border-radius: var(--ui-radius-sm);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), opacity var(--ui-transition-fast);
  white-space: nowrap;
}

.toast-btn-primary {
  background: var(--ui-accent);
  color: var(--ui-text-inverse);
}

.toast-btn-primary:hover {
  background: var(--ui-accent-hover);
}

.toast-btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: inherit;
}

.toast-btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

.toast-btn-close {
  background: transparent;
  color: inherit;
  padding: 4px 6px;
  opacity: 0.7;
}

.toast-btn-close:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.1);
}

.close-icon {
  font-size: 16px;
  line-height: 1;
}

/* Transition animations */
.toast-enter-active {
  animation: toast-in var(--ui-transition-fast) ease-out;
}

.toast-leave-active {
  animation: toast-out var(--ui-transition-fast) ease-in;
}

@keyframes toast-in {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes toast-out {
  from {
    opacity: 1;
    transform: translateX(0);
  }
  to {
    opacity: 0;
    transform: translateX(20px);
  }
}
</style>
