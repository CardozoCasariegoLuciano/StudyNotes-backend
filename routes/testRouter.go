package routes

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers"
	"github.com/labstack/echo/v4"
)

func TestRouter(t *echo.Group) {

	t.GET("", handlers.Helloword)
}
