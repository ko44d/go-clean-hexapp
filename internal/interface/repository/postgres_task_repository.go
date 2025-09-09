package repository

import (
	"context"
	"database/sql"
	"errors"

	domain "github.com/ko44d/go-clean-hexapp/internal/domain/task"
)

type postgresTaskRepository struct {
	db *sql.DB
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

func (r *postgresTaskRepository) Complete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET status = $1, updated_at = $2 WHERE id = $3`,
		domain.StatusComplete, domain.Task{}.UpdatedAt, id,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("task not found")
	}
	return nil
}
