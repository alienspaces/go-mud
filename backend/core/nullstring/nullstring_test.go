package nullstring

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/convert"
)

func TestFromString(t *testing.T) {

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
			converted := FromString(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.String)
		})
	}
}

func TestFromStringPtr(t *testing.T) {

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
			converted := FromStringPtr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.String)
			}
		})
	}
}

func TestToStringPtr(t *testing.T) {

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
			converted := ToStringPtr(tc.ns)

			require.Equal(t, tc.expected, converted)
		})
	}
}
