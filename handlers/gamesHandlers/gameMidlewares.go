package gameshandlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	errorcodes "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/errorCodes"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/labstack/echo/v4"
)

func (game *Game) GetGameByQueryID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(int)
		Game := &dbmodels.Game{}
		paramID, err := strconv.Atoi(c.Param("gameID"))
		if err != nil {
			response := responses.NewResponse(errorcodes.INVALID_ID, "Invalid ID, Must be a number", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		row := game.storage.GetGameById(userID, paramID, Game)
		if row.RowsAffected == 0 {
			response := responses.NewResponse(errorcodes.NOT_FOUND_OR_UNAUTHORIZED, fmt.Sprintf("Game %d not found", paramID), nil)
			return c.JSON(http.StatusNotFound, response)
		}

		c.Set("Game", Game)

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
