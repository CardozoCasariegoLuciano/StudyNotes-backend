package testtools

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(uID int, uRole string, uEmail string) string {
	claims := apimodels.JwtCustomClaims{
		Id:    uID,
		Email: uEmail,
		Role:  uRole,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ValidToken, _ := token.SignedString([]byte(environment.GetJwtSecretKey()))
	return ValidToken
}
