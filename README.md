# Boanerges

Local-first desktop study Bible built with [MyGo](https://mygo.egoist.dev) (Go services + Svelte 5 frontend).

## Develop

Requires Go 1.27+, Bun, and a working platform webview (`mygo doctor`).

```bash
bun install
bun run dev
```

`bun run dev` runs MyGo, which starts Vite on port **5173** (`dev:web`). If you see **Port 5173 is already in use**, a stale Vite or old **Boanerges Dev** session is still running — quit the dev app or run `bun scripts/free-dev-port.ts`, then try again.

Optional two-terminal workflow: `bun run dev:web` in one terminal, `bunx mygo dev -skip-dev-command` in another (after Vite is up).

After changing bound Go APIs:

```bash
bun run generate
```

Typecheck:

```bash
bun run typecheck
```

Go tests:

```bash
go test ./...
```

## Build

```bash
bun run build
```

macOS `.app` bundles land in `build/`. Code signing and auto-update are out of scope for this rewrite.

## Data layout

User data (SQLite bible DB, notes, scribe, verse lists, downloaded SWORD zips) lives under MyGo user data. See [MIGRATION.md](MIGRATION.md) for optional import from the old NativePHP app.

## Bundled resources

Shipped under `extras/` (ASV SQLite, SWORD ASV module, translation catalog, cross-references, Webster 1828). MyGo copies them via `mygo.config.ts` `resources`.
