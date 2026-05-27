# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first (cached layer unless go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/main.go

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.21

# CA certificates needed for HTTPS calls (Firebase, MongoDB Atlas, etc.)
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/server .

# Run as non-root
USER app

EXPOSE 8080

ENTRYPOINT ["./server"]
