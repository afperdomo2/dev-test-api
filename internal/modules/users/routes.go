package users

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	users := rg.Group("/users")
	{
		users.GET("", handler.List)
		users.POST("", handler.Create)
		users.GET("/:id", handler.Get)
		users.DELETE("/:id", handler.Delete)
	}
}
