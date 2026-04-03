package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ko44d/go-clean-hexapp/internal/usecase/task"
)

type TaskResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskHandler struct {
	usecase task.Interactor
}

func NewHandler(usecase task.Interactor) *TaskHandler {
	return &TaskHandler{usecase: usecase}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.usecase.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}
	c.JSON(http.StatusOK, toTaskResponses(tasks))
}

func (h *TaskHandler) AddTask(c *gin.Context) {
	type request struct {
		Title string `json:"title"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.usecase.AddTask(c.Request.Context(), req.Title); err != nil {
		if errors.Is(err, task.ErrInvalidTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *TaskHandler) CompleteTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.CompleteTask(c.Request.Context(), id); err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusOK)
}

func toTaskResponses(tasks []task.TaskOutput) []TaskResponse {
	responses := make([]TaskResponse, 0, len(tasks))
	for _, taskOutput := range tasks {
		responses = append(responses, toTaskResponse(taskOutput))
	}
	return responses
}

func toTaskResponse(taskOutput task.TaskOutput) TaskResponse {
	return TaskResponse{
		ID:        taskOutput.ID,
		Title:     taskOutput.Title,
		Status:    taskOutput.Status,
		CreatedAt: taskOutput.CreatedAt,
		UpdatedAt: taskOutput.UpdatedAt,
	}
}
