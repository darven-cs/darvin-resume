# Roadmap: Darvin-Resume

## Overview

从零搭建一款本地化 Markdown 简历工具。先建立项目骨架与数据基础，再实现编辑器核心体验，确保"所见即所得"的渲染一致性基石，随后接入 AI 能力让工具产生质变，接着完善模板系统与 PDF 导出，最后打磨界面交互与健壮性，交付一个可信赖的简历制作工具。

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: 项目骨架与数据层** - Wails项目初始化、SQLite存储层、简历结构化JSON Schema、基础导航框架
- [x] **Phase 2: 核心编辑器** - Monaco编辑器集成、双栏布局、实时预览、统一渲染引擎、A4分页线 (completed 2026-04-06)
- [ ] **Phase 3: AI核心能力** - API对接、流式传输、选区快捷操作、Diff对比、AI对话侧边栏、简历解析
- [ ] **Phase 4: 简历创建与管理** - AI引导式创建、空白页创建、简历列表、回收站、自动保存
- [ ] **Phase 5: 模板与导出** - 4套内置模板、样式调整、自定义CSS、PDF导出、版本快照管理
- [ ] **Phase 6: 数据安全与备份** - API Key加密存储、备份与恢复、数据导入导出
- [ ] **Phase 7: 界面打磨与健壮性** - 深色模式、响应式、快捷键体系、异常兜底、空状态设计、全操作反馈

## Phase Details

### Phase 1: 项目骨架与数据层
**Goal**: 项目可以启动运行，简历数据可以持久化存储和检索，为所有后续功能提供数据基础
**Depends on**: Nothing (first phase)
**Requirements**: EDIT-04, TMPL-07, TMPL-08
**Success Criteria** (what must be TRUE):
  1. Wails应用可以正常启动，显示基础页面框架（侧边栏+主内容区）
  2. SQLite数据库正确初始化，简历结构化JSON Schema完整定义并可读写
  3. 新建一条简历记录后关闭应用重新打开，数据完整保留
  4. JSON到Markdown的正向同步可以正确执行，修改JSON字段后生成的Markdown内容准确
**Plans**: 3 planned,详见 `.planning/phase-1/PLAN.md`

Plans:
- [ ] 01-01: Wails v2项目初始化、Go后端与Vue 3前端项目骨架、基础构建配置、侧边栏布局、Markdown-it引擎初始化
- [ ] 01-02: SQLite存储层实现、简历结构化JSON Schema定义、基础CRUD操作、单元测试
- [ ] 01-03: Bridge层绑定、JSON-Markdown正向同步、前端页面框架与路由、集成测试

### Phase 2: 核心编辑器
**Goal**: 用户可以在双栏界面中流畅编写Markdown并实时看到与导出一致的预览效果
**Depends on**: Phase 1
**Requirements**: EDIT-01, EDIT-02, EDIT-03, EDIT-05, EDIT-06, EDIT-08, EDIT-09, EDIT-10
**Success Criteria** (what must be TRUE):
  1. 用户在Monaco编辑器中编写Markdown，编辑体验与VS Code一致（撤销/重做、多光标、查找替换）
  2. 预览区实时渲染Markdown内容，输入到渲染延迟<200ms
  3. 双栏布局支持拖拽调整宽度比例，预览区固定显示A4分页线标注页面边界
  4. 用户可以折叠/展开Markdown内容块，可以拖拽排序内容块，可以点击行首图标弹出快捷菜单
**Plans**: 4 planned,详见 `.planning/phases/02-core-editor/`
**UI hint**: yes

Plans:
- [x] 02-01: Monaco Editor集成、Markdown语法高亮、VS Code兼容键位配置
- [x] 02-02: 双栏布局组件、拖拽调整栏宽、统一渲染引擎(Markdown-it)与CSS样式表
- [x] 02-03: 实时预览同步(延迟<200ms)、A4分页线显示
- [x] 02-04: 行级交互 — 折叠/展开、拖拽排序、快捷菜单

