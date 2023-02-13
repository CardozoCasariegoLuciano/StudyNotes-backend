package routes

import "github.com/labstack/echo/v4"

func HanddlerRoutes(e *echo.Echo) {

	basePath := "/api/v1"

	//Auth
	authRoutes := e.Group(basePath + "/auth")
	AuthRoutes(authRoutes)
}
