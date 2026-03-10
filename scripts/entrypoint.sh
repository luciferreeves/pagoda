#!/bin/bash
set -e

./shrine &
SHRINE_PID=$!

if [ "$SEED" = "true" ]; then
  RETRIES=0
  MAX_RETRIES=30
  until bash -c "echo > /dev/tcp/localhost/${PORT:-3000}" 2>/dev/null; do
    RETRIES=$((RETRIES + 1))
    if [ "$RETRIES" -ge "$MAX_RETRIES" ]; then
      echo "Server failed to start after ${MAX_RETRIES}s"
      exit 1
    fi
    sleep 1
  done
  bash scripts/seed.sh
fi

wait $SHRINE_PID