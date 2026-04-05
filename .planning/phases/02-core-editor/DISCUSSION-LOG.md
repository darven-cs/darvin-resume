# Phase 2 Discussion Log

**Date:** 2026-04-05
**Mode:** --auto (automated decision selection)
**Phase:** 02-core-editor

## Discussion Summary

This discussion log captures the decision-making process for Phase 2: 核心编辑器 (Core Editor). The `--auto` flag was used, which automatically selected recommended defaults for all gray areas based on project requirements and technical constraints.

## Gray Areas Identified

### 1. Monaco Editor Integration (EDIT-01, EDIT-02)
**Questions Covered:**
- How to load Monaco Editor in Wails environment?
- Theme configuration approach?
- Font size and line-height defaults?
- Language configuration for Chinese users?

**Auto-Selected Decisions:**
- **Monaco Loading:** Use npm package (@monaco-editor/loader) bundled with webpack
  - *Rationale:* Ensures offline functionality and faster startup vs CDN
  - *Alignment with:* PROJECT.md constraint "冷启动<2s" and privacy-first local-only approach
- **Theme:** VS Code Dark theme as default, light theme variant for Phase 7
  - *Rationale:* Matches target user profile (VS Code users)
  - *Alignment with:* PROJECT.md target user: "习惯使用 VS Code" developers
- **Typography:** 14px font size, 1.6 line-height
  - *Rationale:* VS Code defaults for optimal readability
- **Language:** Enable Chinese language support
  - *Rationale:* PROJECT.md specifies Chinese user base

### 2. Dual-Pane Layout (EDIT-03)
**Questions Covered:**
- Default width ratio between editor and preview?
- Minimum width constraints for each pane?
- Drag handle implementation details?
- Mobile/single-pane breakpoint behavior?

**Auto-Selected Decisions:**
- **Default Ratio:** 50:50 split between editor and preview
  - *Rationale:* Balanced starting point for user customization
- **Minimum Widths:** 300px per pane, 1200px total minimum
  - *Rationale:* Ensures usability on smaller screens while maintaining dual-pane utility
  - *Alignment with:* REQUIREMENTS.md EDIT-03: "默认双栏布局"
- **Drag Handle:** Visual feedback with cursor change and hover highlight
  - *Rationale:* Standard UI pattern for resizable panes
- **Responsive Behavior:** Auto-switch to single-pane below 1200px
  - *Rationale:* Graceful degradation for smaller windows
  - *Alignment with:* PROJECT.md constraint: "窗口<1200px自动切换单栏布局"

### 3. Real-time Preview Sync (EDIT-06)
**Questions Covered:**
- Debounce strategy for rendering (<200ms requirement)?
- Scroll synchronization between editor and preview?
- Error handling for malformed Markdown?

**Auto-Selected Decisions:**
- **Debounce Timing:** 150ms delay
  - *Rationale:* Balances responsiveness with performance, stays within <200ms requirement
  - *Alignment with:* REQUIREMENTS.md EDIT-06: "延迟<200ms"
- **Scroll Sync:** Enable bidirectional scroll synchronization
  - *Rationale:* Improves user experience for longer documents
- **Error Handling:** Graceful inline error messages without breaking application
  - *Rationale:* Maintains application stability while providing user feedback

### 4. Line-level Interactions (EDIT-08, EDIT-09, EDIT-10)
**Questions Covered:**
- Visual design for fold/expand icons?
- Drag-and-drop implementation approach?
- Context menu items and shortcuts?
- Performance considerations for large documents?

**Auto-Selected Decisions:**
- **Visual Design:** ▶/▼ icons in gutter for headings and lists
  - *Rationale:* Standard markdown editor convention
- **Drag Implementation:** HTML5 drag API with visual feedback
  - *Rationale:* Native browser API, Vue 3 compatible
- **Context Menu:** Move Up/Down, AI Rewrite (placeholder), Delete
  - *Rationale:* Essential operations, AI placeholder for Phase 3
  - *Alignment with:* REQUIREMENTS.md EDIT-10: "上移/下移/AI重写/删除"
- **Performance:** Optimize for documents up to 10,000 lines using virtual scrolling
  - *Rationale:* Covers typical resume length with headroom

### 5. A4 Page Boundaries (EDIT-05)
**Questions Covered:**
- How to display A4 page boundaries in preview?
- Page break indication approach?

**Auto-Selected Decisions:**
- **Boundary Lines:** Semi-transparent lines showing 210mm A4 width
  - *Rationale:* Visual guidance without obstructing content
  - *Alignment with:* REQUIREMENTS.md EDIT-05: "固定显示A4标准纸张分页线"
- **Page Breaks:** Dashed lines with page numbers in margin
  - *Rationale:* Clear pagination indication for PDF export preview

## Prior Decisions Carried Forward

From Phase 1 analysis:
- **Markdown-it as unified rendering engine** — Ensures preview/export consistency
- **JSON Schema as single source of truth** — Influences editor data binding
- **Local-only data storage** — Confirms no CDN dependencies for Monaco
- **Vue 3 + TypeScript stack** — Confirms Composition API usage

## Deferred Items

The following items were explicitly identified as out-of-scope for this phase:

### Phase 3 (AI Features)
- AI Rewrite functionality (placeholder only in context menu)
- Diff view for AI modifications
- AI-assisted content completion

### Phase 5 (Templates & Export)
- Custom CSS styling
- Template switching
- PDF export implementation

### Phase 7 (UI Polish)
- Dark/light theme toggle
- Advanced responsive design
- Keyboard shortcut customization

### Future Versions
- Performance optimization for >10,000 line documents
- Collaborative editing
- Advanced markdown extensions

## Technical Constraints Applied

All decisions validated against:
- **Performance:** <200ms rendering delay, <200MB memory, <2s cold start
- **Compatibility:** Windows, macOS, Linux support
- **Privacy:** Local-only operation, no external dependencies
- **Core Value:** "编辑器预览与PDF导出100%排版一致"

## Claude's Discretion Areas

The following areas were left to Claude's discretion for planning/implementation:
- Exact Monaco Editor version selection
- Specific animation timings for drag-and-drop
- Exact color schemes for visual elements
- Performance optimization strategies beyond 10,000 lines

---
*Discussion Log: Phase 02-core-editor*
*Generated: 2026-04-05*
*Mode: --auto automated decision selection*