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
				title := "Test Task"

				newTask, err := task.NewTask(title)

				Expect(err).To(BeNil())
				Expect(newTask).NotTo(BeNil())
				Expect(newTask.Title).To(Equal(title))
				Expect(newTask.Status).To(Equal(task.StatusTodo))
				Expect(newTask.ID).NotTo(BeEmpty())
				Expect(newTask.CreatedAt).To(BeTemporally("~", time.Now(), time.Second))
				Expect(newTask.UpdatedAt).To(BeTemporally("~", time.Now(), time.Second))
			})

			It("should generate unique IDs for different tasks", func() {
				task1, err1 := task.NewTask("Task 1")
				task2, err2 := task.NewTask("Task 2")

				Expect(err1).To(BeNil())
				Expect(err2).To(BeNil())
				Expect(task1.ID).NotTo(Equal(task2.ID))
			})
		})

		Context("when title is empty", func() {
			It("should return ErrInvalidTitle", func() {
				newTask, err := task.NewTask("")

				Expect(err).To(Equal(task.ErrInvalidTitle))
				Expect(newTask).To(BeNil())
			})
		})
	})

	Describe("Complete", func() {
		Context("when task is in todo status", func() {
			It("should change status to complete and update timestamp", func() {
				testTask, _ := task.NewTask("Test Task")
				originalUpdatedAt := testTask.UpdatedAt
				time.Sleep(10 * time.Millisecond) // Small delay to ensure timestamp difference

				now := time.Now()
				testTask.Complete(now)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.UpdatedAt).To(Equal(now))
				Expect(testTask.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
			})
		})

		Context("when task is already complete", func() {
			It("should update the timestamp even if already complete", func() {
				testTask, _ := task.NewTask("Test Task")
				firstCompleteTime := time.Now()
				testTask.Complete(firstCompleteTime)

				time.Sleep(10 * time.Millisecond)
				secondCompleteTime := time.Now()
				testTask.Complete(secondCompleteTime)

				Expect(testTask.Status).To(Equal(task.StatusComplete))
				Expect(testTask.UpdatedAt).To(Equal(secondCompleteTime))
				Expect(testTask.UpdatedAt).To(BeTemporally(">", firstCompleteTime))
			})
		})

		Context("when provided a specific time", func() {
			It("should use the provided time exactly", func() {
				testTask, _ := task.NewTask("Test Task")
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
