package routes

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/handlers"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, h handlers.AuthHandler) {
	rg.POST("/auth/login", h.Login)
	rg.POST("/auth/register", h.Register)
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/auth/me", h.GetCurrentUser)
}
