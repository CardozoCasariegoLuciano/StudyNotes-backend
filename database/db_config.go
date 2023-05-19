package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func getDataBaseURI() string {
	var db_password string
	var db_host string
	var db_name string

	err := godotenv.Load(".env")
	if err != nil {
		db_password = "123luciano456"
		db_host = "localhost:3306"
		db_name = "studyNotes"
	} else {
		db_password = os.Getenv("DB_PASSWORD")
		db_host = os.Getenv("DB_HOST")
		db_name = os.Getenv("DB_NAME")
	}

	URI := fmt.Sprintf("root:%s@tcp(%s)/%s?parseTime=true", db_password, db_host, db_name)
	return URI
}
