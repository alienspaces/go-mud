package sql

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/convert"
)

func Test_From(t *testing.T) {

	type testCase struct {
		name         string
		initialSQL   string
		opts         *Options
		expectString string
		expectParams map[string]interface{}
	}

	initialSQL := "SELECT * FROM test WHERE 1 = 1 "

	for _, tc := range []testCase{
		{
			name:         "No params or operators",
			initialSQL:   initialSQL,
			opts:         nil,
			expectString: initialSQL,
			expectParams: nil,
		},
		{
			name:       "string array fields: [multiple values in slice, one value in slice, single value], string field: [single value, multiple values in slice, single value in slice], with Limit, Offset",
			initialSQL: initialSQL,
			opts: &Options{
				Params: []Param{
					{
						Col:   "string_array_field_slice_multiple",
						Op:    OpContains,
						Array: convert.GenericSlice([]string{"a", "b"}),
					},
					{
						Col:   "string_array_field_slice_single",
						Op:    OpContains,
						Array: convert.GenericSlice([]string{"c"}),
					},
					{
						Col: "string_array_field_single",
						Op:  OpAny,
						Val: "d",
					},
					{
						Col: "string_field_single",
						Op:  OpEqualTo,
						Val: "e",
					},
					{
						Col:   "string_field_slice_multiple",
						Op:    OpIn,
						Array: convert.GenericSlice([]string{"f", "g"}),
					},
					{
						Col:   "string_field_slice_single",
						Op:    OpIn,
						Array: convert.GenericSlice([]string{"h"}),
					},
				},
				Limit:  1,
				Offset: 1,
			},
			expectString: initialSQL +
				`AND string_array_field_slice_multiple @> array[:string_array_field_slice_multiple0,:string_array_field_slice_multiple1]
AND string_array_field_slice_single @> array[:string_array_field_slice_single0]
AND :string_array_field_single = ANY(string_array_field_single)
AND string_field_single = :string_field_single
AND string_field_slice_multiple IN (:string_field_slice_multiple0,:string_field_slice_multiple1)
AND string_field_slice_single IN (:string_field_slice_single0)
LIMIT 1
OFFSET 1
`,
			expectParams: map[string]any{
				"string_array_field_slice_multiple0": "a",
				"string_array_field_slice_multiple1": "b",
				"string_array_field_slice_single0":   "c",
				"string_array_field_single":          "d",
				"string_field_single":                "e",
				"string_field_slice_multiple0":       "f",
				"string_field_slice_multiple1":       "g",
				"string_field_slice_single0":         "h",
			},
		},
		{
			name:       "int array fields: [multiple values in slice, one value in slice, single value], int field: [single value, multiple values in slice, single value in slice], with Limit, Offset",
			initialSQL: initialSQL,
			opts: &Options{
				Params: []Param{
					{
						Col:   "int_array_field_slice_multiple",
						Op:    OpContains,
						Array: convert.GenericSlice([]int{1, 2}),
					},
					{
						Col:   "int_array_field_slice_single",
						Op:    OpContains,
						Array: convert.GenericSlice([]int{3}),
					},
					{
						Col: "int_array_field_single",
						Op:  OpAny,
						Val: 4,
					},
					{
						Col: "int_field_single",
						Op:  OpEqualTo,
						Val: 5,
					},
					{
						Col:   "int_field_slice_multiple",
						Op:    OpIn,
						Array: convert.GenericSlice([]int{6, 7}),
					},
					{
						Col:   "int_field_slice_single",
						Op:    OpIn,
						Array: convert.GenericSlice([]int{8}),
					},
				},
				Limit:  1,
				Offset: 1,
			},
			expectString: initialSQL +
				`AND int_array_field_slice_multiple @> array[:int_array_field_slice_multiple0,:int_array_field_slice_multiple1]
AND int_array_field_slice_single @> array[:int_array_field_slice_single0]
AND :int_array_field_single = ANY(int_array_field_single)
AND int_field_single = :int_field_single
AND int_field_slice_multiple IN (:int_field_slice_multiple0,:int_field_slice_multiple1)
AND int_field_slice_single IN (:int_field_slice_single0)
LIMIT 1
OFFSET 1
`,
			expectParams: map[string]any{
				"int_array_field_slice_multiple0": 1,
				"int_array_field_slice_multiple1": 2,
				"int_array_field_slice_single0":   3,
				"int_array_field_single":          4,
				"int_field_single":                5,
				"int_field_slice_multiple0":       6,
				"int_field_slice_multiple1":       7,
				"int_field_slice_single0":         8,
			},
		},
		{
			name:       "status = \"created\" with Limit, OrderByDescending, OrderByAscending, For Update Skip Locked",
			initialSQL: initialSQL,
			opts: &Options{
				Params: []Param{
					{
						Col: "status",
						Op:  OpEqualTo,
						Val: "created",
					},
				},
				Limit: 1,
				OrderBy: []OrderBy{
					{
						Col:       "updated_at",
						Direction: OrderDirectionDESC,
					},
					{
						Col:       "created_at",
						Direction: OrderDirectionASC,
					},
				},
				Lock: ForUpdateSkipLocked,
			},
			expectString: initialSQL + "AND status = :status\nORDER BY\nupdated_at DESC,\ncreated_at ASC\nLIMIT 1\nFOR UPDATE SKIP LOCKED\n",
			expectParams: map[string]any{"status": "created"},
		},
		{
			name:       "\"created\" = ANY(array_field), For Update NOWAIT",
			initialSQL: initialSQL,
			opts: &Options{
				Params: []Param{
					{
						Col: "array_field",
						Op:  OpAny,
						Val: "created",
					},
				},
				Lock: ForUpdateNoWait,
			},
			expectString: initialSQL + "AND :array_field = ANY(array_field)\nFOR UPDATE NOWAIT\n",
			expectParams: map[string]any{"array_field": "created"},
		},
		{
			name:       "number = 1, For Update",
			initialSQL: initialSQL,
			opts: &Options{
				Params: []Param{
					{
						Col: "number",
						Op:  OpEqualTo,
						Val: 1,
					},
				},
				Lock: ForUpdate,
			},
			expectString: initialSQL + "AND number = :number\nFOR UPDATE \n",
			expectParams: map[string]any{"number": 1},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running test >%s<", tc.name)
			resultString, resultParams, err := From(tc.initialSQL, tc.opts)
			require.Equal(t, tc.expectString, resultString, "Result string equals expected")
			require.Equal(t, tc.expectParams, resultParams, "Result params equals expected")
			require.NoError(t, err, "Result error equals expected")
		})
	}
}
