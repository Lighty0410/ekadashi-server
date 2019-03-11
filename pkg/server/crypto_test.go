package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrypto(t *testing.T) {
	tt := []struct {
		testValue   string
		password    string
		expectError error
	}{
		{
			testValue: "casual password",
			password:  "hopeisverything,",
		},
		{
			testValue: "UPPERCASE",
			password:  "WHATSMYSECRETKEY",
		},
		{
			testValue: "empty string",
			password:  "",
		},
		{
			testValue: "short password",
			password:  "shrt",
		},
		{
			testValue: "a lot of numbers",
			password:  "120938219382109381",
		},
		{
			testValue: "very long password",
			password:  "OnMouseMoveFunctionalTestVerticalSplitIndicatorExactlyOnTheLeftBorderOfTheFirstCellOnTheTheWeekViewAndGroupByResourceAndTwoResources",
		},
		{
			testValue: "ASCII symbols",
			password:  "!@#$%^&*()_+",
		},
	}
	for _, tc := range tt {
		hash, err := generateHash(tc.password)
		assert.NoError(t, err, tc.expectError)
		err = compareHash(hash, []byte(tc.password))
		assert.NoError(t, err, tc.expectError)
	}
}
