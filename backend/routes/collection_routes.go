package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterCollectionRoutes(rg *gin.RouterGroup, h handlers.CollectionHandler) {
	rg.GET("/", h.GetCollection)
	rg.POST("/", h.AddCardToCollection)
	rg.DELETE("/:cardId", h.DeleteQuantityFromCollection)
}
