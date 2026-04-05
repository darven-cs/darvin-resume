<template>
  <div class="save-status-indicator">
    <!-- 已保存状态: 绿色圆点 + "已保存" -->
    <Transition name="fade" mode="out-in">
      <span v-if="status === 'saved'" key="saved" class="status-item status-saved">
        <span class="status-dot green" />
        <span class="status-text">已保存</span>
      </span>

      <!-- 保存中状态: 灰色旋转图标 + "保存中..." -->
      <span v-else-if="status === 'saving'" key="saving" class="status-item status-saving">
        <span class="status-spinner" />
        <span class="status-text">保存中...</span>
      </span>

      <!-- 未保存状态: 橙色圆点 + "未保存" -->
      <span v-else-if="status === 'unsaved'" key="unsaved" class="status-item status-unsaved">
        <span class="status-dot orange" />
        <span class="status-text">未保存</span>
      </span>

      <!-- 保存失败状态: 红色圆点 + "保存失败" + 重试按钮 -->
      <span v-else-if="status === 'error'" key="error" class="status-item status-error">
        <span class="status-dot red" />
        <span class="status-text">保存失败</span>
        <button class="retry-btn" @click.stop="$emit('retry')" title="重试保存">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 4v6h6"/>
            <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"/>
          </svg>
        </button>
      </span>
    </Transition>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  status: 'saved' | 'saving' | 'unsaved' | 'error'
  errorMessage?: string
}>()

defineEmits<{
  (e: 'retry'): void
}>()
</script>

<style scoped>
.save-status-indicator {
  display: flex;
  align-items: center;
  font-size: 12px;
}

.status-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-dot.green {
  background: #2ea043;
}

.status-dot.orange {
  background: #d29922;
}

.status-dot.red {
  background: #f85149;
}

.status-spinner {
  width: 12px;
  height: 12px;
  border: 2px solid #555;
  border-top-color: #8b949e;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.status-text {
  color: #8b949e;
  font-size: 12px;
}

.retry-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  padding: 0;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: #f85149;
  cursor: pointer;
  transition: background-color 0.15s;
}

.retry-btn:hover {
  background: rgba(248, 81, 73, 0.15);
}

/* Transition animations */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
