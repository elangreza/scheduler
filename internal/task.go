package internal

import (
	"fmt"
	"time"
)

type (

	// create sql db migration for this task
	// CREATE TABLE tasks (
	//     id INTEGER PRIMARY KEY AUTOINCREMENT,
	//     name TEXT NOT NULL,
	//     description TEXT,
	//     start_time TIMESTAMP NOT NULL,
	//     end_time TIMESTAMP,
	//     repeat_hourly TEXT,
	//     repeat_daily TEXT
	// );
	Task struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"` // optional, can be nil

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func NewTask(name, description string) (*Task, error) {

	task := &Task{
		Name:        name,
		Description: description,
	}

	if task.Name == "" {
		return nil, fmt.Errorf("task name cannot be empty")
	}

	return task, nil
}
