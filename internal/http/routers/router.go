package router

import (
	"dungeons-dragon-service/internal/http/handlers"
	middleware "dungeons-dragon-service/internal/http/middlewares"
	usecase "dungeons-dragon-service/internal/usecases"

	"github.com/labstack/echo/v4"
)

func NewEchoRouter(e *echo.Echo, jwtMW *middleware.JWTMiddleware, auth usecase.AuthUseCase, opt usecase.OptionUseCase, ch usecase.CharacterUseCase, q usecase.QuestUseCase, img usecase.ImageUseCase) {
	apiV1 := e.Group("/api/v1")
	apiV1.GET("/health", func(c echo.Context) error {
		//message with emoji
		return c.String(200, "still alive! 🧙‍♂️")
	})

	// Global JWT parser (non-blocking)
	apiV1.Use(jwtMW.Parse)

	// Auth
	authH := handlers.NewAuthHandler(auth)
	apiV1.POST("/auth/login", authH.Login)
	apiV1.POST("/auth/register", authH.Register)

	// Public and registered access
	charH := handlers.NewCharacterHandler(ch)
	questH := handlers.NewQuestHandler(q)
	optH := handlers.NewOptionHandler(opt)
	imgH := handlers.NewImageHandler(img)

	apiV1.GET("/characters", charH.List) // Public => public only, Registered => all
	apiV1.GET("/quests", questH.List)

	// Options are public for listing
	apiV1.GET("/options/classes", optH.ListClasses)
	apiV1.GET("/options/races", optH.ListRaces)
	apiV1.GET("/options/quest-levels", optH.ListQuestLevels)

	apiV1.GET("/pictures/:filename", imgH.GetImage)

	// Registered users can create/edit/delete their own
	gAuth := apiV1.Group("", middleware.RequireAuth)
	gAuth.POST("/characters", charH.Create)
	gAuth.PUT("/characters/:id", charH.Update)
	gAuth.DELETE("/characters/:id", charH.Delete)

	gAuth.POST("/quests", questH.Create)
	gAuth.PUT("/quests/:id", questH.Update)
	gAuth.DELETE("/quests/:id", questH.Delete)

	gAuth.POST("/characters/:id/images", imgH.UploadCharacterImage)
	gAuth.POST("/quests/:id/images", imgH.UploadQuestImage)

	// Admin option management
	gAdmin := apiV1.Group("/admin", middleware.RequireAuth, middleware.RequireAdmin)
	gAdmin.POST("/options/classes", optH.CreateClass)
	gAdmin.PUT("/options/classes/:id", optH.UpdateClass)
	gAdmin.DELETE("/options/classes/:id", optH.DeleteClass)

	gAdmin.POST("/options/races", optH.CreateRace)
	gAdmin.PUT("/options/races/:id", optH.UpdateRace)
	gAdmin.DELETE("/options/races/:id", optH.DeleteRace)

	gAdmin.POST("/options/quest-levels", optH.CreateQuestLevel)
	gAdmin.PUT("/options/quest-levels/:id", optH.UpdateQuestLevel)
	gAdmin.DELETE("/options/quest-levels/:id", optH.DeleteQuestLevel)
}
