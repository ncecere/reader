FROM golang:1.22-bullseye as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o reader cmd/reader/main.go

FROM debian:bullseye-slim

# Install Chromium and dependencies
RUN apt-get update && \
    apt-get install -y \
    chromium \
    fonts-ipafont-gothic \
    fonts-wqy-zenhei \
    fonts-thai-tlwg \
    fonts-kacst \
    fonts-symbola \
    fonts-noto \
    fonts-freefont-ttf \
    curl \
    ca-certificates \
    dbus \
    xvfb \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && which chromium || (echo "Chromium not found" && exit 1)

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/reader .

# Create directory for config and ensure proper permissions
RUN mkdir -p /app/config && \
    mkdir -p /tmp/.X11-unix && \
    chmod 1777 /tmp/.X11-unix && \
    ln -s /usr/bin/chromium /usr/bin/chromium-browser

# Set environment variables
ENV PATH="/app:${PATH}"
ENV CHROME_PATH="/usr/bin/chromium"
ENV DISPLAY=":99"

# Expose port
EXPOSE 4444

# Start Xvfb and the application
CMD Xvfb :99 -screen 0 1024x768x16 & ./reader run
