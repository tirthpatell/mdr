package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/tirthpatell/mdr/internal/editor"
)

var editCmd = &cobra.Command{
	Use:   "edit <file>",
	Short: "Open a markdown file in the TUI editor with live preview",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := editor.NewModelFromFile(args[0])
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}

		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
