# Use a multi-stage build for a smaller final image
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o reader ./cmd/reader

# Final stage
FROM alpine:latest

# Install Chrome and dependencies
RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    && mkdir /screenshots

# Copy the binary from builder
COPY --from=builder /app/reader /usr/local/bin/

# Copy config file
COPY config.yml /etc/reader/

# Create non-root user
RUN adduser -D -h /home/reader reader && \
    chown -R reader:reader /screenshots && \
    chmod 755 /screenshots

# Switch to non-root user
USER reader

# Set environment variables
ENV CONFIG_PATH=/etc/reader/config.yml
ENV CHROME_PATH=/usr/bin/chromium-browser

# Expose port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["reader"]
