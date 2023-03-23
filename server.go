package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"simple-crud/go/api"
	"simple-crud/go/db"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func main() {
	db.Init("app.db")

	e := echo.New()
	e.GET("/api/helloworld", Hello)
	e.GET("/api/simple", api.GetAll)
	e.GET("/api/simple/:id", api.Get)
	e.POST("/api/simple", api.Create)
	e.DELETE("/api/simple/:id", api.Delete)
	e.Logger.Fatal(e.Start(":1323"))
}