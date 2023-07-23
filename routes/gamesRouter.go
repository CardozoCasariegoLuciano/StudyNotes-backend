package routes

import (
	gameshandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/gamesHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/midlewares"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func GamesRoutes(gameGroup *echo.Group, storage *storage.Storage) {
	games := gameshandlers.NewGame(storage)

	gameGroup.GET("/", games.GetGames, midlewares.ValidateToken)
	gameGroup.POST("/", games.CreateGame, midlewares.ValidateToken)
	gameGroup.GET("/:gameID", games.GetGameByID, midlewares.ValidateToken, games.GetGameByQueryID)
	gameGroup.PUT("/:gameID", games.EditGame, midlewares.ValidateToken, games.GetGameByQueryID)
	gameGroup.DELETE("/:gameID", games.DeleteGame, midlewares.ValidateToken, games.GetGameByQueryID)
}
