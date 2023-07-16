package query

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/collection/set"
	"gitlab.com/alienspaces/go-mud/backend/core/convert"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/querier"
)

type Query struct {
	Config  Config
	Log     logger.Logger
	Tx      *sqlx.Tx
	Prepare preparer.Query
}

type Config struct {
	Name        string
	ArrayFields set.Set[string]
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

	if q.Config.Name == "" {
		return errors.New("query Config Name is empty, cannot initialise")
	}

	if q.ArrayFields() == nil {
		return errors.New("repository ArrayFields is nil, cannot initialise")
	}

	return nil
}

func (q *Query) Name() string {
	return q.Config.Name
}

func (q *Query) ArrayFields() set.Set[string] {
	return q.Config.ArrayFields
}

func (q *Query) Exec(params map[string]interface{}) (sql.Result, error) {
	l := q.Log

	stmt := q.Prepare.Stmt(q)
	stmt = q.Tx.NamedStmt(stmt)

	res, err := stmt.Exec(params)
	if err != nil {
		l.Warn("failed exec >%v<", err)
		return nil, err
	}

	return res, err
}

func (q *Query) GetRows(opts *coresql.Options) (*sqlx.Rows, error) {
	l := q.Log
	tx := q.Tx

	// params
	querySQL := q.Prepare.SQL(q)

	opts, err := q.resolveOpts(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opts: sql >%s< opts >%#v< >%v<", querySQL, opts, err)
	}

	querySQL, queryArgs, err := coresql.From(querySQL, opts)
	if err != nil {
		q.Log.Warn("failed generating query >%v<", err)
		return nil, err
	}

	l.Debug("Resulting SQL >%s< Params >%#v<", querySQL, queryArgs)

	rows, err := tx.NamedQuery(querySQL, queryArgs)
	if err != nil {
		l.Warn("Failed querying row >%v<", err)
		return nil, err
	}

	return rows, err
}

func (q *Query) resolveOpts(opts *coresql.Options) (*coresql.Options, error) {
	if opts == nil {
		return opts, nil
	}

	for i, p := range opts.Params {
		if p.Op != "" {
			// if Op is specified, it is assumed you know what you're doing
			continue
		}

		switch t := p.Val.(type) {
		case []string:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case []int:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		}

		isArrayField := q.ArrayFields().Contains(p.Col)
		if isArrayField {
			if len(p.Array) > 0 {
				p.Op = coresql.OpContains
			} else {
				p.Op = coresql.OpAny
			}
		} else {
			if len(p.Array) > 0 {
				p.Op = coresql.OpIn
			} else {
				p.Op = coresql.OpEqual
			}
		}

		opts.Params[i] = p
	}

	return opts, nil
}

// SQL should be overridden in implementation with the SQL statement specific to the query
func (q *Query) SQL() string {
	return ""
}
