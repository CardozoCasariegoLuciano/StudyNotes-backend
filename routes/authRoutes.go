package routes

import (
	authhandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/authHandlers"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(authGroup *echo.Group) {

	authGroup.POST("/login", authhandlers.Login)
	authGroup.POST("/register", authhandlers.Register)
}
