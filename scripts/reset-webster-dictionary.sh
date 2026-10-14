#!/usr/bin/env bash
# Clear Webster 1828 import state in the user bible DB (same as in-app Reinstall).
set -euo pipefail

if [[ -n "${BOANERGES_DB:-}" ]]; then
  db="$BOANERGES_DB"
elif [[ "$(uname -s)" == Darwin ]]; then
  db="$HOME/Library/Application Support/Boanerges/bible.sqlite"
elif [[ -n "${XDG_DATA_HOME:-}" ]]; then
  db="$XDG_DATA_HOME/Boanerges/bible.sqlite"
else
  db="$HOME/.local/share/Boanerges/bible.sqlite"
fi

if [[ ! -f "$db" ]]; then
  echo "No database at: $db" >&2
  echo "Set BOANERGES_DB to your bible.sqlite path." >&2
  exit 1
fi

if ! command -v sqlite3 >/dev/null 2>&1; then
  echo "sqlite3 not found. Use Boanerges Settings → Reinstall Webster dictionary instead." >&2
  exit 1
fi

sqlite3 "$db" <<'SQL'
DELETE FROM import_meta WHERE key = 'webster1828_dictionary_v1';
DELETE FROM dictionary_entries;
SQL

echo "Webster dictionary cleared in: $db"
echo "Open Boanerges and look up a word (or Settings → Reinstall Webster dictionary)."
