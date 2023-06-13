package gameshandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/roles"
	testtools "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/testTools"
	mock_models "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/testTools/mocks"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const basePath = "/api/v1"

func Test_GetByID_AdminMiddleware_badCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)

	game := NewGame(mockUserRepo)

	notFoundID := 2
	foundID := 1

	mockUserRepo.
		EXPECT().
		GetById(gomock.Eq(notFoundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 0}).
		AnyTimes()

	mockUserRepo.
		EXPECT().
		GetById(gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		gameID          string
		expectedResonse responses.Response
		expectedCode    int
		role            string
	}{
		{
			name:         "No valid id",
			gameID:       "asdasd",
			role:         roles.USER,
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Invalid ID, Must be a number",
				Data:        nil,
			},
		},
		{
			name:         "Not found id",
			role:         roles.USER,
			gameID:       strconv.Itoa(notFoundID),
			expectedCode: http.StatusNotFound,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     fmt.Sprintf("Game %d not found", notFoundID),
				Data:        nil,
			},
		},
		{
			name:         "User is not admin",
			role:         roles.USER,
			gameID:       strconv.Itoa(foundID),
			expectedCode: http.StatusUnauthorized,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "You can´t see this",
				Data:        nil,
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:          basePath + "/games",
			Method:        http.MethodGet,
			ReqBody:       nil,
			ApplyToken:    true,
			TokenUserRole: tc.role,
		}

		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			handler := game.GetGameByQueryIDAdmin(func(c echo.Context) error {
				return nil
			})
			err := handler(context)
			assert.NoError(t, err)

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

func Test_GetByID_AdminMiddleware_goodCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := NewGame(mockUserRepo)

	foundID := 1

	mockUserRepo.
		EXPECT().
		GetById(gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name   string
		gameID string
		role   string
	}{
		{
			name:   "super admin case",
			role:   roles.SUPER_ADMIN,
			gameID: strconv.Itoa(foundID),
		},
		{
			name:   "admin case",
			role:   roles.ADMIN,
			gameID: strconv.Itoa(foundID),
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:          basePath + "/games",
			Method:        http.MethodGet,
			ReqBody:       nil,
			ApplyToken:    true,
			TokenUserRole: tc.role,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			assert.Nil(t, context.Get("Game"))

			handler := game.GetGameByQueryIDAdmin(func(c echo.Context) error {
				return nil
			})
			err := handler(context)
			assert.NoError(t, err)

			// Test Cases
			assert.NotNil(t, context.Get("Game"))
		})
	}
}

func Test_GetByID_FreeMiddleware_badCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := NewGame(mockUserRepo)

	notFoundID := 2

	mockUserRepo.
		EXPECT().
		GetById(gomock.Eq(notFoundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 0}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		gameID          string
		role            string
		expectedResonse responses.Response
		expectedCode    int
	}{
		{
			name:         "No valid id",
			gameID:       "asdasd",
			role:         roles.USER,
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Invalid ID, Must be a number",
				Data:        nil,
			},
		},
		{
			name:         "Not found id",
			gameID:       strconv.Itoa(notFoundID),
			role:         roles.USER,
			expectedCode: http.StatusNotFound,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     fmt.Sprintf("Game %d not found", notFoundID),
				Data:        nil,
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:          basePath + "/games",
			Method:        http.MethodGet,
			ReqBody:       nil,
			ApplyToken:    true,
			TokenUserRole: tc.role,
		}

		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			handler := game.GetGameByQueryIDFree(func(c echo.Context) error {
				return nil
			})
			err := handler(context)
			assert.NoError(t, err)

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

func Test_GetByID_FreeMiddleware_goodCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := NewGame(mockUserRepo)

	foundID := 1

	mockUserRepo.
		EXPECT().
		GetById(gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name   string
		gameID string
		role   string
	}{
		{
			name:   "user case",
			role:   roles.USER,
			gameID: strconv.Itoa(foundID),
		},
		{
			name:   "admin case",
			role:   roles.ADMIN,
			gameID: strconv.Itoa(foundID),
		},
		{
			name:   "superAdmin case",
			role:   roles.SUPER_ADMIN,
			gameID: strconv.Itoa(foundID),
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:          basePath + "/games",
			Method:        http.MethodGet,
			ReqBody:       nil,
			ApplyToken:    true,
			TokenUserRole: tc.role,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			assert.Nil(t, context.Get("Game"))

			handler := game.GetGameByQueryIDFree(func(c echo.Context) error {
				return nil
			})
			err := handler(context)
			assert.NoError(t, err)

			// Test Cases
			assert.NotNil(t, context.Get("Game"))
		})
	}
}
