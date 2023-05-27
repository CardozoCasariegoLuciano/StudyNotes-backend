package noteshandlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/labstack/echo/v4"
)

func (note *Note) GetNoteByQueryID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := uint(c.Get("userID").(int))
		Note := &dbmodels.Note{}
		paramID, err := strconv.Atoi(c.Param("noteId"))
		if err != nil {
			response := responses.NewResponse("ERROR", "Invalid ID, Must be a number", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		row := note.storage.GetNoteByID(paramID, Note)
		if row.RowsAffected == 0 {
			response := responses.NewResponse("ERROR", fmt.Sprintf("Note %d not found", paramID), nil)
			return c.JSON(http.StatusNotFound, response)
		}

		if Note.UserID != userID {
			response := responses.NewResponse("ERROR", "This is not your note", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}

		c.Set("Note", Note)

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
