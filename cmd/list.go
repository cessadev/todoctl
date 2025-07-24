package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/cessadev/tudoctl/internal/storage"
	"github.com/cessadev/tudoctl/utils"
	"github.com/spf13/cobra"
)

var (
	onlyHighPriority  bool
	onlyRegular       bool
	showCompleted     bool
	showAll           bool
	regularCompleted  bool
	priorityCompleted bool
)

var list = &cobra.Command{
	Use:     "list",
	Short:   "Displays all saved tasks",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.NewStore()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Storage could not be loaded: %v\n", err)
			os.Exit(1)
		}

		tasks := store.GetAll()

		/** Sort by most recent date */
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
		})

		/** Validate mutually exclusive flags */
		/**
		flagCount := 0
		if showCompleted {
			flagCount++
		}
		if showAll {
			flagCount++
		}
		if onlyHighPriority {
			flagCount++
		}
		if onlyRegular {
			flagCount++
		}
		*/

		// Validar flags mutuamente excluyentes
		flagCount := 0
		flags := []bool{showCompleted, showAll, onlyHighPriority, onlyRegular, regularCompleted, priorityCompleted}
		for _, f := range flags {
			if f {
				flagCount++
			}
		}

		if flagCount > 1 {
			fmt.Fprintln(os.Stderr, "You can only use one filter at a time")
			os.Exit(1)
		}

		/** Apply filters */
		if regularCompleted {
			tasks = utils.FilterCompleted(tasks)
			tasks = utils.FilterRegularTasks(tasks)
			if len(tasks) == 0 {
				fmt.Println("📭 No regular tasks completed")
				return
			}
			fmt.Println("✅ Regular tasks completed:")
		} else if priorityCompleted {
			tasks = utils.FilterCompleted(tasks)
			tasks = utils.FilterHighPriority(tasks)
			if len(tasks) == 0 {
				fmt.Println("📭 No high priority tasks completed")
				return
			}
			fmt.Println("✅ High priority tasks completed:")
		} else if showCompleted {
			tasks = utils.FilterCompleted(tasks)
			if len(tasks) == 0 {
				fmt.Println("📭 No tasks completed")
				return
			}
			fmt.Println("✅ Tasks completed:")
		} else if showAll {
			if len(tasks) == 0 {
				fmt.Println("📭 No tasks")
				return
			}
			fmt.Println("📋 All tasks:")
		} else {
			/** Default: only pending tasks */
			tasks = utils.FilterPending(tasks)
			if len(tasks) == 0 {
				fmt.Println("🎉 ¡No pending tasks! 🎉")
				return
			}
			fmt.Println("📋 Pending tasks:")
		}

		/** Apply additional filters */
		if onlyHighPriority {
			tasks = utils.FilterHighPriority(tasks)
		}

		if onlyRegular {
			tasks = utils.FilterRegularTasks(tasks)
		}

		/** Show tasks */
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
	list.Flags().BoolVarP(&onlyHighPriority, "high-priority", "p", false, "Show only high priority tasks")
	list.Flags().BoolVarP(&onlyRegular, "regular", "r", false, "Show only regular tasks")
	list.Flags().BoolVarP(&showCompleted, "completed", "c", false, "Show only completed tasks")
	list.Flags().BoolVarP(&showAll, "all", "a", false, "Show all tasks")
	list.Flags().BoolVarP(&regularCompleted, "regular-completed", "R", false, "Show only completed regular tasks")
	list.Flags().BoolVarP(&priorityCompleted, "priority-completed", "P", false, "Show only completed high priority tasks")
	rootCmd.AddCommand(list)
}
