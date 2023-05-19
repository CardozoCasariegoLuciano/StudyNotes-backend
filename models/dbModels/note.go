package dbmodels

import (
	"gorm.io/gorm"
)

type Note struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      uint   `json:"userID" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	gorm.Model
}

type Notes []Note

func NotesMigration(database *gorm.DB) {
	//database.Migrator().DropTable(Note{})
	database.AutoMigrate(Note{})
}
