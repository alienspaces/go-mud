package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sqlFromParamsAndOperator(t *testing.T) {

	type testCase struct {
		name         string
		initialSQL   string
		params       map[string]interface{}
		operators    map[string]string
		expectString string
		expectParams map[string]interface{}
		expectError  bool
	}

	r := Repository{
		Config: Config{TableName: "table1"},
	}

	for _, tc := range []testCase{
		{
			name:         "No params or operators",
			initialSQL:   r.GetManySQL(),
			params:       nil,
			operators:    nil,
			expectString: r.GetManySQL(),
			expectParams: map[string]interface{}{},
			expectError:  false,
		},
		{
			name:         "Params: status = \"created\" with Operators: Limit",
			initialSQL:   r.GetManySQL(),
			params:       map[string]interface{}{"status": "created"},
			operators:    map[string]string{OperatorLimit: "1"},
			expectString: r.GetManySQL() + "AND status = :status\nLIMIT 1\n",
			expectParams: map[string]interface{}{"status": "created"},
			expectError:  false,
		},
		{
			name:         "Params: status = \"created\" with Operators: Limit, OrderByDescending",
			initialSQL:   r.GetManySQL(),
			params:       map[string]interface{}{"status": "created"},
			operators:    map[string]string{OperatorLimit: "1", OperatorOrderByDescending: "updated_at"},
			expectString: r.GetManySQL() + "AND status = :status\nORDER BY updated_at DESC\nLIMIT 1\n",
			expectParams: map[string]interface{}{"status": "created"},
			expectError:  false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running test >%s<", tc.name)
			resultString, resultParams, err := r.sqlFromParamsAndOperator(tc.initialSQL, tc.params, tc.operators)
			require.Equal(t, tc.expectString, resultString, "Result string equals expected")
			require.Equal(t, tc.expectParams, resultParams, "Result params equals expected")
			require.Equal(t, tc.expectError, err != nil, "Result error equals expected")
		})
	}
}
