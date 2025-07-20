package internal

import (
	"fmt"
	"slices"
	"time"
)

type (
	Task struct {
		ID           int64     `json:"id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"` // optional, can be nil
		StartTime    time.Time `json:"start_time"`
		EndTime      time.Time `json:"end_time"`      // optional, can be nil if the task is ongoing
		RepeatHourly string    `json:"repeat_hourly"` // e.g., "1h", "30m", etc.
		RepeatDaily  []int     `json:"repeat_daily"`  // days of the week, e.g., [1, 2, 3] for Mon, Tue, Wed

		// isRoutine indicates if the task is a routine task
		isRoutine bool `json:"is_routine"`
		// repeatInterval is parsed repeatHourly in time.Duration format
		repeatInterval time.Duration `json:"repeat_duration"`
		// nextRunAt indicates the next scheduled run time for the task
		nextRunAt time.Time `json:"next_run_at"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func NewTask(name, description, startTime, endTime, repeatHourly string, repeatDaily []int) (*Task, error) {

	task := &Task{
		Name:         name,
		Description:  description,
		RepeatHourly: repeatHourly,
		RepeatDaily:  repeatDaily,
	}

	var err error
	if startTime == "" {
		return nil, fmt.Errorf("start time cannot be empty")
	}

	task.StartTime, err = time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %v", err)
	}

	if endTime != "" {
		task.EndTime, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start time format: %v", err)
		}
	}

	err = task.isValid()
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Task) isValid() error {
	if s.Name == "" {
		return fmt.Errorf("task name cannot be empty")
	}

	if !s.EndTime.IsZero() && s.StartTime.After(s.EndTime) {
		return fmt.Errorf("task start time cannot be after end time")
	}

	if s.RepeatHourly != "" {
		var err error
		s.repeatInterval, err = time.ParseDuration(s.RepeatHourly)
		if err != nil {
			return fmt.Errorf("invalid repeat hourly format: %v", err)
		}

		if s.repeatInterval <= 0 {
			return fmt.Errorf("repeat hourly cannot be equal or less then 0")
		}

		if !s.EndTime.IsZero() && s.StartTime.Add(s.repeatInterval).After(s.EndTime) {
			return fmt.Errorf("repeat hourly cannot exceed end time")
		}

		s.isRoutine = true
	}

	if len(s.RepeatDaily) > 0 {
		for _, day := range s.RepeatDaily {
			if day < int(time.Sunday) || day > int(time.Saturday) {
				return fmt.Errorf("invalid repeat daily value: %d, must be between 0 (Sunday) and 6 (Saturday)", day)
			}
		}

		slices.Sort(s.RepeatDaily)

		s.isRoutine = true
	}

	return nil
}

func (s *Task) GetNextRunAt(lastRun time.Time) time.Time {
	if !s.isRoutine {
		return time.Time{}
	}

	s.nextRunAt = lastRun
	if s.nextRunAt.IsZero() {
		s.nextRunAt = s.StartTime
	}

	s.nextRunAt = s.nextRunAt.Add(s.repeatInterval)

	// nextRUnAt := s.nextRunAt.Add(s.repeatInterval)
	if !s.EndTime.IsZero() && s.nextRunAt.After(s.EndTime) {
		if len(s.RepeatDaily) > 0 {
			nextDay := generateRunTimeSequence(s.nextRunAt, s.RepeatDaily)
			return time.Date(
				nextDay.Year(),
				nextDay.Month(),
				nextDay.Day(),
				s.StartTime.Hour(),
				s.StartTime.Minute(),
				s.StartTime.Second(),
				s.StartTime.Nanosecond(),
				nextDay.Location())
		}

		return time.Time{}
	}

	return s.nextRunAt
}

func generateRunTimeSequence(lasSeq time.Time, sequence []int) time.Time {
	for _, day := range sequence {
		if day < int(time.Sunday) || day > int(time.Saturday) {
			return time.Time{}
		}
	}

	nextDay := lasSeq
	for {
		nextDay = nextDay.AddDate(0, 0, 1)
		if slices.Contains(sequence, int(nextDay.Weekday())) {
			return nextDay
		}
	}
}
