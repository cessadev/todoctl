package utils

import (
	"github.com/cessadev/todoctl/internal/storage"
)

func FilterHighPriority(tasks []storage.Task) []storage.Task {
	var filtered []storage.Task
	for _, task := range tasks {
		if task.HighPriority {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func FilterRegularTasks(tasks []storage.Task) []storage.Task {
	var filtered []storage.Task
	for _, task := range tasks {
		if !task.HighPriority {
			filtered = append(filtered, task)
		}
	}
	return filtered
}
