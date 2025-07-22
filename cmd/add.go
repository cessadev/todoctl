package cmd

import (
	"fmt"
	"os"

	"github.com/cessadev/todoctl/internal/storage"
	"github.com/spf13/cobra"
)

var highPriority bool

var add = &cobra.Command{
	Use:   "add [message]",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskText := args[0]

		store, err := storage.NewStore()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Storage could not be initialized: %v\n", err)
			os.Exit(1)
		}

		id, err := store.Add(taskText, highPriority)
		if err != nil {
			fmt.Fprintf(os.Stderr, "The task could not be saved: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Task '%s' added with ID #%d\n", taskText, id)
	},
}

func init() {
	add.Flags().BoolVarP(&highPriority, "high-priority", "p", false, "Mark as high priority")
	add.Flags().BoolVarP(&onlyRegular, "regular-task", "r", false, "Show only regular tasks")
	root.AddCommand(add)
}
