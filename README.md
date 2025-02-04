# Reader Service

A high-performance web page content extraction service with support for text extraction, markdown conversion, and screenshot capture.

## Features

- High-performance content extraction with caching
- Multiple output formats (text, markdown, screenshots)
- Memory-optimized browser pool management
- Sub-millisecond cache response times
- Full-page and viewport screenshots
- Parallel request processing
- Advanced metrics and monitoring
- Docker support with multi-stage builds
- Hot reload development environment

## Performance

### Text Extraction
- First Request: ~586ms
- Cached Request: ~39µs (99.997% improvement)
- Average Processing: ~613ms

### Screenshots
- Full-page: ~1.1s
- Viewport: ~2.0s
- Consistent quality and size (~1.5MB)

### Resource Usage
- Total Memory: 27MB
- Heap Usage: 5.72MB
- Active Goroutines: 37
- Efficient GC cycles

## Requirements

- Go 1.21 or higher
- Chrome/Chromium (for local development)
- Docker (for containerized deployment)
- Make

## Quick Start

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/yourusername/reader.git
cd reader
```

2. Install dependencies:
```bash
make mod-tidy
```

3. Run development server with hot reload:
```bash
make dev
```

### Docker Development

Run the development environment with hot reload:
```bash
docker-compose up dev
```

### Production Deployment

Build and run the production container:
```bash
docker-compose up reader
```

## API Endpoints

### Get Text Content
```http
GET /{url}
X-Respond-With: text

Response:
Plain text content of the webpage
Time: ~586ms (first request), ~39µs (cached)
```

### Get Markdown Content
```http
GET /{url}
X-Respond-With: markdown

Response:
Markdown formatted content with:
- Title
- Visit timestamp
- Converted content with links preserved
Time: ~600ms (first request), ~40µs (cached)
```

### Capture Screenshot
```http
GET /{url}
X-Respond-With: screenshot (viewport)
X-Respond-With: pageshot (full-page)

Response:
PNG image of the webpage
Time: ~2.0s (viewport), ~1.1s (full-page)
Size: ~1.5MB
```

### Health Check
```http
GET /health
```

### Metrics
```http
GET /metrics

Response:
Prometheus metrics including:
- Request latencies
- Cache hit rates
- Memory usage
- Pool utilization
- Content sizes
```

## Configuration

Configuration can be handled through:
1. Environment variables
2. config.yml file
3. Command line flags

### Configuration File

The service uses a YAML configuration file (config.yml). A full example with all options is available in `config.example.yml`.

#### Server Configuration
```yaml
server:
  port: 4444                # Port to listen on
  host: "0.0.0.0"          # Host to bind to
  read_timeout: 30         # Read timeout in seconds
  write_timeout: 30        # Write timeout in seconds
  idle_timeout: 60         # Idle timeout in seconds
  max_body_size: 10        # Maximum request body size in MB
  cors_enabled: false      # Enable CORS
  cors_origins: ["*"]      # CORS allowed origins
  compression: true        # Enable compression
```

#### Browser Configuration
```yaml
browser:
  pool_size: 3             # Number of Chrome instances
  chrome_path: ""          # Optional Chrome executable path
  timeout: 30              # Request timeout in seconds
  max_memory_mb: 128       # Maximum memory per instance
  retries: 3               # Number of retries for failed requests
  retry_delay: 1           # Delay between retries
  prewarming: true        # Enable instance pre-warming
  window_width: 1920      # Screenshot window width
  window_height: 1080     # Screenshot window height
```

#### Cache Configuration
```yaml
cache:
  enabled: true           # Enable caching
  ttl: 3600              # Cache TTL in seconds
  max_size_mb: 256       # Maximum cache size
  cleanup_interval: 300   # Cache cleanup interval
  stale_revalidate: true # Enable stale-while-revalidate
  compression_level: 6    # Cache compression level (0-9)
```

#### Performance Configuration
```yaml
performance:
  parallel_processing: true     # Enable parallel processing
  max_parallel_requests: 10     # Maximum parallel requests
  response_compression: true    # Enable response compression
  compression_level: 6         # Compression level (1-9)
  response_caching: true       # Enable response caching
  response_cache_size: 128     # Cache size in MB
```

### Environment Variables

All configuration options can be set via environment variables using the format:
```bash
READER_[SECTION]_[KEY]=value

# Examples:
READER_SERVER_PORT=4444
READER_BROWSER_POOL_SIZE=3
READER_CACHE_TTL=3600
```

### Command Line Flags

Basic configuration can also be provided via command line flags:
```bash
./reader -port=4444 -pool-size=3 -log-level=info
```

For a complete list of configuration options and their defaults, see `config.example.yml`.

## Development

### Available Make Commands

- `make build`: Build the binary
- `make test`: Run tests
- `make coverage`: Generate test coverage report
- `make lint`: Run linter
- `make fmt`: Format code
- `make docker-build`: Build Docker image
- `make docker-run`: Run Docker container
- `make dev`: Run development server with hot reload
- `make help`: Show available commands

### Project Structure

```
.
├── cmd/
│   └── reader/           # Application entry point
├── internal/
│   ├── api/             # API handlers and middleware
│   ├── common/          # Shared utilities and configurations
│   └── core/            # Core business logic
│       ├── browser/     # Browser automation
│       ├── cache/       # Caching layer
│       ├── metrics/     # Metrics collection
│       └── converter/   # Content conversion
├── dashboards/          # Grafana dashboards
├── config.yml          # Configuration file
├── Dockerfile         # Multi-stage Docker build
├── docker-compose.yml # Container orchestration
└── Makefile          # Build automation
```

## Performance Optimization

The service includes several optimizations:

- Efficient caching layer with sub-millisecond response
- Memory-optimized browser pool
- Parallel request processing
- Resource cleanup and management
- Browser instance pre-warming
- Configurable timeouts and retries
- Memory usage optimization
- Efficient garbage collection

## Monitoring

The service exposes detailed Prometheus metrics at `/metrics`:

- Request latencies and counts
- Cache hit rates and sizes
- Memory usage and GC stats
- Pool utilization
- Content processing times
- Error rates and types

Includes a comprehensive Grafana dashboard for visualization.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
