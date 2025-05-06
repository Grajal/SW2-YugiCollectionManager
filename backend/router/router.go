package router

import (
	"net/http"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"

	"github.com/labstack/echo/v4"
)

// New initializes the Echo router and returns it
func New() *echo.Echo {
	e := echo.New()

	// Example route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Echo!")
	})

	// You can group routes or import handler packages here
	e.GET("/health", handlers.HealthHandler)
	e.GET("/getNewCard", handlers.GetNewCard)

	return e
}
