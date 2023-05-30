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
	TokenUserID int
	Path        string
	Method      string
	ReqBody     map[string]interface{}
	ApplyToken  bool
}

type TestData struct {
	Echo         *echo.Echo
	Request      *http.Request
	Recoder      *httptest.ResponseRecorder
	Context      *echo.Context
	CreatedToken string
}

// TODO: Aplicar este metodo al resto de tests <29-05-23, yourname> //
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
	token := GenerateToken(config.TokenUserID)

	request.Header.Set("Content-Type", "application/json")
	if config.ApplyToken {
		request.Header.Set("Authorization", token)
	}

	writer := httptest.NewRecorder()
	context := e.NewContext(request, writer)

	middlewareToken := midlewares.ValidateToken(func(c echo.Context) error {
		return nil
	})
	err = middlewareToken(context)
	if err != nil {
		log.Println("Error al genererar el Token test: ", err)
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
