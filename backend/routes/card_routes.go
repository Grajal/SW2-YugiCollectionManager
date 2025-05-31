package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterCardRoutes(rg *gin.RouterGroup, h handlers.CardHandler) {
	rg.GET("/cards/:param", h.GetCardByParam)
	rg.GET("/cards", h.GetCards)
	rg.GET("/cards/search", h.SearchCards)
}
