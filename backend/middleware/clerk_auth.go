package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type clerkUser struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

func RequireClerkAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		req, _ := http.NewRequest("GET", "https://api.clerk.dev/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "aplication/json")
		req.Header.Set("Clerk-Secret-Key", os.Getenv("CLERK_SECRET_KEY"))

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session token"})
			c.Abort()
			return
		}

		defer resp.Body.Close()

		var session clerkUser
		if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode response"})
			c.Abort()
			return
		}

		if session.Status != "active" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session is not active"})
			c.Abort()
			return
		}

		c.Set("clerkUserID", session.UserID)
		c.Next()
	}
}
