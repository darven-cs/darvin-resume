# Phase 3: AI 核心能力 - Context

**Gathered:** 2026-04-05
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 3 delivers AI capabilities for resume editing: Claude API integration with streaming responses, selection-based AI operations (polish/translate/summarize/rewrite), AI chat sidebar for conversational assistance, one-click resume parsing from plain text, and comprehensive error handling. This phase makes the resume editor intelligent.

**Out of scope:** Template switching (Phase 5), dark/light themes (Phase 7), ATS optimization (v2)
</domain>

<decisions>
## Implementation Decisions

### Architecture: SSE Streaming (AIAI-07)
- **D-01:** Go backend makes HTTP requests to Claude API, streams SSE responses back to frontend via Wails EventsEmit/EventsOn
- **D-02:** Frontend uses `eventsource` npm package for SSE client connection to Wails event bus
- **D-03:** Streaming chunks processed character-by-character for typewriter effect, with 16ms debounce for DOM updates
- **D-04:** Wails Events namespaced as `ai:stream:{operationId}` for multi-operation isolation

### API Location (AIAI-01, AIAI-03)
- **D-05:** All AI API calls go through Go backend via Wails Bridge methods — API Key never exposed to frontend
- **D-06:** AI configuration stored in existing `settings` SQLite table under `ai.*` keys (apiKey, baseUrl, defaultModel, maxTokens, timeoutSeconds)
- **D-07:** Config CRUD via Bridge methods: GetAIConfig, SaveAIConfig, ValidateAPIKey

### Claude API Integration (AIAI-01, AIAI-04)
- **D-08:** Use Claude Messages API (POST /v1/messages) with `stream: true` parameter
- **D-09:** Request body built in Go backend using standard `net/http` — no external HTTP libraries needed
- **D-10:** All prompts structured to output JSON matching Resume JSON Schema for typed responses
- **D-11:** Support custom BaseURL for API proxy compatibility (anthropic.com is default)

### Selection-based AI Operations (EDIT-07)
- **D-12:** Floating toolbar appears above/below selected text when selection is non-empty
- **D-13:** Toolbar position calculated via Monaco's `getScrolledVisiblePosition()` — auto-flips to stay in viewport
- **D-14:** Operations: 润色(polish), 翻译(translate), 缩写(summarize), 重写(rewrite)
- **D-15:** Each operation sends current selection + job target + optional full context to backend
- **D-16:** Loading state shown in toolbar during streaming; operation ID stored for abort capability

