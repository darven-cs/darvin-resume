# Darvin-Resume

## What This Is

Darvin-Resume 是一款基于 Markdown 原生、隐私优先、AI 深度协同、开发者友好的本地化简历制作与管理工具。基于 Wails v2 (Go + WebView) 构建，面向计算机专业应届生及毕业3年内的初级开发者，解决排版错乱、多版本管理难、内容润色效率低、简历隐私泄露四大核心痛点。所有数据100%本地存储于 SQLite，用户自备 AI API Key，平台不截留任何请求数据。

## Core Value

**编辑器预览与PDF导出100%排版一致，所见即所得，零排版焦虑。** 用户在编辑器里看到的效果就是最终导出的效果，这是用户信任的基石。如果这一点做不到，其他所有功能都无意义。

## Current State

### v1.0.0 — Shipped 2026-04-09

**63/63 v1 requirements delivered.** See `.planning/milestones/v1.0.0-ROADMAP.md` for details.

Key shipped features:
- Monaco Editor + markdown-it 统一渲染引擎，预览与 PDF 100% 一致
- Claude Messages API + SSE 流式传输，打字机效果
- 选区 AI 操作（润色/翻译/缩写）+ Diff 对比视图
- AI 对话侧边栏 + 一键解析旧简历
- 4 套内置模板 + 可视化样式调整 + 自定义 CSS 白名单
- PDF 导出（A4 + break-inside:avoid）+ 版本快照管理
- AES-256-GCM 加密存储 + 加密备份恢复
- 深色/浅色主题 + 跟随系统 + 响应式布局
- 快捷键体系 + 自定义快捷键 + Toast 通知 + 全场景异常兜底
- 简历列表 + AI 引导创建 + 回收站 + 30 天恢复 + 自动保存

### Next Milestone

v1.1 — 详见 `/gsd-new-milestone`

## Requirements

### Validated (v1.0.0)

- [x] Monaco Editor 核心Markdown编辑器，100%兼容VS Code默认编辑键位
- [x] 实时预览功能，双栏布局，编辑器预览与PDF导出采用统一渲染引擎确保100%一致
- [x] AI引导式简历创建（分步填写 + AI润色）与空白页创建两种模式
- [x] 选区快捷AI操作（润色/翻译/缩写/重写）+ Diff对比视图
- [x] 行级交互能力（折叠/展开、拖拽排序、快捷菜单）
- [x] AI对话侧边栏（多轮对话、引用选中文本、一键插入）
- [x] AI一键解析纯文本/Markdown旧简历
- [x] Claude Messages API 适配 + 自定义 BaseURL 支持
- [x] API Key AES-256-GCM 加密存储
- [x] 流式传输（SSE协议)，打字机效果实时渲染
- [x] 简历结构化JSON Schema作为唯一可信数据源，JSON-MD双向同步
- [x] 职位目标上下文管理，全文上下文参考开关
- [x] 内置4套标准化模板（极简通用/大厂校招/学术科研/双栏简约）
- [x] 可视化样式调整（主色调/字号/行高/边距/字体）+ 安全白名单CSS自定义
- [x] PDF导出（Wails原生系统打印 + Chromedp可选高级方案）
- [x] A4标准纸张分页线提示，break-inside: avoid 规则
- [x] 简历版本快照管理（手动/自动、Diff对比、一键回滚）
- [x] 数据备份与恢复（手动/自动、加密压缩备份包）
- [x] 简历列表管理（排序/搜索/重命名/复制/删除/回收站）
- [x] 深色/浅色双模式，支持根据系统主题自动切换
- [x] 响应式适配（最小宽度1200px，窗口<1200px自动单栏）
- [x] 自动保存（30秒间隔 + AI操作/页面切换触发）
- [x] AI能力快捷键（Ctrl/Cmd+R润色、+T翻译、+D缩写）
- [x] 全场景异常兜底（网络失败/Token超限/格式异常/中断保留）

### v1.1 Backlog

