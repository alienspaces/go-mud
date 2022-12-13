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

	prepared map[string]bool
	stmtList map[string]*sqlx.NamedStmt
}

var _ preparer.Query = &Query{}

func NewQueryPreparer(l logger.Logger) (*Query, error) {

	pCfg := Query{
		Log:      l,
		prepared: make(map[string]bool),
		stmtList: make(map[string]*sqlx.NamedStmt),
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

	// This function is called on every new Modeller initialisation (i.e., on every new DB transaction).
	// To prevent memory leaks, we must protect against the same SQL statement being prepared multiple times.
	if _, ok := q.prepared[p.Name()]; ok {
		return nil
	}

	stmt, err := q.DB.PrepareNamed(p.SQL())
	if err != nil {
		q.Log.Warn("error preparing QuerySQL statement >%v<", err)
		return err
	}

	q.prepared[p.Name()] = true
	q.stmtList[p.Name()] = stmt

	return nil
}

func (q *Query) Stmt(p preparable.Query) *sqlx.NamedStmt {
	return q.stmtList[p.Name()]
}
