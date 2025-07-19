package storage

import "sync"

type Store struct {
	mu    sync.Mutex
	tasks []Task
	maxID int
	path  string
}
