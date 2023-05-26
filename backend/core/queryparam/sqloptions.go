package queryparam

import (
	"time"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
)

const dateOnly = "2006-01-02"

func ToSQLOptions(qp *QueryParams) *coresql.Options {
	opts := &coresql.Options{
		Limit:  qp.PageSize + 1,
		Offset: (qp.PageNumber - 1) * qp.PageSize,
	}

	for _, sc := range qp.SortColumns {
		dir := coresql.OrderDirectionASC
		if sc.IsDescending {
			dir = coresql.OrderDirectionDESC
		}

		opts.OrderBy = append(opts.OrderBy, coresql.OrderBy{
			Col:       sc.Col,
			Direction: dir,
		})
	}

	for col, values := range qp.Params {
		qOpVals := map[Operator][]string{}
		var sqlParams []coresql.Param

		var sqlOp coresql.Operator
		for _, p := range values {
			switch p.Op {
			case "":
				// "" corresponds to =, IN, ANY, @>
				qOpVals[p.Op] = append(qOpVals[p.Op], p.Val)
				continue
			case OpGreaterThan:
				sqlOp = coresql.OpGreaterThan
			case OpGreaterThanEqual:
				sqlOp = coresql.OpGreaterThanEqual
			case OpLessThan:
				sqlOp = coresql.OpLessThan
			case OpLessThanEqual:
				sqlOp = coresql.OpLessThanEqual

				// TODO Replace with time.DateOnly when using Go1.20: https://github.com/golang/go/issues/52746
				if _, err := time.Parse(dateOnly, p.Val); err == nil {
					// This is necessary because postgres maps a date of 2022-01-01 to 2022-1-01T00:00:00Z
					p.Val += "T23:59:59.999999999Z"
				}
			case OpNotEqual:
				sqlOp = coresql.OpNotEqual
			case OpLike:
				sqlOp = coresql.OpLike
				p.Val = "%" + p.Val + "%"
			case OpILike:
				sqlOp = coresql.OpILike
				p.Val = "%" + p.Val + "%"
			}

			sqlParams = append(sqlParams, coresql.Param{
				Col: col,
				Op:  sqlOp,
				Val: p.Val,
			})
		}

		for op, v := range qOpVals {
			if op == "" {
				var val any
				val = v
				if len(v) == 1 {
					val = v[0]
				}

				sqlParams = append(sqlParams, coresql.Param{
					Col: col,
					Val: val,
				})
			}
		}
		opts.Params = append(opts.Params, sqlParams...)
	}

	return opts
}
