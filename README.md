# Reader

A high-performance web content reader service that extracts text content from web pages and provides AI-powered summaries.

## Features

- **Text Extraction**: Clean, readable text from any web page
- **Markdown Conversion**: Convert web content to well-formatted markdown
- **AI Summaries**: Generate concise summaries using AI
- **Format Options**: Choose between plain text or markdown output
- **Performance Optimized**: Browser pool for parallel processing
- **Caching**: Efficient content caching for faster responses
- **Metrics**: Built-in performance monitoring
- **Flexible Configuration**: Support for config files, environment variables, and command-line flags

## Installation

### From Source

```bash
git clone https://github.com/ncecere/reader.git
cd reader
go build -o reader cmd/reader/main.go
```

### Using Docker

Pull the image from GitHub Container Registry:

```bash
# Pull the latest version
docker pull ghcr.io/ncecere/reader:latest

# Or pull a specific version
docker pull ghcr.io/ncecere/reader:v1.5.0
```

Run with Docker:

```bash
docker run -p 4444:4444 -v $(pwd)/config.yml:/app/config/config.yml ghcr.io/ncecere/reader:latest
```

### Using Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'
services:
  reader:
    image: ghcr.io/ncecere/reader:latest
    ports:
      - "4444:4444"
    volumes:
      - ./config.yml:/app/config/config.yml
    environment:
      - READER_PORT=4444
      - READER_AI_ENDPOINT=https://ai.bitop.dev/v1
      - READER_AI_KEY=your-api-key
    restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

## Configuration

The application supports multiple configuration methods in order of precedence:
1. Command-line flags
2. Environment variables (prefixed with READER_)
3. Configuration file
4. Default values

### Command-Line Usage

```bash
# Show available commands and flags
./reader --help

# Run with default configuration
./reader run

# Run with custom config file
./reader run --config /path/to/config.yml

# Run with command-line flags
./reader run --port 4444 --pool-size 5 --chrome-path /usr/bin/chromium-browser

# Run with environment variables
READER_PORT=4444 READER_CHROME_PATH=/usr/bin/chromium-browser ./reader run
```

### Available Flags

```
Flags:
      --ai-enabled            Enable/disable AI features (default true)
      --ai-endpoint string    AI API endpoint
      --ai-key string         AI API key
      --ai-model string       AI model to use (default "vltr-mistral-small")
      --browser-timeout int   Browser request timeout in seconds (default 30)
      --chrome-path string    Path to Chrome/Chromium executable
      --config string         Config file path (default "./config.yml")
      --max-retries int      Maximum number of retries for browser operations (default 3)
      --pool-size int        Number of browser instances in the pool (default 3)
      --port int             Port to run the server on (default 4444)
```

### Environment Variables

All configuration options can be set via environment variables by prefixing with `READER_` and using uppercase. For example:
- `READER_PORT=4444`
- `READER_POOL_SIZE=5`
- `READER_CHROME_PATH=/usr/bin/chromium-browser`
- `READER_AI_ENDPOINT=https://ai.bitop.dev/v1`
- `READER_AI_KEY=your-api-key`

### Configuration File

The application uses a YAML configuration file. By default, it looks for `config.yml` in the current directory.

```yaml
# Server configuration
server:
  port: 4444
  pool_size: 3

# Browser configuration
browser:
  chrome_path: ""      # Path to Chrome/Chromium executable
  timeout: 30          # Request timeout in seconds
  max_retries: 3      # Maximum retries for browser operations

# AI configuration
ai:
  enabled: true
  api_endpoint: "https://ai.bitop.dev/v1"
  api_key: "your-api-key"
  model: "vltr-mistral-small"
  prompt: |
    As a summarization assistant...

# Logging configuration
logging:
  level: "info"
  json: true
  caller: true
```

## API Usage

### Extract Text

```bash
# Get plain text
curl -s -H "X-Respond-With: text" "http://localhost:4444/https://example.com"

# Get markdown
curl -s -H "X-Respond-With: markdown" "http://localhost:4444/https://example.com"
```

### Generate AI Summary

```bash
# Get plain text summary
curl -s -H "X-Respond-With: text" "http://localhost:4444/summary/https://example.com"

# Get markdown summary
curl -s -H "X-Respond-With: markdown" "http://localhost:4444/summary/https://example.com"
```

[Rest of the README remains unchanged...]
