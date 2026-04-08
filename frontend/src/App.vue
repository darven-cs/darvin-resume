<template>
  <div class="app-layout">
    <aside class="sidebar" :class="{ collapsed: isSidebarCollapsed, 'hover-expanded': isHoverExpanded }" @mouseenter="handleSidebarMouseEnter" @mouseleave="handleSidebarMouseLeave">
      <div class="sidebar-header">
        <h2 v-if="!isSidebarCollapsed || isHoverExpanded" class="sidebar-title">Darvin-Resume</h2>
        <span v-else class="sidebar-logo-icon">D</span>
      </div>
      <nav class="sidebar-nav">
        <RouterLink to="/" class="nav-item">
          <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
          </svg>
          <span v-if="!isSidebarCollapsed || isHoverExpanded" class="nav-label">简历列表</span>
        </RouterLink>
      </nav>
    </aside>
    <main class="main-content">
      <router-view />
    </main>

    <!-- 全局 Toast 通知容器 -->
    <ToastContainer />

    <!-- 全局确认对话框 -->
    <ConfirmDialog />
  </div>
</template>

<script setup lang="ts">
import ToastContainer from './components/ToastContainer.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'
import { ref, computed, onMounted, onUnmounted } from 'vue'

const SIDEBAR_COLLAPSE_BREAKPOINT = 1200

const windowWidth = ref(window.innerWidth)
const isHoverExpanded = ref(false)

const isSidebarCollapsed = computed(() => windowWidth.value < SIDEBAR_COLLAPSE_BREAKPOINT)

function handleResize() {
  windowWidth.value = window.innerWidth
  // 窗口恢复宽度时关闭 hover 状态
  if (!isSidebarCollapsed.value) {
    isHoverExpanded.value = false
  }
}

function handleSidebarMouseEnter() {
  if (isSidebarCollapsed.value) {
    isHoverExpanded.value = true
  }
}

function handleSidebarMouseLeave() {
  isHoverExpanded.value = false
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  width: 100%;
  overflow: hidden;
}

body {
  font-family: var(--ui-font-sans);
  background-color: var(--ui-bg-primary);
  color: var(--ui-text-primary);
}

.app-layout {
  display: flex;
  height: 100%;
  width: 100%;
}

.sidebar {
  width: 240px;
  min-width: 240px;
  background-color: var(--ui-bg-sidebar);
  border-right: 1px solid var(--ui-border);
  display: flex;
  flex-direction: column;
  transition: width var(--ui-transition-fast) ease, min-width var(--ui-transition-fast) ease;
  overflow: hidden;
}

/* 收起为图标栏模式 */
.sidebar.collapsed {
  width: 48px;
  min-width: 48px;
}

/* 收起状态 hover 展开（overlay 模式，不挤压主内容） */
.sidebar.collapsed.hover-expanded {
  width: 240px;
  min-width: 240px;
  position: absolute;
  z-index: 100;
  box-shadow: var(--ui-shadow-md);
}

.sidebar-header {
  padding: 1rem;
  border-bottom: 1px solid var(--ui-border);
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 52px;
}

.sidebar-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--ui-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar-logo-icon {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--ui-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: var(--ui-radius-md);
  background: var(--ui-bg-hover);
}

.sidebar-nav {
  padding: 0.5rem 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0.625rem 1rem;
  color: var(--ui-text-primary);
  text-decoration: none;
  transition: background-color var(--ui-transition-fast);
  white-space: nowrap;
  overflow: hidden;
}

.nav-item:hover,
.nav-item.router-link-active {
  background-color: var(--ui-bg-hover);
}

.nav-icon {
  flex-shrink: 0;
  color: var(--ui-text-secondary);
}

.nav-label {
  font-size: 0.875rem;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
}

.main-content {
  flex: 1;
  overflow: auto;
  background-color: var(--ui-bg-primary);
}
</style>
