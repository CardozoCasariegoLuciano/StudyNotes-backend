package handlers

import (
	"net/http"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/database"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/roles"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	data := models.Login{}
	err := c.Bind(&data)
	if err != nil {
		response := newResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := newResponse("OK", "User successfully logged", data)
	return c.JSON(http.StatusOK, response)
}

func Register(c echo.Context) error {
	data := models.Register{}

	if err := c.Bind(&data); err != nil {
		response := newResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(&data); err != nil {
		response := newResponse("ERROR", "All fields are required", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if data.Password != data.ConfirmPassword {
		response := newResponse("ERROR", "Passwords are not equals", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	var newUser models.User

	//Check email is not taken
	result := database.DataBase.Where("email = ?", data.Email).First(&newUser)
	if result.RowsAffected > 0 {
		response := newResponse("ERROR", "Email already taken", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	//Hashing ths password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		response := newResponse("ERROR", "Error hashing the password", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	//Saving the new user
	newUser = models.User{
		Name:     data.Name,
		Password: string(hashedPass),
		Email:    data.Email,
		Role:     roles.USER,
	}
	database.DataBase.Save(&newUser)

	//Create token
	claims := models.JwtCustomClaims{
		Id:    int(newUser.Id),
		Email: data.Email,
		Role:  roles.USER,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(environment.GetJwtSecretKey()))
	if err != nil {
		response := newResponse("ERROR", "trouble creating a JWT", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	//Good response
	tokenResponse := map[string]string{"token": t}
	response := newResponse("OK", "User created", tokenResponse)
	return c.JSON(http.StatusOK, response)
}
