package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/querier"
	"gitlab.com/alienspaces/go-mud/backend/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

// Model -
type Model struct {
	Config       configurer.Configurer
	Log          logger.Logger
	Store        storer.Storer
	Repositories map[string]repositor.Repositor
	Queries      map[string]querier.Querier
	Tx           *sqlx.Tx
	Err          error

	// composable functions
	RepositoriesFunc func(p preparer.Repository, tx *sqlx.Tx) ([]repositor.Repositor, error)
	QueriesFunc      func(p preparer.Query, tx *sqlx.Tx) ([]querier.Querier, error)
}

var _ modeller.Modeller = &Model{}

// NewModel - intended for testing only, maybe move into test files..
func NewModel(c configurer.Configurer, l logger.Logger, s storer.Storer) (m *Model, err error) {
	m = &Model{
		Config: c,
		Log:    l,
		Store:  s,
		Err:    nil,
	}

	return m, nil
}

// Init -
func (m *Model) Init(pRepo preparer.Repository, pQ preparer.Query, tx *sqlx.Tx) (err error) {

	// tx required
	if tx == nil {
		msg := "failed init, tx is required"
		m.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	if m.RepositoriesFunc == nil {
		m.RepositoriesFunc = m.DefaultRepositoriesFunc
	}

	if m.QueriesFunc == nil {
		m.QueriesFunc = m.DefaultQueriesFunc
	}

	m.Tx = tx

	// Repositories
	repositories, err := m.RepositoriesFunc(pRepo, tx)
	if err != nil {
		m.Log.Warn("failed repositories func >%v<", err)
		return err
	}

	m.Repositories = make(map[string]repositor.Repositor)
	for _, r := range repositories {
		m.Repositories[r.TableName()] = r
	}

	// Queries
	queries, err := m.QueriesFunc(pQ, tx)
	if err != nil {
		m.Log.Warn("failed queries func >%v<", err)
	}

	m.Queries = make(map[string]querier.Querier)
	for _, q := range queries {
		m.Queries[q.Name()] = q
	}

	return nil
}

// DefaultRepositoriesFunc - default repositor.RepositoriesFunc, override this function for custom repositories
func (m *Model) DefaultRepositoriesFunc(p preparer.Repository, tx *sqlx.Tx) ([]repositor.Repositor, error) {

	m.Log.Info("** repositor.Repositories **")

	return nil, nil
}

// DefaultQueriesFunc - default repositor.QueriesFunc, override this function for custom queries
func (m *Model) DefaultQueriesFunc(p preparer.Query, tx *sqlx.Tx) ([]querier.Querier, error) {

	m.Log.Info("** querier.Queriers **")

	return nil, nil
}

// Commit -
func (m *Model) Commit() error {
	if m.Tx == nil {
		msg := "cannot commit, database Tx is nil"
		m.Log.Warn(msg)
		return fmt.Errorf(msg)
	}
	return m.Tx.Commit()
}

// Rollback -
func (m *Model) Rollback() error {
	if m.Tx == nil {
		msg := "cannot rollback, database Tx is nil"
		m.Log.Warn(msg)
		return fmt.Errorf(msg)
	}
	return m.Tx.Rollback()
}
