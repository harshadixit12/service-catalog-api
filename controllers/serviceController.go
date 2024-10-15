package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harshadixit12/service-catalog-api/repository"
)

// GetServices handles GET requests for /services
func GetServices(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "25"))
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("pageNumber", "1"))
	sortField := c.DefaultQuery("sort", "ID")   // Default sort field
	sortOrder := c.DefaultQuery("order", "asc") // Default sort order

	services, err := repository.GetServices(pageSize, pageNumber, sortField, sortOrder)

	if err != nil {
		fmt.Printf("Error fetching services: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch services."})
		return
	}

	// Respond with the list of services
	c.JSON(http.StatusOK, gin.H{"services": services})
}
