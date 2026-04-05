## 全局要求
**全程使用中文回答**：使用中文回答和写文档

<!-- GSD:project-start source:PROJECT.md -->
## Project

**Darvin-Resume**

Darvin-Resume 是一款基于 Markdown 原生、隐私优先、AI 深度协同、开发者友好的本地化简历制作与管理工具。基于 Wails v2 (Go + WebView) 构建，面向计算机专业应届生及毕业3年内的初级开发者，解决排版错乱、多版本管理难、内容润色效率低、简历隐私泄露四大核心痛点。所有数据100%本地存储于 SQLite，用户自备 AI API Key，平台不截留任何请求数据。

**Core Value:** **编辑器预览与PDF导出100%排版一致，所见即所得，零排版焦虑。** 用户在编辑器里看到的效果就是最终导出的效果，这是用户信任的基石。如果这一点做不到，其他所有功能都无意义。

### Constraints

- **Tech Stack**: Wails v2 + Go + Vue 3 + TypeScript + Monaco Editor + SQLite — 已确定不可更换
- **Performance**: 冷启动<2s，内存<200MB，安装包<50MB — Wails 选型的核心优势必须兑现
- **Privacy**: 全量数据本地存储，无强制云端上传 — 产品核心定位的红线
- **Compatibility**: 支持 Windows、macOS、Linux 三平台 — Wails 原生支持
- **AI Dependency**: 用户自备 API Key，产品不内置免费 AI 额度 — 商业模式约束
- **Rendering**: 预览与PDF导出必须100%渲染一致 — 用户信任基石
<!-- GSD:project-end -->

<!-- GSD:stack-start source:STACK.md -->
## Technology Stack

Technology stack not yet documented. Will populate after codebase mapping or first phase.
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, or `.github/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->

<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
