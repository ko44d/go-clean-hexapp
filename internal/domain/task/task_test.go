package task_test

import (
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ko44d/go-clean-hexapp/internal/domain/task"
)

func TestTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Domain Suite")
}

var _ = Describe("Task Domain", func() {
	Describe("NewTask", func() {
		Context("when title is valid", func() {
			It("should create a new task with default values", func() {
				taskID := "task-1"
				title := "Test Task"
				createdAt := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				updatedAt := createdAt.Add(5 * time.Minute)

				newTask, err := task.NewTask(taskID, title, createdAt, updatedAt)

				Expect(err).To(BeNil())
				Expect(newTask).NotTo(BeNil())
				Expect(newTask.ID).To(Equal(taskID))
				Expect(newTask.Title).To(Equal(title))
				Expect(newTask.Status).To(Equal(task.StatusTodo))
				Expect(newTask.CreatedAt).To(Equal(createdAt))
				Expect(newTask.UpdatedAt).To(Equal(updatedAt))
			})

			It("should preserve the provided IDs for different tasks", func() {
				baseTime := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				task1, err1 := task.NewTask("task-1", "Task 1", baseTime, baseTime)
				task2, err2 := task.NewTask("task-2", "Task 2", baseTime, baseTime)

				Expect(err1).To(BeNil())
				Expect(err2).To(BeNil())
				Expect(task1.ID).NotTo(Equal(task2.ID))
			})
		})

		Context("when title is empty", func() {
			It("should return ErrInvalidTitle", func() {
				createdAt := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				updatedAt := createdAt
				newTask, err := task.NewTask("task-1", "", createdAt, updatedAt)

				Expect(err).To(MatchError(task.ErrInvalidTitle))
				Expect(newTask).To(BeNil())
			})
		})

		Context("when title is whitespace only", func() {
			It("should return ErrTitleBlank", func() {
				createdAt := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				updatedAt := createdAt
				newTask, err := task.NewTask("task-1", "   ", createdAt, updatedAt)

				Expect(err).To(MatchError(task.ErrTitleBlank))
				Expect(newTask).To(BeNil())
			})
		})

		Context("when title length is at the boundary", func() {
			It("should create a task for a 200 character title", func() {
				createdAt := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				updatedAt := createdAt
				title := strings.Repeat("a", 200)

				newTask, err := task.NewTask("task-1", title, createdAt, updatedAt)

				Expect(err).To(BeNil())
				Expect(newTask).NotTo(BeNil())
				Expect(newTask.Title).To(Equal(title))
			})

			It("should return ErrTitleTooLong for a 201 character title", func() {
				createdAt := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				updatedAt := createdAt
				title := strings.Repeat("a", 201)

				newTask, err := task.NewTask("task-1", title, createdAt, updatedAt)

				Expect(err).To(MatchError(task.ErrTitleTooLong))
				Expect(newTask).To(BeNil())
			})
		})
	})

	Describe("Complete", func() {
		Context("when task is in todo status", func() {
			It("should change status to complete and update timestamp", func() {
				createdAt := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				initialUpdatedAt := createdAt.Add(1 * time.Minute)
				completedAt := createdAt.Add(2 * time.Minute)
				testTask, _ := task.NewTask("task-1", "Test Task", createdAt, initialUpdatedAt)

				testTask.Complete(completedAt)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.CreatedAt).To(Equal(createdAt))
				Expect(testTask.UpdatedAt).To(Equal(completedAt))
				Expect(testTask.UpdatedAt).To(BeTemporally(">", initialUpdatedAt))
			})
		})

		Context("when task is already complete", func() {
			It("should update the timestamp even if already complete", func() {
				createdAt := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				initialUpdatedAt := createdAt.Add(1 * time.Minute)
				firstCompletedAt := createdAt.Add(2 * time.Minute)
				secondCompletedAt := createdAt.Add(3 * time.Minute)
				testTask, _ := task.NewTask("task-1", "Test Task", createdAt, initialUpdatedAt)

				testTask.Complete(firstCompletedAt)
				testTask.Complete(secondCompletedAt)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.CreatedAt).To(Equal(createdAt))
				Expect(testTask.UpdatedAt).To(Equal(secondCompletedAt))
				Expect(testTask.UpdatedAt).To(BeTemporally(">", firstCompletedAt))
			})
		})

		Context("when provided a specific time", func() {
			It("should use the provided time exactly", func() {
				createdAt := time.Date(2025, 9, 30, 11, 0, 0, 0, time.UTC)
				initialUpdatedAt := createdAt.Add(30 * time.Minute)
				specificTime := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				testTask, _ := task.NewTask("task-1", "Test Task", createdAt, initialUpdatedAt)

				testTask.Complete(specificTime)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.CreatedAt).To(Equal(createdAt))
				Expect(testTask.UpdatedAt).To(Equal(specificTime))
				Expect(testTask.UpdatedAt).To(BeTemporally(">", testTask.CreatedAt))
			})
		})
	})

	Describe("Status Constants", func() {
		It("should have correct status values", func() {
			Expect(task.StatusTodo).To(Equal(task.Status("todo")))
			Expect(task.StatusComplete).To(Equal(task.Status("complete")))
		})
	})

	Describe("Error Constants", func() {
		It("should match sentinel errors", func() {
			Expect(task.ErrTaskNotFound).To(MatchError(task.ErrTaskNotFound))
			Expect(task.ErrInvalidTitle).To(MatchError(task.ErrInvalidTitle))
			Expect(task.ErrTitleBlank).To(MatchError(task.ErrTitleBlank))
			Expect(task.ErrTitleTooLong).To(MatchError(task.ErrTitleTooLong))
		})
	})
})
