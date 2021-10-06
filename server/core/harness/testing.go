package harness

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/prepare"
	"gitlab.com/alienspaces/go-boilerplate/server/core/store"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

// CreateDataFunc - callback function that creates test data
type CreateDataFunc func() error

// RemoveDataFunc - callback function that removes test data
type RemoveDataFunc func() error

// Testing -
type Testing struct {
	Config  configurer.Configurer
	Log     logger.Logger
	Store   storer.Storer
	Prepare preparer.Preparer
	Model   modeller.Modeller

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

	// default dependencies
	c, l, s, err := t.NewDefaultDependencies()
	if err != nil {
		return err
	}

	t.Config = c
	t.Log = l
	t.Store = s

	// preparer
	t.Prepare, err = prepare.NewPrepare(t.Log)
	if err != nil {
		t.Log.Warn("Failed new preparer >%v<", err)
		return err
	}

	db, err := t.Store.GetDb()
	if err != nil {
		t.Log.Warn("Failed getting database handle >%v<", err)
		return err
	}

	err = t.Prepare.Init(db)
	if err != nil {
		t.Log.Warn("Failed preparer init >%v<", err)
		return err
	}

	t.Log.Debug("Preparer ready")

	// modeller
	t.Model, err = t.ModellerFunc()
	if err != nil {
		t.Log.Warn("Failed new modeller >%v<", err)
		return err
	}

	t.Log.Debug("Modeller ready")

	return nil
}

// NewDefaultDependencies -
func (t *Testing) NewDefaultDependencies() (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := config.NewConfig(nil, false)
	if err != nil {
		return nil, nil, nil, err
	}

	configVars := []string{
		// logger
		"APP_SERVER_LOG_LEVEL",
		// database
		"APP_SERVER_DB_HOST",
		"APP_SERVER_DB_PORT",
		"APP_SERVER_DB_NAME",
		"APP_SERVER_DB_USER",
		"APP_SERVER_DB_PASSWORD",
		// schema
		"APP_SERVER_SCHEMA_PATH",
		// jwt signing key
		"APP_SERVER_JWT_SIGNING_KEY",
	}
	for _, key := range configVars {
		err = c.Add(key, true)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, nil, err
	}

	// storer
	s, err := store.NewStore(c, l)
	if err != nil {
		return nil, nil, nil, err
	}

	err = s.Init()
	if err != nil {
		return nil, nil, nil, err
	}

	return c, l, s, nil
}

// InitTx -
func (t *Testing) InitTx(tx *sqlx.Tx) (err error) {

	// initialise our own database tx when none is provided
	if tx == nil {
		t.Log.Debug("Starting database tx")

		tx, err = t.Store.GetTx()
		if err != nil {
			t.Log.Warn("Failed getting database tx >%v<", err)
			return err
		}
	}

	err = t.Model.Init(t.Prepare, tx)
	if err != nil {
		t.Log.Warn("Failed modeller init >%v<", err)
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
		t.Log.Warn("Failed init >%v<", err)
		return err
	}

	// data function is expected to create and manage its own store
	if t.CreateDataFunc != nil {
		t.Log.Debug("Creating test data")
		err := t.CreateDataFunc()
		if err != nil {
			t.Log.Warn("Failed creating data >%v<", err)
			return err
		}
	}

	// commit data when configured, otherwise we are leaving
	// it up to tests to explicitly commit or rollback
	if t.CommitData {
		t.Log.Debug("Committing database tx")
		err = t.CommitTx()
		if err != nil {
			t.Log.Warn("Failed comitting data >%v<", err)
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
		t.Log.Warn("Failed init >%v<", err)
		return err
	}

	// data function is expected to create and manage its own store
	if t.RemoveDataFunc != nil {
		t.Log.Debug("Removing test data")
		err := t.RemoveDataFunc()
		if err != nil {
			t.Log.Warn("Failed removing data >%v<", err)
			return err
		}
	}

	// commit data when configured, otherwise we are leaving
	// it up to tests to explicitly commit or rollback
	if t.CommitData {
		t.Log.Debug("Committing database tx")
		err = t.CommitTx()
		if err != nil {
			t.Log.Warn("Failed comitting data >%v<", err)
			return err
		}
	}

	return nil
}
