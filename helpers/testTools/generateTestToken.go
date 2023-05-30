package testtools

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(userId int) string {
	userID := userId
	userEmail := "Email"
	userRole := "User"
	claims := apimodels.JwtCustomClaims{
		Id:    userID,
		Email: userEmail,
		Role:  userRole,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ValidToken, _ := token.SignedString([]byte(environment.GetJwtSecretKey()))
	return ValidToken
}
