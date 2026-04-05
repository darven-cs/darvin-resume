# Requirements: Open-Resume

**Defined:** 2026-04-05
**Core Value:** 编辑器预览与PDF导出100%排版一致，所见即所得，零排版焦虑

## v1 Requirements

### Editor — 核心编辑器

- [ ] **EDIT-01**: Monaco Editor 核心Markdown编辑器，100%兼容VS Code默认编辑键位，支持Markdown语法高亮与基础自动补全
- [ ] **EDIT-02**: 支持撤销/重做、多光标编辑、块选择、查找替换等VS Code标配编辑功能
- [ ] **EDIT-03**: 实时预览功能，默认双栏布局（左编辑右预览），支持拖拽调整双栏宽度比例
- [ ] **EDIT-04**: 编辑器预览与PDF导出采用完全一致的渲染引擎与CSS样式表，确保100%渲染一致
- [ ] **EDIT-05**: 预览区固定显示A4标准纸张分页线，明确标注页面边界
- [ ] **EDIT-06**: 编辑区内容输入后预览区实时同步渲染，延迟<200ms
- [ ] **EDIT-07**: 选区快捷操作 — 用户框选文本后自动弹出浮动工具栏，包含AI润色/翻译/缩写/重写
- [ ] **EDIT-08**: 行级交互 — Markdown标题/列表行首显示折叠/展开图标，支持点击折叠内容块
- [ ] **EDIT-09**: 行级交互 — 支持长按行首图标拖拽排序，快速调整内容块顺序
- [ ] **EDIT-10**: 行级交互 — 点击行首图标弹出快捷菜单（上移/下移/AI重写/删除）
- [ ] **EDIT-11**: Diff对比视图 — AI修改内容后以Diff视图展示差异（删除标红/新增标绿），用户点击接受/拒绝

### Resume — 简历创建与管理

- [ ] **RESM-01**: AI引导式简历创建 — 前置模块勾选（基础信息/专业技能/项目经历/自我评价默认，校园/实习/获奖/证书可选）
- [ ] **RESM-02**: AI引导式创建 — 按自定义模块顺序分步填写，每步提供填写框架与示例
- [ ] **RESM-03**: AI引导式创建 — 用户每步填写后AI实时润色并生成结构化数据，全部完成后生成完整Markdown简历
- [ ] **RESM-04**: AI引导式创建 — 支持单步骤重生成、单步骤跳过、中途退出保存已填写内容
- [ ] **RESM-05**: 空白页创建 — 直接进入双栏编辑界面，生成空白Markdown模板，可随时调用AI能力
- [ ] **RESM-06**: 简历列表管理 — 软件启动后默认进入简历列表页，支持按修改时间排序、搜索、重命名、复制、删除
- [ ] **RESM-07**: 回收站功能 — 删除的简历进入回收站保留30天，支持恢复与永久删除
- [ ] **RESM-08**: 自动保存 — 编辑区内容每30秒自动保存，AI操作完成/页面切换时自动触发保存

### AI — AI智能体

- [ ] **AIAI-01**: Claude Messages API核心适配，支持用户自定义BaseURL兼容API中转站
- [ ] **AIAI-02**: API Key采用AES-256-GCM算法加密存储，密钥基于设备唯一标识生成
- [ ] **AIAI-03**: 支持用户配置默认模型、单次调用最大Token数、请求超时时间
- [ ] **AIAI-04**: 所有AI生成/修改操作以简历结构化JSON Schema为唯一标准，AI输出必须严格符合Schema
- [ ] **AIAI-05**: 职位目标上下文管理 — 简历编辑页顶部可编辑目标岗位，AI调用默认传入
- [ ] **AIAI-06**: 全文上下文参考开关 — 局部润色/翻译默认关闭，全文重写/结构化生成默认开启
- [ ] **AIAI-07**: 流式传输 — 所有AI生成内容支持流式传输（SSE协议），打字机效果实时渲染
- [ ] **AIAI-08**: AI对话侧边栏 — 编辑器右侧唤起，支持多轮对话，可引用选中文本，输出可一键插入编辑区
- [ ] **AIAI-09**: 一键解析 — 用户粘贴纯文本/Markdown旧简历后AI自动解析生成结构化数据+Markdown内容
- [ ] **AIAI-10**: 网络/API调用失败时弹窗提示错误原因，完整保留原始输入，提供一键重试
- [ ] **AIAI-11**: Token超限时自动检测并分批生成，同时弹窗提示用户精简内容
- [ ] **AIAI-12**: AI返回格式异常时先基于Schema校验，不通过自动重试1次，仍失败则完整返回原始内容
- [ ] **AIAI-13**: 用户主动中断生成时保留已生成的有效内容，允许继续编辑或重新生成

