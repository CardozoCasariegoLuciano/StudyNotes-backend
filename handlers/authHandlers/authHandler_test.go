package authHandlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/authHandlers"
	mock_models "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/mocks"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/customValidators"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const basePath = "/api/v1"

func TestRegister_badCases(t *testing.T) {
	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		path            string
		body            interface{}
		expectedCode    int
		expectedResonse handlers.Response
	}{
		{
			name: "wrong data type",
			path: "/auth/register",
			body: map[string]interface{}{
				"name":            123,
				"email":           "test@example.com",
				"password":        "testpassword",
				"confirmPassword": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: handlers.Response{
				MessageType: "ERROR",
				Message:     "Not valid body information",
				Data:        nil,
			},
		},
		{
			name: "Fileds missing",
			path: "/auth/register",
			body: map[string]interface{}{
				"password":        "testpassword",
				"confirmPassword": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: handlers.Response{
				MessageType: "ERROR",
				Message:     "All fields are required",
				Data:        nil,
			},
		},
		{
			name: "Passwords don´t match",
			path: "/auth/register",
			body: map[string]interface{}{
				"name":            "Test",
				"email":           "test@example.com",
				"password":        "Fake1",
				"confirmPassword": "Fake2",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: handlers.Response{
				MessageType: "ERROR",
				Message:     "Passwords are not equals",
				Data:        nil,
			},
		},
		{
			name: "Email aready taken",
			path: "/auth/register",
			body: map[string]interface{}{
				"name":            "Test",
				"email":           "test@example.com",
				"password":        "testpassword",
				"confirmPassword": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: handlers.Response{
				MessageType: "ERROR",
				Message:     "Email already taken",
				Data:        nil,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	auth := authHandlers.NewAuth(mockUserRepo)

	mockUserRepo.
		EXPECT().
		CreateUser(gomock.AssignableToTypeOf(&models.User{})).
		Return(nil).
		AnyTimes()

	mockUserRepo.
		EXPECT().
		FindUserByEmail(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(&models.User{})).
		Return(&gorm.DB{RowsAffected: 1}).AnyTimes()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//Init objects
			e := echo.New()
			e.Validator = customValidators.NewCustomValidator()

			//Parse map of strings to []bytes of the body
			body, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			//Create new request, recorder(writer) and contetx
			request := httptest.NewRequest(
				http.MethodPost,
				basePath+tc.path,
				strings.NewReader(string(body)),
			)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			context := e.NewContext(request, writer)

			//Call the Register method
			err = auth.Register(context)
			assert.NoError(t, err)

			//Parse the []bytes of the response
			// and fill the result into resp variable
			resp := handlers.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			//Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.Data, resp.Data)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}

func TestRegister_GoodCases(t *testing.T) {
	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		path            string
		body            interface{}
		expectedCode    int
		expectedResonse handlers.Response
	}{
		{
			name: "User successfully created",
			path: "/auth/register",
			body: map[string]interface{}{
				"name":            "Test",
				"email":           "test@example.com",
				"password":        "pass1",
				"confirmPassword": "pass1",
			},
			expectedCode: http.StatusOK,
			expectedResonse: handlers.Response{
				MessageType: "OK",
				Message:     "User created",
				Data:        nil,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	auth := authHandlers.NewAuth(mockUserRepo)

	mockUserRepo.
		EXPECT().
		CreateUser(gomock.AssignableToTypeOf(&models.User{})).
		Return(nil).
		AnyTimes()

	mockUserRepo.
		EXPECT().
		FindUserByEmail(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(&models.User{})).
		Return(&gorm.DB{RowsAffected: 0}).AnyTimes()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init objects
			e := echo.New()
			e.Validator = customValidators.NewCustomValidator()

			// Parse map of strings to []bytes of the body
			body, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			// Create new request, recorder(writer) and contetx
			request := httptest.NewRequest(
				http.MethodPost,
				basePath+tc.path,
				strings.NewReader(string(body)),
			)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			context := e.NewContext(request, writer)

			// Call the Register method
			err = auth.Register(context)
			assert.NoError(t, err)

			// Parse the []bytes of the response
			// and fill the result into resp variable
			resp := handlers.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
			assert.NotNil(t, resp.Data.(map[string]interface{})["token"])
		})
	}
}
