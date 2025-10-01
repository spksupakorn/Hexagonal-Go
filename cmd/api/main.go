package main

import (
	"dungeons-dragon-service/internal/config"
	"dungeons-dragon-service/internal/http/server"
	database "dungeons-dragon-service/internal/infrastructure/db"
)

//go:generate swag init -g cmd/api/main.go -o ./docs

// @title           Dungeon Dragon API Documentation
// @version         1.0
// @description     API for managing D&D characters and quests.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://swagger.io/contact/
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {token}" to authenticate.
func main() {
	config.LoadConfig()
	db := database.NewPostgresDatabase()
	server.NewEchoServer(db).Start()
}
