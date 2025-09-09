package task

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusTodo     Status = "TODO"
	StatusComplete Status = "COMPLETE"
)

type Task struct {
	ID        string
	Title     string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTask(title string) (*Task, error) {
	if title == "" {
		return nil, errors.New("title must not be empty")
	}

	now := time.Now()

	return &Task{
		ID:        uuid.NewString(),
		Title:     title,
		Status:    StatusTodo,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (t *Task) Complete() {
	t.Status = StatusComplete
	t.UpdatedAt = time.Now()
}
