package middleware

import (
	"bytes"
	"k8s-controller/pkg/logger"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestEnhancedRequestLogger(t *testing.T) {
	// Set up a buffer to capture log output
	buffer := new(bytes.Buffer)
	logger.SetOutput(buffer)
	
	// Create test cases
	testCases := []struct {
		name          string
		path          string
		options       *LoggingOptions
		shouldBeLogged bool
		statusCode    int
	}{
		{
			name:          "Regular path",
			path:          "/api/users",
			options:       DefaultLoggingOptions(),
			shouldBeLogged: true,
			statusCode:    200,
		},
		{
			name:          "Health check path - should be skipped",
			path:          "/health",
			options:       DefaultLoggingOptions(),
			shouldBeLogged: false,
			statusCode:    200,
		},
		{
			name:          "Error path",
			path:          "/error",
			options:       DefaultLoggingOptions(),
			shouldBeLogged: true,
			statusCode:    500,
		},
		{
			name: "Custom options with headers",
			path: "/custom",
			options: &LoggingOptions{
				SkipPaths:       []string{"/skip"},
				LogHeaders:      true,
				LogRequestBody:  true,
				LogResponseBody: true,
				MaxBodyLogSize:  100,
				LogTiming:       true,
			},
			shouldBeLogged: true,
			statusCode:    200,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset buffer for each test
			buffer.Reset()
			
			// Create a test handler based on the path
			testHandler := func(ctx *fasthttp.RequestCtx) {
				// Special behavior for error path
				if string(ctx.Path()) == "/error" {
					ctx.SetStatusCode(500)
					ctx.SetBodyString("Internal Server Error")
				} else {
					ctx.SetStatusCode(tc.statusCode)
					ctx.SetBodyString("OK")
				}
			}
			
			// Wrap it with the enhanced request logger middleware
			handler := EnhancedRequestLogger(tc.options)(testHandler)
			
			// Set up a server
			s := &fasthttp.Server{
				Handler: handler,
			}
			
			// Create a test server with in-memory listener
			ln := fasthttputil.NewInmemoryListener()
			defer func() {
				if err := ln.Close(); err != nil {
					t.Errorf("Error closing listener: %v", err)
				}
			}()
			
			// Start server
			go func() {
				if err := s.Serve(ln); err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
			}()
			
			// Wait a moment for server to start
			time.Sleep(10 * time.Millisecond)
			
			// Create a client
			c := &fasthttp.Client{
				Dial: func(addr string) (net.Conn, error) {
					return ln.Dial()
				},
			}
			
			// Create a request
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)
			
			// Set the request URI and method
			req.SetRequestURI("http://localhost" + tc.path)
			req.Header.SetMethod("GET")
			
			// Set headers and body for testing those features
			if tc.options != nil && tc.options.LogHeaders {
				req.Header.Set("X-Test-Header", "test-value")
			}
			
			if tc.options != nil && tc.options.LogRequestBody {
				req.SetBodyString("Test request body")
			}
			
			// Create a response
			resp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseResponse(resp)
			
			// Send request
			if err := c.Do(req, resp); err != nil {
				t.Fatalf("Error sending request: %s", err)
			}
			
			// Verify status code
			if resp.StatusCode() != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, resp.StatusCode())
			}
			
			// Verify logging behavior
			logOutput := buffer.String()
			
			if tc.shouldBeLogged {
				// Check basic logging
				if !strings.Contains(logOutput, "Request received") {
					t.Error("Expected log to contain 'Request received'")
				}
				
				if !strings.Contains(logOutput, "Request completed") {
					t.Error("Expected log to contain 'Request completed'")
				}
				
				// Check for headers if enabled
				if tc.options != nil && tc.options.LogHeaders {
					if !strings.Contains(logOutput, "X-Test-Header") {
						t.Error("Expected log to contain request headers")
					}
				}
				
				// Check for request body if enabled
				if tc.options != nil && tc.options.LogRequestBody {
					if !strings.Contains(logOutput, "Test request body") {
						t.Error("Expected log to contain request body")
					}
				}
				
				// For error paths, check the log level
				if tc.statusCode >= 500 {
					if !strings.Contains(strings.ToUpper(logOutput), "ERROR") {
						t.Error("Expected log to use ERROR level for 5xx status codes")
					}
				}
			} else {
				// For skipped paths, verify no logging happened
				if strings.Contains(logOutput, "Request received") || strings.Contains(logOutput, "Request completed") {
					t.Error("Expected request to be skipped from logging")
				}
			}
		})
	}
}