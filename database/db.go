package database

import (
	"fmt"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dataBaseURI = getDataBaseURI()

// TODO ver lo del beforeeahc para los tests
func NewDataBase() (db *gorm.DB) {
	env := environment.GetEnvirontment()
	if env != "test" {
		db, err := gorm.Open(mysql.Open(dataBaseURI), &gorm.Config{})
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
