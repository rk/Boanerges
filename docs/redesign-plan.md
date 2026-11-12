# Study UI redesign plan

Boanerges Study should feel **designed for long reading sessions**: attractive without spectacle, with a quiet application frame and a warm, legible scripture canvas. This document captures the agreed direction, theme system (two aesthetic poles × light/dark), layout and component patterns, and a phased implementation order.

**Scope:** web Study shell (Svelte + daisyUI). Native MyGo UI is out of scope here; reuse the same tokens and patterns when columns migrate.

---

## Goals

| Goal | Success signal |
|------|----------------|
| Long-session comfort | Users can read for 30+ minutes without eye strain or “admin panel” fatigue |
| Clear hierarchy | Primary Bible reads as the main canvas; tools read as panels |
| Restrained color | One accent role for focus/CTA; neutrals for chrome and static verse chrome |
| Accessibility | Readable defaults (size, line-height, contrast); comfortable tap targets on key chrome |
| Incremental delivery | Each phase ships visible value without a big-bang rewrite |

---

## Current state (baseline)

Strengths to preserve:

- Reader font stack and per-user **readability** settings (size, line-height, font family, justify).
- Multi-column layout and column type picker.
- Native flyouts for book/chapter and verse actions (clip-safe).
- Recent tool UX improvements (verse list combobox, notes save status, cross-ref navigation).

Pain points addressed in this redesign:

- Stock daisyUI **light/dark** only; minimal brand identity.
- **Equal-width columns** and full-height rules.
- **Repeated column headers** and overly small controls.
- **Primary color on every verse number**.

Design mantra: **quiet frame, warm page**.

---

## Design principles

1. **Two surfaces** — **Frame** (toolbar, panel headers, workspace background): `base-200`. **Page** (Bible / Scribe body inside cards): `base-100`.
2. **One accent job** — Reserve `primary` for **current focus** (nav flash, primary CTA). Verse numbers and ranks use **neutral** base-content opacity.
3. **Measure and rhythm** — Cap line length (~68ch), defaults **17px / 1.6** line-height for new installs.
4. **Hierarchy over equality** — Primary Bible `flex-[1.35]`; tools are card panels.
5. **Comfort metrics** — Chrome uses `btn-sm` / `input-sm`; chapter prev/next use taller join buttons (`min-h-11`).

---

## Theme system: two poles × light/dark

### Stored settings (two fields)

| Field | JSON | Values | Role |
|-------|------|--------|------|
| **Appearance** | `theme` | `auto`, `light`, `dark` | Light/dark mode (auto = `prefers-color-scheme`) |
| **Pole** | `themePole` | `monastic`, `modern` | Aesthetic family (warm vs cool chrome) |

Runtime composes daisyUI name: `study-{pole}-{light|dark}` via [`src/lib/studyTheme.ts`](../src/lib/studyTheme.ts) and [`applyReadabilityTheme`](../src/lib/readerTheme.ts). `data-theme` is **always** set (auto listens to `matchMedia` and updates).

### daisyUI theme names

| Theme name | Pole | Mode |
|------------|------|------|
| `study-monastic-light` | Monastic | Light |
| `study-monastic-dark` | Monastic | Dark |
| `study-modern-light` | Modern | Light |
| `study-modern-dark` | Modern | Dark |

Registered in [`src/style.css`](../src/style.css). Default plugin default: **study-monastic-light**; prefers-dark fallback: **study-monastic-dark**.

### Defaults (new install)

- `theme`: `auto`
- `themePole`: `monastic`
- `fontSize`: 17, `lineHeight`: 1.6 (existing saved values are not rewritten on upgrade)

### Legacy migration (on read)

Old single `theme` string values:

- `auto` / `light` / `dark` — unchanged
- Legacy light daisy names (`sepia`, `nord`, `caramellatte`, …) → `light`
- Legacy dark daisy names (`dracula`, `night`, …) → `dark`
- Missing `themePole` → `monastic`

Writes accept only `auto|light|dark` and `monastic|modern`. See [`internal/settings/appearance.go`](../internal/settings/appearance.go).

### Reader-only tints

No fifth global theme for sepia/night on the article. Shell uses study themes only. Print/PDF stays paper-light ([`internal/studyprint/study-print.css`](../internal/studyprint/study-print.css)).

### Settings UI

**Readability** drawer: pole (Monastic / Modern) + appearance (System / Light / Dark).

---

## Layout

- **Workspace:** `bg-base-200`, `gap-2`, `p-2` ([`ColumnLayout.svelte`](../src/components/layout/ColumnLayout.svelte)).
- **Columns:** `rounded-box bg-base-100 shadow-sm`; primary wider (`flex-[1.35]`).
- **Toolbar:** slim wayfinding row; no in-study wordmark ([`StudyToolbar.svelte`](../src/components/toolbar/StudyToolbar.svelte)).
- **Settings:** end **drawer** with overlay ([`Study.svelte`](../src/pages/Study.svelte)).

---

## Components

- **`ColumnHeader`** — icon + title for tools; primary Bible shows translation `select-sm` only.
- **`StudyListRow`** — search (and similar) hit rows.
- **Verse list** — combobox stays compact; **Add verse** / **Delete** use `btn-sm btn-outline`.

---

## Reader

- **Measure:** `.reader-measure` (`max-width: 68ch`) on Bible and Scribe content.
- **Verse numbers:** `.verse-num` uses muted `base-content` (not `primary`).
- **navFlash:** `bg-primary/10` on verse span.

---

## Phased implementation

| Phase | Status | Notes |
|-------|--------|-------|
| **A — Tokens** | Done | Four themes, persistence, settings, tests |
| **B — Layout** | Done | Cards, width bias, toolbar, drawer |
| **C — Reader** | Done | Measure, verse numbers, navFlash |
| **D — Components** | Done | ColumnHeader icons, StudyListRow, control sizes |

Keep **E2E** green (`study-journey`, column smoke) when changing chrome.

---

## Decisions (closed)

| Topic | Decision |
|-------|----------|
| Default pole | **Monastic** |
| Resizable columns | **Deferred** (fixed flex ratio) |
| UI webfont | **Deferred** (system sans) |
| Words of Jesus | Keep **`color-error`** in reader markup |

---

## Anti-patterns (do not ship)

- Gradients/glassmorphism on scripture.
- Equal 33/33/33 as the only layout.
- `primary` on static verse numbers or every list rank.
- Fifth global theme for reader sepia/night.
- Toolbar text clutter on narrow widths without collapse.

---

## Related files

| Area | Files |
|------|--------|
| Themes | `src/style.css`, `src/lib/studyTheme.ts`, `src/lib/readerTheme.ts` |
| Persistence | `internal/settings/readability.go`, `internal/settings/appearance.go` |
| Shell | `src/pages/Study.svelte`, `src/App.svelte`, `StudyToolbar.svelte` |
| Layout | `src/components/layout/ColumnLayout.svelte`, `ColumnHeader.svelte` |
| Reader | `VerseSpan.svelte`, `BibleColumn.svelte`, `ScribeEditor.svelte` |
| Readability UI | `ReadabilitySettings.svelte`, `readability.svelte.ts` |

---

## Revision log

| Date | Note |
|------|------|
| 2026-11-12 | Initial plan: four study themes; phases A–D |
| 2026-11-12 | Locked two-field theme model; defaults; implementation completed on `master` |
