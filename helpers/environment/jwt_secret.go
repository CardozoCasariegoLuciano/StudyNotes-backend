package environment

import (
	"os"

	"github.com/joho/godotenv"
)

func GetJwtSecretKey() string {
	var secret string
	err := godotenv.Load(".env")
	if err != nil {
		secret = "Secret"
	} else {
		secret = os.Getenv("JWT_SECRET")
	}

	return secret
}
