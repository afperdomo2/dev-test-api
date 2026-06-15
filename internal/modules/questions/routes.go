package questions

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	t := rg.Group("/questions")
	{
		t.GET("", h.List)
		t.GET("/:id", h.Get)
	}
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/questions", h.Create)
	rg.PUT("/questions/:id", h.Update)
	rg.DELETE("/questions/:id", h.Delete)
}
