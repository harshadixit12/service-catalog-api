package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harshadixit12/service-catalog-api/repository"
	"github.com/harshadixit12/service-catalog-api/resources"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRepository(t *testing.T) *gorm.DB {
	// Create an in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory SQLite database: %v", err)
	}

	// Migrate the schema (create tables in the in-memory database)
	err = db.AutoMigrate(&repository.Service{}, &repository.Version{}, &repository.Organization{}, &repository.User{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Create a dummy organisation and user, and ignore errors if already present
	orgResult := db.Create(&repository.Organization{Name: "Poppy Corp."})
	if orgResult.Error != nil {
		t.Fatalf("Error creating org: %v\n", orgResult.Error)
	}

	userResult := db.Create(&repository.User{Name: "Poppy Corp.", Email: "user_1@poppycorp.com", OrganizationID: 1})
	if userResult.Error != nil {
		t.Fatalf("Error creating user: %v\n", userResult.Error)
	}

	return db
}

func TestPingRoute(t *testing.T) {
	dbInstance := setupTestRepository(t)
	repository.DBInstance = dbInstance
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf(`Expected HTTP 200 OK from GET /ping, received %d instead`, w.Code)
	}

	var jsonResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Contains(t, jsonResponse, "data", "Response should contain 'data'")
	assert.Contains(t, jsonResponse, "meta", "Response should contain 'meta'")
	assert.Contains(t, jsonResponse, "error", "Response should contain 'error'")

	assert.Equal(t, "pong", jsonResponse["data"].(map[string]interface{})["message"])
	assert.Equal(t, nil, jsonResponse["meta"])
	assert.Equal(t, nil, jsonResponse["error"])
}

func TestServiceCreation(t *testing.T) {
	dbInstance := setupTestRepository(t)
	repository.DBInstance = dbInstance
	router := setupRouter()

	w := httptest.NewRecorder()

	serviceRequest := resources.ServiceRequestBody{Name: "New Test service", Description: "Service used in tests."}

	requestBody, _ := json.Marshal(serviceRequest)
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(requestBody))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf(`Expected HTTP 201 Created from POST /services, received %d instead`, w.Code)
	}

	var jsonResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Contains(t, jsonResponse, "data", "Response should contain 'data'")
	assert.Contains(t, jsonResponse, "meta", "Response should contain 'meta'")
	assert.Contains(t, jsonResponse, "error", "Response should contain 'error'")

	receivedService := jsonResponse["data"].(map[string]interface{})
	assert.Contains(t, receivedService, "ID", "Response should contain 'ID'")
	assert.Contains(t, receivedService, "Name", "Response should contain 'Name'")
	assert.Contains(t, receivedService, "Description", "Response should contain 'ID'")
	assert.Contains(t, receivedService, "VersionCount", "Response should contain 'VersionCount'")

	assert.Equal(t, serviceRequest.Name, receivedService["Name"])
	assert.Equal(t, serviceRequest.Description, receivedService["Description"])
	assert.Equal(t, 0, int(receivedService["VersionCount"].(float64)))
}

