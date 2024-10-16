package repository

import (
	"time"
)

// Represents an Organization.
// Contains hasMany relationship with other entites like Service, Users and Version
// https://gorm.io/docs/has_many.html
type Organization struct {
	ID        int        `gorm:"unique;primaryKey;autoIncrement"`
	Name      string     `gorm:"type:varchar(256);not null"` // Name of the Organization
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `gorm:"default null"`
	Users     []User     `gorm:"foreignKey:OrganizationID"`
	Services  []Service  `gorm:"foreignKey:OrganizationID"`
	Versions  []Version  `gorm:"foreignKey:OrganizationID"`
}
