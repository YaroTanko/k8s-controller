package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	LogLevel   string `mapstructure:"log_level"`
	KubeConfig string `mapstructure:"kubeconfig"`
	Namespace  string `mapstructure:"namespace"`
}

// LoadConfig initializes and loads configuration from environment variables and flags
func LoadConfig() (*Config, error) {
	// Set up Viper with defaults
	v := viper.New()

	// Set default values
	v.SetDefault("log_level", "info")
	v.SetDefault("kubeconfig", "")
	v.SetDefault("namespace", "")

	// Environment variables
	v.SetEnvPrefix("K8S_CONTROLLER")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Bind configuration to struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SetConfigValue allows setting a configuration value at runtime
func SetConfigValue(key string, value interface{}) {
	viper.Set(key, value)
}
