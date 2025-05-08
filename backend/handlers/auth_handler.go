package handlers

import (
	"net/http"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := services.AuthenticateUser(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

// func LoginWithClerk(c *gin.Context) {
// 	var input struct {
// 		SessionToken string `json:"session_token" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	user, err := services.AuthenticateWithClerk(input.SessionToken)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user.Password = ""
// 	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
// }
