package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	// Test initialization with different log levels
	levels := []LogLevel{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, TraceLevel}
	
	for _, level := range levels {
		Init(level)
		// Not much to verify besides that it doesn't panic
		// More extensive testing would require checking the output
	}
}

func TestSetOutput(t *testing.T) {
	buffer := new(bytes.Buffer)
	SetOutput(buffer)
	
	// Log something
	Info().Msg("test message")
	
	// Verify the output contains our message
	if !strings.Contains(buffer.String(), "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", buffer.String())
	}
}

func TestLogLevels(t *testing.T) {
	// Set up a buffer to capture output
	buffer := new(bytes.Buffer)
	SetOutput(buffer)
	
	// Test all log levels
	Debug().Msg("debug message")
	Info().Msg("info message")
	Warn().Msg("warn message")
	Error().Msg("error message")
	Trace().Msg("trace message")
	
	// Check that all messages were logged
	output := buffer.String()
	
	messages := []string{"debug message", "info message", "warn message", "error message", "trace message"}
	for _, msg := range messages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected log output to contain '%s'", msg)
		}
	}
}