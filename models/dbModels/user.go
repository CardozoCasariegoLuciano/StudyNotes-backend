package dbmodels

import (
	"gorm.io/gorm"
)

type User struct {
	Name     string `json:"name" gorm:"type:VARCHAR(255); NOT NULL"`
	Email    string `json:"email" gorm:"type:VARCHAR(255); NOT NULL; UNIQUE"`
	Password string `json:"password" gorm:"type:VARCHAR(255); NOT NULL"`
	Image    string `json:"image" gorm:"type:VARCHAR(255)"`
	Role     string `json:"role" gorm:"type:enum('USER', 'ADMIN', 'SUPER_ADMIN'); NOT NULL; DEFAULT 'USER'"`
	gorm.Model
}

type Users []User

func UsersMigration(database *gorm.DB) {
	database.AutoMigrate(User{})
}

// Se usa en el /helper/migrations
func DropUsersTable(database *gorm.DB) {
	database.Migrator().DropTable(User{})
}
