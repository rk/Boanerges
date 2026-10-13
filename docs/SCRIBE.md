# Scribe column

Personal chapter draft: one **contenteditable span per Bible verse**, optional **paragraph breaks** between verses, debounced save to disk, **preview** as reader prose, and **print** integration.

Sources: legacy tree at git `6bb99d8^` (`legacy/resources/js/components/scribe/*`, `legacy/resources/js/lib/scribe*.ts`, `legacy/app/Services/Scribe/ScribeChapterStore.php`). MyGO rewrite uses `internal/studydata` + `StudyData.ScribeGet` / `ScribePut`.

---

## User-visible behavior

### Layout

- Scribe is a study **column type** (`study.scribe` menu / column picker).
- Header: column title, optional **book chapter** label when chapter loaded, **save indicator** (spinner while saving, checkmark briefly after save), **actions menu** (⋮):
  - **Set Line Breaks from Primary** — copy paragraph-start flags from the primary Bible column’s chapter verses (verse 1 unchanged).
  - **Reset Line Breaks** — remove user paragraph overrides for verses 2+ (keep verse text).
  - **Preview** — modal showing draft as formatted paragraphs (sup verse numbers, `FormattedVerseText`).
- Body: scrollable **inline verse editor** — verses flow like prose, not a single textarea.
- Each verse shows:
  - Small **primary-colored superscript** verse number (CSS `::before` on empty/hydrated span).
  - **Pilcrow** control (disabled for verse 1) to toggle “new paragraph before this verse”.
  - **Contenteditable** text (`plaintext-only`), multiline within the verse (`white-space: pre-wrap`).
- Empty verses show a subtle **placeholder block** (min width) so the line is tappable.
- **Paragraph break before verse N** (N > 1): block spacing via `[data-paragraph-start='true']:before`.
- Uses **readability** settings (font size, line height, font family, justify) like Bible columns.

### Interactions

| Action | Behavior |
|--------|----------|
| Type in verse | Updates in-memory draft; **debounced save** (~500ms) to backend |
| Paste | Plain text only; normalized `\r\n` / `\r` → `\n`; no HTML |
| Tab / Shift+Tab | Move focus to next/previous verse span (within `.scribe-document`) |
| Pilcrow click | Toggle `paragraphStart` override for that verse (not on verse 1) |
| Navigate book/chapter | Reload draft for new location; re-hydrate spans from server |
| Actions menu | Primary paragraph sync, reset breaks, preview modal |

### Keyboard

- **Tab** / **Shift+Tab** between verse fields (custom handler; default tab trap avoided).
- Standard typing/editing inside contenteditable; Enter inserts newline **inside** the verse (intra-verse “paragraphs” appear in **Preview** when text has blank-line gaps).

### Print (`PrintOptionsModal` + `internal/studyprint`)

- **Reading layout** with user work: Scribe column prints **scribe-content** — merged primary chapter structure with draft text and paragraph flags (`scribeVerses` in Go mirrors `toPreviewVerses`).
- **Handwriting / lined layout**: Scribe column prints empty **lined area** (no draft text).
- Notes + Scribe copy in modal matches legacy strings.

---

## Data model

### On disk

Path (user data root):

```text
scribe/{bookId}/{chapter}.json
```

- `bookId`: lowercase alphanumeric, 1–10 chars (e.g. `gen`).
- `chapter`: integer filename `{chapter}.json`.

Legacy Laravel used `storage/app/private/scribe/`; MyGO uses the same relative path under app user data (`MIGRATION.md`: copy from Laravel private storage to user data `scribe/`).

### JSON shape

Array of verse objects (sparse — only verses with text and/or explicit paragraph flag):

```json
[
  { "verse": 1, "text": "My translation of v1", "paragraphStart": true },
  { "verse": 3, "text": "Line one\n\nLine two", "paragraphStart": true }
]
```

