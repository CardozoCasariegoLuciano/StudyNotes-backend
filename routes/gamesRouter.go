package routes

import (
	gameshandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/gamesHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/midlewares"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/storage"
	"github.com/labstack/echo/v4"
)

func GamesRoutes(gameGroup *echo.Group, storage *storage.Storage) {
	games := gameshandlers.NewGame(storage)

	gameGroup.GET("/", games.GetGames)
	gameGroup.POST("/", games.CreateGame, midlewares.ValidateToken, midlewares.ValidateIsAdmin)
	gameGroup.GET("/:gameID", games.GetGameByID, midlewares.ValidateToken, games.GetGameByQueryIDAdmin)
	gameGroup.PUT("/:gameID", games.EditGame, midlewares.ValidateToken, games.GetGameByQueryIDAdmin)
	gameGroup.DELETE("/:gameID", games.DeleteGame, midlewares.ValidateToken, games.GetGameByQueryIDAdmin)
}
