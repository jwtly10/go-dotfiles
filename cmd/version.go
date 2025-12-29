package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the version of the binary`,
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", version)
}
