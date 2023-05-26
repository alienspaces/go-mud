package harness

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/prepare"
	"gitlab.com/alienspaces/go-mud/backend/core/store"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

// CreateDataFunc - callback function that creates test data
type CreateDataFunc func() error

// RemoveDataFunc - callback function that removes test data
type RemoveDataFunc func() error

// Testing -
type Testing struct {
	Config             configurer.Configurer
	Log                logger.Logger
	Store              storer.Storer
	RepositoryPreparer preparer.Repository
	QueryPreparer      preparer.Query
	Model              modeller.Modeller

	// ShouldCommitData is used to determine whether Setup and Teardown should commit data to the DB.
	// This should only be true if changes in one transaction must be visible in another (e.g., handler tests).
	ShouldCommitData bool

	// Modeller function
	ModellerFunc func() (modeller.Modeller, error)

	// Composable functions
	CreateDataFunc CreateDataFunc
	RemoveDataFunc RemoveDataFunc

	tx *sqlx.Tx
}

// NewTesting -
func NewTesting() (t *Testing, err error) {

	t = &Testing{}

	return t, nil
}

// Init -
func (t *Testing) Init() (err error) {

	// configurer
	if t.Config == nil {
		t.Config, err = config.NewConfigWithDefaults(nil, false)
		if err != nil {
			return err
		}
	}

	// logger
	if t.Log == nil {
		t.Log, err = log.NewLogger(t.Config)
		if err != nil {
			return err
		}
	}
	t.Log = t.Log.WithApplicationContext("harness")

	// storer
	if t.Store == nil {
		t.Store, err = store.NewStore(t.Config, t.Log)
		if err != nil {
			return err
		}
	}

	// preparer
	t.RepositoryPreparer, err = prepare.NewRepositoryPreparer(t.Log)
	if err != nil {
		t.Log.Warn("failed new preparer >%v<", err)
		return err
	}

	t.QueryPreparer, err = prepare.NewQueryPreparer(t.Log)
	if err != nil {
		t.Log.Warn("failed new preparer config >%v<", err)
		return err
	}

	db, err := t.Store.GetDb()
	if err != nil {
		t.Log.Warn("failed getting database handle >%v<", err)
		return err
	}

	err = t.RepositoryPreparer.Init(db)
	if err != nil {
		t.Log.Warn("failed preparer init >%v<", err)
		return err
	}

	err = t.QueryPreparer.Init(db)
	if err != nil {
		t.Log.Warn("failed preparer config init >%v<", err)
		return err
	}

	t.Log.Debug("Repository ready")

	// modeller
	t.Model, err = t.ModellerFunc()
	if err != nil {
		t.Log.Warn("failed new modeller >%v<", err)
		return err
	}

	t.Log.Debug("Modeller ready")

	return nil
}

// InitTx -
func (t *Testing) InitTx() (*sqlx.Tx, error) {

	if t.tx != nil {
		return t.tx, nil
	}

	t.Log.Debug("Starting database tx")

	tx, err := t.Store.GetTx()
	if err != nil {
		t.Log.Warn("failed getting database tx >%v<", err)
		return nil, err
	}

	t.tx = tx

	err = t.Model.Init(t.RepositoryPreparer, t.QueryPreparer, t.tx)
	if err != nil {
		t.Log.Warn("failed modeller init >%v<", err)
		return nil, err
	}

	return t.tx, nil
}

// CommitTx -
func (t *Testing) CommitTx() (err error) {

	err = t.tx.Commit()
	if err != nil {
		return err
	}
	t.tx = nil

	return nil
}

// RollbackTx -
func (t *Testing) RollbackTx() (err error) {

	err = t.tx.Rollback()
	if err != nil {
		return err
	}
	t.tx = nil

	return nil
}

// Setup -
//
// If ShouldCommitData is false, the tx is returned. The caller can perform other queries, but must commit or rollback the tx.
func (t *Testing) Setup() (*sqlx.Tx, error) {

	// init
	_, err := t.InitTx()
	if err != nil {
		t.Log.Warn("failed init >%v<", err)
		return nil, err
	}

	// data function is expected to create and manage its own store
	if t.CreateDataFunc != nil {
		t.Log.Debug("creating test data")
		err := t.CreateDataFunc()
		if err != nil {
			t.Log.Warn("failed creating data >%v<", err)
			return nil, err
		}
	}

	// commit data when configured, otherwise we are leaving
	// it up to tests to explicitly commit or rollback
	if t.ShouldCommitData {
		t.Log.Debug("committing database tx")
		err = t.CommitTx()
		if err != nil {
			t.Log.Warn("failed committing data >%v<", err)
			return nil, err
		}

		return nil, nil
	}

	return t.tx, nil
}

// Teardown -
func (t *Testing) Teardown() error {

	_, err := t.InitTx()
	if err != nil {
		t.Log.Warn("failed init >%v<", err)
		return err
	}

	// data function is expected to create and manage its own store
	if t.RemoveDataFunc != nil {
		t.Log.Debug("removing test data")
		err := t.RemoveDataFunc()
		if err != nil {
			t.Log.Warn("failed removing data >%v<", err)
			return err
		}
	}

	if t.ShouldCommitData {
		t.Log.Debug("committing database tx")
		err := t.CommitTx()
		if err != nil {
			t.Log.Warn("failed committing data >%v<", err)
			return err
		}
	} else {
		t.Log.Debug("rollback database tx")
		err := t.RollbackTx()
		if err != nil {
			t.Log.Warn("failed rolling back data >%v<", err)
			return err
		}
	}

	return nil
}

// Shutdown -
func (t *Testing) Shutdown(test *testing.T) {
	db, err := t.Store.GetDb()
	require.NoError(test, err, "getDb should return no error")

	err = db.Close()
	require.NoError(test, err, "close db should return no error")
}
