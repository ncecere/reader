# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make build-base

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git rev-parse --short HEAD) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o build/reader \
    cmd/reader/main.go

# Runtime stage
FROM alpine:latest

# Install Chrome dependencies with optimizations
RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    && rm -rf /var/cache/* /tmp/* /var/tmp/*

# Set Chrome environment variables and flags for performance
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/ \
    CHROME_FLAGS="--disable-dev-shm-usage --no-sandbox --disable-gpu --headless --disable-software-rasterizer"

# Create non-root user
RUN adduser -D -h /home/appuser appuser
USER appuser
WORKDIR /home/appuser

# Create cache and data directories
RUN mkdir -p \
    screenshots \
    cache \
    /tmp/chrome \
    && chown -R appuser:appuser \
    screenshots \
    cache \
    /tmp/chrome

# Set up tmpfs for Chrome
VOLUME ["/tmp/chrome"]

# Copy binary from builder
COPY --from=builder /app/build/reader .

# Copy config file
COPY config.yml .

# Expose ports
EXPOSE 4444

# Set resource limits and runtime options
ENV GOMAXPROCS=4 \
    GOGC=100 \
    GOMEMLIMIT=128MiB

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:4444/health || exit 1

# Set entrypoint with optimized flags
ENTRYPOINT ["./reader"]
