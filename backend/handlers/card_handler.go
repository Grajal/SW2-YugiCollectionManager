package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

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

func (h *cardHandler) SearchCards(c *gin.Context) {
	name := c.Query("name")
	cardType := c.Query("type")
	frameType := c.Query("frameType")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cards, err := h.service.GetFilteredCards(name, cardType, frameType, limit, offset)
	if err != nil {
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
