package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/cessadev/todoctl/internal/storage"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:   "list",
	Short: "Displays all saved tasks",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.NewStore()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Storage could not be loaded: %v\n", err)
			os.Exit(1)
		}

		tasks := store.GetAll()

		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
		})

		if len(tasks) == 0 {
			fmt.Println("📭 There are no pending tasks.")
			return
		}

		for _, task := range tasks {
			status := "⏳ Pending"
			if task.Done {
				status = "✅ Completed"
			}

			priority := "📋 Regular task"
			if task.HighPriority {
				priority = "⚠️  High priority"
			}

			fmt.Printf("- [%d] %s | %s | %s | %s\n",
				task.ID,
				task.Text,
				status,
				priority,
				task.CreatedAt.Format("2006-01-02 15:04"))
		}
	},
}

func init() {
	root.AddCommand(list)
}
