package handlers

import (
	"net/http"
	"strconv"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

type StatsHandler interface {
	GetCollectionStats(c *gin.Context)
	GetDeckStats(c *gin.Context)
}

type statsHandler struct {
	statsService services.StatsService
	deckService  services.DeckService
}

// Constructor
func NewStatsHandler(statsService services.StatsService, deckService services.DeckService) StatsHandler {
	return &statsHandler{
		statsService: statsService,
		deckService:  deckService,
	}
}

// GET /api/stats/collection
func (h *statsHandler) GetCollectionStats(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in cookie"})
		return
	}

	stats, err := h.statsService.CalculateCollectionStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not compute stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GET /api/stats/deck/:deckID
func (h *statsHandler) GetDeckStats(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in cookie"})
		return
	}

	deckID, err := strconv.ParseUint(c.Param("deckID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	deckCards, err := h.deckService.GetCardsByDeck(userID, uint(deckID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch deck: " + err.Error()})
		return
	}

	stats := h.statsService.CalculateDeckStats(deckCards)
	c.JSON(http.StatusOK, stats)
}
