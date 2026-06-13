#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
OUT_DIR="$ROOT_DIR/bin"
OUT_FILE="$OUT_DIR/cognitor.exe"
GOOS_VALUE="${GOOS:-windows}"
GOARCH_VALUE="${GOARCH:-amd64}"

mkdir -p "$OUT_DIR"
cd "$ROOT_DIR"

GOOS="$GOOS_VALUE" GOARCH="$GOARCH_VALUE" go build -trimpath -ldflags="-s -w" -o "$OUT_FILE" ./cmd/cognitor

printf 'built %s for %s/%s\n' "$OUT_FILE" "$GOOS_VALUE" "$GOARCH_VALUE"
