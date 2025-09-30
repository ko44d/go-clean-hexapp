package task_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ko44d/go-clean-hexapp/internal/domain/task/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
	"github.com/ko44d/go-clean-hexapp/internal/usecase/task"
)

func TestTaskInteractor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Interactor Suite")
}

var _ = Describe("Task Interactor", func() {
	var (
		ctrl       *gomock.Controller
		mockRepo   *mocks.MockRepository
		interactor task.Interactor
		ctx        context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepository(ctrl)
		interactor = task.NewInteractor(mockRepo)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("GetTasks", func() {
		Context("when repository returns tasks successfully", func() {
			It("should return all tasks", func() {
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

				mockRepo.EXPECT().FindAll(ctx).Return(expectedTasks, nil)

				tasks, err := interactor.GetTasks(ctx)

				Expect(err).To(BeNil())
				Expect(tasks).To(HaveLen(2))
				Expect(tasks[0].Title).To(Equal("Test Task 1"))
				Expect(tasks[1].Title).To(Equal("Test Task 2"))
			})
		})

		Context("when repository returns an error", func() {
			It("should return the error", func() {
				expectedError := errors.New("database error")
				mockRepo.EXPECT().FindAll(ctx).Return(nil, expectedError)

				tasks, err := interactor.GetTasks(ctx)

				Expect(err).To(Equal(expectedError))
				Expect(tasks).To(BeNil())
			})
		})

		Context("when repository returns empty list", func() {
			It("should return empty list", func() {
				mockRepo.EXPECT().FindAll(ctx).Return([]*domain.Task{}, nil)

				tasks, err := interactor.GetTasks(ctx)

				Expect(err).To(BeNil())
				Expect(tasks).To(BeEmpty())
			})
		})
	})

	Describe("AddTask", func() {
		Context("when title is valid", func() {
			It("should create a new task successfully", func() {
				title := "New Task"
				mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, task *domain.Task) error {
						Expect(task.Title).To(Equal(title))
						Expect(task.Status).To(Equal(domain.StatusTodo))
						Expect(task.ID).NotTo(BeEmpty())
						return nil
					},
				)

				err := interactor.AddTask(ctx, title)

				Expect(err).To(BeNil())
			})
		})

		Context("when title is empty", func() {
			It("should return validation error", func() {
				err := interactor.AddTask(ctx, "")

				Expect(err).To(Equal(domain.ErrInvalidTitle))
			})
		})

		Context("when repository returns an error", func() {
			It("should return the error", func() {
				title := "New Task"
				expectedError := errors.New("database error")
				mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedError)

				err := interactor.AddTask(ctx, title)

				Expect(err).To(Equal(expectedError))
			})
		})
	})

	Describe("CompleteTask", func() {
		Context("when task exists", func() {
			It("should mark task as complete and update it", func() {
				taskID := "task-1"
				existingTask := &domain.Task{
					ID:        taskID,
					Title:     "Test Task",
					Status:    domain.StatusTodo,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				mockRepo.EXPECT().FindByID(ctx, taskID).Return(existingTask, nil)
				mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, task *domain.Task) error {
						Expect(task.ID).To(Equal(taskID))
						Expect(task.Status).To(Equal(domain.StatusComplete))
						return nil
					},
				)

				err := interactor.CompleteTask(ctx, taskID)

				Expect(err).To(BeNil())
			})
		})

		Context("when task does not exist", func() {
			It("should return task not found error", func() {
				taskID := "non-existent-id"
				mockRepo.EXPECT().FindByID(ctx, taskID).Return(nil, domain.ErrTaskNotFound)

				err := interactor.CompleteTask(ctx, taskID)

				Expect(err).To(Equal(domain.ErrTaskNotFound))
			})
		})

		Context("when FindByID returns an error", func() {
			It("should return the error", func() {
				taskID := "task-1"
				expectedError := errors.New("database error")
				mockRepo.EXPECT().FindByID(ctx, taskID).Return(nil, expectedError)

				err := interactor.CompleteTask(ctx, taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when Update returns an error", func() {
			It("should return the error", func() {
				taskID := "task-1"
				existingTask := &domain.Task{
					ID:        taskID,
					Title:     "Test Task",
					Status:    domain.StatusTodo,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				expectedError := errors.New("update failed")

				mockRepo.EXPECT().FindByID(ctx, taskID).Return(existingTask, nil)
				mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(expectedError)

				err := interactor.CompleteTask(ctx, taskID)

				Expect(err).To(Equal(expectedError))
			})
		})
	})
})
