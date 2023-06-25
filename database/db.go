package database

import (
	"fmt"
	"sync"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dataBaseURI = getDataBaseURI()
var db *gorm.DB
var once sync.Once

func GetDataBase() *gorm.DB {
	once.Do(func() {
		db = newDataBase()
	})
	return db
}

func newDataBase() (db *gorm.DB) {
	env := environment.GetEnvirontment()
	if env != "test" {
		var err error
		db, err = gorm.Open(mysql.Open(dataBaseURI), &gorm.Config{})
		if err != nil {
			fmt.Println("Error en la conexion", err)
			panic(err)
		} else {
			fmt.Println("Conexon de la base de datos exitosa")
			return db
		}
	} else {
		return nil
	}
}
