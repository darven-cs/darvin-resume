<template>
  <div class="shortcut-settings">
    <!-- 快捷键列表 -->
    <div class="shortcut-list">
      <div
        v-for="group in groupedShortcuts"
        :key="group.label"
        class="shortcut-group"
      >
        <h4 class="group-label">{{ group.label }}</h4>
        <div
          v-for="item in group.items"
          :key="item.id"
          class="shortcut-item"
        >
          <span class="shortcut-label">{{ item.label }}</span>
          <div class="shortcut-right">
            <kbd
              :class="['shortcut-key', {
                'listening': listeningId === item.id,
                'has-conflict': getConflict(item.id),
                'is-custom': isCustom(item.id),
              }]"
              @click="startListening(item.id)"
              :title="listeningId === item.id ? '按下新快捷键...' : '点击修改快捷键'"
            >
              <template v-if="listeningId === item.id">
                <span class="listening-hint">按下新快捷键...</span>
              </template>
              <template v-else>
                {{ formatKeyDisplay(item.effectiveKey) }}
              </template>
            </kbd>
            <span v-if="getConflict(item.id)" class="conflict-badge" :title="getConflict(item.id) ?? undefined">
              !
            </span>
            <button
              v-if="isCustom(item.id)"
              class="reset-btn"
              @click="handleReset(item.id)"
              title="恢复默认"
            >
              重置
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 全部重置按钮 -->
    <div class="shortcut-footer">
      <button class="reset-all-btn" @click="handleResetAll">
        全部恢复默认
      </button>
    </div>

    <!-- 监听模式的全局遮罩（ESC 取消） -->
    <div
      v-if="listeningId"
      class="listening-overlay"
      @click="cancelListening"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useKeyboard, formatKeyDisplay as fmtDisplay, DEFAULT_SHORTCUTS, normalizeKey } from '../composables/useKeyboard'

const keyboard = useKeyboard()

const listeningId = ref<string | null>(null)

interface ShortcutItem {
  id: string
  label: string
  effectiveKey: string
  scope: string
}

interface ShortcutGroup {
  label: string
  items: ShortcutItem[]
}

// 直接从 DEFAULT_SHORTCUTS 构建列表，不依赖 register() 后的 bindings Map
const groupedShortcuts = computed<ShortcutGroup[]>(() => {
  const currentOverrides = keyboard.overrides.value
  const groups: ShortcutGroup[] = []

  const groupMap: Record<string, { label: string; items: ShortcutItem[] }> = {
    ai: { label: 'AI 操作', items: [] },
    view: { label: '视图', items: [] },
    file: { label: '文件', items: [] },
  }

  for (const def of DEFAULT_SHORTCUTS) {
    const prefix = def.id.split('.')[0]
    const group = groupMap[prefix]
    if (group) {
      const customKey = currentOverrides[def.id]
      group.items.push({
        id: def.id,
        label: def.label,
        effectiveKey: customKey || def.defaultKey,
        scope: def.scope,
      })
    }
  }

  for (const key of Object.keys(groupMap)) {
    if (groupMap[key].items.length > 0) {
      groups.push(groupMap[key])
    }
  }

  return groups
})

function formatKeyDisplay(keyStr: string): string {
  return fmtDisplay(keyStr)
}

function isCustom(id: string): boolean {
  return !!keyboard.overrides.value[id]
}

function getConflict(id: string): string | null {
  const currentOverrides = keyboard.overrides.value
  const myKey = currentOverrides[id] || DEFAULT_SHORTCUTS.find(s => s.id === id)?.defaultKey
  if (!myKey) return null

  const normalizedMyKey = normalizeKey(myKey)
  for (const def of DEFAULT_SHORTCUTS) {
    if (def.id === id) continue
    const theirKey = currentOverrides[def.id] || def.defaultKey
    if (normalizeKey(theirKey) === normalizedMyKey) {
      return `与"${def.label}"冲突`
    }
  }
  return null
}

function startListening(id: string) {
  listeningId.value = id
}

function cancelListening() {
  listeningId.value = null
}

