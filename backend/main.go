package main

import (
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/routes"
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

	router := routes.SetupRouter()

	router.Run(":" + port)
}
