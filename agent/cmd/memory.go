package cmd

import (
	"fmt"
	"strings"
	"time"
	"venvi/agent/internal/memory"

	"github.com/spf13/cobra"
)

const defaultDataDir = ".agent_data"

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage agent memory (skills and lessons)",
	Long:  `Add, search, and recall skills or lessons learned by the agent.`,
}

var addCmd = &cobra.Command{
	Use:   "add <topic> <content> [tags...]",
	Short: "Add a new skill or lesson",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]
		content := args[1]
		tags := args[2:]

		store := memory.NewStore(defaultDataDir)
		skill := memory.Skill{
			Topic:   topic,
			Content: content,
			Tags:    tags,
		}

		if err := store.Save(skill); err != nil {
			return fmt.Errorf("error saving memory: %w", err)
		}

		cmd.Printf("Memory added: %s\n", topic)
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for skills or lessons",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		store := memory.NewStore(defaultDataDir)

		results, err := store.Search(query)
		if err != nil {
			return fmt.Errorf("error searching memory: %w", err)
		}

		if len(results) == 0 {
			cmd.Println("No matching memories found.")
			return nil
		}

		cmd.Printf("Found %d memories:\n", len(results))
		for _, skill := range results {
			cmd.Printf("\n--- [%s] %s ---\n", skill.CreatedAt.Format(time.RFC3339), skill.Topic)
			cmd.Printf("Tags: %s\n", strings.Join(skill.Tags, ", "))
			cmd.Printf("Content:\n%s\n", skill.Content)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(memoryCmd)
	memoryCmd.AddCommand(addCmd)
	memoryCmd.AddCommand(searchCmd)
}
