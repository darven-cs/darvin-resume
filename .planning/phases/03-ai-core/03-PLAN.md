# Plan 03: AI 核心能力

**Phase:** 3 of 7
**Goal:** 用户可以在编辑过程中随时调用AI能力润色、翻译、改写内容，并通过AI对话获取简历建议
**Depends on:** Phase 2 complete
**Requirements:** AIAI-01~13, EDIT-07, EDIT-11

## Success Criteria (what must be TRUE)
1. 用户配置API Key和BaseURL后，AI调用可以成功连接并返回结果，支持流式打字机效果渲染
2. 用户框选编辑器文本后弹出浮动工具栏，可以执行AI润色/翻译/缩写/重写操作，修改以Diff视图展示差异
3. 用户可以从编辑器右侧唤起AI对话侧边栏，进行多轮对话，引用选中文本，一键插入AI输出
4. 用户粘贴旧简历文本后，AI可以自动解析生成结构化数据和Markdown内容
5. 网络失败、Token超限、格式异常、用户中断等异常场景均有兜底处理，不丢失已输入内容

## Plans

| Plan | Name | Status | Dependencies |
|------|------|--------|-------------|
| 03-01 | Claude API 适配与 SSE 流式传输 | ⏳ Planned | — |
| 03-02 | 选区浮动工具栏 | ⏳ Planned | 03-01 |
| 03-03 | Diff 对比视图 | ⏳ Planned | 03-01, 03-02 |
| 03-04 | AI 对话侧边栏 | ⏳ Planned | 03-01 |
| 03-05 | 一键解析旧简历 & 上下文管理 | ⏳ Planned | 03-01 |
| 03-06 | AI 异常兜底处理 | ⏳ Planned | 03-01 |

## New Dependencies

### npm (frontend)
```json
{
  "eventsource": "^2.3.7",
  "diff": "^5.2.0"
}
```

### go.mod (backend)
- 无需新增（标准库 `net/http` + `encoding/json` 足够）

## New Files

### Backend (Go)
```
backend/internal/ai/
├── client.go       # Claude API 客户端、流式传输
├── config.go        # AIConfig 模型
├── error.go         # 错误类型映射
└── prompt.go        # Prompt 模板
```

### Frontend (Vue/TS)
```
frontend/src/
├── components/
│   ├── AIFloatingToolbar.vue   # 选区浮动工具栏
│   ├── AIDiffView.vue          # Diff 对比视图
│   ├── AIChatSidebar.vue       # AI 对话侧边栏
│   ├── AIConfigModal.vue       # AI 配置弹窗
│   ├── AIErrorToast.vue        # 错误通知
│   ├── ResumeParserModal.vue   # 简历解析弹窗
│   └── JobTargetChip.vue        # 职位目标 Chip
├── composables/
│   ├── useAIStream.ts          # SSE 流处理
│   ├── useAISelection.ts       # 选区状态管理
│   ├── useAIConfig.ts          # 配置管理
│   └── useAIError.ts           # 统一错误处理
├── services/
│   └── ai.ts                   # AI 服务层
└── types/
    └── ai.ts                   # AI 类型定义
```

## Execution Order

```
Wave 1 (并行):
  03-01: SSE 基础设施 ←——— 所有其他计划的依赖

Wave 2 (串行):
  03-02: 浮动工具栏
  03-03: Diff 对比视图
  03-04: 对话侧边栏
  03-05: 简历解析
  03-06: 异常兜底（可并行于 Wave 2 任意计划）
```
