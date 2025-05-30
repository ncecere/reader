package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/browser"
	"github.com/ncecere/reader-go/internal/core/cache"
	"github.com/ncecere/reader-go/internal/core/extractors"
	cachex "github.com/ncecere/reader-go/internal/core/extractors/cache"
	"github.com/ncecere/reader-go/internal/core/metrics"
	"github.com/ncecere/reader-go/internal/core/parallel"
	"go.uber.org/zap"
)

// Service manages browser operations through specialized components
type Service struct {
	pool     *browser.Pool
	text     *cachex.CachedTextExtractor
	html     *extractors.HTMLExtractor
	parallel *parallel.ParallelProcessor
	metrics  *metrics.Metrics
}

// NewService creates a new browser service
func NewService(opts *browser.BrowserOptions) (*Service, error) {
	if opts == nil {
		opts = browser.DefaultOptions()
	}

	logger.Log.Info("Creating browser pool",
		zap.Int("pool_size", opts.PoolSize),
		zap.String("chrome_path", opts.ChromePath),
		zap.Int("max_memory_mb", opts.MaxMemoryMB))

	serviceMetrics := metrics.New()

	pool, err := browser.NewPool(opts, serviceMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create browser pool: %w", err)
	}

	textExtractor := extractors.NewTextExtractor(pool)
	cachedTextExtractor := cachex.NewCachedTextExtractor(textExtractor, &cache.Options{
		MaxAge:   1 * time.Hour,
		MaxItems: 1000,
	}, serviceMetrics)

	svc := &Service{
		pool:    pool,
		text:    cachedTextExtractor,
		html:    extractors.NewHTMLExtractor(pool),
		metrics: serviceMetrics,
	}

	svc.parallel = parallel.NewParallelProcessor(svc, opts.PoolSize)

	return svc, nil
}

// GetText extracts text content from a URL
func (s *Service) GetText(ctx context.Context, url string) (string, error) {
	start := time.Now()
	content, err := s.text.ExtractText(ctx, url)
	s.metrics.RecordRequest(time.Since(start), err == nil)
	return content, err
}

// GetHTML retrieves the HTML content from a URL
func (s *Service) GetHTML(ctx context.Context, url string) (string, error) {
	return s.html.ExtractHTML(ctx, url)
}

// ProcessURLs processes multiple URLs in parallel
func (s *Service) ProcessURLs(ctx context.Context, urls []string) []parallel.Result {
	return s.parallel.ProcessURLs(ctx, urls)
}

// ProcessURLsBatch processes URLs in batches
func (s *Service) ProcessURLsBatch(ctx context.Context, urls []string, batchSize int) []parallel.Result {
	return s.parallel.ProcessURLsBatch(ctx, urls, batchSize)
}

// GetMetrics returns current service metrics
func (s *Service) GetMetrics() metrics.Stats {
	return s.metrics.GetStats()
}

// Close closes all browser instances
func (s *Service) Close() error {
	s.pool.Close()
	return nil
}
