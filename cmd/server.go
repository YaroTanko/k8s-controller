package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"k8s-controller/pkg/logger"
)

var (
	serverPort int
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

		// Set up handler
		handler := func(ctx *fasthttp.RequestCtx) {
			path := string(ctx.Path())

			switch path {
			case "/":
				ctx.SetContentType("text/plain")
				fmt.Fprintf(ctx, "Welcome to the k8s-controller HTTP server!")
			case "/health":
				ctx.SetContentType("application/json")
				fmt.Fprintf(ctx, `{"status":"healthy"}`)
			default:
				ctx.SetStatusCode(404)
				fmt.Fprintf(ctx, "Not found")
			}

			logger.Info().
				Str("method", string(ctx.Method())).
				Str("path", path).
				Int("status", ctx.Response.StatusCode()).
				Msg("HTTP request")
		}

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

	// Bind flag to environment variable
	viper.BindPFlag("server_port", serverCmd.Flags().Lookup("port"))

	// Set environment variable name
	viper.SetDefault("server_port", 8080)

	// Set up automatic binding to SERVER_PORT env var
	viper.AutomaticEnv()
}
