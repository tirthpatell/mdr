package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set via ldflags at build time
var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of mdr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mdr %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
