package tips

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	t := rg.Group("/tips")
	{
		t.GET("", h.List)
		// Debe ir antes de "/:id" para que "today" no se parsee como UUID.
		t.GET("/today", h.GetOfTheDay)
		t.GET("/:id", h.Get)
	}
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	tips := rg.Group("/tips")
	{
		tips.POST("", h.Create)
		tips.PUT("/:id", h.Update)
		tips.DELETE("/:id", h.Delete)
		tips.POST("/import", h.Import)
	}
}
