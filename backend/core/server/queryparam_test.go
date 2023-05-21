package server

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
)

func TestBuildQueryParams(t *testing.T) {
	type args struct {
		queryParams url.Values
	}
	tests := []struct {
		name string
		args args
		want *queryparam.QueryParams
	}{
		{
			name: "multiple values, pagination, sort_column",
			args: args{
				queryParams: url.Values{
					"customer_countries": []string{"US", "CA", "AU"},
					"contract_country":   []string{"US"},
					"sort_column":        []string{"-created_at", "updated_at"},
					"page_size":          []string{"10"},
					"page_number":        []string{"1"},
				},
			},
			want: &queryparam.QueryParams{
				Params: map[string][]string{
					"customer_countries": {"US", "CA", "AU"},
					"contract_country":   {"US"},
				},
				SortColumns: []queryparam.SortColumn{
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
			want: &queryparam.QueryParams{
				Params: map[string][]string{
					"customer_countries": {"US"},
					"contract_country":   {"US"},
				},
				SortColumns: []queryparam.SortColumn{
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
			want: &queryparam.QueryParams{
				Params: map[string][]string{},
				SortColumns: []queryparam.SortColumn{
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
			actualQueryParams, err := buildQueryParams(nil, tt.args.queryParams)
			require.NoError(t, err, "buildQueryParams should not return error")
			require.Equalf(t, tt.want, actualQueryParams, " should equal")
		})
	}
}
