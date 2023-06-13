package migrations

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/database"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
)

func MakeAllMigrations() {
	database := database.GetDataBase()

	////// Descomentar para borrar la tabla y volver a comentarlo
	//dbmodels.DropGamesTable(database)
	//dbmodels.DropUsersTable(database)

	dbmodels.UsersMigration(database)
	dbmodels.GameMigration(database)
}
