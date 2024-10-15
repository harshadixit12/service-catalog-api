package repository

import (
	"time"
)

// Represents an Organization.
// Todo: Indices
type Organization struct {
	ID        uint   `gorm:"unique;primaryKey;autoIncrement"`
	Name      string `gorm:"size:255;not null"` // Name of the Organization
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"default null"`
}
