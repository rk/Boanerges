# Boanerges

Local-first desktop study Bible, rewritten in [MyGo](https://mygo.egoist.dev) (Go + web frontend).

The previous Laravel / NativePHP implementation lives in [`legacy/`](legacy/) as the behavior spec until cutover.

## Develop

Requires Go 1.27+, Bun, and the platform webview (see `mygo doctor`).

```bash
bun install
bun run dev
```

## Build

```bash
bun run build
```

Packaged apps appear under `build/`.

## Legacy app

From `legacy/`:

```bash
composer install && npm ci
php artisan test
```

Legacy tests are not CI gates on `go-rewrite`.
