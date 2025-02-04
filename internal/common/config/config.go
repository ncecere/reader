package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port    int `yaml:"port"`
		Timeout int `yaml:"timeout"`
	} `yaml:"server"`

	Browser struct {
		PoolSize   int    `yaml:"poolSize"`
		ChromePath string `yaml:"chromePath"`
		Timeout    int    `yaml:"timeout"`
		MaxRetries int    `yaml:"maxRetries"`
	} `yaml:"browser"`

	Screenshots struct {
		StoragePath string `yaml:"storagePath"`
		Quality     int    `yaml:"quality"`
		DefaultType string `yaml:"defaultType"`
	} `yaml:"screenshots"`

	Logging struct {
		Level  string `yaml:"level"`
		JSON   bool   `yaml:"json"`
		Caller bool   `yaml:"caller"`
	} `yaml:"logging"`
}

// Load loads configuration from a YAML file
func Load() (*Config, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Set defaults
	if config.Server.Port == 0 {
		config.Server.Port = 4000
	}
	if config.Server.Timeout == 0 {
		config.Server.Timeout = 60
	}
	if config.Browser.PoolSize == 0 {
		config.Browser.PoolSize = 2
	}
	if config.Browser.Timeout == 0 {
		config.Browser.Timeout = 30
	}
	if config.Screenshots.Quality == 0 {
		config.Screenshots.Quality = 90
	}
	if config.Screenshots.StoragePath == "" {
		config.Screenshots.StoragePath = "screenshots"
	}
	if config.Screenshots.DefaultType == "" {
		config.Screenshots.DefaultType = "viewport"
	}

	return &config, nil
}
