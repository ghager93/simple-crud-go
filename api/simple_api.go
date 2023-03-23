package api

import (
	"net/http"
	"simple-crud/go/db"
	"simple-crud/go/models"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func Create(c echo.Context) error {
	var simple models.Simple	
	
	if err := c.Bind(&simple); err != nil {
		return c.String(http.StatusBadRequest, "Bad request.")
	}

	validate := validator.New()
	if err := validate.Struct(&simple); err != nil {
		return c.String(http.StatusBadRequest, "Bad request.")
	}

	db := db.DbManager()
	db.Create(&simple)

	return c.String(http.StatusCreated, "Record created")
}

func GetAll(c echo.Context) error {
	var simples []models.Simple

	db := db.DbManager()
	
	if err := db.Find(&simples).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Error accessing Database.")
	}
	
	return c.JSON(http.StatusOK, &simples)
}

func Get(c echo.Context) error {
	var simple models.Simple

	db := db.DbManager()

	if err := db.First(&simple, c.Param("id")).Error; err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	return c.JSON(http.StatusOK, &simple)
}

func Delete(c echo.Context) error {
	var simple models.Simple

	db := db.DbManager()

	if err := db.Clauses(clause.Returning{}).Delete(&simple, c.Param("id")).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Error connecting to database.")
	}

	if simple.ID == 0 {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	return c.JSON(http.StatusOK, &simple)
}