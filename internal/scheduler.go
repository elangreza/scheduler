package internal

import "time"

const (
	StatusCanceled ActionStatus = iota - 1
	StatusCreated
	StatusSending
	StatusFailed
	StatusSuccess
)

type (
	ActionStatus int8

	Schedule struct {
		ID       int64        `json:"id"`
		TaskID   int64        `json:"task_id"`
		Status   ActionStatus `json:"action_status"`
		NotifyAt time.Time    `json:"notify_at"`
		DoneAt   time.Time    `json:"done_at"`
		IsDone   bool         `json:"is_done"` // indicates if the scheduler has completed its action callback via email or API calls

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func NewSchedule(taskID int64, notifyAt time.Time) *Schedule {
	return &Schedule{
		TaskID:   taskID,
		Status:   StatusCreated,
		NotifyAt: notifyAt,
		IsDone:   false,
	}
}
