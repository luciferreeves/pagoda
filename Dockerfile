FROM golang:1.25 AS builder

WORKDIR /shrine

RUN apt-get update && apt-get install -y gcc libc6-dev make git

ENV CGO_ENABLED=1

COPY shrine/go.mod shrine/go.sum* ./
RUN go mod download

COPY shrine/ .
RUN make build

FROM debian:bookworm-slim

WORKDIR /shrine

RUN apt-get update && \
    apt-get install -y ca-certificates tzdata curl sqlite3 apache2-utils \
    chromium && \
    rm -rf /var/lib/apt/lists/*

ENV CHROME_PATH=/usr/bin/chromium

COPY --from=builder /shrine/bin/shrine .
COPY --from=builder /shrine/templates ./templates
COPY scripts/ ./scripts/
COPY seed/ ./seed/
RUN chmod +x /shrine/scripts/entrypoint.sh

CMD ["bash", "scripts/entrypoint.sh"]
