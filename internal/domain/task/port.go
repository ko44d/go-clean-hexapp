package task

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Task, error)
	Create(ctx context.Context, task *Task) error
	Complete(ctx context.Context, id string) error
}
