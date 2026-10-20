package notifications

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	n := rg.Group("/notifications")
	{
		n.GET("", h.List)
		n.GET("/unread-count", h.UnreadCount)
		// Rutas estáticas antes de la paramétrica (mismo patrón que progress).
		n.POST("/read-all", h.MarkAllRead)
		n.POST("/:id/read", h.MarkRead)
	}
}
