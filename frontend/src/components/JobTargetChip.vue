<template>
  <div class="job-target-chip" :class="{ 'job-target-chip--editing': isEditing }">
    <span v-if="!isEditing" class="chip-display" @click="startEditing" title="点击修改目标岗位">
      <span class="chip-icon">🎯</span>
      <span class="chip-text">{{ displayText }}</span>
      <span v-if="modelValue" class="chip-edit-hint">✎</span>
    </span>

    <div v-else class="chip-edit">
      <span class="chip-icon">🎯</span>
      <input
        ref="inputRef"
        v-model="editValue"
        class="chip-input"
        type="text"
        placeholder="输入目标岗位..."
        @blur="confirmEdit"
        @keydown.enter="confirmEdit"
        @keydown.escape="cancelEdit"
      />
      <button class="chip-confirm" @mousedown.prevent="confirmEdit" title="确认">✓</button>
      <button class="chip-cancel" @mousedown.prevent="cancelEdit" title="取消">✕</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'

const props = defineProps<{
  modelValue: string
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  change: [value: string]
}>()

const isEditing = ref(false)
const editValue = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

const displayText = computed(() => {
  if (props.modelValue) {
    return `目标岗位: ${props.modelValue}`
  }
  return '添加目标岗位'
})

async function startEditing() {
  if (props.disabled) return
  editValue.value = props.modelValue
  isEditing.value = true
  await nextTick()
  inputRef.value?.select()
}

function confirmEdit() {
  const value = editValue.value.trim()
  emit('update:modelValue', value)
  emit('change', value)
  isEditing.value = false
}

function cancelEdit() {
  isEditing.value = false
  editValue.value = ''
}

// Reset editing state if modelValue changes externally
watch(() => props.modelValue, () => {
  if (!isEditing.value) {
    editValue.value = props.modelValue
  }
})
</script>

<style scoped>
.job-target-chip {
  display: inline-flex;
  align-items: center;
}

.chip-display {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: rgba(14, 99, 156, 0.25);
  border: 1px solid rgba(0, 120, 212, 0.4);
  border-radius: 16px;
  cursor: pointer;
  font-size: 12px;
  color: #79c0ff;
  transition: background-color 0.15s, border-color 0.15s;
  user-select: none;
}

.chip-display:hover {
  background: rgba(14, 99, 156, 0.4);
  border-color: rgba(0, 120, 212, 0.6);
}

.chip-icon {
  font-size: 11px;
}

.chip-text {
  color: #c9d1d9;
}

.chip-edit-hint {
  font-size: 10px;
  color: #5a5a5a;
  margin-left: 2px;
}

.chip-display:hover .chip-edit-hint {
  color: #8b949e;
}

/* Edit Mode */
.chip-edit {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 4px 2px 6px;
  background: #1e1e1e;
  border: 1px solid #0078d4;
  border-radius: 16px;
}

.chip-edit .chip-icon {
  font-size: 11px;
  flex-shrink: 0;
}

.chip-input {
  background: transparent;
  border: none;
  color: #e0e0e0;
  font-size: 12px;
  width: 160px;
  padding: 2px 0;
  outline: none;
}

.chip-input::placeholder {
  color: #5a5a5a;
}

.chip-confirm,
.chip-cancel {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 2px 4px;
  font-size: 11px;
  border-radius: 50%;
  transition: background-color 0.1s;
}

.chip-confirm {
  color: #3fb950;
}

.chip-confirm:hover {
  background: rgba(63, 185, 80, 0.2);
}

.chip-cancel {
  color: #8b949e;
}

.chip-cancel:hover {
  background: rgba(139, 148, 158, 0.2);
}
</style>
