package db

import (
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"simple-crud/go/models"
	"simple-crud/go/utils/projectpath"
)

var db *gorm.DB
var err error

func Init(dbName string) {
	dbPath := filepath.Join(projectpath.Root, dbName)
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}

	db.AutoMigrate(&models.Simple{})
}

func DbManager() *gorm.DB {
	return db
}

func DropTables() {
	db.Migrator().DropTable(&models.Simple{})
}