package utils

import "github.com/gin-gonic/gin"

func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := val.(uint)
	return userID, ok
}
