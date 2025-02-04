package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/browser"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger.Log = log

	// Create browser service
	service, err := browser.NewService(browser.DefaultOptions())
	if err != nil {
		panic(fmt.Sprintf("Failed to create service: %v", err))
	}
	defer service.Close()

	url := "https://fastapi.tiangolo.com/"
	ctx := context.Background()

	// First request (cold)
	start := time.Now()
	content, err := service.GetText(ctx, url)
	if err != nil {
		panic(fmt.Sprintf("Failed to get text: %v", err))
	}
	fmt.Printf("First request took: %v\n", time.Since(start))
	fmt.Printf("Content length: %d bytes\n", len(content))
	fmt.Printf("\nFirst few lines:\n%s\n", content[:min(500, len(content))])

	// Second request (cached)
	start = time.Now()
	content, err = service.GetText(ctx, url)
	if err != nil {
		panic(fmt.Sprintf("Failed to get text: %v", err))
	}
	fmt.Printf("\nSecond request took: %v\n", time.Since(start))

	// Get metrics
	metrics := service.GetMetrics()
	fmt.Printf("\nMetrics:\n")
	fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("Cache Hit Rate: %.2f%%\n", metrics.CacheHitRate*100)
	fmt.Printf("Average Processing Time: %.2fms\n", metrics.AvgProcessingMs)
	fmt.Printf("Pool Utilization: %.2f%%\n", metrics.InstanceUtilPct)
	fmt.Printf("Current Memory: %dMB\n", metrics.CurrentMemoryMB)
	fmt.Printf("Peak Memory: %dMB\n", metrics.PeakMemoryMB)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
