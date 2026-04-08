<template>
  <div class="home-view">
    <!-- 顶部栏 -->
    <header class="top-bar">
      <div class="top-bar-left">
        <h1 class="app-title">Darvin-Resume</h1>
      </div>
      <div class="top-bar-right">
        <!-- 搜索框 -->
        <div class="search-box">
          <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索简历..."
            class="search-input"
          />
        </div>
        <!-- 排序按钮 -->
        <button class="sort-btn" @click="toggleSort" :title="sortOrder === 'newest' ? '切换为最早优先' : '切换为最新优先'">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path v-if="sortOrder === 'newest'" d="M12 5v14M5 12l7-7 7 7"/>
            <path v-else d="M12 19V5M5 12l7 7 7-7"/>
          </svg>
        </button>
        <!-- 新建按钮 -->
        <button class="create-btn" @click="showCreateModal = true">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          新建简历
        </button>
        <!-- 设置按钮 -->
        <button class="settings-btn" @click="showSettings = true" title="设置">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </button>
        <!-- 备份按钮 -->
        <button class="backup-btn" @click="showBackupManager = true" title="数据备份">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
        </button>
      </div>
    </header>

    <!-- 卡片网格区域 -->
    <main class="resume-grid-area">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <span>加载中...</span>
      </div>

      <!-- 空状态 + 模板 Demo 预览 -->
      <div v-else-if="filteredResumes.length === 0 && !searchQuery" class="empty-area">
        <EmptyState
          title="还没有简历"
          description="选择一个模板快速开始，或创建空白简历"
          action-label="空白简历"
          @action="handleCreateMode('blank')"
        >
          <template #icon>
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
            </svg>
          </template>
        </EmptyState>

        <div class="template-preview-section">
          <h3 class="template-section-title">内置模板</h3>
          <div class="template-preview-grid">
            <div
              v-for="tpl in builtinTemplates"
              :key="tpl.id"
              class="template-preview-card"
              @click="handleCreateFromTemplate(tpl.id)"
            >
              <div class="template-color-bar" :style="{ background: getTemplateColor(tpl.id) }" />
              <div class="template-preview-body">
                <div class="template-preview-line long" />
                <div class="template-preview-line short" />
                <div class="template-preview-line medium" />
                <div class="template-preview-line long" />
                <div class="template-preview-line short" />
              </div>
              <div class="template-preview-footer">
                <span class="template-preview-name">{{ tpl.name }}</span>
                <span class="template-preview-tag">{{ tpl.tag }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 搜索无结果 -->
      <EmptyState
        v-else-if="filteredResumes.length === 0 && searchQuery"
        title="未找到匹配的简历"
        description="尝试调整搜索关键词"
      >
        <template #icon>
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            <line x1="8" y1="11" x2="14" y2="11"/>
          </svg>
        </template>
      </EmptyState>

      <!-- 卡片网格 -->
      <div v-else class="resume-grid">
        <ResumeCard
          v-for="resume in filteredResumes"
          :key="resume.id"
          :resume="resume"
          @open="openEditor"
          @rename="handleRename"
          @duplicate="handleDuplicate"
          @delete="handleDelete"
        />
      </div>
    </main>

    <!-- 回收站折叠区 -->
    <RecycleBinSection @restore="handleRecycleAction" @permanent-delete="handleRecycleAction" />

    <!-- 创建模式选择弹窗 -->
    <CreateModeModal
      :visible="showCreateModal"
      @close="showCreateModal = false"
      @select-mode="handleCreateMode"
    />

    <!-- 设置弹窗 -->
    <SettingsDialog
      :visible="showSettings"
      @close="showSettings = false"
    />
    <!-- 备份管理弹窗 -->
    <BackupManager
      :visible="showBackupManager"
      @close="showBackupManager = false"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { CreateResume, UpdateResume } from '../wailsjs/wailsjs/go/main/App'
import { useResumeList } from '../composables/useResumeList'
import { BUILTIN_TEMPLATES, type TemplateDef } from '../composables/useTemplate'
import ResumeCard from '../components/ResumeCard.vue'
import CreateModeModal from '../components/CreateModeModal.vue'
import RecycleBinSection from '../components/RecycleBinSection.vue'
import SettingsDialog from '../components/SettingsDialog.vue'
import BackupManager from '../components/BackupManager.vue'
import EmptyState from '../components/EmptyState.vue'

const router = useRouter()
const {
  searchQuery,
  sortOrder,
  filteredResumes,
  loading,
  fetchResumes,
  renameResume,
  duplicateResume,
  deleteResume
} = useResumeList()

const showCreateModal = ref(false)
const showSettings = ref(false)
const showBackupManager = ref(false)

// 内置模板列表
const builtinTemplates: TemplateDef[] = BUILTIN_TEMPLATES

// 模板对应的主题色
const TEMPLATE_COLORS: Record<string, string> = {
  minimal: '#007acc',
  'dual-col': '#238636',
  academic: '#6e40c9',
  campus: '#d32f2f',
}

function getTemplateColor(id: string): string {
  return TEMPLATE_COLORS[id] || '#007acc'
}

// 从模板创建简历
async function handleCreateFromTemplate(templateId: string) {
  try {
    const resume = await CreateResume('未命名简历')
    await UpdateResume(resume.id, JSON.stringify({
      markdownContent: BLANK_TEMPLATE,
      jobTarget: '',
      templateId
    }))
    router.push(`/editor/${resume.id}`)
  } catch (err) {
    console.error('从模板创建简历失败:', err)
  }
}

// 页面加载时获取简历列表
onMounted(() => {
  fetchResumes()
})

// 切换排序
function toggleSort() {
  sortOrder.value = sortOrder.value === 'newest' ? 'oldest' : 'newest'
}

// 打开编辑器
function openEditor(id: string) {
  router.push(`/editor/${id}`)
}

// 处理重命名
async function handleRename(id: string, title: string) {
  try {
    await renameResume(id, title)
  } catch {
    // 错误已在 composable 中处理
  }
}

// 处理复制
async function handleDuplicate(id: string) {
  try {
    await duplicateResume(id)
  } catch {
    // 错误已在 composable 中处理
  }
}

// 处理删除
async function handleDelete(id: string) {
  try {
    await deleteResume(id)
  } catch {
    // 错误已在 composable 中处理
  }
}

// 回收站操作后刷新列表
function handleRecycleAction() {
  fetchResumes()
}

// 空白简历模板 per D-02 / RESM-05
const BLANK_TEMPLATE = `# 你的姓名

📞 电话 | ✉️ 邮箱 | 🌐 GitHub

---

## 专业技能

-

## 项目经历

### 项目名称
**角色** | 起始日期 - 结束日期

- 描述

## 自我评价

-`

// 处理创建模式选择
async function handleCreateMode(mode: 'wizard' | 'blank') {
  try {
    if (mode === 'wizard') {
      // AI 引导模式：创建简历后带 wizard 参数跳转
      const resume = await CreateResume('我的简历')
      router.push(`/editor/${resume.id}?wizard=true`)
    } else {
      // 空白页模式 per D-02 / RESM-05
      // 创建新简历后，填入空白模板并跳转
      const resume = await CreateResume('未命名简历')
      // 使用 UpdateResume 填入空白模板内容
      await UpdateResume(resume.id, JSON.stringify({
        markdownContent: BLANK_TEMPLATE,
        jobTarget: ''
      }))
      router.push(`/editor/${resume.id}`)
    }
  } catch (err) {
    console.error('创建简历失败:', err)
  }
}
</script>

<style scoped>
.home-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--ui-bg-primary);
  overflow: hidden;
}

