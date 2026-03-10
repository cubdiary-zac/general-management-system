#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
EMAIL="${SMOKE_EMAIL:-admin@gms.local}"
PASSWORD="${SMOKE_PASSWORD:-admin123}"

json() { jq -r "$1"; }

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required for smoke_api.sh" >&2
  exit 1
fi

echo "[smoke] health"
HEALTH=$(curl -sS -f "$BASE_URL/api/health")
echo "$HEALTH" | json '.status' >/dev/null

echo "[smoke] login"
LOGIN=$(curl -sS -f -X POST "$BASE_URL/api/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
TOKEN=$(echo "$LOGIN" | json '.token')
if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "login failed: missing token" >&2
  exit 1
fi

echo "[smoke] create project"
PROJECT=$(curl -sS -f -X POST "$BASE_URL/api/pm/projects" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Smoke Script Project"}')
PROJECT_ID=$(echo "$PROJECT" | json '.id')

echo "[smoke] create task"
TASK=$(curl -sS -f -X POST "$BASE_URL/api/pm/tasks" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"projectId\":$PROJECT_ID,\"title\":\"Smoke Script Task\"}")
TASK_ID=$(echo "$TASK" | json '.id')

echo "[smoke] status transition"
for STEP in in_progress in_review done; do
  curl -sS -f -X PATCH "$BASE_URL/api/pm/tasks/$TASK_ID/status" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"status\":\"$STEP\"}" >/dev/null
  echo "  -> $STEP"
done

echo "[smoke] ok"
