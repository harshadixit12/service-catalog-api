package resources

// Represents the request body for creating a service
type ServiceRequestBody struct {
	Name        string `json:"name" binding:"required,min=1,max=255"` // Name is a string, required should be less than 255 chars long
	Description string `json:"description" binding:"max=1024"`        // Description is not required, and can be up to 1024 characters
}
