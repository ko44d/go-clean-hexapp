package task

import (
	"time"

	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
)

type TaskOutput struct {
	ID        string
	Title     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func toTaskOutputs(tasks []*domain.Task) []TaskOutput {
	outputs := make([]TaskOutput, 0, len(tasks))
	for _, task := range tasks {
		outputs = append(outputs, toTaskOutput(task))
	}
	return outputs
}

func toTaskOutput(task *domain.Task) TaskOutput {
	if task == nil {
		return TaskOutput{}
	}

	return TaskOutput{
		ID:        task.ID,
		Title:     task.Title,
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
