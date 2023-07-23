package routes

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func HanddlerRoutes(e *echo.Echo) {
	basePath := "/api/v1"
	st := storage.GetStorage()

	//Auth
	authRoutes := e.Group(basePath + "/auth")
	AuthRoutes(authRoutes, st)

	//Games
	gamesRoutes := e.Group(basePath + "/games")
	GamesRoutes(gamesRoutes, st)

	//User
	userRoutes := e.Group(basePath + "/user")
	UserRoutes(userRoutes, st)
}
