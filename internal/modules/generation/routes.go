package generation

import "github.com/gin-gonic/gin"

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/generation")
	{
		g.POST("", h.CreateJob)
		g.GET("", h.ListJobs)
		g.GET("/stats", h.Stats)
		g.GET("/:id", h.GetJob)
	}
}
