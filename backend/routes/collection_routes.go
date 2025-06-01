package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCollectionRoutes(rg *gin.RouterGroup, h handlers.CollectionHandler) {
	rg = rg.Group("/collections")
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/", h.GetCollection)
	rg.GET("/:cardId", h.GetCollectionCard)
	rg.POST("/", h.AddCardToCollection)
	rg.DELETE("/:cardId", h.DeleteQuantityFromCollection)
}
