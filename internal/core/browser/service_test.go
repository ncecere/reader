package browser

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServicePerformance(t *testing.T) {
	service := setupTestService(t)

	// Test URLs - using http to avoid SSL issues in tests
	urls := []string{
		"http://example.com",
		"http://example.org",
		"http://example.net",
	}

	// Test parallel processing
	t.Run("Parallel Processing", func(t *testing.T) {
		ctx := context.Background()
		results := service.ProcessURLs(ctx, urls)
		assert.Equal(t, len(urls), len(results))
	})

	// Test batch processing
	t.Run("Batch Processing", func(t *testing.T) {
		ctx := context.Background()
		results := service.ProcessURLsBatch(ctx, urls, 2)
		assert.Equal(t, len(urls), len(results))
	})

	// Test caching
	t.Run("Cache Hit", func(t *testing.T) {
		ctx := context.Background()
		url := "http://example.com"

		// First request - should miss cache
		content1, err := service.GetText(ctx, url)
		require.NoError(t, err)
		require.NotEmpty(t, content1)

		// Second request - should hit cache
		start := time.Now()
		content2, err := service.GetText(ctx, url)
		require.NoError(t, err)
		assert.Equal(t, content1, content2)
		assert.Less(t, time.Since(start), 100*time.Millisecond) // Cache hit should be fast
	})

	// Test metrics
	t.Run("Metrics Collection", func(t *testing.T) {
		metrics := service.GetMetrics()
		assert.NotZero(t, metrics.TotalRequests)
		assert.NotZero(t, metrics.SuccessfulReqs)
	})

	// Test pool utilization
	t.Run("Pool Utilization", func(t *testing.T) {
		ctx := context.Background()
		url := "http://example.com"

		// Make concurrent requests up to pool size
		done := make(chan struct{})
		for i := 0; i < 2; i++ { // Using pool size from testing.go
			go func() {
				defer func() { done <- struct{}{} }()
				_, err := service.GetText(ctx, url)
				require.NoError(t, err)
			}()
		}

		// Wait for all requests
		for i := 0; i < 2; i++ {
			<-done
		}

		// Check metrics
		metrics := service.GetMetrics()
		assert.NotZero(t, metrics.InstanceUtilPct)
	})

	// Test error handling
	t.Run("Error Handling", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetText(ctx, "invalid-url")
		assert.Error(t, err)

		metrics := service.GetMetrics()
		assert.NotZero(t, metrics.FailedReqs)
	})

	// Test memory management
	t.Run("Memory Management", func(t *testing.T) {
		metrics := service.GetMetrics()
		assert.NotZero(t, metrics.CurrentMemoryMB)
		assert.NotZero(t, metrics.PeakMemoryMB)
		assert.GreaterOrEqual(t, metrics.PeakMemoryMB, metrics.CurrentMemoryMB)
	})
}

func TestServiceConcurrency(t *testing.T) {
	service := setupTestService(t)
	ctx := context.Background()

	// Generate test URLs
	urls := make([]string, 20)
	for i := range urls {
		urls[i] = "http://example.com"
	}

	// Test different batch sizes
	batchSizes := []int{1, 5, 10}
	for _, size := range batchSizes {
		t.Run("BatchSize_"+string(rune(size)), func(t *testing.T) {
			start := time.Now()
			results := service.ProcessURLsBatch(ctx, urls, size)
			duration := time.Since(start)

			assert.Equal(t, len(urls), len(results))
			t.Logf("Batch size %d took %v", size, duration)
		})
	}
}

func BenchmarkService(b *testing.B) {
	service := setupTestService(b)
	ctx := context.Background()
	url := "http://example.com"

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := service.GetText(ctx, url)
			require.NoError(b, err)
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		urls := make([]string, b.N)
		for i := range urls {
			urls[i] = url
		}
		results := service.ProcessURLs(ctx, urls)
		for _, result := range results {
			require.NoError(b, result.Error)
		}
	})

	b.Run("Batched", func(b *testing.B) {
		urls := make([]string, b.N)
		for i := range urls {
			urls[i] = url
		}
		results := service.ProcessURLsBatch(ctx, urls, 10)
		for _, result := range results {
			require.NoError(b, result.Error)
		}
	})
}
