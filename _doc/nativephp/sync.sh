#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
DEST="$ROOT/docs"
REPO="https://github.com/NativePHP/nativephp.com.git"
SPARSE_PATH="resources/views/docs/desktop/2"

echo "Syncing NativePHP Desktop v2 docs into $DEST ..."

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

git clone --depth 1 --filter=blob:none --sparse "$REPO" "$TMP/repo"
git -C "$TMP/repo" sparse-checkout set "$SPARSE_PATH"

rm -rf "$DEST"
mkdir -p "$DEST"
cp -R "$TMP/repo/$SPARSE_PATH/." "$DEST/"

echo "Done. $(find "$DEST" -name '*.md' | wc -l | tr -d ' ') markdown files in $DEST"
