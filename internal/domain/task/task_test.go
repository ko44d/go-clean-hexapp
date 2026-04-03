package task_test

import (
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
				now := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				newTask, err := task.NewTask("task-1", "", now, now)

				Expect(err).To(Equal(task.ErrInvalidTitle))
				Expect(newTask).To(BeNil())
			})
		})
	})

	Describe("Complete", func() {
		Context("when task is in todo status", func() {
			It("should change status to complete and update timestamp", func() {
				baseTime := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				testTask, _ := task.NewTask("task-1", "Test Task", baseTime, baseTime)
				originalUpdatedAt := testTask.UpdatedAt

				now := baseTime.Add(10 * time.Millisecond)
				testTask.Complete(now)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.UpdatedAt).To(Equal(now))
				Expect(testTask.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
			})
		})

		Context("when task is already complete", func() {
			It("should update the timestamp even if already complete", func() {
				baseTime := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)
				testTask, _ := task.NewTask("task-1", "Test Task", baseTime, baseTime)
				firstCompleteTime := baseTime.Add(10 * time.Millisecond)
				testTask.Complete(firstCompleteTime)

				secondCompleteTime := baseTime.Add(20 * time.Millisecond)
				testTask.Complete(secondCompleteTime)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.UpdatedAt).To(Equal(secondCompleteTime))
				Expect(testTask.UpdatedAt).To(BeTemporally(">", firstCompleteTime))
			})
		})

		Context("when provided a specific time", func() {
			It("should use the provided time exactly", func() {
				baseTime := time.Date(2025, 9, 30, 11, 0, 0, 0, time.UTC)
				testTask, _ := task.NewTask("task-1", "Test Task", baseTime, baseTime)
				specificTime := time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC)

				testTask.Complete(specificTime)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.UpdatedAt).To(Equal(specificTime))
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
		It("should have correct error messages", func() {
			Expect(task.ErrTaskNotFound.Error()).To(Equal("task not found"))
			Expect(task.ErrInvalidTitle.Error()).To(Equal("title must not be empty"))
		})
	})
})
