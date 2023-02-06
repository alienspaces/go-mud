package modeller

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
)

// Modeller -
type Modeller interface {
	Init(pRepo preparer.Repository, pQ preparer.Query, tx *sqlx.Tx) (err error)
	Commit() error
	Rollback() error
}
