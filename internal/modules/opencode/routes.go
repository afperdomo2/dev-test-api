package opencode

import "github.com/gin-gonic/gin"

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/opencode")
	{
		g.GET("/usage", h.GetUsage)
	}
}
