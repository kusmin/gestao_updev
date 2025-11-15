#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 ]]; then
  echo "Uso: $0 <diretorio> <script> [args...]" >&2
  exit 1
fi

TARGET_DIR="$1"
SCRIPT_NAME="$2"
shift 2

if [[ ! -f "$TARGET_DIR/package.json" ]]; then
  echo "Ignorando npm run $SCRIPT_NAME em $TARGET_DIR (package.json não encontrado)." >&2
  exit 0
fi

pnpm --prefix "$TARGET_DIR" run "$SCRIPT_NAME" "$@"
