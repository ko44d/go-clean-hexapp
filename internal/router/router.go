package router

import (
	"net/http"

	"github.com/ko44d/go-clean-hexapp/internal/container"
)

// NewRouter sets up the HTTP routing using handlers from the DI container.
func NewRouter(c *container.Container) http.Handler {
	mux := http.NewServeMux()

	return mux
}
