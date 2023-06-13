package dbmodels

import "gorm.io/gorm"

type Game struct {
	Title       string `json:"title" gorm:"type:VARCHAR(150); NOT NULL"`
	Description string `json:"description" gorm:"type:VARCHAR(555); NOT NULL"`
	gorm.Model
}

type Games []Game

func GameMigration(database *gorm.DB) {
	database.AutoMigrate(Game{})
}

// Se usa en el /helper/migrations
func DropGamesTable(database *gorm.DB) {
	database.Migrator().DropTable(Game{})
}
