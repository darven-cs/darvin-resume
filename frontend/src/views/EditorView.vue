<template>
  <div class="editor-view">
    <!-- Editor Toolbar -->
    <div class="editor-toolbar">
      <!-- 返回按钮 per D-29 -->
      <button class="toolbar-btn" @click="handleBack" title="返回 (有未保存内容时将提示)">
        <span class="btn-icon">&#8592;</span>
        <span class="btn-label">返回</span>
      </button>
      <div class="toolbar-divider" />

      <!-- 简历标题编辑 per D-30 -->
      <div class="title-edit-wrapper">
        <input
          v-if="isEditingTitle"
          ref="titleInputRef"
          v-model="editingTitle"
          class="title-input"
          @keydown.enter="confirmTitleEdit"
          @keydown.escape="cancelTitleEdit"
          @blur="confirmTitleEdit"
        />
        <button
          v-else
          class="title-display"
          @click="startTitleEdit"
          :title="'点击修改标题: ' + resumeTitle"
        >
          <span class="title-text">{{ resumeTitle }}</span>
          <svg class="title-edit-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
          </svg>
        </button>
      </div>

      <div class="toolbar-divider" />

      <!-- 导入按钮 -->
      <button class="toolbar-btn" @click="showParserModal = true" title="导入旧简历">
        <span class="btn-icon">📥</span>
        <span class="btn-label">导入</span>
      </button>

      <JobTargetChip
        v-model="jobTarget"
        @change="handleJobTargetChange"
      />
      <div class="toolbar-spacer" />

      <!-- 保存状态指示器 per D-25 -->
      <SaveStatusIndicator
        :status="saveStatus"
        :error-message="errorMessage"
        @retry="handleRetrySave"
      />

      <button
        class="toolbar-btn"
        @click="showAIConfigModal = true"
        title="AI 设置"
      >
        <span class="btn-icon">⚙</span>
        <span class="btn-label">设置</span>
      </button>
      <button
        class="toolbar-btn"
        :class="{ active: chatSidebarVisible }"
        @click="chatSidebarVisible = !chatSidebarVisible"
        title="AI 对话 (Ctrl+Shift+A)"
      >
        <span class="btn-icon">💬</span>
        <span class="btn-label">AI 助手</span>
      </button>
    </div>

    <!-- 双栏模式 (窗口宽度 >= 1200px) per D-09 -->
    <template v-if="!isSinglePane">
      <SplitPane :default-ratio="40" :min-width="300">
        <template #left>
          <div class="editor-pane">
            <div class="pane-header">
              <span class="pane-title">编辑</span>
            </div>
            <div class="editor-wrapper">
              <MonacoEditor
                ref="monacoRef"
                v-model="content"
                :job-target="jobTarget"
                @change="handleContentChange"
              />
            </div>
          </div>
        </template>

        <template #right>
          <div class="preview-pane">
            <div class="pane-header">
              <span class="pane-title">预览</span>
            </div>
            <A4Page :content="debouncedContent" />
          </div>
        </template>
      </SplitPane>
    </template>

    <!-- 单栏模式 (窗口宽度 < 1200px) per D-09 -->
    <template v-else>
      <div class="single-pane-mode">
        <div class="view-tabs">
          <button
            :class="{ active: activeView === 'editor' }"
            @click="activeView = 'editor'"
          >
            编辑
          </button>
          <button
            :class="{ active: activeView === 'preview' }"
            @click="activeView = 'preview'"
          >
            预览
          </button>
        </div>

        <div class="single-pane-content">
          <div v-if="activeView === 'editor'" class="editor-wrapper">
            <MonacoEditor
              ref="monacoRef"
              v-model="content"
              :job-target="jobTarget"
              @change="handleContentChange"
            />
          </div>
          <div v-else class="preview-wrapper">
            <A4Page :content="debouncedContent" />
          </div>
        </div>
      </div>
    </template>

    <!-- Resume Parser Modal -->
    <ResumeParserModal
      :visible="showParserModal"
      :resume-id="resumeId"
      @close="showParserModal = false"
      @import="handleImport"
    />

    <!-- AI Chat Sidebar -->
    <AIChatSidebar
      ref="chatSidebarRef"
      :visible="chatSidebarVisible"
      :resume-id="resumeId"
      :job-target="jobTarget"
      :editor-content="content"
      @close="chatSidebarVisible = false"
      @insert-text="handleInsertText"
    />

    <!-- AI Config Modal -->
    <AIConfigModal
      :visible="showAIConfigModal"
      @close="showAIConfigModal = false"
    />

    <!-- Resume Wizard Sidebar -->
    <ResumeWizardSidebar
      :visible="wizardVisible"
      :resume-id="resumeId"
      :job-target="jobTarget"
      @close="wizardVisible = false"
      @complete="handleWizardComplete"
      @save-draft="handleWizardSaveDraft"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import MonacoEditor from '../components/MonacoEditor.vue'