func TestGetServiceByID(t *testing.T) {
	dbInstance := setupTestRepository(t)
	repository.DBInstance = dbInstance
	router := setupRouter()

	createdService := repository.Service{Name: "New Test service", Description: "Service used in tests."}
	result, err := repository.CreateService(&createdService)

	if err != nil {
		t.Fatalf(`Failed to create service in DB for test`)
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/services/"+result.ID, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf(`Expected HTTP 200 Created from GET /services/:id, received %d instead`, w.Code)
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Contains(t, jsonResponse, "data", "Response should contain 'data'")
	assert.Contains(t, jsonResponse, "meta", "Response should contain 'meta'")
	assert.Contains(t, jsonResponse, "error", "Response should contain 'error'")

	receivedService := jsonResponse["data"].(map[string]interface{})
	assert.Contains(t, receivedService, "ID", "Response should contain 'ID'")
	assert.Contains(t, receivedService, "Name", "Response should contain 'Name'")
	assert.Contains(t, receivedService, "Description", "Response should contain 'ID'")
	assert.Contains(t, receivedService, "VersionCount", "Response should contain 'VersionCount'")

	assert.Equal(t, createdService.Name, receivedService["Name"])
	assert.Equal(t, createdService.Description, receivedService["Description"])
}

func TestCreateServiceVersion(t *testing.T) {
	dbInstance := setupTestRepository(t)
	repository.DBInstance = dbInstance
	router := setupRouter()

	createdService := repository.Service{Name: "New Test service", Description: "Service used in tests."}
	result, err := repository.CreateService(&createdService)

	if err != nil {
		t.Fatalf(`Failed to create service in DB for test`)
	}

	w := httptest.NewRecorder()

	versionRequest := resources.VersionRequestBody{Name: "v1.0.0"}
	versionRequestBody, _ := json.Marshal(versionRequest)
	req, _ := http.NewRequest("POST", "/services/"+result.ID+"/versions", bytes.NewBuffer(versionRequestBody))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf(`Expected HTTP 201 Created from POST /services/:id/versions, received %d instead`, w.Code)
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Contains(t, jsonResponse, "data", "Response should contain 'data'")
	assert.Contains(t, jsonResponse, "meta", "Response should contain 'meta'")
	assert.Contains(t, jsonResponse, "error", "Response should contain 'error'")

	receivedService := jsonResponse["data"].(map[string]interface{})
	assert.Contains(t, receivedService, "ID", "Response should contain 'ID'")
	assert.Contains(t, receivedService, "Name", "Response should contain 'Name'")

	assert.Equal(t, versionRequest.Name, receivedService["Name"])
}

func TestGetServiceList(t *testing.T) {
	dbInstance := setupTestRepository(t)
	repository.DBInstance = dbInstance
	router := setupRouter()

	serviceFirst := repository.Service{Name: "New Test service - 1", Description: "Service used in tests."}
	serviceSecond := repository.Service{Name: "New Test service - 2", Description: "Service used in tests."}
	_, errFirst := repository.CreateService(&serviceFirst)
	_, errSecond := repository.CreateService(&serviceSecond)

	if errFirst != nil || errSecond != nil {
		t.Fatalf(`Failed to create service in DB for test`)
	}

	w := httptest.NewRecorder()

	// We will test the pagination by specifying page size as 1
	req, _ := http.NewRequest("GET", "/services?page_size=1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf(`Expected HTTP 200 Created from GET /services, received %d instead`, w.Code)
	}

	var jsonResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Contains(t, jsonResponse, "data", "Response should contain 'data'")
	assert.Contains(t, jsonResponse, "meta", "Response should contain 'meta'")
	assert.Contains(t, jsonResponse, "error", "Response should contain 'error'")

	jsonMeta := jsonResponse["meta"].(map[string]interface{})

	assert.Contains(t, jsonMeta, "PageSizeLimit", "Response Meta should contain 'PageSizeLimit'")
	assert.Contains(t, jsonMeta, "PageSize", "Response Meta should contain 'PageSize'")
	assert.Contains(t, jsonMeta, "PageNumber", "Response Meta should contain 'PageNumber'")
	assert.Equal(t, 1, int(jsonMeta["PageNumber"].(float64)))
}

func TestGetServiceVersionsList(t *testing.T) {
	dbInstance := setupTestRepository(t)
	repository.DBInstance = dbInstance
	router := setupRouter()

	service := repository.Service{Name: "New Test service - 1", Description: "Service used in tests."}

	createdService, err := repository.CreateService(&service)

	if err != nil {
		t.Fatalf(`Failed to create service in DB for test`)
	}

	version := repository.Version{Name: "v1.0.0", ServiceID: createdService.ID}
	secondVersion := repository.Version{Name: "v2.0.0", ServiceID: createdService.ID}

	_, err = repository.CreateVersion(&version)

	if err != nil {
		t.Fatalf(`Failed to create first version in DB for test`)
	}

	_, err = repository.CreateVersion(&secondVersion)

	if err != nil {
		t.Fatalf(`Failed to create second version in DB for test`)
	}

	w := httptest.NewRecorder()

	// We will test the pagination by specifying page size as 1
	req, _ := http.NewRequest("GET", "/services/"+createdService.ID+"/versions?page_size=1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf(`Expected HTTP 200 Created from GET /services/:id/versions, received %d instead`, w.Code)
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Contains(t, jsonResponse, "data", "Response should contain 'data'")
	assert.Contains(t, jsonResponse, "meta", "Response should contain 'meta'")
	assert.Contains(t, jsonResponse, "error", "Response should contain 'error'")

	jsonMeta := jsonResponse["meta"].(map[string]interface{})

	assert.Contains(t, jsonMeta, "PageSizeLimit", "Response Meta should contain 'PageSizeLimit'")
	assert.Contains(t, jsonMeta, "PageSize", "Response Meta should contain 'PageSize'")
	assert.Contains(t, jsonMeta, "PageNumber", "Response Meta should contain 'PageNumber'")
	assert.Equal(t, 1, int(jsonMeta["PageNumber"].(float64)))
}
