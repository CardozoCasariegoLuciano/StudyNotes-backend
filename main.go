package main

import (
	"fmt"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/database"
	migrations "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/Migrations"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/customValidators"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database := database.NewDataBase()
	migrations.MakeAllMigrations(database)
	port := environment.GetApplicationPort()

	e := echo.New()
	e.Validator = customValidators.NewCustomValidator()

	//CURL

	//Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	//Routes
	routes.HanddlerRoutes(e, database)

	//Starting App
	fmt.Printf("Server runnin on port http://localhost%s", port)
	e.Logger.Fatal(e.Start(port))
}