import SplitPane from '../components/SplitPane.vue'
import A4Page from '../components/A4Page.vue'
import JobTargetChip from '../components/JobTargetChip.vue'
import ResumeParserModal from '../components/ResumeParserModal.vue'
import AIChatSidebar from '../components/AIChatSidebar.vue'
import AIConfigModal from '../components/AIConfigModal.vue'
import ResumeWizardSidebar from '../components/ResumeWizardSidebar.vue'
import SaveStatusIndicator from '../components/SaveStatusIndicator.vue'
import { useAutoSave } from '../composables/useAutoSave'
import { GetResume, RenameResume } from '../wailsjs/wailsjs/go/main/App'
import type { Resume } from '../types/resume'

const route = useRoute()
const router = useRouter()
const monacoRef = ref<InstanceType<typeof MonacoEditor> | null>(null)

// 当前简历 ID
const resumeId = computed(() => route.params.id as string)

// 编辑器内容
const content = ref('')
const debouncedContent = ref('')
const jobTarget = ref('')

// 简历标题 per D-30
const resumeTitle = ref('未命名简历')
const isEditingTitle = ref(false)
const editingTitle = ref('')
const titleInputRef = ref<HTMLInputElement | null>(null)

// Modal 状态
const showParserModal = ref(false)
const showAIConfigModal = ref(false)

// AI Chat Sidebar 状态
const chatSidebarVisible = ref(false)
const chatSidebarRef = ref<InstanceType<typeof AIChatSidebar> | null>(null)

// Resume Wizard 状态 — 检测路由 query ?wizard=true 自动打开向导
const wizardVisible = ref(route.query.wizard === 'true')

// 响应式状态 per D-09
const windowWidth = ref(window.innerWidth)
const isSinglePane = computed(() => windowWidth.value < 1200)
const activeView = ref<'editor' | 'preview'>('editor')

// 键盘快捷键处理 — Ctrl+S / Cmd+S 保存
function handleKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    triggerSave()
  }
}

// 自动保存 composable per D-25~D-28
const { saveStatus, errorMessage, markDirty, triggerSave, startAutoSave, stopAutoSave } = useAutoSave({
  resumeId,
  getData: () => ({ markdownContent: content.value, jobTarget: jobTarget.value })
})

// 加载简历数据
onMounted(async () => {
  const id = resumeId.value
  if (id) {
    try {
      const resume = await GetResume(id)
      content.value = resume.markdownContent || ''
      debouncedContent.value = resume.markdownContent || ''
      jobTarget.value = (resume as any).jobTarget || ''
      resumeTitle.value = resume.title || '未命名简历'
    } catch (err) {
      console.error('Failed to load resume:', err)
    }
  }

  // 启动自动保存定时器 per D-26
  startAutoSave()

  // 注册 Ctrl+S / Cmd+S 保存快捷键
  window.addEventListener('keydown', handleKeydown)

  // 监听窗口宽度变化
  window.addEventListener('resize', handleResize)

  // 编辑器与预览滚动同步 per D-11
  setTimeout(setupScrollSync, 500)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  window.removeEventListener('keydown', handleKeydown)
  stopAutoSave()
})

function handleResize() {
  windowWidth.value = window.innerWidth
}

// 内容变化处理 per D-10 (150ms debounce for preview)
// 预览防抖和自动保存是两个独立机制，不互相干扰 per D-10 注释
function handleContentChange(newContent: string) {
  content.value = newContent
  // 标记内容已修改，触发自动保存
  markDirty()
}

// 职位目标变更时保存
async function handleJobTargetChange(value: string) {
  jobTarget.value = value
  markDirty()
}

// 重试保存
function handleRetrySave() {
  triggerSave()
}

// 返回按钮处理 per D-29
function handleBack() {
  // 有未保存内容时弹确认框
  if (saveStatus.value === 'unsaved' || saveStatus.value === 'error') {
    if (confirm('有未保存的内容，确定要离开吗？')) {
      triggerSave().then(() => router.push('/'))
      return
    }
  }
  router.push('/')
}

// 标题编辑 per D-30
function startTitleEdit() {
  editingTitle.value = resumeTitle.value
  isEditingTitle.value = true
  // DOM 更新后聚焦 input
  setTimeout(() => {
    titleInputRef.value?.focus()
    titleInputRef.value?.select()
  }, 0)
}

async function confirmTitleEdit() {
  if (!isEditingTitle.value) return
  isEditingTitle.value = false

  const newTitle = editingTitle.value.trim()
  if (!newTitle || newTitle === resumeTitle.value) return

  try {
    await RenameResume(resumeId.value, newTitle)
    resumeTitle.value = newTitle
  } catch (err) {
    console.error('重命名失败:', err)
  }
}

function cancelTitleEdit() {
  isEditingTitle.value = false
  editingTitle.value = ''
}

// 导入旧简历处理
async function handleImport(markdown: string, importedJobTarget: string) {
  // 填充编辑器内容
  content.value = markdown
  debouncedContent.value = markdown

  // 如果导入了职位目标，更新 chip
  if (importedJobTarget) {
    jobTarget.value = importedJobTarget
  }

  // 触发保存 per D-26 (AI操作完成时)
  await triggerSave()

  // 关闭 modal
  showParserModal.value = false

  // 聚焦编辑器
  monacoRef.value?.focus()
}

