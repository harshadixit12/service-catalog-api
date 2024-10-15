package resources

import "github.com/oklog/ulid/v2"

// Represents the request body for creating a service version
type VersionRequest struct {
	Name      string    `json:"name" binding:"required,min=1,max=255"` // Name is a string, required should be less than 255 chars long
	ServiceID ulid.ULID `json:"description" binding:"max=26"`          // ServiceID is required, and can be up to 26 characters
}
