package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ko44d/go-clean-hexapp/internal/usecase/task"
)

type Handler interface {
	GetTasks(c *gin.Context)
	AddTask(c *gin.Context)
	CompleteTask(c *gin.Context)
}

type taskHandler struct {
	usecase task.Interactor
}

func NewHandler(usecase task.Interactor) Handler {
	return &taskHandler{usecase: usecase}
}

func (h *taskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.usecase.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *taskHandler) AddTask(c *gin.Context) {
	type request struct {
		Title string `json:"title"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.usecase.AddTask(c.Request.Context(), req.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add task"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *taskHandler) CompleteTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := h.usecase.CompleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete task"})
		return
	}
	c.Status(http.StatusNoContent)
}
