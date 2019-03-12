package server

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	tt := []struct {
		name     string
		password string
	}{
		{
			name:     "casual password",
			password: "hopeisverything,",
		},
		{
			name:     "UPPERCASE",
			password: "WHATSMYSECRETKEY",
		},
		{
			name:     "empty string",
			password: "",
		},
		{
			name:     "short password",
			password: "shrt",
		},
		{
			name:     "a lot of numbers",
			password: "120938219382109381",
		},
		{
			name:     "very long password",
			password: "OnMouseMoveFunctionalTestVerticalSplitIndicatorExactlyOnTheLeftBorderOfTheFirstCellOnTheTheWeekViewAndGroupByResourceAndTwoResources",
		},
		{
			name:     "ASCII symbols",
			password: "!@#$%^&*()_+",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := generateHash(tc.password)
			require.NoError(t, err)
			err = compareHash(hash, []byte(tc.password))
			assert.NoError(t, err)
		})
	}
}
