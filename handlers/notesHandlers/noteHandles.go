package noteshandlers

import (
	"net/http"
	"strings"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/labstack/echo/v4"
)

type Note struct {
	storage apimodels.Istorage
}

func NewAuth(store apimodels.Istorage) *Note {
	return &Note{storage: store}
}

func (note *Note) GetUserNotes(c echo.Context) error {
	userID := c.Get("userID")
	AllNotes := &dbmodels.Notes{}

	note.storage.GetAllNotes(userID.(int), AllNotes)

	response := responses.NewResponse("OK", "All user notes", AllNotes)
	return c.JSON(http.StatusOK, response)
}

func (note *Note) CreateNote(c echo.Context) error {
	noteData := apimodels.CreateNoteData{}

	if err := c.Bind(&noteData); err != nil {
		response := responses.NewResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(&noteData); err != nil {
		if strings.Contains(err.Error(), "'min' tag") {
			response := responses.NewResponse("ERROR", "Title field must have more than 3 characters", nil)
			return c.JSON(http.StatusBadRequest, response)
		} else {
			response := responses.NewResponse("ERROR", "All fields are required", nil)
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	userID := c.Get("userID")
	if userID == nil {
		response := responses.NewResponse("ERROR", "unknow author, check the token", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	newNote := &dbmodels.Note{
		Title:       noteData.Title,
		Description: noteData.Description,
		UserID:      uint(userID.(int)),
	}

	note.storage.Create(newNote)

	response := responses.NewResponse("OK", "Note created", newNote)
	return c.JSON(http.StatusCreated, response)
}
