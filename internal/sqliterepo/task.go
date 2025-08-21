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

func (r *taskRepository) ListTasks(ctx context.Context) ([]internal.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []internal.Task
	for rows.Next() {
		var task internal.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (r *taskRepository) UpdateTask(ctx context.Context, id int, req internal.UpdateTaskParams) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tasks SET name = ?, description = ? WHERE id = ?",
		req.Name,
		req.Description,
		id,
	)
	return err
}
