package moderation

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/reports", h.Create)
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	reports := rg.Group("/reports")
	{
		reports.GET("", h.List)
		reports.PUT("/:id/review", h.Review)
	}
}
