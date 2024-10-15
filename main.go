package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/harshadixit12/service-catalog-api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(dbInstance *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
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
