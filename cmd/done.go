package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cessadev/tudoctl/internal/storage"
	"github.com/spf13/cobra"
)

var done = &cobra.Command{
	Use:     "done [id]",
	Short:   "Mark a task as completed",
	Args:    cobra.ExactArgs(1),
	Example: "tudoctl done 3",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "The ID must be a valid number\n")
			os.Exit(1)
		}

		store, err := storage.NewStore()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Storage could not be loaded: %v\n", err)
			os.Exit(1)
		}

		err = store.MarkDone(id)
		if err != nil {
			if err == storage.ErrTaskNotFound {
				fmt.Fprintf(os.Stderr, "No task found with ID %d\n", id)
			} else {
				fmt.Fprintf(os.Stderr, "Error when marking the task: %v\n", err)
			}
			os.Exit(1)
		}

		fmt.Printf("Task #%d marked as completed\n", id)
	},
}

func init() {
	rootCmd.AddCommand(done)
}
