package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harshadixit12/service-catalog-api/repository"
	"github.com/harshadixit12/service-catalog-api/resources"
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

func CreateService(c *gin.Context) {
	// Load user and organization IDs from auth
	userID, userExists := c.Get("userID")
	orgID, orgExists := c.Get("organizationID")
	fmt.Printf("%s %s", userID, orgID)
	if !userExists || !orgExists {
		resources.SendError(c, http.StatusUnauthorized, gin.H{"message": "User is not authorized."})
		return
	}

	var serviceRequestInstance resources.ServiceRequest

	// Bind the JSON request body to the service struct
	if err := c.ShouldBindJSON(&serviceRequestInstance); err != nil {
		resources.SendError(c, http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	service := repository.Service{Name: serviceRequestInstance.Name, Description: serviceRequestInstance.Description, UserID: 1, OrganizationID: 1}

	// Todo - is uint the right choice?
	createdService, err := repository.CreateService(&service)

	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create service."})
		return
	}

	resources.SendSuccess(c, http.StatusCreated, createdService)
}
