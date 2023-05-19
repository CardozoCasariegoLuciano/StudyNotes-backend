package routes

import (
	noteshandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/notesHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func NotesRoutes(noteGroup *echo.Group, storage *storage.Storage) {
	notes := noteshandlers.NewAuth(storage)

	noteGroup.GET("/", notes.GetUserNotes)
	noteGroup.POST("/", notes.CreateNote)
}
