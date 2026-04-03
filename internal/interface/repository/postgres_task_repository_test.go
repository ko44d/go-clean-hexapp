package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

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
			ctx        context.Context
			repo       *postgresTaskRepository
			testTask   *domain.Task
			execState  *stubExecState
			database   *sql.DB
			closeDBErr error
		)

		BeforeEach(func() {
			ctx = context.Background()
			execState = &stubExecState{}

			var err error
			database, err = openStubDB(execState)
			Expect(err).NotTo(HaveOccurred())

			repo = &postgresTaskRepository{db: database}
			testTask = &domain.Task{
				ID:        "task-1",
				Title:     "Test Task",
				Status:    domain.StatusComplete,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		})

		AfterEach(func() {
			if database != nil {
				closeDBErr = database.Close()
				Expect(closeDBErr).NotTo(HaveOccurred())
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

				Expect(err).To(Equal(expectedErr))
			})
		})
	})
})

type stubExecState struct {
	rowsAffected int64
	execErr      error
}

type stubDriver struct {
	execState *stubExecState
}

type stubConn struct {
	execState *stubExecState
}

func openStubDB(execState *stubExecState) (*sql.DB, error) {
	return sql.OpenDB(&stubConnector{execState: execState}), nil
}

type stubConnector struct {
	execState *stubExecState
}

func (c *stubConnector) Connect(context.Context) (driver.Conn, error) {
	return &stubConn{execState: c.execState}, nil
}

func (c *stubConnector) Driver() driver.Driver {
	return &stubDriver{execState: c.execState}
}

func (d *stubDriver) Open(string) (driver.Conn, error) {
	return &stubConn{execState: d.execState}, nil
}

func (c *stubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *stubConn) Close() error {
	return nil
}

func (c *stubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

func (c *stubConn) ExecContext(_ context.Context, _ string, _ []driver.NamedValue) (driver.Result, error) {
	if c.execState.execErr != nil {
		return nil, c.execState.execErr
	}
	return driver.RowsAffected(c.execState.rowsAffected), nil
}
