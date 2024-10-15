package repository

import (
	"time"
)

// User represents a customer who belongs to an Organization.
// Todo: Indices
type User struct {
	ID             int    `gorm:"unique;primaryKey;autoIncrement"`
	Name           string `gorm:"size:255;not null"`
	Email          string `gorm:"size:512;not null"`
	OrganizationID uint   `gorm:"type:bigint;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"default null"`
	Organization   Organization
}
