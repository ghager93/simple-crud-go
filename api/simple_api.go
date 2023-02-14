package api

import (
	"net/http"
	"simple-crud/go/db"
	"simple-crud/go/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) error {
	name := c.FormValue("name")
	number, _ := strconv.Atoi(c.FormValue("number"))

	simple := models.Simple{Name: name, Number: number}
	
	db := db.DbManager()

	db.Create(&simple)

	c.Bind(&simple)

	return c.String(http.StatusCreated, "Record created")
}