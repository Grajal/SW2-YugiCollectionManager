package handlers

import (
	"net/http"
	"strconv"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"github.com/gin-gonic/gin"
)

type AddCardInput struct {
	CardID   uint `json:"card_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

func GetColletion(c *gin.Context) {
	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticate"})
		return
	}

	collection, err := services.GetCollectionByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"collection": collection})
}

func AddCardToCollection(c *gin.Context) {
	var input AddCardInput
	if err := c.ShouldBindJSON(&input); err != nil || input.CardID == 0 || input.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := services.AddCardToCollection(userID, input.CardID, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card added to collection successfully"})
}

func DeleteCardFromCollection(c *gin.Context) {
	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	cardIDParam := c.Param("card_id")
	cardID, err := strconv.ParseUint(cardIDParam, 10, 64)
	if err != nil || cardID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	if err := services.DeleteCardFromCollection(userID, uint(cardID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card from collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deleted from collection successfully"})
}
