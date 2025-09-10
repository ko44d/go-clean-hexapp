package repository

import (
	"context"
	"database/sql"

	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
)

type postgresTaskRepository struct {
	db *sql.DB
}

func (r *postgresTaskRepository) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, title, status, created_at, updated_at FROM tasks WHERE id = $1`, id)

	var task domain.Task
	err := row.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *postgresTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tasks SET title = $1, status = $2, updated_at = $3 WHERE id = $4`,
		task.Title, task.Status, task.UpdatedAt, task.ID)
	return err
}

func NewTaskRepository(db *sql.DB) domain.Repository {
	return &postgresTaskRepository{db: db}
}

func (r *postgresTaskRepository) FindAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, status, created_at, updated_at FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		t := &domain.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *postgresTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (id, title, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		task.ID, task.Title, task.Status, task.CreatedAt, task.UpdatedAt,
	)
	return err
}
