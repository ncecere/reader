package extractors

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/browser"
	"go.uber.org/zap"
)

// HTMLExtractor handles HTML content extraction
type HTMLExtractor struct {
	pool *browser.Pool
}

// NewHTMLExtractor creates a new HTML extractor
func NewHTMLExtractor(pool *browser.Pool) *HTMLExtractor {
	return &HTMLExtractor{pool: pool}
}

// ExtractHTML retrieves the HTML content from a URL
func (e *HTMLExtractor) ExtractHTML(ctx context.Context, url string) (string, error) {
	logger.Log.Info("Getting HTML content", zap.String("url", url))

	var html string
	err := e.pool.Execute(ctx, func(ctx context.Context) error {
		if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
			return fmt.Errorf("failed to navigate: %w", err)
		}

		if err := chromedp.Run(ctx, chromedp.WaitReady("body", chromedp.ByQuery)); err != nil {
			return fmt.Errorf("failed to wait for page load: %w", err)
		}

		if err := chromedp.Run(ctx, chromedp.OuterHTML("html", &html)); err != nil {
			return fmt.Errorf("failed to get HTML: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.Log.Error("Failed to get HTML",
			zap.String("url", url),
			zap.Error(err))
		return "", fmt.Errorf("HTML extraction failed: %w", err)
	}

	logger.Log.Info("Successfully retrieved HTML",
		zap.String("url", url),
		zap.Int("length", len(html)))

	return html, nil
}
