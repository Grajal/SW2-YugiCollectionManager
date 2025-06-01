package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddCardInput defines the structure for adding cards to the collection.
type AddCardInput struct {
	CardID   uint `json:"card_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

// DeleteCardInput defines the structure for deleting a quantity of cards.
type DeleteCardInput struct {
	Quantity int `json:"quantity" binding:"required"`
}

// CollectionHandler defines the interface for collection-related HTTP operations.
type CollectionHandler interface {
	GetCollection(c *gin.Context)
	GetCollectionCard(c *gin.Context)
	AddCardToCollection(c *gin.Context)
	DeleteQuantityFromCollection(c *gin.Context)
}

type collectionHandler struct {
	service services.CollectionService
}

// Constructor
func NewCollectionHandler(service services.CollectionService) CollectionHandler {
	return &collectionHandler{service: service}
}

// GET /collection
func (h *collectionHandler) GetCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	collection, err := h.service.GetUserCollection(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"collection": collection})
}

func (h *collectionHandler) GetCollectionCard(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	cardIDParam := c.Param("cardId")
	cardIDUint64, err := strconv.ParseUint(cardIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}
	cardID := uint(cardIDUint64)

	// Llamar al servicio
	userCard, err := h.service.GetUserCard(userID, cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Card not found in collection"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve card"})
		return
	}

	// Devolver la carta si se encuentra
	c.JSON(http.StatusOK, userCard)
}

// POST /collection
func (h *collectionHandler) AddCardToCollection(c *gin.Context) {
	var input AddCardInput
	if err := c.ShouldBindJSON(&input); err != nil || input.CardID == 0 || input.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.MustGet("user_id").(uint)

	if err := h.service.AddCardToCollection(userID, input.CardID, input.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card added to collection successfully"})
}

// DELETE /collection/:cardId/quantity
func (h *collectionHandler) DeleteQuantityFromCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	cardIDParam := c.Param("cardId")
	cardIDUint, err := strconv.ParseUint(cardIDParam, 10, 64)
	if err != nil || cardIDUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	var input DeleteCardInput
	if err := c.ShouldBindJSON(&input); err != nil || input.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	if err := h.service.DecreaseCardQuantity(userID, uint(cardIDUint), input.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card quantity updated or removed successfully"})
}
