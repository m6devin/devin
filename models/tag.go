package models

import (
	"time"
)

type Tag struct {
	ID          uint64
	Title       string
	CompanyID   uint64 `doc:"This tag created on this company. All tags of a company will be shared in all of its projects and other modules"`
	Company     *User
	CreatedByID uint64 `doc:"Who creates this tag?"`
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
