package storage

import (
	"sync"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

var stor *Storage
var once sync.Once

func GetStorage(db *gorm.DB) *Storage {
	once.Do(func() {
		stor = newStorage(db)
	})

	return stor
}

func newStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}
