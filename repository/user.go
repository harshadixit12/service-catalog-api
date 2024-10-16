package repository

import (
	"time"
)

// User represents a customer who belongs to an Organization.
type User struct {
	ID             int        `gorm:"unique;primaryKey;autoIncrement"`
	Name           string     `gorm:"type:varchar(256);not null"`
	Email          string     `gorm:"type:varchar(512);not null"`
	OrganizationID int        `gorm:"type:int;not null"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `gorm:"default null"`
	Services       []Service  `gorm:"foreignKey:UserID"`
	Versions       []Version  `gorm:"foreignKey:UserID"`
}
