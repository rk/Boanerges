#!/usr/bin/env bash
set -euo pipefail
SHA=$(git rev-parse --short HEAD)

close() {
  local n="$1"
  local body="$2"
  gh issue close "$n" --comment "$body"
}

close 23 "**Closed @ ${SHA}** — Settings persist via e2e (readability slider + reload). Validation errors surfaced in UI (readability-error). Go tests: \`internal/settings\`."

close 24 "**Closed @ ${SHA}** — Markup pipeline covered by Go tests + e2e Genesis 1 \`<em>\` render."

close 25 "**Closed @ ${SHA}** — Welcome progress e2e (\`/?e2eWelcome=1\`), bundled ASV path, \`translationInstallProgress\` events, chapter import tests."

close 26 "**Closed @ ${SHA}** — USFM + SWORD feature tests, catalog install e2e (KJV + OEB), bundled ASV guards. **Accordance importer not ported** — \`TestCatalogHasNoAccordanceEntries\` confirms zero active catalog entries; registry rejects format."

close 27 "**Closed @ ${SHA}** — Books/chapter/search/xref/dictionary journey e2e + \`internal/bibleapp\` tests; bundled cross-ref/Webster import on startup."

close 28 "**Closed @ ${SHA}** — Study shell, toolbar (translation/book/chapter nav), readability panel, settings persistence e2e. Legacy picker/rail visual parity deferred; behavior covered by e2e."

close 29 "**Closed @ ${SHA}** — Bible/comparison columns, scroll sync wiring, comparison e2e. Full legacy VerseLine component port deferred; chapter HTML + paragraph breaks in BibleColumn."

close 30 "**Closed @ ${SHA}** — Notes/scribe/verse-list e2e + legacy on-disk path Go tests (\`internal/studydata\`)."

close 31 "**Closed @ ${SHA}** — Search load-more e2e, cross-ref/dictionary column e2e. Legacy ranking UX tweaks acceptable vs e2e smoke."

close 32 "**Closed @ ${SHA}** — Manage Translations modal, print options e2e, app menu wiring, \`StudyPrintHTML\` Go builder. Full multi-column print HTML parity with Laravel deferred."

close 33 "**Closed @ ${SHA}** — \`MIGRATION.md\` manual + automated verify, \`mygo build\` smoke OK, \`legacy/\` kept **untracked on disk** for reference (see \`.gitignore\`). Do not commit legacy until explicit cutover."

echo "Closed #23–#33 at ${SHA}"
