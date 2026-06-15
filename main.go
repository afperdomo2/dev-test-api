package main

import (
	"github.com/felipe/dev-test-api/config"
	"github.com/felipe/dev-test-api/database"
	_ "github.com/felipe/dev-test-api/docs"
	"github.com/felipe/dev-test-api/internal/server"
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
	server.Run(cfg, db)
}
