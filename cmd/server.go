package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"k8s-controller/pkg/logger"
	"k8s-controller/pkg/middleware"
	"os"
)

var (
	serverPort int
	debugMode bool
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server using fasthttp.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info().Msg("Starting HTTP server...")

		// Get server port from flag or environment variable
		port := serverPort
		logger.Info().Int("port", port).Msg("Server configuration")

		// Create base handler
		baseHandler := func(ctx *fasthttp.RequestCtx) {
			path := string(ctx.Path())

			switch path {
			case "/":
				ctx.SetContentType("text/plain")
				_, err := fmt.Fprintf(ctx, "Welcome to the k8s-controller HTTP server!")
				if err != nil {
					logger.Error().Err(err).Msg("Failed to write response")
				}
			case "/health":
				ctx.SetContentType("application/json")
				_, err := fmt.Fprintf(ctx, `{"status":"healthy"}`)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to write health response")
				}
			default:
				ctx.SetStatusCode(404)
				_, err := fmt.Fprintf(ctx, "Not found")
				if err != nil {
					logger.Error().Err(err).Msg("Failed to write not found response")
				}
			}
		}
		
		// Set up logging options
		loggingOptions := middleware.DefaultLoggingOptions()
		// Enable header logging for development
		if debugMode {
			loggingOptions.LogHeaders = true
			loggingOptions.LogRequestBody = true
			logger.Info().Msg("Debug mode enabled: detailed request logging activated")
		}
		
		// Wrap base handler with request logging middleware
		handler := middleware.EnhancedRequestLogger(loggingOptions)(baseHandler)

		// Start server
		addr := fmt.Sprintf(":%d", port)
		logger.Info().Str("address", addr).Msg("HTTP server is running")

		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Add server-specific flags
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "HTTP server port")
	serverCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode with detailed request logging")

	// Bind flag to environment variable
	if err := viper.BindPFlag("server_port", serverCmd.Flags().Lookup("port")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding server_port flag: %v\n", err)
		os.Exit(1)
	}

	// Set environment variable name
	viper.SetDefault("server_port", 8080)

	// Set up automatic binding to SERVER_PORT env var
	viper.AutomaticEnv()
}