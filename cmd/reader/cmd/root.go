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
	viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("browser.pool_size", rootCmd.PersistentFlags().Lookup("pool-size"))
	viper.BindPFlag("browser.chrome_path", rootCmd.PersistentFlags().Lookup("chrome-path"))
	viper.BindPFlag("browser.timeout", rootCmd.PersistentFlags().Lookup("browser-timeout"))
	viper.BindPFlag("browser.max_retries", rootCmd.PersistentFlags().Lookup("max-retries"))
	viper.BindPFlag("ai.enabled", rootCmd.PersistentFlags().Lookup("ai-enabled"))
	viper.BindPFlag("ai.api_endpoint", rootCmd.PersistentFlags().Lookup("ai-endpoint"))
	viper.BindPFlag("ai.api_key", rootCmd.PersistentFlags().Lookup("ai-key"))
	viper.BindPFlag("ai.model", rootCmd.PersistentFlags().Lookup("ai-model"))

	// Environment variables
	viper.SetEnvPrefix("READER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables
	viper.BindEnv("server.port", "READER_PORT")
	viper.BindEnv("browser.pool_size", "READER_POOL_SIZE")
	viper.BindEnv("browser.chrome_path", "READER_CHROME_PATH")
	viper.BindEnv("browser.timeout", "READER_BROWSER_TIMEOUT")
	viper.BindEnv("browser.max_retries", "READER_MAX_RETRIES")
	viper.BindEnv("ai.enabled", "READER_AI_ENABLED")
	viper.BindEnv("ai.api_endpoint", "READER_AI_ENDPOINT")
	viper.BindEnv("ai.api_key", "READER_AI_KEY")
	viper.BindEnv("ai.model", "READER_AI_MODEL")

	viper.AutomaticEnv()
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
