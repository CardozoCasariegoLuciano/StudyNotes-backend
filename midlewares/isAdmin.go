package midlewares

import (
	"net/http"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	errorcodes "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/errorCodes"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/roles"
	"github.com/labstack/echo/v4"
)

func ValidateIsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get("userRole")
		if userRole == roles.USER {
			response := responses.NewResponse(errorcodes.NO_ROLE, "Don´t have the role", nil)
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
