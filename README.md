# Reader

Reader is a high-performance web content extraction service that converts web pages into various formats including markdown, plain text, and screenshots. It uses Chrome in headless mode with an efficient browser pool for optimal performance.

## Features

- Convert web pages to clean markdown
- Extract plain text content
- Capture full or viewport screenshots
- High-performance browser pool management
- Prometheus metrics integration
- Structured JSON logging
- RESTful API endpoints

## Performance

The service is optimized for performance with real-time metrics:
- Browser instance pooling with automatic recycling
- Efficient resource management
- Response times (v1.1.0):
  - Text: ~900ms
  - Markdown: ~904ms
  - Screenshots: ~1.9s

## Requirements

- Go 1.21 or later
- Chrome/Chromium browser
- Docker (optional, for containerized deployment)
- Prometheus (optional, for metrics collection)

## Installation

### Local Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/reader.git
cd reader
```

2. Install dependencies:
```bash
go mod download
```

3. Create a configuration file (config.yml):
```yaml
server:
  port: 8080
  timeout: 60

browser:
  poolSize: 2
  chromePath: ""  # System default
  timeout: 30
  maxRetries: 3

screenshots:
  storagePath: "screenshots"
  quality: 90
  defaultType: "viewport"

logging:
  level: "info"
  json: true
  caller: true
```

4. Run the application:
```bash
go run cmd/reader/main.go
```

### Docker Setup

1. Build the Docker image:
```bash
docker build -t reader .
```

2. Run the container:
```bash
docker run -p 8080:8080 reader
```

## Configuration

The application is configured through a YAML file (config.yml). Key configuration sections:

### Server Settings
| Key         | Description                           | Default     |
|-------------|---------------------------------------|-------------|
| port        | HTTP server port                      | 8080        |
| timeout     | Request timeout in seconds            | 60          |

### Browser Settings
| Key         | Description                           | Default     |
|-------------|---------------------------------------|-------------|
| poolSize    | Number of Chrome instances            | 2           |
| chromePath  | Path to Chrome executable             | (system)    |
| timeout     | Operation timeout in seconds          | 30          |
| maxRetries  | Failed operation retries              | 3           |

### Screenshot Settings
| Key         | Description                           | Default     |
|-------------|---------------------------------------|-------------|
| storagePath | Screenshot storage directory          | screenshots |
| quality     | Image quality (1-100)                 | 90          |
| defaultType | Default capture type                  | viewport    |

### Logging Settings
| Key         | Description                           | Default     |
|-------------|---------------------------------------|-------------|
| level       | Log level (debug,info,warn,error)     | info        |
| json        | Enable JSON formatting                | true        |
| caller      | Include caller information            | true        |
| file        | Log file path (empty for stdout)      | ""          |

See [config.yml](config.yml) for detailed configuration options and examples.

## API Usage

### Convert to Markdown

```bash
curl -H "X-Respond-With: markdown" http://localhost:8080/https://example.com
```

### Extract Text

```bash
curl -H "X-Respond-With: text" http://localhost:8080/https://example.com
```

### Capture Screenshot

```bash
# Viewport screenshot
curl -H "X-Respond-With: screenshot" http://localhost:8080/https://example.com

# Full-page screenshot
curl -H "X-Respond-With: pageshot" http://localhost:8080/https://example.com
```

## Monitoring

### Available Metrics

The service exposes Prometheus metrics at `/metrics`. Key metrics include:

#### HTTP Metrics
- `reader_http_requests_total` - Total requests by endpoint and status code
- `reader_http_request_duration_seconds` - Request latency distribution
- `reader_http_request_size_bytes` - Request size distribution
- `reader_http_response_size_bytes` - Response size distribution
- `reader_http_in_flight_requests` - Current number of in-flight requests

#### Operation Metrics
- `reader_content_processing_duration_seconds` - Content processing time by type
- `reader_content_size_bytes` - Processed content size by type
- `reader_content_processing_errors_total` - Processing errors by type

#### Business Metrics
- `reader_url_processing_total` - URLs processed by domain
- `reader_url_content_types_total` - Content types encountered
- `reader_url_sizes_bytes` - URL content size distribution

### Prometheus Configuration

Example Prometheus scrape config:

```yaml
scrape_configs:
  - job_name: 'reader'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboard

A sample Grafana dashboard is available in [dashboards/reader.json](dashboards/reader.json).

## Development

### Available Make Commands

```bash
# Build the application
make build

# Run tests
make test

# Run with hot reload
make dev

# Format code
make fmt

# Run linters
make lint

# Build Docker image
make docker-build

# Show all commands
make help
```

### Project Structure

```
reader/
├── cmd/
│   └── reader/          # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP request handlers
│   │   └── middleware/  # HTTP middleware components
│   ├── common/
│   │   ├── config/     # Configuration management
│   │   ├── logger/     # Structured logging
│   │   └── metrics/    # Prometheus metrics
│   └── core/
│       ├── browser/    # Chrome browser management
│       └── converter/  # Content conversion utilities
└── screenshots/        # Screenshot storage
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
