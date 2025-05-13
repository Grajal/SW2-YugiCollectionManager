package handlers

import (
	"net/http"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

const apiBaseURL = "https://db.ygoprodeck.com/api/v7/cardinfo.php"

// GetNewCard handles the creation of a new card from the YGOProDeck API
func GetNewCard(c *gin.Context) {
	cardName := c.Query("name")
	if cardName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'name' query parameter"})
		return
	}

	// Get card data from the API
	cardData, err := services.GetCardByName(cardName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create the card in the database
	createdCard, err := services.CreateCard(cardData)
	if err != nil {
		// Handle error from database or if card already exists
		if err.Error() == "card already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Respond with the created card
	c.JSON(http.StatusOK, createdCard)
}
