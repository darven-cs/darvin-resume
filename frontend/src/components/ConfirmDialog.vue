<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="dialogState.visible" class="confirm-overlay" @click.self="onCancel">
        <div class="confirm-box" role="dialog" aria-modal="true" aria-labelledby="confirm-dialog-title">
          <!-- 标题 -->
          <div class="confirm-header">
            <span id="confirm-dialog-title" class="confirm-title">{{ dialogState.title }}</span>
          </div>

          <!-- 内容 -->
          <div class="confirm-body">
            <p class="confirm-message">{{ dialogState.message }}</p>
          </div>

          <!-- 按钮 -->
          <div class="confirm-footer">
            <button class="confirm-btn confirm-btn-cancel" @click="onCancel">
              {{ dialogState.cancelLabel }}
            </button>
            <button
              class="confirm-btn"
              :class="dialogState.type === 'danger' ? 'confirm-btn-danger' : 'confirm-btn-primary'"
              @click="onConfirm"
            >
              {{ dialogState.confirmLabel }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useConfirm } from '../composables/useConfirm'

const { dialogState, onConfirm, onCancel } = useConfirm()
</script>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  z-index: 10001;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--ui-overlay-bg);
  backdrop-filter: blur(2px);
}

.confirm-box {
  background: var(--ui-bg-secondary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  box-shadow: var(--ui-shadow-lg);
  width: 90%;
  max-width: 400px;
  overflow: hidden;
}

.confirm-header {
  padding: 16px 20px 12px;
  border-bottom: 1px solid var(--ui-border);
}

.confirm-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--ui-text-primary);
}

.confirm-body {
  padding: 16px 20px;
}

.confirm-message {
  font-size: 13px;
  color: var(--ui-text-secondary);
  line-height: 1.6;
  margin: 0;
}

.confirm-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 20px 16px;
}

.confirm-btn {
  padding: 6px 16px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-primary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), border-color var(--ui-transition-fast);
}

.confirm-btn:hover {
  background: var(--ui-border);
  border-color: var(--ui-border-hover);
}

.confirm-btn-primary {
  background: var(--ui-accent);
  border-color: var(--ui-accent);
  color: var(--ui-text-inverse);
}

.confirm-btn-primary:hover {
  background: var(--ui-accent-hover);
  border-color: var(--ui-accent-hover);
}

.confirm-btn-danger {
  background: var(--ui-danger);
  border-color: var(--ui-danger);
  color: var(--ui-text-inverse);
}

.confirm-btn-danger:hover {
  background: var(--ui-danger-hover);
  border-color: var(--ui-danger-hover);
}

/* 过渡动画 */
.modal-enter-active {
  animation: confirm-in var(--ui-transition-fast) ease-out;
}

.modal-leave-active {
  animation: confirm-out 0.15s ease-in;
}

@keyframes confirm-in {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes confirm-out {
  from {
    opacity: 1;
    transform: scale(1);
  }
  to {
    opacity: 0;
    transform: scale(0.95);
  }
}
</style>
