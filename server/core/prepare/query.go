package prepare

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparable"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
)

type Query struct {
	Log logger.Logger
	DB  *sqlx.DB

	stmt *sqlx.NamedStmt
	sql  string
}

var _ preparer.Query = &Query{}

func NewQueryPreparer(l logger.Logger) (*Query, error) {

	pCfg := Query{
		Log: l,
	}

	return &pCfg, nil
}

// Init - Initialise preparer with database tx
func (q *Query) Init(db *sqlx.DB) error {

	if db == nil {
		msg := "database db is nil, cannot init"
		q.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	q.DB = db

	return nil
}

func (q *Query) Prepare(p preparable.Query) error {
	sql := p.SQL()

	stmt, err := q.DB.PrepareNamed(sql)
	if err != nil {
		q.Log.Warn("error preparing QuerySQL statement >%v<", err)
		return err
	}

	q.sql = sql
	q.stmt = stmt

	return nil
}

func (q *Query) Stmt() *sqlx.NamedStmt {
	return q.stmt
}

func (q *Query) SQL() string {
	return q.sql
}
