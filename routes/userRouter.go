package routes

import (
	userhandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/userHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/midlewares"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func UserRoutes(userGroup *echo.Group, storage *storage.Storage) {
	user := userhandlers.NewUser(storage)

	userGroup.GET("/", user.GetUser, midlewares.ValidateToken)
	userGroup.PUT("/", user.EditUser, midlewares.ValidateToken)
}
