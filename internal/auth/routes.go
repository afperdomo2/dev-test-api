package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/setup", handler.Setup)
		auth.POST("/login", handler.Login)
	}
}
