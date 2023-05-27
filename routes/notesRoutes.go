package routes

import (
	noteshandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/notesHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func NotesRoutes(noteGroup *echo.Group, storage *storage.Storage) {
	notes := noteshandlers.NewNote(storage)

	noteGroup.GET("/", notes.GetUserNotes)
	noteGroup.POST("/", notes.CreateNote)
	noteGroup.GET("/:noteId", notes.GetNoteByID, notes.GetNoteByQueryID)
	noteGroup.DELETE("/:noteId", notes.DeleteNoteByID, notes.GetNoteByQueryID)
	noteGroup.PUT("/:noteId", notes.UpdateNoteByID, notes.GetNoteByQueryID)
}
