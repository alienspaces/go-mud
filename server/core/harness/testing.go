package harness

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/config"
	"gitlab.com/alienspaces/go-mud/server/core/log"
	"gitlab.com/alienspaces/go-mud/server/core/prepare"
	"gitlab.com/alienspaces/go-mud/server/core/store"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

// CreateDataFunc - callback function that creates test data
type CreateDataFunc func() error

// RemoveDataFunc - callback function that removes test data
type RemoveDataFunc func() error

// Testing -
type Testing struct {
	Config            configurer.Configurer
	Log               logger.Logger
	Store             storer.Storer
	PrepareRepository preparer.Repository
	PrepareQuery      preparer.Query
	Model             modeller.Modeller

	// Configuration
	CommitData bool

	// Modeller function
	ModellerFunc func() (modeller.Modeller, error)

	// Composable functions
	CreateDataFunc CreateDataFunc
	RemoveDataFunc RemoveDataFunc

	// Private
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

	// storer
	if t.Store == nil {
		t.Store, err = store.NewStore(t.Config, t.Log)
		if err != nil {
			return err
		}
	}

	err = t.Store.Init()
	if err != nil {
		t.Log.Warn("failed storer init >%v<", err)
		return err
	}

	// preparer
	t.PrepareRepository, err = prepare.NewPrepare(t.Log)
	if err != nil {
		t.Log.Warn("failed new preparer >%v<", err)
		return err
	}

	t.PrepareQuery, err = prepare.NewQuery(t.Log)
	if err != nil {
		t.Log.Warn("failed new preparer config >%v<", err)
		return err
	}

	db, err := t.Store.GetDb()
	if err != nil {
		t.Log.Warn("failed getting database handle >%v<", err)
		return err
	}

	err = t.PrepareRepository.Init(db)
	if err != nil {
		t.Log.Warn("failed preparer init >%v<", err)
		return err
	}

	err = t.PrepareQuery.Init(db)
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
func (t *Testing) InitTx(tx *sqlx.Tx) (err error) {

	// initialise our own database tx when none is provided
	if tx == nil {
		t.Log.Debug("Starting database tx")

		tx, err = t.Store.GetTx()
		if err != nil {
			t.Log.Warn("failed getting database tx >%v<", err)
			return err
		}
	}

	err = t.Model.Init(t.PrepareRepository, t.PrepareQuery, tx)
	if err != nil {
		t.Log.Warn("failed modeller init >%v<", err)
		return err
	}

	t.tx = tx

	return nil
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
func (t *Testing) Setup() (err error) {

	// init
	err = t.InitTx(nil)
	if err != nil {
		t.Log.Warn("failed init >%v<", err)
		return err
	}

	// data function is expected to create and manage its own store
	if t.CreateDataFunc != nil {
		t.Log.Debug("Creating test data")
		err := t.CreateDataFunc()
		if err != nil {
			t.Log.Warn("failed creating data >%v<", err)
			return err
		}
	}

	// commit data when configured, otherwise we are leaving
	// it up to tests to explicitly commit or rollback
	if t.CommitData {
		t.Log.Debug("Committing database tx")
		err = t.CommitTx()
		if err != nil {
			t.Log.Warn("failed comitting data >%v<", err)
			return err
		}
	}

	return nil
}

// Teardown -
func (t *Testing) Teardown() (err error) {

	// init
	err = t.InitTx(nil)
	if err != nil {
		t.Log.Warn("failed init >%v<", err)
		return err
	}

	// data function is expected to create and manage its own store
	if t.RemoveDataFunc != nil {
		t.Log.Debug("Removing test data")
		err := t.RemoveDataFunc()
		if err != nil {
			t.Log.Warn("failed removing data >%v<", err)
			return err
		}
	}

	// commit data when configured, otherwise we are leaving
	// it up to tests to explicitly commit or rollback
	if t.CommitData {
		t.Log.Debug("Committing database tx")
		err = t.CommitTx()
		if err != nil {
			t.Log.Warn("failed committing data >%v<", err)
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
