package browser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ncecere/reader-go/internal/core/browser"
	"github.com/ncecere/reader-go/internal/core/extractors"
)

func TestHTMLExtractor_ExtractHTML(t *testing.T) {
	// Create a test server
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head><title>Test Page</title></head>
		<body>
			<h1>Test Content</h1>
			<p>This is a test paragraph.</p>
		</body>
		</html>
	`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(testHTML))
	}))
	defer ts.Close()

	// Create HTML extractor using test pool
	pool := browser.SetupTestPool(t)
	extractor := extractors.NewHTMLExtractor(pool)

	// Test cases
	tests := []struct {
		name       string
		url        string
		wantHTML   string
		wantErr    bool
		checkEmpty bool
	}{
		{
			name:     "Valid HTML page",
			url:      ts.URL,
			wantHTML: "Test Content",
			wantErr:  false,
		},
		{
			name:       "Invalid URL",
			url:        "http://invalid.url.that.does.not.exist",
			wantHTML:   "",
			wantErr:    true,
			checkEmpty: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			got, err := extractor.ExtractHTML(ctx, tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTMLExtractor.ExtractHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkEmpty {
				if got != "" {
					t.Errorf("HTMLExtractor.ExtractHTML() = %v, want empty string", got)
				}
				return
			}

			if !strings.Contains(got, tt.wantHTML) {
				t.Errorf("HTMLExtractor.ExtractHTML() = %v, want to contain %v", got, tt.wantHTML)
			}
		})
	}
}
