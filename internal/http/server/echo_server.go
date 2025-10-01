package server

import (
	"context"
	"dungeons-dragon-service/docs"
	"dungeons-dragon-service/internal/config"
	"dungeons-dragon-service/internal/http/handlers"
	"dungeons-dragon-service/internal/http/middlewares"
	router "dungeons-dragon-service/internal/http/routers"
	"dungeons-dragon-service/internal/repositories"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	database "dungeons-dragon-service/internal/infrastructure/db"

	usecase "dungeons-dragon-service/internal/usecases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type echoServer struct {
	app *echo.Echo
	db  database.Database
}

var (
	once sync.Once
	app  *echoServer
)

func NewEchoServer(db database.Database) Server {
	echoApp := echo.New()
	echoApp.HideBanner = true
	echoApp.Logger.SetLevel(log.DEBUG)

	once.Do(func() {
		app = &echoServer{
			app: echoApp,
			db:  db,
		}
	})
	return app
}

func (s *echoServer) Start() {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())
	s.app.Use(middleware.CORS())

	s.initializeRouter()
	s.httpListening()
}

func (s *echoServer) httpListening() {
	// Start server in a goroutine
	serverUrl := fmt.Sprintf(":%d", config.GetConfigInt("PORT"))

	go func() {
		if err := s.app.Start(serverUrl); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func (s *echoServer) initializeRouter() {
	gormDB := s.db.ConnectDB()

	// Repositories
	userRepo := repositories.NewUserRepo(gormDB)
	classRepo := repositories.NewClassRepo(gormDB)
	raceRepo := repositories.NewRaceRepo(gormDB)
	questLevelRepo := repositories.NewQuestLevelRepo(gormDB)
	charRepo := repositories.NewCharacterRepo(gormDB)
	questRepo := repositories.NewQuestRepo(gormDB)
	imageRepo := repositories.NewImageRepo(gormDB)

	// Use cases
	authUC := usecase.NewAuthUsecase(userRepo, config.GetConfigString("JWT_SECRET"))
	optUC := usecase.NewOptionUseCase(classRepo, raceRepo, questLevelRepo, charRepo, questRepo)
	charUC := usecase.NewCharacterUsecase(charRepo, classRepo, raceRepo)
	questUC := usecase.NewQuestUsecase(questRepo, questLevelRepo)
	imageUC := usecase.NewImageUsecase(imageRepo, charRepo, questRepo)

	// Middlewares
	jwtMW := middlewares.NewJWTMiddleware(config.GetConfigString("JWT_SECRET"))
	// Swagger setup
	docs.SwaggerInfo.Title = "Dungeon Dragon API Documentation"
	docs.SwaggerInfo.Description = "API for managing D&D characters and quests."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.GetConfigString("PORT")
	if config.GetConfigString("PORT") == "" || config.GetConfigString("PORT") == "8080" {
		docs.SwaggerInfo.Host = "localhost:8080"
	}
	docs.SwaggerInfo.BasePath = "/api/v1"

	//*if setup Swagger UI
	// s.app.GET("/swagger/*", echoSwagger.WrapHandler)

	//*if setup rapiDoc UI
	// Serve raw OpenAPI JSON (consumed by RapiDoc)
	s.app.GET("/openapi.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json; charset=utf-8", []byte(docs.SwaggerInfo.ReadDoc()))
	})

	// Serve RapiDoc UI
	s.app.GET("/rapidoc", handlers.RapiDoc)

	// Routes
	router.NewEchoRouter(s.app, jwtMW, authUC, optUC, charUC, questUC, imageUC)
}