### Template — 模板与渲染

- [ ] **TMPL-01**: 内置4套标准化模板 — 极简通用风、大厂校招风、学术科研风、双栏简约风
- [ ] **TMPL-02**: 模板切换仅修改渲染样式，底层结构化JSON数据完全复用，切换后内容无丢失
- [ ] **TMPL-03**: 支持用户将自定义样式保存为个人模板，可复用至新建简历
- [ ] **TMPL-04**: 可视化样式调整 — 滑块/选择器调整主色调、基础字号、行高、页面边距、字体选择
- [ ] **TMPL-05**: 自定义CSS能力 — 安全白名单机制（仅开放字体/颜色/边距/行高/列表样式属性）
- [ ] **TMPL-06**: 一键重置默认样式功能
- [ ] **TMPL-07**: JSON→Markdown正向同步自动执行（结构化JSON修改后自动生成Markdown）
- [ ] **TMPL-08**: Markdown→JSON反向同步需用户手动触发，AI解析+Diff视图展示修改差异，用户确认后执行

### Export — 数据管理与导出

- [ ] **EXPT-01**: PDF导出 — 默认采用Wails原生系统打印接口，无需内置Chromium
- [ ] **EXPT-02**: PDF导出保留Chromedp无头浏览器渲染作为高级可选模式
- [ ] **EXPT-03**: PDF导出严格遵循A4标准纸张尺寸，自动处理break-inside: avoid规则
- [ ] **EXPT-04**: 支持用户自定义导出页码范围、是否隐藏分页线、导出DPI参数
- [ ] **EXPT-05**: 版本快照 — 支持用户手动创建版本快照，可自定义版本标签和备注
- [ ] **EXPT-06**: 版本快照 — 每次成功导出PDF后自动创建版本快照
- [ ] **EXPT-07**: 版本管理 — 查看全量历史版本列表，支持两版本Diff对比
- [ ] **EXPT-08**: 版本管理 — 一键回滚到指定历史版本，回滚前自动为当前内容创建快照
- [ ] **EXPT-09**: 版本快照完整存储结构化JSON、Markdown内容、模板信息、自定义CSS
- [ ] **EXPT-10**: 手动备份 — 一键导出全量数据备份包（加密压缩格式），仅可本软件解密导入
- [ ] **EXPT-11**: 自动备份 — 支持设置自动备份周期（每日/每周/每月），可设置最大备份数量
- [ ] **EXPT-12**: 数据恢复 — 选择备份包一键导入，恢复前自动为当前数据创建备份

### UIUX — 界面与交互

- [ ] **UIUX-01**: 极简主义视觉风格，低饱和度专业色系，聚焦文字内容
- [ ] **UIUX-02**: 深色/浅色双模式适配，支持根据系统主题自动切换
- [ ] **UIUX-03**: 编辑区等宽代码字体，预览区无衬线系统字体，支持用户自定义字体
- [ ] **UIUX-04**: 响应式适配 — 支持窗口大小自由缩放，最小适配宽度1200px
- [ ] **UIUX-05**: 窗口宽度<1200px时自动切换为单栏布局，支持手动切换编辑/预览视图
- [ ] **UIUX-06**: 快捷键体系 — 核心编辑快捷键100%兼容VS Code，AI快捷键（Ctrl/Cmd+R润色/T翻译/D缩写）
- [ ] **UIUX-07**: 支持用户自定义所有快捷键，配置本地存储
- [ ] **UIUX-08**: 空状态设计 — 无简历时首页展示内置模板Demo预览+醒目新建按钮
- [ ] **UIUX-09**: 全异步操作提供明确加载提示，AI流式生成提供打字机动效
- [ ] **UIUX-10**: 全操作即时反馈（toast/弹窗），无无反馈操作
- [ ] **UIUX-11**: 全场景异常状态覆盖（网络异常/API失败/渲染异常/数据加载失败），提供错误提示+可执行方案

## v2 Requirements

### Advanced AI

