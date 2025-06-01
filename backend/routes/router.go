// Package routes contains API route configuration
package routes

import (
	"strings"
	"time"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main application router
func SetupRouter() *gin.Engine {
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

	api := router.Group("/api")
	RegisterAuthRoutes(api, authHandler)
	RegisterCardRoutes(api, cardHandler)
	RegisterDeckRoutes(api, deckHandler)
	RegisterStatsRoutes(api, statsHandler)
	RegisterCollectionRoutes(api, collectionHandler)

	return router
}
