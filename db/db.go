package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
	connString := os.Getenv("DB_CONN")
	if connString == "" {
		log.Fatal("DB_CONN environment variable not set")
	}

	var err error
	DB, err = connectWithRetries(connString, 5)
	if err != nil {
		log.Fatalf("Could not connect to the database after multiple retries: %v", err)
	}

	log.Println("Successfully connected to the database")
}

func connectWithRetries(connString string, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open("postgres", connString)
		if err == nil {
			return db, nil
		}

		log.Printf("Failed to connect to the database (attempt %d of %d): %v", i, maxRetries, err)
		time.Sleep(time.Duration(i) * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to the database after %d attempts", maxRetries)
}
