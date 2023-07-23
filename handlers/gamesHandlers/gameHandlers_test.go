package gameshandlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	gameshandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/gamesHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
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

	userOwner := 100

	expectedGames := &dbmodels.Games{
		dbmodels.Game{
			Title:       "test",
			Description: "test",
			UserID:      userOwner,
		},
		dbmodels.Game{
			Title:       "test2",
			Description: "test2",
			UserID:      userOwner,
		},
	}

	mockUserRepo.
		EXPECT().
		GetAllGames(gomock.Eq(userOwner), gomock.AssignableToTypeOf(&dbmodels.Games{})).
		Do(func(userOwner int, games *dbmodels.Games) {
			*games = *expectedGames
		}).
		Times(1)

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		expectedCode    int
		userID          int
		expectedResonse responses.Response
	}{
		{
			name:         "when user is owner",
			expectedCode: http.StatusOK,
			userID:       userOwner,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "All games",
				Data:        expectedGames,
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/games",
			Method:      http.MethodGet,
			ApplyToken:  true,
			TokenUserID: tc.userID,
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

	userID := 100

	mockUserRepo.
		EXPECT().
		GetAllGames(gomock.Eq(userID), gomock.AssignableToTypeOf(&dbmodels.Games{})).
		Return(nil).
		Times(1)

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		userID          int
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "When no have elements",
			expectedCode: http.StatusOK,
			userID:       100,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "All games",
				Data:        []interface{}{},
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/games",
			Method:      http.MethodGet,
			ApplyToken:  true,
			TokenUserID: tc.userID,
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

			// Test Cases
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

	userID := 100

	foundID := 1
	mockUserRepo.
		EXPECT().
		GetGameById(gomock.Eq(userID), gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		Times(1)

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		gameID          string
		userID          int
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "Must show the game",
			gameID:       strconv.Itoa(foundID),
			userID:       userID,
			expectedCode: http.StatusOK,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "Game Selected",
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/games",
			Method:      http.MethodGet,
			ReqBody:     nil,
			ApplyToken:  true,
			TokenUserID: userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			midleware := game.GetGameByQueryID(func(c echo.Context) error {
				return nil
			})
			err := midleware(context)
			assert.NoError(t, err)

			err = game.GetGameByID(context)
			assert.NoError(t, err)

			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Contains(t, resp.Data, "title")
			assert.Contains(t, resp.Data, "id")
			assert.Contains(t, resp.Data, "user_id")
			assert.Contains(t, resp.Data, "description")
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}

// CreateGame endpoint
func Test_CreateGame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := gameshandlers.NewGame(mockUserRepo)

	mockUserRepo.
		EXPECT().
		Save(gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(nil).
		Times(1)

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		isBadCase       bool
		reqBody         map[string]interface{}
		userID          int
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "empty field in reqBody",
			isBadCase:    true,
			expectedCode: http.StatusBadRequest,
			reqBody:      map[string]interface{}{},
			userID:       100,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "All fields are required",
				Data:        nil,
			},
		},
		{
			name:         "Wrong data in reqBody",
			isBadCase:    true,
			expectedCode: http.StatusBadRequest,
			userID:       100,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Not valid body information",
				Data:        nil,
			},
			reqBody: map[string]interface{}{
				"title":       123,
				"description": 123,
			},
		},
		{
			name:         "Short title",
			isBadCase:    true,
			expectedCode: http.StatusBadRequest,
			userID:       100,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Title field must have more than 3 characters",
				Data:        nil,
			},
			reqBody: map[string]interface{}{
				"title":       "ho",
				"description": "",
			},
		},
		{
			name:         "OK",
			isBadCase:    false,
			expectedCode: http.StatusCreated,
			userID:       100,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "Game created",
				Data:        dbmodels.Game{},
			},
			reqBody: map[string]interface{}{
				"title":       "valid Title",
				"description": "description",
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/games",
			Method:      http.MethodPost,
			ReqBody:     tc.reqBody,
			ApplyToken:  true,
			TokenUserID: tc.userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			err := game.CreateGame(context)
			assert.NoError(t, err)

			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// Test Cases
			if tc.isBadCase {
				assert.Equal(t, tc.expectedCode, writer.Code)
				assert.Equal(t, tc.expectedResonse.Data, resp.Data)
			} else {
				assert.Contains(t, resp.Data, "title")
				assert.Contains(t, resp.Data, "id")
				assert.Contains(t, resp.Data, "description")
				assert.Contains(t, resp.Data, "user_id")
			}
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}

