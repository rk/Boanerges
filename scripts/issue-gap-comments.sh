#!/usr/bin/env bash
# Posts gap audit comments and reopens MyGO rewrite issues #23–#33.
set -euo pipefail
cd "$(dirname "$0")/.."

post() {
  local num="$1"
  local body="$2"
  gh issue comment "$num" --body "$body"
  gh issue reopen "$num" --comment "Reopened for gap closure sweep."
}

post 23 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Mostly complete.

**Gaps:**
- No **e2e** coverage that readability/study settings persist across app restart and drive UI (legacy Pest/browser parity).
- Confirm all legacy `SettingsController` / store edge cases (validation errors surfaced to UI).

**Reference:** `legacy/tests/` settings-related feature tests; `legacy/app/Services/` settings stores.
EOF
)"

post 24 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Mostly complete (canon, schema, markup unit tests).

**Gaps:**
- **E2e:** none that a installed translation renders formatted verse HTML matching legacy pipeline for a known verse (markup regressions).
- Optional: parity test count vs legacy markup converter tests (not required if e2e covers rendering).

**Reference:** `legacy/tests/Unit/*Markup*`, `legacy/app/Services/Bible/Markup/`.
EOF
)"

post 25 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Core complete (sword reader, prebuilt sqlite, welcome progress, chapter read test).

**Gaps:**
- **E2e:** full first-launch welcome → study transition with progress UI (not just Go unit tests).
- Progress steps are coarse vs legacy `TranslationInstallStep` enum during long catalog imports (bundled copy is fast).

**Reference:** `legacy/resources/js/pages/Welcome.svelte`, `legacy/app/Events/TranslationInstallProgress.php`.
EOF
)"

post 26 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial.

**Gaps:**
- **Accordance importer not ported** (legacy has `AccordanceImporter` / `AccordanceFormatImporter`; defer only if catalog has zero active entries — verify catalog).
- **SWORD catalog install:** needs **feature tests** with fixture zip (legacy KJV/verify patterns), not only ASV bundled guards.
- **USFM:** feature test exists; add **e2e** path: install OEB (or fixture) → read chapter via UI/API.
- **E2e:** Translation manager install/uninstall flow (422 bundled ASV guard visible to user).
- Download/network failure and conflict (409) UX vs legacy `TranslationManagerTest`.

**Reference:** `legacy/tests/Feature/TranslationManagerTest.php`, `legacy/tests/Unit/UsfmImporterTest.php`, `legacy/app/Services/Bible/Import/*`.
EOF
)"

post 27 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial (Go services + unit/feature-style tests).

**Gaps:**
- **E2e:** books list → chapter load → FTS search → cross-ref panel → dictionary lookup (single smoke journey).
- Not all legacy `BibleControllerTest` / `DictionaryTest` error shapes (404 translation, 422 chapter) asserted in **e2e**.
- Cross-ref/Webster imports skip in CI when extras missing — ensure CI always has bundled extras or embed minimal fixtures.

**Reference:** `legacy/tests/Feature/BibleControllerTest.php`, `CrossReferenceTest.php`, `DictionaryTest.php`.
EOF
)"

post 28 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial UI port.

**Gaps vs legacy `Study.svelte` / toolbar:**
- Book/chapter **picker** (legacy `BookChapterPicker`, `ChapterNavRail`) — current simple selects only.
- **TranslationSelect**, column headers, `AppLayout` / themes / fonts (`ReaderThemeParityTest`).
- **E2e:** 1–3 columns, toggle column types from menu, readability settings change reader appearance.

**Reference:** `legacy/resources/js/pages/Study.svelte`, `legacy/resources/js/components/toolbar/*`, `legacy/tests/Unit/ReaderThemeParityTest.php`.
EOF
)"

post 29 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial.

**Gaps:**
- Reader uses simplified verse HTML; missing **VerseLine**, **FormattedVerseText**, paragraph breaks, chapter headings/nav (legacy `reader/*`).
- **Comparison column** aligned verses parity vs legacy `ComparisonColumn.svelte`.
- **Scroll sync** toggle exists but **not wired** between column scroll containers.
- **E2e:** primary + secondary columns show same chapter; comparison shows aligned text.

**Reference:** `legacy/resources/js/components/columns/BibleColumn.svelte`, `ComparisonColumn.svelte`, `legacy/resources/js/components/reader/*`.
EOF
)"

post 30 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial.

**Gaps:**
- **E2e:** notes CRUD per chapter, scribe save/load, verse list create/rename/add verses (on-disk format must stay compatible with legacy paths).
- UI parity vs legacy column components (preview modal, save status, etc.).
- Go tests may not cover all legacy `NotesControllerTest` / `ScribeControllerTest` / `VerseListControllerTest` scenarios.

**Reference:** `legacy/tests/Feature/NotesControllerTest.php`, `ScribeControllerTest.php`, `VerseListControllerTest.php`.
EOF
)"

post 31 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial.

**Gaps:**
- Search **pagination**, ranking, translation filter UX vs legacy `SearchColumn.svelte` + `search.svelte.ts`.
- Cross-refs **rank** ordering and navigation; dictionary **suggest** behavior.
- **E2e:** run search query, open hit → navigates study location; dictionary lookup from selected word.

**Reference:** `legacy/resources/js/components/columns/SearchColumn.svelte`, `CrossReferencesColumn.svelte`, `DictionaryColumn.svelte`.
EOF
)"

post 32 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial.

**Gaps:**
- **File → Manage Translations** menu item not wired to modal (legacy `TranslationManagerModal.svelte`).
- Print: missing legacy **StudyPrintHtmlBuilder** / printable column types / **PrintOptionsModal** (current `Page().Print` / PDF to user dir only).
- Hotkeys parity vs `legacy/app/Support/ApplicationMenuBuilder.php` and `StudyColumnType`.
- **E2e:** menu actions change column count; print dialog invoked without crash.

**Reference:** `legacy/tests/Feature/StudyPrintTest.php`, `legacy/app/Services/Study/StudyPrintService.php`.
EOF
)"

post 33 "$(cat <<'EOF'
## Gap audit (2026-10-12)

**Status:** Partial / premature close.

**Gaps:**
- **`legacy/` must remain on disk as untracked reference** until all other MyGO issues have no gaps; do **not** delete from workspace. (Removed from git in `6bb99d8` — OK if untracked copy kept; add `/legacy/` to `.gitignore`.)
- **Migration doc** (`MIGRATION.md`) not verified against sample NativePHP data dir.
- **`mygo build`** macOS smoke not recorded in CI or issue.
- Cutover (delete legacy from repo) only when **all** issues #23–#32 are complete with **e2e** confidence.

**Reference:** `legacy/README.md`, `legacy/app/Console/Concerns/ConfiguresSqliteDatabase.php`.
EOF
)"

echo "Done: commented and reopened #23–#33"
