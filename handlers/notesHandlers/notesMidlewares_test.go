package noteshandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	mock_models "github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/mocks"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/handlers/responses"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/environment"
	testtools "github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/testTools"
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const basePath = "/api/v1"

func TestNotesByIDMidlewares_bad(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)

	note := NewNote(mockUserRepo)

	notFoundID := 2
	foundID := 1
	noAuthorID := 123

	mockUserRepo.
		EXPECT().
		GetNoteByID(gomock.Eq(notFoundID), gomock.AssignableToTypeOf(&dbmodels.Note{})).
		Return(&gorm.DB{RowsAffected: 0}).
		AnyTimes()

	mockUserRepo.
		EXPECT().
		GetNoteByID(gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Note{})).
		Return(&gorm.DB{RowsAffected: 1}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name            string
		noteID          string
		userID          int
		expectedResonse responses.Response
		expectedCode    int
	}{
		{
			name:         "No valid id",
			noteID:       "asdasd",
			userID:       noAuthorID,
			expectedCode: http.StatusBadRequest,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "Invalid ID, Must be a number",
				Data:        nil,
			},
		},
		{
			name:         "Not found id",
			noteID:       strconv.Itoa(notFoundID),
			userID:       noAuthorID,
			expectedCode: http.StatusNotFound,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     fmt.Sprintf("Note %d not found", notFoundID),
				Data:        nil,
			},
		},
		{
			name:         "user not owner",
			noteID:       strconv.Itoa(foundID),
			userID:       noAuthorID,
			expectedCode: http.StatusUnauthorized,
			expectedResonse: responses.Response{
				MessageType: "ERROR",
				Message:     "This is not your note",
				Data:        nil,
			},
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        basePath + "/notes",
			Method:      http.MethodGet,
			ReqBody:     nil,
			ApplyToken:  true,
			TokenUserID: tc.userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context
		writer := testData.Recoder

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:noteId")
			context.SetParamNames("noteId")
			context.SetParamValues(tc.noteID)

			handler := note.GetNoteByQueryID(func(c echo.Context) error {
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

func TestNotesByIDMidlewares_good(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_models.NewMockIstorage(ctrl)
	note := NewNote(mockUserRepo)

	foundID := 1
	authorID := 0

	mockUserRepo.
		EXPECT().
		GetNoteByID(gomock.Eq(foundID), gomock.AssignableToTypeOf(&dbmodels.Note{})).
		Return(&gorm.DB{RowsAffected: 1}).
		AnyTimes()

	environment.SetTestEnvirontment()
	testCases := []struct {
		name   string
		noteID string
		userID int
	}{
		{
			name:   "good case",
			noteID: strconv.Itoa(foundID),
			userID: authorID,
		},
	}

	for _, tc := range testCases {
		testConfig := testtools.InitTestConfig{
			Path:        basePath + "/notes",
			Method:      http.MethodGet,
			ReqBody:     nil,
			ApplyToken:  true,
			TokenUserID: tc.userID,
		}
		testData := testtools.SetGenericTestData(&testConfig)
		context := *testData.Context

		t.Run(tc.name, func(t *testing.T) {
			context.SetPath("/:noteId")
			context.SetParamNames("noteId")
			context.SetParamValues(tc.noteID)

			assert.Nil(t, context.Get("Note"))

			handler := note.GetNoteByQueryID(func(c echo.Context) error {
				return nil
			})
			err := handler(context)
			assert.NoError(t, err)

			//Test Cases
			assert.NotNil(t, context.Get("Note"))
		})
	}
}
