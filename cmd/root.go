package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "go-pass",
	Version:      "1.0.0",
	Short:        "CLI application",
	Long:         "CLI application",
	SilenceUsage: true,
}

// Execute _
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
