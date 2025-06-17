package config

import (
	"os"
	"testing"
	
	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	// Test with environment variables
	os.Setenv("K8S_CONTROLLER_LOG_LEVEL", "debug")
	os.Setenv("K8S_CONTROLLER_NAMESPACE", "test-namespace")
	
	// Load config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify the config is loaded
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel to be 'debug', got %s", cfg.LogLevel)
	}
	
	if cfg.Namespace != "test-namespace" {
		t.Errorf("Expected Namespace to be 'test-namespace', got %s", cfg.Namespace)
	}
	
	// Clean up
	os.Unsetenv("K8S_CONTROLLER_LOG_LEVEL")
	os.Unsetenv("K8S_CONTROLLER_NAMESPACE")
}

func TestLoadConfigDefaults(t *testing.T) {
	// Load config with defaults
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify defaults are set
	if cfg.LogLevel != "info" {
		t.Errorf("Expected default LogLevel to be 'info', got %s", cfg.LogLevel)
	}
	
	if cfg.Namespace != "" {
		t.Errorf("Expected default Namespace to be empty, got %s", cfg.Namespace)
	}
}

func TestSetConfigValue(t *testing.T) {
	// Set a config value
	SetConfigValue("log_level", "trace")
	
	// Verify the value was set
	if val := viper.GetString("log_level"); val != "trace" {
		t.Errorf("Expected log_level to be 'trace', got %s", val)
	}
}