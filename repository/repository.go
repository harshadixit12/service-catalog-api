package repository

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initializes the database connection - SQLite for the first iteration of development.
// We will make use of GORM as the ORM.
func InitDatabase() (*gorm.DB, error) {
	// Open a connection to the SQLite database file (it will be created if it doesn't exist)
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// creates tables if they don't exist)
	err = db.AutoMigrate(&Organization{}, &User{}, &Service{}, &Version{})
	if err != nil {
		fmt.Printf("Error automigrating schema: %v\n", err)
		return nil, err
	}

	return db, nil
}
