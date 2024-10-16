package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/harshadixit12/service-catalog-api/controllers"
	"github.com/harshadixit12/service-catalog-api/middleware"
	"github.com/harshadixit12/service-catalog-api/repository"
	"github.com/harshadixit12/service-catalog-api/resources"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		resources.SendSuccess(c, http.StatusOK, gin.H{"message": "pong"}, nil)
	})

	r.GET("/services", controllers.GetServices)
	r.POST("/services", controllers.CreateService)
	r.GET("/services/:serviceId", controllers.GetServiceByID)
	r.GET("/services/:serviceId/versions", controllers.GetServiceVersions)
	r.POST("/services/:serviceId/versions", controllers.CreateVersion)
	return r
}

func main() {
	// Set timezone as UTC so we store timestamps in UTC in the database using gorm
	utcTime, _ := time.LoadLocation("UTC")

	_, err := repository.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	fmt.Printf("SQLite database initialized successfully at: %s", utcTime.String())

	router := setupRouter()

	router.Run("localhost:8080")
}