| Field | Type | Meaning |
|-------|------|---------|
| `verse` | number | Verse number (1-based) |
| `text` | string | Draft plain text; may contain `\n` |
| `paragraphStart` | boolean, optional | **Override** only: if present, forces paragraph break before this verse in preview/print. Omitted = default rules (verse 1 starts a paragraph; others continue unless primary sync or user toggle set override) |

### In-memory draft (`ScribeDraftEntry`)

Per verse number:

```ts
{
  text: string;
  paragraphStartOverride?: boolean; // maps to JSON paragraphStart when set
}
```

**Serialize** (save): walk **all chapter verse numbers** from Bible API; emit entry only if `text.trim()` non-empty **or** `paragraphStartOverride` defined.

**Effective paragraph start** for UI:

- If `paragraphStartOverride` defined → use it.
- Else verse 1 → `true`, else → `false`.

---

## API / backend contract

### Legacy (Laravel)

- `GET /scribe/{book}/{chapter}` → `{ verses: ScribeVerse[] }`
- `PUT /scribe/{book}/{chapter}` body `{ verses: [...] }` → `{ verses }`
- Validated: `verses.*.verse` 1–200, `text` string, optional `paragraphStart` boolean.

### MyGO

Go: `internal/studydata.Service`

- `ScribeGet(bookID, chapter) ([]ScribeVerse, error)` — missing file → `[]`.
- `ScribePut(bookID, chapter, verses) error` — writes compact JSON (not pretty-print; legacy used pretty-print — compatible for read).

Bound to frontend via `StudyData.scribeGet` / `scribePut` in `src/mygo.ts`.

E2E mock: `sessionStorage` key `e2e.scribe.{bookId}:{chapter}`.

---

## Acceptance checklist (legacy parity)

- [ ] One editable field per verse in current chapter (from primary translation chapter fetch).
- [ ] Verse numbers visible; empty verses still show target area.
- [ ] Pilcrow toggles paragraph before verse (not v1); visual break matches preview.
- [ ] Tab moves between verses.
- [ ] Paste strips formatting; keeps plain text newlines.
- [ ] Debounced autosave; saving/saved indicator in header.
- [ ] Draft survives navigation away and app reload (same book/chapter).
- [ ] “Set Line Breaks from Primary” / “Reset Line Breaks” / “Preview” in header menu.
- [ ] Preview modal: paragraph grouping + intra-verse `\n\n` splits.
- [ ] Print reading mode includes scribe draft when column is Scribe.
- [ ] On-disk JSON under `scribe/{bookId}/{chapter}.json` matches legacy field names.

---

## Remaining gaps vs legacy

| Area | Legacy | MyGO rewrite (current) |
|------|--------|-------------------------|
| localStorage migration | One-time uplift `boanerges.scribe.*` → API | **Not ported** (users who never opened Scribe post-Laravel may still have keys) |
| JSON pretty-print | Laravel `JSON_PRETTY_PRINT` | Compact JSON (same schema; readable either way) |
| Global save toast | `SaveStatus.svelte` in app layout | Save icons only in Scribe column header |
| Chapter heading component | `ChapterHeading.svelte` | Compact span in column header |

---

## Key files (rewrite target)

| Role | Path |
|------|------|
| Column shell | `src/components/columns/ScribeColumn.svelte` |
| Editor | `src/components/scribe/ScribeEditor.svelte`, `ScribeVerseSpan.svelte`, `ScribePreviewModal.svelte`, `ScribeParagraphText.svelte` |
| Logic | `src/lib/scribe.svelte.ts`, `scribeParagraphBreaks.ts`, `scribeDocument.ts`, `paragraphs.ts` |
| Styles | `src/style.css` (`.scribe-*`) |
| Backend | `internal/studydata/service.go` |
| Print | `internal/studyprint/columns.go`, `render.go` |
| Tests | `src/lib/scribeParagraphBreaks.test.ts`, `src/lib/scribeDocument.test.ts`, `e2e/study-journey.spec.ts` |