- **AAI-01**: AI智能内容补全 — 编辑器中智能建议下一步内容
- **AAI-02**: 简历内容查重 — 检测重复/冗余表述
- **AAI-03**: ATS简历优化 — 针对ATS系统优化简历关键词与格式
- **AAI-04**: 岗位匹配度分析 — 分析简历与目标岗位的匹配程度

### Advanced Input

- **AINP-01**: PDF/Word简历OCR解析 — 从现有PDF/Word简历中提取内容
- **AINP-02**: 资深开发者专属简历模板与引导逻辑

## Out of Scope

| Feature | Reason |
|---------|--------|
| 云端简历存储/共享/投递 | 坚守本地化隐私定位，不做任何云端强制上传 |
| 招聘平台对接/内推信息整合 | 不涉足招聘业务 |
| 企业端简历筛选/人才管理 | 纯C端工具定位 |
| 实时协作编辑 | 本地化工具，无需多人协作 |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| EDIT-01 | Phase 2 | Pending |
| EDIT-02 | Phase 2 | Pending |
| EDIT-03 | Phase 2 | Pending |
| EDIT-04 | Phase 1 | Pending |
| EDIT-05 | Phase 2 | Pending |
| EDIT-06 | Phase 2 | Pending |
| EDIT-07 | Phase 3 | Pending |
| EDIT-08 | Phase 2 | Pending |
| EDIT-09 | Phase 2 | Pending |
| EDIT-10 | Phase 2 | Pending |
| EDIT-11 | Phase 3 | Pending |
| RESM-01 | Phase 4 | Pending |
| RESM-02 | Phase 4 | Pending |
| RESM-03 | Phase 4 | Pending |
| RESM-04 | Phase 4 | Pending |
| RESM-05 | Phase 4 | Pending |
| RESM-06 | Phase 4 | Pending |
| RESM-07 | Phase 4 | Pending |
| RESM-08 | Phase 4 | Pending |
| AIAI-01 | Phase 3 | Pending |
| AIAI-02 | Phase 6 | Pending |
| AIAI-03 | Phase 3 | Pending |
| AIAI-04 | Phase 3 | Pending |
| AIAI-05 | Phase 3 | Pending |
| AIAI-06 | Phase 3 | Pending |
| AIAI-07 | Phase 3 | Pending |
| AIAI-08 | Phase 3 | Pending |
| AIAI-09 | Phase 3 | Pending |
| AIAI-10 | Phase 3 | Pending |
| AIAI-11 | Phase 3 | Pending |
| AIAI-12 | Phase 3 | Pending |
| AIAI-13 | Phase 3 | Pending |
| TMPL-01 | Phase 5 | Pending |
| TMPL-02 | Phase 5 | Pending |
| TMPL-03 | Phase 5 | Pending |
| TMPL-04 | Phase 5 | Pending |
| TMPL-05 | Phase 5 | Pending |
| TMPL-06 | Phase 5 | Pending |
| TMPL-07 | Phase 1 | Pending |
| TMPL-08 | Phase 1 | Pending |
| EXPT-01 | Phase 5 | Pending |
| EXPT-02 | Phase 5 | Pending |
| EXPT-03 | Phase 5 | Pending |
| EXPT-04 | Phase 5 | Pending |
| EXPT-05 | Phase 5 | Pending |
| EXPT-06 | Phase 5 | Pending |
| EXPT-07 | Phase 5 | Pending |
| EXPT-08 | Phase 5 | Pending |
| EXPT-09 | Phase 5 | Pending |
| EXPT-10 | Phase 6 | Pending |
| EXPT-11 | Phase 6 | Pending |
| EXPT-12 | Phase 6 | Pending |
| UIUX-01 | Phase 7 | Pending |
| UIUX-02 | Phase 7 | Pending |
| UIUX-03 | Phase 7 | Pending |
| UIUX-04 | Phase 7 | Pending |
| UIUX-05 | Phase 7 | Pending |
| UIUX-06 | Phase 7 | Pending |
| UIUX-07 | Phase 7 | Pending |
| UIUX-08 | Phase 7 | Pending |
| UIUX-09 | Phase 7 | Pending |
| UIUX-10 | Phase 7 | Pending |
| UIUX-11 | Phase 7 | Pending |

**Coverage:**
- v1 requirements: 63 total
- Mapped to phases: 63
- Unmapped: 0

---
*Requirements defined: 2026-04-05*
*Last updated: 2026-04-05 after roadmap creation*
