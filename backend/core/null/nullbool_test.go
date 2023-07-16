package null

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/convert"
)

func TestFromBool(t *testing.T) {

	tests := map[string]struct {
		value         bool
		expectedValid bool
	}{
		"valid true value": {
			value:         true,
			expectedValid: true,
		},
		"valid false value": {
			value:         false,
			expectedValid: true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullBoolFromBool(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.Bool)
		})
	}
}

func TestFromBoolPtr(t *testing.T) {

	tests := map[string]struct {
		value         *bool
		expectedValid bool
	}{
		"valid true value": {
			value:         convert.Boolp(true),
			expectedValid: true,
		},
		"valid false value": {
			value:         convert.Boolp(false),
			expectedValid: true,
		},
		"invalid nil": {
			value:         nil,
			expectedValid: false,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullBoolFromBoolPtr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.Bool)
			}
		})
	}
}

func TestToBoolPtr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullBool
		expected *bool
	}{
		"valid true": {
			ns: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			expected: convert.Boolp(true),
		},
		"valid false": {
			ns: sql.NullBool{
				Bool:  false,
				Valid: true,
			},
			expected: convert.Boolp(false),
		},
		"invalid": {
			ns: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullBoolToBoolPtr(tc.ns)

			require.Equal(t, tc.expected, converted)
		})
	}
}
