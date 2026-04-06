import { ref, computed, type Ref } from 'vue'

/**
 * 模板定义
 */
export interface TemplateDef {
  id: string       // 'minimal' | 'dual-col' | 'academic' | 'campus'
  name: string     // 显示名称
  nameEn: string   // 英文名
  category: 'builtin' | 'user'
  /** 适用场景标签 */
  tag: string
}

/**
 * 内置模板列表
 */
export const BUILTIN_TEMPLATES: TemplateDef[] = [
  { id: 'minimal', name: '极简通用', nameEn: 'Minimal', category: 'builtin', tag: '通用' },
  { id: 'dual-col', name: '双栏简约', nameEn: 'Dual Column', category: 'builtin', tag: '技术岗' },
  { id: 'academic', name: '学术科研', nameEn: 'Academic', category: 'builtin', tag: '学术' },
  { id: 'campus', name: '大厂校招', nameEn: 'Campus', category: 'builtin', tag: '校招' },
]

/**
 * 模板 CSS 文件映射
 */
const TEMPLATE_CSS_MAP: Record<string, string> = {
  minimal: 'template-minimal.css',
  'dual-col': 'template-dual-col.css',
  academic: 'template-academic.css',
  campus: 'template-campus.css',
}

/**
 * 获取模板 CSS 文件名
 */
export function getTemplateCSSFile(templateId: string): string {
  return TEMPLATE_CSS_MAP[templateId] || 'template-minimal.css'
}

/**
 * 获取模板对应的 CSS 类名
 */
export function getTemplateCSSClass(templateId: string): string {
  return `template-${templateId}`
}

/**
 * useTemplate - 模板状态管理 composable
 * 管理当前简历的模板选择、自定义 CSS 加载与保存
 */
export function useTemplate(resumeId: Ref<string>) {
  /** 当前模板 ID */
  const currentTemplateId = ref<string>('minimal')
  /** 用户自定义 CSS */
  const customCss = ref<string>('')
  /** 加载状态 */
  const isLoading = ref(false)

  /**
   * 加载简历的模板设置（从后端获取 templateId 和 customCss）
   */
  async function loadTemplate() {
    isLoading.value = true
    try {
      const resume = await window.go.main.App.GetResume(resumeId.value)
      currentTemplateId.value = (resume as any).templateId || 'minimal'
      customCss.value = (resume as any).customCss || ''
    } catch (err) {
      console.error('[useTemplate] Failed to load template:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 选择模板（更新 templateId 并保存到后端）
   */
  async function selectTemplate(templateId: string) {
    if (templateId === currentTemplateId.value) return
    currentTemplateId.value = templateId
    try {
      await window.go.main.App.UpdateResumeTemplate(resumeId.value, templateId)
    } catch (err) {
      console.error('[useTemplate] Failed to update template:', err)
    }
  }

  /**
   * 保存为个人模板（更新 customCss）
   */
  async function saveAsTemplate(css: string) {
    customCss.value = css
    try {
      await window.go.main.App.UpdateResumeCustomCSS(resumeId.value, css)
    } catch (err) {
      console.error('[useTemplate] Failed to save custom CSS:', err)
    }
  }

  /**
   * 获取当前模板 CSS 文件名
   */
  const currentTemplateCSS = computed(() => {
    return getTemplateCSSFile(currentTemplateId.value)
  })

  /**
   * 获取当前模板 CSS 类名
   */
  const currentTemplateClass = computed(() => {
    return getTemplateCSSClass(currentTemplateId.value)
  })

  return {
    /** 当前模板 ID */
    currentTemplateId,
    /** 用户自定义 CSS */
    customCss,
    /** 加载状态 */
    isLoading,
    /** 加载简历模板设置 */
    loadTemplate,
    /** 选择内置模板 */
    selectTemplate,
    /** 保存自定义 CSS */
    saveAsTemplate,
    /** 当前模板 CSS 文件名 */
    currentTemplateCSS,
    /** 当前模板 CSS 类名 */
    currentTemplateClass,
    /** 所有模板定义 */
    templates: BUILTIN_TEMPLATES,
  }
}
