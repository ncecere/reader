package browser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTextExtractor_ExtractText(t *testing.T) {
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

	// Create text extractor using test pool
	pool := setupTestPool(t)
	extractor := NewTextExtractor(pool)

	// Test cases
	tests := []struct {
		name    string
		url     string
		want    []string
		wantErr bool
	}{
		{
			name: "Valid HTML page",
			url:  ts.URL,
			want: []string{
				"Test Content",
				"This is a test paragraph",
			},
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "http://invalid.url.that.does.not.exist",
			want:    nil,
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			got, err := extractor.ExtractText(ctx, tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("TextExtractor.ExtractText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check each expected string is present
				for _, want := range tt.want {
					if !strings.Contains(got, want) {
						t.Errorf("TextExtractor.ExtractText() = %v, want to contain %v", got, want)
					}
				}
			}
		})
	}
}
