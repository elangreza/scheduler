package service

import (
	"context"

	"github.com/elangreza/scheduler/internal"
)

type (
	sqlRepo interface {
		CreateTask(ctx context.Context, task internal.Task) (err error)
		ListTasks(ctx context.Context) ([]internal.Task, error)
		DeleteTask(ctx context.Context, id int) error
		UpdateTask(ctx context.Context, id int, req internal.UpdateTaskParams) error
		// CreateSchedule(task *internal.Schedule) error
	}

	TaskService struct {
		sqlRepo sqlRepo
	}
)

func NewTaskService(sqlRepo sqlRepo) *TaskService {
	return &TaskService{sqlRepo: sqlRepo}
}

func (s *TaskService) CreateTask(ctx context.Context, req internal.CreateTaskParams) error {
	task, err := internal.NewTask(
		req.Name,
		req.Description,
	)
	if err != nil {
		return err
	}

	err = s.sqlRepo.CreateTask(ctx, *task)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) ListTask(ctx context.Context) ([]internal.Task, error) {
	tasks, err := s.sqlRepo.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []internal.Task{}, nil
	}

	return tasks, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	return s.sqlRepo.DeleteTask(ctx, id)
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, req internal.UpdateTaskParams) error {
	return s.sqlRepo.UpdateTask(ctx, id, req)
}
