package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cessadev/tudoctl/internal/storage"
	"github.com/spf13/cobra"
)

var update = &cobra.Command{
	Use:     "update [id] [new description]",
	Short:   "Update the description of a pending task",
	Args:    cobra.ExactArgs(2),
	Example: "tudoctl update 3 \"New task description\"",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "The ID must be a valid number\n")
			os.Exit(1)
		}

		newDescription := args[1]
		if newDescription == "" {
			fmt.Fprintf(os.Stderr, "The description cannot be empty\n")
			os.Exit(1)
		}

		store, err := storage.NewStore()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Storage could not be loaded: %v\n", err)
			os.Exit(1)
		}

		/** Verify that the task exists and is pending */
		task, err := store.GetByID(id)
		if err != nil {
			if err == storage.ErrTaskNotFound {
				fmt.Fprintf(os.Stderr, "No task with ID %d.\n", id)
			} else {
				fmt.Fprintf(os.Stderr, "❌ Error searching for task: %v\n", err)
			}
			os.Exit(1)
		}

		/** Only allow updating of pending tasks */
		if task.Done {
			fmt.Fprintf(os.Stderr, "You cannot update a task that has already been completed.\n")
			os.Exit(1)
		}

		/** Update description */
		err = store.UpdateDescription(id, newDescription)
		if err != nil {
			if err == storage.ErrTaskNotFound {
				fmt.Fprintf(os.Stderr, "No task with ID %d.\n", id)
			} else {
				fmt.Fprintf(os.Stderr, "❌ Error updating task: %v\n", err)
			}
			os.Exit(1)
		}

		fmt.Printf("✏️ Task #%d successfully updated.\n", id)
		fmt.Printf("   Before: %s\n", task.Text)
		fmt.Printf("   Now: %s\n", newDescription)
	},
}

func init() {
	rootCmd.AddCommand(update)
}
