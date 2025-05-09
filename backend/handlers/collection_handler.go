package handlers

import (
	"net/http"
	"strconv"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"github.com/gin-gonic/gin"
)

// AddCardInput defines the structure for the input payload when adding a card to the collection.
// It requires a CardID and a Quantity, both of which are mandatory fields.
type AddCardInput struct {
	CardID   uint `json:"card_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

// GetColletion retrieves the user's card collection.
// It first checks if the user is authenticated by extracting the user ID from the context.
// If the user is not authenticated, it returns a 401 Unauthorized response.
// Otherwise, it fetches the collection using the service layer and returns it as a JSON response.
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

// AddCardToCollection adds a card to the user's collection.
// It validates the input payload to ensure the CardID and Quantity are valid.
// If the user is not authenticated, it returns a 401 Unauthorized response.
// If the input is invalid or the service layer fails, appropriate error responses are returned.
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

// DeleteCardFromCollection removes a card from the user's collection.
// It validates the card ID from the URL parameter and ensures the user is authenticated.
// If the card ID is invalid or the service layer fails, appropriate error responses are returned.
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
