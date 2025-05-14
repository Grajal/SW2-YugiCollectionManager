package handlers

import (
	"errors"
	"net/http"

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
		if errors.Is(err, services.ErrDeckAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deck"})
		}
		return
	}

	c.JSON(http.StatusCreated, deck)
}
