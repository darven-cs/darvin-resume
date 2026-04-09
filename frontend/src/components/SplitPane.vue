<template>
  <div
    class="split-pane"
    :class="{ 'is-dragging': isDragging }"
    :style="containerStyle"
  >
    <!-- 左侧面板 -->
    <div
      class="pane left-pane"
      :style="leftPaneStyle"
    >
      <slot name="left" />
    </div>

    <!-- 分割线 -->
    <div
      class="divider"
      :class="{ dragging: isDragging }"
      @mousedown="startDrag"
      @mouseenter="hoverDivider"
      @mouseleave="leaveDivider"
    >
      <div class="divider-handle" />
    </div>

    <!-- 右侧面板 -->
    <div
      class="pane right-pane"
      :style="rightPaneStyle"
    >
      <slot name="right" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'

const props = withDefaults(defineProps<{
  /** 默认分割比例，左侧占百分比 */
  defaultRatio?: number
  /** 每栏最小宽度（px） */
  minWidth?: number
}>(), {
  defaultRatio: 50,
  minWidth: 300,
})

const splitRatio = ref(props.defaultRatio)
const isDragging = ref(false)

const containerStyle = computed((): Record<string, string> => ({
  display: 'flex',
  height: '100%',
  width: '100%',
  userSelect: isDragging.value ? 'none' : 'auto',
}))

const leftPaneStyle = computed((): Record<string, string> => ({
  width: `${splitRatio.value}%`,
  minWidth: `${props.minWidth}px`,
  overflow: 'hidden',
}))

const rightPaneStyle = computed((): Record<string, string> => ({
  width: `${100 - splitRatio.value}%`,
  minWidth: `${props.minWidth}px`,
  overflow: 'hidden',
}))

function startDrag(e: MouseEvent) {
  e.preventDefault()
  isDragging.value = true
  const startX = e.clientX
  const startRatio = splitRatio.value

  const container = document.querySelector('.split-pane') as HTMLElement
  if (!container) return
  const containerWidth = container.offsetWidth

  const onMove = (e: MouseEvent) => {
    const delta = e.clientX - startX
    const deltaRatio = (delta / containerWidth) * 100
    const newRatio = startRatio + deltaRatio
    // 限制范围在 20% - 80%
    splitRatio.value = Math.max(20, Math.min(80, newRatio))
  }

  const onUp = () => {
    isDragging.value = false
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  }

  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}

function hoverDivider() {
  if (!isDragging.value) {
    document.body.style.cursor = 'col-resize'
  }
}

function leaveDivider() {
  if (!isDragging.value) {
    document.body.style.cursor = ''
  }
}

onUnmounted(() => {
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
})
</script>

<style scoped>
.split-pane {
  position: relative;
  width: 100%;
  height: 100%;
}

.pane {
  height: 100%;
  overflow: hidden;
}

.divider {
  position: relative;
  width: 6px;
  height: 100%;
  cursor: col-resize;
  background: transparent;
  flex-shrink: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color var(--ui-transition-fast);
}

.divider:hover,
.divider.dragging {
  background: rgba(0, 120, 212, 0.15);
}

.divider-handle {
  width: 2px;
  height: 40px;
  background: var(--ui-text-secondary);
  border-radius: 1px;
  transition: background-color var(--ui-transition-fast), height var(--ui-transition-fast);
}

.divider:hover .divider-handle,
.divider.dragging .divider-handle {
  background: var(--ui-accent);
  height: 48px;
}

.is-dragging {
  cursor: col-resize;
}
</style>
