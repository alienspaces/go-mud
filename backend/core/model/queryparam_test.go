package model

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
)

func TestResolveQueryParams(t *testing.T) {
	type args struct {
		queryParams map[string]any
	}
	tests := []struct {
		name         string
		args         args
		wantOpts     *coresql.Options
		wantPageSize int
	}{
		{
			name: "multiple values, pagination, sort_column",
			args: args{
				queryParams: map[string]any{
					"customer_countries": []string{"US", "CA", "AU"},
					"contract_country":   []string{"US"},
					"sort_column":        []string{"-created_at", "updated_at"},
					"page_size":          []string{"10"},
					"page_number":        []string{"1"},
				},
			},
			wantOpts: &coresql.Options{
				Params: []coresql.Param{
					{
						Col: "customer_countries",
						Val: []string{"US", "CA", "AU"},
					},
					{
						Col: "contract_country",
						Val: "US",
					},
				},
				OrderBy: []coresql.OrderBy{
					{
						Col:       "created_at",
						Direction: coresql.OrderDirectionDESC,
					},
					{
						Col:       "updated_at",
						Direction: coresql.OrderDirectionASC,
					},
				},
				Limit:  11,
				Offset: 0,
			},
			wantPageSize: 10,
		},
		{
			name: "single value, pagination, sort_column",
			args: args{
				queryParams: map[string]any{
					"customer_countries": []string{"US"},
					"contract_country":   []string{"US"},
					"sort_column":        []string{"created_at", "-updated_at"},
					"page_size":          []string{"10"},
					"page_number":        []string{"2"},
				},
			},
			wantOpts: &coresql.Options{
				Params: []coresql.Param{
					{
						Col: "customer_countries",
						Val: "US",
					},
					{
						Col: "contract_country",
						Val: "US",
					},
				},
				OrderBy: []coresql.OrderBy{
					{
						Col:       "created_at",
						Direction: coresql.OrderDirectionASC,
					},
					{
						Col:       "updated_at",
						Direction: coresql.OrderDirectionDESC,
					},
				},
				Limit:  11,
				Offset: 10,
			},
			wantPageSize: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOpts, actualPageSize, err := ResolveQueryParams(nil, tt.args.queryParams)
			// must sort params to reliably test, since query params is a map
			sort.Slice(actualOpts.Params, func(i, j int) bool {
				return actualOpts.Params[i].Col < actualOpts.Params[j].Col
			})
			sort.Slice(tt.wantOpts.Params, func(i, j int) bool {
				return tt.wantOpts.Params[i].Col < tt.wantOpts.Params[j].Col
			})
			require.NoError(t, err, "ResolveQueryParams should not return error")
			require.Equalf(t, tt.wantOpts, actualOpts, "coresql.Options should equal")
			require.Equalf(t, tt.wantPageSize, actualPageSize, "page size should equal")
		})
	}
}
