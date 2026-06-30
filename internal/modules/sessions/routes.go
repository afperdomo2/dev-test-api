package sessions

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	s := rg.Group("/sessions")
	{
		s.GET("", h.List)
		s.POST("", h.Create)
		s.GET("/:id", h.Get)
		s.GET("/:id/summary", h.Summary)
		s.PUT("/:id/finish", h.Finish)
		s.GET("/:id/next", h.NextQuestion)
		s.POST("/:id/answer", h.Answer)
		s.DELETE("/:id", h.Delete)
	}
}
