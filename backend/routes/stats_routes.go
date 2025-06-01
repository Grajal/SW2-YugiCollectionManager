package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterStatsRoutes(rg *gin.RouterGroup, h handlers.StatsHandler) {
	rg = rg.Group("/stats")
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/collection", h.GetCollectionStats)
	rg.GET("/deck/:deckID", h.GetDeckStats)
}
