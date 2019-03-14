package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	tt := []struct {
		testValue     string
		value         loginRequest
		expectedError error
	}{
		{
			testValue:     "ordinary username and password",
			value:         loginRequest{Username: "Username", Password: "Password"},
			expectedError: nil,
		},
		{
			testValue:     "numbers only",
			value:         loginRequest{Username: "1234567890", Password: "numbernununm"},
			expectedError: nil,
		},
		{
			testValue:     "numbers in revert order",
			value:         loginRequest{Username: "0987654321", Password: "numnumn"},
			expectedError: nil,
		},
		{
			testValue:     "spaces only",
			value:         loginRequest{Username: "   ", Password: "    "},
			expectedError: fmt.Errorf("field username contain latin characters and numbers without space only"),
		},
		{
			testValue:     "ASCII symbols in username",
			value:         loginRequest{Username: "@(*#&!", Password: "passwordexist"},
			expectedError: fmt.Errorf("field username contain latin characters and numbers without space only"),
		},
		{
			testValue:     "ASCII symbols in password",
			value:         loginRequest{Username: "passwordexist", Password: "@(*#&!"},
			expectedError: fmt.Errorf("field password contain latin characters and numbers without space only"),
		},
		{
			testValue:     "non-latin symbols in username",
			value:         loginRequest{Username: "待望のオリジナルアル", Password: "whatsthis"},
			expectedError: fmt.Errorf("field username contain latin characters and numbers without space only"),
		},
		{
			testValue:     "non-latin symbols in password ",
			value:         loginRequest{Username: "stilldontknow", Password: "のオリジナルアル"},
			expectedError: fmt.Errorf("field password contain latin characters and numbers without space only"),
		},
		{
			testValue:     "password less than 6 characters",
			value:         loginRequest{Username: "whoiswho", Password: "iam"},
			expectedError: fmt.Errorf("field password could not be less than 6 characters"),
		}, {
			testValue:     "username less than 6 characters",
			value:         loginRequest{Username: "yikes", Password: "smtwentwrong"},
			expectedError: fmt.Errorf("field username could not be less than 6 characters"),
		},
		{
			testValue:     "space between characters",
			value:         loginRequest{Username: "kekeke", Password: "kek kek"},
			expectedError: fmt.Errorf("field password contain latin characters and numbers without space only"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.testValue, func(t *testing.T) {
			err := tc.value.validateRequest()
			assert.Equal(t, err, tc.expectedError)
		})
	}
}
