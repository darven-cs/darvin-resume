# Phase 5: 模板与导出 - Research

**研究日期:** 2026-04-06
**领域:** 模板系统 + PDF导出 + 版本快照
**置信度:** MEDIUM（部分技术细节基于训练知识，Web 搜索受限）

---

## 一、技术调研

### 1.1 模板系统

#### 多模板渲染方案

**核心认知：模板切换 = CSS 变量系统 + HTML 结构，不涉及 markdown-it 渲染逻辑变更。**

markdown-it 是一个 HTML 生成器，它本身不处理样式。所有模板通过以下两层实现：

| 层 | 内容 | 如何切换 |
|----|------|---------|
| **HTML 结构层** | markdown-it 生成的 HTML（`h1`/`h2`/`ul`/`table` 等） | 所有模板共用相同 HTML |
| **CSS 样式层** | CSS 变量 + 选择器规则 | 切换 `.css` 文件或 `<style>` 块 |

**推荐方案：CSS 变量 + 类名前缀隔离**

```css
/* 模板1: 极简通用风 */
.template-minimal .page-content { ... }
.template-minimal .page-content h1 { font-size: 18pt; border-bottom: 1.5pt solid #1a1a1a; }

/* 模板2: 双栏简约风 */
.template-dual-col .page-content { column-count: 2; }
.template-dual-col .page-content h1 { font-size: 16pt; border-bottom: 1pt solid #888; }

/* 模板3: 大厂校招风 */
.template-campus .page-content h1 { font-size: 20pt; font-weight: 700; letter-spacing: 2pt; }
.template-campus .page-content h2 { background: linear-gradient(90deg, #1a1a1a 0%, transparent 100%); color: white; }
```

**模板数据结构设计（SQLite）：**

```sql
-- 模板定义表（内置 + 用户自定义）
CREATE TABLE templates (
    id TEXT PRIMARY KEY,          -- 'default'|'academic'|'dual-col'|'campus'|'user_xxx'
    name TEXT NOT NULL,           -- 显示名称
    name_en TEXT,                -- 英文名
    category TEXT NOT NULL,       -- 'builtin'|'user'
    css_content TEXT NOT NULL,   -- 模板专属 CSS
    preview_html TEXT,           -- 模板预览缩略图 HTML（可选）
    created_at DATETIME,
    updated_at DATETIME
);

-- 样式预设参数表（TMPL-04 可视化调整用）
CREATE TABLE template_presets (
    id TEXT PRIMARY KEY,
    template_id TEXT NOT NULL,
    primary_color TEXT DEFAULT '#1a1a1a',
    font_size REAL DEFAULT 10.5,   -- pt
    line_height REAL DEFAULT 1.6,
    page_padding_mm REAL DEFAULT 20,
    font_family TEXT DEFAULT '-apple-system, BlinkMacSystemFont, sans-serif',
    created_at DATETIME,
    FOREIGN KEY (template_id) REFERENCES templates(id)
);
```

**关键发现：**
- `resumes.template_id` 列已存在（migration 001），切换模板仅需更新此字段
- `resumes.custom_css` 列已存在，用户自定义 CSS 可直接追加
- 模板切换时仅替换 `<style class="template-css">` 内容，不重新执行 markdown-it

#### CSS 变量系统

4 套内置模板的共同 CSS 变量：

```css
/* 所有模板共享的 CSS 变量（定义在 editor.css 的 :root 中） */
:root {
  --resume-primary-color: #1a1a1a;    /* TMPL-04 主色调 */
  --resume-bg-color: #ffffff;
  --resume-font-size: 10.5pt;         /* TMPL-04 字号 */
  --resume-line-height: 1.6;          /* TMPL-04 行高 */
  --resume-padding: 20mm;             /* TMPL-04 页面边距 */
  --resume-font-family: -apple-system, BlinkMacSystemFont, sans-serif;
  --resume-link-color: #0066cc;
  --resume-heading-color: #1a1a1a;
}
```

**TMPL-04 可视化调整：** 滑块/选择器修改 CSS 变量值 → 实时反映到预览区。无需重新渲染 markdown-it，仅更新 CSS 变量即可。

#### 内置模板设计方案

