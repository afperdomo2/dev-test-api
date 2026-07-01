package questions

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	t := rg.Group("/questions")
	{
		t.GET("", h.List)
		t.GET("/:id", h.Get)
		t.POST("", h.Create)
		t.PUT("/:id", h.Update)
		t.DELETE("/:id", h.Delete)
	}
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {}
