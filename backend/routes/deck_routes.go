package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterDeckRoutes(rg *gin.RouterGroup, h handlers.DeckHandler) {
	rg = rg.Group("/decks")
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/", h.GetUserDecks)
	rg.POST("/", h.CreateDeck)
	rg.POST("/import/:deckId", h.ImportDeckHandler)
	rg.POST("/export/:deckId", h.ExportDeckHandler)
	rg.GET("/:deckId/cards", h.GetCardByDeck)
	rg.POST("/:deckId/cards", h.AddCardToDeck)
	rg.DELETE("/:deckId", h.DeleteDeck)
	rg.DELETE("/:deckId/cards/:cardId", h.RemoveCardFromDeck)
}
