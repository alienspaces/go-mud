package null

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/backend/core/convert"
)

func TestNullInt64FromInt64(t *testing.T) {

	tests := map[string]struct {
		value         int64
		expectedValid bool
	}{
		"valid value": {
			value:         464,
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
			converted := NullInt64FromInt64(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.Int64)
		})
	}
}

func TestNullInt64FromInt64Ptr(t *testing.T) {

	tests := map[string]struct {
		value         *int64
		expectedValid bool
	}{
		"valid value": {
			value:         convert.Int64p(654),
			expectedValid: true,
		},
		"valid empty": {
			value:         convert.Int64p(0),
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
			converted := NullInt64FromInt64Ptr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.Int64)
			}
		})
	}
}

func TestNullInt64ToInt64(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt64
		expected int64
	}{
		"valid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: true,
			},
			expected: 654,
		},
		"invalid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: false,
			},
			expected: 0,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullInt64ToInt64(tc.ns)
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}

func TestNullInt64ToInt64Ptr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt64
		expected *int64
	}{
		"valid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: true,
			},
			expected: convert.Int64p(654),
		},
		"invalid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := NullInt64ToInt64Ptr(tc.ns)
			require.Equal(t, tc.expected, converted)
		})
	}
}