async function handleReset(id: string) {
  await keyboard.resetToDefault(id)
}

async function handleResetAll() {
  await keyboard.resetAll()
}

// 监听模式：捕获按键
function handleKeyCapture(e: KeyboardEvent) {
  if (!listeningId.value) return

  // ESC 取消监听
  if (e.key === 'Escape') {
    e.preventDefault()
    e.stopPropagation()
    listeningId.value = null
    return
  }

  // 忽略单独的修饰键
  if (['Control', 'Shift', 'Alt', 'Meta'].includes(e.key)) {
    return
  }

  e.preventDefault()
  e.stopPropagation()

  // 构建快捷键字符串
  const parts: string[] = []
  if (e.ctrlKey || e.metaKey) parts.push('Ctrl')
  if (e.shiftKey) parts.push('Shift')
  if (e.altKey) parts.push('Alt')
  parts.push(e.key.length === 1 ? e.key.toUpperCase() : e.key)

  const newKey = parts.join('+')
  const id = listeningId.value
  listeningId.value = null

  keyboard.setCustomKey(id, newKey)
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyCapture, true)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyCapture, true)
})
</script>

<style scoped>
.shortcut-settings {
  display: flex;
  flex-direction: column;
  gap: var(--ui-spacing-md);
}

.shortcut-list {
  display: flex;
  flex-direction: column;
  gap: var(--ui-spacing-lg);
}

.shortcut-group {
  display: flex;
  flex-direction: column;
  gap: var(--ui-spacing-xs);
}

.group-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--ui-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0;
  padding-bottom: var(--ui-spacing-xs);
  border-bottom: 1px solid var(--ui-border);
}

.shortcut-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
  min-height: 32px;
}

.shortcut-label {
  font-size: 13px;
  color: var(--ui-text-primary);
}

.shortcut-right {
  display: flex;
  align-items: center;
  gap: var(--ui-spacing-sm);
}

.shortcut-key {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 90px;
  padding: 4px 10px;
  background: var(--ui-bg-tertiary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  font-family: inherit;
  font-size: 12px;
  color: var(--ui-text-primary);
  cursor: pointer;
  transition: border-color var(--ui-transition-fast), background var(--ui-transition-fast), box-shadow var(--ui-transition-fast);
  user-select: none;
}

.shortcut-key:hover {
  border-color: var(--ui-border-hover);
  background: var(--ui-bg-hover);
}

.shortcut-key.listening {
  border-color: var(--ui-accent);
  box-shadow: 0 0 0 2px var(--ui-accent-focus);
  background: var(--ui-bg-active);
  animation: pulse-border 1.5s ease-in-out infinite;
}

.shortcut-key.has-conflict {
  border-color: var(--ui-danger);
  color: var(--ui-danger);
}

.shortcut-key.is-custom {
  border-color: var(--ui-accent);
  color: var(--ui-accent);
}

.listening-hint {
  font-size: 11px;
  color: var(--ui-accent);
}

.conflict-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--ui-danger);
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  cursor: help;
}

.reset-btn {
  padding: 2px 8px;
  background: transparent;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-secondary);
  font-size: 11px;
  cursor: pointer;
  transition: color var(--ui-transition-fast), border-color var(--ui-transition-fast);
}

.reset-btn:hover {
  color: var(--ui-text-primary);
  border-color: var(--ui-border-hover);
}

.shortcut-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: var(--ui-spacing-sm);
  border-top: 1px solid var(--ui-border);
}

.reset-all-btn {
  padding: 4px 12px;
  background: transparent;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: color var(--ui-transition-fast), border-color var(--ui-transition-fast);
}

.reset-all-btn:hover {
  color: var(--ui-text-primary);
  border-color: var(--ui-border-hover);
}

.listening-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: transparent;
  cursor: default;
}

@keyframes pulse-border {
  0%, 100% {
    box-shadow: 0 0 0 2px var(--ui-accent-focus);
  }
  50% {
    box-shadow: 0 0 0 4px var(--ui-accent-focus);
  }
}
</style>
