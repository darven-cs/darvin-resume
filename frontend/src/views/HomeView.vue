<template>
  <div class="home-view">
    <!-- 顶部栏 -->
    <header class="top-bar">
      <div class="top-bar-left">
        <h1 class="app-title">Open-Resume</h1>
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
      </div>
    </header>

    <!-- 卡片网格区域 -->
    <main class="resume-grid-area">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <span>加载中...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="filteredResumes.length === 0 && !searchQuery" class="empty-state">
        <svg class="empty-icon" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
        </svg>
        <p class="empty-title">还没有简历</p>
        <p class="empty-desc">点击上方「新建简历」按钮开始创建</p>
        <button class="empty-create-btn" @click="showCreateModal = true">新建简历</button>
      </div>

      <!-- 搜索无结果 -->
      <div v-else-if="filteredResumes.length === 0 && searchQuery" class="empty-state">
        <p class="empty-title">未找到匹配的简历</p>
        <p class="empty-desc">尝试其他关键词搜索</p>
      </div>

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
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { CreateResume, UpdateResume } from '../wailsjs/wailsjs/go/main/App'
import { useResumeList } from '../composables/useResumeList'
import ResumeCard from '../components/ResumeCard.vue'
import CreateModeModal from '../components/CreateModeModal.vue'
import RecycleBinSection from '../components/RecycleBinSection.vue'

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
  background: #181818;
  overflow: hidden;
}

/* 顶部栏 */
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #2d2d2d;
  flex-shrink: 0;
}

.top-bar-left {
  display: flex;
  align-items: center;
}

.app-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #e6edf3;
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
  color: #8b949e;
  pointer-events: none;
}

.search-input {
  padding: 6px 12px 6px 32px;
  font-size: 0.875rem;
  color: #e6edf3;
  background: #2d2d2d;
  border: 1px solid #3c3c3c;
  border-radius: 6px;
  outline: none;
  width: 200px;
  transition: border-color 0.2s, width 0.2s;
}

.search-input::placeholder {
  color: #8b949e;
}

.search-input:focus {
  border-color: #58a6ff;
  width: 260px;
}

/* 排序按钮 */
.sort-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid #3c3c3c;
  border-radius: 6px;
  background: #2d2d2d;
  color: #8b949e;
  cursor: pointer;
  transition: border-color 0.2s, color 0.2s;
}

.sort-btn:hover {
  border-color: #58a6ff;
  color: #e6edf3;
}

/* 新建按钮 */
.create-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  background: #238636;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.create-btn:hover {
  background: #2ea043;
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
  color: #8b949e;
  font-size: 0.875rem;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 300px;
  text-align: center;
}

.empty-icon {
  color: #3c3c3c;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #e6edf3;
  margin: 0 0 8px 0;
}

.empty-desc {
  font-size: 0.875rem;
  color: #8b949e;
  margin: 0 0 24px 0;
}

.empty-create-btn {
  padding: 8px 24px;
  font-size: 0.95rem;
  font-weight: 500;
  color: #fff;
  background: #238636;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.empty-create-btn:hover {
  background: #2ea043;
}
</style>
