package repository

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Service represents a service in the User's organization.
// To-do
// Indices, unique name constraints
type Version struct {
	ID             ulid.ULID `gorm:"primaryKey;size:36"`
	Name           string    `gorm:"size:255;not null"`
	Description    string    `gorm:"type:text;size:1023"`
	ServiceID      ulid.ULID `gorm:"size:36;not null"`
	UserID         uint      `gorm:"type:bigint;not null"`
	OrganizationID uint      `gorm:"type:bigint;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time `gorm:"default null"`
	Service        Service
	Organization   Organization
	User           User
}
