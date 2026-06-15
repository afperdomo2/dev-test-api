package topics

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	t := rg.Group("/topics")
	{
		t.GET("", h.List)
		t.GET("/:id", h.Get)
	}
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/topics", h.Create)
	rg.PUT("/topics/:id", h.Update)
	rg.DELETE("/topics/:id", h.Delete)
}
