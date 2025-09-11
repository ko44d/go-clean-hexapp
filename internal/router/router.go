package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ko44d/go-clean-hexapp/internal/interface/handler"
)

func NewRouter(h handler.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", h.GetTasks)
	r.POST("/tasks", h.AddTask)
	r.POST("/tasks/complete", h.CompleteTask)

	return r
}
