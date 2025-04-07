package extractors

import (
	"context"

	"github.com/ncecere/reader-go/internal/core/browser"
)

// TextExtractor extracts plain text from HTML pages
type TextExtractor struct {
	pool *browser.Pool
}

// NewTextExtractor creates a new text extractor
func NewTextExtractor(pool *browser.Pool) *TextExtractor {
	return &TextExtractor{pool: pool}
}

// ExtractText extracts text content from a URL
func (e *TextExtractor) ExtractText(ctx context.Context, url string) (string, error) {
	var extracted string
	err := e.pool.Execute(ctx, func(ctx context.Context) error {
		return browser.ExtractTextFromPage(ctx, url, &extracted)
	})
	if err != nil {
		return "", err
	}
	return extracted, nil
}
