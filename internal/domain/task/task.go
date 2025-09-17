package task

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusTodo     Status = "todo"
	StatusComplete Status = "complete"
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
		return nil, ErrInvalidTitle
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

func (t *Task) Complete(now time.Time) {
	t.Status = StatusComplete
	t.UpdatedAt = now
}
