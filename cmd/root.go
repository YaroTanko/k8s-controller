package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s-controller/pkg/config"
	"k8s-controller/pkg/logger"
)

var (
	logLevel   string
	kubeconfig string
	namespace  string
	cfg        *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "k8s-controller",
	Short: "A Kubernetes controller",
	Long: `A Kubernetes controller application that manages custom resources
and performs operations based on Kubernetes events.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return err
		}

		// Override with command line flags if provided
		if cmd.Flags().Changed("log-level") {
			cfg.LogLevel = logLevel
		}
		if cmd.Flags().Changed("kubeconfig") {
			cfg.KubeConfig = kubeconfig
		}
		if cmd.Flags().Changed("namespace") {
			cfg.Namespace = namespace
		}

		// Initialize logger
		logger.Init(logger.LogLevel(cfg.LogLevel))
		logger.Debug().Msg("Debug logging enabled")
		logger.Debug().Interface("config", cfg).Msg("Configuration loaded")

		return nil
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
	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace to operate in")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (trace, debug, info, warn, error)")

	// Bind flags to environment variables
	// Use K8S_CONTROLLER_KUBECONFIG instead of --kubeconfig
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	// Configure Viper
	viper.SetEnvPrefix("K8S_CONTROLLER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
