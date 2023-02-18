package routes

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func HanddlerRoutes(e *echo.Echo, database *gorm.DB) {
	basePath := "/api/v1"
	st := storage.NewStorage(database)

	//Auth
	authRoutes := e.Group(basePath + "/auth")
	AuthRoutes(authRoutes, st)
}
