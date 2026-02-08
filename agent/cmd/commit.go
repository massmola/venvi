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
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")

		changed, err := git.HasChanges()
		if err != nil {
			fmt.Printf("Error checking git status: %v\n", err)
			return
		}

		if !changed {
			fmt.Println("No changes to commit.")
			return
		}

		if err := git.StageAll(); err != nil {
			fmt.Printf("Error staging changes: %v\n", err)
			return
		}

		if err := git.Commit(message); err != nil {
			fmt.Printf("Error committing changes: %v\n", err)
			return
		}

		fmt.Printf("Committed changes with message: agent: %s\n", message)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
