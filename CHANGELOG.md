# Changelog

All notable changes to this project will be documented in this file.

## [v1.5.1] - 2025-02-04

### Added
- Cobra CLI support for flexible configuration
- Environment variable support with READER_ prefix
- Command-line flags for all configuration options

### Changed
- Improved Docker configuration for better Chrome support
- Enhanced browser pool initialization and error handling
- Updated configuration system to use Viper
- Moved compose.yml to local-only file

### Fixed
- Browser pool nil pointer issues
- Chrome path handling in container
- Configuration loading edge cases

## [v1.5.0] - 2025-02-03

### Removed
- Screenshot and pageshot functionality to streamline core features
- Screenshot-related code and directories
- Screenshot configuration options

### Changed
- Improved SSL/Security handling in Chrome flags
- Optimized browser options for better performance
- Enhanced context timeout handling in request handlers

### Fixed
- SSL protocol errors in web requests
- Context timeout issues in handlers
- Memory optimization in browser pool

## [v1.4.0] - 2025-01-15

### Added
- AI summarization feature with OpenAI API integration
- Support for both text and markdown summary formats
- Configurable AI settings (API endpoint, key, model, prompt)
- Caching for AI summaries
- New /summary endpoint

### Changed
- Enhanced error handling for API requests
- Improved response formatting
- Updated configuration structure for AI settings

## [v1.3.0] - 2024-12-10

### Added
- Markdown conversion support
- HTML to Markdown transformation
- Format selection via X-Respond-With header
- Enhanced text extraction

### Changed
- Improved HTML parsing
- Better error handling
- Updated documentation

## [v1.2.0] - 2024-11-20

### Added
- Browser pool for parallel processing
- Caching system for text extraction
- Performance metrics collection
- Health check endpoints

### Changed
- Optimized memory usage
- Improved error handling
- Enhanced logging

## [v1.1.0] - 2024-10-15

### Added
- Screenshot capture functionality
- Support for full page screenshots
- Custom screenshot quality settings
- Screenshot storage management

### Changed
- Enhanced browser control
- Improved error handling
- Better logging

## [v1.0.0] - 2024-09-01

### Added
- Initial release
- Basic text extraction
- Web page processing
- Simple API endpoints
- Basic error handling
