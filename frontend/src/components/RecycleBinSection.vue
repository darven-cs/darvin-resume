<template>
  <div class="recycle-bin">
    <!-- 标题行：点击展开/折叠 -->
    <div class="recycle-bin-header" @click="toggleExpand">
      <svg
        class="recycle-bin-arrow"
        :class="{ expanded: isExpanded }"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <polyline points="9 18 15 12 9 6"/>
      </svg>
      <svg
        class="recycle-bin-icon"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <polyline points="3 6 5 6 21 6"/>
        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
        <line x1="10" y1="11" x2="10" y2="17"/>
        <line x1="14" y1="11" x2="14" y2="17"/>
      </svg>
      <span class="recycle-bin-title">回收站</span>
      <span v-if="deletedResumes.length > 0" class="recycle-bin-count">
        {{ deletedResumes.length }}
      </span>
    </div>

    <!-- 展开内容 -->
    <Transition name="expand">
      <div v-if="isExpanded" class="recycle-bin-content">
        <!-- 加载中 -->
        <div v-if="loading" class="recycle-bin-loading">
          加载中...
        </div>

        <!-- 回收站为空 -->
        <div v-else-if="deletedResumes.length === 0" class="recycle-bin-empty">
          回收站为空
        </div>

        <!-- 已删除简历列表 -->
        <div
          v-else
          v-for="resume in deletedResumes"
          :key="resume.id"
          class="recycle-bin-item"
        >
          <div class="recycle-bin-item-info">
            <span class="recycle-bin-item-title">{{ resume.title }}</span>
            <span class="recycle-bin-item-date">{{ formatDate(resume.deletedAt || resume.updatedAt) }}</span>
          </div>
          <div class="recycle-bin-item-actions">
            <button class="action-btn restore-btn" @click="handleRestore(resume.id)" title="恢复">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
                <path d="M3 3v5h5"/>
              </svg>
              恢复
            </button>
            <button class="action-btn delete-btn" @click="handlePermanentDelete(resume.id)" title="永久删除">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
              删除
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ListDeletedResumes, RestoreResume, PermanentDeleteResume } from '../wailsjs/wailsjs/go/main/App'
import type { ResumeListItem } from '../types/resume'

// 展开/折叠状态
const isExpanded = ref(false)
const deletedResumes = ref<ResumeListItem[]>([])
const loading = ref(false)

// 切换展开/折叠
async function toggleExpand() {
  isExpanded.value = !isExpanded.value
  // 展开时加载数据
  if (isExpanded.value && deletedResumes.value.length === 0) {
    await fetchDeletedResumes()
  }
}

// 获取已删除简历列表
async function fetchDeletedResumes() {
  loading.value = true
  try {
    const list = await ListDeletedResumes()
    deletedResumes.value = (list || []).map(item => ({
      id: item.id,
      title: item.title,
      updatedAt: (item as any).deletedAt || item.updatedAt || '',
      deletedAt: (item as any).deletedAt || undefined
    }))
  } catch (err) {
    console.error('加载回收站失败:', err)
  } finally {
    loading.value = false
  }
}

// 恢复简历
async function handleRestore(id: string) {
  try {
    await RestoreResume(id)
    // 从列表中移除已恢复的简历
    deletedResumes.value = deletedResumes.value.filter(r => r.id !== id)
    // 通知父组件刷新主列表（通过自定义事件）
    emit('restore')
  } catch (err) {
    console.error('恢复简历失败:', err)
  }
}

// 永久删除
async function handlePermanentDelete(id: string) {
  if (!confirm('永久删除后无法恢复，确定要删除吗？')) return
  try {
    await PermanentDeleteResume(id)
    // 从列表中移除已删除的简历
    deletedResumes.value = deletedResumes.value.filter(r => r.id !== id)
    // 通知父组件刷新主列表（通过自定义事件）
    emit('permanent-delete')
  } catch (err) {
    console.error('永久删除失败:', err)
  }
}

// 格式化删除时间
function formatDate(dateStr: string | undefined): string {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

const emit = defineEmits<{
  (e: 'restore'): void
  (e: 'permanent-delete'): void
}>()
</script>

<style scoped>
.recycle-bin {
  border-top: 1px solid #3c3c3c;
  margin-top: 24px;
}

.recycle-bin-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  cursor: pointer;
  user-select: none;
  transition: background-color 0.15s;
}

.recycle-bin-header:hover {
  background: rgba(255, 255, 255, 0.03);
}

.recycle-bin-arrow {
  color: #8b949e;
  transition: transform 0.2s;
  flex-shrink: 0;
}

.recycle-bin-arrow.expanded {
  transform: rotate(90deg);
}

.recycle-bin-icon {
  color: #8b949e;
  flex-shrink: 0;
}

.recycle-bin-title {
  font-size: 13px;
  font-weight: 500;
  color: #8b949e;
}

.recycle-bin-count {
  font-size: 11px;
  padding: 1px 6px;
  background: #3c3c3c;
  border-radius: 10px;
  color: #8b949e;
  margin-left: 4px;
}

.recycle-bin-content {
  overflow: hidden;
}

.recycle-bin-loading,
.recycle-bin-empty {
  padding: 12px 16px 12px 40px;
  font-size: 12px;
  color: #8b949e;
}

.recycle-bin-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px 8px 40px;
  border-bottom: 1px solid #2d2d2d;
  transition: background-color 0.15s;
}

.recycle-bin-item:last-child {
  border-bottom: none;
}

.recycle-bin-item:hover {
  background: rgba(255, 255, 255, 0.02);
}

.recycle-bin-item-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}

.recycle-bin-item-title {
  font-size: 13px;
  color: #cccccc;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.recycle-bin-item-date {
  font-size: 11px;
  color: #8b949e;
}

.recycle-bin-item-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  margin-left: 12px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  font-size: 12px;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  background: transparent;
  cursor: pointer;
  transition: border-color 0.15s, background-color 0.15s;
}

.restore-btn {
  color: #58a6ff;
  border-color: rgba(88, 166, 255, 0.3);
}

.restore-btn:hover {
  background: rgba(88, 166, 255, 0.1);
  border-color: #58a6ff;
}

.delete-btn {
  color: #f85149;
  border-color: rgba(248, 81, 73, 0.3);
}

.delete-btn:hover {
  background: rgba(248, 81, 73, 0.1);
  border-color: #f85149;
}

/* 展开动画 */
.expand-enter-active,
.expand-leave-active {
  transition: max-height 0.25s ease, opacity 0.2s ease;
  max-height: 500px;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>
