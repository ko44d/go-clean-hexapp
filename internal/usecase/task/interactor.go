//go:generate mockgen -source=interactor.go -destination=mocks/mock_interactor.go -package=mocks

package task

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
)

type Interactor interface {
	GetTasks(ctx context.Context) ([]TaskOutput, error)
	AddTask(ctx context.Context, title string) error
	CompleteTask(ctx context.Context, id string) error
}

type interactor struct {
	repo domain.Repository
}

func New(repo domain.Repository) Interactor {
	return &interactor{repo: repo}
}

func (i *interactor) GetTasks(ctx context.Context) ([]TaskOutput, error) {
	tasks, err := i.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetTasks: %w", err)
	}
	return toTaskOutputs(tasks), nil
}

func (i *interactor) AddTask(ctx context.Context, title string) error {
	now := time.Now()
	task, err := domain.New(uuid.New().String(), title, now, now)
	if err != nil {
		switch err {
		case domain.ErrInvalidTitle, domain.ErrTitleBlank, domain.ErrTitleTooLong:
			return err
		default:
			return fmt.Errorf("AddTask: %w", err)
		}
	}
	if err := i.repo.Create(ctx, task); err != nil {
		return fmt.Errorf("AddTask: %w", err)
	}
	return nil
}

func (i *interactor) CompleteTask(ctx context.Context, id string) error {
	task, err := i.repo.FindByID(ctx, id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			return err
		}
		return fmt.Errorf("CompleteTask: %w", err)
	}
	task.Complete(time.Now())
	if err := i.repo.Update(ctx, task); err != nil {
		return fmt.Errorf("CompleteTask: %w", err)
	}
	return nil
}
