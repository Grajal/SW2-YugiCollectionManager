// Package routes contains API route configuration
package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main application router
// Define todas las rutas de la API agrupadas por funcionalidad:
// - /api/users: Gestión de usuarios (crear, obtener, eliminar)
// - /api/cards: Gestión de cartas (obtener nueva carta)
// - /api/auth: Autenticación (login normal y con Clerk)
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/", handlers.CreateUser)
			users.GET("/:username", handlers.GetUserByName)
			users.DELETE("/:username", handlers.DeleteUser)
		}

		cards := api.Group("/cards")
		{
			cards.GET("/getNewCard", handlers.GetNewCard) // Get new card
		}
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("register", handlers.Register)
		}

		auth = api.Group("/auth")
		{
			auth.Use(middleware.AuthMiddleware())
			auth.GET("/me", handlers.GetCurrentUser) // Get current user
		}

		collections := api.Group("/collection")
		collections.Use(middleware.AuthMiddleware())
		{
			collections.GET("/", handlers.GetColletion) // Get collection
		}
	}

	return r
}
