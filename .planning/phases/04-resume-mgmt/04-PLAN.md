# Plan 04: 简历创建与管理

**Phase:** 4 of 7
**Goal:** 用户可以通过AI引导或空白页两种方式创建简历，并对简历列表进行完整的生命周期管理
**Depends on:** Phase 3 complete
**Requirements:** RESM-01~08

## Success Criteria (what must be TRUE)
1. 用户可以通过AI引导模式分步创建简历（勾选模块、分步填写、AI实时润色），完成后生成完整Markdown简历
2. 用户可以通过空白页模式直接进入双栏编辑界面，随时调用AI能力
3. 软件启动后进入简历列表页，用户可以排序、搜索、重命名、复制、删除简历，删除的简历进入回收站可恢复
4. 编辑过程中内容每30秒自动保存，AI操作完成和页面切换时自动触发保存，无数据丢失

## Plans

| Plan | Name | Status | Dependencies |
|------|------|--------|-------------|
| 04-01 | 简历列表页 — 卡片网格、排序、搜索、CRUD | ⏳ Planned | — |
| 04-02 | AI引导式创建流程 — 侧边栏向导 | ⏳ Planned | 04-01 |
| 04-03 | 空白页创建、回收站、自动保存 | ⏳ Planned | 04-01, 04-02 |

## New Dependencies

### npm (frontend)
无新增依赖

### go.mod (backend)
无新增依赖

## New Files

### Backend (Go)
```
internal/
├── service/
│   └── resume.go          # 新增: RenameResume, DuplicateResume, RestoreResume, PermanentDeleteResume, ListDeletedResumes
└── database/
    └── db.go              # 新增: 对应 SQL 查询
app.go                     # 新增: Bridge 绑定
```

### Frontend (Vue/TS)
```
frontend/src/
├── components/
│   ├── CreateModeModal.vue        # 创建模式选择弹窗（AI引导/空白页）
│   ├── ResumeWizardSidebar.vue    # AI引导式创建侧边栏向导
│   ├── ResumeCard.vue             # 简历卡片组件
│   ├── RecycleBinSection.vue      # 回收站折叠区
│   └── SaveStatusIndicator.vue    # 保存状态指示器
├── composables/
│   ├── useAutoSave.ts             # 自动保存 composable
│   └── useResumeList.ts           # 简历列表状态管理
├── views/
│   ├── HomeView.vue               # 重构: 卡片网格 + 搜索 + 回收站
│   └── EditorView.vue             # 重构: 抽取自动保存 + 添加导航 + 标题编辑
└── router/
    └── index.ts                   # 更新路由守卫
```

## Execution Order

```
Wave 1:
  04-01: 简历列表页 (后端新方法 + 前端卡片网格) ← 基础设施

Wave 2 (串行，04-02 先于 04-03):
  04-02: AI引导式创建流程 ← 依赖 04-01 的 CreateModeModal（修改 EditorView 集成向导）
  04-03: 回收站 + 自动保存 ← 依赖 04-01 的后端方法（修改 EditorView 抽取保存逻辑，需在 04-02 之后避免编辑冲突）
```
