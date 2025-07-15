package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const fileName = "tasks.json"

var (
	ErrTaskNotFound = fmt.Errorf("Task not found")
)

type Store struct {
	mu    sync.Mutex
	tasks []Task
	maxID int
	path  string
}

func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	storeDir := filepath.Join(homeDir, ".todoctl")
	filePath := filepath.Join(storeDir, fileName)

	/** Create the directory if it does not exist */
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			return nil, err
		}
	}

	store := &Store{path: filePath}

	/** Load existing tasks */
	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			s.tasks = []Task{}
			return nil
		}
		return err
	}

	if len(data) == 0 {
		s.tasks = []Task{}
		return nil
	}

	return json.Unmarshal(data, &s.tasks)
}

func (s *Store) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func (s *Store) Add(text string) (int, error) {
	s.maxID++
	newTask := Task{
		ID:        s.maxID,
		Text:      text,
		Done:      false,
		CreatedAt: time.Now(),
	}

	s.tasks = append(s.tasks, newTask)
	if err := s.save(); err != nil {
		return 0, err
	}

	return newTask.ID, nil
}
