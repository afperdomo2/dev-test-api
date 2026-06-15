package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/profile", h.GetProfile)
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	users := rg.Group("/users")
	{
		users.GET("", h.List)
		users.POST("", h.Create)
		users.GET("/:id", h.Get)
		users.DELETE("/:id", h.Delete)
	}
}
