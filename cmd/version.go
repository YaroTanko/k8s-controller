package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s-controller/pkg/logger"
)

var (
	// These variables will be set at build time
	version   = "0.1.0"
	commit    = "none"
	buildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version, commit, and build date information for the k8s-controller.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Log structured version info
		logger.Info().
			Str("version", version).
			Str("commit", commit).
			Str("buildDate", buildDate).
			Msg("Version information")

		// Also print to stdout for CLI usage
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
