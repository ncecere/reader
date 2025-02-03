package main

import (
	"flag"
	"fmt"
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
	// Parse command line flags
	configPath := flag.String("config", "config.yml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Logging.Level, cfg.Logging.JSON, cfg.Logging.Caller); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Initialize browser service
	browserService, err := browser.NewService(&browser.BrowserOptions{
		PoolSize:   cfg.Browser.PoolSize,
		ChromePath: cfg.Browser.ChromePath,
		Timeout:    cfg.Browser.Timeout,
	})
	if err != nil {
		logger.Log.Fatal("Failed to initialize browser service", zap.Error(err))
	}
	defer browserService.Close()

	// Initialize and start server
	srv := server.NewServer(cfg, browserService)
	go func() {
		logger.Log.Info("Starting server",
			zap.String("port", cfg.Server.Port),
			zap.Int("pool_size", cfg.Browser.PoolSize),
			zap.String("chrome_path", cfg.Browser.ChromePath))
		if err := srv.Start(); err != nil {
			logger.Log.Fatal("Server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	logger.Log.Info("Shutting down server...")
	if err := srv.Stop(); err != nil {
		logger.Log.Error("Error during shutdown", zap.Error(err))
	}
}
