package resources

import (
	"github.com/gin-gonic/gin"
)

// Sends a standard response, containing the result for request.
func SendSuccess(c *gin.Context, status int, data interface{}, meta interface{}) {
	c.JSON(status, Response{
		Meta:  meta,
		Data:  data,
		Error: nil,
	})
}

// Sends a standard error response.
func SendError(c *gin.Context, status int, err interface{}) {
	c.JSON(status, Response{
		Meta:  nil,
		Data:  nil,
		Error: err,
	})
}
