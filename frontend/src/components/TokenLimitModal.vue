<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-box" role="dialog" aria-modal="true">
          <div class="modal-header">
            <div class="modal-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2L22 20H2L12 2Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
                <path d="M12 9v5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                <circle cx="12" cy="16.5" r="0.75" fill="currentColor"/>
              </svg>
            </div>
            <h2 class="modal-title">{{ title }}</h2>
          </div>
          <div class="modal-body">
            <p class="modal-content">{{ content }}</p>
            <ul class="suggestions-list">
              <li v-for="(suggestion, index) in suggestions" :key="index">
                {{ suggestion }}
              </li>
            </ul>
          </div>
          <div class="modal-footer">
            <button
              v-for="action in actions"
              :key="action.label"
              class="modal-btn"
              :class="{ 'modal-btn-primary': action.action === 'focusEditor' || action.action === 'splitProcess' }"
              @click="handleAction(action.action)"
            >
              {{ action.label }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Action {
  label: string
  action: string
}

interface Props {
  visible: boolean
  title?: string
  content?: string
  suggestions?: string[]
  actions?: Action[]
}

const props = withDefaults(defineProps<Props>(), {
  title: '内容过长',
  content: '简历内容超出了模型处理的 Token 限制。',
  suggestions: () => [
    '1. 精简项目描述，删除重复内容',
    '2. 移除不必要的细节描述',
    '3. 适当缩写过长的公司/项目名称',
    '4. 拆分处理，分批优化各部分',
  ],
  actions: () => [
    { label: '手动精简', action: 'focusEditor' },
    { label: '分批处理', action: 'splitProcess' },
  ],
})

const emit = defineEmits<{
  action: [action: string]
  close: []
}>()

function handleAction(action: string) {
  emit('action', action)
}

function handleClose() {
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--ui-overlay-bg);
  backdrop-filter: blur(2px);
}

.modal-box {
  background: var(--ui-bg-secondary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  box-shadow: var(--ui-shadow-lg);
  width: 90%;
  max-width: 440px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 20px 0;
}

.modal-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  color: var(--ui-warning);
}

.modal-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--ui-text-inverse);
  margin: 0;
}

.modal-body {
  padding: 16px 20px;
}

.modal-content {
  color: var(--ui-text-secondary);
  font-size: 13px;
  line-height: 1.6;
  margin: 0 0 12px;
}

.suggestions-list {
  margin: 0;
  padding: 0 0 0 20px;
  color: var(--ui-text-secondary);
  font-size: 13px;
  line-height: 1.8;
}

.suggestions-list li {
  margin: 0;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 0 20px 20px;
}

.modal-btn {
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

.modal-btn:hover {
  background: var(--ui-border);
  border-color: var(--ui-text-tertiary);
}

.modal-btn-primary {
  background: var(--ui-accent);
  border-color: var(--ui-accent);
  color: var(--ui-text-inverse);
}

.modal-btn-primary:hover {
  background: var(--ui-accent-hover);
  border-color: var(--ui-accent-hover);
}

/* Transition */
.modal-enter-active {
  animation: modal-in var(--ui-transition-normal) ease-out;
}

.modal-leave-active {
  animation: modal-out var(--ui-transition-fast) ease-in;
}

@keyframes modal-in {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes modal-out {
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
