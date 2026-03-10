#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR"

echo "[gate] backend unit tests"
make test-unit

echo "[gate] backend blackbox tests"
make test-blackbox

echo "[gate] backend regression tests"
make test-regression

echo "[gate] backend smoke tests"
make test-smoke

echo "[gate] frontend tests"
(cd apps/web && npm run test)

echo "[gate] frontend build"
(cd apps/web && npm run build)

echo "[gate] all checks passed"
