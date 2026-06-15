package main

import (
	"log"
	"net/http"

	_ "github.com/felipe/dev-test-api/docs"
	"github.com/felipe/dev-test-api/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           dev-test-api
// @version         1.0
// @description     REST API con Go + Gin.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Felipe
// @contact.email  felipe@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @schemes  http
func main() {
	r := gin.Default()

	r.Use(middleware.Logger())

	r.GET("/health", healthCheck)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
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
