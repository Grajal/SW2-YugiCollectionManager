// Package routes contains API route configuration
package routes

import (
	"strings"
	"time"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main application router
func SetupRouter() *gin.Engine {
	r := gin.Default()

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
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/", handlers.CreateUser)
			users.GET("/:username", handlers.GetUserByName)
			users.DELETE("/:username", handlers.DeleteUser)
		}

		cards := api.Group("/cards")
		cards.Use(middleware.AuthMiddleware())
		{
			cards.GET("/", handlers.GetCards)
			cards.GET("/search", handlers.SearchCards)
			cards.GET("/:param", handlers.GetOrFetchCard) // Get new card
		}
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
		}

		auth = api.Group("/auth")
		{
			auth.Use(middleware.AuthMiddleware())
			auth.GET("/me", handlers.GetCurrentUser) // Get current user
		}

		collections := api.Group("/collections")
		collections.Use(middleware.AuthMiddleware())
		{
			collections.GET("/", handlers.GetColletion) // Get collection
			collections.POST("/", handlers.AddCardToCollection)
			collections.DELETE("/:cardId", handlers.DeleteQuantityCardsFromCollcetion)
		}

		decks := api.Group("/decks")
		decks.Use(middleware.AuthMiddleware())
		{
			decks.GET("/", handlers.GetUserDecks)
			decks.POST("/", handlers.CreateDeck)
			decks.POST("/import/:deckId", handlers.ImportDeckHandler)
			decks.POST("/export/:deckId", handlers.ExportDeckHandler)
			decks.GET("/:deckId/cards", handlers.GetCardByDeck)
			decks.POST("/:deckId/cards", handlers.AddCardToDeck)
			decks.DELETE("/:deckId", handlers.DeleteDeck)
			decks.DELETE("/:deckId/cards/:cardId", handlers.RemoveCardFromDeck)
		}
	}

	return r
}
