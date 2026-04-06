---
phase: "05-templates-export"
plan: "04"
subsystem: "快照管理"
tags: ["snapshot", "version-control", "diff", "rollback"]
dependency_graph:
  requires: ["05-03"]
  provides: ["EXPT-05", "EXPT-06", "EXPT-07", "EXPT-08", "EXPT-09"]
  affects: ["EditorView", "ExportDialog"]
tech_stack:
  added: ["github.com/sergi/go-diff/diffmatchpatch"]
  patterns: ["Service Layer", "Composable Pattern", "Wails Bridge"]
key_files:
  created:
    - "internal/database/migrations/004_create_snapshots_table.sql"
    - "internal/model/snapshot.go"
    - "internal/service/snapshot.go"
    - "frontend/src/types/snapshot.ts"
    - "frontend/src/composables/useSnapshot.ts"
    - "frontend/src/components/SnapshotSidebar.vue"
  modified:
    - "app.go"
    - "frontend/src/views/EditorView.vue"
    - "frontend/src/wailsjs/go/main/App.d.ts"
    - "frontend/src/wailsjs/wailsjs/go/main/App.d.ts"
    - "frontend/src/wailsjs/wailsjs/go/main/App.js"
    - "wailsjs/go/main/App.d.ts"
decisions:
  - id: "snap-01"
    decision: "快照存储完整数据（JSON + Markdown + Template + CSS），支持完整回滚"
    rationale: "确保回滚时可以恢复到精确的模板和样式状态"
  - id: "snap-02"
    decision: "Diff 使用 diffmatchpatch 库，后端计算 delta，前端使用 diff npm 包展示"
    rationale: "后端生成 delta 减少传输量，前端展示兼容已有的 AI Diff 视图样式"
  - id: "snap-03"
    decision: "回滚前自动创建 rollback 类型快照（标签格式：回滚至「xxx」前自动备份）"
    rationale: "防止误操作导致数据丢失，提供安全保护"
  - id: "snap-04"
    decision: "每简历最多 50 个快照，超出自动删除最旧的"
    rationale: "平衡存储空间和历史价值，防止 DoS（T-05-08）"
metrics:
  duration: "约15分钟"
  completed: "2026-04-06"
---

# Phase 05 Plan 04: 版本快照功能 Summary

**版本快照管理：手动创建快照（标签+备注）、PDF导出后自动创建快照、版本历史列表、两版本 Diff 对比、一键回滚（回滚前自动快照当前）、快照存储完整数据（JSON+Markdown+模板+CSS）、50个上限自动清理。**

## 实现概述

### EXPT-05: 手动创建快照（标签+备注）
- 用户可在工具栏点击「版本快照」按钮打开侧边栏
- 点击「创建快照」按钮，输入标签（必填，最大20字符）和备注（可选，最大100字符）
- 快照保存完整数据：json_data、markdown_content、template_id、custom_css

### EXPT-06: PDF导出后自动创建快照
- ExportDialog 的 `exported` 事件触发 `autoCreateSnapshot`
- 自动快照标签格式：`PDF 导出快照 YYYY/MM/DD`
- triggerType 设置为 `auto_pdf_export`

### EXPT-07: 版本列表 + 两版本 Diff 对比
- 快照列表按时间倒序显示
- 支持多选两个快照进行对比
- 使用 diff npm 包的 `diffLines` 计算可读的 diffLines
- 显示增加行数/删除行数统计

### EXPT-08: 一键回滚（回滚前自动快照当前）
- 点击回滚按钮弹出确认对话框
- 显示「回滚后将自动创建当前版本备份」警告
- 确认后调用 `RollbackToSnapshot`，后端自动创建 rollback 类型快照
- 回滚成功后重新加载简历内容

### EXPT-09: 快照完整存储 + 50个上限
- 快照存储 json_data、markdown_content、template_id、custom_css
- `CleanupOldSnapshots` 在每次创建快照后自动调用
- 保留最新的 50 个快照，删除更旧的

## 提交记录

| 任务 | Commit | 描述 |
| ---- | ------ | ---- |
| Task 1 | e0a34e7 | 创建 snapshots 表迁移和 Snapshot model |
| Task 2 | c2d48ce | 实现 SnapshotService 完整 CRUD + Diff + Rollback |
| Task 3 | b6e0ac0 | 添加 Snapshot bridge 方法到 app.go |
| Task 4 | 2277ec6 | 创建 useSnapshot composable 和类型定义 |
| Task 5 | 1c564e8 | 创建 SnapshotSidebar.vue 侧边栏组件 |
| Task 6 | 88a88e4 | 集成快照功能到 EditorView |
| Chore | be69813 | 更新 Wails 绑定类型和 go.mod |

## 技术细节

### 数据模型
- **Snapshot**: 完整快照数据，用于回滚和查看详情
- **SnapshotListItem**: 列表项，不含完整数据，减少传输量
- **DiffResult**: 对比结果，包含 delta 字符串和统计信息

### 触发类型标签
- `manual`: 蓝色，用户手动创建
- `auto_pdf_export`: 绿色，PDF导出自动创建
- `rollback`: 橙色，回滚前自动备份

### 安全考虑
- T-05-08: MaxSnapshotsPerResume = 50，超出自动清理
- T-05-09: 回滚前自动创建 rollback 类型快照
- T-05-10: SQLite WAL 模式保证数据完整性

## 验证结果

- `go build ./...`: 通过
- `npx vue-tsc --noEmit`: 通过
- `grep -c "useSnapshot\|SnapshotSidebar\|CreateSnapshot" EditorView.vue`: 返回 6

## 已知限制

- Wails 绑定类型文件（App.d.ts/App.js）需要手动更新，因为 `wails generate bindings` 不可用
- 前端 diff 展示使用 diff npm 包（与 AI Diff 视图复用样式），后端 diffmatchpatch 用于生成 delta
