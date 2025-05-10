package handlers

import (
	"net/http"

	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

var input struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "Login succesful", "user": user})
}

func LoginWithGoogle(c *gin.Context) {
	var input struct {
		IDToken string `json:"id_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	user, err := services.AuthenticateWithGoogle(input.IDToken, clientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}
