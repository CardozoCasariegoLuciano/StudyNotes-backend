package authHandlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/authHandlers"
	mock_models "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/mocks"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/customValidators"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const basePath = "/api/v1"

// Register Tests
func TestRegister_badCases(t *testing.T) {
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
		FindUserByEmail(gomock.Eq("mailtaken@example.com"), gomock.AssignableToTypeOf(&models.User{})).
		Return(&gorm.DB{RowsAffected: 1}).AnyTimes()

	mockUserRepo.
		EXPECT().
		FindUserByEmail(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(&models.User{})).
		Return(&gorm.DB{RowsAffected: 0}).AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		path            string
		body            interface{}
		expectedCode    int
		expectedResonse responses.Response
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
			expectedResonse: responses.Response{
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
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "All fields are required",
				Data:        nil,
			},
		},
		{
			name: "No valid email",
			path: "/auth/register",
			body: map[string]interface{}{
				"name":            "test",
				"email":           "novalidEmail",
				"password":        "testpassword",
				"confirmPassword": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "email field must be a valid email",
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
			expectedResonse: responses.Response{
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
				"email":           "mailtaken@example.com",
				"password":        "testpassword",
				"confirmPassword": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Email already taken",
				Data:        nil,
			},
		},
	}

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
			resp := responses.Response{}
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

	testCases := []struct {
		name            string
		path            string
		body            interface{}
		expectedCode    int
		expectedResonse responses.Response
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
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "User created",
			},
		},
	}

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
			resp := responses.Response{}
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

// Login Tests
func TestLogin_badCases(t *testing.T) {
	environment.SetTestEnvirontment()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	auth := authHandlers.NewAuth(mockUserRepo)

	mockUserRepo.
		EXPECT().
		FindUserByEmail(gomock.Eq("noExist@email.com"), gomock.AssignableToTypeOf(&models.User{})).
		Return(&gorm.DB{RowsAffected: 0}).AnyTimes()

	mockUserRepo.
		EXPECT().
		FindUserByEmail(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(&models.User{})).
		Return(&gorm.DB{RowsAffected: 1}).AnyTimes()

	mockUserRepo.
		EXPECT().
		ComparePasswords(gomock.AssignableToTypeOf(""), gomock.Eq("testpassword")).
		Return(errors.New("Error")).AnyTimes()

	testCases := []struct {
		name            string
		path            string
		body            interface{}
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name: "wrong data type",
			path: "/auth/login",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": 12323,
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Not valid body information",
				Data:        nil,
			},
		},
		{
			name: "Fileds missing",
			path: "/auth/login",
			body: map[string]interface{}{
				"password": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "All fields are required",
				Data:        nil,
			},
		},
		{
			name: "No valid email",
			path: "/auth/login",
			body: map[string]interface{}{
				"email":           "novalidEmail",
				"confirmPassword": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "email field must be a valid email",
				Data:        nil,
			},
		},
		{
			name: "Email dont exist",
			path: "/auth/login",
			body: map[string]interface{}{
				"email":    "noExist@email.com",
				"password": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Wrong email or password",
				Data:        nil,
			},
		},
		{
			name: "Passwords dont match",
			path: "/auth/login",
			body: map[string]interface{}{
				"email":    "test@test.com",
				"password": "testpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Wrong email or password",
				Data:        nil,
			},
		},
	}

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
			err = auth.Login(context)
			assert.NoError(t, err)

			//Parse the []bytes of the response
			// and fill the result into resp variable
			resp := responses.Response{}
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

func TestLogin_GoodCases(t *testing.T) {
	environment.SetTestEnvirontment()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	auth := authHandlers.NewAuth(mockUserRepo)

	mockUserRepo.
		EXPECT().
		FindUserByEmail(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(&models.User{})).
		Return(&gorm.DB{RowsAffected: 1}).AnyTimes()

	mockUserRepo.
		EXPECT().
		ComparePasswords(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf("")).
		Return(nil).AnyTimes()

	testCases := []struct {
		name            string
		path            string
		body            interface{}
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name: "User successfully Logged",
			path: "/auth/login",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "pass1",
			},
			expectedCode: http.StatusOK,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "User Logged",
			},
		},
	}

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
			err = auth.Login(context)
			assert.NoError(t, err)

			// Parse the []bytes of the response
			// and fill the result into resp variable
			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
			assert.NotNil(t, resp.Data.(map[string]interface{})["token"])
			assert.NotNil(t, resp.Data.(map[string]interface{})["email"])
			assert.NotNil(t, resp.Data.(map[string]interface{})["userName"])
		})
	}
}
