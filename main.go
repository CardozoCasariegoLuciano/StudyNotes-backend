package main

import (
	"fmt"
	"os"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//PORT .env
	var port string
	err := godotenv.Load(".env")
	if err != nil {
		port = ":3000"
	} else {
		port = os.Getenv("PORT")
	}

	e := echo.New()

	//CURL

	//Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	//Routes
	testRoute := e.Group("/test")
	routes.TestRouter(testRoute)

	//Starting App
	fmt.Printf("Server runnin on port http://localhost%s", port)
	e.Logger.Fatal(e.Start(port))
}
