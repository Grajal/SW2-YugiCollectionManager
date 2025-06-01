// Package main is the application entry point
package main

import (
	"os"
	"strings"
	"time"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/routes"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.DBConnect()

	if err := database.DB.AutoMigrate(
		models.User{},
		models.Card{},
		models.SpellTrapCard{},
		models.MonsterCard{},
		models.LinkMonsterCard{},
		models.PendulumMonsterCard{},
		models.UserCard{},
		models.Deck{},
		models.DeckCard{},
	); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	userRepo := repository.NewUserRepository()
	authServices := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authServices)

	cardRepo := repository.NewCardRepository()
	cardFactory := services.NewCardFactory()
	cardService := services.NewCardService(cardRepo, cardFactory)
	cardHandler := handlers.NewCardHandler(cardService)

	deckRepo := repository.NewDeckRepository()
	deckCardRepo := repository.NewDeckCardRepository()
	deckCardService := services.NewDeckCardService(deckCardRepo)
	deckService := services.NewDeckService(deckRepo, cardService, deckCardService)
	deckHandler := handlers.NewDeckHandler(deckService)

	collectionRepo := repository.NewCollectionRepository()
	collectionService := services.NewCollectionService(collectionRepo)
	collectionHandler := handlers.NewCollectionHandler(collectionService)

	statsService := services.NewStatsService(collectionService, deckService)
	statsHandler := handlers.NewStatsHandler(statsService, deckService)

	router := gin.Default()
	allowedOrigins := "http://localhost:5173,https://sw-2-yugi-collection-manager.vercel.app"

	origins := strings.Split(allowedOrigins, ",")

	config := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	api := router.Group("/api")
	routes.RegisterAuthRoutes(api, authHandler)
	routes.RegisterCardRoutes(api, cardHandler)
	routes.RegisterDeckRoutes(api, deckHandler)
	routes.RegisterStatsRoutes(api, statsHandler)

	collections := api.Group("/collections")
	collections.Use(middleware.AuthMiddleware())
	routes.RegisterCollectionRoutes(collections, collectionHandler)

	if err := router.Run(":" + port); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
