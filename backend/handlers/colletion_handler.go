package handlers

import (
	"net/http"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"github.com/gin-gonic/gin"
)

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
