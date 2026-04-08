<template>
  <div class="resume-card" @click="handleClick">
    <!-- 操作按钮组 - hover 时显示 -->
    <div class="card-actions" @click.stop>
      <button class="action-btn" title="重命名" @click="startRename">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
        </svg>
      </button>
      <button class="action-btn" title="复制" @click="handleDuplicate">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
        </svg>
      </button>
      <button class="action-btn action-btn-danger" title="删除" @click="handleDelete">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="3 6 5 6 21 6"/>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
        </svg>
      </button>
    </div>

    <!-- 卡片内容 -->
    <div class="card-body">
      <!-- 标题区域 -->
      <div class="card-title-area">
        <input
          v-if="isRenaming"
          ref="renameInput"
          v-model="renameTitle"
          class="rename-input"
          @keyup.enter="confirmRename"
          @keyup.escape="cancelRename"
          @blur="confirmRename"
          @click.stop
        />
        <h3 v-else class="card-title">{{ resume.title }}</h3>
      </div>
      <!-- 修改时间 -->
      <span class="card-time">{{ relativeTime }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import type { ResumeListItem } from '../types/resume'

const props = defineProps<{
  resume: ResumeListItem
}>()

const emit = defineEmits<{
  (e: 'open', id: string): void
  (e: 'rename', id: string, title: string): void
  (e: 'duplicate', id: string): void
  (e: 'delete', id: string): void
}>()

// 重命名状态
const isRenaming = ref(false)
const renameTitle = ref('')
const renameInput = ref<HTMLInputElement | null>(null)

// 相对时间格式化
const relativeTime = computed(() => {
  const date = new Date(props.resume.updatedAt)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSeconds = Math.floor(diffMs / 1000)
  const diffMinutes = Math.floor(diffSeconds / 60)
  const diffHours = Math.floor(diffMinutes / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffSeconds < 60) return '刚刚'
  if (diffMinutes < 60) return `${diffMinutes}分钟前`
  if (diffHours < 24) return `${diffHours}小时前`
  if (diffDays < 30) return `${diffDays}天前`
  // 超过30天显示日期
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
})

// 点击卡片主体 → 打开编辑器
function handleClick() {
  if (!isRenaming.value) {
    emit('open', props.resume.id)
  }
}

// 开始重命名
function startRename() {
  isRenaming.value = true
  renameTitle.value = props.resume.title
  nextTick(() => {
    renameInput.value?.focus()
    renameInput.value?.select()
  })
}

// 确认重命名
function confirmRename() {
  if (!isRenaming.value) return
  const trimmed = renameTitle.value.trim()
  if (trimmed && trimmed !== props.resume.title) {
    emit('rename', props.resume.id, trimmed)
  }
  isRenaming.value = false
}

// 取消重命名
function cancelRename() {
  isRenaming.value = false
}

// 复制
function handleDuplicate() {
  emit('duplicate', props.resume.id)
}

// 删除（由父组件处理确认）
function handleDelete() {
  emit('delete', props.resume.id)
}
</script>

<style scoped>
.resume-card {
  position: relative;
  height: 160px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  background: var(--ui-bg-primary);
  cursor: pointer;
  overflow: hidden;
  transition: border-color var(--ui-transition-fast), box-shadow var(--ui-transition-fast);
  display: flex;
  flex-direction: column;
}

.resume-card:hover {
  border-color: var(--ui-accent);
  box-shadow: 0 2px 8px rgba(88, 166, 255, 0.15);
}

.card-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: var(--ui-radius-sm);
  opacity: 0;
  transition: opacity var(--ui-transition-fast);
  z-index: 2;
}

.resume-card:hover .card-actions {
  opacity: 1;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: var(--ui-radius-sm);
  background: rgba(255, 255, 255, 0.1);
  color: var(--ui-text-secondary);
  cursor: pointer;
  transition: background var(--ui-transition-fast), color var(--ui-transition-fast);
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: var(--ui-text-inverse);
}

.action-btn-danger:hover {
  background: rgba(248, 81, 73, 0.3);
  color: var(--ui-danger);
}

.card-body {
  flex: 1;
  padding: 16px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-title-area {
  min-height: 24px;
}

.card-title {
  margin: 0;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--ui-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rename-input {
  width: 100%;
  padding: 2px 6px;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--ui-text-primary);
  background: var(--ui-bg-tertiary);
  border: 1px solid var(--ui-accent);
  border-radius: var(--ui-radius-sm);
  outline: none;
  box-sizing: border-box;
}

.card-time {
  font-size: 0.75rem;
  color: var(--ui-text-tertiary);
}
</style>
