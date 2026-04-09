# Phase 7: 界面打磨与健壮性 — 总执行计划

**创建时间**: 2026-04-08
**Phase目标**: 应用界面专业美观，支持深色/浅色切换，窗口自适应，所有异常场景有兜底，交互流畅完整
**依赖**: Phase 1-6 全部完成

## 计划拆分

| Plan | 名称 | Wave | 文件 | 状态 |
|------|------|------|------|------|
| 07-01 | 设计令牌与主题系统 | 1 | [07-01-PLAN.md](./07-01-PLAN.md) | ⏳ 待执行 |
| 07-02 | 响应式布局完善与空状态设计 | 2 | [07-02-PLAN.md](./07-02-PLAN.md) | ⏳ 待执行 |
| 07-03 | 快捷键体系与自定义快捷键 | 2 | [07-03-PLAN.md](./07-03-PLAN.md) | ⏳ 待执行 |
| 07-04 | 通用Toast系统与全场景异常兜底 | 3 | [07-04-PLAN.md](./07-04-PLAN.md) | ⏳ 待执行 |

## 执行顺序

```
Wave 1: 07-01 (设计令牌与主题) ← 基础设施，其他计划依赖CSS变量体系
    ↓
Wave 2: 07-02 (响应式+空状态) + 07-03 (快捷键) ← 可并行
    ↓
Wave 3: 07-04 (Toast+异常兜底) ← 收尾，依赖前面的交互模式
```

## 需求覆盖

| Plan | 覆盖需求 |
|------|----------|
| 07-01 | UIUX-01, UIUX-02, UIUX-03 |
| 07-02 | UIUX-04, UIUX-05, UIUX-08 |
| 07-03 | UIUX-06, UIUX-07 |
| 07-04 | UIUX-09, UIUX-10, UIUX-11 |

## 新增文件预估

### 后端 (Go)
- `internal/theme/` — 主题设置管理（使用现有settings基础设施）
- `internal/shortcut/` — 快捷键设置管理
- `app.go` 新增Bridge方法: GetTheme, SetTheme, GetShortcuts, SetShortcuts

### 前端 (Vue/TS)
- `frontend/src/styles/design-tokens.css` — UI设计令牌定义
- `frontend/src/styles/themes/light.css` — 浅色主题覆盖
- `frontend/src/styles/themes/dark.css` — 深色主题覆盖
- `frontend/src/composables/useTheme.ts` — 主题切换逻辑
- `frontend/src/composables/useKeyboard.ts` — 全局快捷键注册
- `frontend/src/composables/useToast.ts` — 通用通知系统
- `frontend/src/components/ToastContainer.vue` — Toast容器组件
- `frontend/src/components/EmptyState.vue` — 通用空状态组件
- `frontend/src/components/SettingsDialog.vue` — 通用设置对话框（多Tab：AI配置/外观/快捷键）
- `frontend/src/components/EmptyState.vue` — 通用空状态组件
- `frontend/src/components/ShortcutSettingsPanel.vue` — 快捷键配置Tab面板

## 关键技术决策

1. **CSS变量分层**: `--ui-*` 控制UI外壳，`--resume-*` 控制简历内容，互不干扰
2. **主题策略**: `data-theme="light|dark"` 属性切换 + CSS变量覆盖，Monaco `vs`/`vs-dark` 同步
3. **快捷键存储**: JSON序列化存入settings表，前端维护默认映射+用户覆盖
4. **Toast系统**: 基于事件总线的全局单例，支持队列、自动消失、手动关闭
5. **后端最小变更**: 仅新增settings key + Bridge方法，复用 `internal/settings` 基础设施
6. **设置集中化**: 升级 `AIConfigModal` → `SettingsDialog`（多Tab：AI配置/外观/快捷键），删除 EditorView 工具栏设置按钮，设置入口仅保留首页
