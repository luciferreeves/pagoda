#!/bin/bash
set -e

if [ "$SEED" = "true" ]; then
  echo "[entrypoint] SEED=true, wiping database..."
  if [ "$DB_DRIVER" = "sqlite" ]; then
    rm -f "$DSN"
    echo "[entrypoint] Deleted $DSN"
  elif [ "$DB_DRIVER" = "libsql" ]; then
    TURSO_URL=$(echo "$DSN" | sed 's|^libsql://|https://|; s|?.*||')
    TURSO_TOKEN=$(echo "$DSN" | sed -n 's/.*authToken=\(.*\)/\1/p')
    INDEXES=$(curl -s "${TURSO_URL}/v2/pipeline" \
      -H "Authorization: Bearer ${TURSO_TOKEN}" \
      -H "Content-Type: application/json" \
      -d '{"requests":[{"type":"execute","stmt":{"sql":"SELECT name FROM sqlite_master WHERE type='\''index'\'' AND name NOT LIKE '\''sqlite_%'\''"}},{"type":"close"}]}' \
      | grep -o '"type":"text","value":"[^"]*"' | sed 's/.*"value":"//;s/"//')
    TABLES=$(curl -s "${TURSO_URL}/v2/pipeline" \
      -H "Authorization: Bearer ${TURSO_TOKEN}" \
      -H "Content-Type: application/json" \
      -d '{"requests":[{"type":"execute","stmt":{"sql":"SELECT name FROM sqlite_master WHERE type='\''table'\'' AND name NOT LIKE '\''sqlite_%'\'' AND name NOT LIKE '\''_litestream%'\''"}},{"type":"close"}]}' \
      | grep -o '"type":"text","value":"[^"]*"' | sed 's/.*"value":"//;s/"//')
    STMTS=""
    for INDEX in $INDEXES; do
      echo "[entrypoint] Dropping index: $INDEX"
      STMTS="${STMTS}{\"type\":\"execute\",\"stmt\":{\"sql\":\"DROP INDEX IF EXISTS ${INDEX}\"}},"
    done
    for TABLE in $TABLES; do
      echo "[entrypoint] Dropping table: $TABLE"
      STMTS="${STMTS}{\"type\":\"execute\",\"stmt\":{\"sql\":\"DROP TABLE IF EXISTS ${TABLE}\"}},"
    done
    if [ -n "$STMTS" ]; then
      STMTS="${STMTS%,}"
      curl -s "${TURSO_URL}/v2/pipeline" \
        -H "Authorization: Bearer ${TURSO_TOKEN}" \
        -H "Content-Type: application/json" \
        -d "{\"requests\":[${STMTS},{\"type\":\"close\"}]}" > /dev/null
      echo "[entrypoint] All indexes and tables dropped"
    else
      echo "[entrypoint] Nothing to drop"
    fi
  fi
fi

echo "[entrypoint] Starting shrine server..."
./shrine &
SHRINE_PID=$!

if [ "$SEED" = "true" ]; then
  echo "[entrypoint] Waiting for server on port ${PORT:-3000}..."
  RETRIES=0
  MAX_RETRIES=30
  until bash -c "echo > /dev/tcp/localhost/${PORT:-3000}" 2>/dev/null; do
    RETRIES=$((RETRIES + 1))
    if [ "$RETRIES" -ge "$MAX_RETRIES" ]; then
      echo "[entrypoint] Server failed to start after ${MAX_RETRIES}s"
      exit 1
    fi
    echo "[entrypoint] Waiting... ($RETRIES/$MAX_RETRIES)"
    sleep 1
  done
  echo "[entrypoint] Server is up, running seed script..."
  bash scripts/seed.sh
  echo "[entrypoint] Seed complete"
else
  echo "[entrypoint] SEED not set, skipping seed"
fi

wait $SHRINE_PID