package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s-controller/pkg/logger"
)

var (
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   "k8s-controller",
	Short: "A Kubernetes controller",
	Long: `A Kubernetes controller application that manages custom resources
and performs operations based on Kubernetes events.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger
		logger.Init(logger.LogLevel(logLevel))
		logger.Debug().Msg("Debug logging enabled")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info().Msg("Starting k8s-controller")
		fmt.Println("Welcome to the Kubernetes controller! Use --help for available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "Path to kubeconfig file")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Kubernetes namespace to operate in")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (trace, debug, info, warn, error)")
}