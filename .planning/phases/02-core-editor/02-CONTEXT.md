# Phase 2: 核心编辑器 - Context

**Gathered:** 2026-04-05
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 2 delivers the core Markdown editing experience with Monaco Editor integration, real-time preview in dual-pane layout, and basic line-level interactions. This phase focuses on establishing the fundamental editing interface that ensures WYSIWYG consistency between editor preview and PDF export.

**Out of scope:** AI features (Phase 3), template management (Phase 5), advanced styling (Phase 5)
</domain>

<decisions>
## Implementation Decisions

### Monaco Editor Integration (EDIT-01, EDIT-02)
- **D-01:** Use Monaco Editor via npm package (@monaco-editor/loader) bundled with webpack, not CDN — ensures offline functionality and faster startup
- **D-02:** Default to VS Code Dark theme, add light theme variant for Phase 7 dark mode implementation
- **D-03:** Font size 14px, line-height 1.6 for optimal readability matching VS Code defaults
- **D-04:** Enable Chinese language support for Monaco Editor UI and spell checking
- **D-05:** Enable all VS Code standard editing features: undo/redo, multi-cursor, block selection, find/replace

### Dual-Pane Layout (EDIT-03)
- **D-06:** Default 50:50 width ratio between editor and preview panes
- **D-07:** Minimum width constraint: 300px for each pane, with total minimum window width of 1200px
- **D-08:** Implement drag handle with visual feedback (cursor change, hover highlight) in the divider between panes
- **D-09:** Below 1200px window width, automatically switch to single-pane mode with toggle button between editor/preview views

### Real-time Preview Sync (EDIT-06)
- **D-10:** Implement 150ms debounce for rendering updates (within <200ms requirement)
- **D-11:** Add scroll synchronization between editor and preview for better user experience
- **D-12:** Graceful error handling for malformed Markdown — show inline error message in preview pane without breaking the application

### Line-level Interactions (EDIT-08, EDIT-09, EDIT-10)
- **D-13:** Display fold/expand icons (▶/▼) in the gutter area for headings (# ## ###) and list items (- * +)
- **D-14:** Implement drag-and-drop using HTML5 drag API on the gutter icons with visual feedback during drag
- **D-15:** Context menu on gutter icon click with options: Move Up/Down, AI Rewrite (Phase 3 placeholder), Delete
- **D-16:** Optimize performance for documents up to 10,000 lines using virtual scrolling techniques

### A4 Page Boundaries (EDIT-05)
- **D-17:** Display semi-transparent A4 page boundary lines in preview pane (210mm width approximation)
- **D-18:** Page breaks indicated by dashed lines with page numbers in margin

### Claude's Discretion
- Exact Monaco Editor version selection (use latest stable compatible with Vue 3)
- Specific animation timing for drag-and-drop feedback
- Exact color scheme for page boundary lines (ensure visibility without distraction)
- Performance optimization strategies for large documents beyond 10,000 lines

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Specifications
- `.planning/REQUIREMENTS.md` — Full requirements specification including EDIT-01 through EDIT-11
- `.planning/PROJECT.md` — Project vision, technical stack constraints, and key decisions
- `.planning/ROADMAP.md` — Phase 2 success criteria and dependency relationships

### Phase 1 Context
- `.planning/phases/01-01/01-03-SUMMARY.md` — Bridge layer implementation, TypeScript types, markdown-it initialization
- `.planning/phase-1/CONTEXT.md` — Phase 1 decisions affecting editor implementation

### Code References
- `frontend/src/utils/markdown.ts` — Markdown-it rendering engine instance (unified rendering)
- `frontend/src/views/EditorView.vue` — Current editor placeholder to be replaced
- `frontend/src/types/resume.ts` — Resume type definitions with markdownContent field
- `app.go` — Bridge CRUD methods for loading/saving resume content

### Technical Documentation
- Monaco Editor API documentation: https://microsoft.github.io/monaco-editor/api/index.html
- @monaco-editor/loader documentation: https://github.com/suren-atoyan/monaco-react#usage
- Markdown-it documentation: https://github.com/markdown-it/markdown-it
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `frontend/src/utils/markdown.ts` — Markdown-it instance (md) and renderMarkdown() function for unified rendering
- `frontend/src/types/resume.ts` — Resume interface with markdownContent field for type-safe content handling
- `app.go` Bridge methods — CreateResume, GetResume, UpdateResume for content persistence

### Established Patterns
- Wails Bridge pattern: Go methods exposed to frontend via wailsjs bindings
- Service layer pattern: ResumeService interface for testable business logic
- Context.Background() usage in Bridge methods for Wails context lifecycle

### Integration Points
- EditorView.vue: Replace placeholder with Monaco Editor + preview pane implementation
- App.go: May need additional Bridge methods for auto-save functionality (Phase 4)
- markdown.ts: Extend with additional plugins for line-level interactions (folding, etc.)

### Technical Constraints from Codebase
- Vue 3 Composition API with TypeScript
- Vite build system (requires Monaco Editor webpack configuration)
- Wails runtime for Go-JavaScript interop
</code_context>

<specifics>
## Specific Ideas

**User Experience Priorities:**
- "编辑器预览与PDF导出100%排版一致" — Editor preview must match PDF export exactly (Core Value)
- Target users: VS Code users — maintain familiar editing experience
- Chinese language support for spell checking and UI localization

**Performance Requirements:**
- <200ms rendering delay (EDIT-06)
- <200MB total memory footprint (PROJECT.md constraint)
- Cold start <2s (PROJECT.md constraint)

**Layout Constraints:**
- A4 page boundaries must be visible in preview (EDIT-05)
- Minimum window width: 1200px for dual-pane, auto-switch to single-pane below
- Offline functionality — no CDN dependencies for Monaco Editor
</specifics>

<deferred>
## Deferred Ideas

**Deferred to Phase 3 (AI Core Features):**
- AI Rewrite functionality in context menu (EDIT-10) — placeholder only
- Diff view for AI modifications (EDIT-11)
- AI-assisted content completion and suggestions

**Deferred to Phase 5 (Templates & Export):**
- Custom CSS styling for preview pane
- Template switching mechanism
- PDF export functionality (preview rendering preparation only)

**Deferred to Phase 7 (UI Polish):**
- Dark/light theme toggle implementation
- Advanced responsive design beyond 1200px breakpoint
- Keyboard shortcut customization

**Deferred to Future Versions:**
- Advanced performance optimization for documents >10,000 lines
- Collaborative editing features
- Advanced markdown extensions (diagrams, math formulas)
</deferred>

---
*Phase: 02-core-editor*