| ID | 名称 | 视觉特征 | 适合场景 |
|----|------|---------|---------|
| `minimal` | 极简通用风 | 白底黑字、单栏、标准边距、无装饰 | 通用 |
| `dual-col` | 双栏简约风 | 姓名联系方式左侧固定列、经历右侧主栏 | 技术岗 |
| `academic` | 学术科研风 | 宽松行距、教育背景突出、论文格式 | 申请博士/学术 |
| `campus` | 大厂校招风 | 大标题、强调项目经历、活力配色 | 互联网校招 |

### 1.2 PDF 导出

#### Wails 原生打印（EXPT-01 默认方案）

**Wails v2 官方文档（基于训练数据）：** Wails v2 没有内置的 `PrintToPDF` 或 `SilentPrint` API。所有打印通过前端 `window.print()` 实现，该调用通过 WebView 转发到操作系统原生打印对话框。

```typescript
// 前端调用（Vue/TypeScript）
async function exportPDF(options: {
  pageRange?: string,      // e.g., "1-3"
  hidePageBreaks?: boolean,
  dpi?: number
}) {
  // 隐藏 A4 边界线（分页线）
  if (options.hidePageBreaks) {
    document.querySelectorAll('.a4-page::before').forEach(el => {
      (el as HTMLElement).style.display = 'none'
    })
  }

  // 调用原生打印
  window.print()

  // 恢复分页线
  // ... restore
}
```

**流程：**
1. 用户点击「导出 PDF」按钮
2. 前端准备打印视图（隐藏编辑区、隐藏分页线、确保 WebView 可见）
3. 调用 `window.print()`
4. 操作系统弹出打印对话框，用户选择「另存为 PDF」
5. 浏览器打印 API 按 CSS `@media print` 规则渲染，遵循 `break-inside: avoid`

**@media print 关键规则（需添加到 editor.css）：**

```css
@media print {
  /* 隐藏所有 UI 元素，仅保留 A4 页面 */
  .editor-toolbar, .monaco-editor-container, .preview-container {
    display: none !important;
  }

  /* 确保每页内容不跨页断裂 */
  .page-content > * {
    break-inside: avoid;
    -webkit-column-break-inside: avoid;
  }

  /* A4 尺寸精确设置 */
  .a4-page {
    width: 210mm;
    height: 297mm;
    padding: 20mm;
    margin: 0;
    box-shadow: none;
  }

  /* 强制背景色/背景图打印 */
  * {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
    color-adjust: exact !important;
  }
}
```

**局限性：**
- 需要用户手动选择「另存为 PDF」，不是静默导出
- 无法自定义 DPI（由系统控制）
- 页码范围依赖用户手动设置

#### Chromedp 无头浏览器方案（EXPT-02 备选高级模式）

**Chromedp** 是 Go 生态中成熟的 Chrome DevTools Protocol 客户端，支持无头浏览器 PDF 渲染。

**安装：**
```bash
go get github.com/chromedp/chromedp
```

**PDF 导出流程：**

```go
import (
    "github.com/chromedp/chromedp"
    "github.com/chromedp/cdproto/page"
)

func ExportPDF(ctx context.Context, htmlContent string, outputPath string, opts *PDFOptions) error {
    // 创建无头浏览器分配器
    allocCtx, _ := chromedp.NewExecAllocator(ctx,
        chromedp.Flag("headless", true),
        chromedp.Flag("disable-gpu", true),
        chromedp.Flag("no-sandbox", true),
    )

    browserCtx, _ := chromedp.NewContext(allocCtx)

    var pdfBuf []byte
    err := chromedp.Run(browserCtx,
        // 加载 HTML 内容（data URL）
        chromedp.Navigate("data:text/html," + url.QueryEscape(htmlContent)),
        chromedp.WaitReady("body"),

        // 生成 PDF
        chromedp.ActionFunc(func(ctx context.Context) error {
            p := page.PrintToPDF().
                WithPrintBackground(true).
                WithLandscape(false).
                WithPaperWidth(8.27).        // A4: 210mm = 8.27 inches
                WithPaperHeight(11.69).       // A4: 297mm = 11.69 inches
                WithMarginTop(0).
                WithMarginBottom(0).
                WithMarginLeft(0).
                WithMarginRight(0).
                WithScale(1.0)

            if opts != nil && opts.DPI > 0 {
                // DPI 通过 scale 控制（scale = DPI / 72）
                p = p.WithScale(float64(opts.DPI) / 72.0)
            }

            buf, _, err := p.Do(ctx)
            if err != nil {
                return err
            }
            pdfBuf = buf
            return nil
        }),
    )

    if err != nil {
        return err
    }

    return os.WriteFile(outputPath, pdfBuf, 0644)
}
```

