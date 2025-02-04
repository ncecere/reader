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
docker run -p 4444:4444 -v $(pwd)/config.yml:/app/config.yml ghcr.io/ncecere/reader:latest
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
      - ./config.yml:/app/config.yml
    restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

## Configuration

The application is configured through a YAML file. By default, it looks for `config.yml` in the current directory, but you can specify a different path using the `--config` flag.

### Configuration Options

```yaml
# Server configuration
server:
  port: 4444           # Port to run the server on
  pool_size: 3         # Number of browser instances in the pool

# Browser configuration
browser:
  chrome_path: ""      # Optional: Path to Chrome/Chromium executable
  timeout: 30          # Request timeout in seconds
  max_memory: 128      # Maximum memory per browser instance in MB
  user_agent: ""       # Optional: Custom user agent string

# AI configuration
ai:
  enabled: true                    # Enable/disable AI features
  api_endpoint: "..."             # AI API endpoint
  api_key: "your-api-key"        # API key for authentication
  model: "gpt-3.5-turbo"         # Model to use for summarization
  prompt: "..."                   # Custom prompt for summarization

# Cache configuration
cache:
  max_age: 3600       # Maximum age of cached items in seconds
  max_items: 1000     # Maximum number of items to keep in cache

# Metrics configuration
metrics:
  enabled: true       # Enable/disable Prometheus metrics
  path: "/metrics"    # Path for metrics endpoint
```

### Configuration Details

#### Server Options
- `port`: The port number the server will listen on
- `pool_size`: Number of concurrent browser instances to maintain. Higher numbers allow more parallel processing but use more memory

#### Browser Options
- `chrome_path`: Optional path to Chrome/Chromium executable. If empty, uses system default
- `timeout`: Maximum time in seconds for a request to complete
- `max_memory`: Maximum memory in MB allocated per browser instance
- `user_agent`: Optional custom user agent string for requests

#### AI Options
- `enabled`: Enable or disable AI summarization features
- `api_endpoint`: URL of the AI API endpoint
- `api_key`: Authentication key for the AI service
- `model`: AI model to use for summarization
- `prompt`: Custom system prompt for better summarization results

#### Cache Options
- `max_age`: How long to keep items in cache (in seconds)
- `max_items`: Maximum number of items to store in cache

#### Metrics Options
- `enabled`: Enable or disable Prometheus metrics
- `path`: Endpoint path for accessing metrics

## Usage

### Start the Server

```bash
# Using default config.yml
./reader

# Using custom config file
./reader --config /path/to/config.yml
```

### Extract Text

```bash
# Get plain text
curl http://localhost:4444/https://example.com

# Get markdown
curl -H "X-Respond-With: markdown" http://localhost:4444/https://example.com
```

### Generate AI Summary

```bash
# Get plain text summary
curl http://localhost:4444/summary/https://example.com

# Get markdown summary
curl -H "X-Respond-With: markdown" http://localhost:4444/summary/https://example.com
```

## Response Formats

### Text Format
Plain text extraction of the web page content, with:
- Clean, readable text
- Preserved structure
- Removed ads and clutter
- Maintained paragraph formatting

### Markdown Format
Structured markdown version of the content, including:
- Headers (h1-h6)
- Lists (ordered and unordered)
- Links with proper formatting
- Basic text formatting (bold, italic)
- Tables (when present)
- Code blocks (with syntax highlighting hints)

### AI Summary Format
Concise summary focusing on:
- Main ideas and key points
- Factual accuracy
- Important context
- Clear language
- Structured format (with markdown option)

## Performance

The service is optimized for performance through several mechanisms:

- **Browser Pool**: Multiple browser instances handle requests in parallel
- **Caching System**: Frequently accessed content is cached for faster retrieval
- **Memory Management**: Optimized browser instances with controlled memory usage
- **Efficient Processing**: Streamlined content extraction and processing
- **Connection Pooling**: Reuse of connections for better performance

## Metrics

Access Prometheus metrics at `/metrics` for monitoring:

- **Request Metrics**:
  - Total requests
  - Request durations
  - Response sizes
  - Error rates

- **Cache Metrics**:
  - Hit/miss rates
  - Cache size
  - Item age distribution

- **System Metrics**:
  - Memory usage
  - Browser pool status
  - Process statistics

## Development

### Prerequisites

- Go 1.21 or later
- Chrome/Chromium browser
- Make (optional, for using Makefile commands)

### Build

```bash
make build
# or
go build -o reader cmd/reader/main.go
```

### Test

```bash
make test
# or
go test ./...
```

### Development Server

```bash
make dev
# or
go run cmd/reader/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
