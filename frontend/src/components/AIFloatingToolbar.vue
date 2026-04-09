<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="ai-floating-toolbar"
      :style="toolbarStyle"
    >
      <button
        v-for="op in operations"
        :key="op.id"
        class="toolbar-btn"
        :disabled="loading || !selectedText"
        :title="op.tooltip"
        @click="handleOperation(op.id)"
      >
        <span v-if="loading && currentOperation === op.id" class="spinner" />
        <template v-else>
          <span class="btn-icon">{{ op.icon }}</span>
          <span class="btn-label">{{ op.label }}</span>
        </template>
      </button>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface AIOperation {
  id: string
  label: string
  icon: string
  tooltip: string
}

const operations: AIOperation[] = [
  { id: 'polish', label: '润色', icon: '✦', tooltip: '润色选区文字' },
  { id: 'translate', label: '翻译', icon: '🌐', tooltip: '翻译成英文' },
  { id: 'summarize', label: '缩写', icon: '◈', tooltip: '压缩内容' },
  { id: 'rewrite', label: '重写', icon: '↻', tooltip: '重写选区文字' },
]

interface Props {
  visible: boolean
  position: { top: number; left: number }
  selectedText: string
  loading: boolean
  currentOperation: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'operate', operation: string): void
}>()

// 工具栏尺寸
const TOOLBAR_HEIGHT = 36
const TOOLBAR_WIDTH = 200

const toolbarStyle = computed(() => {
  let top = props.position.top
  let left = props.position.left

  // 向上偏移（显示在选区上方）
  top -= TOOLBAR_HEIGHT + 8

  // 右边界溢出处理
  if (left + TOOLBAR_WIDTH > window.innerWidth - 8) {
    left = window.innerWidth - TOOLBAR_WIDTH - 8
  }
  // 左边界溢出处理
  if (left < 8) {
    left = 8
  }
  // 顶部溢出处理（选区在顶部附近）
  if (top < 8) {
    // 改为向下显示
    top = props.position.top + 24
  }

  return {
    top: `${top}px`,
    left: `${left}px`,
  }
})

function handleOperation(operation: string) {
  if (props.loading || !props.selectedText) return
  emit('operate', operation)
}
</script>

<style scoped>
.ai-floating-toolbar {
  position: fixed;
  z-index: 10001;
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 4px 6px;
  background: var(--ui-bg-primary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-md);
  box-shadow: var(--ui-shadow-md);
  min-height: 36px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: var(--ui-radius-sm);
  padding: 4px 10px;
  min-width: 42px;
  height: 28px;
  background: transparent;
  border: none;
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.1s, color 0.1s;
  white-space: nowrap;
}

.toolbar-btn:hover:not(:disabled) {
  background: var(--ui-accent);
  color: var(--ui-text-inverse);
}

.toolbar-btn:active:not(:disabled) {
  background: var(--ui-accent-hover);
}

.toolbar-btn:disabled {
  color: var(--ui-text-tertiary);
  cursor: not-allowed;
}

.btn-icon {
  font-size: 11px;
}

.btn-label {
  font-size: 12px;
}

.spinner {
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid var(--ui-border);
  border-top-color: var(--ui-text-secondary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
