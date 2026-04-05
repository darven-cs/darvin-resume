# Darvin-Resume

## What This Is

Darvin-Resume 是一款基于 Markdown 原生、隐私优先、AI 深度协同、开发者友好的本地化简历制作与管理工具。基于 Wails v2 (Go + WebView) 构建，面向计算机专业应届生及毕业3年内的初级开发者，解决排版错乱、多版本管理难、内容润色效率低、简历隐私泄露四大核心痛点。所有数据100%本地存储于 SQLite，用户自备 AI API Key，平台不截留任何请求数据。

## Core Value

**编辑器预览与PDF导出100%排版一致，所见即所得，零排版焦虑。** 用户在编辑器里看到的效果就是最终导出的效果，这是用户信任的基石。如果这一点做不到，其他所有功能都无意义。

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Monaco Editor 核心Markdown编辑器，100%兼容VS Code默认编辑键位
- [ ] 实时预览功能，双栏布局，编辑器预览与PDF导出采用统一渲染引擎确保100%一致
- [ ] AI引导式简历创建（分步填写 + AI润色）与空白页创建两种模式
- [ ] 选区快捷AI操作（润色/翻译/缩写/重写）+ Diff对比视图
- [ ] 行级交互能力（折叠/展开、拖拽排序、快捷菜单）
- [ ] AI对话侧边栏（多轮对话、引用选中文本、一键插入）
- [ ] AI一键解析纯文本/Markdown旧简历
- [ ] Claude Messages API 适配 + 自定义 BaseURL 支持
- [ ] API Key AES-256-GCM 加密存储
- [ ] 流式传输（SSE协议），打字机效果实时渲染
- [ ] 简历结构化JSON Schema作为唯一可信数据源，JSON-MD双向同步
- [ ] 职位目标上下文管理，全文上下文参考开关
- [ ] 内置4套标准化模板（极简通用/大厂校招/学术科研/双栏简约）
- [ ] 可视化样式调整（主色调/字号/行高/边距/字体）+ 安全白名单CSS自定义
- [ ] PDF导出（Wails原生系统打印 + Chromedp可选高级方案）
- [ ] A4标准纸张分页线提示，break-inside: avoid 规则
- [ ] 简历版本快照管理（手动/自动、Diff对比、一键回滚）
- [ ] 数据备份与恢复（手动/自动、加密压缩备份包）
- [ ] 简历列表管理（排序/搜索/重命名/复制/删除/回收站）
- [ ] 深色/浅色双模式，支持根据系统主题自动切换
- [ ] 响应式适配（最小宽度1200px，窗口<1200px自动单栏）
- [ ] 自动保存（30秒间隔 + AI操作/页面切换触发）
- [ ] AI能力快捷键（Ctrl/Cmd+R润色、+T翻译、+D缩写）
- [ ] 全场景异常兜底（网络失败/Token超限/格式异常/中断保留）

### Out of Scope

- 云端简历存储、共享、投递功能 — 坚守本地化隐私定位，不做任何云端强制上传
- 招聘平台对接、内推信息整合 — 不涉足招聘业务
- 企业端简历筛选、人才管理 — 纯 C 端工具定位
- PDF/Word OCR解析 — 放入 v2.0 迭代
- AI智能内容补全、查重 — 放入 v1.1 迭代
- ATS简历优化、岗位匹配度分析 — 放入 v1.1 迭代
- 资深开发者专属模板与引导 — 放入 v1.1 迭代

## Context

### 技术栈已确定
- **跨平台框架**：Wails v2 (Golang + WebView2/WebKit)，替代 Electron，安装包<50MB
- **后端**：Golang，负责文件IO、SQLite事务、AI API调用、系统级PDF打印、核心业务逻辑
- **前端**：Vue 3 + TypeScript，Monaco Editor 内核
- **渲染引擎**：Markdown-it，统一负责预览与导出渲染
- **本地数据库**：SQLite 3，全量本地存储
- **PDF导出**：Wails原生系统打印为主，Chromedp无头浏览器为备选高级方案
- **AI接口**：Claude Messages API（首选），支持自定义 BaseURL 兼容 API 中转站

### 架构分层已定义
1. **Bridge层**：Wails绑定层，仅负责前后端通信路由
2. **Domain业务层**：ResumeManager、AIAgentService、TemplateRender、ExportService
3. **Infrastructure基础设施层**：SQLiteStore、PDFRenderer、HTTPClient、EncryptService
4. **Provider适配层**：AI接口适配器、PDF渲染适配器

### 数据规范已明确
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
| Wails v2 替代 Electron | 安装包<50MB vs Electron动辄150MB+，内存占用<200MB vs Electron 500MB+ | — Pending |
| Monaco Editor 作为编辑器内核 | 100%兼容VS Code编辑键位，贴合开发者习惯 | — Pending |
| JSON Schema 为唯一数据源 | 避免Markdown直接解析导致的数据错乱，确保模板切换内容不丢失 | — Pending |
| Wails原生打印为主PDF方案 | 无需内置Chromium，压缩安装包体积，降低内存占用 | — Pending |
| Claude Messages API 首选 | 需求文档明确指定，支持自定义BaseURL兼容中转站 | — Pending |
| Markdown-it 统一渲染 | 预览与导出共用同一渲染引擎+CSS，确保100%一致性 | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

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
*Last updated: 2026-04-05 after initialization*
