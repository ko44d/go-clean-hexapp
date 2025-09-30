package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
	"github.com/ko44d/go-clean-hexapp/internal/interface/handler"
	"github.com/ko44d/go-clean-hexapp/internal/usecase/task/mocks"
)

func TestTaskHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Handler Suite")
}

var _ = Describe("Task Handler", func() {
	var (
		ctrl           *gomock.Controller
		mockInteractor *mocks.MockInteractor
		taskHandler    handler.Handler
		router         *gin.Engine
		recorder       *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		ctrl = gomock.NewController(GinkgoT())
		mockInteractor = mocks.NewMockInteractor(ctrl)
		taskHandler = handler.NewHandler(mockInteractor)
		router = gin.New()
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("GetTasks", func() {
		Context("when tasks are retrieved successfully", func() {
			It("should return 200 with tasks list", func() {
				expectedTasks := []*domain.Task{
					{
						ID:        "task-1",
						Title:     "Test Task 1",
						Status:    domain.StatusTodo,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        "task-2",
						Title:     "Test Task 2",
						Status:    domain.StatusComplete,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}

				mockInteractor.EXPECT().GetTasks(gomock.Any()).Return(expectedTasks, nil)

				router.GET("/tasks", taskHandler.GetTasks)
				req, _ := http.NewRequest("GET", "/tasks", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				var response []*domain.Task
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response).To(HaveLen(2))
				Expect(response[0].Title).To(Equal("Test Task 1"))
				Expect(response[1].Title).To(Equal("Test Task 2"))
			})
		})

		Context("when usecase returns an error", func() {
			It("should return 500 with error message", func() {
				mockInteractor.EXPECT().GetTasks(gomock.Any()).Return(nil, errors.New("database error"))

				router.GET("/tasks", taskHandler.GetTasks)
				req, _ := http.NewRequest("GET", "/tasks", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("failed to get tasks"))
			})
		})

		Context("when there are no tasks", func() {
			It("should return 200 with empty list", func() {
				mockInteractor.EXPECT().GetTasks(gomock.Any()).Return([]*domain.Task{}, nil)

				router.GET("/tasks", taskHandler.GetTasks)
				req, _ := http.NewRequest("GET", "/tasks", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				var response []*domain.Task
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response).To(BeEmpty())
			})
		})
	})

	Describe("AddTask", func() {
		Context("when request body is valid", func() {
			It("should return 201", func() {
				requestBody := map[string]string{"title": "New Task"}
				jsonBody, _ := json.Marshal(requestBody)

				mockInteractor.EXPECT().AddTask(gomock.Any(), "New Task").Return(nil)

				router.POST("/tasks", taskHandler.AddTask)
				req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 with error message", func() {
				router.POST("/tasks", taskHandler.AddTask)
				req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("invalid json")))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("invalid request body"))
			})
		})

		Context("when title is empty", func() {
			It("should return 400 with error message", func() {
				requestBody := map[string]string{"title": ""}
				jsonBody, _ := json.Marshal(requestBody)

				router.POST("/tasks", taskHandler.AddTask)
				req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("invalid request body"))
			})
		})

		Context("when usecase returns invalid title error", func() {
			It("should return 400 with error message", func() {
				requestBody := map[string]string{"title": "Test Task"}
				jsonBody, _ := json.Marshal(requestBody)

				mockInteractor.EXPECT().AddTask(gomock.Any(), "Test Task").Return(domain.ErrInvalidTitle)

				router.POST("/tasks", taskHandler.AddTask)
				req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("invalid title"))
			})
		})

		Context("when usecase returns internal error", func() {
			It("should return 500 with error message", func() {
				requestBody := map[string]string{"title": "Test Task"}
				jsonBody, _ := json.Marshal(requestBody)

				mockInteractor.EXPECT().AddTask(gomock.Any(), "Test Task").Return(errors.New("database error"))

				router.POST("/tasks", taskHandler.AddTask)
				req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("internal server error"))
			})
		})
	})

	Describe("CompleteTask", func() {
		Context("when task ID is provided and task exists", func() {
			It("should return 204", func() {
				taskID := "task-1"

				mockInteractor.EXPECT().CompleteTask(gomock.Any(), taskID).Return(nil)

				router.PUT("/tasks/complete", taskHandler.CompleteTask)
				req, _ := http.NewRequest("PUT", "/tasks/complete?id="+taskID, nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when task ID is missing", func() {
			It("should return 400 with error message", func() {
				router.PUT("/tasks/complete", taskHandler.CompleteTask)
				req, _ := http.NewRequest("PUT", "/tasks/complete", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("missing id"))
			})
		})

		Context("when task is not found", func() {
			It("should return 404 with error message", func() {
				taskID := "non-existent-id"

				mockInteractor.EXPECT().CompleteTask(gomock.Any(), taskID).Return(domain.ErrTaskNotFound)

				router.PUT("/tasks/complete", taskHandler.CompleteTask)
				req, _ := http.NewRequest("PUT", "/tasks/complete?id="+taskID, nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("task not found"))
			})
		})

		Context("when usecase returns internal error", func() {
			It("should return 500 with error message", func() {
				taskID := "task-1"

				mockInteractor.EXPECT().CompleteTask(gomock.Any(), taskID).Return(errors.New("database error"))

				router.PUT("/tasks/complete", taskHandler.CompleteTask)
				req, _ := http.NewRequest("PUT", "/tasks/complete?id="+taskID, nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				Expect(err).To(BeNil())
				Expect(response["error"]).To(Equal("internal server error"))
			})
		})
	})
})
