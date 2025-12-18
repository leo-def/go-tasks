# syntax=docker/dockerfile:1

# --- Builder stage ---
FROM golang:1.25 AS builder
WORKDIR /app

# Enable reproducible builds and faster downloads
ENV CGO_ENABLED=0 \
    GOFLAGS="-trimpath"

# Pre-cache deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source and build
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o /out/go-tasks ./cmd/api

# --- Dev stage (with Go toolchain for hot reload) ---
FROM golang:1.25 AS dev
WORKDIR /app
ENV CGO_ENABLED=0 \
    GOTOOLCHAIN=auto
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .

# --- Runtime stage ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 appuser
WORKDIR /home/appuser
COPY --from=builder /out/go-tasks /usr/local/bin/go-tasks
USER appuser
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/go-tasks"]