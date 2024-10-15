package middleware

import (
	"github.com/gin-gonic/gin"
)

// Middleware function to add user details to the request context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ideal case, we'd parse user's access token from request headers, and validate it, and add
		// user info to the request

		// Add the user info to the request context
		c.Set("userID", 1)
		c.Set("organizationID", 1)

		// Continue to the next middleware or handler
		c.Next()
	}
}