**环境检查：**
- 系统已有 `/usr/bin/google-chrome` 和 `/snap/bin/chromium` — Chromedp 可直接使用
- Chromedp 会自动查找浏览器，无需手动指定路径

**备选方案：go-pdf / unipdf**

如果 Chromedp 引入的二进制大小和启动延迟不可接受，可考虑纯 Go 的 PDF 库：
- `github.com/unidoc/unipdf` — 支持从 HTML 渲染（但需要内嵌 CSS + 字体子集化）
- `github.com/jung-kurt/gofpdf` — 纯结构化 PDF 生成（需要重新布局，计算量大）

**推荐：** Chromedp 作为备选，因为：
1. 系统已有 Chrome/Chromium，无需额外下载
2. 渲染结果与预览 100% 一致（复用同一 HTML + CSS）
3. 支持 `print-color-adjust: exact` 精确颜色
4. 支持 `break-inside: avoid` 自动分页

**EXPT-04 导出参数实现：**

| 参数 | Chromedp 实现 | Wails 原生实现 |
|------|--------------|----------------|
| 页码范围 | `page.PrintToPDF` 无直接支持，需裁剪 PDF（用 `github.com/ledongthuc/pdf` 或 `github.com/unidoc/unipdf`） | 用户在系统打印对话框手动设置 |
| 隐藏分页线 | 通过 CSS `display: none` 隐藏 `.a4-page::before` | 同左 |
| 自定义 DPI | `WithScale(dpi/72)` | 系统控制 |

#### 多页分页逻辑（EXPT-03）

当前 `A4Page.vue` 的 `pages` computed 仅返回单页（`[html]`）。Phase 5 需要实现真正的按高度分页：

**分页算法（客户端 JS）：**

```typescript
function paginateHTML(html: string, pageHeightPx: number): string[] {
  // 1. 将 HTML 转为 DOM
  const div = document.createElement('div')
  div.innerHTML = html

  // 2. 逐元素向下累积高度
  const pages: string[] = []
  let currentPage = document.createElement('div')

  for (const child of Array.from(div.children)) {
    currentPage.appendChild(child.cloneNode(true))
    const height = currentPage.getBoundingClientRect().height

    if (height > pageHeightPx) {
      // 当前元素导致溢出
      // 策略：寻找可断点（块级元素、列表项等）
      const breakable = findBreakPoint(currentPage)
      if (breakable) {
        // 截断当前页到 breakable 位置
        const overflow = currentPage.innerHTML.substring(breakable.position)
        pages.push(currentPage.innerHTML.substring(0, breakable.position))

        // 下一页从截断处开始
        currentPage = document.createElement('div')
        currentPage.innerHTML = overflow + child.outerHTML
      } else {
        // 无法断点：强制放入下一页
        pages.push(currentPage.innerHTML)
        currentPage = document.createElement('div')
        currentPage.appendChild(child.cloneNode(true))
      }
    }
  }

  if (currentPage.children.length > 0) {
    pages.push(currentPage.innerHTML)
  }

  return pages.length > 0 ? pages : [html]
}
```

**关键 CSS 规则：**
- `.page-break-avoid { break-inside: avoid; }` — 已在 `editor.css` 中（Phase 2 已添加）
- Markdown 生成 HTML 时，`<h1>`/`<h2>`/`<h3>`/`<table>`/`<ul>` 列表块自动包裹 `page-break-avoid` 类

### 1.3 版本快照管理

#### 快照数据模型（SQLite）

