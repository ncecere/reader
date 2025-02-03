# Changelog

All notable changes to the Reader project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[1.2.0]: https://github.com/yourusername/reader/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/yourusername/reader/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/yourusername/reader/releases/tag/v1.0.0
