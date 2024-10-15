package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshadixit12/service-catalog-api/resources"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure all middleware have run before handling uncaught error
		c.Next()

		for _, err := range c.Errors {
			fmt.Println("Uncaught error: %w", err)
		}

		resources.SendError(c, http.StatusInternalServerError, gin.H{"message": "Something went wrong."})
	}
}
