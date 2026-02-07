package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "Venvi Autonomous Agent CLI",
	Long: `A CLI tool to manage the memory, logs, and skills of the Venvi autonomous agent.
This tool facilitates the Perception-Reasoning-Action-Reflection loop by providing
local storage for agent state.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
