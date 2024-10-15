package repository

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Service represents a service in the User's organization.
// Todo: Indices, unique name constraints
type Service struct {
	ID             ulid.ULID `gorm:"primaryKey;size:26"`  // ULID as the primary key - size 26 chars
	Name           string    `gorm:"size:255;not null"`   // Name of the service
	Description    string    `gorm:"type:text;size:1023"` // Description about the service
	UserID         int       `gorm:"type:int;not null"`   // ID of user who created the service
	OrganizationID int       `gorm:"type:int;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"default null"`
	Organization   Organization
	User           User
}

// BeforeCreate GORM hook to generate a ULID before inserting a new service
func (s *Service) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = ulid.Make()
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

// Loads all services and returns
func GetServices(pageSize int, pageNo int, sortField string, sortOrder string) ([]Service, error) {
	var services []Service
	if err := DBInstance.Order(sortField + " " + sortOrder).Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}
