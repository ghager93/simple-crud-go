package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"simple-crud/go/models"
)

var db *gorm.DB
var err error

func Init() {
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}

	db.AutoMigrate(&models.Simple{})
}

func DbManager() *gorm.DB {
	return db
}