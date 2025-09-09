package task

import (
	"context"
	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
)

type Interactor interface {
	GetTasks(ctx context.Context) ([]*domain.Task, error)
	AddTask(ctx context.Context, title string) error
	CompleteTask(ctx context.Context, id string) error
}

type interactor struct {
	repo domain.Repository
}

func NewInteractor(repo domain.Repository) Interactor {
	return &interactor{repo: repo}
}

func (i *interactor) GetTasks(ctx context.Context) ([]*domain.Task, error) {
	return i.repo.FindAll(ctx)
}

func (i *interactor) AddTask(ctx context.Context, title string) error {
	task, err := domain.NewTask(title)
	if err != nil {
		return err
	}
	return i.repo.Create(ctx, task)
}

func (i *interactor) CompleteTask(ctx context.Context, id string) error {
	return i.repo.Complete(ctx, id)
}
