package router

import (
	"net/http"

	"github.com/ko44d/go-clean-hexapp/internal/container"
)

func NewRouter(c *container.Container) http.Handler {
	mux := http.NewServeMux()

	return mux
}
