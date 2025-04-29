package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetNewCard(c echo.Context) error {
	// Example response
	return c.JSON(http.StatusOK, []string{"Alice", "Bob"})
}
