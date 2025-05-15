package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

type CreateDeckRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
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
	deckIDStr := c.Param("id")

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
