package browser

import (
	"context"
	"sync"

	"github.com/ncecere/reader-go/internal/common/logger"
	"go.uber.org/zap"
)

// Result represents the result of processing a single URL
type Result struct {
	URL     string
	Content string
	Error   error
}

// ParallelProcessor handles concurrent URL processing
type ParallelProcessor struct {
	service *Service
	workers int
}

// NewParallelProcessor creates a new parallel processor
func NewParallelProcessor(service *Service, workers int) *ParallelProcessor {
	if workers <= 0 {
		workers = 3 // Default to 3 workers
	}
	return &ParallelProcessor{
		service: service,
		workers: workers,
	}
}

// ProcessURLs processes multiple URLs concurrently
func (p *ParallelProcessor) ProcessURLs(ctx context.Context, urls []string) []Result {
	results := make([]Result, len(urls))
	jobs := make(chan int, len(urls))
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				url := urls[idx]
				content, err := p.service.GetText(ctx, url)

				results[idx] = Result{
					URL:     url,
					Content: content,
					Error:   err,
				}

				if err != nil {
					logger.Log.Error("Failed to process URL",
						zap.String("url", url),
						zap.Error(err))
				} else {
					logger.Log.Info("Successfully processed URL",
						zap.String("url", url),
						zap.Int("content_length", len(content)))
				}
			}
		}()
	}

	// Send jobs to workers
	for i := range urls {
		jobs <- i
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()

	return results
}

// ProcessURLsBatch processes URLs in batches to control memory usage
func (p *ParallelProcessor) ProcessURLsBatch(ctx context.Context, urls []string, batchSize int) []Result {
	if batchSize <= 0 {
		batchSize = 10 // Default batch size
	}

	results := make([]Result, 0, len(urls))

	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}

		batch := urls[i:end]
		batchResults := p.ProcessURLs(ctx, batch)
		results = append(results, batchResults...)

		logger.Log.Info("Batch processed",
			zap.Int("batch_size", len(batch)),
			zap.Int("total_processed", len(results)),
			zap.Int("total_urls", len(urls)))
	}

	return results
}

// GetSuccessfulResults filters and returns only successful results
func (p *ParallelProcessor) GetSuccessfulResults(results []Result) []Result {
	successful := make([]Result, 0, len(results))
	for _, result := range results {
		if result.Error == nil {
			successful = append(successful, result)
		}
	}
	return successful
}

// GetFailedResults filters and returns only failed results
func (p *ParallelProcessor) GetFailedResults(results []Result) []Result {
	failed := make([]Result, 0)
	for _, result := range results {
		if result.Error != nil {
			failed = append(failed, result)
		}
	}
	return failed
}
