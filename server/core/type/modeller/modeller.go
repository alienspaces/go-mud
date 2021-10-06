package modeller

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
)

// Modeller -
type Modeller interface {
	Init(p preparer.Preparer, tx *sqlx.Tx) (err error)
	Commit() error
	Rollback() error
}
