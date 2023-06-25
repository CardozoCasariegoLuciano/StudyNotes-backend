package testtools

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/helpers/customValidators"
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/midlewares"
	"github.com/labstack/echo/v4"
)

type InitTestConfig struct {
	TokenUserID    int
	TokenUserRole  string
	TokenUserEmail string
	Path           string
	Method         string
	ReqBody        map[string]interface{}
	ApplyToken     bool
}

type TestData struct {
	Echo         *echo.Echo
	Request      *http.Request
	Recoder      *httptest.ResponseRecorder
	Context      *echo.Context
	CreatedToken string
}

// PARA MAS ADELANTE: ver si hace falta crear/modificar este para que acepte otros middlewares
func SetGenericTestData(config *InitTestConfig) *TestData {
	//Init objects
	e := echo.New()
	e.Validator = customValidators.NewCustomValidator()

	//Parse map of strings to []bytes of the body
	body, err := json.Marshal(config.ReqBody)
	if err != nil {
		log.Println("Error at marshal ReqBody")
	}

	//Create new request, recorder(writer) and contetx
	request := httptest.NewRequest(
		config.Method,
		config.Path,
		strings.NewReader(string(body)),
	)

	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	context := e.NewContext(request, writer)

	var token string
	if config.ApplyToken {
		token = GenerateToken(
			config.TokenUserID,
			config.TokenUserRole,
			config.TokenUserEmail,
		)
		request.Header.Set("Authorization", token)

		middlewareToken := midlewares.ValidateToken(func(c echo.Context) error {
			return nil
		})
		err = middlewareToken(context)
		if err != nil {
			log.Println("Error al genererar el Token test: ", err)
		}
	}

	returnValues := &TestData{
		Echo:         e,
		Request:      request,
		Context:      &context,
		Recoder:      writer,
		CreatedToken: token,
	}

	return returnValues
}
