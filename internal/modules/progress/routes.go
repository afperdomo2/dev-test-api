package progress

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	p := rg.Group("/progress")
	{
		p.POST("/:question_id/answer", h.Answer)
		p.GET("/upcoming", h.Upcoming)
		p.GET("/saved", h.Saved)
		p.POST("/:question_id/toggle-save", h.ToggleSave)
	}
}
