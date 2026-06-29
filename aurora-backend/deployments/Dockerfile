# ─── Build stage ───
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go.mod first, then download (go.sum generated inside container)
COPY go.mod ./
RUN go mod download

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o aurora ./cmd/aurora

# ─── Run stage ───
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /build/aurora .
COPY configs/ ./configs/

RUN addgroup -S aurora && adduser -S aurora -G aurora
USER aurora

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

ENTRYPOINT ["./aurora", "-config", "configs/config.yaml"]
