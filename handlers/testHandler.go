package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Helloword(c echo.Context) error {
	text := c.Param("text")
	if text == "" {
		text = "hola mundo"
	}

	fmt.Printf("Limpia %s\n", text)

	return c.String(http.StatusOK, text)
}
