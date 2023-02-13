package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}

	db.AutoMigrate()

	e := echo.New()
	e.GET("/api/helloworld", Hello)
	e.Logger.Fatal(e.Start(":1323"))
}