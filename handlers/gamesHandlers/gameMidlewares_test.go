package gameshandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	errorcodes "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/errorCodes"
	testtools "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/testTools"
	mock_models "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/testTools/mocks"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_GetByID_Middleware_badCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := NewGame(mockUserRepo)

	notFoundID := 2
	notOwnerID := 999

	mockUserRepo.
		EXPECT().
		GetGameById(gomock.Eq(notOwnerID), gomock.Eq(notFoundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 0}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		gameID          string
		userID          int
		expectedResonse responses.Response
		expectedCode    int
	}{
		{
			name:         "No valid id",
			gameID:       "asdasd",
			userID:       1,
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: errorcodes.INVALID_ID,
				Message:     "Invalid ID, Must be a number",
				Data:        nil,
			},
		},
		{
			name:         "Not found id",
			gameID:       strconv.Itoa(notFoundID),
			userID:       notOwnerID,
			expectedCode: http.StatusNotFound,
			expectedResonse: responses.Response{
				MessageType: errorcodes.NOT_FOUND_OR_UNAUTHORIZED,
				Message:     fmt.Sprintf("Game %d not found", notFoundID),
				Data:        nil,
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/games",
			Method:      http.MethodGet,
			ReqBody:     nil,
			ApplyToken:  true,
			TokenUserID: tc.userID,
		}

		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			handler := game.GetGameByQueryID(func(c echo.Context) error {
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

func Test_GetByID_Middleware_goodCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := NewGame(mockUserRepo)

	foundID := 1
	userOwner := 19
	notOwner := 1999

	mockUserRepo.
		EXPECT().
		GetGameById(gomock.Eq(userOwner), gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		AnyTimes()

	mockUserRepo.
		EXPECT().
		GetGameById(gomock.Eq(notOwner), gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 0}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name   string
		gameID string
		userID int
	}{
		{
			name:   "owner case",
			userID: userOwner,
			gameID: strconv.Itoa(foundID),
		},
		{
			name:   "Not owner case",
			userID: notOwner,
			gameID: strconv.Itoa(foundID),
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/games",
			Method:      http.MethodGet,
			ReqBody:     nil,
			ApplyToken:  true,
			TokenUserID: tc.userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			assert.Nil(t, context.Get("Game"))

			handler := game.GetGameByQueryID(func(c echo.Context) error {
				return nil
			})
			err := handler(context)
			assert.NoError(t, err)

			resp := responses.Response{
				MessageType: "",
				Message:     "",
				Data:        nil,
			}

			// Test Cases
			if tc.userID == userOwner {
				assert.NotNil(t, context.Get("Game"))
			} else {
				err = json.Unmarshal(writer.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, http.StatusNotFound, writer.Code)
				assert.Equal(t, nil, resp.Data)
				assert.Equal(t, errorcodes.NOT_FOUND_OR_UNAUTHORIZED, resp.MessageType)
				assert.Equal(t, fmt.Sprintf("Game %d not found", foundID), resp.Message)
			}
		})
	}
}
