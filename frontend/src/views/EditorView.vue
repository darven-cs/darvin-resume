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

      <!-- PDF 导出按钮 -->
      <button class="toolbar-btn" @click="showExportDialog = true" title="导出 PDF">
        <span class="btn-icon">PDF</span>
        <span class="btn-label">导出</span>
      </button>

      <!-- 版本快照按钮 -->
      <button
        class="toolbar-btn"
        :class="{ active: snapshotSidebarVisible }"
        @click="snapshotSidebarVisible = !snapshotSidebarVisible"
        title="版本快照"
      >
        <span class="btn-icon">&#128247;</span>
        <span class="btn-label">快照</span>
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

      <!-- 视图模式切换按钮组 -->
      <div class="view-mode-group">
        <button
          class="view-mode-btn"
          :class="{ active: viewMode === 'split' }"
          @click="viewMode = 'split'"
          title="双栏模式"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2"/>
            <line x1="12" y1="3" x2="12" y2="21"/>
          </svg>
        </button>
        <button
          class="view-mode-btn"
          :class="{ active: viewMode === 'editor' }"
          @click="viewMode = 'editor'"
          title="仅编辑"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="12" height="18" rx="2"/>
            <line x1="7" y1="8" x2="11" y2="8"/>
            <line x1="7" y1="12" x2="11" y2="12"/>
          </svg>
        </button>
        <button
          class="view-mode-btn"
          :class="{ active: viewMode === 'preview' }"
          @click="viewMode = 'preview'"
          title="仅预览"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="9" y="3" width="12" height="18" rx="2"/>
            <circle cx="15" cy="12" r="2"/>
          </svg>
        </button>
      </div>

      <button
        class="toolbar-btn"
        :class="{ active: chatSidebarVisible }"
        @click="chatSidebarVisible = !chatSidebarVisible"
        title="AI 对话 (Ctrl+Shift+A)"
      >
        <span class="btn-icon">💬</span>
        <span class="btn-label">AI 助手</span>
      </button>
      <button
        class="toolbar-btn"
        :class="{ active: styleEditorVisible }"
        @click="styleEditorVisible = !styleEditorVisible"
        title="样式调整"
      >
        <span class="btn-icon">🎨</span>
        <span class="btn-label">样式</span>
      </button>
    </div>

    <!-- 双栏模式 (窗口宽度 >= 1200px) per D-09 -->
    <template v-if="!isSinglePane && effectiveViewMode === 'split'">
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
            <div class="preview-content-wrapper">
              <A4Page :content="debouncedContent" :template-id="currentTemplateId" :custom-css="customCss" />
              <StyleEditor
                v-if="styleEditorVisible"
                v-model="styleEditorVisible"
                :resume-id="resumeId"
                :template-id="currentTemplateId"
                :initial-css="customCss"
              />
            </div>
          </div>
        </template>
      </SplitPane>
    </template>

    <!-- 单栏模式 (窗口宽度 < 1200px) per D-09 -->
    <template v-else-if="isSinglePane || effectiveViewMode !== 'split'">
      <div class="single-pane-mode">
        <div class="view-tabs" v-if="effectiveViewMode === 'split'">
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
          <div v-if="effectiveSingleView === 'editor'" class="editor-wrapper">
            <MonacoEditor
              ref="monacoRef"
              v-model="content"
              :job-target="jobTarget"
              @change="handleContentChange"
            />
          </div>
          <div v-else class="preview-wrapper">
            <A4Page :content="debouncedContent" :template-id="currentTemplateId" :custom-css="customCss" />
            <StyleEditor
              v-if="styleEditorVisible"
              v-model="styleEditorVisible"
              :resume-id="resumeId"
              :template-id="currentTemplateId"
              :initial-css="customCss"
            />
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

    <!-- Settings Dialog -->
    <SettingsDialog
      :visible="showSettingsDialog"
      @close="showSettingsDialog = false"
    />

    <!-- Export Dialog -->
    <ExportDialog
      :visible="showExportDialog"
      :resume-id="resumeId"
      @close="showExportDialog = false"
      @exported="handleExportCompleted"
    />

    <!-- Snapshot Sidebar -->
    <SnapshotSidebar
      :visible="snapshotSidebarVisible"
      :resume-id="resumeId"
      @close="snapshotSidebarVisible = false"
      @rollback="handleRollback"
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
import SettingsDialog from '../components/SettingsDialog.vue'
import ExportDialog from '../components/ExportDialog.vue'
import SnapshotSidebar from '../components/SnapshotSidebar.vue'
import ResumeWizardSidebar from '../components/ResumeWizardSidebar.vue'
import SaveStatusIndicator from '../components/SaveStatusIndicator.vue'
import StyleEditor from '../components/StyleEditor.vue'
import { useAutoSave } from '../composables/useAutoSave'
import { useTemplate } from '../composables/useTemplate'
import { useSnapshot } from '../composables/useSnapshot'
import { useKeyboard } from '../composables/useKeyboard'
import { useAISelection } from '../composables/useAISelection'
import { GetResume, RenameResume } from '../wailsjs/wailsjs/go/main/App'
import type { Resume } from '../types/resume'
import type { AIOperationType } from '../types/ai'

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
const showSettingsDialog = ref(false)

