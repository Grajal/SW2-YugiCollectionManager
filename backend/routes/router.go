package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/gin-gonic/gin"
)

// New initializes the Echo router and returns it
// func New() *echo.Echo {
// 	e := echo.New()

// 	// Example route
// 	e.GET("/", func(c echo.Context) error {
// 		return c.String(http.StatusOK, "Hello from Echo!")
// 	})

// 	// You can group routes or import handler packages here
// 	e.GET("/health", handlers.HealthHandler)
// 	e.GET("/getNewCard", handlers.GetNewCard)

// 	return e
// }

func SetupRouter() *gin.Engine {
	r := gin.Default()

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
			cards.GET("/getNewCard", handlers.GetNewCard)
		}
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
		}
	}

	return r
}
