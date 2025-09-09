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

	c := container.New(cfg)
	r := router.NewRouter(c.Handler)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	log.Printf("server starting at %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
