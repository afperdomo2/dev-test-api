package main

import (
	"log"
	"net/http"

	"github.com/felipe/dev-test-api/config"
	"github.com/felipe/dev-test-api/database"
	_ "github.com/felipe/dev-test-api/docs"
	"github.com/felipe/dev-test-api/internal/auth"
	"github.com/felipe/dev-test-api/internal/questions"
	"github.com/felipe/dev-test-api/internal/users"
	"github.com/felipe/dev-test-api/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           dev-test-api
// @version         1.0
// @description     REST API con Go + Gin · Auth JWT · PostgreSQL · Swagger
// @termsOfService  http://swagger.io/terms/

// @contact.name   Felipe
// @contact.email  felipe@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 JWT token. Prefix with "Bearer ". Example: "Bearer eyJhbG..."

// @schemes  http
func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

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
