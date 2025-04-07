package service

import (
	"os"
	"testing"
	"time"

	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/browser"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	testTimeout = 60 * time.Second
	testRetries = 3
)

// SetupTestService initializes a Service instance for testing
func SetupTestService(tb testing.TB) *Service {
	// Setup logger
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	config.OutputPaths = []string{os.DevNull}

	log, err := config.Build()
	if err != nil {
		tb.Fatalf("Failed to create test logger: %v", err)
	}
	logger.Log = log
	tb.Cleanup(func() { _ = log.Sync() })

	opts := &browser.BrowserOptions{
		PoolSize: 2,
		Timeout:  int(testTimeout.Seconds()),
	}

	svc, err := NewService(opts)
	if err != nil {
		tb.Fatalf("Failed to create service: %v", err)
	}

	tb.Cleanup(func() {
		svc.Close()
	})

	return svc
}
