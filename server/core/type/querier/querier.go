package querier

import "database/sql"

type Querier interface {
	Init() error
	Name() string

	Exec(params map[string]interface{}) (sql.Result, error)
}
