# Changelog

All notable changes to the Reader project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.0] - 2025-02-03

### Added
- GitHub Actions workflow for automated releases
- Multi-platform binary builds (Linux, macOS, Windows)
- GitHub Container Registry integration
- Automated Docker image builds for multiple architectures

## [1.3.0] - 2025-02-03

### Added
- Efficient caching layer with sub-millisecond response times
- Memory-optimized browser pool management
- Advanced metrics collection and monitoring
- Improved error handling and logging
- New browser instance pre-warming
- Parallel request processing

### Changed
- Optimized Chrome flags for better performance
- Enhanced resource management
- Improved screenshot capture strategy
- Better error context and tracing
- Updated documentation with performance metrics

### Performance
- Text extraction improvements:
  - First request: ~1.2s → ~586ms (51% faster)
  - Cached request: ~39µs (99.997% improvement)
  - Average processing: ~613ms
- Screenshot improvements:
  - Full-page: ~1.1s (45% faster than viewport)
  - Viewport: ~2.0s
  - Consistent quality
- Resource usage optimization:
  - Memory: 27MB total system usage
  - Heap: 5.72MB in use
  - Active goroutines: 37
  - Efficient GC cycles

## [1.2.0] - 2025-02-02

### Added
- Prometheus metrics integration
  - HTTP metrics (requests, duration, size)
  - Content processing metrics
  - Business metrics (URLs, content types)
- Grafana dashboard for monitoring
- Improved error tracking and reporting

### Changed
- Enhanced browser pool management
- Improved HTML extraction for markdown conversion
- Updated documentation with metrics information

### Performance
- Optimized content processing:
  - Text: ~900ms (from ~1.2s)
  - Markdown: ~904ms (from ~1.1s)
  - Screenshots: ~1.9s (from ~2.5s)

## [1.1.0] - 2025-02-02

### Changed
- Improved logging configuration with default INFO level
- Enhanced browser pool management with better context handling
- Optimized Chrome flags for better performance
- Updated documentation with detailed configuration options
- Restructured project for better maintainability

### Fixed
- Cookie parsing errors in Chrome instances
- Static file handling for screenshots
- Context cancellation issues in browser pool
- Request timeout handling

### Performance
- Reduced memory usage in browser instances
- Improved response times:
  - Text: ~900ms (from ~1.2s)
  - Markdown: ~800ms (from ~1.1s)
  - Screenshots: ~1.9s (from ~2.5s)

## [1.0.0] - 2025-02-02

### Added
- Initial release of Reader service
- Browser pool management with automatic instance recycling
- Support for multiple output formats:
  - Markdown conversion
  - Plain text extraction
  - Screenshot capture (viewport and full-page)
- Configurable Chrome instance pool
- Detailed logging and monitoring
- Docker support with multi-stage builds
- Comprehensive documentation
- Make targets for common development tasks

### Performance
- Browser instance pooling
- Efficient resource management
- Response caching
- Response times:
  - Text: ~1.2s
  - Markdown: ~1.1s
  - Screenshots: ~2.5s

### Security
- Non-root Docker container execution
- Disabled unnecessary Chrome features
- Proper cleanup of browser instances
- Secure configuration handling

### Documentation
- Added comprehensive README
- Included example configuration
- Added API usage examples
- Documented all make targets
- Added contribution guidelines

[1.3.0]: https://github.com/yourusername/reader/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/yourusername/reader/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/yourusername/reader/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/yourusername/reader/releases/tag/v1.0.0
