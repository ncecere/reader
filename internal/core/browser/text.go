package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ncecere/reader-go/internal/common/logger"
	"go.uber.org/zap"
)

// TextExtractor handles text extraction operations
type TextExtractor struct {
	pool *Pool
}

// NewTextExtractor creates a new text extractor
func NewTextExtractor(pool *Pool) *TextExtractor {
	return &TextExtractor{pool: pool}
}

// ExtractText extracts text content from a URL
func (e *TextExtractor) ExtractText(ctx context.Context, url string) (string, error) {
	logger.Log.Info("Getting text content", zap.String("url", url))

	var text string
	err := e.pool.Execute(ctx, func(ctx context.Context) error {
		// Extract text with retry
		var err error
		for i := 0; i < 5; i++ {
			timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			err = chromedp.Run(timeoutCtx,
				chromedp.Navigate(url),
				chromedp.WaitReady("body", chromedp.ByQuery),
				// Wait for network idle to ensure dynamic content is loaded
				chromedp.ActionFunc(func(ctx context.Context) error {
					return chromedp.Evaluate(`
						new Promise((resolve) => {
							if (document.readyState === 'complete') {
								resolve();
							} else {
								window.addEventListener('load', resolve);
							}
						})
					`, nil).Do(ctx)
				}),
				chromedp.Text("body", &text, chromedp.NodeVisible, chromedp.ByQuery),
			)
			if err == nil {
				break
			}
			logger.Log.Warn("Retrying text extraction",
				zap.String("url", url),
				zap.Int("attempt", i+1),
				zap.Error(err))
			time.Sleep(time.Second)
		}
		if err != nil {
			return fmt.Errorf("failed to extract text: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.Log.Error("Failed to get text",
			zap.String("url", url),
			zap.Error(err))
		return "", fmt.Errorf("text extraction failed: %w", err)
	}

	logger.Log.Info("Successfully retrieved text",
		zap.String("url", url),
		zap.Int("length", len(text)))

	return text, nil
}
