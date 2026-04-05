---
phase: 03-ai-core
plan: 03-05
subsystem: ai
tags: [vue3, typescript, markdown, claude-api, sse, wails]

# Dependency graph
requires:
  - phase: 03-01
    provides: AI SSE infrastructure, sendMessage/sendMessageSync, useAIStream composable
provides:
  - ResumeParserModal for importing and parsing old resumes via AI
  - JobTargetChip inline-editable component for target job position
  - JSON to Markdown conversion utility (jsonToMarkdown, extractBasicInfo)
  - ParsedResume validation function (validateParsedResume)
  - jobTarget field on Resume type and model
affects: [03-02, 03-04] (AI context injection uses jobTarget + includeFullContext)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - AI-powered resume parsing via sendMessageSync with structured JSON output
    - JSON extraction from AI text (code block or bare JSON detection)
    - Resume structured data to Markdown conversion
    - Inline-editable chip component pattern (click to edit, Enter/Escape confirm/cancel)

key-files:
  created:
    - frontend/src/components/ResumeParserModal.vue - Main parser modal with paste area, job target input, full context toggle, preview, import
    - frontend/src/components/JobTargetChip.vue - Inline editable target job chip
    - frontend/src/utils/resume.ts - jsonToMarkdown, extractBasicInfo utilities
  modified:
    - frontend/src/types/resume.ts - Added jobTarget, ParsedResume interface, validateParsedResume()
    - frontend/src/views/EditorView.vue - Integrated toolbar, ResumeParserModal, JobTargetChip, auto-save
    - internal/model/resume.go - Added JobTarget field to Resume struct

key-decisions:
  - "Used sendMessageSync (non-streaming) for parsing since we need the full JSON response before preview"
  - "JSON extraction handles both code block format (```json ...```) and bare JSON in AI response"
  - "Resume model JobTarget field enables job target persistence without schema migration"
  - "Inline ID generation (Date.now + random string) avoids uuid package dependency"

patterns-established:
  - "Modal pattern with Teleport to body, overlay click-to-close, reset on hide"
  - "AI parsing with schema validation, graceful degradation (show preview even on partial validation)"
  - "Editor toolbar with dark theme buttons matching VS Code aesthetic"

requirements-completed: [AIAI-05, AIAI-06, AIAI-09]

# Metrics
duration: 5min
completed: 2026-04-05
---

# Plan 03-05: 一键解析旧简历 & 上下文管理

**AI-powered resume parsing: users paste plain text or Markdown old resumes, AI extracts structured data and converts to Markdown; JobTargetChip inline editor with full context toggle**

## Performance

- **Duration:** 5 min
- **Started:** 2026-04-05T20:26:09+08:00
- **Completed:** 2026-04-05T20:30:40+08:00
- **Tasks:** 6 (5 task commits + build)
- **Files modified:** 7 files (3 created frontend, 2 modified frontend, 1 modified backend model, 1 modified frontend view)

## Accomplishments

- ResumeParserModal: paste old resumes, AI parses to structured JSON, preview before import, converts to Markdown
- JobTargetChip: inline-editable chip in editor toolbar showing "目标岗位: {target}", persists to resume record
- Full context toggle: includeFullContext flag passed to all AI operations (OFF by default to save tokens)
- ParsedResume type + validateParsedResume() function for schema validation of AI output
- jsonToMarkdown() utility converts structured resume data to app's Markdown format
- Backend Resume model extended with JobTarget field for persistence

## Task Commits

1. **Task 1: Resume types extension** - `ea3a0f0` (feat): jobTarget field, ParsedResume/ExperienceItem interfaces, validateParsedResume()
2. **Task 2: JSON to Markdown utility** - `1dbe859` (feat): jsonToMarkdown() and extractBasicInfo() functions
3. **Task 3: ResumeParserModal component** - `c7bc873` (feat): full modal with paste area, AI parsing, JSON extraction, preview, error handling
4. **Task 4: JobTargetChip component** - `a84d2b9` (feat): inline-editable chip with click-to-edit, Enter/Escape keys
5. **Task 5: EditorView integration** - `9b32196` (feat): toolbar with import button, ResumeParserModal integration, JobTargetChip, auto-save on change

**Build verification:** `9b32196` (frontend npm run build passed)

## Files Created/Modified

### Frontend (Vue 3 + TypeScript)
- `frontend/src/components/ResumeParserModal.vue` (692 lines) - Main parser modal with paste area, job target, full context toggle, AI parsing with JSON extraction and validation, preview display, import confirmation
- `frontend/src/components/JobTargetChip.vue` (179 lines) - Inline editable chip with click-to-edit, Enter/Escape/cancel support, dark theme
- `frontend/src/utils/resume.ts` (147 lines) - jsonToMarkdown() converting ParsedResume to formatted Markdown, extractBasicInfo() helper
- `frontend/src/types/resume.ts` - Added jobTarget to Resume, ParsedResume interface, ExperienceItem interface, validateParsedResume() function
- `frontend/src/views/EditorView.vue` - Added toolbar with import button + JobTargetChip, ResumeParserModal integration, auto-save with debounce, jobTarget sync

### Backend (Go)
- `internal/model/resume.go` - Added JobTarget string field to Resume struct

## Decisions Made

- **sendMessageSync over streaming for parsing**: Parsing needs the complete JSON before showing preview, so non-streaming is simpler
- **JSON extraction regex**: Handles both ` ```json ... ``` ` code blocks and bare `{...}` JSON in plain text
- **Schema validation with graceful degradation**: Even if validation fails, show parsed data so user can review and manually fix
- **Inline ID generation**: `Date.now() + random string` instead of uuid package dependency
- **Resume model JobTarget**: Added field directly to struct (not migration) since JSON is stored as text

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - all tasks completed without issues. Build passed on first attempt.

## Self-Check: PASSED

- [x] All created files exist (ResumeParserModal.vue, JobTargetChip.vue, resume.ts, resume.ts utility)
- [x] All 5 task commits verified in git history
- [x] Build passes (npm run build, vue-tsc noEmit)
- [x] SUMMARY.md created with substantive content
- [x] Requirements [AIAI-05, AIAI-06, AIAI-09] marked complete

---

## Next Phase Readiness

- ResumeParserModal ready for use by end users with AI API key configured
- JobTargetChip component available for reuse in other plans (03-04 AI Chat sidebar)
- jsonToMarkdown utility reusable for other AI-to-resume conversion flows
- includeFullContext flag already wired into sendMessage/sendMessageSync service calls

---
*Plan: 03-05*
*Completed: 2026-04-05*
