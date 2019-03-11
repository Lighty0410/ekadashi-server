package server

import (
	"bytes"
	"encoding/json"
	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
	"gotest.tools/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type message struct {
	Reason string `json:"reason"`
}

const ekadashiURL = "EKADASHI_MONGO_URL"

func createHandler() *EkadashiServer {
	connectionURL := os.Getenv(ekadashiURL)
	if connectionURL == "" {
		log.Fatalf("Innapropriate %v variable for mongoDB connection", ekadashiURL)
	}
	mongoService, err := mongo.NewService(connectionURL)
	if err != nil {
		log.Fatalf("Could not create mongo service: %v", err)
	}
	testEkadashi, err := NewEkadashiServer(mongoService)
	if err != nil {
		log.Fatalf("Could not create ekadashi server: %v", err)
	}
	return testEkadashi
}

func TestRegisterFunc(t *testing.T) {
	tt := []struct {
		name             string
		request          string
		expectedResponse int
		expectedMessage  string
	}{
		{
			name: "status 200 with latin symbols only",
			request: `{"username":"passwordpass",
              "password":"asdqweas"}`,
			expectedResponse: http.StatusOK,
			expectedMessage:  "",
		},
		{
			name: "status 200 with numbers inside both fields",
			request: `{"username":"1234username",
              "password":"12414password"}`,
			expectedResponse: http.StatusOK,
			expectedMessage:  "",
		},
		{
			name: "status 400 with space inside field username",
			request: `{"username":"1234u sername",
              "password":"12414password"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field username contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with space inside field password",
			request: `{"username":"1234username",
              "password":"12414pa  ssword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field password contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with ASCII symbols inside field password",
			request: `{"username":"1234username",
              "password":"12414pa  ssword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field password contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with ASCII symbols inside field username",
			request: `{"username":"123user@&!*",
              "password":"12414password"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field username contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with ASCII symbols inside field password",
			request: `{"username":"iamuserhah",
              "password":"123pas@(&!(s"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field password contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with empty username",
			request: `{"username":"    ",
              "password":"12414password"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field username contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with empty password",
			request: `{"username":"usernamedadadada",
              "password":"   "}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field password contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with ASCII symbols, numbers and space inside field username",
			request: `{"username":"123pa  s@(&!(s",
              "password":"whatsthispassword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field username contain latin characters and numbers without space only",
		},
		{
			name: "status 400 with ASCII symbols, numbers and space inside field password",
			request: `{"username":"whoisusername",
              "password":"123pa  s@(&!(s"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field password contain latin characters and numbers without space only",
		},
		{
			name: "status 400 when field username is too short",
			request: `{"username":"abc",
              "password":"whatsthispassword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field username could not be less than 6 characters",
		},
		{
			name: "status 400 when field password is too short",
			request: `{"username":"iamnotauser",
              "password":"abc"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "field password could not be less than 6 characters",
		},
	}
	handler := createHandler()
	for _, tc := range tt {
		rec := httptest.NewRecorder()
		bodyreq := bytes.NewBuffer([]byte(tc.request))
		req, _ := http.NewRequest("POST", "/register", bodyreq)
		handler.handleRegistration(rec, req)
		assert.Equal(t, rec.Result().StatusCode, tc.expectedResponse, tc.name)
		var jsonMessage message
		json.NewDecoder(rec.Body).Decode(&jsonMessage)
		assert.Equal(t, jsonMessage.Reason, tc.expectedMessage, tc.expectedMessage, "\n", tc.name)
	}
}

func TestLoginAndNextEkadashiFunc(t *testing.T) {
	tt := []struct {
		name             string
		requst           string
		expectedResponse int
		expectedMessage  string
	}{
		{
			name: "status 200 with existing username and password",
			requst: `{"username":"passwordpass",
              "password":"asdqweas"}`,
			expectedResponse: http.StatusOK,
			expectedMessage:  "",
		},
		{
			name: "another status 200 with existing username and password",
			requst: `{"username":"1234username",
              "password":"12414password"}`,
			expectedResponse: http.StatusOK,
			expectedMessage:  "",
		},
		{
			name: "status 400 with with incorrect body request",
			requst: `{"username":"1234username",
              "password":"12414password"`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "can not decode the request: unexpected EOF",
		},
		{
			name: "another status 400 with incorrect body request",
			requst: `"1234username",
              "password":"12414password"`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "can not decode the request: json: cannot unmarshal string into Go value of type server.loginRequest",
		},
		{
			name: "status 401 when user doesn't exist",
			requst: `{"username":"newusernameus",
              "password":"newpasswordus"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "incorrect username or password: mongo: no documents in result",
		},
		{
			name: "another status 401 when user doesn't exist ",
			requst: `{"username":"newusernameus",
              "password":"userpass"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "incorrect username or password: mongo: no documents in result",
		},
		{
			name: "status 401 when field password is incorrect",
			requst: `{"username":"passwordpass",
              "password":"qwihoaslk"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "incorrect username or password: crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			name: "another status 401 when field password is incorrect",
			requst: `{"username":"1234username",
              "password":"asda21"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "incorrect username or password: crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
	}
	handler := createHandler()
	for _, tc := range tt {
		rec := httptest.NewRecorder()
		bodyreq := bytes.NewBuffer([]byte(tc.requst))
		req, _ := http.NewRequest("POST", "/register", bodyreq)
		handler.handleLogin(rec, req)
		assert.Equal(t, rec.Result().StatusCode, tc.expectedResponse, tc.name)
		var jsonMessage message
		json.NewDecoder(rec.Body).Decode(&jsonMessage)
		assert.Equal(t, jsonMessage.Reason, tc.expectedMessage, tc.expectedMessage, "\n", tc.name)

		if rec.Result().StatusCode == 200 {
			handler.nextEkadashiHandler(rec, req)
		} // this seems do be undone for a moment. Need to ask
	}
}
