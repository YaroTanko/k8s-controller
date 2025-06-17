package middleware

import (
	"fmt"
	"k8s-controller/pkg/logger"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

// LoggingOptions configures how requests are logged
type LoggingOptions struct {
	// SkipPaths lists paths that should not be logged (e.g. health checks)
	SkipPaths []string
	// LogHeaders determines whether to log HTTP headers
	LogHeaders bool
	// LogRequestBody determines whether to log request body (for debugging)
	LogRequestBody bool
	// LogResponseBody determines whether to log response body (for debugging)
	LogResponseBody bool
	// MaxBodyLogSize limits the body size to log
	MaxBodyLogSize int
	// LogTiming enables detailed timing information
	LogTiming bool
}

// DefaultLoggingOptions returns default logging options
func DefaultLoggingOptions() *LoggingOptions {
	return &LoggingOptions{
		SkipPaths:       []string{"/health", "/metrics"},
		LogHeaders:      false,
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodyLogSize:  1024, // 1KB max for body logging
		LogTiming:       true,
	}
}

// EnhancedRequestLogger creates a middleware with configurable options
func EnhancedRequestLogger(options *LoggingOptions) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	if options == nil {
		options = DefaultLoggingOptions()
	}
	
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			path := string(ctx.Path())
			
			// Skip logging for specified paths
			for _, skipPath := range options.SkipPaths {
				if strings.HasPrefix(path, skipPath) {
					next(ctx)
					return
				}
			}
			
			// Record start time
			start := time.Now()
			
			// Create request ID for tracing
			requestID := fmt.Sprintf("%016X", start.UnixNano())
			
			// Extract request information
			method := string(ctx.Method())
			uri := string(ctx.RequestURI())
			clientIP := ctx.RemoteIP().String()
			
			// Log request start with basic info
			logEvent := logger.Debug().
				Str("request_id", requestID).
				Str("client_ip", clientIP).
				Str("method", method).
				Str("uri", uri)
			
			// Log headers if enabled
			if options.LogHeaders {
				headers := make(map[string]string)
				ctx.Request.Header.VisitAll(func(key, value []byte) {
					headers[string(key)] = string(value)
				})
				logEvent.Interface("headers", headers)
			}
			
			// Log request body if enabled and present
			if options.LogRequestBody && ctx.Request.Header.ContentLength() > 0 {
				body := string(ctx.Request.Body())
				if len(body) > options.MaxBodyLogSize {
					body = body[:options.MaxBodyLogSize] + "... (truncated)"
				}
				logEvent.Str("request_body", body)
			}
			
			logEvent.Msg("Request received")
			
			// Process request
			next(ctx)
			
			// Calculate duration
			duration := time.Since(start)
			
			// Extract response information
			statusCode := ctx.Response.StatusCode()
			responseSize := len(ctx.Response.Body())
			userAgent := string(ctx.UserAgent())
			
			// Determine log level based on status code
			var completeLogEvent *zerolog.Event
			if statusCode >= 500 {
				completeLogEvent = logger.Error()
			} else if statusCode >= 400 {
				completeLogEvent = logger.Warn()
			} else {
				completeLogEvent = logger.Info()
			}
			
			// Add common fields
			completeLogEvent.
				Str("request_id", requestID).
				Str("client_ip", clientIP).
				Str("method", method).
				Str("path", path).
				Int("status", statusCode).
				Int("size", responseSize).
				Str("user_agent", userAgent)
			
			// Add timing information if enabled
			if options.LogTiming {
				completeLogEvent.Dur("duration_ms", duration)
			}
			
			// Log response body if enabled
			if options.LogResponseBody {
				body := string(ctx.Response.Body())
				if len(body) > options.MaxBodyLogSize {
					body = body[:options.MaxBodyLogSize] + "... (truncated)"
				}
				completeLogEvent.Str("response_body", body)
			}
			
			completeLogEvent.Msg("Request completed")
		}
	}
}