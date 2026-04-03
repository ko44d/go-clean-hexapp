package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ko44d/go-clean-hexapp/internal/interface/handler"
)

func NewRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", taskHandler.GetTasks)
	r.POST("/tasks", taskHandler.AddTask)
	r.POST("/tasks/complete", taskHandler.CompleteTask)

	return r
}