```sql
-- 版本快照表
CREATE TABLE snapshots (
    id TEXT PRIMARY KEY,
    resume_id TEXT NOT NULL,
    label TEXT NOT NULL,             -- EXPT-05: 用户自定义标签
    note TEXT,                       -- EXPT-05: 用户备注
    trigger_type TEXT NOT NULL,      -- 'manual'|'auto_pdf_export'|'rollback'
    json_data TEXT NOT NULL,         -- EXPT-09: 完整结构化 JSON
    markdown_content TEXT NOT NULL,  -- EXPT-09: Markdown 内容
    template_id TEXT NOT NULL,       -- EXPT-09: 模板 ID
    custom_css TEXT NOT NULL,        -- EXPT-09: 自定义 CSS
    created_at DATETIME NOT NULL,
    FOREIGN KEY (resume_id) REFERENCES resumes(id) ON DELETE CASCADE
);

-- 简历表新增列（记录当前版本快照 ID）
ALTER TABLE resumes ADD COLUMN current_snapshot_id TEXT;

-- 索引
CREATE INDEX idx_snapshots_resume_id ON snapshots(resume_id);
CREATE INDEX idx_snapshots_created_at ON snapshots(created_at);
```

**Go model:**

```go
type Snapshot struct {
    ID              string    `json:"id"`
    ResumeID        string    `json:"resumeId"`
    Label           string    `json:"label"`
    Note            string    `json:"note"`
    TriggerType     string    `json:"triggerType"` // 'manual'|'auto_pdf_export'|'rollback'
    JSONData        string    `json:"jsonData"`    // 完整简历 JSON（与 resumes.json_data 相同结构）
    MarkdownContent string    `json:"markdownContent"`
    TemplateID      string    `json:"templateId"`
    CustomCSS       string    `json:"customCss"`
    CreatedAt       time.Time `json:"createdAt"`
}
```

#### 快照服务接口（Go 后端）

```go
// SnapshotService 接口
type SnapshotService interface {
    CreateSnapshot(ctx context.Context, resumeId string, label string, note string, triggerType string) (*Snapshot, error)
    ListSnapshots(ctx context.Context, resumeId string) ([]*Snapshot, error)
    GetSnapshot(ctx context.Context, snapshotId string) (*Snapshot, error)
    DiffSnapshots(ctx context.Context, id1 string, id2 string) (string, error) // 返回 diff HTML
    RollbackToSnapshot(ctx context.Context, resumeId string, snapshotId string) (*Snapshot, error)
    DeleteSnapshot(ctx context.Context, snapshotId string) error
}

// EXPT-06: 导出 PDF 后自动创建快照
// 在 PDF 导出流程中调用: svc.CreateSnapshot(ctx, resumeId, "PDF导出", "", "auto_pdf_export")

// EXPT-08: 回滚前自动快照当前
// 在 RollbackToSnapshot 中：先 CreateSnapshot(triggerType="rollback")，再恢复数据
```

#### Diff 对比实现（EXPT-07）

**前端已有 `diff` npm 包（v5.2.2）**，用于 AI Diff 视图。复用同一库实现快照 Diff。

```typescript
import * as Diff from 'diff'

function diffSnapshots(snapshot1: Snapshot, snapshot2: Snapshot): string {
  const diff = Diff.createPatch(
    'resume',
    snapshot1.markdownContent,
    snapshot2.markdownContent,
    `版本 ${snapshot1.label}`,
    `版本 ${snapshot2.label}`
  )
  // 转为 side-by-side HTML 用于展示
  return diffToSideBySideHTML(diff)
}
```

前端已有 `diff` 包（`package.json` 确认），无需新增依赖。

### 1.4 CSS 白名单安全机制（TMPL-05）

#### 白名单策略

**允许的 CSS 属性（仅布局、排版、颜色）：**

