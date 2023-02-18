package migrations

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"gorm.io/gorm"
)

func MakeAllMigrations(database *gorm.DB) {
	models.UsersMigration(database)
}
