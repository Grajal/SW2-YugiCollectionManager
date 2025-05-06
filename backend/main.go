package main

import (
	"net/http"
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/router"

	echo "github.com/labstack/echo/v4"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}

	database.DBConnect()

	if err := database.DB.AutoMigrate(models.User{}, models.Card{}, models.SpellTrapCard{}, models.MonsterCard{}, models.Collection{}, models.Deck{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	e := router.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/health", handlers.HealthHandler)

	e.Logger.Fatal(e.Start(":" + port))
}
