<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click="handleClose">
      <div class="modal-container" @click.stop>
        <!-- 关闭按钮 -->
        <button class="modal-close" @click="handleClose">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>

        <h2 class="modal-title">新建简历</h2>
        <p class="modal-subtitle">选择创建方式</p>

        <div class="mode-options">
          <!-- AI 引导创建 -->
          <button class="mode-card" @click="selectMode('wizard')">
            <div class="mode-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M12 2L2 7l10 5 10-5-10-5z"/>
                <path d="M2 17l10 5 10-5"/>
                <path d="M2 12l10 5 10-5"/>
              </svg>
            </div>
            <div class="mode-info">
              <span class="mode-name">AI 引导创建</span>
              <span class="mode-desc">分步填写，AI 润色优化</span>
            </div>
          </button>

          <!-- 空白页创建 -->
          <button class="mode-card" @click="selectMode('blank')">
            <div class="mode-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10 9 9 9 8 9"/>
              </svg>
            </div>
            <div class="mode-info">
              <span class="mode-name">空白页创建</span>
              <span class="mode-desc">自由编辑，随时调用 AI</span>
            </div>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'select-mode', mode: 'wizard' | 'blank'): void
}>()

function handleClose() {
  emit('close')
}

function selectMode(mode: 'wizard' | 'blank') {
  emit('select-mode', mode)
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--ui-overlay-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  position: relative;
  width: 480px;
  max-width: 90vw;
  padding: 32px;
  background: var(--ui-bg-secondary);
  border: 1px solid var(--ui-border);
  border-radius: 12px;
  box-shadow: var(--ui-shadow-lg);
}

.modal-close {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--ui-radius-md);
  background: transparent;
  color: var(--ui-text-tertiary);
  cursor: pointer;
  transition: background var(--ui-transition-fast), color var(--ui-transition-fast);
}

.modal-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--ui-text-primary);
}

.modal-title {
  margin: 0 0 4px 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--ui-text-primary);
}

.modal-subtitle {
  margin: 0 0 24px 0;
  font-size: 0.875rem;
  color: var(--ui-text-tertiary);
}

.mode-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.mode-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  background: var(--ui-bg-primary);
  cursor: pointer;
  transition: border-color var(--ui-transition-fast), background var(--ui-transition-fast);
  text-align: left;
  width: 100%;
}

.mode-card:hover {
  border-color: var(--ui-accent);
  background: var(--ui-bg-hover);
}

.mode-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: var(--ui-radius-lg);
  background: rgba(88, 166, 255, 0.1);
  color: var(--ui-accent);
  flex-shrink: 0;
}

.mode-info {
  display: flex;
  flex-direction: column;
  gap: var(--ui-radius-sm);
}

.mode-name {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--ui-text-primary);
}

.mode-desc {
  font-size: 0.8rem;
  color: var(--ui-text-tertiary);
}
</style>
