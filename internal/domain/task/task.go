package task

import "time"

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

func NewTask(id string, title string, createdAt time.Time, updatedAt time.Time) (*Task, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}

	return &Task{
		ID:        id,
		Title:     title,
		Status:    StatusTodo,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (t *Task) Complete(now time.Time) {
	t.Status = StatusComplete
	t.UpdatedAt = now
}
