package queryparam

import (
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
)

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

	for col, v := range qp.Params {
		var val any
		val = v
		if len(v) == 1 {
			val = v[0]
		}

		opts.Params = append(opts.Params, coresql.Param{
			Col: col,
			Val: val,
		})
	}

	return opts
}
