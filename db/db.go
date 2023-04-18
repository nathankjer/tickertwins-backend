package db

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
	connString := os.Getenv("DB_CONN")
	var err error
	DB, err = gorm.Open("postgres", connString)
	if err != nil {
		panic("Failed to connect to the database.")
	}
}
