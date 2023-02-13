package environment

import (
	"os"

	"github.com/joho/godotenv"
)

func GetApplicationPort() string {
	var port string
	err := godotenv.Load(".env")
	if err != nil {
		port = ":3000"
	} else {
		port = os.Getenv("PORT")
	}

	return port
}
