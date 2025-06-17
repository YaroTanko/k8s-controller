package middleware

import (
	"fmt"
	"k8s-controller/pkg/logger"
	"time"

	"github.com/valyala/fasthttp"
)

// RequestLogger is a middleware that logs HTTP requests with detailed information
func RequestLogger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Record start time
		start := time.Now()
		
		// Create request ID for tracing
		requestID := fmt.Sprintf("%016X", start.UnixNano())
		
		// Get client IP address
		clientIP := ctx.RemoteIP().String()
		
		// Log request details before handling
		logger.Debug().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("method", string(ctx.Method())).
			Str("uri", string(ctx.RequestURI())).
			Msg("Request received")
		
		// Process request
		next(ctx)
		
		// Calculate duration
		duration := time.Since(start)
		
		// Extract request and response information
		path := string(ctx.Path())
		method := string(ctx.Method())
		statusCode := ctx.Response.StatusCode()
		responseSize := len(ctx.Response.Body())
		userAgent := string(ctx.UserAgent())
		
		// Log the completed request
		logEvent := logger.Info()
		
		// For errors (4xx, 5xx), use higher log level
		if statusCode >= 400 && statusCode < 500 {
			logEvent = logger.Warn()
		} else if statusCode >= 500 {
			logEvent = logger.Error()
		}
		
		logEvent.
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Int("size", responseSize).
			Str("user_agent", userAgent).
			Dur("duration", duration).
			Msg("HTTP request completed")
	}
}