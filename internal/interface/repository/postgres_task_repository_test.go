package repository

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPostgresTaskRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Task Repository Suite")
}

var _ = Describe("postgresTaskRepository", func() {
	Describe("Update", func() {
		var (
			ctx       context.Context
			repo      *postgresTaskRepository
			testTask  *domain.Task
			execState *stubExecState
		)

		BeforeEach(func() {
			ctx = context.Background()
			execState = &stubExecState{}
			repo = &postgresTaskRepository{db: &stubQueryExecutor{execState: execState}}
			testTask = &domain.Task{
				ID:        "task-1",
				Title:     "Test Task",
				Status:    domain.StatusComplete,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		})

		Context("when the task exists", func() {
			It("returns nil", func() {
				execState.rowsAffected = 1

				err := repo.Update(ctx, testTask)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the task does not exist", func() {
			It("returns task not found", func() {
				execState.rowsAffected = 0

				err := repo.Update(ctx, testTask)

				Expect(err).To(Equal(domain.ErrTaskNotFound))
			})
		})

		Context("when ExecContext fails", func() {
			It("returns the execution error", func() {
				expectedErr := errors.New("exec failed")
				execState.execErr = expectedErr

				err := repo.Update(ctx, testTask)

				Expect(err).To(MatchError(MatchRegexp("save task")))
				Expect(errors.Is(err, expectedErr)).To(BeTrue())
			})
		})
	})
})

type stubExecState struct {
	rowsAffected int64
	execErr      error
}

type stubQueryExecutor struct {
	execState *stubExecState
}

func (s *stubQueryExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	if s.execState.execErr != nil {
		return pgconn.CommandTag{}, s.execState.execErr
	}
	return pgconn.NewCommandTag("UPDATE " + strconv.FormatInt(s.execState.rowsAffected, 10)), nil
}

func (s *stubQueryExecutor) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("not implemented")
}

func (s *stubQueryExecutor) QueryRow(context.Context, string, ...any) pgx.Row {
	return stubRow{}
}

type stubRow struct{}

func (stubRow) Scan(...any) error {
	return errors.New("not implemented")
}
