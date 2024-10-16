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

func CreateVersion(c *gin.Context) {
	// Load user and organization IDs from auth
	userID, userExists := c.Get("userID")
	orgID, orgExists := c.Get("organizationID")
	fmt.Printf("%s %s", userID, orgID)
	if !userExists || !orgExists {
		resources.SendError(c, http.StatusUnauthorized, gin.H{"message": "User is not authorized."})
		return
	}

	serviceId := c.Param("serviceId")

	serviceULID, err := ulid.Parse(serviceId)
	if err != nil {
		resources.SendError(c, http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var versionRequestInstance resources.VersionRequest

	// Bind the JSON request body to the service struct
	if err := c.ShouldBindJSON(&versionRequestInstance); err != nil {
		resources.SendError(c, http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	version := repository.Version{Name: versionRequestInstance.Name, ServiceID: serviceULID, UserID: 1, OrganizationID: 1}

	createdVersion, err := repository.CreateVersion(&version)

	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		resources.SendError(c, http.StatusInternalServerError, gin.H{"error": "Unable to create version."})
		return
	}

	resources.SendSuccess(c, http.StatusCreated, createdVersion, nil)
}

func GetServiceVersions(c *gin.Context) {
	serviceULID, err := ulid.Parse(c.Param("serviceId"))

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size_limit", "25"))
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page_number", "1"))

	if pageNumber < 1 {
		resources.SendError(c, http.StatusBadRequest, "Invalid page_number - must be greater than 1.")
		return
	}

	if pageSize < 1 || pageSize > 100 {
		resources.SendError(c, http.StatusBadRequest, "Invalid page_size_limit - must be greater than 1 and less than 101.")
		return
	}

	if err != nil {
		resources.SendError(c, http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	version := repository.Version{ServiceID: serviceULID}
	versions, err := repository.GetServiceVersions(version, pageNumber, pageSize)

	if err != nil {
		fmt.Printf("Error fetching services: %v\n", err)
		resources.SendError(c, http.StatusInternalServerError, gin.H{"message": "Unable to fetch service versions."})
		return
	}

	resources.SendSuccess(c, http.StatusOK, versions, gin.H{"PageNumber": pageNumber, "PageSize": len(versions), "PageSizeLimit": pageSize})
}
