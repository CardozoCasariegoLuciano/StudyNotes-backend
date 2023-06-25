package midlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/customValidators"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/roles"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func generateToken(role string) string {
	// Create token
	userID := 123
	userEmail := "Email"
	userRole := role
	claims := apimodels.JwtCustomClaims{
		Id:    userID,
		Email: userEmail,
		Role:  userRole,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	userToken, _ := token.SignedString([]byte(environment.GetJwtSecretKey()))
	return userToken
}

func TestIsAdminMiddleware_good(t *testing.T) {
	userToken := generateToken(roles.USER)
	adminToken := generateToken(roles.ADMIN)

	environment.SetTestEnvirontment()
	testCases := []struct {
		name                string
		path                string
		token               string
		expectedCode        int
		isBadCase           bool
		expectedResonseData responses.Response
	}{
		{
			name:         "user role",
			path:         "/games/",
			token:        userToken,
			expectedCode: http.StatusUnauthorized,
			isBadCase:    true,
			expectedResonseData: responses.Response{
				MessageType: "ERROR",
				Message:     "Don´t have the role",
				Data:        nil,
			},
		},
		{
			name:         "Admin role",
			isBadCase:    false,
			path:         "/games/",
			token:        adminToken,
			expectedCode: http.StatusOK,
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

			//validate Token
			middleware := ValidateToken(func(c echo.Context) error {
				return nil
			})
			err := middleware(context)
			assert.NoError(t, err)

			//validate Role
			isAdmin := ValidateIsAdmin(func(c echo.Context) error {
				return nil
			})
			err = isAdmin(context)
			assert.NoError(t, err)

			//Test cases
			if tc.isBadCase {
				resp := responses.Response{
					MessageType: "",
					Message:     "",
					Data:        nil,
				}
				err = json.Unmarshal(writer.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedResonseData.Data, resp.Data)
				assert.Equal(t, tc.expectedResonseData.MessageType, resp.MessageType)
				assert.Equal(t, tc.expectedResonseData.Message, resp.Message)
			}

			assert.Equal(t, tc.expectedCode, writer.Code)
		})
	}
}
