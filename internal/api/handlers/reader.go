package handlers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/common/metrics"
	"github.com/ncecere/reader-go/internal/core/browser"
	"github.com/ncecere/reader-go/internal/core/converter"
	"go.uber.org/zap"
)

// ReaderHandler handles web content extraction requests
type ReaderHandler struct {
	browser *browser.Service
}

// NewReaderHandler creates a new reader handler
func NewReaderHandler(browser *browser.Service) *ReaderHandler {
	return &ReaderHandler{
		browser: browser,
	}
}

// HandleRequest processes URL requests and returns content in requested format
func (h *ReaderHandler) HandleRequest(c *fiber.Ctx) error {
	url := strings.TrimPrefix(c.Path(), "/")
	format := c.Get("X-Respond-With", "text")

	// Start timing
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.ContentProcessingDuration.WithLabelValues(format).Observe(duration)
	}()

	var (
		content string
		err     error
	)

	switch format {
	case "text":
		content, err = h.browser.GetText(c.Context(), url)
		if err != nil {
			logger.Log.Error("Failed to extract text",
				zap.String("url", url),
				zap.Error(err))
			metrics.ContentProcessingErrors.WithLabelValues(format, "extraction_failed").Inc()
			return c.SendString("Failed to extract text")
		}

	case "markdown":
		html, err := h.browser.GetHTML(c.Context(), url)
		if err != nil {
			logger.Log.Error("Failed to get HTML",
				zap.String("url", url),
				zap.Error(err))
			metrics.ContentProcessingErrors.WithLabelValues(format, "html_extraction_failed").Inc()
			return c.SendString("Failed to get HTML")
		}

		content, err = converter.HTMLToMarkdown(html)
		if err != nil {
			logger.Log.Error("Failed to convert to markdown",
				zap.String("url", url),
				zap.Error(err))
			metrics.ContentProcessingErrors.WithLabelValues(format, "conversion_failed").Inc()
			return c.SendString("Failed to convert to markdown")
		}

	case "screenshot", "pageshot":
		outputPath := filepath.Join("screenshots", fmt.Sprintf("%d.png", time.Now().Unix()))
		err = h.browser.CaptureScreenshot(c.Context(), url, outputPath, format == "pageshot")
		if err != nil {
			logger.Log.Error("Failed to capture screenshot",
				zap.String("url", url),
				zap.Error(err))
			metrics.ContentProcessingErrors.WithLabelValues(format, "capture_failed").Inc()
			return c.SendString("Failed to capture screenshot")
		}
		return c.SendFile(outputPath)

	default:
		return c.Status(400).SendString("Invalid format")
	}

	// Record content size
	metrics.ContentSize.WithLabelValues(format).Observe(float64(len(content)))

	// Record URL processing
	domain := extractDomain(url)
	metrics.URLProcessing.WithLabelValues(domain).Inc()
	metrics.URLContentTypes.WithLabelValues(format).Inc()
	metrics.URLSizes.WithLabelValues(domain).Observe(float64(len(content)))

	return c.SendString(content)
}

// extractDomain extracts the domain from a URL
func extractDomain(url string) string {
	// Simple domain extraction
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}
