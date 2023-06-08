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
		m.RepositoriesFunc = m.NewRepositories
	}

	if m.QueriesFunc == nil {
		m.QueriesFunc = m.NewQueriers
	}

	m.Tx = tx

	// repositories
	repositories, err := m.RepositoriesFunc(pRepo, tx)
	if err != nil {
		m.Log.Warn("failed repositories func >%v<", err)
		return err
	}

	m.Repositories = make(map[string]repositor.Repositor)
	for _, r := range repositories {
		m.Repositories[r.TableName()] = r
	}

	queriers, err := m.QueriesFunc(pQ, tx)
	if err != nil {
		m.Log.Warn("failed queries func >%v<", err)
		return err
	}

	m.Queries = make(map[string]querier.Querier)
	for _, q := range queriers {
		m.Queries[q.Name()] = q
	}

	return nil
}

// NewRepositories - default repositor.RepositoriesFunc, override this function for custom repositories
func (m *Model) NewRepositories(p preparer.Repository, tx *sqlx.Tx) ([]repositor.Repositor, error) {

	m.Log.Info("** repositor.Repositories **")

	return nil, nil
}

// NewQueriers - default repositor.QueriesFunc, override this function for custom queriers
func (m *Model) NewQueriers(p preparer.Query, tx *sqlx.Tx) ([]querier.Querier, error) {

	m.Log.Info("** querier.Queriers **")

	return nil, nil
}

// Commit -
func (m *Model) Commit() error {
	if m.Tx != nil {
		m.Tx.Commit()
		return nil
	}
	msg := "cannot commit, database Tx is nil"
	m.Log.Warn(msg)
	return fmt.Errorf(msg)
}

// SetTxLockTimeout -
func (m *Model) SetTxLockTimeout(timeoutSecs float64) error {
	if m.Tx == nil {
		err := fmt.Errorf("cannot set transaction lock timeout seconds, database Tx is nil")
		m.Log.Warn(err.Error())
		return err
	}

	// If we SET, instead of SET LOCAL, lock_timeout would be at the session-level.
	// Since we use connection pooling, this would mean that different sessions (and therefore requests)
	// would have different, unknown lock_timeout parameters.

	timeoutMs := timeoutSecs * 1000
	_, err := m.Tx.Exec(fmt.Sprintf("SET LOCAL lock_timeout = %d", int(timeoutMs)))
	if err != nil {
		m.Log.Warn("failed setting transaction lock timeout seconds >%v<", err)
		return err
	}

	m.Log.Debug("lock timeout seconds set to >%fs<", timeoutSecs)

	return nil
}

// Rollback -
func (m *Model) Rollback() error {
	if m.Tx != nil {
		return m.Tx.Rollback()
	}
	msg := "cannot rollback, database Tx is nil"
	m.Log.Warn(msg)
	return fmt.Errorf(msg)
}
