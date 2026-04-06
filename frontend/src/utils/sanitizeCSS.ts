/**
 * sanitizeCSS.ts - CSS 白名单安全校验工具
 *
 * TMPL-05: 确保用户自定义 CSS 经过严格白名单过滤后
 * 才注入到预览区，防止 XSS、选择器逃逸、危险属性注入。
 *
 * 安全策略：
 * 1. 属性名白名单：仅允许布局、排版、颜色、间距等安全属性
 * 2. 属性值黑名单：拒绝 url()、expression()、javascript:、@import、@keyframes 等危险模式
 * 3. 选择器前缀：仅允许 .page-content 开头的选择器
 * 4. 全部 @ 规则自动移除
 */

import postcss from 'postcss'

// 白名单属性（仅允许布局、排版、颜色、间距）
const ALLOWED_CSS_PROPERTIES = new Set([
  // 字体与文本
  'color', 'font', 'font-family', 'font-size', 'font-weight', 'font-style',
  'text-align', 'text-decoration', 'text-indent', 'text-transform',
  'letter-spacing', 'line-height', 'word-spacing',
  // 间距
  'margin', 'margin-top', 'margin-right', 'margin-bottom', 'margin-left',
  'padding', 'padding-top', 'padding-right', 'padding-bottom', 'padding-left',
  // 布局
  'display', 'flex-direction', 'justify-content', 'align-items', 'flex-wrap',
  'flex', 'flex-grow', 'flex-shrink', 'flex-basis',
  'column-count', 'column-gap', 'column-rule',
  'gap',
  // 尺寸
  'width', 'min-width', 'max-width', 'height', 'min-height', 'max-height',
  // 边框
  'border', 'border-width', 'border-style', 'border-color',
  'border-top', 'border-bottom', 'border-left', 'border-right',
  'border-radius', 'border-top-width', 'border-bottom-width',
  // 背景
  'background-color', 'background',
  // 列表
  'list-style', 'list-style-type',
  // 溢出
  'overflow', 'overflow-x', 'overflow-y',
  // 断行
  'break-inside', 'page-break-inside', '-webkit-column-break-inside',
  // 垂直对齐
  'vertical-align',
  // 空白
  'white-space',
])

// 危险属性（拒绝全部）
const DANGEROUS_PATTERNS = [
  /url\s*\(/i,              // 外部资源
  /expression\s*\(/i,        // CSS 表达式
  /javascript\s*:/i,         // JS 伪协议
  /@import/i,               // 外部样式表
  /@keyframes/i,            // 动画
  /@supports/i,             // 特性查询
  /@media/i,                // 媒体查询（防止逃逸）
  /--\w+:.*url/i,           // CSS 变量中的 url
  /data\s*:/i,              // data: 协议
]

// 允许的选择器前缀（防止选择器逃逸到页面任意元素）
const ALLOWED_SELECTOR_PREFIXES = [
  '.page-content',
  '.page-content ',
]

/**
 * 默认 CSS 变量值（与 editor.css 中的 :root 定义保持同步）
 */
export const DEFAULT_CSS_VARS: Record<string, string> = {
  '--resume-primary-color': '#1a1a1a',
  '--resume-heading-color': '#1a1a1a',
  '--resume-link-color': '#0066cc',
  '--resume-bg-color': '#ffffff',
  '--resume-font-size': '10.5pt',
  '--resume-line-height': '1.6',
  '--resume-padding': '20mm',
  '--resume-font-family': '-apple-system, BlinkMacSystemFont, sans-serif',
}

/**
 * 校验并清理用户输入的自定义 CSS
 * - 移除所有 @ 规则
 * - 仅保留白名单属性
 * - 拒绝包含危险值的声明
 * - 仅允许 .page-content 前缀的选择器
 *
 * @param input 用户输入的原始 CSS 字符串
 * @returns 清理后的安全 CSS 字符串（失败时返回空字符串）
 */
export function sanitizeCustomCSS(input: string): string {
  if (!input || typeof input !== 'string') {
    return ''
  }

  try {
    const root = postcss.parse(input)

    // 移除所有 @ 规则（@import, @keyframes, @media, @supports 等）
    root.walkAtRules((rule) => {
      rule.remove()
    })

    // 遍历所有声明
    root.walkDecls((decl) => {
      const prop = decl.prop.toLowerCase()
      const value = decl.value

      // 1. 属性名白名单检查
      if (!ALLOWED_CSS_PROPERTIES.has(prop)) {
        decl.remove()
        return
      }

      // 2. 危险值模式检查
      for (const pattern of DANGEROUS_PATTERNS) {
        if (pattern.test(value)) {
          decl.remove()
          return
        }
      }

      // 3. 选择器前缀检查（如果不在允许的前缀内，拒绝）
      const parent = decl.parent
      if (parent && parent.type === 'rule') {
        const selector = parent.selector
        const isAllowed = ALLOWED_SELECTOR_PREFIXES.some(prefix =>
          selector === prefix || selector.startsWith(prefix)
        )
        if (!isAllowed) {
          decl.remove()
        }
      }
    })

    return root.toString()
  } catch {
    // 解析失败时返回空字符串（完全拒绝）
    return ''
  }
}

/**
 * 将样式参数转换为 CSS 变量声明字符串
 *
 * @param vars 部分样式参数（Partial<typeof DEFAULT_CSS_VARS>）
 * @returns 格式化的 CSS :root 块字符串
 *
 * @example
 * buildCSSVars({ '--resume-font-size': '12pt', '--resume-primary-color': '#333' })
 * // => ":root {\n  --resume-font-size: 12pt;\n  --resume-primary-color: #333;\n}"
 */
export function buildCSSVars(vars: Partial<Record<string, string>>): string {
  const entries = Object.entries(vars)
  if (entries.length === 0) return ''
  return `:root {\n${entries.map(([k, v]) => `  ${k}: ${v};`).join('\n')}\n}`
}
