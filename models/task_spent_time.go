package models

import "time"

type TaskSpentTime struct {
	talbelName  struct{} `sql:"public.task_spent_times"`
	ID          uint64
	SpentByID   uint64 `doc:"automatically set to current user"`
	SpentBy     *User
	TaskID      uint64
	Task        *Task
	StartDate   time.Time
	EndDate     *time.Time
	IsBillable  bool
	IsBilled    bool
	Description string
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
