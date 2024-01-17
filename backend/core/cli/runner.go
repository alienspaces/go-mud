package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-mud/backend/core/prepare"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

// Runner - implements the runnerer interface

type Runner struct {
	Config             configurer.Configurer
	Log                logger.Logger
	Store              storer.Storer
	RepositoryPreparer preparer.Repository
	QueryPreparer      preparer.Query
	Model              modeller.Modeller

	// cli configuration - https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	App *cli.App

	// RepositoryPreparerFunc returns a repository preparer. Note that the runner InitDB method
	// must be called beforehand to assign the default query preparer.
	RepositoryPreparerFunc func() (preparer.Repository, error)
	// QueryPreparerFunc returns a query preparer. Note that the runner InitDB method must be
	// called beforehand to assign the default query preparer.
	QueryPreparerFunc func() (preparer.Query, error)
	ModellerFunc      func() (modeller.Modeller, error)

	// Initialisation will be deferred, it becomes the responsiblity
	// of the runner implementation to call init
	DeferModelInitialisation bool
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// Init - override to perform custom initialization
func (rnr *Runner) Init(s storer.Storer) error {

	rnr.Log.Info("** Init **")

	// Storer
	rnr.Store = s

	if !rnr.DeferModelInitialisation && rnr.Store == nil {
		msg := "store undefined, cannot init runner"
		rnr.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	return nil
}

// InitDB initialises a database connection, the repository preparer and the query preparer
func (rnr *Runner) InitDB() error {

	rnr.Log.Info("** InitDB **")

	db, err := rnr.Store.GetDb()
	if err != nil {
		rnr.Log.Warn("Failed getting database handle >%v<", err)
		return err
	}

	// Repository preparer
	if rnr.RepositoryPreparerFunc == nil {
		rnr.RepositoryPreparerFunc = rnr.defaultRepositoryPreparerFunc
	}

	repoPreparer, err := rnr.RepositoryPreparerFunc()
	if err != nil {
		rnr.Log.Warn("Failed getting repository preparer func >%v<", err)
		return err
	}

	rnr.RepositoryPreparer = repoPreparer
	if rnr.RepositoryPreparer == nil {
		rnr.Log.Warn("RepositoryPreparer is nil, cannot continue")
		return err
	}

	if err = rnr.RepositoryPreparer.Init(db); err != nil {
		rnr.Log.Warn("Failed preparer init >%v<", err)
		return err
	}

	// Query preparer
	if rnr.QueryPreparerFunc == nil {
		rnr.QueryPreparerFunc = rnr.defaultQueryPreparerFunc
	}

	queryPreparer, err := rnr.QueryPreparerFunc()
	if err != nil {
		rnr.Log.Warn("Failed getting query preparer func >%v<", err)
		return err
	}

	rnr.QueryPreparer = queryPreparer
	if rnr.QueryPreparer == nil {
		rnr.Log.Warn("QueryPreparer config is nil, cannot continue")
		return err
	}

	if err = rnr.QueryPreparer.Init(db); err != nil {
		rnr.Log.Warn("Failed preparer config init >%v<", err)
		return err
	}

	return nil
}

// Run - Runs the CLI application.
func (rnr *Runner) Run(args map[string]interface{}) (err error) {

	rnr.Log.Debug("** Run **")

	if !rnr.DeferModelInitialisation {
		err := rnr.InitModel()
		if err != nil {
			rnr.Log.Warn("failed model init >%v<", err)
			return err
		}
	}

	// Run
	err = rnr.App.Run(os.Args)
	if err != nil {
		rnr.Log.Warn("failed running app >%v<", err)

		// Rollback database transaction on error
		if rnr.Model != nil {
			rnr.Log.Warn("rolling back database transaction")
			rnr.Model.Rollback()
		}

		return err
	}

	// Commit database transaction
	if rnr.Model != nil {
		rnr.Log.Info("committing database transaction")
		err = rnr.Model.Commit()
		if err != nil {
			rnr.Log.Warn("Failed model commit >%v<", err)
			return err
		}
	}

	return nil
}

// InitModel iniitialises a database connection, sources a new model and initialises the model
// with a new database transaction.
func (rnr *Runner) InitModel() error {

	rnr.Log.Info("** Initialising model **")

	err := rnr.InitDB()
	if err != nil {
		rnr.Log.Warn("failed initialising db >%v<", err)
		return err
	}

	if rnr.ModellerFunc == nil {
		rnr.ModellerFunc = rnr.DefaultModellerFunc
	}

	m, err := rnr.ModellerFunc()
	if err != nil {
		rnr.Log.Warn("Failed getting modeller >%v<", err)
		return err
	}

	if m == nil {
		rnr.Log.Warn("Model is nil, cannot continue")
		return err
	}

	rnr.Model = m

	err = rnr.InitModelTx()
	if err != nil {
		rnr.Log.Warn("Failed model init tx >%v<", err)
		return err
	}

	return nil
}

// InitModelTx initialises the model with a new database transaction.
func (rnr *Runner) InitModelTx() error {

	rnr.Log.Info("** Initialising model tx **")

	tx, err := rnr.Store.GetTx()
	if err != nil {
		rnr.Log.Warn("Failed getting store transaction >%v<", err)
		return err
	}

	err = rnr.Model.Init(rnr.RepositoryPreparer, rnr.QueryPreparer, tx)
	if err != nil {
		rnr.Log.Warn("Failed model init >%v<", err)
		return err
	}

	return nil
}

// defaultRepositoryPreparerFunc - returns a default uninitialised repository preparer, set the
// property RepositoryPreparerFunc to provide your own custom repository preparer.
func (rnr *Runner) defaultRepositoryPreparerFunc() (preparer.Repository, error) {

	rnr.Log.Debug("** defaultRepositoryPreparerFunc **")

	p, err := prepare.NewRepositoryPreparer(rnr.Log)
	if err != nil {
		rnr.Log.Warn("Failed new prepare repository >%v<", err)
		return nil, err
	}

	return p, nil
}

// defaultQueryPreparerFunc - returns a default uninitialised query preparer, set the
// property QueryPreparerFunc to provide your own custom query preparer.
func (rnr *Runner) defaultQueryPreparerFunc() (preparer.Query, error) {

	rnr.Log.Debug("** defaultQueryPreparerFunc **")

	p, err := prepare.NewQueryPreparer(rnr.Log)
	if err != nil {
		rnr.Log.Warn("Failed new prepare query >%v<", err)
		return nil, err
	}

	return p, nil
}

// DefaultModellerFunc does not provide a modeller, set the property ModellerFunc to
// provide your own custom modeller.
func (rnr *Runner) DefaultModellerFunc() (modeller.Modeller, error) {

	rnr.Log.Debug("** DefaultModellerFunc **")

	return nil, nil
}
