package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sqlFromParamsAndOperators(t *testing.T) {

	type testCase struct {
		name         string
		initialSQL   string
		params       map[string]interface{}
		operators    map[string]string
		expectString string
		expectParams map[string]interface{}
		expectError  bool
	}

	initialSQL := "SELECT * FROM test WHERE 1 = 1 "

	for _, tc := range []testCase{
		{
			name:         "No params or operators",
			initialSQL:   initialSQL,
			params:       nil,
			operators:    nil,
			expectString: initialSQL,
			expectParams: map[string]interface{}{},
			expectError:  false,
		},
		{
			name:         "Params: status = \"created\" with Operators: Limit",
			initialSQL:   initialSQL,
			params:       map[string]interface{}{"status": "created"},
			operators:    map[string]string{OperatorLimit: "1"},
			expectString: initialSQL + "AND status = :status\nLIMIT 1\n",
			expectParams: map[string]interface{}{"status": "created"},
			expectError:  false,
		},
		{
			name:         "Params: status = \"created\" with Operators: Limit, OrderByDescending",
			initialSQL:   initialSQL,
			params:       map[string]interface{}{"status": "created"},
			operators:    map[string]string{OperatorLimit: "1", OperatorOrderByDescending: "updated_at"},
			expectString: initialSQL + "AND status = :status\nORDER BY updated_at DESC\nLIMIT 1\n",
			expectParams: map[string]interface{}{"status": "created"},
			expectError:  false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running test >%s<", tc.name)
			resultString, resultParams, err := FromParamsAndOperators(tc.initialSQL, tc.params, tc.operators)
			require.Equal(t, tc.expectString, resultString, "Result string equals expected")
			require.Equal(t, tc.expectParams, resultParams, "Result params equals expected")
			require.Equal(t, tc.expectError, err != nil, "Result error equals expected")
		})
	}
}
