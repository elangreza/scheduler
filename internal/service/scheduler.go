package service

import (
	"context"

	"github.com/elangreza/scheduler/internal"
)

type (
	taskRepository interface {
		CreateTask(task *internal.Task) error
	}

	SchedulerService struct {
		taskRepo taskRepository
	}
)

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{}
}

func (s *SchedulerService) CreateTask(ctx context.Context, task *internal.Task) error {
	return s.taskRepo.CreateTask(task)
}
