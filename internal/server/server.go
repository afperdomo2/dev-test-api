package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/internal/middleware"
	"github.com/felipe/dev-test-api/internal/modules/auth"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/felipe/dev-test-api/internal/modules/sessions"
	"github.com/felipe/dev-test-api/internal/modules/topics"
	"github.com/felipe/dev-test-api/internal/modules/users"
	"github.com/felipe/dev-test-api/internal/services/ai"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *config.Config, db *gorm.DB) {
	topicStore := topics.NewStore(db)
	topicService := topics.NewService(topicStore)
	topicHandler := topics.NewHandler(topicService)

	questionStore := questions.NewStore(db)
	questionService := questions.NewService(questionStore)
	questionHandler := questions.NewHandler(questionService)

	progressStore := progress.NewStore(db)
	progressService := progress.NewService(progressStore)
	progressHandler := progress.NewHandler(progressService)

	userStore := users.NewStore(db)
	userService := users.NewService(userStore)
	userHandler := users.NewHandler(userService)

	sessionStore := sessions.NewStore(db)
	aiGenerator := ai.NewGenerator(db, cfg.AI)
	sessionService := sessions.NewService(sessionStore, progressService, aiGenerator)
	sessionHandler := sessions.NewHandler(sessionService)

	authService := auth.NewService(userStore, cfg.JWT.SecretBytes(), cfg.JWT.ExpiryHrs)
	authHandler := auth.NewHandler(authService)

	r := gin.Default()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg.Cors.AllowedOrigins))

	r.GET("/health", healthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	auth.RegisterRoutes(api, authHandler)

	protected := api.Group("")
	protected.Use(middleware.Auth(cfg.JWT.SecretBytes()))
	{
		questions.RegisterRoutes(protected, questionHandler)
		topics.RegisterRoutes(protected, topicHandler)
		users.RegisterRoutes(protected, userHandler)
		progress.RegisterRoutes(protected, progressHandler)
		sessions.RegisterRoutes(protected, sessionHandler)

		admin := protected.Group("")
		admin.Use(middleware.AdminOnly())
		{
			users.RegisterAdminRoutes(admin, userHandler)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("🚀 Server starting on :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ failed to run server: %v", err)
		}
	}()

	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
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
