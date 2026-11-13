# Agentic Studio (compile-time feature)

Agentic Studio is a **separate window** for long-running, local **Ollama** jobs that distill dictionary entries and topical references. It is **not** in release builds.

## Develop

```bash
bun run dev:agentic
```

Requires Ollama on `http://127.0.0.1:11434`.

Open **File → Agentic Studio…** or load `/?view=studio`. Dictionary and topical project folders are remembered separately on this machine (local storage).

## Distributable

```bash
bun run build
```

Default `go build` and `bun run build:web` omit the `agentic` tag and `VITE_AGENTIC`.

## Project layout

| Path | Purpose |
|------|---------|
| `project.json` | `kind` (`dictionary` \| `topical`), `model` — each kind exposes different orchestrator tools (topical has no source/entry APIs) |
| `run.json` | Checkpoint: status, steps, instruction |
| `run.raw.log` | Full model replies and tool I/O (debug) |
| `sources/` | Public-domain source files (dictionary) |
| `source-guides.json` | Per-file rules for spotting headwords (from first 10k of each source) |
| `entries/*.json` | Modern-English dictionary entries |
| `topics/*.json` | Topical groups and scripture refs |

Drop source files into `sources/` with Finder before starting a dictionary run. Dictionary distillation is **multi-batch**: one “Start run” keeps going across many Ollama turns until sources are fully read (large files are chunked at ~12k characters per `read_source` call). Pause when needed; start again to continue from checkpoint offsets in `run.json`.

**Source formats:** plain UTF-8 text (`.txt`, `.md`, `.html`) works — the distiller reads raw bytes in chunks. Avoid PDF, DOCX, and EPUB unless you convert them to `.txt` first (binary or markup noise confuses the model).

## Go tests (agentic)

```bash
go test -tags agentic ./internal/agentic/...
```