### Diff Comparison (EDIT-11)
- **D-17:** After AI modification, show diff inline below the edited paragraph (not modal, not replace preview)
- **D-18:** Use `diff` npm package to compute line-level diffs
- **D-19:** Accept replaces original with AI version (via editor.executeEdits); Reject closes diff view, keeps original
- **D-20:** Diff colors: deletions in red (#ffcccc bg), insertions in green (#ccffcc bg)

### AI Chat Sidebar (AIAI-08)
- **D-21:** Slide-in panel on right side of editor, 400px wide, independent of SplitPane layout
- **D-22:** Chat messages stored per-resume in `ai_messages` table, linked by resume_id
- **D-23:** Support "引用选中文本" — inserts selected text as user message with quote formatting
- **D-24:** Assistant output has "插入编辑区" button — calls editor.executeEdits at cursor position
- **D-25:** Job target prepended as system context to each conversation

### Job Target Context (AIAI-05)
- **D-26:** Job target displayed as editable chip at top of editor, persisted with resume
- **D-27:** Default empty string; if empty, AI prompts omit job-specific optimizations
- **D-28:** Job target included in all AI operation requests as system context

### Full Context Toggle (AIAI-06)
- **D-29:** Toggle switch in AI config panel: "包含全文参考" — default OFF for selection operations
- **D-30:** When ON: append full resume Markdown content to user message as reference
- **D-31:** Global setting persisted per-user (not per-resume)

### One-click Resume Parsing (AIAI-09)
- **D-32:** Triggered from HomeView paste area or EditorView import button
- **D-33:** User pastes plain text/Markdown into modal; AI returns structured JSON matching Resume Schema
- **D-34:** JSON auto-populates resume data model; user confirms before generating Markdown

### Error Handling (AIAI-10, -11, -12, -13)
- **D-35:** Network errors: toast notification + "重试" button, original text preserved
- **D-36:** Token limit: split request into chunks + warning modal suggesting content reduction
- **D-37:** Format errors: auto-retry once; on second failure, return raw response with warning
- **D-38:** Abort: stop streaming via operationId cancel, preserve accumulated text
- **D-39:** All errors logged to frontend console with operationId for debugging

### New npm Dependencies
- **D-40:** `eventsource` — SSE client for frontend
- **D-41:** `diff` — diff computation for EDIT-11

### Claude's Discretion
- Exact streaming debounce interval (16ms default, adjustable)
- Sidebar animation timing and easing
- Diff view styling specifics
- Operation abort implementation details
- Chat message UI design details
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Specifications
- `.planning/REQUIREMENTS.md` — AIAI-01 through AIAI-13, EDIT-07, EDIT-11
- `.planning/PROJECT.md` — Architecture decisions, tech stack constraints
- `.planning/ROADMAP.md` — Phase 3 success criteria and dependencies

### Phase 2 Context
- `.planning/phases/02-core-editor/02-SUMMARY.md` — Completed editor components
- `.planning/phases/02-core-editor/02-CONTEXT.md` — Editor implementation decisions

### Code References
- `frontend/src/components/MonacoEditor.vue` — Selection APIs, executeEdits, cursor events
- `frontend/src/views/EditorView.vue` — EditorView layout, current state management
- `frontend/src/types/resume.ts` — Resume type definitions (AI output must match)
- `frontend/src/utils/markdown.ts` — Markdown rendering for preview
- `frontend/src/composables/useDebounce.ts` — Debounce utilities (reusable for streaming)
- `app.go` — Bridge methods, Wails runtime EventsEmit/EventsOn usage
- `internal/database/db.go` — SQLite settings table access pattern

### Technical Documentation
- Claude Messages API: https://docs.anthropic.com/en/api/messages
- SSE Protocol: https://html.spec.whatwg.org/multipage/server-sent-events.html
- Monaco Selection API: https://microsoft.github.io/monaco-editor/api/modules/monaco.editor.html
- Wails Events: https://wails.io/docs/reference/runtime/events
</canonical_refs>

<file_plan>
## Phase 3 File Plan

### Backend (Go)
```
backend/
├── internal/
│   ├── ai/
│   │   ├── client.go        # Claude API HTTP client, streaming, config
│   │   └── prompt.go        # Prompt templates for each operation
│   └── model/
│       └── config.go        # AIConfig struct
└── app.go                   # Add AI Bridge methods
```

### Frontend (Vue/TS)
```
frontend/src/
├── components/
│   ├── AIFloatingToolbar.vue    # Selection-based floating toolbar
│   ├── AIDiffView.vue           # Inline diff comparison
│   ├── AIChatSidebar.vue        # Conversational AI panel
│   └── AIConfigModal.vue        # API key/model configuration
├── composables/
│   ├── useAIStream.ts           # SSE event handling, streaming state
│   ├── useAISelection.ts        # Selection detection, toolbar logic
│   └── useAIConfig.ts           # Config state and validation
├── types/
│   └── ai.ts                    # AI types, operation types, error codes
└── views/
    └── EditorView.vue            # Integrate AI components
```
</file_plan>

<deferred>
## Deferred Ideas

**Deferred to Phase 5 (Templates & Export):**
- AI-generated template recommendations
- AI-assisted style optimization

**Deferred to Phase 6 (Security):**
- API Key AES-256-GCM encryption (AIAI-02) — store plaintext this phase, encrypt Phase 6
- Secure key derivation from device ID

**Deferred to Phase 7 (UI Polish):**
- Dark/light mode theming for AI components
- Keyboard shortcuts (Ctrl+R/T/D) for AI operations
- Loading animation design

**Deferred to v2:**
- ATS optimization suggestions
- Content duplication detection
- Job matching analysis
</deferred>
