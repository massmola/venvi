package cmd

import (
	"fmt"
	"os"
	"strings"
	"venvi/agent/internal/log"
	"venvi/agent/internal/prompts"

	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:       "prompt [role] [args...]",
	Short:     "Get a prompt template for a specific role",
	Long:      `Generates a prompt for the Orchestrator or Critic role.`,
	ValidArgs: []string{"orchestrator", "critic"},
	Args:      cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		role := strings.ToLower(args[0])

		switch role {
		case "orchestrator":
			if len(args) < 2 {
				_, _ = fmt.Fprintln(os.Stdout, "Usage: agent prompt orchestrator <goal>")
				return
			}
			goal := strings.Join(args[1:], " ")
			_, _ = fmt.Fprintln(os.Stdout, prompts.GetOrchestratorPrompt(goal))

		case "critic":
			if len(args) < 2 {
				_, _ = fmt.Fprintln(os.Stdout, "Usage: agent prompt critic <session_id>")
				return
			}
			sessionID := args[1]
			logger := log.NewLogger(".agent_data")
			session, err := logger.GetSession(sessionID)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stdout, "Error retrieving session logs: %v\n", err)
				return
			}

			// Format logs for the prompt
			var logContent strings.Builder
			for _, entry := range session.Entries {
				_, _ = fmt.Fprintf(&logContent, "[%s] %s: %s\n", entry.Role, entry.Timestamp.Format("15:04:05"), entry.Content)
			}

			_, _ = fmt.Fprintln(os.Stdout, prompts.GetCriticPrompt(logContent.String()))

		default:
			_, _ = fmt.Fprintf(os.Stdout, "Unknown role: %s. Valid roles: orchestrator, critic\n", role)
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
