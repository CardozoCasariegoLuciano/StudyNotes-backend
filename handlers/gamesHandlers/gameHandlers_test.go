package gameshandlers_test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	gameshandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/gamesHandlers"
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

// AllGames endpoint
func Test_AllGames_HaveGames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := gameshandlers.NewGame(mockUserRepo)

	expectedGames := &dbmodels.Games{
		dbmodels.Game{
			Title:       "test",
			Description: "test",
		},
		dbmodels.Game{
			Title:       "test2",
			Description: "test2",
		},
	}

	mockUserRepo.
		EXPECT().
		GetAll(gomock.AssignableToTypeOf(&dbmodels.Games{})).
		Do(func(games *dbmodels.Games) {
			*games = *expectedGames
		}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "When have elements",
			expectedCode: http.StatusOK,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "All games",
				Data:        expectedGames,
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:       "/api/v1" + "/games",
			Method:     http.MethodGet,
			ApplyToken: false,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			err := game.GetGames(context)
			assert.NoError(t, err)

			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			//Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, len(resp.Data.([]interface{})), len(*expectedGames))
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}

func Test_AllGames_EmptyGames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := gameshandlers.NewGame(mockUserRepo)

	mockUserRepo.
		EXPECT().
		GetAll(gomock.AssignableToTypeOf(&dbmodels.Games{})).
		Return(nil).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "When no have elements",
			expectedCode: http.StatusOK,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "All games",
				Data:        []interface{}{},
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:       "/api/v1" + "/games",
			Method:     http.MethodGet,
			ApplyToken: false,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			err := game.GetGames(context)
			assert.NoError(t, err)

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

// GetGameByID endpoint
func Test_GameByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := gameshandlers.NewGame(mockUserRepo)

	foundID := 1
	mockUserRepo.
		EXPECT().
		GetById(gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		gameID          string
		userRole        string
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "Must show the game",
			gameID:       strconv.Itoa(foundID),
			userRole:     roles.ADMIN,
			expectedCode: http.StatusOK,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "Game Selected",
				Data:        nil,
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:          "/api/v1" + "/games",
			Method:        http.MethodGet,
			ReqBody:       nil,
			ApplyToken:    true,
			TokenUserRole: tc.userRole,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			midleware := game.GetGameByQueryIDAdmin(func(c echo.Context) error {
				return nil
			})
			err := midleware(context)
			assert.NoError(t, err)

			err = game.GetGameByID(context)
			assert.NoError(t, err)

			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			//Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Contains(t, resp.Data, "title")
			assert.Contains(t, resp.Data, "ID")
			assert.Contains(t, resp.Data, "description")
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}
