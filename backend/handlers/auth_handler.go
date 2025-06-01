package handlers

import (
	"net/http"
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthHandler defines the interface for authentication-related endpoints.
type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetCurrentUser(c *gin.Context)
}

type authHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler with the given AuthService.
func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// Login handles user login, validates credentials, generates a JWT token,
// and sets it as a cookie in the response.
func (h *authHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.authService.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	domain := os.Getenv("COOKIE_DOMAIN")

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   domain,
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user})
}

// Register handles new user registration, validates the input,
// and creates a new user if username/email is not already taken.
func (h *authHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.authService.Register(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// GetCurrentUser retrieves the currently authenticated user from the request context.
func (h *authHandler) GetCurrentUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
