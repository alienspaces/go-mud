package nullint

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/convert"
)

func TestFromInt16(t *testing.T) {

	tests := map[string]struct {
		value         int16
		expectedValid bool
	}{
		"valid value": {
			value:         1234,
			expectedValid: true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := FromInt16(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.Int16)
		})
	}
}

func TestFromInt16Ptr(t *testing.T) {

	tests := map[string]struct {
		value         *int16
		expectedValid bool
	}{
		"valid non-zero": {
			value:         convert.Int16p(1234),
			expectedValid: true,
		},
		"valid zero": {
			value:         convert.Int16p(0),
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
			converted := FromInt16Ptr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.Int16)
			}
		})
	}
}

func TestToInt16Ptr(t *testing.T) {

	tests := map[string]struct {
		ni       sql.NullInt16
		expected *int16
	}{
		"valid": {
			ni: sql.NullInt16{
				Int16: 1234,
				Valid: true,
			},
			expected: convert.Int16p(1234),
		},
		"invalid": {
			ni: sql.NullInt16{
				Valid: false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := ToInt16Ptr(tc.ni)

			require.Equal(t, tc.expected, converted)
		})
	}
}
