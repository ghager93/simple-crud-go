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
	db.Init()

	e := echo.New()
	e.GET("/api/helloworld", Hello)
	e.POST("/api/simple", api.Create)
	e.Logger.Fatal(e.Start(":1323"))
}