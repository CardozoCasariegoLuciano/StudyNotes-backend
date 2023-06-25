package midlewares

import (
	"fmt"
	"net/http"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			response := responses.NewResponse("ERROR", "Dont have a token", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}
		dataToken := &apimodels.JwtCustomClaims{}
		secret := environment.GetJwtSecretKey()
		tkn, err := jwt.ParseWithClaims(token, dataToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			fmt.Println(err)
			response := responses.NewResponse("ERROR", "Wrong token", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}
		if !tkn.Valid {
			response := responses.NewResponse("ERROR", "Invalid token", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}

		c.Set("userID", dataToken.Id)
		c.Set("userEmail", dataToken.Email)
		c.Set("userRole", dataToken.Role)

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
