package null

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/backend/core/convert"
)

func TestNullInt16FromInt16(t *testing.T) {

	tests := map[string]struct {
		value         int16
		expectedValid bool
	}{
		"valid value": {
			value:         416,
			expectedValid: true,
		},
		"valid zero value": {
			value:         0,
			expectedValid: true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullInt16FromInt16(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.Int16)
		})
	}
}

func TestNullInt16FromInt16Ptr(t *testing.T) {

	tests := map[string]struct {
		value         *int16
		expectedValid bool
	}{
		"valid value": {
			value:         convert.Int16p(654),
			expectedValid: true,
		},
		"valid empty": {
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
			converted := NullInt16FromInt16Ptr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.Int16)
			}
		})
	}
}

func TestNullInt16ToInt16(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt16
		expected int16
	}{
		"valid": {
			ns: sql.NullInt16{
				Int16: 654,
				Valid: true,
			},
			expected: 654,
		},
		"invalid": {
			ns: sql.NullInt16{
				Int16: 654,
				Valid: false,
			},
			expected: 0,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullInt16ToInt16(tc.ns)
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}

func TestNullInt16ToInt16Ptr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt16
		expected *int16
	}{
		"valid": {
			ns: sql.NullInt16{
				Int16: 654,
				Valid: true,
			},
			expected: convert.Int16p(654),
		},
		"invalid": {
			ns: sql.NullInt16{
				Int16: 654,
				Valid: false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullInt16ToInt16Ptr(tc.ns)
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}
