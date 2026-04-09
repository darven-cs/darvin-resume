<template>
  <div class="wizard-generate">
    <div class="generate-header">
      <h3 class="generate-title">简历预览</h3>
      <p class="generate-desc">确认以下模块内容，点击生成简历</p>
    </div>

    <div class="preview-list">
      <div
        v-for="mod in filledModules"
        :key="mod.type"
        class="preview-module"
      >
        <div class="module-header">
          <span class="module-icon">{{ mod.icon }}</span>
          <span class="module-name">{{ mod.label }}</span>
          <span :class="['module-status', mod.hasData ? 'filled' : 'empty']">
            {{ mod.hasData ? '已填写' : '未填写' }}
          </span>
        </div>
        <div v-if="mod.hasData" class="module-preview">
          <pre class="preview-text">{{ mod.preview }}</pre>
        </div>
      </div>
    </div>

    <div class="generate-footer">
      <button class="btn-secondary" @click="$emit('back')">返回修改</button>
      <button
        class="btn-primary generate-btn"
        :disabled="isGenerating"
        @click="$emit('generate')"
      >
        <span v-if="isGenerating" class="spinner"></span>
        <span v-else>生成简历</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// ============================================================
// Props & Emits
// ============================================================

interface Props {
  selectedModules: string[]
  formData: Record<string, any>
  isGenerating: boolean
}

const props = defineProps<Props>()

defineEmits<{
  (e: 'back'): void
  (e: 'generate'): void
}>()

// ============================================================
// 模块元数据
// ============================================================

const moduleMeta: Record<string, { label: string; icon: string }> = {
  basicInfo: { label: '基础信息', icon: '\u{1F464}' },
  skills: { label: '专业技能', icon: '\u{1F6E0}' },
  projects: { label: '项目经历', icon: '\u{1F4C2}' },
  evaluation: { label: '自我评价', icon: '\u{270D}' },
  campus: { label: '校园经历', icon: '\u{1F393}' },
  internship: { label: '实习经历', icon: '\u{1F4BC}' },
  awards: { label: '获奖', icon: '\u{1F3C6}' },
  certificates: { label: '证书', icon: '\u{1F4DC}' },
}

// ============================================================
// 计算已填写的模块
// ============================================================

const filledModules = computed(() => {
  return props.selectedModules.map(type => {
    const meta = moduleMeta[type] || { label: type, icon: '?' }
    const data = props.formData[type]
    const hasData = checkHasData(type, data)
    const preview = hasData ? generatePreview(type, data) : ''

    return {
      type,
      label: meta.label,
      icon: meta.icon,
      hasData,
      preview,
    }
  })
})

function checkHasData(type: string, data: any): boolean {
  if (!data) return false

  switch (type) {
    case 'basicInfo':
      return Object.entries(data).some(([, v]) => typeof v === 'string' && v.trim() !== '')
    case 'evaluation':
      return !!data.content?.trim()
    case 'skills':
      return (data.categories || []).some(
        (cat: any) => (cat.category && cat.category.trim()) || (cat.skillsText && cat.skillsText.trim())
      )
    default:
      // 列表类模块
      if (data.items) {
        return data.items.some((item: any) =>
          Object.values(item).some((v: any) => typeof v === 'string' && v.trim() !== '')
        )
      }
      return false
  }
}

function generatePreview(type: string, data: any): string {
  switch (type) {
    case 'basicInfo': {
      const parts: string[] = []
      if (data.name) parts.push(`姓名: ${data.name}`)
      if (data.phone) parts.push(`电话: ${data.phone}`)
      if (data.email) parts.push(`邮箱: ${data.email}`)
      if (data.github) parts.push(`GitHub: ${data.github}`)
      if (data.website) parts.push(`网站: ${data.website}`)
      if (data.address) parts.push(`地址: ${data.address}`)
      if (data.summary) parts.push(`简介: ${data.summary}`)
      return parts.join('\n')
    }

    case 'skills': {
      return (data.categories || [])
        .filter((cat: any) => cat.category || cat.skillsText)
        .map((cat: any) => `${cat.category || '未分类'}: ${cat.skillsText || '(空)'}`)
        .join('\n')
    }

    case 'evaluation': {
      return data.content || ''
    }

    default: {
      // 列表类模块
      if (!data.items) return ''
      return data.items
        .map((item: any, idx: number) => {
          const lines = Object.entries(item)
            .filter(([, v]) => typeof v === 'string' && v.trim())
            .map(([k, v]) => `${k}: ${v}`)
          return `[${idx + 1}] ` + lines.join(' | ')
        })
        .join('\n')
    }
  }
}
</script>

<style scoped>
.wizard-generate {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.generate-header {
  padding: 20px 20px 12px;
}

.generate-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--ui-text-primary);
  margin: 0 0 6px 0;
}

.generate-desc {
  font-size: 13px;
  color: var(--ui-text-tertiary);
  margin: 0;
}

/* 预览列表 */
.preview-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preview-module {
  background: var(--ui-bg-tertiary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-md);
  overflow: hidden;
}

.module-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-bottom: 1px solid var(--ui-border);
}

.module-icon {
  font-size: 16px;
}

.module-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--ui-text-primary);
  flex: 1;
}

.module-status {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
}

.module-status.filled {
  background: var(--ui-bg-active);
  color: var(--ui-success);
}

.module-status.empty {
  background: rgba(139, 148, 158, 0.15);
  color: var(--ui-text-tertiary);
}

.module-preview {
  padding: 8px 12px;
}

.preview-text {
  margin: 0;
  font-size: 12px;
  color: var(--ui-text-tertiary);
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
}

/* 底部操作 */
.generate-footer {
  padding: 12px 20px;
  border-top: 1px solid var(--ui-border);
  display: flex;
  justify-content: space-between;
  background: var(--ui-bg-secondary);
}

.generate-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 100px;
  justify-content: center;
}

.btn-primary {
  padding: 8px 20px;
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

.btn-secondary {
  padding: 8px 16px;
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-secondary);
  border: none;
  border-radius: var(--ui-radius-sm);
  font-size: 13px;
  cursor: pointer;
  transition: background-color var(--ui-transition-fast);
}

.btn-secondary:hover {
  background: var(--ui-border);
}

/* 加载动画 */
.spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: var(--ui-text-inverse);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
