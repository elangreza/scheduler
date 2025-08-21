package sqliterepo

import (
	"context"
	"database/sql"

	"github.com/elangreza/scheduler/internal"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *taskRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) CreateTask(ctx context.Context, task internal.Task) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO tasks (name, description) VALUES (?, ?)",
		task.Name,
		task.Description,
	)
	return err
}
