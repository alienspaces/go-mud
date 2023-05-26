package queryparam

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
)

func TestToSQLOptions(t *testing.T) {
	type args struct {
		qp *QueryParams
	}
	tests := []struct {
		name           string
		args           args
		wantOpts       *coresql.Options
		isBetweenQuery bool
	}{
		{
			name: "like, ilike, multiple values, pagination, sort_column",
			args: args{
				qp: &QueryParams{
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
						"client_legal_name": {
							{
								Val: "ABC",
								Op:  OpLike,
							},
						},
						"name": {
							{
								Val: "abc",
								Op:  OpILike,
							},
						},
						"rabbits": {
							{
								Val: "10",
								Op:  OpLessThan,
							},
						},
						"tortoises": {
							{
								Val: "5",
								Op:  OpGreaterThan,
							},
						},
						"bananas": {
							{
								Val: "5",
								Op:  OpNotEqual,
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
					{
						Col: "client_legal_name",
						Op:  coresql.OpLike,
						Val: "%ABC%",
					},
					{
						Col: "name",
						Op:  coresql.OpILike,
						Val: "%abc%",
					},
					{
						Col: "rabbits",
						Op:  coresql.OpLessThan,
						Val: "10",
					},
					{
						Col: "tortoises",
						Op:  coresql.OpGreaterThan,
						Val: "5",
					},
					{
						Col: "bananas",
						Op:  coresql.OpNotEqual,
						Val: "5",
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
		{
			name:           "created_at:lte, updated_at:lte, updated_at:gte, created_at:gte",
			isBetweenQuery: true,
			args: args{
				qp: &QueryParams{
					Params: map[string][]QueryParam{
						"created_at": {
							{
								Val: "2018-01-01T23:59:59.99999Z",
								Op:  OpGreaterThanEqual,
							},
							{
								Val: "2019-12-12",
								Op:  OpLessThanEqual,
							},
						},
						"updated_at": {
							{
								Val: "2020-01-01",
								Op:  OpGreaterThanEqual,
							},
							{
								Val: "2020-05-31T23:59:59Z",
								Op:  OpLessThanEqual,
							},
						},
					},
					PageSize:   10,
					PageNumber: 1,
				},
			},
			wantOpts: &coresql.Options{
				Params: []coresql.Param{
					{
						Col: "created_at",
						Op:  coresql.OpGreaterThanEqual,
						Val: "2018-01-01T23:59:59.99999Z",
					},
					{
						Col: "created_at",
						Op:  coresql.OpLessThanEqual,
						Val: "2019-12-12T23:59:59.999999999Z",
					},
					{
						Col: "updated_at",
						Op:  coresql.OpGreaterThanEqual,
						Val: "2020-01-01",
					},
					{
						Col: "updated_at",
						Op:  coresql.OpLessThanEqual,
						Val: "2020-05-31T23:59:59Z",
					},
				},
				Limit:  11,
				Offset: 0,
			},
		},
		{
			name: "updated_at:gte, created_at:gte",
			args: args{
				qp: &QueryParams{
					Params: map[string][]QueryParam{
						"created_at": {
							{
								Val: "2018-01-01T23:59:59.99999Z",
								Op:  OpGreaterThanEqual,
							},
						},
						"updated_at": {
							{
								Val: "2020-01-01",
								Op:  OpGreaterThanEqual,
							},
						},
					},
					PageSize:   10,
					PageNumber: 1,
				},
			},
			wantOpts: &coresql.Options{
				Params: []coresql.Param{
					{
						Col: "created_at",
						Op:  coresql.OpGreaterThanEqual,
						Val: "2018-01-01T23:59:59.99999Z",
					},
					{
						Col: "updated_at",
						Op:  coresql.OpGreaterThanEqual,
						Val: "2020-01-01",
					},
				},
				Limit:  11,
				Offset: 0,
			},
		},
		{
			name: "created_at:lte, updated_at:lte",
			args: args{
				qp: &QueryParams{
					Params: map[string][]QueryParam{
						"created_at": {
							{
								Val: "2019-12-12",
								Op:  OpLessThanEqual,
							},
						},
						"updated_at": {
							{
								Val: "2020-05-31T23:59:59Z",
								Op:  OpLessThanEqual,
							},
						},
					},
					PageSize:   10,
					PageNumber: 1,
				},
			},
			wantOpts: &coresql.Options{
				Params: []coresql.Param{
					{
						Col: "created_at",
						Op:  coresql.OpLessThanEqual,
						Val: "2019-12-12T23:59:59.999999999Z",
					},
					{
						Col: "updated_at",
						Op:  coresql.OpLessThanEqual,
						Val: "2020-05-31T23:59:59Z",
					},
				},
				Limit:  11,
				Offset: 0,
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

			if tt.isBetweenQuery {
				for i, param := range actualOpts.Params {
					require.NotEqualf(t, param.Val, param.ValB, "coresql.Options.Params.Val should not equal ValB for >%s<", param.Col)
					paramValEqualsOneOfExpectedValOrValB := param.Val == tt.wantOpts.Params[i].Val || param.Val == tt.wantOpts.Params[i].ValB
					require.Truef(t, paramValEqualsOneOfExpectedValOrValB, "coresql.Options.Params.Val should equal one of expected Val or ValB for >%s< Val >%v< ValB >%v<", param.Col, param.Val, param.ValB)

					paramValBEqualsOneOfExpectedValOrValB := param.ValB == tt.wantOpts.Params[i].Val || param.ValB == tt.wantOpts.Params[i].ValB
					require.Truef(t, paramValBEqualsOneOfExpectedValOrValB, "coresql.Options.Params.ValB should equal one of expected Val or ValB expected values for >%s< Val >%v< ValB >%v<", param.Col, param.Val, param.ValB)
				}
				require.Equalf(t, tt.wantOpts.OrderBy, actualOpts.OrderBy, "coresql.Options.OrderBy should equal")
				require.Equalf(t, tt.wantOpts.Offset, actualOpts.Offset, "coresql.Options.Offset should equal")
				require.Equalf(t, tt.wantOpts.Limit, actualOpts.Limit, "coresql.Options.Limit should equal")
			} else {
				require.Equalf(t, tt.wantOpts, actualOpts, "coresql.Options should equal")
			}
		})
	}
}
