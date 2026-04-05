# Phase 3: AI 核心能力 - 讨论日志

## 讨论记录

### 2026-04-05 — Phase 3 规划

**参与者:** Claude (Auto 模式)
**决策点:**

1. **SSE 架构**: Go 后端作为代理，通过 Wails EventsEmit 传输流式数据
2. **API 位置**: 所有 AI 调用走 Go 后端，API Key 不暴露到前端
3. **Diff 对比**: 内联显示在工具栏下方，使用 diff npm 包计算
4. **对话侧边栏**: 400px 右侧面板，消息持久化到 SQLite
5. **依赖安装**: eventsource + diff 两个 npm 包

**技术方案已确定，无需进一步讨论。**
