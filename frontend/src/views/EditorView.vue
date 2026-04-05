<template>
  <div class="editor-view">
    <!-- Editor Toolbar -->
    <div class="editor-toolbar">
      <button class="toolbar-btn" @click="showParserModal = true" title="导入旧简历">
        <span class="btn-icon">📥</span>
        <span class="btn-label">导入</span>
      </button>
      <div class="toolbar-divider" />
      <JobTargetChip
        v-model="jobTarget"
        @change="handleJobTargetChange"
      />
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import MonacoEditor from '../components/MonacoEditor.vue'
import SplitPane from '../components/SplitPane.vue'
import A4Page from '../components/A4Page.vue'
import JobTargetChip from '../components/JobTargetChip.vue'
import ResumeParserModal from '../components/ResumeParserModal.vue'
import { GetResume, UpdateResume } from '../wailsjs/wailsjs/go/main/App'
import type { Resume } from '../types/resume'

const route = useRoute()
const monacoRef = ref<InstanceType<typeof MonacoEditor> | null>(null)

// 当前简历 ID
const resumeId = computed(() => route.params.id as string)

// 编辑器内容
const content = ref('')
const debouncedContent = ref('')
const jobTarget = ref('')

// Modal 状态
const showParserModal = ref(false)

// 响应式状态 per D-09
const windowWidth = ref(window.innerWidth)
const isSinglePane = computed(() => windowWidth.value < 1200)
const activeView = ref<'editor' | 'preview'>('editor')

// 防抖保存 timer
let saveTimer: ReturnType<typeof setTimeout> | null = null

// 加载简历数据
onMounted(async () => {
  const id = resumeId.value
  if (id) {
    try {
      const resume = await GetResume(id)
      content.value = resume.markdownContent || ''
      debouncedContent.value = resume.markdownContent || ''
      jobTarget.value = (resume as any).jobTarget || ''
    } catch (err) {
      console.error('Failed to load resume:', err)
    }
  }

  // 监听窗口宽度变化
  window.addEventListener('resize', handleResize)

  // 编辑器与预览滚动同步 per D-11
  setTimeout(setupScrollSync, 500)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (saveTimer) clearTimeout(saveTimer)
})

function handleResize() {
  windowWidth.value = window.innerWidth
}

// 内容变化处理 per D-10 (150ms debounce)
function handleContentChange(newContent: string) {
  content.value = newContent
  debouncedSave()
}

// 150ms 防抖 timer per D-10
function debouncedSave() {
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = setTimeout(() => {
    debouncedContent.value = content.value
    saveResume()
  }, 150)
}

// 保存简历到后端
async function saveResume() {
  const id = resumeId.value
  if (!id) return

  try {
    // 构建更新数据（包含 markdownContent 和 jobTarget）
    const resumeData: Record<string, unknown> = {
      markdownContent: content.value,
      jobTarget: jobTarget.value,
    }
    await UpdateResume(id, JSON.stringify(resumeData))
  } catch (err) {
    console.error('Failed to save resume:', err)
  }
}

// 职位目标变更时保存
async function handleJobTargetChange(value: string) {
  jobTarget.value = value
  debouncedSave()
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

  // 保存到后端
  await saveResume()

  // 关闭 modal
  showParserModal.value = false

  // 聚焦编辑器
  monacoRef.value?.focus()
}

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
