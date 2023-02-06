package preparer

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/type/preparable"
)

type ExcludePreparation struct {
	GetOne          bool
	GetOneForUpdate bool
	CreateOne       bool
	UpdateOne       bool
	DeleteOne       bool
	RemoveOne       bool
}

type Repository interface {
	Init(tx *sqlx.DB) (err error)
	Prepare(m preparable.Repository, e ExcludePreparation) error

	GetOneStmt(m preparable.Repository) *sqlx.Stmt
	GetOneForUpdateStmt(m preparable.Repository) *sqlx.Stmt
	CreateOneStmt(m preparable.Repository) *sqlx.NamedStmt
	UpdateOneStmt(m preparable.Repository) *sqlx.NamedStmt
	DeleteOneStmt(m preparable.Repository) *sqlx.NamedStmt
	RemoveOneStmt(m preparable.Repository) *sqlx.NamedStmt

	GetOneSQL(m preparable.Repository) string
	CreateSQL(m preparable.Repository) string
	UpdateOneSQL(m preparable.Repository) string
	DeleteOneSQL(m preparable.Repository) string
	RemoveOneSQL(m preparable.Repository) string
}

type Query interface {
	Init(tx *sqlx.DB) (err error)
	Prepare(q preparable.Query) error
	Stmt(q preparable.Query) *sqlx.NamedStmt
}
