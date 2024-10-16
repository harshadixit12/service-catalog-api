package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harshadixit12/service-catalog-api/repository"
	"github.com/harshadixit12/service-catalog-api/resources"
	"github.com/oklog/ulid/v2"
)

// List of fields supported for filtering
var allowedFilterFields = map[string]bool{
	"name":        true,
	"description": true,
}

var allowedSortFields = map[string]bool{
	"ID":            true,
	"id":            true,
	"Name":          true,
	"name":          true,
	"created_at":    true,
	"updated_at":    true,
	"version_count": true,
}

var allowedSortOrder = map[string]bool{
	"asc":  true,
	"desc": true,
	"ASC":  true,
	"DESC": true,
}

func GetServices(c *gin.Context) {
	_, userExists := c.Get("userID")
	orgID, orgExists := c.Get("organizationID")
	if !userExists || !orgExists {
		resources.SendError(c, http.StatusUnauthorized, gin.H{"message": "User is not authorized."})
		return
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size_limit", "25"))
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page_number", "1"))
	sortField := c.DefaultQuery("sort", "ID")   // Default sort field
	sortOrder := c.DefaultQuery("order", "asc") // Default sort order
	filterField := c.DefaultQuery("filter_field", "")
	filterValue := c.DefaultQuery("filter_value", "")

	if pageNumber < 1 {
		resources.SendError(c, http.StatusBadRequest, "Invalid page_number - must be greater than 1.")
		return
	}

	if pageSize < 1 || pageSize > 100 {
		resources.SendError(c, http.StatusBadRequest, "Invalid page_size_limit - must be greater than 1 and less than 101.")
		return
	}

	if !allowedSortFields[sortField] {
		resources.SendError(c, http.StatusBadRequest, "Invalid sort_field - must be one of [id, name, created_at, updated_at, version_count].")
		return
	}

	if !allowedSortOrder[sortOrder] {
		resources.SendError(c, http.StatusBadRequest, "Invalid sort_order - must be one of [asc, desc].")
		return
	}

	if filterField != "" || filterValue != "" {
		if !allowedFilterFields[filterField] {
			resources.SendError(c, http.StatusBadRequest, "Invalid filter_field: must be one of [name, description]")
			return
		}
	}

	services, err := repository.GetServices(orgID.(int), pageSize, pageNumber, sortField, sortOrder, filterField, filterValue)

	if err != nil {
		fmt.Printf("Error loading services: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load services."})
		return
	}

	resources.SendSuccess(c, http.StatusOK, services, gin.H{"PageNumber": pageNumber, "PageSize": len(services), "PageSizeLimit": pageSize})
}

func GetServiceByID(c *gin.Context) {
	_, userExists := c.Get("userID")
	_, orgExists := c.Get("organizationID")
	if !userExists || !orgExists {
		resources.SendError(c, http.StatusUnauthorized, gin.H{"message": "User is not authorized."})
		return
	}

	serviceId := c.Param("serviceId")

	serviceULID, err := ulid.Parse(serviceId)

	if err != nil {
		fmt.Printf("Invalid service ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service ID is not valid."})
		return
	}

	service, err := repository.GetServiceByID(serviceULID.String())

	if err != nil {
		fmt.Printf("Error loading service: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load service."})
		return
	}

	resources.SendSuccess(c, http.StatusOK, service, nil)
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

	var serviceRequestInstance resources.ServiceRequestBody

	// Bind the JSON request body to the service struct
	if err := c.ShouldBindJSON(&serviceRequestInstance); err != nil {
		resources.SendError(c, http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	service := repository.Service{Name: serviceRequestInstance.Name, Description: serviceRequestInstance.Description, UserID: 1, OrganizationID: 1}

	// Todo - is uint the right choice?
	createdService, err := repository.CreateService(&service)

	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create service."})
		return
	}

	resources.SendSuccess(c, http.StatusCreated, createdService, nil)
}
