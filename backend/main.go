// Package main is the application entry point
package main

import (
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/routes"
)

// port is the server port, defaults to 8080 if not set in environment
var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}

	database.DBConnect()

	if err := database.DB.AutoMigrate(models.User{}, models.Card{}, models.SpellTrapCard{}, models.MonsterCard{}, models.UserCard{}, models.Deck{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	router := routes.SetupRouter()

	if err := router.Run(":" + port); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
