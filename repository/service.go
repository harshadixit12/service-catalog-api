package repository

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Service represents a service in the User's organization.
// Could improve: unique name constraints for services within an org might be needed.
// Contains hasMany relationship with Version
// https://gorm.io/docs/has_many.html
type Service struct {
	ID             string     `gorm:"primaryKey;type:char(36)"`   // ULID as the primary key - size 26 chars
	Name           string     `gorm:"type:varchar(256);not null"` // Name of the service
	Description    string     `gorm:"type:varchar(1024)"`         // Description about the service
	UserID         int        `gorm:"type:int;not null"`          // ID of user who created the service
	OrganizationID int        `gorm:"type:int;not null"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `gorm:"default null"`
	VersionCount   int        `gorm:"type:int;not null;default:0"`
	Versions       []Version  `gorm:"foreignKey:ServiceID"`
}

// BeforeCreate GORM hook to generate a ULID before inserting a new service
func (s *Service) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = ulid.Make().String()
	return
}

// Creates a Service and inserts into DB
func CreateService(service *Service) (*Service, error) {
	value := DBInstance.Create(service)
	if value.Error != nil {
		return nil, value.Error
	}
	return service, nil
}

// Loads all non-deleted services and returns an array.
func GetServices(organizationID int, pageSize int, pageNo int, sortField string, sortOrder string, filterField string, filterValue string) ([]Service, error) {
	var services []Service

	tx := DBInstance.Session(&gorm.Session{})

	if filterField != "" && filterValue != "" {
		tx = tx.Where(filterField+" = ?", filterValue)
	}

	if err := tx.Where("deleted_at IS NULL").Where("organization_id = ?", organizationID).Order(sortField + " " + sortOrder).Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

// Loads a single service by ID
func GetServiceByID(serviceId string) (*Service, error) {
	var service Service

	tx := DBInstance.Session(&gorm.Session{})

	if err := tx.First(&service, "id=?", serviceId).Error; err != nil {
		return nil, err
	}

	return &service, nil
}
