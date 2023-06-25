package storage

import (
	"sync"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/database"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

var stor *Storage
var once sync.Once

func GetStorage() *Storage {
	once.Do(func() {
		stor = newStorage()
	})

	return stor
}

func newStorage() *Storage {
	database := database.GetDataBase()
	return &Storage{db: database}
}
