# Boanerges scratchpad

Informal backlog — not a commitment order. Check items off as they land.

---

## Study print (broken — needs proper design)

**Docs:** [MyGo — Printing, PDF and screenshots](https://mygo.egoist.dev/docs/windows#printing-pdf-and-screenshots)

MyGo exposes first-class APIs on **`win.Page()`**:

- `Page().Print()` — system print dialog (print CSS / `@page` apply)
- `Page().PrintToPDF(PDFOptions)` — layout for print; `PageSize`, `Landscape`, `Margins`, `Background`
- `CapturePage()` — PNG of what the window shows (not the same as print layout)

Native UI windows (`Content: ui.View(...)`) have their **own** inspector and platform print path; see the same doc section and native UI guides. That may be a better long-term fit for “real” print control than hacking a throwaway webview window.

### Current state (go-rewrite)

| Piece | Status |
| --- | --- |
| `internal/studyprint/*` | Legacy `StudyPrintHtmlBuilder` parity — multi-column HTML, lined notes/scribe, markup spans |
| `PrintOptionsModal` + `src/lib/printStudy.ts` | UX wired; calls `AppChrome.PrintStudy` |
| `AppChrome.PrintStudy` (`internal/appchrome/print_study.go`) | **Broken** — opens a secondary window, loads preview HTML; print dialog often does not appear; devtools not useful on that window; `window.print()` / `afterprint` workarounds are unreliable |

**Do not stack more webview workarounds** until we pick an approach below and validate it against MyGo on macOS (and other targets we ship).

### Symptoms observed

- Preview window shows formatted HTML but **no print sheet**
- Print helper window hard to inspect (`DevTools` / menu not practical)
- Earlier: modal `$effect` loop (fixed via `syncDialog` action); separate from print pipeline

### Legacy reference

- Laravel: `StudyPrintService` → `System::print($html, $settings)` (NativePHP) with landscape + A4 + backgrounds
- HTML template: `resources/views/print/study.blade.php` (removed with legacy tree; logic lives in `internal/studyprint/render.go`)

### Direction (pick one, then implement cleanly)

1. **Main-window print host (minimal spike)**  
   Load study print HTML into the **main** `Boanerges` window (or a dedicated long-lived print host window with `PageOptions{DevTools: DevToolsEnabled}`), call `Page().Print()` once, restore study URL. Prove MyGo’s `Print()` on a stable window — no immediate `Close()`.

2. **PDF-first UX (matches MyGo docs literally)**  
   `PrintToPDF` with options derived from study (`columnCount > 1` → `Landscape`, `PageA4`, `Background: true`, margins ~12mm). User prints or saves via system PDF workflow; optional “Save PDF…” in modal instead of fighting the sheet.

3. **Native UI print plugin / extension**  
   Research whether a small **MyGo plugin** (or upstream contribution) can wrap `NSPrintOperation` / platform print with: blocking until dismiss, landscape from options, HTML or paginated native layout. Goal: **real print control**, not `eval(window.print())`.

4. **Hybrid**  
   Keep HTML builder; render to PDF in-app for preview + “Export PDF”; use `Print()` only when (1) or (3) is solid.

### Concrete next steps

- [ ] Re-read MyGo `Page.Print` / `PrintToPDF` source (`github.com/egoist/mygo`) — confirm async behavior on darwin and whether `Print()` must run on main window
- [ ] Spike (1): main window or devtools-enabled print window + `LoadFile` + `Print()` only — document exact behavior
- [ ] Spike (2): `PrintToPDF` from generated HTML loaded in page; wire “Save PDF” in modal
- [ ] Map `PDFOptions` / print settings to `study.columnCount`, readability font metrics (legacy landscape rule)
- [ ] Decide: stay web-page print vs native UI renderer for study columns
- [ ] Remove or gate broken `PrintStudy` path until replacement ships (avoid shipping silent failure)
- [ ] E2E: keep mock `AppChrome.PrintStudy`; add integration test when real path exists

### Code touchpoints

- `internal/appchrome/print_study.go`, `service.go` (`Print`, `PrintToPDF`, `StudyPrintHTML`, `WriteStudyPrintFile`)
- `internal/studyprint/`
- `src/components/print/PrintOptionsModal.svelte`, `src/lib/printStudy.ts`

---

## Other (add as needed)

- [ ] …
