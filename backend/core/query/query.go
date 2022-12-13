package query

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/querier"
)

type Query struct {
	Config   Config
	Log      logger.Logger
	Tx       *sqlx.Tx
	Preparer preparer.Query
	Alias    string
}

type Config struct {
	Name string
}

var _ querier.Querier = &Query{}

func (q *Query) Name() string {
	return q.Config.Name
}

func (q *Query) Init() error {

	q.Log.Debug("initialising query %s", q.Name())

	if q.Tx == nil {
		return errors.New("query Tx is nil, cannot initialise")
	}

	if q.Preparer == nil {
		return errors.New("query Preparer is nil, cannot initialise")
	}

	return nil
}

func (q *Query) Exec(params map[string]interface{}) (sql.Result, error) {
	l := q.Log
	tx := q.Tx

	stmt := q.Preparer.Stmt(q)
	stmt = tx.NamedStmt(stmt)

	res, err := stmt.Exec(params)
	if err != nil {
		l.Warn("failed exec >%v<", err)
		return nil, err
	}

	return res, err
}

func (q *Query) GetRows(sql string, params map[string]interface{}, operators map[string]string) (*sqlx.Rows, error) {
	l := q.Log
	tx := q.Tx

	// params
	querySQL, queryParams, err := coresql.FromParamsAndOperators(q.Alias, sql, params, operators)
	if err != nil {
		l.Warn("failed generating query >%v<", err)
		return nil, err
	}

	l.Debug("SQL >%s< Params >%#v<", querySQL, queryParams)

	rows, err := tx.NamedQuery(querySQL, queryParams)
	if err != nil {
		l.Warn("failed querying rows >%v<", err)
		return nil, err
	}

	return rows, err
}

// SQL should be overridden in a custom implementation with the SQL statement specific to the query
func (q *Query) SQL() string {
	return ""
}