// EditGame endpoint
func Test_EditGame(t *testing.T) {
	// Test Constatns
	gameID := 1
	userID := 100

	// TestCases
	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		isBadCase       bool
		gameID          string
		userID          int
		reqBody         map[string]interface{}
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "empty field in reqBody",
			isBadCase:    true,
			userID:       userID,
			gameID:       strconv.Itoa(gameID),
			expectedCode: http.StatusBadRequest,
			reqBody:      map[string]interface{}{},
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "All fields are required",
				Data:        nil,
			},
		},
		{
			name:         "Wrong data in reqBody",
			isBadCase:    true,
			userID:       userID,
			gameID:       strconv.Itoa(gameID),
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Not valid body information",
				Data:        nil,
			},
			reqBody: map[string]interface{}{
				"title":       123,
				"description": 123,
			},
		},
		{
			name:         "Short title",
			isBadCase:    true,
			gameID:       strconv.Itoa(gameID),
			userID:       userID,
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Title field must have more than 3 characters",
				Data:        nil,
			},
			reqBody: map[string]interface{}{
				"title":       "ho",
				"description": "",
			},
		},
		{
			name:         "OK",
			isBadCase:    false,
			userID:       userID,
			gameID:       strconv.Itoa(gameID),
			expectedCode: http.StatusCreated,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "Game edited",
			},
			reqBody: map[string]interface{}{
				"title":       "valid Title",
				"description": "description",
			},
		},
	}

	// Mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := gameshandlers.NewGame(mockUserRepo)

	mockUserRepo.
		EXPECT().
		Save(gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(nil).
		Times(1)

	mockUserRepo.
		EXPECT().
		GetGameById(gomock.Eq(userID), gomock.Eq(gameID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		Times(len(testCases))

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/games",
			Method:      http.MethodPut,
			ReqBody:     tc.reqBody,
			ApplyToken:  true,
			TokenUserID: userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:gameID")
			context.SetParamNames("gameID")
			context.SetParamValues(tc.gameID)

			midleware := game.GetGameByQueryID(func(c echo.Context) error {
				return nil
			})
			err := midleware(context)
			assert.NoError(t, err)

			err = game.EditGame(context)
			assert.NoError(t, err)

			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// Test Cases
			if tc.isBadCase {
				assert.Equal(t, tc.expectedCode, writer.Code)
				assert.Equal(t, tc.expectedResonse.Data, resp.Data)
			} else {
				assert.Contains(t, resp.Data, "title")
				assert.Contains(t, resp.Data, "description")
				assert.Contains(t, resp.Data, "id")
				assert.Contains(t, resp.Data, "user_id")
			}
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}

// DeleteGame endpoint
func Test_DeletGame(t *testing.T) {
	foundID := 1
	userID := 100
	NouserID := 200

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		gameID          string
		userID          int
		isAuthor        bool
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "Must delete the game",
			gameID:       strconv.Itoa(foundID),
			userID:       userID,
			isAuthor:     true,
			expectedCode: http.StatusOK,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "Game Deleted",
			},
		},
		{
			name:         "Must not find the game",
			gameID:       strconv.Itoa(foundID),
			userID:       NouserID,
			isAuthor:     false,
			expectedCode: http.StatusNotFound,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     fmt.Sprintf("Game %d not found", foundID),
				Data:        nil,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	game := gameshandlers.NewGame(mockUserRepo)

	mockUserRepo.
		EXPECT().
		GetGameById(gomock.Eq(userID), gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 1}).
		Times(1)

	mockUserRepo.
		EXPECT().
		GetGameById(gomock.Eq(NouserID), gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(&gorm.DB{RowsAffected: 0}).
		Times(1)

	mockUserRepo.
		EXPECT().
		DeleteByID(gomock.AssignableToTypeOf(foundID), gomock.AssignableToTypeOf(&dbmodels.Game{})).
		Return(nil).
		Times(1)

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

			midleware := game.GetGameByQueryID(func(c echo.Context) error {
				return nil
			})
			err := midleware(context)
			assert.NoError(t, err)

			if tc.isAuthor {
				err = game.DeleteGame(context)
				assert.NoError(t, err)
			}

			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// Test Cases
			if !tc.isAuthor {
				assert.Equal(t, tc.expectedResonse.Data, resp.Data)
			} else {
				assert.Contains(t, resp.Data, "title")
				assert.Contains(t, resp.Data, "description")
				assert.Contains(t, resp.Data, "id")
				assert.Contains(t, resp.Data, "user_id")
			}
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
		})
	}
}
