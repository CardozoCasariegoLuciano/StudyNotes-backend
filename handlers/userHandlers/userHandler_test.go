package userhandlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	userhandlers "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/userHandlers"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	testtools "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/testTools"
	mock_models "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/testTools/mocks"
	apimodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/apiModels"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_GetUser(t *testing.T) {
	userID := 10

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		userID          int
		isBadCase       bool
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "get user data",
			expectedCode: http.StatusOK,
			userID:       userID,
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "User data",
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	user := userhandlers.NewUser(mockUserRepo)

	mockUserRepo.
		EXPECT().
		GetUserByID(gomock.Eq(userID), gomock.AssignableToTypeOf(&dbmodels.User{})).
		Return(&gorm.DB{RowsAffected: 1}).
		Times(1)

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/user",
			Method:      http.MethodGet,
			ApplyToken:  true,
			TokenUserID: tc.userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			err := user.GetUser(context)
			assert.NoError(t, err)

			resp := responses.Response{}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			//Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
			assert.Contains(t, resp.Data, "email")
			assert.Contains(t, resp.Data, "role")
			assert.Contains(t, resp.Data, "name")
			assert.Contains(t, resp.Data, "id")
			assert.NotContains(t, resp.Data, "password")
		})
	}
}

func Test_EditUser(t *testing.T) {
	userID := 10
	newName := "fulano"
	newImage := "imagen"
	expectedData := map[string]interface{}{
		"name":  newName,
		"email": "initialEmail",
		"image": newImage,
		"role":  "initialRole",
		"id":    float64(0),
	}

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		userID          int
		reqBody         map[string]interface{}
		expectedCode    int
		expectedResonse responses.Response
	}{
		{
			name:         "OK edit user",
			expectedCode: http.StatusOK,
			userID:       userID,
			reqBody: map[string]interface{}{
				"name":  newName,
				"image": newImage,
			},
			expectedResonse: responses.Response{
				MessageType: "OK",
				Message:     "User edited",
				Data:        expectedData,
			},
		},
		{
			name:         "Wrong body data",
			expectedCode: http.StatusBadRequest,
			userID:       userID,
			reqBody: map[string]interface{}{
				"name":  123,
				"image": 123,
			},
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Not valid body information",
				Data:        nil,
			},
		},
		{
			name:         "Empty body data",
			expectedCode: http.StatusBadRequest,
			userID:       userID,
			reqBody:      map[string]interface{}{},
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "There are required fields",
				Data:        nil,
			},
		},
		{
			name:         "max limit in Name field",
			expectedCode: http.StatusBadRequest,
			userID:       userID,
			reqBody: map[string]interface{}{
				"name": "qweqqweqweqweqweasdqwasdqwdasdqwdasdwe"},
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Name field must have less than 30 characters",
				Data:        nil,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	user := userhandlers.NewUser(mockUserRepo)

	mockUserRepo.
		EXPECT().
		GetUserByID(gomock.Eq(userID), gomock.AssignableToTypeOf(&dbmodels.User{})).
		DoAndReturn(func(userID int, user *dbmodels.User) *gorm.DB {
			*user = dbmodels.User{
				Name:  "initialName",
				Email: "initialEmail",
				Image: "initialImagen",
				Role:  "initialRole",
			}
			return &gorm.DB{RowsAffected: 1}
		}).
		Times(len(testCases))

	mockUserRepo.
		EXPECT().
		Save(gomock.AssignableToTypeOf(&dbmodels.User{})).
		Return(nil).
		Times(1)

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        "/api/v1" + "/user",
			Method:      http.MethodPut,
			ApplyToken:  true,
			ReqBody:     tc.reqBody,
			TokenUserID: tc.userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			err := user.EditUser(context)
			assert.NoError(t, err)

			resp := responses.Response{
				MessageType: "",
				Message:     "",
				Data:        apimodels.User{},
			}
			err = json.Unmarshal(writer.Body.Bytes(), &resp)
			assert.NoError(t, err)

			//Test Cases
			assert.Equal(t, tc.expectedCode, writer.Code)
			assert.Equal(t, tc.expectedResonse.MessageType, resp.MessageType)
			assert.Equal(t, tc.expectedResonse.Message, resp.Message)
			assert.Equal(t, resp.Data, tc.expectedResonse.Data)
		})
	}
}
