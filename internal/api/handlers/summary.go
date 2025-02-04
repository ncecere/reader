package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/common/metrics"
	"github.com/ncecere/reader-go/internal/core/ai"
	"github.com/ncecere/reader-go/internal/core/browser"
	"github.com/ncecere/reader-go/internal/core/converter"
	"go.uber.org/zap"
)

// SummaryHandler handles web content summarization requests
type SummaryHandler struct {
	browser *browser.Service
	ai      *ai.Service
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(browser *browser.Service, ai *ai.Service) *SummaryHandler {
	return &SummaryHandler{
		browser: browser,
		ai:      ai,
	}
}

// HandleRequest processes URL requests and returns summarized content
func (h *SummaryHandler) HandleRequest(c *fiber.Ctx) error {
	url := strings.TrimPrefix(c.Path(), "/summary/")
	format := c.Get("X-Respond-With", "text")

	// Start timing
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.ContentProcessingDuration.WithLabelValues("summary").Observe(duration)
	}()

	// Get the text content first
	text, err := h.browser.GetText(c.Context(), url)
	if err != nil {
		logger.Log.Error("Failed to extract text for summary",
			zap.String("url", url),
			zap.Error(err))
		metrics.ContentProcessingErrors.WithLabelValues("summary", "extraction_failed").Inc()
		return c.SendString("Failed to extract text for summary")
	}

	// Generate summary
	summary, err := h.ai.Summarize(c.Context(), text)
	if err != nil {
		logger.Log.Error("Failed to generate summary",
			zap.String("url", url),
			zap.Error(err))
		metrics.ContentProcessingErrors.WithLabelValues("summary", "summarization_failed").Inc()
		return c.SendString("Failed to generate summary")
	}

	var content string
	switch format {
	case "markdown":
		content, err = converter.TextToMarkdown(summary)
		if err != nil {
			logger.Log.Error("Failed to convert summary to markdown",
				zap.String("url", url),
				zap.Error(err))
			metrics.ContentProcessingErrors.WithLabelValues("summary", "conversion_failed").Inc()
			return c.SendString("Failed to convert summary to markdown")
		}
	case "text":
		content = summary
	default:
		return c.Status(400).SendString("Invalid format")
	}

	// Record metrics
	metrics.ContentSize.WithLabelValues("summary").Observe(float64(len(content)))
	domain := extractDomain(url)
	metrics.URLProcessing.WithLabelValues(domain).Inc()
	metrics.URLContentTypes.WithLabelValues("summary").Inc()
	metrics.URLSizes.WithLabelValues(domain).Observe(float64(len(content)))

	return c.SendString(content)
}
