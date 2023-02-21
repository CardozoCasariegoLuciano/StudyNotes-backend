package routes

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/authHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(authGroup *echo.Group, storage *storage.Storage) {
	auth := authHandlers.NewAuth(storage)

	authGroup.POST("/login", auth.Login)
	authGroup.POST("/register", auth.Register)
}
