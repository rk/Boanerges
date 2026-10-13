#!/usr/bin/env bash
# Restore large bundled files from git if missing (Webster dictionary, cross-refs, etc.).
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

paths=(
  extras/websters1828/dictionary.json
  extras/websters1828/ATTRIBUTION.md
  extras/cross-references/cross_references.txt
)

missing=()
for p in "${paths[@]}"; do
  if [[ ! -f "$p" ]]; then
    missing+=("$p")
  fi
done

if ((${#missing[@]} == 0)); then
  echo "Bundled extras present."
  wc -c extras/websters1828/dictionary.json | awk '{print "  dictionary.json:", $1, "bytes"}'
  exit 0
fi

echo "Missing: ${missing[*]}"
echo "Restoring from HEAD…"
git checkout HEAD -- "${missing[@]}"
wc -c extras/websters1828/dictionary.json | awk '{print "  dictionary.json:", $1, "bytes"}'
