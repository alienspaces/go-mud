package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	tests := map[string]struct {
		testErrors     []error
		expectedCount  int
		expectedString string
	}{
		"0 errors": {
			testErrors:     []error{},
			expectedCount:  0,
			expectedString: "",
		},
		"1 error": {
			testErrors: []error{
				errors.New("e1"),
			},
			expectedCount:  1,
			expectedString: "e1\n",
		},
		"multiple errors": {
			testErrors: []error{
				errors.New("e1"),
				errors.New("e2"),
			},
			expectedCount:  2,
			expectedString: "e1\ne2\n",
		},
	}

	for tcName, tc := range tests {
		t.Logf("Running test >%s<", tcName)

		t.Run(tcName, func(t *testing.T) {
			allErrors := Error{}

			for _, e := range tc.testErrors {
				allErrors.Add(e)
			}

			require.Equal(t, tc.expectedCount, len(allErrors.Errors))
			require.Equal(t, tc.expectedCount, allErrors.Count())
			require.Equal(t, tc.expectedString, allErrors.Error())
		})
	}
}