/* 顶部栏 */
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid var(--ui-border);
  flex-shrink: 0;
}

.top-bar-left {
  display: flex;
  align-items: center;
}

.app-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--ui-text-primary);
  margin: 0;
}

.top-bar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 搜索框 */
.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 10px;
  color: var(--ui-text-tertiary);
  pointer-events: none;
}

.search-input {
  padding: 6px 12px 6px 32px;
  font-size: 0.875rem;
  color: var(--ui-text-primary);
  background: var(--ui-input-bg);
  border: 1px solid var(--ui-input-border);
  border-radius: var(--ui-radius-md);
  outline: none;
  width: 200px;
  transition: border-color var(--ui-transition-fast), width var(--ui-transition-fast);
}

.search-input::placeholder {
  color: var(--ui-text-tertiary);
}

.search-input:focus {
  border-color: var(--ui-border-focus);
  width: 260px;
}

/* 排序按钮 */
.sort-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-md);
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-tertiary);
  cursor: pointer;
  transition: border-color var(--ui-transition-fast), color var(--ui-transition-fast);
}

.sort-btn:hover {
  border-color: var(--ui-border-focus);
  color: var(--ui-text-primary);
}

/* 新建按钮 */
.create-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--ui-text-inverse);
  background: var(--ui-success);
  border: none;
  border-radius: var(--ui-radius-md);
  cursor: pointer;
  transition: background-color var(--ui-transition-fast);
}

