package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cessadev/todoctl/internal/storage"
	"github.com/spf13/cobra"
)

var delete = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete a task by ID",
	Args:    cobra.ExactArgs(1),
	Example: "todoctl delete 3",
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

		/** Get the task before deleting it (to check priority) */
		task, err := store.GetByID(id)
		if err != nil {
			if err == storage.ErrTaskNotFound {
				fmt.Fprintf(os.Stderr, "No task found with ID %d\n", id)
			} else {
				fmt.Fprintf(os.Stderr, "Error when searching for the task: %v\n", err)
			}
			os.Exit(1)
		}

		/**  If it is a high priority task, ask for confirmation */
		if task.HighPriority {
			fmt.Printf("⚠️ Task #%d is of HIGH PRIORITY: \"%s\"\n", id, task.Text)
			fmt.Print("Are you sure you want to delete it (Y/n): ")

			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				response := strings.TrimSpace(strings.ToLower(scanner.Text()))
				if response == "n" || response == "no" {
					fmt.Println("❌ Elimination cancelled")
					return
				}
			}
		}

		err = store.Delete(id)
		if err != nil {
			if err == storage.ErrTaskNotFound {
				fmt.Fprintf(os.Stderr, "No task found with ID %d.\n", id)
			} else {
				fmt.Fprintf(os.Stderr, "Error deleting the task: %v\n", err)
			}
			os.Exit(1)
		}

		fmt.Printf("🗑️ Task #%d successfully deleted\n", id)
	},
}

func init() {
	rootCmd.AddCommand(delete)
}
