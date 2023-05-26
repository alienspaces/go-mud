package null

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/backend/core/convert"
)

func TestNullStringFromString(t *testing.T) {

	tests := map[string]struct {
		value         string
		expectedValid bool
	}{
		"valid value": {
			value:         "abc",
			expectedValid: true,
		},
		"invalid value": {
			value:         "",
			expectedValid: false,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullStringFromString(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.String)
		})
	}
}

func TestNullStringFromStringPtr(t *testing.T) {

	tests := map[string]struct {
		value         *string
		expectedValid bool
	}{
		"valid value": {
			value:         convert.Stringp("abc"),
			expectedValid: true,
		},
		"valid empty": {
			value:         convert.Stringp(""),
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
			converted := NullStringFromStringPtr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.String)
			}
		})
	}
}

func TestNullStringToStringPtr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullString
		expected *string
	}{
		"valid": {
			ns: sql.NullString{
				String: "abc",
				Valid:  true,
			},
			expected: convert.Stringp("abc"),
		},
		"invalid": {
			ns: sql.NullString{
				String: "abc",
				Valid:  false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullStringToStringPtr(tc.ns)

			require.Equal(t, tc.expected, converted)
		})
	}
}
