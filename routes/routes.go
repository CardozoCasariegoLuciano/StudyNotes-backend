package routes

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/database"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func HanddlerRoutes(e *echo.Echo) {
	basePath := "/api/v1"
	st := storage.NewStorage(database.DataBase)

	//Auth
	authRoutes := e.Group(basePath + "/auth")
	AuthRoutes(authRoutes, st)
}
