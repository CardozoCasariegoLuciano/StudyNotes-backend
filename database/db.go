package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dataBaseURI = getDataBaseURI()

var DataBase = func() (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open(dataBaseURI), &gorm.Config{})
	if err != nil {
		fmt.Println("Error en la conexion", err)
		panic(err)
	} else {
		fmt.Println("Conexon de la base de datos exitosa")
		return db
	}
}()
