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

RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

COPY --from=builder /shrine/bin/shrine .

CMD ["./shrine"]
