// Package main is the application entry point
package main

import (
	"log"
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.DBConnect()

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate databse: %v", err)
	}

	router := routes.SetupRouter()
	if err := router.Run(":" + port); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
