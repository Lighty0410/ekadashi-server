package server

import "testing"

func TestCrypto(t *testing.T) {
	tt := []struct {
		testValue   string
		password    string
		expectError error
	}{
		{
			testValue:   "casual password",
			password:    "hopeisverything,",
			expectError: nil,
		},
		{
			testValue:   "UPPERCASE",
			password:    "WHATSMYSECRETKEY",
			expectError: nil,
		},
		{
			testValue:   "empty string",
			password:    "",
			expectError: nil,
		},
		{
			testValue:   "short password",
			password:    "shrt",
			expectError: nil,
		},
		{
			testValue:   "a lot of numbers",
			password:    "120938219382109381",
			expectError: nil,
		},
		{
			testValue:   "very long password",
			password:    "OnMouseMoveFunctionalTestVerticalSplitIndicatorExactlyOnTheLeftBorderOfTheFirstCellOnTheTheWeekViewAndGroupByResourceAndTwoResources",
			expectError: nil,
		},
		{
			testValue:   "ASCII symbols",
			password:    "!@#$%^&*()_+",
			expectError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testValue, func(t *testing.T) {
			hash, err := generateHash(tc.password)
			if err != nil {
				t.Fatal(
					"for: ", tc.testValue,
					"expected: ", tc.expectError,
					"got: %v", err)
			}
			err = compareHash(hash, []byte(tc.password))
			if err != nil {
				t.Fatal(
					"for: ", tc.testValue,
					"expected: ", tc.expectError,
					"got: %v", err)
			}
		})
	}
}
