import MarkdownIt from 'markdown-it'

// 统一渲染引擎实例 — 预览和导出共用
export const md = new MarkdownIt({
  html: false,        // 安全：不解析HTML标签
  breaks: true,       // 换行符转 <br>
  linkify: true,      // 自动识别URL
  typographer: true,  // 优化排版
})

// 导出渲染方法供后续Phase使用
export function renderMarkdown(content: string): string {
  return md.render(content)
}
