package queryparam

import (
	"fmt"
	"net/url"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildQueryParams(t *testing.T) {
	type args struct {
		queryParams url.Values
	}
	tests := []struct {
		name string
		args args
		want *QueryParams
	}{
		{
			name: "multiple values, pagination, sort_column",
			args: args{
				queryParams: url.Values{
					"customer_countries": []string{"US", "CA", "AU"},
					"contract_country":   []string{"US"},
					fmt.Sprintf("%s:%s", "created_at", OpGreaterThanEqual): []string{"2020-01-01"},
					fmt.Sprintf("%s:%s", "created_at", OpLessThanEqual):    []string{"2022-01-01"},
					"sort_column": []string{"-created_at", "updated_at"},
					"page_size":   []string{"10"},
					"page_number": []string{"1"},
				},
			},
			want: &QueryParams{
				Params: map[string][]QueryParam{
					"customer_countries": {
						{
							Val: "US",
						},
						{
							Val: "CA",
						},
						{
							Val: "AU",
						},
					},
					"contract_country": {
						{
							Val: "US",
						},
					},
					"created_at": {
						{
							Val: "2020-01-01",
							Op:  OpGreaterThanEqual,
						},
						{
							Val: "2022-01-01",
							Op:  OpLessThanEqual,
						},
					},
				},
				SortColumns: []SortColumn{
					{
						Col:          "created_at",
						IsDescending: true,
					},
					{
						Col:          "updated_at",
						IsDescending: false,
					},
				},
				PageSize:   10,
				PageNumber: 1,
			},
		},
		{
			name: "single value, pagination, sort_column",
			args: args{
				queryParams: url.Values{
					"customer_countries": []string{"US"},
					"contract_country":   []string{"US"},
					"sort_column":        []string{"created_at", "-updated_at"},
					"page_size":          []string{"10"},
					"page_number":        []string{"2"},
				},
			},
			want: &QueryParams{
				Params: map[string][]QueryParam{
					"customer_countries": {
						{
							Val: "US",
						},
					},
					"contract_country": {
						{
							Val: "US",
						},
					},
				},
				SortColumns: []SortColumn{
					{
						Col:          "created_at",
						IsDescending: false,
					},
					{
						Col:          "updated_at",
						IsDescending: true,
					},
				},
				PageSize:   10,
				PageNumber: 2,
			},
		},
		{
			name: "no query params",
			args: args{
				queryParams: url.Values{},
			},
			want: &QueryParams{
				Params: map[string][]QueryParam{},
				SortColumns: []SortColumn{
					{
						Col:          DefaultOrderDescendingColumn,
						IsDescending: true,
					},
				},
				PageSize:   DefaultPageSizeInt,
				PageNumber: DefaultPageNumberInt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualQueryParams, err := BuildQueryParams(nil, tt.args.queryParams)
			require.NoError(t, err, "buildQueryParams should not return error")

			// url.Values (map[string][]string) processing order is not stable
			for k, expected := range tt.want.Params {
				sort.Slice(expected, func(i, j int) bool {
					if expected[i].Val == expected[j].Val {
						return expected[i].Op < expected[j].Op
					}
					return expected[i].Val < expected[j].Val
				})

				actual := actualQueryParams.Params[k]
				sort.Slice(actual, func(i, j int) bool {
					if actual[i].Val == actual[j].Val {
						return actual[i].Op < actual[j].Op
					}
					return actual[i].Val < actual[j].Val
				})

				require.Equalf(t, expected, actual, "expected should equal")
			}

			require.Equalf(t, tt.want.PageSize, actualQueryParams.PageSize, "page size should equal")
			require.Equalf(t, tt.want.SortColumns, actualQueryParams.SortColumns, "sort columns should equal")
			require.Equalf(t, tt.want.PageNumber, actualQueryParams.PageNumber, "page number should equal")
		})
	}
}
