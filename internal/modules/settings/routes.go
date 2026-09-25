package settings

import "github.com/gin-gonic/gin"

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	s := rg.Group("/settings")
	{
		s.GET("", h.List)
		s.PUT("", h.Update)
	}
}