.create-btn:hover {
  background: var(--ui-success-hover);
}

/* 设置按钮 */
.settings-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-md);
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-tertiary);
  cursor: pointer;
  transition: border-color var(--ui-transition-fast), color var(--ui-transition-fast);
}

.settings-btn:hover {
  border-color: var(--ui-border-focus);
  color: var(--ui-text-primary);
}

/* 备份按钮 */
.backup-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-md);
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-tertiary);
  cursor: pointer;
  transition: border-color var(--ui-transition-fast), color var(--ui-transition-fast);
}

.backup-btn:hover {
  border-color: var(--ui-border-focus);
  color: var(--ui-text-primary);
}

/* 卡片网格区域 */
.resume-grid-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

/* 卡片网格 */
.resume-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

/* 加载状态 */
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: var(--ui-text-tertiary);
  font-size: 0.875rem;
}

/* 空状态区域 */
.empty-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  overflow-y: auto;
}

/* 模板 Demo 预览 */
.template-preview-section {
  width: 100%;
  max-width: 720px;
  padding: 0 var(--ui-spacing-xl);
  margin-top: var(--ui-spacing-xl);
}

.template-section-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--ui-text-secondary);
  margin: 0 0 var(--ui-spacing-md) 0;
  padding-left: var(--ui-spacing-sm);
}

.template-preview-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--ui-spacing-md);
}

@media (max-width: 600px) {
  .template-preview-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.template-preview-card {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  overflow: hidden;
  cursor: pointer;
  transition: border-color var(--ui-transition-fast), box-shadow var(--ui-transition-fast), transform var(--ui-transition-fast);
  background: var(--ui-bg-primary);
}

.template-preview-card:hover {
  border-color: var(--ui-border-focus);
  box-shadow: var(--ui-shadow-md);
  transform: translateY(-2px);
}

.template-color-bar {
  height: 4px;
  flex-shrink: 0;
}

.template-preview-body {
  padding: 12px 10px 8px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.template-preview-line {
  height: 4px;
  border-radius: 2px;
  background: var(--ui-bg-tertiary);
}

.template-preview-line.long {
  width: 100%;
}

.template-preview-line.medium {
  width: 75%;
}

.template-preview-line.short {
  width: 50%;
}

.template-preview-footer {
  padding: 6px 10px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.template-preview-name {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--ui-text-primary);
}

.template-preview-tag {
  font-size: 0.625rem;
  padding: 1px 6px;
  background: var(--ui-bg-secondary);
  border-radius: 8px;
  color: var(--ui-text-tertiary);
}
</style>