// AI Chat Sidebar 状态
const chatSidebarVisible = ref(false)
const chatSidebarRef = ref<InstanceType<typeof AIChatSidebar> | null>(null)

// Resume Wizard 状态 — 检测路由 query ?wizard=true 自动打开向导
const wizardVisible = ref(route.query.wizard === 'true')

// Export Dialog 状态
const showExportDialog = ref(false)

// Snapshot Sidebar 状态
const snapshotSidebarVisible = ref(false)
const {
  loadSnapshots,
  autoCreateSnapshot,
  rollback,
} = useSnapshot()

// 响应式状态 per D-09
const windowWidth = ref(window.innerWidth)
const isSinglePane = computed(() => windowWidth.value < 1200)
const activeView = ref<'editor' | 'preview'>('editor')

// 手动视图模式 (Task 3)
const viewMode = ref<'split' | 'editor' | 'preview'>('split')

/** 实际生效的视图模式：在窄屏时强制为 split（即单栏 Tab 切换） */
const effectiveViewMode = computed(() => {
  if (isSinglePane.value) return 'split'
  return viewMode.value
})

/** 单栏模式下实际显示哪个面板 */
const effectiveSingleView = computed(() => {
  if (viewMode.value === 'editor') return 'editor'
  if (viewMode.value === 'preview') return 'preview'
  return activeView.value
})

// 全局快捷键系统 (07-03)
const keyboard = useKeyboard()

// AI 选区操作 (07-03: 用于快捷键触发的 AI 操作)
const aiSelection = useAISelection(
  () => monacoRef.value?.getEditor?.(),
  jobTarget.value
)

// 注册快捷键 handler
function registerShortcuts() {
  keyboard.register('file.save', () => {
    triggerSave()
  })

  keyboard.register('view.togglePreview', () => {
    if (effectiveViewMode.value === 'split') {
      viewMode.value = 'editor'
    } else {
      viewMode.value = 'split'
    }
  })

  keyboard.register('view.toggleChat', () => {
    chatSidebarVisible.value = !chatSidebarVisible.value
  })

  keyboard.register('ai.polish', () => {
    handleAIShortcut('polish')
  })

  keyboard.register('ai.translate', () => {
    handleAIShortcut('translate')
  })

  keyboard.register('ai.shorten', () => {
    handleAIShortcut('summarize')
  })
}

// AI 快捷键处理：检查选中文本后执行 AI 操作
async function handleAIShortcut(operation: AIOperationType) {
  const editor = monacoRef.value?.getEditor?.()
  if (!editor) return

  const selection = editor.getSelection?.()
  const model = editor.getModel?.()
  if (!selection || !model) return

  const selectedText = model.getValueInRange(selection)
  if (!selectedText || !selectedText.trim()) {
    console.info(`[Shortcut] AI ${operation}: 没有选中文本，跳过`)
    return
  }

  try {
    const result = await aiSelection.performAIOperation(operation, selectedText)
    if (result && selection) {
      // 用 AI 结果替换选区
      editor.executeEdits('ai-shortcut', [{
        range: selection,
        text: result,
      }])
      markDirty()
    }
  } catch (err) {
    console.error(`AI ${operation} 快捷键操作失败:`, err)
  }
}

// 自动保存 composable per D-25~D-28
const { saveStatus, errorMessage, markDirty, triggerSave, startAutoSave, stopAutoSave } = useAutoSave({
  resumeId,
  getData: () => ({ markdownContent: content.value, jobTarget: jobTarget.value })
})

