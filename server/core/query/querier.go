package query

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/querier"
)

type Query struct {
	Config  Config
	Log     logger.Logger
	Tx      *sqlx.Tx
	Prepare preparer.Query
}

type Config struct {
	Name string
}

var _ querier.Querier = &Query{}

func (q *Query) Init() error {

	q.Log.Debug("initialising query %s", q.Name())

	if q.Tx == nil {
		return errors.New("query Tx is nil, cannot initialise")
	}

	if q.Prepare == nil {
		return errors.New("query Prepare is nil, cannot initialise")
	}

	return nil
}

func (q *Query) Name() string {
	return q.Config.Name
}

func (q *Query) Exec(params map[string]interface{}) (sql.Result, error) {
	stmt := q.Prepare.Stmt()
	stmt = q.Tx.NamedStmt(stmt)

	res, err := stmt.Exec(params)
	if err != nil {
		q.Log.Warn("failed exec >%v<", err)
		return nil, err
	}

	return res, err
}

// SQL should be overridden in a custom implementation with the SQL statement specific to the query
func (q *Query) SQL() string {
	return ""
}
