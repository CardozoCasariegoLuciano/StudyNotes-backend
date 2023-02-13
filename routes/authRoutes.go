package routes

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(authGroup *echo.Group) {

	authGroup.POST("/login", handlers.Login)
	authGroup.POST("/register", handlers.Register)
}
