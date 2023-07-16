package modeller

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
)

// Modeller -
type Modeller interface {
	Init(repoPreparer preparer.Repository, queryPreparer preparer.Query, tx *sqlx.Tx) (err error)
	SetTxLockTimeout(timeoutSecs float64) error
	Commit() error
	Rollback() error
}
