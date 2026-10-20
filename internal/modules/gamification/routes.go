package gamification

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/gamification")
	{
		g.GET("/overview", h.Overview)
		g.PUT("/daily-goal", h.SetDailyGoal)
		g.POST("/streak/freeze", h.Freeze)
		g.GET("/badges", h.Badges)
	}
}
