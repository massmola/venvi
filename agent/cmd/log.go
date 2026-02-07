package cmd

import (
	"fmt"
	"time"
	"venvi/agent/internal/log"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage agent task logs",
	Long:  `Start, append to, and view task logs for the autonomous agent.`,
}

var logStartCmd = &cobra.Command{
	Use:   "start [session_id]",
	Short: "Start a new logging session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := args[0]
		logger := log.NewLogger(".agent_data")

		if err := logger.StartSession(sessionID); err != nil {
			fmt.Printf("Error starting session: %v\n", err)
			return
		}

		// Set current session pointer
		// In a real app we might store this in a file, for now user must remember ID
		fmt.Printf("Session '%s' started.\n", sessionID)
	},
}

var logAppendCmd = &cobra.Command{
	Use:   "append [session_id] [role] [content]",
	Short: "Append an entry to a session log",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := args[0]
		role := args[1]
		content := args[2]

		logger := log.NewLogger(".agent_data")

		if err := logger.AppendEntry(sessionID, role, content); err != nil {
			fmt.Printf("Error appending log: %v\n", err)
			return
		}

		fmt.Println("Log entry added.")
	},
}

var logShowCmd = &cobra.Command{
	Use:   "show [session_id]",
	Short: "Show a session log",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := args[0]
		logger := log.NewLogger(".agent_data")

		session, err := logger.GetSession(sessionID)
		if err != nil {
			fmt.Printf("Error reading session: %v\n", err)
			return
		}

		fmt.Printf("Session ID: %s\nStarted: %s\n\n", session.ID, session.StartTime.Format(time.RFC3339))
		for _, entry := range session.Entries {
			fmt.Printf("[%s] %s: %s\n", entry.Timestamp.Format("15:04:05"), entry.Role, entry.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.AddCommand(logStartCmd)
	logCmd.AddCommand(logAppendCmd)
	logCmd.AddCommand(logShowCmd)
}
