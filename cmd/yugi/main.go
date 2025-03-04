package main

import (
	"net/http"
	"os"

	echo "github.com/labstack/echo/v4"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":" + port))
}
