package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ncecere/reader-go/internal/common/config"
	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/ai"
	"github.com/ncecere/reader-go/internal/core/browser"
	"github.com/ncecere/reader-go/internal/core/service"
	"github.com/ncecere/reader-go/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the reader server",
	Long:  `Start the reader server with the specified configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create browser service
		browserService, err := service.NewService(&browser.BrowserOptions{
			PoolSize:   viper.GetInt("browser.pool_size"),
			ChromePath: viper.GetString("browser.chrome_path"),
			Timeout:    viper.GetInt("browser.timeout"),
		})
		if err != nil {
			logger.Log.Fatal("Failed to create browser service", zap.Error(err))
		}

		// Create config for AI service
		cfg := &config.Config{}
		cfg.AI.Enabled = viper.GetBool("ai.enabled")
		cfg.AI.APIEndpoint = viper.GetString("ai.api_endpoint")
		cfg.AI.APIKey = viper.GetString("ai.api_key")
		cfg.AI.Model = viper.GetString("ai.model")
		cfg.AI.Prompt = viper.GetString("ai.prompt")

		// Create AI service
		aiService := ai.NewService(cfg)

		// Create and start server
		srv := server.New(&server.Config{
			Port: viper.GetInt("server.port"),
		}, browserService, aiService)

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
			zap.Int("port", viper.GetInt("server.port")),
			zap.Int("pool_size", viper.GetInt("browser.pool_size")),
			zap.String("chrome_path", viper.GetString("browser.chrome_path")),
			zap.Bool("ai_enabled", viper.GetBool("ai.enabled")))

		return srv.Start()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
