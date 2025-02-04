package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ncecere/reader-go/internal/common/config"
	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/browser"
	"github.com/ncecere/reader-go/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.Init(); err != nil {
		panic(err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Create browser service
	browserService, err := browser.NewService(&browser.BrowserOptions{
		PoolSize:   cfg.Browser.PoolSize,
		ChromePath: cfg.Browser.ChromePath,
		Timeout:    cfg.Browser.Timeout,
	})
	if err != nil {
		logger.Log.Fatal("Failed to create browser service", zap.Error(err))
	}

	// Create and start server
	srv := server.New(&server.Config{
		Port: cfg.Server.Port,
	}, browserService)

	// Handle shutdown gracefully
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Log.Info("Shutting down server...")
		browserService.Close()
		logger.Log.Info("Server shutdown complete")
		os.Exit(0)
	}()

	// Start server
	logger.Log.Info("Starting server",
		zap.Int("port", cfg.Server.Port),
		zap.Int("pool_size", cfg.Browser.PoolSize),
		zap.String("chrome_path", cfg.Browser.ChromePath))

	if err := srv.Start(); err != nil {
		logger.Log.Fatal("Server error", zap.Error(err))
	}
}
