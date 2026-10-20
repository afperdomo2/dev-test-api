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
	"github.com/felipe/dev-test-api/internal/modules/gamification"
	"github.com/felipe/dev-test-api/internal/modules/generation"
	"github.com/felipe/dev-test-api/internal/modules/insights"
	"github.com/felipe/dev-test-api/internal/modules/moderation"
	"github.com/felipe/dev-test-api/internal/modules/notifications"
	opencodeModule "github.com/felipe/dev-test-api/internal/modules/opencode"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/felipe/dev-test-api/internal/modules/sessions"
	"github.com/felipe/dev-test-api/internal/modules/settings"
	"github.com/felipe/dev-test-api/internal/modules/tips"
	"github.com/felipe/dev-test-api/internal/modules/topics"
	"github.com/felipe/dev-test-api/internal/modules/users"
	"github.com/felipe/dev-test-api/internal/services/ai"
	opencodeService "github.com/felipe/dev-test-api/internal/services/opencode"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *config.Config, db *gorm.DB) {
	topicStore := topics.NewStore(db)
	topicService := topics.NewService(topicStore)
	topicHandler := topics.NewHandler(topicService)

	userStore := users.NewStore(db)
	userService := users.NewService(userStore)
	userHandler := users.NewHandler(userService)

	questionStore := questions.NewStore(db)
	questionService := questions.NewService(questionStore, topicStore, userStore)
	questionHandler := questions.NewHandler(questionService)

	tipStore := tips.NewStore(db)
	tipService := tips.NewService(tipStore, topicStore)
	tipHandler := tips.NewHandler(tipService)

	moderationStore := moderation.NewStore(db)
	moderationService := moderation.NewService(moderationStore, questionStore, tipStore)
	moderationHandler := moderation.NewHandler(moderationService)

	gamificationStore := gamification.NewStore(db)
	gamificationService := gamification.NewService(gamificationStore)
	gamificationHandler := gamification.NewHandler(gamificationService)

	progressStore := progress.NewStore(db)
	progressService := progress.NewService(progressStore, gamificationService)
	progressHandler := progress.NewHandler(progressService)

	insightsStore := insights.NewStore(db)
	insightsService := insights.NewService(insightsStore, progressService, gamificationService)
	insightsHandler := insights.NewHandler(insightsService)

	sessionStore := sessions.NewStore(db)

	// settings es la fuente de verdad de los ajustes tunables de IA. Se crea
	// antes que el generador porque implementa ai.SettingsProvider.
	settingsStore := settings.NewStore(db)
	settingsService := settings.NewService(settingsStore, cfg.AI)
	settingsHandler := settings.NewHandler(settingsService)

	aiGenerator := ai.NewGenerator(db, cfg.AI, settingsService)

	generationStore := generation.NewStore(db)
	generationService := generation.NewService(generationStore, aiGenerator, settingsService)
	generationHandler := generation.NewHandler(generationService)

	concurrency, err := settingsService.WorkerConcurrency()
	if err != nil {
		log.Printf("❌ No se pudo leer la concurrencia inicial del pool: %v", err)
		concurrency = 1
	}
	generationWorker := generation.NewWorker(generationService, concurrency, func() int {
		n, cErr := settingsService.WorkerConcurrency()
		if cErr != nil {
			return concurrency
		}
		return n
	})

	sessionService := sessions.NewService(sessionStore, progressService, gamificationService, tipService)
	sessionHandler := sessions.NewHandler(sessionService)

	authService := auth.NewService(userStore, cfg.JWT.SecretBytes(), cfg.JWT.ExpiryHrs)
	authHandler := auth.NewHandler(authService)

	// pusher nil en web-only; push lo inyecta en su ola (D-28).
	notificationStore := notifications.NewStore(db)
	notificationService := notifications.NewService(notificationStore, nil)
	notificationHandler := notifications.NewHandler(notificationService)
	notificationScheduler := notifications.NewScheduler(notificationService)

	opencodeClient := opencodeService.NewClient(cfg.AI, cfg.OpenCode)
	opencodeSvc := opencodeModule.NewService(opencodeClient)
	opencodeHandler := opencodeModule.NewHandler(opencodeSvc)

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
		tips.RegisterRoutes(protected, tipHandler)
		topics.RegisterRoutes(protected, topicHandler)
		users.RegisterRoutes(protected, userHandler)
		progress.RegisterRoutes(protected, progressHandler)
		sessions.RegisterRoutes(protected, sessionHandler)
		moderation.RegisterRoutes(protected, moderationHandler)
		gamification.RegisterRoutes(protected, gamificationHandler)
		insights.RegisterRoutes(protected, insightsHandler)
		notifications.RegisterRoutes(protected, notificationHandler)

		admin := protected.Group("")
		admin.Use(middleware.AdminOnly())
		{
			users.RegisterAdminRoutes(admin, userHandler)
			tips.RegisterAdminRoutes(admin, tipHandler)
			moderation.RegisterAdminRoutes(admin, moderationHandler)
			opencodeModule.RegisterAdminRoutes(admin, opencodeHandler)
			settings.RegisterAdminRoutes(admin, settingsHandler)
			generation.RegisterAdminRoutes(admin, generationHandler)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	poolCtx, stopPool := context.WithCancel(context.Background())
	poolDone := make(chan struct{})
	go func() {
		defer close(poolDone)
		generationWorker.Start(poolCtx)
	}()

	schedCtx, stopScheduler := context.WithCancel(context.Background())
	schedDone := make(chan struct{})
	go func() {
		defer close(schedDone)
		notificationScheduler.Start(schedCtx)
	}()

	go func() {
		log.Println("🚀 Server starting on :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ failed to run server: %v", err)
		}
	}()

	<-quit
	log.Println("🛑 Shutting down server...")

	// Primero se corta el pool: así ningún trabajo nuevo se reclama mientras el
	// servidor cierra sus listeners.
	stopPool()
	select {
	case <-poolDone:
	case <-time.After(15 * time.Second):
		log.Println("⚠️ El pool de generación no terminó a tiempo")
	}

	// El scheduler se detiene antes de cerrar el HTTP server.
	stopScheduler()
	select {
	case <-schedDone:
	case <-time.After(5 * time.Second):
		log.Println("⚠️ El scheduler de notificaciones no terminó a tiempo")
	}

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
