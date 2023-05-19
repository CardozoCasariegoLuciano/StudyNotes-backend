package migrations

import (
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"gorm.io/gorm"
)

func MakeAllMigrations(database *gorm.DB) {
	dbmodels.UsersMigration(database)
	dbmodels.NotesMigration(database)
}
