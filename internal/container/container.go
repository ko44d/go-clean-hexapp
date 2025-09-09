package container

import (
	"log"

	"github.com/ko44d/go-clean-hexapp/config"
	"github.com/ko44d/go-clean-hexapp/internal/infrastructure/db"
	"github.com/ko44d/go-clean-hexapp/internal/interface/handler"
	"github.com/ko44d/go-clean-hexapp/internal/interface/repository"
	"github.com/ko44d/go-clean-hexapp/internal/usecase/task"
)

type Container struct {
	Handler handler.Handler
}

func New(cfg *config.Config) *Container {
	dbpool, err := db.NewDBPool(cfg.GetDSN())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	repo := repository.NewTaskRepository(dbpool)
	usecase := task.NewInteractor(repo)
	h := handler.NewHandler(usecase)

	return &Container{
		Handler: h,
	}
}
