package models

import (
	"time"
)

type TaskList struct {
	ID           uint64
	Name         string
	Description  string
	AllowedUsers []*User
	MilestoneID  uint64     `doc:"This task is belong to which milestone"`
	Milestone    *Milestone ``
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}