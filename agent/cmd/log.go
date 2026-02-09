package cmd

import (
	"fmt"
	"time"
	"venvi/agent/internal/log"

	"github.com/spf13/cobra"
)

const agentDataDir = ".agent_data"

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage agent task logs",
	Long:  `Start, append to, and view task logs for the autonomous agent.`,
}

var logStartCmd = &cobra.Command{
	Use:   "start [session_id]",
	Short: "Start a new logging session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		logger := log.NewLogger(agentDataDir)

		if err := logger.StartSession(sessionID); err != nil {
			return fmt.Errorf("error starting session: %w", err)
		}

		// Set current session pointer
		// In a real app we might store this in a file, for now user must remember ID
		cmd.Printf("Session '%s' started.\n", sessionID)
		return nil
	},
}

var logAppendCmd = &cobra.Command{
	Use:   "append [session_id] [role] [content]",
	Short: "Append an entry to a session log",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		role := args[1]
		content := args[2]

		logger := log.NewLogger(agentDataDir)

		if err := logger.AppendEntry(sessionID, role, content); err != nil {
			return fmt.Errorf("error appending log: %w", err)
		}

		cmd.Println("Log entry added.")
		return nil
	},
}

var logShowCmd = &cobra.Command{
	Use:   "show [session_id]",
	Short: "Show a session log",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		logger := log.NewLogger(agentDataDir)

		session, err := logger.GetSession(sessionID)
		if err != nil {
			return fmt.Errorf("error reading session: %w", err)
		}

		cmd.Printf("Session ID: %s\nStarted: %s\n\n", session.ID, session.StartTime.Format(time.RFC3339))
		for _, entry := range session.Entries {
			cmd.Printf("[%s] %s: %s\n", entry.Timestamp.Format("15:04:05"), entry.Role, entry.Content)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.AddCommand(logStartCmd)
	logCmd.AddCommand(logAppendCmd)
	logCmd.AddCommand(logShowCmd)
}
