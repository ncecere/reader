package browser

import (
	"os"
	"testing"
	"time"

	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/metrics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	testTimeout = 60 * time.Second
	testRetries = 3
)

func setupTestLogger(tb testing.TB) func() {
	// Create a test logger configuration
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	config.OutputPaths = []string{os.DevNull} // Suppress output during tests

	// Create and set the logger
	log, err := config.Build()
	if err != nil {
		tb.Fatalf("Failed to create test logger: %v", err)
	}

	// Set the global logger
	logger.Log = log

	// Return cleanup function
	return func() {
		_ = log.Sync()
	}
}

func SetupTestPool(tb testing.TB) *Pool {
	cleanup := setupTestLogger(tb)
	tb.Cleanup(cleanup)

	opts := &BrowserOptions{
		PoolSize: 1,
		Timeout:  int(testTimeout.Seconds()),
	}

	testMetrics := metrics.New()
	pool, err := NewPool(opts, testMetrics)
	if err != nil {
		tb.Fatalf("Failed to create browser pool: %v", err)
	}

	tb.Cleanup(pool.Close)
	return pool
}
