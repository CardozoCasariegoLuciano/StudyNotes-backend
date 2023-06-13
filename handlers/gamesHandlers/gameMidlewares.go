package gameshandlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/roles"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/labstack/echo/v4"
)

func (game *Game) GetGameByQueryIDFree(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		Game := &dbmodels.Game{}
		paramID, err := strconv.Atoi(c.Param("gameID"))
		if err != nil {
			response := responses.NewResponse("ERROR", "Invalid ID, Must be a number", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		row := game.storage.GetById(paramID, Game)
		if row.RowsAffected == 0 {
			response := responses.NewResponse("ERROR", fmt.Sprintf("Game %d not found", paramID), nil)
			return c.JSON(http.StatusNotFound, response)
		}

		c.Set("Game", Game)

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

func (game *Game) GetGameByQueryIDAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get("userRole")
		Game := &dbmodels.Game{}
		paramID, err := strconv.Atoi(c.Param("gameID"))
		if err != nil {
			response := responses.NewResponse("ERROR", "Invalid ID, Must be a number", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		row := game.storage.GetById(paramID, Game)
		if row.RowsAffected == 0 {
			response := responses.NewResponse("ERROR", fmt.Sprintf("Game %d not found", paramID), nil)
			return c.JSON(http.StatusNotFound, response)
		}

		if userRole == roles.USER {
			response := responses.NewResponse("ERROR", "You can´t see this", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}

		c.Set("Game", Game)

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