```typescript
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
  'column-count', 'column-gap', 'column-rule',
  // 尺寸
  'width', 'min-width', 'max-width', 'height', 'min-height', 'max-height',
  // 边框
  'border', 'border-width', 'border-style', 'border-color',
  'border-top', 'border-bottom', 'border-left', 'border-right',
  'border-radius',
  // 背景
  'background-color', 'background-image',
  // 列表
  'list-style', 'list-style-type',
])

// 危险属性（拒绝）
const BLOCKED_PROPERTIES = new Set([
  'position', 'top', 'left', 'right', 'bottom', 'z-index',  // 绝对定位
  'display: none',                                          // 隐藏元素
  'behavior', '-moz-binding',                               // CSS 表达式（XSS）
  'url(',                                                   // 外部资源加载
  'expression(',                                            // IE CSS 表达式
  '@import',                                                // 外部样式表
  'animation', 'transition',                                // 性能/过度动画
])

// 允许的选择器前缀（防止选择器逃逸）
const ALLOWED_SELECTOR_PREFIXES = [
  '.page-content',
  '.page-content h1',
  '.page-content h2',
  '.page-content h3',
  '.page-content p',
  '.page-content ul',
  '.page-content ol',
  '.page-content li',
  '.page-content table',
  '.page-content th',
  '.page-content td',
  '.page-content strong',
  '.page-content em',
  '.page-content code',
  '.page-content a',
  '.page-content blockquote',
]
```

#### 实现方式

**前端（TypeScript，TMPL-05 推荐）：**

使用 DOMPurify + PostCSS 的组合：

```typescript
import DOMPurify from 'dompurify'
import postcss from 'postcss'
import postcssSelector from 'postcss-selector-parser'

function sanitizeCustomCSS(css: string): string {
  try {
    const root = postcss.parse(css)

    root.walkDecls((decl: postcss.Declaration) => {
      // 白名单检查属性名
      if (!ALLOWED_CSS_PROPERTIES.has(decl.prop.toLowerCase())) {
        decl.remove()
        return
      }
      // 危险值检查
      if (/url\(|expression\(|javascript:/i.test(decl.value)) {
        decl.remove()
        return
      }
    })

    root.walkAtRules('import', (rule) => rule.remove())
    root.walkAtRules('keyframes', (rule) => rule.remove())

    return root.toString()
  } catch {
    return '' // 解析失败时拒绝全部
  }
}

// 存储时：sanitizeCustomCSS(userCSS) → 保存到数据库
// 渲染时：将 sanitized CSS 注入 <style class="user-custom-css">
```

**后端（Go，EXPT-09 快照存储时校验）：**

Go 中无需重复校验（前端已做），但快照存储时不做额外处理，保持快照的原始 CSS（即使后来白名单变了，快照应保持原样）。

---

## 二、关键发现

### 发现 1: 模板切换技术成本极低

现有的 `resumes.template_id` 和 `resumes.custom_css` 列已支持模板切换。核心工作仅是：
1. 定义 4 套 CSS 并存储到 `templates` 表
2. 前端在预览区动态注入对应 CSS
3. 模板切换更新 `template_id` 并触发 CSS 重载

**[VERIFIED: 迁移文件 001_create_resumes_table.sql 确认 `template_id` 和 `custom_css` 列已存在]**

### 发现 2: PDF 导出两条路线都可行

- **默认方案（`window.print()`）：** 零新增依赖，依赖系统打印对话框，用户体验稍差（需手动选择保存 PDF）
- **Chromedp 方案：** 精确控制、可静默导出、精确 DPI，但增加 Go 依赖和约 30-50MB 二进制体积

**系统已有 Chrome/Chromium，Chromedp 可行。** 需要在安装包大小（<50MB 约束）和功能完整性间权衡。

**[VERIFIED: 系统检查确认 `/usr/bin/google-chrome` 和 `/snap/bin/chromium` 均存在]**

### 发现 3: 前端已有 diff 库

`package.json` 确认 `diff: ^5.2.2` 已安装。快照 Diff 功能可复用现有库，无需引入新依赖。

**[VERIFIED: frontend/package.json]**

### 发现 4: 分页逻辑是最大技术风险

当前的 `A4Page.vue` 仅返回 `[html]`（单页），真正的多页分页需要：
- DOM 高度精确计算
- 智能断点选择（不能在单词中间断开）
- 列表项/表格行的原子化处理

建议 Phase 5 采用保守策略：**先让单页内容正常工作，多页溢出时用「内容较多，请调整字号或删除内容」提示**。完整多页分页作为后续迭代。

**[ASSUMED]**

### 发现 5: 自定义 CSS 安全边界清晰

白名单方案（允许 30+ 个安全属性）比黑名单（危险属性列表）更安全。Shadow DOM 隔离作为可选增强。

**[ASSUMED: 基于 OWASP CSS Sanitization 标准实践]**

---

