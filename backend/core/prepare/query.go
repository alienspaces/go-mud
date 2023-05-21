package prepare

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparable"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
)

type Query struct {
	Log logger.Logger
	DB  *sqlx.DB

	stmt map[string]*sqlx.NamedStmt
	sql  map[string]string
}

var _ preparer.Query = &Query{}

func NewQueryPreparer(l logger.Logger) (*Query, error) {

	pCfg := Query{
		Log: l,

		stmt: make(map[string]*sqlx.NamedStmt),
		sql:  make(map[string]string),
	}

	return &pCfg, nil
}

// Init - Initialise preparer with database tx
func (p *Query) Init(db *sqlx.DB) error {
	l := p.Log.WithFunctionContext("Init")

	if db == nil {
		msg := "database db is nil, cannot init"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	p.DB = db

	return nil
}

func (p *Query) Prepare(q preparable.Query) error {
	l := p.Log.WithFunctionContext("Prepare")

	name := q.Name()
	sql := q.SQL()

	// This function is called on every new Modeller initialisation (i.e., on every new DB transaction).
	// To prevent memory leaks, we must protect against the same SQL statement being prepared multiple times.
	if _, ok := p.stmt[name]; ok {
		return nil
	}

	stmt, err := p.DB.PrepareNamed(sql)
	if err != nil {
		l.Warn("error preparing QuerySQL statement >%v<", err)
		return err
	}

	p.sql[name] = sql
	p.stmt[name] = stmt

	return nil
}

func (p *Query) Stmt(q preparable.Query) *sqlx.NamedStmt {
	return p.stmt[q.Name()]
}

func (p *Query) SQL(q preparable.Query) string {
	return p.sql[q.Name()]
}
