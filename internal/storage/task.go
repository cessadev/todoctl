package storage

import "time"

type Task struct {
	ID           int       `json:"id"`
	Text         string    `json:"text"`
	Done         bool      `json:"done"`
	HighPriority bool      `json:"high_priority"`
	CreatedAt    time.Time `json:"created_at"`
}
