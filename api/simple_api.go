package api

import (
	"net/http"
	"simple-crud/go/db"
	"simple-crud/go/models"

	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) error {
	var simple models.Simple	
	
	err := c.Bind(&simple)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request.")
	}

	db := db.DbManager()
	db.Create(&simple)

	return c.String(http.StatusCreated, "Record created")
}