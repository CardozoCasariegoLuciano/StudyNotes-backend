package gameshandlers

import (
	"net/http"
	"strings"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/labstack/echo/v4"
)

type Game struct {
	storage apimodels.Istorage
}

func NewGame(store apimodels.Istorage) *Game {
	return &Game{storage: store}
}

// GetGames
// @Summary Get all user games
// @Description List all user games
// @Tags Games
// @Accept json
// @Produce json
// @Param Authorization header string false "Token to validate user"
// @Success 200 {object} responses.Response{data=[]dbmodels.Game}
// @Failure 400 {object} responses.Response{data=object}
// @Router /games/ [get]
func (game *Game) GetGames(c echo.Context) error {
	userID := c.Get("userID").(int)
	AllGames := &dbmodels.Games{}
	game.storage.GetAllGames(userID, AllGames)

	response := responses.NewResponse("OK", "All games", AllGames)
	return c.JSON(http.StatusOK, response)
}

// GetGameByID
// @Summary Get game by ID
// @Description Get one game by param id
// @Tags Games
// @Accept json
// @Produce json
// @Param gameID path string true "ID to find the game"
// @Param Authorization header string false "Token to validate user"
// @Success 200 {object} responses.Response{data=dbmodels.Game}
// @Failure 400 {object} responses.Response{data=object}
// @Router /games/{gameID} [get]
func (game *Game) GetGameByID(c echo.Context) error {
	Game := c.Get("Game")
	response := responses.NewResponse("OK", "Game Selected", Game)
	return c.JSON(http.StatusOK, response)
}

// CreateGame
// @Summary Create new game
// @Description Create new game
// @Tags Games
// @Accept json
// @Produce json
// @Param GameData body apimodels.CreateGameData true "request body"
// @Param Authorization header string false "Token to validate user"
// @Success 200 {object} responses.Response{data=dbmodels.Game}
// @Failure 400 {object} responses.Response{data=object}
// @Router /games/ [post]
func (game *Game) CreateGame(c echo.Context) error {
	userID := c.Get("userID").(int)
	gameData := apimodels.CreateGameData{}

	if err := c.Bind(&gameData); err != nil {
		response := responses.NewResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(&gameData); err != nil {
		if strings.Contains(err.Error(), "'min' tag") {
			response := responses.NewResponse("ERROR", "Title field must have more than 3 characters", nil)
			return c.JSON(http.StatusBadRequest, response)
		} else {
			response := responses.NewResponse("ERROR", "All fields are required", nil)
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	newGame := &dbmodels.Game{
		Title:       gameData.Title,
		Description: gameData.Description,
		UserID:      userID,
	}

	game.storage.Save(newGame)

	response := responses.NewResponse("OK", "Game created", newGame)
	return c.JSON(http.StatusCreated, response)
}

// EditGame
// @Summary Edit game
// @Description Edit game
// @Tags Games
// @Accept json
// @Produce json
// @Param gameID path string true "ID to find the game"
// @Param GameData body apimodels.CreateGameData true "request body"
// @Param Authorization header string false "Token to validate user"
// @Success 200 {object} responses.Response{data=dbmodels.Game}
// @Failure 400 {object} responses.Response{data=object}
// @Router /games/{gameID} [PUT]
func (game *Game) EditGame(c echo.Context) error {
	contextGame := c.Get("Game").(*dbmodels.Game)

	reqData := &apimodels.CreateGameData{}

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

	contextGame.Title = reqData.Title
	contextGame.Description = reqData.Description

	game.storage.Save(contextGame)

	response := responses.NewResponse("OK", "Game edited", contextGame)
	return c.JSON(http.StatusOK, response)
}

// DeleteGame
// @Summary Delete game
// @Description Delete game
// @Tags Games
// @Accept json
// @Produce json
// @Param gameID path string true "ID to find the game"
// @Param Authorization header string false "Token to validate user"
// @Success 200 {object} responses.Response{data=dbmodels.Game}
// @Failure 400 {object} responses.Response{data=object}
// @Router /games/{gameID} [DELETE]
func (game *Game) DeleteGame(c echo.Context) error {
	contextGame := c.Get("Game").(*dbmodels.Game)

	game.storage.DeleteByID(int(contextGame.ID), contextGame)

	response := responses.NewResponse("OK", "Game Deleted", contextGame)
	return c.JSON(http.StatusOK, response)
}
