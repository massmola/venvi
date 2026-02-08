package cmd

import (
	"fmt"
	"strings"
	"venvi/agent/internal/git"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit [message]",
	Short: "Auto-commit changes with an 'agent:' prefix",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := strings.Join(args, " ")

		changed, err := git.HasChanges()
		if err != nil {
			return fmt.Errorf("error checking git status: %w", err)
		}

		if !changed {
			cmd.Println("No changes to commit.")
			return nil
		}

		if err := git.StageAll(); err != nil {
			return fmt.Errorf("error staging changes: %w", err)
		}

		if err := git.Commit("agent: " + message); err != nil {
			return fmt.Errorf("error committing changes: %w", err)
		}

		cmd.Printf("Committed changes with message: agent: %s\n", message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
