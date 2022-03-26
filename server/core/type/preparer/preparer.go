package preparer

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/type/preparable"
)

type ExcludePreparation struct {
	GetOne          bool
	GetOneForUpdate bool
	GetMany         bool
	CreateOne       bool
	CreateMany      bool
	UpdateOne       bool
	UpdateMany      bool
	DeleteOne       bool
	DeleteMany      bool
	RemoveOne       bool
	RemoveMany      bool
}

type Repository interface {
	Init(tx *sqlx.DB) (err error)
	Prepare(m preparable.Repository, e ExcludePreparation) error
	GetOneStmt(m preparable.Repository) *sqlx.Stmt
	GetOneForUpdateStmt(m preparable.Repository) *sqlx.Stmt
	GetManyStmt(m preparable.Repository) *sqlx.NamedStmt
	CreateOneStmt(m preparable.Repository) *sqlx.NamedStmt
	UpdateOneStmt(m preparable.Repository) *sqlx.NamedStmt
	UpdateManyStmt(m preparable.Repository) *sqlx.NamedStmt
	DeleteOneStmt(m preparable.Repository) *sqlx.NamedStmt
	DeleteManyStmt(m preparable.Repository) *sqlx.NamedStmt
	RemoveOneStmt(m preparable.Repository) *sqlx.NamedStmt
	RemoveManyStmt(m preparable.Repository) *sqlx.NamedStmt
	GetOneSQL(m preparable.Repository) string
	GetManySQL(m preparable.Repository) string
	CreateSQL(m preparable.Repository) string
	UpdateOneSQL(m preparable.Repository) string
	UpdateManySQL(m preparable.Repository) string
	DeleteOneSQL(m preparable.Repository) string
	DeleteManySQL(m preparable.Repository) string
	RemoveOneSQL(m preparable.Repository) string
	RemoveManySQL(m preparable.Repository) string
}

type Query interface {
	Init(tx *sqlx.DB) (err error)
	Prepare(q preparable.Query) error

	Stmt() *sqlx.NamedStmt
	SQL() string
}
