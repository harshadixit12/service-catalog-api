package repository

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Service represents a service in the User's organization.
// To-do
// Indices, unique name constraints
type Version struct {
	ID             ulid.ULID `gorm:"primaryKey;size:36"`
	Name           string    `gorm:"size:255;not null"`
	ServiceID      ulid.ULID `gorm:"size:36;not null"`
	UserID         int       `gorm:"type:int;not null"`
	OrganizationID int       `gorm:"type:int;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"default null"`
	Service        Service
	Organization   Organization
	User           User
}

// BeforeCreate GORM hook to generate a ULID before inserting a new service
func (v *Version) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = ulid.Make()
	return
}

// Creates a Service Version and inserts into DB
func CreateVersion(version *Version) (*Version, error) {
	value := DBInstance.Create(version)
	if value.Error != nil {
		return nil, value.Error
	}
	return version, nil
}

func GetServiceVersions(dbInstance *gorm.DB, version Version) ([]Version, error) {
	var versions []Version
	value := dbInstance.Where("Service_ID", version.ServiceID.String()).Find(&versions)

	if value.Error != nil {
		return nil, value.Error
	}

	return versions, nil
}
