# Boanerges

Boanerges is a small passion project: a desktop study Bible app, that may one day grow into something more. It's tailored to my wants, so it might not be what you're looking for.

It's named for the "sons of thunder", James and John. Jesus dubbed them this name in Mark 3:17.

## Purpose

- To provide a basic study Bible app for free.
- To provide a unified notes structure across all translations.
- To provide enhanced search capabilities.

## Structure

```
- _doc - internal planning documents, ignored in .git.
- _doc/wiki - internal documentation of the structure.
- _doc/nativephp - NativePHP Desktop documentation.
- _doc/[library] - any critical dependency documentation.
```

## Architecture

- Local first, no authentication or cloud.
- Self-contained as much as possible, with no 3rd party connections or telemetry.
- Minimal dependency tree.
- Tests and components first.

## Bible module setup

The ASV translation is bundled under `extras/sword/` and ships with NativePHP builds. Additional translations are downloaded to the local disk at `modules/bible/` and persist in app data.

Use the **Manage translations** button (square plus icon) in the sidebar to install or remove English translations from the bundled catalog (`extras/translations.json`).

Verify the bundled ASV module:

```bash
php artisan bible:verify-asv
```