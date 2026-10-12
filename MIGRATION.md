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

## Verify

After copying, launch the MyGo app, open Study, and confirm translations, notes, and dictionary lookup work. Re-install catalog translations from **File → Manage Translations** if module tables are missing.
