package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "reader",
		Short: "A web page reader and summarizer",
		Long: `Reader is a service that extracts text content from web pages
and provides summarization capabilities using AI.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")

	// Server flags
	rootCmd.PersistentFlags().Int("port", 4444, "Port to run the server on")

	// Browser flags
	rootCmd.PersistentFlags().Int("pool-size", 3, "Number of browser instances in the pool")
	rootCmd.PersistentFlags().String("chrome-path", "", "Path to Chrome/Chromium executable")
	rootCmd.PersistentFlags().Int("browser-timeout", 30, "Browser request timeout in seconds")
	rootCmd.PersistentFlags().Int("max-retries", 3, "Maximum number of retries for browser operations")

	// AI flags
	rootCmd.PersistentFlags().Bool("ai-enabled", true, "Enable/disable AI features")
	rootCmd.PersistentFlags().String("ai-endpoint", "", "AI API endpoint")
	rootCmd.PersistentFlags().String("ai-key", "", "AI API key")
	rootCmd.PersistentFlags().String("ai-model", "vltr-mistral-small", "AI model to use")

	// Bind flags to viper
	bindFlags()

	// Environment variables
	viper.SetEnvPrefix("READER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables
	bindEnvs()

	viper.AutomaticEnv()
}

// bindFlags binds command line flags to viper configuration
func bindFlags() {
	flags := map[string]string{
		"server.port":         "port",
		"browser.pool_size":   "pool-size",
		"browser.chrome_path": "chrome-path",
		"browser.timeout":     "browser-timeout",
		"browser.max_retries": "max-retries",
		"ai.enabled":          "ai-enabled",
		"ai.api_endpoint":     "ai-endpoint",
		"ai.api_key":          "ai-key",
		"ai.model":            "ai-model",
	}

	for configKey, flagName := range flags {
		if err := viper.BindPFlag(configKey, rootCmd.PersistentFlags().Lookup(flagName)); err != nil {
			logger.Log.Error("Failed to bind flag",
				zap.String("flag", flagName),
				zap.Error(err))
		}
	}
}

// bindEnvs binds environment variables to viper configuration
func bindEnvs() {
	envs := map[string]string{
		"server.port":         "READER_PORT",
		"browser.pool_size":   "READER_POOL_SIZE",
		"browser.chrome_path": "READER_CHROME_PATH",
		"browser.timeout":     "READER_BROWSER_TIMEOUT",
		"browser.max_retries": "READER_MAX_RETRIES",
		"ai.enabled":          "READER_AI_ENABLED",
		"ai.api_endpoint":     "READER_AI_ENDPOINT",
		"ai.api_key":          "READER_AI_KEY",
		"ai.model":            "READER_AI_MODEL",
	}

	for configKey, envVar := range envs {
		if err := viper.BindEnv(configKey, envVar); err != nil {
			logger.Log.Error("Failed to bind environment variable",
				zap.String("env", envVar),
				zap.Error(err))
		}
	}
}

func initConfig() {
	if err := logger.Init(); err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Log.Info("Using config file", zap.String("file", viper.ConfigFileUsed()))
	} else {
		logger.Log.Info("No config file found, using defaults and environment variables")
	}

	// Debug: Print all settings
	logger.Log.Info("Configuration",
		zap.Int("port", viper.GetInt("server.port")),
		zap.Int("pool_size", viper.GetInt("browser.pool_size")),
		zap.String("chrome_path", viper.GetString("browser.chrome_path")),
		zap.Int("timeout", viper.GetInt("browser.timeout")),
		zap.Bool("ai_enabled", viper.GetBool("ai.enabled")))
}
