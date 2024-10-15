package controllers

import (
	"fmt"
	"net/http"

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
	}

	version := repository.Version{Name: versionRequestInstance.Name, ServiceID: serviceULID, UserID: 1, OrganizationID: 1}

	// Todo - is uint the right choice?
	createdVersion, err := repository.CreateVersion(&version)

	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		resources.SendError(c, http.StatusInternalServerError, gin.H{"error": "Unable to create version."})
		return
	}

	resources.SendSuccess(c, http.StatusCreated, createdVersion)
}
