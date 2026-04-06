<template>
  <div class="preview-container" ref="containerRef">
    <div
      v-for="(pageContent, index) in pages"
      :key="index"
      class="a4-page"
      :class="templateClass"
      :data-page="index + 1"
    >
      <div class="page-content" v-html="pageContent" />
      <div class="page-number">{{ index + 1 }}</div>
    </div>
    <div v-if="pages.length === 0" class="preview-empty">
      开始编写 Markdown，预览将实时显示在这里
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, onUnmounted, ref } from 'vue'
import { renderMarkdown } from '../utils/markdown'
import '../styles/editor.css'
// 静态导入所有模板 CSS（确保打包）
import '../styles/templates/template-minimal.css'
import '../styles/templates/template-dual-col.css'
import '../styles/templates/template-academic.css'
import '../styles/templates/template-campus.css'

const props = defineProps<{
  /** Markdown 内容 */
  content: string
  /** 当前模板 ID，如 'minimal', 'dual-col', 'academic', 'campus' */
  templateId?: string
  /** 用户自定义 CSS */
  customCss?: string
}>()

const containerRef = ref<HTMLElement | null>(null)
const dynamicStyleId = 'a4page-template-dynamic-css'

/** 动态注入自定义 CSS */
function injectCustomCSS(css: string) {
  // 移除旧标签
  const old = document.getElementById(dynamicStyleId)
  if (old) old.remove()

  if (css && css.trim()) {
    const style = document.createElement('style')
    style.id = dynamicStyleId
    style.textContent = css
    containerRef.value?.appendChild(style)
  }
}

// 监听 templateId 和 customCss 变化
watch([() => props.templateId, () => props.customCss], ([, css]) => {
  injectCustomCSS(css || '')
}, { immediate: true })

onMounted(() => {
  injectCustomCSS(props.customCss || '')
})

onUnmounted(() => {
  const old = document.getElementById(dynamicStyleId)
  if (old) old.remove()
})

// A4 尺寸（mm转px, 96dpi: 1mm = 3.78px）
// A4: 210mm × 297mm → 794px × 1123px（近似值）
const A4_WIDTH_MM = 210
const A4_HEIGHT_MM = 297
const MM_TO_PX = 3.78

void A4_WIDTH_MM
void MM_TO_PX
void A4_HEIGHT_MM

// 按 A4 高度分页
const pages = computed(() => {
  if (!props.content) return []
  const html = renderMarkdown(props.content)

  // Phase 2: 仅实现单页 A4 边界线（满足 EDIT-05 页面边界显示要求）
  // Phase 5: 多页分页逻辑将在 PDF 导出时完善（实际页面换行计算）
  return [html]
})

/** 根据 templateId 生成模板 CSS 类名 */
const templateClass = computed(() => {
  const validIds = ['minimal', 'dual-col', 'academic', 'campus']
  const tid = props.templateId
  if (tid && validIds.includes(tid)) {
    return `template-${tid}`
  }
  return 'template-minimal'
})

/**
 * 获取用于导出的 HTML 内容
 * 返回当前渲染的 HTML 字符串
 */
function getHTMLForExport(): string {
  return pages.value.join('')
}

// 暴露方法给父组件
defineExpose({
  getHTMLForExport,
})
</script>

<style scoped>
.preview-container {
  width: 100%;
  height: 100%;
  overflow: auto;
  background: #f0f0f0;
  padding: 20px;
  box-sizing: border-box;
}

.a4-page {
  position: relative;
  width: 210mm;
  min-height: 297mm;
  margin: 0 auto 20px;
  padding: 20mm;
  background: white;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  box-sizing: border-box;
}

/* A4 页面边界线 */
.a4-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border: 1px solid rgba(0, 120, 212, 0.25);
  pointer-events: none;
  z-index: 1;
}

/* 分页底部标注 */
.a4-page::after {
  content: attr(data-page);
  position: absolute;
  bottom: -24px;
  right: 0;
  font-size: 11px;
  color: rgba(0, 120, 212, 0.5);
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

/* 分页符参考线 */
.page-break {
  position: relative;
  border-bottom: 1px dashed rgba(0, 120, 212, 0.3);
  margin-bottom: 8px;
}

.page-number {
  position: absolute;
  bottom: 8mm;
  right: 20mm;
  font-size: 10pt;
  color: #888;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

.preview-empty {
  text-align: center;
  color: #8b949e;
  padding: 3rem;
  font-size: 0.9rem;
}
</style>
