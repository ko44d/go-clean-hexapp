package container

import (
	"fmt"

	"github.com/ko44d/go-clean-hexapp/config"
	"github.com/ko44d/go-clean-hexapp/internal/infrastructure/db"
	"github.com/ko44d/go-clean-hexapp/internal/interface/handler"
	"github.com/ko44d/go-clean-hexapp/internal/interface/repository"
	"github.com/ko44d/go-clean-hexapp/internal/usecase/task"
)

type Container struct {
	Handler *handler.TaskHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
	dbPool, err := db.NewDB(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	repo := repository.NewTaskRepository(dbPool)
	usecase := task.NewInteractor(repo)
	h := handler.NewHandler(usecase)

	return &Container{
		Handler: h,
	}, nil
}
