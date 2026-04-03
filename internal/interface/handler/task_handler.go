package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
	"github.com/ko44d/go-clean-hexapp/internal/usecase/task"
)

type TaskResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Status    domain.Status `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

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
	c.JSON(http.StatusOK, toTaskResponses(tasks))
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
		if errors.Is(err, domain.ErrInvalidTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusOK)
}

func toTaskResponses(tasks []*domain.Task) []TaskResponse {
	responses := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		responses = append(responses, toTaskResponse(task))
	}
	return responses
}

func toTaskResponse(task *domain.Task) TaskResponse {
	if task == nil {
		return TaskResponse{}
	}

	return TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		Status:    task.Status,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
