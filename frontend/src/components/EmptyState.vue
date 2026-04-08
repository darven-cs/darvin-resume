<template>
  <div class="empty-state">
    <div class="empty-state-icon" v-if="$slots.icon">
      <slot name="icon" />
    </div>
    <div class="empty-state-icon default-icon" v-else>
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
    </div>
    <h3 v-if="title" class="empty-state-title">{{ title }}</h3>
    <p v-if="description" class="empty-state-desc">{{ description }}</p>
    <slot name="action">
      <button
        v-if="actionLabel"
        class="empty-state-action"
        @click="$emit('action')"
      >
        {{ actionLabel }}
      </button>
    </slot>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  /** 标题文字 */
  title?: string
  /** 描述文字 */
  description?: string
  /** 操作按钮文字（为空则不显示按钮） */
  actionLabel?: string
}>()

defineEmits<{
  (e: 'action'): void
}>()
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 240px;
  text-align: center;
  padding: var(--ui-spacing-xl);
}

.empty-state-icon {
  margin-bottom: var(--ui-spacing-lg);
  color: var(--ui-text-tertiary);
  opacity: 0.5;
}

.default-icon {
  opacity: 0.3;
}

.empty-state-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--ui-text-primary);
  margin: 0 0 var(--ui-spacing-sm) 0;
}

.empty-state-desc {
  font-size: 0.875rem;
  color: var(--ui-text-tertiary);
  margin: 0 0 var(--ui-spacing-xl) 0;
  max-width: 320px;
  line-height: 1.5;
}

.empty-state-action {
  padding: 8px 24px;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--ui-text-inverse);
  background: var(--ui-accent);
  border: none;
  border-radius: var(--ui-radius-md);
  cursor: pointer;
  transition: background-color var(--ui-transition-fast);
}

.empty-state-action:hover {
  background: var(--ui-accent-hover);
}
</style>
