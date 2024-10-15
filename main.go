package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/harshadixit12/service-catalog-api/middleware"
	"github.com/harshadixit12/service-catalog-api/repository"
	"github.com/harshadixit12/service-catalog-api/resources"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(dbInstance *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		resources.SendSuccess(c, http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}

func main() {
	dbInstance, err := repository.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	fmt.Println("SQLite database initialized successfully")

	router := setupRouter(dbInstance)

	router.Run("localhost:8080")
}
