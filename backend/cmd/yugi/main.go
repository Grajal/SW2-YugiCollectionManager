package main

import (
	"net/http"
	"os"
	"yugi/internal/handlers"

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

	e.GET("/health", handlers.HealthHandler)

	e.Logger.Fatal(e.Start(":" + port))
}
