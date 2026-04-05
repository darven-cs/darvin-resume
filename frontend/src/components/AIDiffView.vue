<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="ai-diff-overlay"
    >
      <div class="ai-diff-view" :style="diffStyle">
        <!-- Header -->
        <div class="diff-header">
          <span class="diff-title">AI 修改对比</span>
          <span class="diff-operation">{{ operationLabel }}</span>
          <div class="diff-actions">
            <button
              class="diff-btn diff-btn-accept"
              :disabled="streaming"
              @click="handleAccept"
            >
              {{ streaming ? '生成中...' : '接受' }}
            </button>
            <button
              class="diff-btn diff-btn-reject"
              @click="handleReject"
            >
              拒绝
            </button>
          </div>
        </div>

        <!-- Diff Content -->
        <div class="diff-content" ref="contentRef">
          <div
            v-for="(line, index) in diffLines"
            :key="index"
            :class="getLineClass(line.type)"
          >
            <span class="line-prefix">{{ getLinePrefix(line.type) }}</span>
            <span class="line-text">{{ line.content }}</span>
          </div>
          <div v-if="diffLines.length === 0 && streaming" class="diff-loading">
            AI 正在生成修改内容...
          </div>
          <div v-if="diffLines.length === 0 && !streaming" class="diff-empty">
            暂无差异内容
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { diffLines as computeDiffLines, Change } from 'diff'
import type { AIOperationType } from '../types/ai'
import { OPERATION_LABELS } from '../composables/useAISelection'

interface DiffLine {
  type: 'added' | 'removed' | 'unchanged'
  content: string
}

interface Props {
  visible: boolean
  originalText: string
  modifiedText: string
  operationType: AIOperationType
  streaming: boolean
  position: { top: number; left: number }
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'accept'): void
  (e: 'reject'): void
}>()

const contentRef = ref<HTMLElement | null>(null)

// 操作类型标签
const operationLabel = computed(() => {
  return OPERATION_LABELS[props.operationType] || 'AI 修改'
})

// 计算 diff 行
const diffLines = computed<DiffLine[]>(() => {
  if (!props.originalText && !props.modifiedText) return []

  const changes: Change[] = computeDiffLines(props.originalText, props.modifiedText)
  const result: DiffLine[] = []

  for (const change of changes) {
    const lines = change.value.replace(/\n$/, '').split('\n')
    for (const line of lines) {
      if (line === '' && !change.value) continue
      result.push({
        type: change.added ? 'added' : change.removed ? 'removed' : 'unchanged',
        content: line,
      })
    }
  }

  return result
})

// 定位：在工具栏下方显示
const DIFF_WIDTH = 480
const DIFF_MAX_HEIGHT = 320

const diffStyle = computed(() => {
  let top = props.position.top + 44 // 工具栏高度 + 间距
  let left = props.position.left

  // 右边界溢出处理
  if (left + DIFF_WIDTH > window.innerWidth - 16) {
    left = window.innerWidth - DIFF_WIDTH - 16
  }
  // 左边界溢出处理
  if (left < 16) {
    left = 16
  }
  // 底部溢出处理
  if (top + DIFF_MAX_HEIGHT > window.innerHeight - 16) {
    top = Math.max(16, props.position.top - DIFF_MAX_HEIGHT - 8)
  }

  return {
    top: `${top}px`,
    left: `${left}px`,
    maxWidth: `${DIFF_WIDTH}px`,
  }
})

// 行样式类名
function getLineClass(type: DiffLine['type']): string {
  switch (type) {
    case 'added': return 'diff-line diff-line-added'
    case 'removed': return 'diff-line diff-line-removed'
    case 'unchanged': return 'diff-line diff-line-unchanged'
  }
}

// 行前缀符号
function getLinePrefix(type: DiffLine['type']): string {
  switch (type) {
    case 'added': return '+'
    case 'removed': return '-'
    case 'unchanged': return ' '
  }
}

function handleAccept() {
  if (props.streaming) return
  emit('accept')
}

function handleReject() {
  emit('reject')
}
</script>

<style scoped>
.ai-diff-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 10002;
  pointer-events: none;
}

.ai-diff-view {
  position: fixed;
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  border-left: 3px solid #409eff;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.6);
  overflow: hidden;
  pointer-events: auto;
  display: flex;
  flex-direction: column;
}

/* Header */
.diff-header {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
  gap: 8px;
}

.diff-title {
  font-size: 13px;
  font-weight: 600;
  color: #e0e0e0;
}

.diff-operation {
  font-size: 11px;
  color: #409eff;
  background: rgba(64, 158, 255, 0.15);
  padding: 1px 6px;
  border-radius: 3px;
}

.diff-actions {
  margin-left: auto;
  display: flex;
  gap: 6px;
}

.diff-btn {
  padding: 4px 14px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s, opacity 0.15s;
}

.diff-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.diff-btn-accept {
  background: #409eff;
  color: #ffffff;
}

.diff-btn-accept:hover:not(:disabled) {
  background: #337ecc;
}

.diff-btn-reject {
  background: #3c3c3c;
  color: #cccccc;
  border: 1px solid #4c4c4c;
}

.diff-btn-reject:hover {
  background: #4c4c4c;
  color: #ffffff;
}

/* Diff Content */
.diff-content {
  max-height: 280px;
  overflow-y: auto;
  padding: 4px 0;
  font-family: 'JetBrains Mono', 'Fira Code', Consolas, monospace;
  font-size: 13px;
  line-height: 1.5;
}

.diff-content::-webkit-scrollbar {
  width: 6px;
}

.diff-content::-webkit-scrollbar-thumb {
  background: #4c4c4c;
  border-radius: 3px;
}

.diff-content::-webkit-scrollbar-track {
  background: transparent;
}

/* Diff Lines */
.diff-line {
  display: flex;
  align-items: flex-start;
  padding: 1px 12px;
  min-height: 22px;
}

.line-prefix {
  flex-shrink: 0;
  width: 16px;
  font-weight: 600;
  user-select: none;
  text-align: center;
}

.line-text {
  white-space: pre-wrap;
  word-break: break-all;
}

.diff-line-added {
  background: rgba(0, 255, 0, 0.08);
}

.diff-line-added .line-prefix {
  color: #4caf50;
}

.diff-line-added .line-text {
  color: #b5e8b5;
}

.diff-line-removed {
  background: rgba(255, 0, 0, 0.08);
}

.diff-line-removed .line-prefix {
  color: #ef5350;
}

.diff-line-removed .line-text {
  color: #f5b7b7;
}

.diff-line-unchanged {
  opacity: 0.5;
}

.diff-line-unchanged .line-prefix {
  color: #6e7681;
}

.diff-line-unchanged .line-text {
  color: #8b949e;
}

/* Loading & Empty */
.diff-loading,
.diff-empty {
  padding: 20px 12px;
  text-align: center;
  color: #6e7681;
  font-size: 13px;
}
</style>
