package queryparam

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
)

func TestQueryParamsToSQLOptions(t *testing.T) {
	type args struct {
		qp *QueryParams
	}
	tests := []struct {
		name     string
		args     args
		wantOpts *coresql.Options
	}{
		{
			name: "multiple values, pagination, sort_column",
			args: args{
				qp: &QueryParams{
					Params: map[string][]string{
						"customer_countries": {"US", "CA", "AU"},
						"contract_country":   {"US"},
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
		},
		{
			name: "single value, pagination, sort_column",
			args: args{
				qp: &QueryParams{
					Params: map[string][]string{
						"customer_countries": {"US"},
						"contract_country":   {"US"},
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOpts := ToSQLOptions(tt.args.qp)
			// must sort params to reliably test, since query params is a map
			sort.Slice(actualOpts.Params, func(i, j int) bool {
				return actualOpts.Params[i].Col < actualOpts.Params[j].Col
			})
			sort.Slice(tt.wantOpts.Params, func(i, j int) bool {
				return tt.wantOpts.Params[i].Col < tt.wantOpts.Params[j].Col
			})
			require.Equalf(t, tt.wantOpts, actualOpts, "coresql.Options should equal")
		})
	}
}
