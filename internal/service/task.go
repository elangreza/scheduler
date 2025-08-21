package service

import (
	"context"

	"github.com/elangreza/scheduler/internal"
)

type (
	sqlRepo interface {
		CreateTask(ctx context.Context, task internal.Task) (err error)
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
