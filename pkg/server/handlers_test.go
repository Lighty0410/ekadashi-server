package server

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Lighty0410/ekadashi-server/pkg/mongo"
	"github.com/stretchr/testify/assert"
)

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
			expectedMessage:  "{\"reason\":\"field username contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with space inside field password",
			request: `{"username":"1234username",
              "password":"12414pa  ssword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field password contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with ASCII symbols inside field password",
			request: `{"username":"1234username",
              "password":"12414pa  ssword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field password contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with ASCII symbols inside field username",
			request: `{"username":"123user@&!*",
              "password":"12414password"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field username contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with ASCII symbols inside field password",
			request: `{"username":"iamuserhah",
              "password":"123pas@(&!(s"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field password contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with empty username",
			request: `{"username":"    ",
              "password":"12414password"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field username contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with empty password",
			request: `{"username":"usernamedadadada",
              "password":"   "}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field password contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with ASCII symbols, numbers and space inside field username",
			request: `{"username":"123pa  s@(&!(s",
              "password":"whatsthispassword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field username contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 with ASCII symbols, numbers and space inside field password",
			request: `{"username":"whoisusername",
              "password":"123pa  s@(&!(s"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field password contain latin characters and numbers without space only\"}\n",
		},
		{
			name: "status 400 when field username is too short",
			request: `{"username":"abc",
              "password":"whatsthispassword"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field username could not be less than 6 characters\"}\n",
		},
		{
			name: "status 400 when field password is too short",
			request: `{"username":"iamnotauser",
              "password":"abc"}`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"field password could not be less than 6 characters\"}\n",
		},
		{
			name: "status 409 when username already exists",
			request: `{"username":"passwordpass",
              "password":"asdqweas"}`,
			expectedResponse: http.StatusConflict,
			expectedMessage:  "{\"reason\":\"user already exists\"}\n",
		},
		{
			name: "another status 409 when username already exists",
			request: `{"username":"1234username",
              "password":"12414password"}`,
			expectedResponse: http.StatusConflict,
			expectedMessage:  "{\"reason\":\"user already exists\"}\n",
		},
	}
	handler := createHandler()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			bodyreq := bytes.NewBuffer([]byte(tc.request))
			req := httptest.NewRequest("POST", "/register", bodyreq) // todo user httptest
			handler.handleRegistration(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			responseBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, res.StatusCode, tc.expectedResponse, tc.name)
			assert.Equal(t, string(responseBody), tc.expectedMessage, "unexpected body")
		})
	}
}

func TestLogin(t *testing.T) {
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
			expectedMessage:  "{\"reason\":\"can not decode the request: unexpected EOF\"}\n",
		},
		{
			name: "another status 400 with incorrect body request",
			requst: `"1234username",
              "password":"12414password"`,
			expectedResponse: http.StatusBadRequest,
			expectedMessage:  "{\"reason\":\"can not decode the request: json: cannot unmarshal string into Go value of type server.loginRequest\"}\n",
		},
		{
			name: "status 401 when user doesn't exist",
			requst: `{"username":"newusernameus",
              "password":"newpasswordus"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "{\"reason\":\"incorrect username or password: mongo: no documents in result\"}\n",
		},
		{
			name: "another status 401 when user doesn't exist ",
			requst: `{"username":"newusernameus",
              "password":"userpass"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "{\"reason\":\"incorrect username or password: mongo: no documents in result\"}\n",
		},
		{
			name: "status 401 when field password is incorrect",
			requst: `{"username":"passwordpass",
				"password":"qwihoaslk"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "{\"reason\":\"incorrect username or password: crypto/bcrypt: hashedPassword is not the hash of the given password\"}\n",
		},
		{
			name: "another status 401 when field password is incorrect",
			requst: `{"username":"1234username",
              "password":"asda21"}`,
			expectedResponse: http.StatusUnauthorized,
			expectedMessage:  "{\"reason\":\"incorrect username or password: crypto/bcrypt: hashedPassword is not the hash of the given password\"}\n",
		},
	}
	handler := createHandler()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			bodyreq := bytes.NewBuffer([]byte(tc.requst))
			req := httptest.NewRequest("POST", "/register", bodyreq)
			handler.handleLogin(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			responseBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, res.StatusCode, tc.expectedResponse, tc.name)
			assert.Equal(t, string(responseBody), tc.expectedMessage, "unexpected body")
		})
	}
}
