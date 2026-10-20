package insights

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	i := rg.Group("/insights")
	{
		i.GET("/summary", h.Summary)
	}
}
