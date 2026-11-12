#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"
OUT="${OUT:-web}"
mkdir -p "$OUT"

RL="$(go list -m -f '{{.Dir}}' github.com/gen2brain/raylib-go/raylib)"
GOROOT="$(go env GOROOT)"

GOOS=js GOARCH=wasm CGO_ENABLED=0 go build -tags "${TAGS:-}" -ldflags="-s -w" -o "$OUT/app.wasm" .

if [ -f "$GOROOT/lib/wasm/wasm_exec.js" ]; then
  cp "$GOROOT/lib/wasm/wasm_exec.js" "$OUT/"
else
  cp "$GOROOT/misc/wasm/wasm_exec.js" "$OUT/"
fi

MODULE=raylib
TAGS="${TAGS:-}"
case " ${TAGS//,/ } " in *" es3 "*) MODULE=raylib-es3 ;; esac

gzip -dc "$RL/libs/wasm/$MODULE.js.gz" > "$OUT/raylib.js"
gzip -dc "$RL/libs/wasm/$MODULE.wasm.gz" > "$OUT/raylib.wasm"
cp "$RL/libs/wasm/raylib-go.js" index.html "$OUT/"

if [ -d resources ]; then
  tar -czf "$OUT/resources.tar.gz" resources
fi
