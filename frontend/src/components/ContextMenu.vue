<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="context-menu"
      :style="menuStyle"
      @click.stop
    >
      <button
        v-for="item in items"
        :key="item.id"
        class="menu-item"
        :class="{ disabled: item.disabled, danger: item.danger }"
        :disabled="item.disabled"
        @click="handleClick(item)"
      >
        <span class="menu-icon">{{ item.icon }}</span>
        <span class="menu-label">{{ item.label }}</span>
        <span v-if="item.shortcut" class="menu-shortcut">{{ item.shortcut }}</span>
      </button>
    </div>
    <!-- 点击外部关闭 -->
    <div
      v-if="visible"
      class="context-menu-overlay"
      @click="close"
      @contextmenu.prevent="close"
    />
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'

export interface MenuItem {
  id: string
  label: string
  icon?: string
  shortcut?: string
  disabled?: boolean
  danger?: boolean
  action: () => void
}

const props = defineProps<{
  items: MenuItem[]
  x: number
  y: number
  visible: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const menuStyle = computed(() => {
  // 确保菜单位于视口内
  const menuWidth = 180
  const menuHeight = props.items.length * 32 + 8

  let left = props.x
  let top = props.y

  // 右侧超出
  if (left + menuWidth > window.innerWidth) {
    left = window.innerWidth - menuWidth - 8
  }
  // 底部超出
  if (top + menuHeight > window.innerHeight) {
    top = window.innerHeight - menuHeight - 8
  }

  return {
    left: `${Math.max(4, left)}px`,
    top: `${Math.max(4, top)}px`,
  }
})

function handleClick(item: MenuItem) {
  if (item.disabled) return
  item.action()
  emit('close')
}

function close() {
  emit('close')
}

// ESC 键关闭
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.context-menu {
  position: fixed;
  z-index: 10000;
  min-width: 160px;
  background: #252526;
  border: 1px solid #3c3c3c;
  border-radius: 6px;
  padding: 4px 0;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
}

.menu-item {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 6px 12px;
  background: transparent;
  border: none;
  color: #cccccc;
  font-size: 13px;
  text-align: left;
  cursor: pointer;
  gap: 8px;
  transition: background-color 0.1s;
}

.menu-item:hover:not(.disabled) {
  background: #094771;
  color: #ffffff;
}

.menu-item.danger:hover:not(.disabled) {
  background: #5a1d1d;
  color: #ff6b6b;
}

.menu-item.disabled {
  color: #5a5a5a;
  cursor: not-allowed;
}

.menu-icon {
  width: 16px;
  font-size: 12px;
}

.menu-label {
  flex: 1;
}

.menu-shortcut {
  font-size: 11px;
  color: #8b949e;
}

.context-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
}
</style>
