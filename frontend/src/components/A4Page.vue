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
    <EmptyState
      v-if="pages.length === 0"
      title="开始编写简历"
      description="在编辑器中输入 Markdown，预览将实时显示在这里"
    >
      <template #icon>
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <path d="M12 20h9"/>
          <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>
        </svg>
      </template>
    </EmptyState>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, onUnmounted, ref } from 'vue'
import { renderMarkdown } from '../utils/markdown'
import EmptyState from './EmptyState.vue'
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
const A4_HEIGHT_PX = Math.floor(A4_HEIGHT_MM * MM_TO_PX) // ~1123px

// 估算单行高度的辅助函数（基于默认行高 1.5）
const LINE_HEIGHT_PX = 22
const CHARS_PER_LINE = 45 // A4 宽 210mm，约 45 中文字符/行

/**
 * 估算一段文本渲染后的高度（像素）
 * 原理：计算字符数 / 每行字符数 = 行数，行数 * 行高 = 高度
 */
function estimateTextHeight(text: string): number {
  const lines = Math.ceil(text.length / CHARS_PER_LINE)
  return lines * LINE_HEIGHT_PX
}

// 多页内容分割
const pages = computed(() => {
  if (!props.content) return []

  try {
    // 方案：按自然节分割 Markdown 内容（标题/分隔线/段落边界）
    // 原理：将内容分成若干"节"，每个节尽量填满一页 A4
    const raw = props.content.trim()
    if (!raw) return []

    // 估算总高度，如果能塞进一页就直接返回
    const estimatedTotal = estimateTextHeight(raw)
    if (estimatedTotal <= A4_HEIGHT_PX) {
      return [renderMarkdown(raw)]
    }

    // 多页分割：找到所有自然节边界（h1/h2/h3/hr/空行）
    const sections = splitMarkdownIntoSections(raw)
    if (sections.length <= 1) {
      // 无法分割，退化成单页（内容可能被截断）
      return [renderMarkdown(raw)]
    }

    // 贪心装箱：将节分配到各页，每页尽量填满 A4_HEIGHT_PX
    const pageContents: string[] = []
    let currentPage = ''
    let currentPageHeight = 0

    for (const section of sections) {
      const sectionHeight = estimateTextHeight(section.text)

      // 如果当前节单独能超过一页，直接强制放入（可能溢出）
      if (sectionHeight > A4_HEIGHT_PX && currentPage.length > 0) {
        pageContents.push(currentPage.trim())
        currentPage = ''
        currentPageHeight = 0
      }

      if (currentPageHeight + sectionHeight <= A4_HEIGHT_PX) {
        // 能放入当前页
        currentPage += (currentPage ? '\n\n' : '') + section.text
        currentPageHeight += sectionHeight
      } else {
        // 放不下，先保存当前页
        if (currentPage.trim()) {
          pageContents.push(currentPage.trim())
        }
        currentPage = section.text
        currentPageHeight = sectionHeight
      }
    }

    // 最后一页
    if (currentPage.trim()) {
      pageContents.push(currentPage.trim())
    }

    if (pageContents.length === 0) {
      return [renderMarkdown(raw)]
    }

    return pageContents.map(c => renderMarkdown(c))
  } catch (err) {
    // 渲染异常：显示原始文本 + 错误提示
    console.error('Markdown 渲染异常:', err)
    return [renderMarkdown(props.content)]
  }
})

/**
 * 将 Markdown 内容按自然节边界分割
 * 返回每节文本和类型信息
 */
interface Section {
  text: string
  type: 'h1' | 'h2' | 'h3' | 'hr' | 'blank' | 'content'
}

function splitMarkdownIntoSections(markdown: string): Section[] {
  const lines = markdown.split('\n')
  const sections: Section[] = []
  let currentBlock: string[] = []

  function flushBlock() {
    if (currentBlock.length > 0) {
      const text = currentBlock.join('\n')
      sections.push({ text, type: 'content' })
      currentBlock = []
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]

    // 标题行 → 新节开始
    if (/^#{1,3}\s+/.test(line)) {
      flushBlock()
      sections.push({ text: line, type: 'h1' })
      continue
    }

    // 分隔线 → 新节开始
    if (/^(-{3,}|\*{3,}|_{3,})$/.test(line.trim())) {
      flushBlock()
      sections.push({ text: line, type: 'hr' })
      continue
    }

    // 连续空行 → 合并为段落分隔（但不立即开新节）
    if (line.trim() === '') {
      // 保留空行在当前块中
      currentBlock.push(line)
      continue
    }

    // 非空内容行
    currentBlock.push(line)
  }

  flushBlock()
  return sections
}

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
  background: var(--ui-bg-toolbar);
  padding: 20px;
  box-sizing: border-box;
}

.a4-page {
  position: relative;
  width: 210mm;
  min-height: 297mm;
  margin: 0 auto 20px;
  background: white;
  box-shadow: var(--ui-shadow-md);
  box-sizing: border-box;
  /* Bug修复：移除硬编码 padding，改用 editor.css 全局 CSS 变量 --resume-padding */
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
  font-family: var(--ui-font-sans);
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
  color: var(--ui-text-tertiary);
  font-family: var(--ui-font-sans);
}
</style>