### Phase 3: AI核心能力
**Goal**: 用户可以在编辑过程中随时调用AI能力润色、翻译、改写内容，并通过AI对话获取简历建议
**Depends on**: Phase 2
**Requirements**: AIAI-01, AIAI-03, AIAI-04, AIAI-05, AIAI-06, AIAI-07, AIAI-08, AIAI-09, AIAI-10, AIAI-11, AIAI-12, AIAI-13, EDIT-07, EDIT-11
**Success Criteria** (what must be TRUE):
  1. 用户配置API Key和BaseURL后，AI调用可以成功连接并返回结果，支持流式打字机效果渲染
  2. 用户框选编辑器文本后弹出浮动工具栏，可以执行AI润色/翻译/缩写/重写操作，修改以Diff视图展示差异
  3. 用户可以从编辑器右侧唤起AI对话侧边栏，进行多轮对话，引用选中文本，一键插入AI输出
  4. 用户粘贴旧简历文本后，AI可以自动解析生成结构化数据和Markdown内容
  5. 网络失败、Token超限、格式异常、用户中断等异常场景均有兜底处理，不丢失已输入内容
**Plans**: 6 planned,详见 `.planning/phases/03-ai-core/`
**UI hint**: yes

Plans:
- [x] 03-01: Claude Messages API适配、自定义BaseURL、流式传输(SSE)
- [x] 03-02: 选区快捷AI操作(润色/翻译/缩写/重写)、浮动工具栏
- [x] 03-03: Diff对比视图、AI修改接受/拒绝交互
- [x] 03-04: AI对话侧边栏、多轮对话、引用文本、一键插入
- [x] 03-05: 一键解析旧简历、职位目标上下文、全文上下文开关
- [x] 03-06: AI调用全场景异常兜底(网络/Token/格式/中断)

### Phase 4: 简历创建与管理
**Goal**: 用户可以通过AI引导或空白页两种方式创建简历，并对简历列表进行完整的生命周期管理
**Depends on**: Phase 3
**Requirements**: RESM-01, RESM-02, RESM-03, RESM-04, RESM-05, RESM-06, RESM-07, RESM-08
**Success Criteria** (what must be TRUE):
  1. 用户可以通过AI引导模式分步创建简历（勾选模块、分步填写、AI实时润色），完成后生成完整Markdown简历
  2. 用户可以通过空白页模式直接进入双栏编辑界面，随时调用AI能力
  3. 软件启动后进入简历列表页，用户可以排序、搜索、重命名、复制、删除简历，删除的简历进入回收站30天内可恢复
  4. 编辑过程中内容每30秒自动保存，AI操作完成和页面切换时自动触发保存，无数据丢失
**Plans**: 3 planned,详见 `.planning/phases/04-resume-mgmt/`
**UI hint**: yes

Plans:
- [x] 04-01: 简历列表页 — 列表展示、排序、搜索、重命名、复制、删除
- [x] 04-02: AI引导式创建流程 — 模块勾选、分步填写、AI润色、完成生成
- [x] 04-03: 空白页创建、回收站功能、自动保存机制

### Phase 5: 模板与导出
**Goal**: 用户可以选择模板、调整样式并导出排版精确的PDF，同时可以管理简历的版本历史
**Depends on**: Phase 4
**Requirements**: TMPL-01, TMPL-02, TMPL-03, TMPL-04, TMPL-05, TMPL-06, EXPT-01, EXPT-02, EXPT-03, EXPT-04, EXPT-05, EXPT-06, EXPT-07, EXPT-08, EXPT-09
**Success Criteria** (what must be TRUE):
  1. 用户可以在4套内置模板间自由切换，切换后简历内容完整保留，仅渲染样式变化
  2. 用户可以通过滑块/选择器调整主色调、字号、行高、边距、字体，调整后实时在预览区看到效果
  3. 用户可以将当前样式保存为个人模板，可以输入自定义CSS（仅限白名单属性），可以一键重置默认样式
  4. 用户可以导出PDF（默认系统打印/可选Chromedp），导出严格遵循A4尺寸和分页规则，与预览100%一致
  5. 用户可以创建版本快照（手动/自动），查看版本历史，对比两版本差异，一键回滚
