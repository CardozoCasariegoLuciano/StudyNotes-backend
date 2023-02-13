package migrations

import "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"

func MakeAllMigrations() {

	models.UsersMigration()
}
