package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
)

type queryExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type postgresTaskRepository struct {
	db queryExecutor
}

func (r *postgresTaskRepository) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	row := r.db.QueryRow(ctx, `SELECT id, title, status, created_at, updated_at FROM tasks WHERE id = $1`, id)

	var task domain.Task
	err := row.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find task by id %q: %w", id, err)
	}
	return &task, nil
}

func (r *postgresTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	result, err := r.db.Exec(ctx, `UPDATE tasks SET title = $1, status = $2, updated_at = $3 WHERE id = $4`,
		task.Title, task.Status, task.UpdatedAt, task.ID)
	if err != nil {
		return fmt.Errorf("save task %q: %w", task.ID, err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func NewTaskRepository(db *pgxpool.Pool) domain.Repository {
	return &postgresTaskRepository{db: db}
}

func (r *postgresTaskRepository) FindAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, status, created_at, updated_at FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		t := &domain.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("list tasks: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return tasks, nil
}

func (r *postgresTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO tasks (id, title, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		task.ID, task.Title, task.Status, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save task %q: %w", task.ID, err)
	}
	return nil
}
