<template>
  <div class="module-select">
    <div class="select-header">
      <h3 class="select-title">选择简历模块</h3>
      <p class="select-desc">选择你想在简历中包含的模块，至少选择1个</p>
    </div>

    <div class="module-list">
      <label
        v-for="mod in wizardModules"
        :key="mod.type"
        :class="['module-item', { checked: selected.includes(mod.type) }]"
      >
        <input
          type="checkbox"
          :checked="selected.includes(mod.type)"
          class="module-checkbox"
          @change="toggleModule(mod.type)"
        />
        <span class="module-icon">{{ mod.icon }}</span>
        <div class="module-info">
          <span class="module-name">{{ mod.label }}</span>
          <span class="module-desc">{{ mod.description }}</span>
        </div>
      </label>
    </div>

    <div class="select-footer">
      <span class="selected-count">已选择 {{ selected.length }} 个模块</span>
      <button
        class="btn-primary start-btn"
        :disabled="selected.length === 0"
        @click="$emit('startFill')"
      >
        开始填写
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
// Props & Emits
interface Props {
  selected: string[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:selected', value: string[]): void
  (e: 'startFill'): void
}>()

// 模块定义（含默认勾选标记）
const wizardModules = [
  {
    type: 'basicInfo',
    label: '基础信息',
    icon: '\u{1F464}',
    description: '姓名、联系方式、GitHub 等',
    defaultChecked: true,
  },
  {
    type: 'skills',
    label: '专业技能',
    icon: '\u{1F6E0}',
    description: '按类别展示技术技能',
    defaultChecked: true,
  },
  {
    type: 'projects',
    label: '项目经历',
    icon: '\u{1F4C2}',
    description: '项目描述、角色和亮点',
    defaultChecked: true,
  },
  {
    type: 'evaluation',
    label: '自我评价',
    icon: '\u{270D}',
    description: '个人优势和自我介绍',
    defaultChecked: true,
  },
  {
    type: 'campus',
    label: '校园经历',
    icon: '\u{1F393}',
    description: '校内活动、社团经历',
    defaultChecked: false,
  },
  {
    type: 'internship',
    label: '实习经历',
    icon: '\u{1F4BC}',
    description: '公司实习和职位描述',
    defaultChecked: false,
  },
  {
    type: 'awards',
    label: '获奖',
    icon: '\u{1F3C6}',
    description: '竞赛获奖和荣誉',
    defaultChecked: false,
  },
  {
    type: 'certificates',
    label: '证书',
    icon: '\u{1F4DC}',
    description: '资格证书和考试成绩',
    defaultChecked: false,
  },
]

// 初始化默认勾选（仅首次）
const initialized = { value: false }
if (!initialized.value) {
  initialized.value = true
  const defaults = wizardModules.filter(m => m.defaultChecked).map(m => m.type)
  emit('update:selected', defaults)
}

function toggleModule(type: string) {
  const current = [...props.selected]
  const idx = current.indexOf(type)
  if (idx >= 0) {
    current.splice(idx, 1)
  } else {
    current.push(type)
  }
  emit('update:selected', current)
}
</script>

<style scoped>
.module-select {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.select-header {
  padding: 20px 20px 12px;
}

.select-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--ui-text-primary);
  margin: 0 0 6px 0;
}

.select-desc {
  font-size: 13px;
  color: var(--ui-text-tertiary);
  margin: 0;
}

/* 模块列表 */
.module-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.module-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--ui-bg-tertiary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-md);
  cursor: pointer;
  transition: border-color var(--ui-transition-fast), background-color var(--ui-transition-fast);
}

.module-item:hover {
  border-color: var(--ui-border-hover);
  background: var(--ui-bg-tertiary);
}

.module-item.checked {
  border-color: var(--ui-accent);
  background: rgba(0, 120, 212, 0.08);
}

.module-checkbox {
  width: 16px;
  height: 16px;
  accent-color: var(--ui-accent);
  cursor: pointer;
  flex-shrink: 0;
}

.module-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.module-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.module-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--ui-text-primary);
}

.module-desc {
  font-size: 12px;
  color: var(--ui-text-tertiary);
}

/* 底部操作 */
.select-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--ui-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--ui-bg-secondary);
}

.selected-count {
  font-size: 13px;
  color: var(--ui-text-tertiary);
}

.start-btn {
  padding: 8px 24px;
}

.btn-primary {
  padding: 8px 16px;
  background: var(--ui-accent);
  color: var(--ui-text-inverse);
  border: none;
  border-radius: var(--ui-radius-sm);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color var(--ui-transition-fast);
}

.btn-primary:hover:not(:disabled) {
  background: var(--ui-accent-hover);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
