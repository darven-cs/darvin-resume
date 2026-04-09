<template>
  <div class="template-selector">
    <div class="selector-header">
      <span class="selector-title">选择模板</span>
    </div>
    <div class="template-grid">
      <div
        v-for="tmpl in templates"
        :key="tmpl.id"
        class="template-card"
        :class="{ active: tmpl.id === modelValue }"
        @click="handleSelect(tmpl.id)"
        :title="`${tmpl.name} (${tmpl.nameEn})`"
      >
        <!-- 微型简历预览 -->
        <div class="template-preview" :class="`preview-${tmpl.id}`">
          <!-- 通用预览内容（各模板有不同样式变体） -->
          <div class="preview-content">
            <div class="preview-header">
              <div class="preview-name">张三</div>
              <div class="preview-contact">📧 邮箱 | 📱 电话</div>
            </div>
            <div class="preview-section">
              <div class="preview-section-title">工作经历</div>
              <div class="preview-line short" />
              <div class="preview-line medium" />
              <div class="preview-line short" />
            </div>
            <div class="preview-section">
              <div class="preview-section-title">项目经历</div>
              <div class="preview-line medium" />
              <div class="preview-line short" />
            </div>
          </div>
        </div>
        <!-- 模板信息 -->
        <div class="template-info">
          <span class="template-name">{{ tmpl.name }}</span>
          <span class="template-tag">{{ tmpl.tag }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BUILTIN_TEMPLATES, type TemplateDef } from '../composables/useTemplate'

const props = defineProps<{
  /** 当前选中的模板 ID */
  modelValue: string
}>()

const emit = defineEmits<{
  /** 模板选中事件 */
  'update:modelValue': [templateId: string]
  /** 模板切换事件 */
  'change': [templateId: string]
}>()

/** 模板列表 */
const templates: TemplateDef[] = BUILTIN_TEMPLATES

/**
 * 处理模板选择
 */
function handleSelect(templateId: string) {
  if (templateId === props.modelValue) return
  emit('update:modelValue', templateId)
  emit('change', templateId)
}
</script>

<style scoped>
.template-selector {
  padding: 12px;
  user-select: none;
}

.selector-header {
  margin-bottom: 10px;
}

.selector-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--ui-text-tertiary);
}

/* 网格布局：大屏4列，中屏2列，小屏1列 */
.template-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}

@media (max-width: 600px) {
  .template-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* 卡片基础样式 */
.template-card {
  display: flex;
  flex-direction: column;
  border-radius: var(--ui-radius-md);
  cursor: pointer;
  transition: transform var(--ui-transition-fast), box-shadow var(--ui-transition-fast), border-color var(--ui-transition-fast), background-color var(--ui-transition-fast);
  border: 2px solid transparent;
  background: #fff;
  overflow: hidden;
}

.template-card:hover {
  transform: scale(1.03);
  box-shadow: var(--ui-shadow-md);
}

.template-card.active {
  border-color: var(--ui-accent);
  background-color: rgba(0, 120, 212, 0.06);
}

/* 微型简历预览容器 */
.template-preview {
  width: 100%;
  height: 120px;
  background: #fff;
  padding: 6px;
  box-sizing: border-box;
  overflow: hidden;
}

.template-preview .preview-content {
  width: 100%;
  height: 100%;
  background: #fff;
  border: 1px solid var(--ui-border);
  padding: 6px;
  box-sizing: border-box;
}

.preview-header {
  margin-bottom: 5px;
}

.preview-name {
  font-size: 8px;
  font-weight: 700;
  margin-bottom: 2px;
}

.preview-contact {
  font-size: 5px;
  color: #666;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.preview-section {
  margin-bottom: 4px;
}

.preview-section-title {
  font-size: 6px;
  font-weight: 600;
  border-bottom: 0.5px solid var(--ui-text-secondary);
  margin-bottom: 2px;
  padding-bottom: 1px;
}

.preview-line {
  height: 3px;
  background: var(--ui-border);
  border-radius: 1px;
  margin-bottom: 2px;
}

.preview-line.short {
  width: 60%;
}

.preview-line.medium {
  width: 80%;
}

/* ===== 各模板预览的视觉差异 ===== */

/* 极简通用风：单栏，标准样式 */
.preview-minimal .preview-content {
  font-size: 5px;
  line-height: 1.4;
}

.preview-minimal .preview-name {
  border-bottom: 0.5px solid #1a1a1a;
  padding-bottom: 1px;
}

/* 双栏简约风：双栏分割线 */
.preview-dual-col .preview-content {
  display: flex;
  gap: 3px;
  font-size: 4.5px;
}

.preview-dual-col .preview-header {
  width: 35%;
  border-right: 0.5px solid var(--ui-border);
  padding-right: 3px;
  flex-shrink: 0;
}

.preview-dual-col .preview-section {
  flex: 1;
}

/* 学术科研风：宽松行距，居中标题 */
.preview-academic .preview-name {
  text-align: center;
  font-size: 7px;
  border-bottom: 0.3px solid #1a1a1a;
}

.preview-academic .preview-section-title {
  border-left: 1.5px solid #2563eb;
  padding-left: 2px;
  font-style: italic;
}

/* 大厂校招风：左侧色条装饰，大标题 */
.preview-campus .preview-name {
  font-size: 8px;
  font-weight: 700;
  letter-spacing: 0.5px;
  border-left: 2px solid #1a1a1a;
  padding-left: 3px;
}

.preview-campus .preview-section-title {
  border-left: 1.5px solid #1a1a1a;
  padding-left: 2px;
}

/* 模板信息区域 */
.template-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  background: #fafafa;
}

.template-card.active .template-info {
  background: rgba(0, 120, 212, 0.08);
}

.template-name {
  font-size: 11px;
  font-weight: 500;
  color: var(--ui-text-primary);
}

.template-card.active .template-name {
  color: var(--ui-accent);
  font-weight: 600;
}

.template-tag {
  font-size: 9px;
  color: var(--ui-text-tertiary);
  background: #f0f0f0;
  padding: 1px 5px;
  border-radius: 10px;
}

.template-card.active .template-tag {
  background: rgba(0, 120, 212, 0.15);
  color: var(--ui-accent);
}
</style>
