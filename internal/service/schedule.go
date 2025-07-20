package service

import (
	"context"

	"github.com/elangreza/scheduler/internal"
)

type (
	taskRepository interface {
		CreateTask(task *internal.Task) (taskID int64, err error)
	}

	scheduleRepository interface {
		CreateSchedule(task *internal.Schedule) error
	}

	SchedulerService struct {
		taskRepo     taskRepository
		scheduleRepo scheduleRepository
	}
)

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{}
}

// TODO change to params
// TODO add validation
func (s *SchedulerService) CreateTask(ctx context.Context, req internal.CreateTaskParams) error {
	task, err := internal.NewTask(
		req.Name,
		req.Description,
		req.StartTime,
		req.EndTime,
		req.RepeatHourly,
		req.RepeatDaily,
	)
	if err != nil {
		return err
	}

	taskID, err := s.taskRepo.CreateTask(task)
	if err != nil {
		return err
	}

	if err = s.scheduleRepo.CreateSchedule(internal.NewSchedule(
		taskID,
		task.StartTime,
	)); err != nil {
		return err
	}

	return nil
}
