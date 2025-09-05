package container

import "github.com/ko44d/go-clean-hexapp/config"

type Container struct {
}

func NewContainer(cfg *config.Config) (*Container, error) {
	return &Container{}, nil
}
