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

func NewNote(store apimodels.Istorage) *Note {
	return &Note{storage: store}
}

//TODO testear todos estos metodos y midlewares
//TODO Documentar todo esto

//TODO Actualizar el Jira, el Docs y el Figma con el nuevo aprouch del proyecto
//TODO Crear un script para Actualizar los mocks y el swager mas facilmente
//TODO Actualizar el readme con el uso de esos scrips e indicando como conectar a la base de datos

//TODO en el front, separar los componentes en otra libreria como zeta
//TODO Crear codigos de ERROR para reemplazar en los endpoints

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

	note.storage.Save(newNote)

	response := responses.NewResponse("OK", "Note created", newNote)
	return c.JSON(http.StatusCreated, response)
}

func (note *Note) GetNoteByID(c echo.Context) error {
	Note := c.Get("Note")
	response := responses.NewResponse("OK", "Note Selected", Note)
	return c.JSON(http.StatusOK, response)
}

func (note *Note) DeleteNoteByID(c echo.Context) error {
	contextNote := c.Get("Note").(*dbmodels.Note)
	note.storage.DeleteNoteByID(int(contextNote.ID), contextNote)

	response := responses.NewResponse("OK", "Note Deleted", contextNote)
	return c.JSON(http.StatusOK, response)
}

func (note *Note) UpdateNoteByID(c echo.Context) error {
	contextNote := c.Get("Note").(*dbmodels.Note)
	reqData := &apimodels.CreateNoteData{}

	if err := c.Bind(reqData); err != nil {
		response := responses.NewResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(reqData); err != nil {
		if strings.Contains(err.Error(), "'min' tag") {
			response := responses.NewResponse("ERROR", "Title field must have more than 3 characters", nil)
			return c.JSON(http.StatusBadRequest, response)
		} else {
			response := responses.NewResponse("ERROR", "All fields are required", nil)
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	contextNote.Title = reqData.Title
	contextNote.Description = reqData.Description

	note.storage.Save(contextNote)

	response := responses.NewResponse("OK", "Note Deleted", contextNote)
	return c.JSON(http.StatusOK, response)
}
