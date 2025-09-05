package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ko44d/go-clean-hexapp/config"
	"github.com/ko44d/go-clean-hexapp/internal/container"
	"github.com/ko44d/go-clean-hexapp/internal/router"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	c, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatalf("failed to initialize container: %v", err)
	}

	r := router.NewRouter(c)

	port := fmt.Sprintf(":%s", cfg.HTTP.Port)
	log.Printf("server starting on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
