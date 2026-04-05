# Phase 1 Context: 项目骨架与数据层

## Phase Meta
- **Phase**: 1 of 7
- **Goal**: 项目可以启动运行，简历数据可以持久化存储和检索，为所有后续功能提供数据基础
- **Depends on**: Nothing (first phase)
- **Requirements**: EDIT-04, TMPL-07, TMPL-08

## Success Criteria
1. Wails应用可以正常启动，显示基础页面框架（侧边栏+主内容区）
2. SQLite数据库正确初始化，简历结构化JSON Schema完整定义并可读写
3. 新建一条简历记录后关闭应用重新打开，数据完整保留
4. JSON到Markdown的正向同步可以正确执行，修改JSON字段后生成的Markdown内容准确

## Requirements Detail (from REQUIREMENTS.md)

### EDIT-04: 渲染一致性
编辑器预览与PDF导出采用完全一致的渲染引擎与CSS样式表，确保100%渲染一致
→ 本阶段需要建立渲染引擎基础设施（Markdown-it），为Phase 2的完整编辑器做准备

### TMPL-07: JSON→Markdown正向同步
JSON→Markdown正向同步自动执行（结构化JSON修改后自动生成Markdown）
→ 本阶段核心任务，需要实现完整的JSON→Markdown转换逻辑

### TMPL-08: Markdown→JSON反向同步
Markdown→JSON反向同步需用户手动触发，AI解析+Diff视图展示修改差异，用户确认后执行
→ 本阶段仅建立数据模型基础，反向同步在Phase 3 AI能力接入后实现

## Plans Overview
1. **01-01**: Wails v2项目初始化、Go后端与Vue 3前端项目骨架、基础构建配置
2. **01-02**: SQLite存储层实现、简历结构化JSON Schema定义、基础CRUD操作
3. **01-03**: Bridge层绑定、JSON-Markdown正向同步、基础前端页面框架与路由

## Technology Stack
- **框架**: Wails v2 (Go + WebView)
- **后端**: Go (文件IO、SQLite事务、业务逻辑)
- **前端**: Vue 3 + TypeScript
- **数据库**: SQLite 3
- **渲染**: Markdown-it（统一渲染引擎）

## Architecture
1. **Bridge层**: Wails绑定层，仅负责前后端通信路由
2. **Domain业务层**: ResumeManager、TemplateRender
3. **Infrastructure层**: SQLiteStore
4. **前端**: Vue 3组件 + 路由

## Project Constraints
- 冷启动<2s，内存<200MB，安装包<50MB
- 全量数据本地存储
- 支持 Windows、macOS、Linux 三平台
- 预览与PDF导出100%渲染一致

## Current State
- 项目目录已创建，但没有任何代码文件
- 仅有规划文档：PROJECT.md、ROADMAP.md、REQUIREMENTS.md、STATE.md
- 需要从零搭建整个项目
