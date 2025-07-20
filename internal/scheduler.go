package internal

import "time"

const (
	SchedulerStatusCanceled ActionStatus = iota - 1
	SchedulerStatusInSchedule
	SchedulerStatusCompleted
)

const (
	NotificationStatusInProgress NotificationStatus = iota
	NotificationStatusIsSending
	NotificationStatusFailed
	NotificationStatusSent
)

type (
	ActionStatus int8

	NotificationStatus int8

	Scheduler struct {
		ID                 int64              `json:"id"`
		TaskID             int64              `json:"task_id"`
		ActionStatus       ActionStatus       `json:"action_status"`
		NotificationStatus NotificationStatus `json:"notification_status"`
		CreatedAt          time.Time          `json:"created_at"`
		UpdatedAt          time.Time          `json:"updated_at"`
		IsDone             bool               `json:"is_done"` // indicates if the scheduler has completed its action callback via email or API calls
	}
)
