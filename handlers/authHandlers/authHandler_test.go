package authHandlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/authHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/customValidators"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const basePath = "/api/v1"

// TODO ver si puedo sacar de aca este FakeDatabase
type fakeDatabase struct {
}

func (f *fakeDatabase) CreateUser(user *models.User) *gorm.DB {
	return nil
}

func (f *fakeDatabase) FindUserByEmail(email string, model *models.User) *gorm.DB {
	return &gorm.DB{}
}

func TestRegister_badCases(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = customValidators.NewCustomValidator()
			auth := authHandlers.NewAuth(&fakeDatabase{})

			body, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			request := httptest.NewRequest(
				http.MethodPost,
				basePath+tc.path,
				strings.NewReader(string(body)),
			)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			httptest.NewRecorder()
			context := e.NewContext(request, writer)

			err = auth.Register(context)
			assert.NoError(t, err)

			resp := handlers.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.Data, resp.Data)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}

func TestRegister_GoodCases(t *testing.T) {
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

	//TODO comentar este test con las distinatas partes

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = customValidators.NewCustomValidator()
			auth := authHandlers.NewAuth(&fakeDatabase{})

			body, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			request := httptest.NewRequest(
				http.MethodPost,
				basePath+tc.path,
				strings.NewReader(string(body)),
			)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			httptest.NewRecorder()
			context := e.NewContext(request, writer)

			err = auth.Register(context)
			assert.NoError(t, err)

			resp := handlers.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
			assert.NotNil(t, response["data"].(map[string]interface{})["token"])
		})
	}
}
