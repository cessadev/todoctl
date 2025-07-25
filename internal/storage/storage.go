package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const fileName = "tasks.json"

var (
	ErrTaskNotFound = fmt.Errorf("Task not found")
)

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
			s.maxID = 0
			return nil
		}
		return err
	}

	if len(data) == 0 {
		s.tasks = []Task{}
		s.maxID = 0
		return nil
	}

	if err := json.Unmarshal(data, &s.tasks); err != nil {
		return err
	}

	for _, task := range s.tasks {
		if task.ID > s.maxID {
			s.maxID = task.ID
		}
	}

	return nil
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

func (s *Store) Add(text string, highPriority bool) (int, error) {
	s.maxID++
	newTask := Task{
		ID:           s.maxID,
		Text:         text,
		Done:         false,
		HighPriority: highPriority,
		CreatedAt:    time.Now(),
	}

	s.tasks = append(s.tasks, newTask)
	if err := s.save(); err != nil {
		return 0, err
	}

	return newTask.ID, nil
}

func (s *Store) UpdateDescription(id int, newText string) error {
	s.mu.Lock()
	index := -1
	for i, task := range s.tasks {
		if task.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		s.mu.Unlock()
		return ErrTaskNotFound
	}

	/** Verify that the task is pending */
	if s.tasks[index].Done {
		s.mu.Unlock()
		return fmt.Errorf("cannot update a completed task")
	}

	/** Update description */
	s.tasks[index].Text = newText
	s.mu.Unlock()

	return s.save()
}

func (s *Store) MarkDone(id int) error {
	s.mu.Lock()
	var task *Task
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Done = true
			task = &s.tasks[i]
			break
		}
	}
	s.mu.Unlock()

	if task == nil {
		return ErrTaskNotFound
	}

	return s.save()
}

/** Returns a copy of all tasks */
func (s *Store) GetAll() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	/** Avoids overmodification */
	tasks := make([]Task, len(s.tasks))
	copy(tasks, s.tasks)
	return tasks
}

func (s *Store) GetByID(id int) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			/** Return a copy to avoid external modifications */
			task := s.tasks[i]
			return &task, nil
		}
	}

	return nil, ErrTaskNotFound
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	index := -1
	for i, task := range s.tasks {
		if task.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		s.mu.Unlock()
		return ErrTaskNotFound
	}

	/** Remove slice element */
	s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
	s.mu.Unlock()

	return s.save()
}
