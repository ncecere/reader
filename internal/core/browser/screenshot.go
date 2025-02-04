package browser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ncecere/reader-go/internal/common/logger"
	"go.uber.org/zap"
)

// ScreenshotOptions configures screenshot capture
type ScreenshotOptions struct {
	Quality     int
	StoragePath string
	FullPage    bool
}

// ScreenshotCapture handles webpage screenshot operations
type ScreenshotCapture struct {
	pool *Pool
	opts *ScreenshotOptions
}

// NewScreenshotCapture creates a new screenshot capture instance
func NewScreenshotCapture(pool *Pool, opts *ScreenshotOptions) *ScreenshotCapture {
	if opts == nil {
		opts = &ScreenshotOptions{
			Quality:     90,
			StoragePath: "screenshots",
			FullPage:    false,
		}
	}
	return &ScreenshotCapture{
		pool: pool,
		opts: opts,
	}
}

// CaptureScreenshot takes a screenshot of a webpage
func (s *ScreenshotCapture) CaptureScreenshot(ctx context.Context, url string, outputPath string) error {
	logger.Log.Info("Capturing screenshot",
		zap.String("url", url),
		zap.String("output", outputPath))

	var buf []byte
	err := s.pool.Execute(ctx, func(ctx context.Context) error {
		// Navigate to page
		if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
			return fmt.Errorf("failed to navigate: %w", err)
		}

		// Wait for page to load with enhanced conditions
		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.WaitReady("body", chromedp.ByQuery),
			chromedp.Sleep(1 * time.Second), // Give time for dynamic content to load
		}); err != nil {
			return fmt.Errorf("failed to wait for page load: %w", err)
		}

		// Set viewport size for consistent screenshots
		if err := chromedp.Run(ctx, chromedp.EmulateViewport(1280, 800)); err != nil {
			return fmt.Errorf("failed to set viewport: %w", err)
		}

		// Capture screenshot
		if s.opts.FullPage {
			if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, s.opts.Quality)); err != nil {
				return fmt.Errorf("failed to capture full screenshot: %w", err)
			}
		} else {
			if err := chromedp.Run(ctx, chromedp.Screenshot("body", &buf, chromedp.NodeVisible)); err != nil {
				return fmt.Errorf("failed to capture screenshot: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("screenshot capture failed: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write screenshot to file
	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return fmt.Errorf("failed to write screenshot: %w", err)
	}

	logger.Log.Info("Screenshot captured successfully",
		zap.String("url", url),
		zap.String("output", outputPath),
		zap.Int("size", len(buf)))

	return nil
}
