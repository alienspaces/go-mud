package nulltime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFromTime(t *testing.T) {

	tests := map[string]struct {
		time          time.Time
		expectedValid bool
	}{
		"valid time": {
			time:          time.Now(),
			expectedValid: true,
		},
		"invalid time": {
			time:          time.Time{},
			expectedValid: false,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := FromTime(tc.time)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.time, converted.Time)
		})
	}
}
