package server

import (
	"log"
	"net/http"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/internal/middleware"
	"github.com/felipe/dev-test-api/internal/modules/auth"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/felipe/dev-test-api/internal/modules/users"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *config.Config, db *gorm.DB) {
	userStore := users.NewStore(db)
	userService := users.NewService(userStore)
	userHandler := users.NewHandler(userService)

	authService := auth.NewService(userStore, cfg.JWT.SecretBytes(), cfg.JWT.ExpiryHrs)
	authHandler := auth.NewHandler(authService)

	r := gin.Default()
	r.Use(middleware.Logger())

	r.GET("/health", healthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	auth.RegisterRoutes(api, authHandler)

	protected := api.Group("")
	protected.Use(middleware.Auth(cfg.JWT.SecretBytes()))
	{
		users.RegisterRoutes(protected, userHandler)
		questions.RegisterRoutes(protected, nil)
	}

	log.Println("Server starting on :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// @Summary      Health check
// @Description  Verifica que el servidor esté corriendo
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
