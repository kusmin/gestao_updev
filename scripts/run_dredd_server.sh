#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR/backend"

DATABASE_URL="${DATABASE_URL:-postgres://testuser:testpassword@localhost:5433/testdb?sslmode=disable}"

OTEL_ENABLED=false \
OTEL_RESOURCE_ATTRIBUTES="" \
DATABASE_URL="$DATABASE_URL" \
go run cmd/api/main.go
