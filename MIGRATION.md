# Migrating from NativePHP Boanerges

MyGo stores data under `mygo.PathUserData` (macOS: `~/Library/Application Support/Boanerges/`).

## Optional one-time copy

If you used the NativePHP app, you can copy SQLite and study files once:

1. **Bible database** — copy `legacy/storage/app/database/nativephp.sqlite` (or your configured SQLite path) to:
   - `~/Library/Application Support/Boanerges/bible.sqlite` (macOS)
   - Adjust for Linux/Windows user data dirs from MyGo docs.

2. **Notes** — copy `storage/app/private/notes/` into user data as `notes/` (same `bookId/chapter.txt` layout).

3. **Scribe** — copy `storage/app/private/scribe/` to `scribe/` (`bookId/chapter.json`).

4. **Verse lists** — copy `storage/app/private/verse-lists/` to `verse-lists/` (`index.json` + list JSON files).

5. **SWORD modules** — copy downloaded modules from Laravel `storage/app/private/sword/` to `~/Library/Application Support/Boanerges/sword/`.

Bundled ASV, cross-references, and Webster 1828 re-import automatically on first launch if missing from the destination database.

## Verify (manual)

After copying:

1. Launch the MyGo app (`bun run dev` or the built `.app`).
2. Confirm the welcome progress reaches **ready** and Study opens.
3. Open **Genesis 1** — verse text renders with markup (e.g. italics via `<em>`).
4. Toggle **File → Manage Translations** — bundled ASV shows installed; install another module or confirm catalog loads.
5. Add a chapter note, reload the app, confirm the note persists under `notes/gen/1.txt` in user data.
6. Run dictionary lookup on a common word (e.g. *love*) if Webster was imported.

If module tables are missing after copy, re-install catalog translations from **File → Manage Translations**.

## Verify (automated)

From the repo root:

```bash
go test ./...
bun run typecheck
bun run test          # unit tests (scroll sync, etc.)
bun run test:e2e      # Playwright against Vite + MyGO mock (`index.e2e.html`)
```

## macOS build smoke

On macOS with MyGo CLI installed:

```bash
bun run build         # mygo build — produces build/Boanerges.app (or platform output under build/)
```

`mygo.config.ts` lists only `extras` under `resources` (the `resources/` tree is included automatically). Record the git SHA when validating cutover.

CI runs `go test ./...` and `bun run typecheck`; run `bun run test:e2e` locally before closing cutover. macOS `bun run build` smoke is recommended when changing packaging or embedded assets.

## Legacy reference tree

The untracked `legacy/` directory on disk (see `.gitignore`) is optional local reference only (restore from git history before `6bb99d8` if needed). Do not commit `legacy/` to git.
