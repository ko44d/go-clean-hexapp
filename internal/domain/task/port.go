//go:generate mockgen -source=port.go -destination=mocks/mock_repository.go -package=mocks

package task

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Task, error)
	FindByID(ctx context.Context, id string) (*Task, error)
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
}