- [ ] AI智能内容补全
- [ ] 简历内容查重
- [ ] ATS简历优化、岗位匹配度分析
- [ ] PDF/Word简历OCR解析
- [ ] 资深开发者专属模板与引导

### Out of Scope

- 云端简历存储、共享、投递功能 — 坚守本地化隐私定位，不做任何云端强制上传
- 招聘平台对接、内推信息整合 — 不涉足招聘业务
- 企业端简历筛选、人才管理 — 纯 C 端工具定位
- 实时协作编辑 — 本地化工具，无需多人协作

## Context

### 技术栈
- **跨平台框架**：Wails v2 (Golang + WebView2/WebKit)，替代 Electron，安装包<50MB
- **后端**：Golang，负责文件IO、SQLite事务、AI API调用、系统级PDF打印、核心业务逻辑
- **前端**：Vue 3 + TypeScript，Monaco Editor 内核
- **渲染引擎**：Markdown-it，统一负责预览与导出渲染
- **本地数据库**：SQLite 3，全量本地存储
- **PDF导出**：Wails原生系统打印为主，Chromedp无头浏览器为备选高级方案
- **AI接口**：Claude Messages API（首选），支持自定义 BaseURL 兼容 API 中转站

### 架构分层
1. **Bridge层**：Wails绑定层，仅负责前后端通信路由
2. **Domain业务层**：ResumeManager、AIAgentService、TemplateRender、ExportService
3. **Infrastructure基础设施层**：SQLiteStore、PDFRenderer、HTTPClient、EncryptService
4. **Provider适配层**：AI接口适配器、PDF渲染适配器

### 数据规范
- 简历结构化JSON Schema为全系统唯一可信数据源
- 正向同步（JSON→Markdown）自动执行
- 反向同步（Markdown→JSON）需用户手动触发，AI解析+Diff确认
- 敏感数据 AES-256-GCM 加密存储，密钥基于设备唯一标识生成

### 量化验收标准
- 冷启动<2s，内存<200MB，单页PDF导出<3s
- 编辑器预览与PDF导出排版100%一致
- 核心功能崩溃率<0.1%，无数据丢失场景

### 目标用户画像
计算机相关专业应届生/毕业3年内的初级开发者，聚焦校招、社招初级岗位。他们习惯使用 VS Code，期望专业级的编辑体验，对隐私高度敏感，需要AI辅助降低优质简历制作门槛。

## Constraints

- **Tech Stack**: Wails v2 + Go + Vue 3 + TypeScript + Monaco Editor + SQLite — 已确定不可更换
- **Performance**: 冷启动<2s，内存<200MB，安装包<50MB — Wails 选型的核心优势必须兑现
- **Privacy**: 全量数据本地存储，无强制云端上传 — 产品核心定位的红线
- **Compatibility**: 支持 Windows、macOS、Linux 三平台 — Wails 原生支持
- **AI Dependency**: 用户自备 API Key，产品不内置免费 AI 额度 — 商业模式约束
- **Rendering**: 预览与PDF导出必须100%渲染一致 — 用户信任基石

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Wails v2 替代 Electron | 安装包<50MB vs Electron动辄150MB+，内存占用<200MB vs Electron 500MB+ | ✅ Shipped in v1.0.0 |
| Monaco Editor 作为编辑器内核 | 100%兼容VS Code编辑键位，贴合开发者习惯 | ✅ Shipped in v1.0.0 |
| JSON Schema 为唯一数据源 | 避免Markdown直接解析导致的数据错乱，确保模板切换内容不丢失 | ✅ Shipped in v1.0.0 |
| Wails原生打印为主PDF方案 | 无需内置Chromium，压缩安装包体积，降低内存占用 | ✅ Shipped in v1.0.0 |
| Claude Messages API 首选 | 需求文档明确指定，支持自定义BaseURL兼容中转站 | ✅ Shipped in v1.0.0 |
| Markdown-it 统一渲染 | 预览与导出共用同一渲染引擎+CSS，确保100%一致性 | ✅ Shipped in v1.0.0 |

## Evolution

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-04-09 after v1.0.0 milestone completion*
