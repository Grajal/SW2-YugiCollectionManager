package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

type CreateDeckRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AddCardRequest struct {
	CardID   uint `json:"card_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

type RemoveCardRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

func CreateDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req CreateDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Invalid request body"})
		return
	}

	deck, err := services.CreateDeck(userID, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, services.ErrDeckAlreadyExists) || errors.Is(err, services.ErrMaximumNumberOfDecks) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deck"})
		}
		return
	}

	c.JSON(http.StatusCreated, deck)
}

func GetUserDecks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	decks, err := services.GetDecksByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch decks"})
		return
	}

	c.JSON(http.StatusOK, decks)
}

func DeleteDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	deckIDStr := c.Param("deckId")

	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	err = services.DeleteDeck(uint(deckID), userID)
	if err != nil {
		if errors.Is(err, services.ErrDeckNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deck"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deck deleted successfully"})
}

func GetCardByDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	deckIDStr := c.Param("deckId")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	cards, err := services.GetCardsByDeck(userID, uint(deckID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive deck cards"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func AddCardToDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	deckIDStr := c.Param("deckId")
	fmt.Println("Param deckId:", c.Param("deckId"))
	fmt.Println("full path:", c.FullPath())
	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	var req AddCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := services.AddCardToDeck(userID, uint(deckID), req.CardID, req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCardNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
			return
		case errors.Is(err, services.ErrCardCopyLimitExceeded),
			errors.Is(err, services.ErrDeckLimitReached),
			errors.Is(err, services.ErrExtraDeckLimitReached):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Println("Error al añadir cartal al deck:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to deck"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func RemoveCardFromDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	deckIDStr := c.Param("deckId")
	cardIDStr := c.Param("cardId")

	deckID, err1 := strconv.ParseUint(deckIDStr, 10, 64)
	cardID, err2 := strconv.ParseUint(cardIDStr, 10, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IDs"})
		return
	}

	var req RemoveCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := services.RemoveCardFromDeck(userID, uint(deckID), uint(cardID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card removed from deck"})
}

// Allows you to export a deck in .ydk format for use in clients such as EDOPro.
// according to their number in the deck. (card_ygo_id) of each card, repeated
// according to their quantity in the deck.
func ExportDeckHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	deckIDStr := c.Param("deckId")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	ydkContent, err := services.ExportDeckAsYDK(userID, uint(deckID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=deck.ydk")
	c.Data(http.StatusOK, "text/plain", []byte(ydkContent))
}
