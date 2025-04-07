package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ncecere/reader-go/internal/common/logger"
	"go.uber.org/zap"
)

// ExtractTextFromPage navigates to the URL and extracts visible text from the page body.
// It retries on failure up to 5 times.
func ExtractTextFromPage(ctx context.Context, url string, out *string) error {
	logger.Log.Info("Getting text content", zap.String("url", url))

	var lastErr error
	for i := 0; i < 5; i++ {
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err := chromedp.Run(timeoutCtx,
			chromedp.Navigate(url),
			chromedp.WaitReady("body", chromedp.ByQuery),
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
			chromedp.Text("body", out, chromedp.NodeVisible, chromedp.ByQuery),
		)
		if err == nil {
			logger.Log.Info("Successfully retrieved text",
				zap.String("url", url),
				zap.Int("length", len(*out)))
			return nil
		}

		lastErr = err
		logger.Log.Warn("Retrying text extraction",
			zap.String("url", url),
			zap.Int("attempt", i+1),
			zap.Error(err))
		time.Sleep(time.Second)
	}

	logger.Log.Error("Failed to get text",
		zap.String("url", url),
		zap.Error(lastErr))
	return fmt.Errorf("text extraction failed: %w", lastErr)
}
