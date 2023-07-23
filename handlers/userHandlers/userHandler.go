package userhandlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	errorcodes "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/errorCodes"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/labstack/echo/v4"
)

type User struct {
	storage apimodels.Istorage
}

func NewUser(store apimodels.Istorage) *User {
	return &User{storage: store}
}

// Get user
// @Summary Get user
// @Description Get user using the token in the header
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string false "Token to validate user"
// @Success 200 {object} responses.Response{data=[]apimodels.User}
// @Failure 400 {object} responses.Response{data=object}
// @Router /user/ [get]
func (user *User) GetUser(c echo.Context) error {
	userID := c.Get("userID").(int)
	userDB := &dbmodels.User{}

	row := user.storage.GetUserByID(userID, userDB)
	if row.RowsAffected == 0 {
		response := responses.NewResponse(
			errorcodes.NOT_FOUND,
			fmt.Sprintf("User %d not found", userID),
			nil,
		)
		return c.JSON(http.StatusNotFound, response)
	}

	userData := apimodels.User{
		Id:    userDB.ID,
		Email: userDB.Email,
		Name:  userDB.Name,
		Image: userDB.Image,
		Role:  userDB.Role,
	}

	response := responses.NewResponse("OK", "User data", userData)
	return c.JSON(http.StatusOK, response)
}

// Edit user
// @Summary Edit user
// @Description Edit the user using the token in header
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string false "Token to validate user"
// @Param UserData body apimodels.EditUserData true "request body"
// @Success 200 {object} responses.Response{data=[]apimodels.User}
// @Failure 400 {object} responses.Response{data=object}
// @Router /user/ [put]
func (user *User) EditUser(c echo.Context) error {
	userID := c.Get("userID").(int)
	userDB := &dbmodels.User{}
	reqData := &apimodels.EditUserData{}

	row := user.storage.GetUserByID(userID, userDB)
	if row.RowsAffected == 0 {
		response := responses.NewResponse(
			errorcodes.NOT_FOUND,
			fmt.Sprintf("User %d not found", userID),
			nil,
		)
		return c.JSON(http.StatusNotFound, response)
	}

	if err := c.Bind(reqData); err != nil {
		response := responses.NewResponse(
			errorcodes.BODY_TYPES_ERROR,
			"Not valid body information",
			nil,
		)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(reqData); err != nil {
		if strings.Contains(err.Error(), "'max' tag") {
			response := responses.NewResponse(
				errorcodes.BODY_VALIDATION_ERROR,
				"Name field must have less than 30 characters",
				nil,
			)
			return c.JSON(http.StatusBadRequest, response)
		} else {
			response := responses.NewResponse(
				errorcodes.BODY_VALIDATION_ERROR,
				"There are required fields",
				nil,
			)
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	userDB.Name = reqData.Name
	userDB.Image = reqData.Image

	user.storage.Save(userDB)

	userData := apimodels.User{
		Id:    userDB.ID,
		Email: userDB.Email,
		Name:  userDB.Name,
		Image: userDB.Image,
		Role:  userDB.Role,
	}

	response := responses.NewResponse("OK", "User edited", userData)
	return c.JSON(http.StatusOK, response)
}
