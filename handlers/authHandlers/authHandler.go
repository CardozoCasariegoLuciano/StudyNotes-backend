package authHandlers

import (
	"net/http"
	"strings"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
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

// Login godoc
// @Summary Login a user
// @Description Login a user and get his token
// @Tags Auth
// @Accept json
// @Param Register body models.Login true "request body"
// @Produce json
// @Success 200 {object} responses.Response{data=swaggertypes.SwaggerCustomTypes{token=string}}
// @Failure 400 {object} responses.Response{data=object}
// @Router /auth/login [post]
func (auth *Auth) Login(c echo.Context) error {
	data := models.Login{}

	if err := c.Bind(&data); err != nil {
		response := responses.NewResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(&data); err != nil {
		if strings.Contains(err.Error(), "'email' tag") {
			response := responses.NewResponse("ERROR", "email field must be a valid email", nil)
			return c.JSON(http.StatusBadRequest, response)
		} else {
			response := responses.NewResponse("ERROR", "All fields are required", nil)
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	var userLogged models.User
	//Check email exist
	result := auth.storage.FindUserByEmail(data.Email, &userLogged)
	if result.RowsAffected == 0 {
		response := responses.NewResponse("ERROR", "Wrong email or password", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(userLogged.Password), []byte(data.Password))
	if err != nil {
		response := responses.NewResponse("ERROR", "Wrong email or password", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	//Create token
	claims := models.JwtCustomClaims{
		Id:    int(userLogged.ID),
		Email: userLogged.Email,
		Role:  userLogged.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(environment.GetJwtSecretKey()))
	if err != nil {
		response := responses.NewResponse("ERROR", "trouble creating a JWT", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	//Good response
	tokenResponse := map[string]string{"token": t}
	response := responses.NewResponse("OK", "User created", tokenResponse)

	return c.JSON(http.StatusOK, response)
}

// Register godoc
// @Summary Register new user
// @Description Charge new user into the database
// @Tags Auth
// @Accept json
// @Param Register body models.Register true "request body"
// @Produce json
// @Success 200 {object} responses.Response{data=swaggertypes.SwaggerCustomTypes{token=string}}
// @Failure 400 {object} responses.Response{data=object}
// @Router /auth/register [post]
func (auth *Auth) Register(c echo.Context) error {
	data := models.Register{}

	if err := c.Bind(&data); err != nil {
		response := responses.NewResponse("ERROR", "Not valid body information", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(&data); err != nil {
		if strings.Contains(err.Error(), "'email' tag") {
			response := responses.NewResponse("ERROR", "email field must be a valid email", nil)
			return c.JSON(http.StatusBadRequest, response)
		} else {
			response := responses.NewResponse("ERROR", "All fields are required", nil)
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	if data.Password != data.ConfirmPassword {
		response := responses.NewResponse("ERROR", "Passwords are not equals", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	//Hashing ths password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		response := responses.NewResponse("ERROR", "Error hashing the password", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	var newUser models.User
	//Check email is not taken
	result := auth.storage.FindUserByEmail(data.Email, &newUser)
	if result.RowsAffected > 0 {
		response := responses.NewResponse("ERROR", "Email already taken", nil)
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
		response := responses.NewResponse("ERROR", "trouble creating a JWT", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	//Good response
	tokenResponse := map[string]string{"token": t}
	response := responses.NewResponse("OK", "User created", tokenResponse)

	return c.JSON(http.StatusOK, response)
}
