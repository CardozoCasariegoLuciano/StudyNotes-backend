package midlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/customValidators"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	errorcodes "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/errorCodes"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidateTokenMiddleware_bad(t *testing.T) {
	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		path            string
		token           string
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "No Authorization header",
			path:         "/games/",
			token:        "",
			expectedCode: http.StatusUnauthorized,
			expectedResonse: responses.Response{
				MessageType: errorcodes.NO_TOKEN,
				Message:     "Dont have a token",
				Data:        nil,
			},
		},
		{
			name:         "Wrong header",
			path:         "/games/",
			token:        "asd123asd",
			expectedCode: http.StatusUnauthorized,
			expectedResonse: responses.Response{
				MessageType: errorcodes.WRONG_TOKEN,
				Message:     "Wrong token",
				Data:        nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init objects
			e := echo.New()
			e.Validator = customValidators.NewCustomValidator()

			// Create new request, recorder(writer) and contetx
			request := httptest.NewRequest(
				http.MethodGet,
				"/api/v1"+tc.path,
				nil,
			)

			request.Header.Set("Authorization", tc.token)
			request.Header.Set("Content-Type", "application/json")

			writer := httptest.NewRecorder()
			context := e.NewContext(request, writer)

			handler := ValidateToken(func(c echo.Context) error {
				return nil
			})
			err := handler(context)

			// Parse the []bytes of the response
			// and fill the result into resp variable
			resp := responses.Response{
				MessageType: "",
				Message:     "",
				Data:        nil,
			}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.Data, resp.Data)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}

func TestValidateTokenMiddleware_good(t *testing.T) {
	// Create token
	userID := 123
	userEmail := "Email"
	userRole := "User"
	claims := apimodels.JwtCustomClaims{
		Id:    userID,
		Email: userEmail,
		Role:  userRole,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ValidToken, _ := token.SignedString([]byte(environment.GetJwtSecretKey()))

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		path            string
		token           string
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "Ok token",
			path:         "/games/",
			token:        ValidToken,
			expectedCode: http.StatusOK,
			expectedResonse: responses.Response{
				MessageType: errorcodes.WRONG_TOKEN,
				Message:     "Wrong token",
				Data:        nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = customValidators.NewCustomValidator()

			// Create new request, recorder(writer) and contetx
			request := httptest.NewRequest(
				http.MethodGet,
				"/api/v1"+tc.path,
				nil,
			)

			request.Header.Set("Authorization", tc.token)
			request.Header.Set("Content-Type", "application/json")

			writer := httptest.NewRecorder()
			context := e.NewContext(request, writer)

			assert.Nil(t, context.Get("userID"))
			assert.Nil(t, context.Get("userEmail"))
			assert.Nil(t, context.Get("userRole"))

			middleware := ValidateToken(func(c echo.Context) error {
				return nil
			})
			err := middleware(context)
			assert.NoError(t, err)

			// Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			//assert.Equal(t, context.Get("userID"), userID)
			//assert.Equal(t, context.Get("userEmail"), userEmail)
			//assert.Equal(t, context.Get("userRole"), userRole)
		})
	}
}
