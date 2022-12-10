package querier

import "database/sql"

type Querier interface {
	Init() error
	Name() string
	SQL() string
	Exec(params map[string]interface{}) (sql.Result, error)
}
