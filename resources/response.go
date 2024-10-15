package resources

// Standard API response structure
type Response struct {
	Data  interface{} `json:"data"`  // actual response data
	Meta  interface{} `json:"meta"`  // metadata about the response, such as pagination information
	Error interface{} `json:"error"` // error details if status is "error"
}
