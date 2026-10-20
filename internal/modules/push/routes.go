package push

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	p := rg.Group("/push")
	{
		p.POST("/tokens", h.Register)
		p.POST("/tokens/unregister", h.Unregister)
	}
}
