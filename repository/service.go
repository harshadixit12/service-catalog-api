package repository

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Service represents a service in the User's organization.
// Todo: Indices, unique name constraints
type Service struct {
	ID             ulid.ULID `gorm:"primaryKey;size:26"`   // ULID as the primary key - size 26 chars
	Name           string    `gorm:"size:255;not null"`    // Name of the service
	Description    string    `gorm:"type:text;size:1023"`  // Description about the service
	UserID         uint      `gorm:"type:bigint;not null"` // ID of user who created the service
	OrganizationID uint      `gorm:"type:bigint;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time `gorm:"default null"`
	Organization   Organization
	User           User
}
