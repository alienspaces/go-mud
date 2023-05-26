package querier

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
)

type Querier interface {
	Init() error
	Name() string
	GetRows(opts *coresql.Options) (*sqlx.Rows, error)
	Exec(params map[string]interface{}) (sql.Result, error)
}
