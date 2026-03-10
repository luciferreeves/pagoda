#!/bin/bash
set -e

echo "[entrypoint] Starting shrine server..."
./shrine &
SHRINE_PID=$!

if [ "$SEED" = "true" ]; then
  echo "[entrypoint] SEED=true, waiting for server on port ${PORT:-3000}..."
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

echo "[entrypoint] Waiting on server process..."
wait $SHRINE_PID