**Plans**: 4 planned,详见 `.planning/phases/05-templates-export/`
**UI hint**: yes

Plans:
- [x] 05-01: 4套内置模板、模板切换机制、个人模板保存
- [x] 05-02: 可视化样式调整、自定义CSS白名单、一键重置
- [x] 05-03: PDF导出 — Wails系统打印、Chromedp备选、A4分页规则、导出参数
- [x] 05-04: 版本快照管理 — 创建/自动创建、历史列表、Diff对比、回滚

### Phase 6: 数据安全与备份
**Goal**: 用户数据受到加密保护，用户可以备份和恢复全部数据，不担心数据丢失
**Depends on**: Phase 5
**Requirements**: AIAI-02, EXPT-10, EXPT-11, EXPT-12
**Success Criteria** (what must be TRUE):
  1. 用户输入的API Key经过AES-256-GCM加密后存储，应用重启后可以正确解密使用
  2. 用户可以一键导出全量数据为加密压缩备份包，可以设置自动备份周期
  3. 用户可以选择备份包一键导入恢复，恢复前自动为当前数据创建备份，恢复后所有数据完整
**Plans**: 3 planned,详见 `.planning/phases/06-security-backup/`

Plans:
- [ ] 06-01: API Key加密存储 — AES-256-GCM、设备唯一标识密钥生成
- [ ] 06-02: 手动备份与恢复 — 加密压缩备份包导出、选择备份导入、恢复前自动备份
- [ ] 06-03: 自动备份 — 周期设置(每日/每周/每月)、最大备份数量

### Phase 7: 界面打磨与健壮性
**Goal**: 应用界面专业美观，支持深色/浅色切换，窗口自适应，所有异常场景有兜底，交互流畅完整
**Depends on**: Phase 6
**Requirements**: UIUX-01, UIUX-02, UIUX-03, UIUX-04, UIUX-05, UIUX-06, UIUX-07, UIUX-08, UIUX-09, UIUX-10, UIUX-11
**Success Criteria** (what must be TRUE):
  1. 应用呈现极简专业视觉风格，深色/浅色模式均可正常使用，支持跟随系统自动切换
  2. 窗口自由缩放时布局正确适配，宽度<1200px自动切换单栏布局，用户可手动切换编辑/预览视图
  3. 所有AI操作支持快捷键（Ctrl/Cmd+R/T/D），用户可自定义快捷键，配置本地持久化
  4. 无简历时首页展示模板Demo预览和醒目新建按钮，所有异步操作有加载提示，所有操作有即时反馈
  5. 全场景异常状态（网络异常/API失败/渲染异常/数据加载失败）均有错误提示和可执行恢复方案
**Plans**: 4 planned,详见 `.planning/phases/07-ui-polish/`
**UI hint**: yes

Plans:
- [ ] 07-01: 视觉风格统一、深色/浅色主题、系统主题跟随
- [ ] 07-02: 响应式布局适配、单栏自动切换
- [ ] 07-03: 快捷键体系、自定义快捷键、空状态设计
- [ ] 07-04: 全局反馈机制(toast/弹窗)、全场景异常兜底

## Progress

**Execution Order:**
Phases execute in numeric order: 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. 项目骨架与数据层 | 3/3 | ✅ Complete | 2026-04-05 |
| 2. 核心编辑器 | 5/4 | Complete   | 2026-04-06 |
| 3. AI核心能力 | 6/6 | ✅ Complete | 2026-04-05 |
| 4. 简历创建与管理 | 3/3 | ✅ Complete | 2026-04-06 |
| 5. 模板与导出 | 4/5 | In Progress|  |
| 6. 数据安全与备份 | 0/3 | Not started | - |
| 7. 界面打磨与健壮性 | 0/4 | Not started | - |