// 模板样式 composable (05-02)
const resumeIdRef = computed(() => resumeId.value)
const { currentTemplateId, customCss, loadTemplate, saveAsTemplate } = useTemplate(resumeIdRef)

// 样式编辑器面板显示状态 (05-02 TMPL-04)
const styleEditorVisible = ref(false)

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

  // 加载模板设置 (05-02 TMPL-04)
  await loadTemplate()

  // 启动自动保存定时器 per D-26
  startAutoSave()

  // 启动全局快捷键系统并注册 handler (07-03)
  registerShortcuts()
  await keyboard.start()

  // 监听窗口宽度变化
  window.addEventListener('resize', handleResize)

  // 加载快照列表
  await loadSnapshots(resumeId.value)

  // 编辑器与预览滚动同步 per D-11
  setTimeout(setupScrollSync, 500)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  // 停止全局快捷键系统 (07-03)
  keyboard.stop()
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

// 导出完成后自动创建快照
async function handleExportCompleted() {
  const dateStr = new Date().toLocaleDateString('zh-CN')
  await autoCreateSnapshot(resumeId.value, `PDF 导出快照 ${dateStr}`, '')
  console.log('PDF 导出成功，已自动创建版本快照')
}

// 回滚后重新加载简历内容
async function handleRollback(snapshotId: string) {
  // 重新加载简历数据
  await loadResume()
  snapshotSidebarVisible.value = false
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
  background: var(--ui-bg-toolbar);
  border-bottom: 1px solid var(--ui-border);
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
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), border-color var(--ui-transition-fast);
}

.toolbar-btn:hover {
  background: var(--ui-bg-hover);
  border-color: var(--ui-border-hover);
}

.toolbar-btn.active {
  background: var(--ui-accent);
  border-color: var(--ui-accent);
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
  background: var(--ui-border);
  margin: 0 4px;
}

.toolbar-spacer {
  flex: 1;
}

/* 视图模式按钮组 */
.view-mode-group {
  display: flex;
  align-items: center;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  overflow: hidden;
}

.view-mode-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 26px;
  background: transparent;
  border: none;
  color: var(--ui-text-tertiary);
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), color var(--ui-transition-fast);
}

.view-mode-btn:not(:last-child) {
  border-right: 1px solid var(--ui-border);
}

.view-mode-btn:hover {
  background: var(--ui-bg-hover);
  color: var(--ui-text-secondary);
}

.view-mode-btn.active {
  background: var(--ui-accent);
  color: #ffffff;
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
  border-radius: var(--ui-radius-sm);
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), border-color var(--ui-transition-fast);
  max-width: 200px;
}

.title-display:hover {
  background: var(--ui-bg-hover);
  border-color: var(--ui-border);
}

.title-text {
  font-size: 12px;
  color: var(--ui-text-primary);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 180px;
}

.title-edit-icon {
  color: var(--ui-text-tertiary);
  flex-shrink: 0;
  opacity: 0;
  transition: opacity var(--ui-transition-fast);
}

.title-display:hover .title-edit-icon {
  opacity: 1;
}

.title-input {
  padding: 2px 6px;
  font-size: 12px;
  font-weight: 500;
  color: var(--ui-text-primary);
  background: var(--ui-input-bg);
  border: 1px solid var(--ui-border-focus);
  border-radius: var(--ui-radius-sm);
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

.preview-content-wrapper {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.pane-header {
  height: 36px;
  min-height: 36px;
  padding: 0 12px;
  background: var(--ui-bg-secondary);
  border-bottom: 1px solid var(--ui-border);
  display: flex;
  align-items: center;
}

.pane-title {
  font-size: 12px;
  color: var(--ui-text-secondary);
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
  background: var(--ui-bg-secondary);
  border-bottom: 1px solid var(--ui-border);
}

.view-tabs button {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--ui-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: color var(--ui-transition-fast), background-color var(--ui-transition-fast);
}

.view-tabs button.active {
  color: var(--ui-text-primary);
  background: var(--ui-bg-primary);
  border-bottom: 2px solid var(--ui-accent);
}

.view-tabs button:hover:not(.active) {
  color: var(--ui-text-primary);
  background: var(--ui-bg-hover);
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
  display: flex;
  flex: 1;
}
</style>
