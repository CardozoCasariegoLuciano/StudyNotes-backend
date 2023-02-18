package authHandlers

import (
	"net/http"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/roles"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	storage models.Istorage
}

func NewAuth(store models.Istorage) *Auth {
	return &Auth{storage: store}
}

func (auth *Auth) Login(c echo.Context) error {
	data := models.Login{}
	err := c.Bind(&data)
	if err != nil {
		response := handlers.NewResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	response := handlers.NewResponse("OK", "User successfully logged", data)
	return c.JSON(http.StatusOK, response)
}

func (auth *Auth) Register(c echo.Context) error {
	data := models.Register{}

	if err := c.Bind(&data); err != nil {
		response := handlers.NewResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(&data); err != nil {
		response := handlers.NewResponse("ERROR", "All fields are required", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if data.Password != data.ConfirmPassword {
		response := handlers.NewResponse("ERROR", "Passwords are not equals", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	//Hashing ths password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		response := handlers.NewResponse("ERROR", "Error hashing the password", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	var newUser models.User
	//Check email is not taken
	result := auth.storage.FindUserByEmail(data.Email, &newUser)
	if result.RowsAffected > 0 {
		response := handlers.NewResponse("ERROR", "Email already taken", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	//Saving the new user
	newUser = models.User{
		Name:     data.Name,
		Password: string(hashedPass),
		Email:    data.Email,
		Role:     roles.USER,
	}
	auth.storage.CreateUser(&newUser)

	//Create token
	claims := models.JwtCustomClaims{
		Id:    int(newUser.ID),
		Email: data.Email,
		Role:  roles.USER,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(environment.GetJwtSecretKey()))
	if err != nil {
		response := handlers.NewResponse("ERROR", "trouble creating a JWT", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	//Good response
	tokenResponse := map[string]string{"token": t}
	response := handlers.NewResponse("OK", "User created", tokenResponse)
	return c.JSON(http.StatusOK, response)
}
