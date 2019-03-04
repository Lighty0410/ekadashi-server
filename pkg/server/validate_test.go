package server

import (
	"fmt"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	tt := []struct {
		testValue     string
		value         loginRequest
		expectedError error
	}{
		{
			testValue:     "spaces only",
			value:         loginRequest{Username: "   ", Password: "    "},
			expectedError: fmt.Errorf("fields username contain latin characters and numbers without space only"),
		},
		{
			testValue:     "space in username",
			value:         loginRequest{Username: "  ", Password: "hahaimcrying"},
			expectedError: fmt.Errorf("fields username contain latin characters and numbers without space only"),
		},
		{
			testValue:     "space in password",
			value:         loginRequest{Username: "imstillcrying", Password: "   "},
			expectedError: fmt.Errorf("fields password contain latin characters and numbers without space only"),
		},
		{
			testValue:     "ASCII symbols in username",
			value:         loginRequest{Username: "@(*#&!", Password: "passwordexist"},
			expectedError: fmt.Errorf("fields username contain latin characters and numbers without space only"),
		},
		{
			testValue:     "ASCII symbols in password",
			value:         loginRequest{Username: "passwordexist", Password: "@(*#&!"},
			expectedError: fmt.Errorf("fields password contain latin characters and numbers without space only"),
		},
		{
			testValue:     "non-latin symbols in username",
			value:         loginRequest{Username: "待望のオリジナルアル", Password: "whatsthis"},
			expectedError: fmt.Errorf("fields username contain latin characters and numbers without space only"),
		},
		{
			testValue:     "non-latin symbols in password ",
			value:         loginRequest{Username: "stilldontknow", Password: "のオリジナルアル"},
			expectedError: fmt.Errorf("fields password contain latin characters and numbers without space only"),
		},
		{
			testValue:     "password less than 6 characters",
			value:         loginRequest{Username: "whoiswho", Password: "iam"},
			expectedError: fmt.Errorf("field password could not be less than 6 characters"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.testValue, func(t *testing.T) {
			err := tc.value.validateRequest().Error()
			if err != tc.expectedError.Error() {
				t.Fatal("\nFor: ", tc.testValue,
					"\nExpected: ", tc.expectedError,
					"\nGot: ", err)
			}
		})
	}
}
