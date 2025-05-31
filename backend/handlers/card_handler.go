package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/client"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

// CardHandler defines the interface for handling HTTP requests related to cards.
type CardHandler interface {
	GetCardByParam(c *gin.Context)
	GetCards(c *gin.Context)
	SearchCards(c *gin.Context)
}

type cardHandler struct {
	service services.CardService
}

func NewCardHandler(service services.CardService) CardHandler {
	return &cardHandler{service}
}

// GetCardByParam handles GET requests to retrieve a card by ID (numeric) or name (string).
// Example routes:
// - GET /cards/42 → by ID
// - GET /cards/Dark%20Magician → by name
// Returns 200 with the card if found, 404 if not, and 400 if param is missing.
func (h *cardHandler) GetCardByParam(c *gin.Context) {
	param := strings.TrimSpace(c.Param("param"))
	if param == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter (ID or name)"})
		return
	}

	if parsedID, err := strconv.ParseUint(param, 10, 64); err == nil {
		id := uint(parsedID)
		card, err := h.service.GetCardByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, card)
		return
	}

	card, err := h.service.GetCardByName(param)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}

// GetCards handles GET requests to retrieve a paginated list of all cards.
// Query params:
// - limit (default: 20): max number of results
// - offset (default: 0): number of results to skip
// Returns 200 with total count and array of cards, or 500 on error.
func (h *cardHandler) GetCards(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cards, err := h.service.GetCards(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards"})
		return
	}

	total, err := h.service.CountAllCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count all cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalCards": total,
		"cards":      cards,
	})
}

// SearchCards handles GET requests to retrieve cards that match optional filters.
// Query params:
// - name: partial or full name of the card
// - type: card type (e.g. "Spell Card", "Normal Monster")
// - frameType: card frame type (e.g. "normal", "link", "pendulum")
// - limit (default: 20): max number of results
// - offset (default: 0): number of results to skip
// Returns 200 with total count and results, 400 if the API call was invalid,
// or 500 if an internal error occurred.
func (h *cardHandler) SearchCards(c *gin.Context) {
	name := c.Query("name")
	cardType := c.Query("type")
	frameType := c.Query("frameType")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cards, err := h.service.GetFilteredCards(name, cardType, frameType, limit, offset)
	if err != nil {
		if errors.Is(err, client.ErrBadRequestFromAPI) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search term"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards from database"})
		return
	}

	total, err := h.service.CountFilteredCards(name, cardType, frameType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count filtered cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalCards": total,
		"cards":      cards,
	})
}