// AI Sidebar text insertion
function handleInsertText(text: string) {
  monacoRef.value?.insertAtCursor(text)
}

// 向导完成后重新加载简历内容（后端已自动生成 Markdown）
async function handleWizardComplete() {
  wizardVisible.value = false
  await loadResume()
  monacoRef.value?.focus()
}

// 向导保存草稿后重新加载简历内容
async function handleWizardSaveDraft() {
  wizardVisible.value = false
  await loadResume()
}

// 加载简历数据（提取为可复用函数）
async function loadResume() {
  const id = resumeId.value
  if (!id) return
  try {
    const resume = await GetResume(id)
    content.value = resume.markdownContent || ''
    debouncedContent.value = resume.markdownContent || ''
    jobTarget.value = (resume as any).jobTarget || ''
    resumeTitle.value = resume.title || '未命名简历'
  } catch (err) {
    console.error('Failed to load resume:', err)
  }
}

// 路由离开守卫 per D-31
onBeforeRouteLeave((to, from, next) => {
  if (saveStatus.value === 'unsaved') {
    const confirmed = confirm('有未保存的内容，确定要离开吗？')
    if (confirmed) {
      triggerSave().then(() => next())
      return
    } else {
      next(false)
    }
  } else {
    next()
  }
})

// 150ms 防抖更新预览 per D-10
let debounceTimer: ReturnType<typeof setTimeout> | null = null
watch(content, (newVal) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    debouncedContent.value = newVal
  }, 150)
})

// 编辑器与预览滚动同步 per D-11
function setupScrollSync() {
  const editorInstance = monacoRef.value?.getEditor?.()
  if (!editorInstance) return

  const container = document.querySelector('.preview-container')
  if (!container) return

  editorInstance.onDidScrollChange((e: any) => {
    const scrollTop = e.scrollTop
    const scrollHeight = e.scrollHeight
    const clientHeight = e.height
    if (scrollHeight <= clientHeight) return

    const ratio = scrollTop / (scrollHeight - clientHeight)
    container.scrollTop = ratio * (container.scrollHeight - container.clientHeight)
  })
}
</script>

<style scoped>
.editor-view {
  height: 100%;
  width: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* Toolbar */
.editor-toolbar {
  height: 40px;
  min-height: 40px;
  padding: 0 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #3c3c3c;
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: transparent;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #cccccc;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.15s, border-color 0.15s;
}

.toolbar-btn:hover {
  background: #3c3c3c;
  border-color: #4c4c4c;
}

.toolbar-btn.active {
  background: #0078d4;
  border-color: #0078d4;
  color: #ffffff;
}

.btn-icon {
  font-size: 13px;
}

.btn-label {
  font-size: 12px;
}

.toolbar-divider {
  width: 1px;
  height: 20px;
  background: #3c3c3c;
  margin: 0 4px;
}

.toolbar-spacer {
  flex: 1;
}

/* 简历标题编辑 per D-30 */
.title-edit-wrapper {
  display: flex;
  align-items: center;
}

.title-display {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 6px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.15s, border-color 0.15s;
  max-width: 200px;
}

.title-display:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: #3c3c3c;
}

.title-text {
  font-size: 12px;
  color: #e6edf3;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 180px;
}

.title-edit-icon {
  color: #8b949e;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.15s;
}

.title-display:hover .title-edit-icon {
  opacity: 1;
}

.title-input {
  padding: 2px 6px;
  font-size: 12px;
  font-weight: 500;
  color: #e6edf3;
  background: #1e1e1e;
  border: 1px solid #58a6ff;
  border-radius: 4px;
  outline: none;
  width: 180px;
}

/* 双栏模式 */
.editor-pane,
.preview-pane {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.pane-header {
  height: 36px;
  min-height: 36px;
  padding: 0 12px;
  background: #f3f3f3;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
}

.pane-title {
  font-size: 12px;
  color: #555;
  font-weight: 500;
}

.editor-wrapper {
  flex: 1;
  overflow: hidden;
}

/* 单栏模式 per D-09 */
.single-pane-mode {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.view-tabs {
  display: flex;
  height: 36px;
  min-height: 36px;
  background: #f3f3f3;
  border-bottom: 1px solid #e0e0e0;
}

.view-tabs button {
  flex: 1;
  background: transparent;
  border: none;
  color: #555;
  font-size: 13px;
  cursor: pointer;
  transition: color 0.15s, background-color 0.15s;
}

.view-tabs button.active {
  color: #1a1a1a;
  background: #ffffff;
  border-bottom: 2px solid #0078d4;
}

.view-tabs button:hover:not(.active) {
  color: #1a1a1a;
  background: #e8e8e8;
}

.single-pane-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.single-pane-content .editor-wrapper,
.single-pane-content .preview-wrapper {
  height: 100%;
  overflow: hidden;
}
</style>