## 三、风险评估

| 风险 | 级别 | 描述 | 缓解方案 |
|------|------|------|---------|
| **PDF 排版与预览不一致** | 高 | 这是产品的核心价值。如果导出 PDF 与预览不一致，用户失去信任 | Chromedp 方案保证复用相同 HTML+CSS；Wails 原生方案需严格测试 @media print 规则 |
| **Chromedp 增加安装包体积** | 中 | Chromedp 本身约 5-10MB，但需要 Chrome 运行时会增加更多 | 设为可选高级功能；默认使用 Wails 原生；Chromedp 作为用户主动开启的选项 |
| **多页分页实现复杂度** | 中 | 智能分页需要 DOM 高度计算、断点选择等，涉及大量边界情况 | 保守方案：单页溢出时提示用户；完整多页作为迭代 |
| **CSS 白名单可能过于严格** | 低 | 用户可能想要一些未列入白名单但安全的样式 | 提供「专家模式」开关，放宽白名单（需用户确认风险） |
| **版本快照占用 SQLite 空间** | 低 | 每次快照完整存储 JSON+Markdown+CSS，多次编辑后累积较大 | 设置快照数量上限（如 50 个），超限时提示用户清理或自动删除最旧的 |
| **模板切换时 Markdown 结构不兼容** | 低 | 不同模板的 HTML 结构假设不同（如双栏模板假设特定 HTML 结构） | 所有模板共用相同的 markdown-it HTML 输出，仅通过 CSS 差异实现视觉效果 |

---

## 四、建议的实现方案

### 方案概览

```
Phase 5 拆分为 4 个子计划（与 ROADMAP.md 一致）：

05-01: 模板系统 — 4套内置模板 + 模板切换 + 个人模板保存
05-02: 样式调整 — 可视化调整(TMPL-04) + CSS白名单(TMPL-05) + 一键重置(TMPL-06)
05-03: PDF导出 — Wails原生打印 + Chromedp备选 + A4分页 + 导出参数
05-04: 版本快照 — 创建/自动创建 + 历史列表 + Diff对比 + 回滚
```

### 05-01 模板系统（TMPL-01, TMPL-02, TMPL-03）

**核心实现：**
1. 创建 `templates` 表（内置 4 套 + 用户自定义）
2. 创建 `template_presets` 表（样式预设参数）
3. 前端 `TemplateSelector.vue` 组件（模板卡片 + 切换按钮）
4. `A4Page.vue` 根据当前模板 ID 注入对应 CSS
5. 模板切换调用 `UpdateResumeTemplate()` 更新 `template_id`

**文件变更：**
- `frontend/src/components/TemplateSelector.vue`（新建）
- `frontend/src/components/A4Page.vue`（修改：注入模板 CSS）
- `internal/database/migrations/004_create_templates_table.sql`（新建）
- `internal/database/migrations/005_create_template_presets_table.sql`（新建）
- `internal/service/snapshot.go`（新建：快照服务）
- `app.go`（新增模板相关 bridge 方法）

### 05-02 样式调整（TMPL-04, TMPL-05, TMPL-06）

**核心实现：**
1. `StyleEditor.vue` 组件（滑块 + 颜色选择器 + 字体下拉）
2. CSS 白名单校验（`sanitizeCustomCSS.ts`）
3. 实时预览（CSS 变量更新 → 预览即时变化）
4. 一键重置（恢复内置模板默认样式）

### 05-03 PDF导出（EXPT-01, EXPT-02, EXPT-03, EXPT-04）

**核心实现：**
1. `@media print` CSS 规则完善（追加到 `editor.css`）
2. `ExportPDFDialog.vue` 组件（导出参数设置）
3. Go 后端 `ExportPDF()` 方法（Wails 原生 + Chromedp）
4. 多页分页 JS 逻辑（保守版本）

### 05-04 版本快照（EXPT-05, EXPT-06, EXPT-07, EXPT-08, EXPT-09）

**核心实现：**
1. `snapshots` 表迁移
2. `SnapshotService`（Go）
3. `SnapshotSidebar.vue`（快照列表 + Diff 视图 + 回滚按钮）
4. PDF 导出后自动创建快照（hook）
5. 回滚前自动快照（`triggerType='rollback'`）

