package browser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScreenshotCapture_CaptureScreenshot(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head><title>Test Page</title></head>
			<body>
				<h1>Test Content</h1>
				<p>This is a test paragraph.</p>
			</body>
			</html>
		`))
	}))
	defer ts.Close()

	// Create temporary directory for test screenshots
	tempDir, err := os.MkdirTemp("", "screenshot_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create screenshot capture using test pool
	pool := setupTestPool(t)
	screenshotOpts := &ScreenshotOptions{
		Quality:     90,
		StoragePath: tempDir,
	}
	capture := NewScreenshotCapture(pool, screenshotOpts)

	// Test cases
	tests := []struct {
		name     string
		url      string
		fullPage bool
		wantErr  bool
	}{
		{
			name:     "Valid page - viewport screenshot",
			url:      ts.URL,
			fullPage: false,
			wantErr:  false,
		},
		{
			name:     "Valid page - full page screenshot",
			url:      ts.URL,
			fullPage: true,
			wantErr:  false,
		},
		{
			name:     "Invalid URL",
			url:      "http://invalid.url.that.does.not.exist",
			fullPage: false,
			wantErr:  true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			outputPath := filepath.Join(tempDir, "test.png")
			capture.opts.FullPage = tt.fullPage

			err := capture.CaptureScreenshot(ctx, tt.url, outputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScreenshotCapture.CaptureScreenshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if file exists and has content
				info, err := os.Stat(outputPath)
				if err != nil {
					t.Errorf("Failed to stat screenshot file: %v", err)
					return
				}
				if info.Size() == 0 {
					t.Error("Screenshot file is empty")
				}
			}
		})
	}
}
