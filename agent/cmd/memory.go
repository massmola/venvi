package cmd

import (
	"fmt"
	"strings"
	"time"
	"venvi/agent/internal/memory"

	"github.com/spf13/cobra"
)

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage agent memory (skills and lessons)",
	Long:  `Add, search, and recall skills or lessons learned by the agent.`,
}

var addCmd = &cobra.Command{
	Use:   "add [topic] [content] [tags...]",
	Short: "Add a new skill or lesson",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		topic := args[0]
		content := args[1]
		tags := args[2:]

		store := memory.NewStore(".agent_data")
		skill := memory.Skill{
			Topic:   topic,
			Content: content,
			Tags:    tags,
		}

		if err := store.Save(skill); err != nil {
			fmt.Printf("Error saving memory: %v\n", err)
			return
		}

		fmt.Printf("Memory added: %s\n", topic)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for skills or lessons",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		store := memory.NewStore(".agent_data")

		results, err := store.Search(query)
		if err != nil {
			fmt.Printf("Error searching memory: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("No matching memories found.")
			return
		}

		fmt.Printf("Found %d memories:\n", len(results))
		for _, skill := range results {
			fmt.Printf("\n--- [%s] %s ---\n", skill.CreatedAt.Format(time.RFC3339), skill.Topic)
			fmt.Printf("Tags: %s\n", strings.Join(skill.Tags, ", "))
			fmt.Printf("Content:\n%s\n", skill.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(memoryCmd)
	memoryCmd.AddCommand(addCmd)
	memoryCmd.AddCommand(searchCmd)
}
