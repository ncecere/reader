package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/cache"
	"github.com/ncecere/reader-go/internal/core/extractors"
	"github.com/ncecere/reader-go/internal/core/metrics"
	"go.uber.org/zap"
)

// CachedTextExtractor adds caching to text extraction
type CachedTextExtractor struct {
	extractor *extractors.TextExtractor
	cache     *cache.Cache
	metrics   *metrics.Metrics
}

// NewCachedTextExtractor creates a new cached text extractor
func NewCachedTextExtractor(extractor *extractors.TextExtractor, opts *cache.Options, metrics *metrics.Metrics) *CachedTextExtractor {
	if opts == nil {
		opts = &cache.Options{
			MaxAge:   1 * time.Hour, // Cache content for 1 hour by default
			MaxItems: 1000,          // Store up to 1000 items
		}
	}

	return &CachedTextExtractor{
		extractor: extractor,
		cache:     cache.New(opts),
		metrics:   metrics,
	}
}

// ExtractText attempts to get text from cache before falling back to actual extraction
func (e *CachedTextExtractor) ExtractText(ctx context.Context, url string) (string, error) {
	key := e.generateKey(url)

	if content, found := e.cache.Get(key); found {
		logger.Log.Info("Cache hit",
			zap.String("url", url),
			zap.Int("content_length", len(content)))
		e.metrics.RecordCacheAccess(true)
		return content, nil
	}

	logger.Log.Info("Cache miss, extracting text", zap.String("url", url))
	e.metrics.RecordCacheAccess(false)
	content, err := e.extractor.ExtractText(ctx, url)
	if err != nil {
		return "", fmt.Errorf("text extraction failed: %w", err)
	}

	e.cache.Set(key, content)

	logger.Log.Info("Cached extracted text",
		zap.String("url", url),
		zap.Int("content_length", len(content)))

	return content, nil
}

func (e *CachedTextExtractor) generateKey(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}

// GetStats returns cache statistics
func (e *CachedTextExtractor) GetStats() cache.Stats {
	return e.cache.Stats()
}

// ClearCache removes all entries from the cache
func (e *CachedTextExtractor) ClearCache() {
	e.cache.Clear()
}
