package repository

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Service represents a service in the User's organization.
// Unique name constraint might be useful for <Name, ServiceId>
type Version struct {
	ID             string     `gorm:"primaryKey;type:char(36)"`
	Name           string     `gorm:"type:varchar(256);not null"`
	ServiceID      ulid.ULID  `gorm:"type:char(36);not null"`
	UserID         int        `gorm:"type:int;not null"`
	OrganizationID int        `gorm:"type:int;not null"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `gorm:"default null"`
}

// BeforeCreate GORM hook to generate a ULID before inserting a new service
func (v *Version) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = ulid.Make().String()
	return
}

// Creates a Service Version and inserts into DB, also updates the version count
func CreateVersion(version *Version) (*Version, error) {
	// Use a transaction to keep version count in Service consistent.
	err := DBInstance.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(version).Error; err != nil {
			// Return error to rollback
			return err
		}

		// Increment the version_count in the service
		if err := tx.Model(&Service{}).
			Where("id = ?", version.ServiceID.String()).
			Update("version_count", gorm.Expr("version_count + ?", 1)).
			Error; err != nil {
			// If error, return to rollback
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return version, nil
}

func GetServiceVersions(version Version, pageNumber int, pageSize int) ([]Version, error) {
	var versions []Version
	tx := DBInstance.Session(&gorm.Session{})

	value := tx.Where("deleted_at IS NULL").Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&versions)

	if value.Error != nil {
		return nil, value.Error
	}

	return versions, nil
}
