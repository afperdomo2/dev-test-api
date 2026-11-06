package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes monta /auth. El limiter se aplica solo a setup y login:
// /status lo llama el guard del SPA en cada arranque y no debe consumir cuota.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authLimit gin.HandlerFunc) {
	if authLimit == nil {
		authLimit = func(c *gin.Context) { c.Next() }
	}

	auth := rg.Group("/auth")
	{
		auth.GET("/status", handler.Status)
		auth.POST("/setup", authLimit, handler.Setup)
		auth.POST("/login", authLimit, handler.Login)
	}
}
