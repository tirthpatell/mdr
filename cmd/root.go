package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdr",
	Short: "Markdown renderer, editor, and linter for the terminal",
	Long:  "mdr is a CLI tool for viewing, editing, and linting markdown files directly in your terminal.",
	// Runtime errors (unreadable file, lint findings) shouldn't dump usage
	// or be printed twice; Execute prints them once.
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
