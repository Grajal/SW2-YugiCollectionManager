package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCardRoutes(rg *gin.RouterGroup, h handlers.CardHandler) {
	rg = rg.Group("/cards")
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/:param", h.GetCardByParam)
	rg.GET("/", h.GetCards)
	rg.GET("/search", h.SearchCards)
}
