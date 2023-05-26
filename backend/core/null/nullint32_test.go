package null

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/backend/core/convert"
)

func TestNullInt32FromInt32(t *testing.T) {

	tests := map[string]struct {
		value         int32
		expectedValid bool
	}{
		"valid value": {
			value:         432,
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
			converted := NullInt32FromInt32(tc.value)
			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.Int32)
		})
	}
}

func TestNullInt32FromInt32Ptr(t *testing.T) {

	tests := map[string]struct {
		value         *int32
		expectedValid bool
	}{
		"valid value": {
			value:         convert.Int32p(654),
			expectedValid: true,
		},
		"valid empty": {
			value:         convert.Int32p(0),
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
			converted := NullInt32FromInt32Ptr(tc.value)
			require.Equal(t, tc.expectedValid, converted.Valid)
			if tc.value != nil {
				require.Equal(t, *tc.value, converted.Int32)
			}
		})
	}
}

func TestNullInt32ToInt32(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt32
		expected int32
	}{
		"valid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: true,
			},
			expected: 654,
		},
		"invalid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: false,
			},
			expected: 0,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullInt32ToInt32(tc.ns)
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}

func TestNullInt32ToInt32Ptr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt32
		expected *int32
	}{
		"valid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: true,
			},
			expected: convert.Int32p(654),
		},
		"invalid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullInt32ToInt32Ptr(tc.ns)
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}
