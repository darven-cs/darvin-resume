<template>
  <Teleport to="body">
    <div class="toast-container" aria-live="polite" aria-label="通知">
      <TransitionGroup name="toast" tag="div">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="toast-item"
          :class="`toast-${toast.type}`"
          role="alert"
        >
          <!-- 类型图标 -->
          <div class="toast-icon" v-html="getIcon(toast.type)" />

          <!-- 内容 -->
          <div class="toast-content">
            <div class="toast-title">{{ toast.title }}</div>
            <div v-if="toast.message" class="toast-message">{{ toast.message }}</div>
            <button
              v-if="toast.action"
              class="toast-action-btn"
              @click="handleAction(toast)"
            >
              {{ toast.action.label }}
            </button>
          </div>

          <!-- 关闭按钮 -->
          <button
            class="toast-close-btn"
            @click="handleClose(toast.id)"
            aria-label="关闭通知"
          >
            <span>&times;</span>
          </button>

          <!-- 进度条（仅自动消失的 toast 显示） -->
          <div
            v-if="toast.duration > 0"
            class="toast-progress"
            :style="{ animationDuration: toast.duration + 'ms' }"
          />
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast, type ToastType } from '../composables/useToast'

const { toasts, remove } = useToast()

function handleClose(id: string) {
  remove(id)
}

function handleAction(toast: { id: string; action?: { label: string; handler: () => void } }) {
  if (toast.action) {
    try {
      toast.action.handler()
    } catch (e) {
      console.error('[Toast] Action handler error:', e)
    }
  }
  remove(toast.id)
}

function getIcon(type: ToastType): string {
  switch (type) {
    case 'success':
      return `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
        <polyline points="22 4 12 14.01 9 11.01"/>
      </svg>`
    case 'error':
      return `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"/>
        <line x1="15" y1="9" x2="9" y2="15"/>
        <line x1="9" y1="9" x2="15" y2="15"/>
      </svg>`
    case 'warning':
      return `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
        <line x1="12" y1="9" x2="12" y2="13"/>
        <line x1="12" y1="17" x2="12.01" y2="17"/>
      </svg>`
    case 'info':
    default:
      return `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="16" x2="12" y2="12"/>
        <line x1="12" y1="8" x2="12.01" y2="8"/>
      </svg>`
  }
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 10000;
  display: flex;
  flex-direction: column;
  gap: 8px;
  pointer-events: none;  /* 允许点击穿透到下层 */
  max-width: 360px;
  width: 100%;
}

/* 每个 toast 项 */
.toast-item {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  min-width: 280px;
  max-width: 360px;
  border-radius: var(--ui-radius-md);
  box-shadow: var(--ui-shadow-md), 0 0 0 1px rgba(255, 255, 255, 0.05);
  font-family: var(--ui-font-sans);
  font-size: 13px;
  line-height: 1.5;
  pointer-events: auto;   /* 恢复可点击 */
  overflow: hidden;
  border: 1px solid transparent;
}

/* 类型样式 */
.toast-success {
  background: var(--ui-bg-secondary);
  border-color: var(--ui-success);
  color: var(--ui-success);
}

.toast-error {
  background: var(--ui-bg-secondary);
  border-color: var(--ui-danger);
  color: var(--ui-danger);
}

.toast-warning {
  background: var(--ui-bg-secondary);
  border-color: var(--ui-warning);
  color: var(--ui-warning);
}

.toast-info {
  background: var(--ui-bg-secondary);
  border-color: var(--ui-info);
  color: var(--ui-info);
}

/* 图标 */
.toast-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 1px;
}

/* 内容 */
.toast-content {
  flex: 1;
  min-width: 0;
}

.toast-title {
  font-weight: 600;
  font-size: 13px;
  color: var(--ui-text-primary);
  line-height: 1.4;
}

.toast-message {
  font-size: 12px;
  color: var(--ui-text-secondary);
  margin-top: 2px;
  line-height: 1.4;
  word-break: break-word;
}

.toast-action-btn {
  margin-top: 6px;
  padding: 3px 8px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  background: transparent;
  color: var(--ui-accent);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), border-color var(--ui-transition-fast);
}

.toast-action-btn:hover {
  background: var(--ui-bg-hover);
  border-color: var(--ui-accent);
}

/* 关闭按钮 */
.toast-close-btn {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border: none;
  border-radius: var(--ui-radius-sm);
  background: transparent;
  color: var(--ui-text-tertiary);
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  padding: 0;
  opacity: 0.7;
  transition: opacity var(--ui-transition-fast), background-color var(--ui-transition-fast);
}

.toast-close-btn:hover {
  opacity: 1;
  background: var(--ui-bg-hover);
}

/* 进度条 */
.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 2px;
  background: currentColor;
  opacity: 0.5;
  animation: toast-progress linear forwards;
  width: 100%;
  transform-origin: left;
}

@keyframes toast-progress {
  from {
    transform: scaleX(1);
  }
  to {
    transform: scaleX(0);
  }
}

/* 过渡动画 */
.toast-enter-active {
  animation: toast-slide-in 200ms ease-out;
}

.toast-leave-active {
  animation: toast-slide-out 200ms ease-in;
}

.toast-move {
  transition: transform 200ms ease;
}

@keyframes toast-slide-in {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes toast-slide-out {
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
