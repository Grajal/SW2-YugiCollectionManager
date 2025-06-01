package handlers

import (
	"errors"
	"fmt"
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

// DeckHandler defines the handler interface for deck-related routes.
type DeckHandler interface {
	CreateDeck(c *gin.Context)
	GetUserDecks(c *gin.Context)
	DeleteDeck(c *gin.Context)
	GetCardByDeck(c *gin.Context)
	AddCardToDeck(c *gin.Context)
	RemoveCardFromDeck(c *gin.Context)
	ExportDeckHandler(c *gin.Context)
	ImportDeckHandler(c *gin.Context)
}

type deckHandler struct {
	deckService services.DeckService
}

// NewDeckHandler creates a new instance of DeckHandler with the provided service.
func NewDeckHandler(deckService services.DeckService) DeckHandler {
	return &deckHandler{
		deckService: deckService,
	}
}

// CreateDeck handles the creation of a new deck for the authenticated user.
func (h *deckHandler) CreateDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req CreateDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	deck, err := h.deckService.CreateDeck(userID, req.Name, req.Description)
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

// GetUserDecks returns all decks associated with the authenticated user.
func (h *deckHandler) GetUserDecks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	decks, err := h.deckService.GetDecksByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch decks"})
		return
	}

	c.JSON(http.StatusOK, decks)
}

// DeleteDeck deletes a deck by its ID, ensuring it belongs to the authenticated user.
func (h *deckHandler) DeleteDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	deckIDStr := c.Param("deckId")

	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	err = h.deckService.DeleteDeck(uint(deckID), userID)
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

// GetCardByDeck returns all cards associated with a given deck.
func (h *deckHandler) GetCardByDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	deckIDStr := c.Param("deckId")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	cards, err := h.deckService.GetCardsByDeck(userID, uint(deckID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive deck cards"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

// AddCardToDeck adds a card to a deck, respecting quantity and deck constraints.
func (h *deckHandler) AddCardToDeck(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	deckIDStr := c.Param("deckId")
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

	err = h.deckService.AddCardToDeck(userID, req.CardID, uint(deckID), req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCardCopyLimitExceeded),
			errors.Is(err, services.ErrDeckLimitReached),
			errors.Is(err, services.ErrExtraDeckLimitReached):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to deck"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card added successfully"})
}

// RemoveCardFromDeck removes a specific quantity of a card from a deck.
func (h *deckHandler) RemoveCardFromDeck(c *gin.Context) {
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

	err := h.deckService.RemoveCardFromDeck(userID, uint(deckID), uint(cardID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card removed from deck"})
}

// Allows you to export a deck in .ydk format for use in clients such as EDOPro.
// according to their number in the deck. (card_ygo_id) of each card, repeated
// according to their quantity in the deck.
func (h *deckHandler) ExportDeckHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	deckIDStr := c.Param("deckId")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	ydkContent, err := h.deckService.ExportDeckAsYDK(userID, uint(deckID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=deck.ydk")
	c.Data(http.StatusOK, "text/plain", []byte(ydkContent))
}

// ImportDeckHandler imports a .ydk file into an existing deck, adding the cards accordingly.
func (h *deckHandler) ImportDeckHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	deckIDStr := c.Param("deckId")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	err = h.deckService.ImportDeckFromYDK(userID, uint(deckID), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to import deck: %v", err)})
		return
	}

	c.Status(http.StatusNoContent)
}
