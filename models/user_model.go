package models

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/database"
	"gorm.io/gorm"
)

// TODO ver que mas campos le puedo agregar aca de GORM
type User struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Role        string `json:"role"`
	gorm.Model
}

type Users []User

func UsersMigration() {
	database.DataBase.AutoMigrate(User{})
}
