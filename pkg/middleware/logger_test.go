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

func TestRequestLogger(t *testing.T) {
	// Set up a buffer to capture log output
	buffer := new(bytes.Buffer)
	logger.SetOutput(buffer)
	
	// Create a test handler
	testHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(200)
		ctx.SetBodyString("OK")
	}
	
	// Wrap it with the request logger middleware
	handler := RequestLogger(testHandler)
	
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
	req.SetRequestURI("http://localhost/test")
	req.Header.SetMethod("GET")
	
	// Create a response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	
	// Send request
	if err := c.Do(req, resp); err != nil {
		t.Fatalf("Error sending request: %s", err)
	}
	
	// Verify the response
	if resp.StatusCode() != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode())
	}
	
	// Verify the log output
	logOutput := buffer.String()
	
	// Check that both request received and completed logs are present
	if !strings.Contains(logOutput, "Request received") {
		t.Error("Expected log to contain 'Request received'")
	}
	
	if !strings.Contains(logOutput, "HTTP request completed") {
		t.Error("Expected log to contain 'HTTP request completed'")
	}
	
	// Check that method and path are logged
	if !strings.Contains(strings.ToLower(logOutput), "method") {
		t.Error("Expected log to contain the request method")
	}
	
	if !strings.Contains(strings.ToLower(logOutput), "path") {
		t.Error("Expected log to contain the request path")
	}
}