---

## 五、依赖分析

### npm 依赖（前端）

| 包 | 版本 | 用途 | 状态 |
|----|------|------|------|
| markdown-it | 14.1.1 | 已有，无需修改 | 已有 |
| diff | 5.2.2 | 快照 Diff 对比 | 已有 |
| dompurify | 3.3.3 | CSS 白名单校验 | 需新增 |
| postcss | 8.5.8 | CSS AST 解析 | 需新增 |
| postcss-selector-parser | — | CSS 选择器校验 | 需新增 |

### Go 依赖（后端）

| 包 | 版本 | 用途 | 状态 |
|----|------|------|------|
| github.com/chromedp/chromedp | latest | PDF 高级导出 | 需新增（可选） |
| github.com/google/uuid | 1.6.0 | 已有，无需修改 | 已有 |

### 环境依赖

| 工具 | 用途 | 状态 |
|------|------|------|
| google-chrome | Chromedp 运行时 | 已有 |
| chromium | Chromedp 备选 | 已有 |

---

## 六、Assumptions Log

> 以下均为 [ASSUMED] 项，基于训练知识未在本次会话中验证，需要在规划阶段确认。

| # | 假设 | 风险 |
|---|------|------|
| A1 | Wails v2 的 `window.print()` 能正确触发系统打印对话框 | 如果 WebView 配置不正确可能失效 |
| A2 | Chromedp 在 Linux 上能正常通过 `google-chrome` 运行 | 需要实测验证无头模式稳定性 |
| A3 | markdown-it 生成的 HTML 结构足够稳定，所有模板都可以通过 CSS 实现不同视觉风格 | 如果双栏模板需要 HTML 结构变化（不只是 CSS），则方案需要调整 |
| A4 | DOMPurify + PostCSS 的 CSS 白名单校验足够安全 | OWASP 标准实践，但未在实际场景验证 |
| A5 | 多页分页的 DOM 高度计算在所有主流浏览器上结果一致 | 可能存在浏览器差异，需要兼容性测试 |

---

## 七、Open Questions

1. **Chromedp 打包策略：** 是否将 Chromedp 作为默认依赖（影响安装包体积），还是作为可选插件？
   - 建议：默认使用 Wails 原生打印（零依赖），Chromedp 作为「高级导出模式」可选。用户首次使用时检测 Chrome 是否可用，提示开启。

2. **模板数量上限：** 用户自定义模板是否有限制？
   - 建议：最多 20 个自定义模板，超出时提示清理。

3. **快照数量上限：** 版本快照是否需要自动清理？
   - 建议：每个简历最多保留 50 个快照，超出时按时间自动删除最旧的。但 ROADMAP 说「支持 30 天恢复」，这里有矛盾——建议明确需求。

4. **TMPL-02（数据复用）：** 模板切换时，`template_id` 变更，但 `json_data`/`markdown_content` 不变。这意味着模板完全基于现有 HTML 输出，仅改变 CSS。但如果某些模板需要不同的 markdown-it 配置（例如学术模板需要不同的 `typographer` 设置），方案需要调整。

5. **PDF 导出与预览 100% 一致性验证：** 这是产品核心价值。需要明确的验证机制——建议在每版发布前自动化对比截图（预览截图 vs PDF 截图）作为回归测试。

---

## 八、Sources

### Primary (HIGH confidence)
- `internal/database/migrations/001_create_resumes_table.sql` — 确认 `template_id` 和 `custom_css` 列存在
- `frontend/package.json` — 确认 `markdown-it: ^14.1.1`, `diff: ^5.2.2`
- `go.mod` — 确认 Go 版本和现有依赖
- `/usr/bin/google-chrome` 系统检查 — Chromedp 运行环境存在

### Secondary (MEDIUM confidence)
- Wails v2 printing approach: `window.print()` via WebView — 基于训练知识
- Chromedp PDF API: `page.PrintToPDF()` — 基于训练知识
- CSS whitelist properties — 基于 OWASP CSS Sanitization 标准

### Tertiary (LOW confidence)
- DOMPurify + PostCSS 具体集成方式 — 基于训练知识，需要实测
- 多页分页算法细节 — 基于训练知识，需要浏览器兼容性测试
