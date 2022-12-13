package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-mud/server/core/prepare"

	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

// Runner - implements the runnerer interface
type Runner struct {
	Config            configurer.Configurer
	Log               logger.Logger
	Store             storer.Storer
	PrepareRepository preparer.Repository
	PrepareQuery      preparer.Query
	Model             modeller.Modeller

	// cli configuration - https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	App *cli.App

	// composable functions
	PreparerRepositoryFunc func() (preparer.Repository, error)
	PreparerQueryFunc      func() (preparer.Query, error)
	ModellerFunc           func() (modeller.Modeller, error)
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// Init - override to perform custom initialization
func (rnr *Runner) Init(s storer.Storer) error {

	if rnr.Log == nil {
		return fmt.Errorf("logger is nil, cannot initialise CLI runner")
	}

	rnr.Log.Info("** Initialise **")

	// Storer
	rnr.Store = s
	if rnr.Store == nil {
		msg := "storer undefined, cannot init runner"
		rnr.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	// Initialise storer
	err := rnr.Store.Init()
	if err != nil {
		rnr.Log.Warn("Failed store init >%v<", err)
		return err
	}

	// Repository
	if rnr.PreparerRepositoryFunc == nil {
		rnr.PreparerRepositoryFunc = rnr.PreparerRepository
	}

	p, err := rnr.PreparerRepositoryFunc()
	if err != nil {
		rnr.Log.Warn("Failed preparer func >%v<", err)
		return err
	}

	rnr.PrepareRepository = p
	if rnr.PrepareRepository == nil {
		rnr.Log.Warn("Repository is nil, cannot continue")
		return err
	}

	if rnr.PreparerQueryFunc == nil {
		rnr.PreparerQueryFunc = rnr.PreparerQuery
	}

	pCfg, err := rnr.PreparerQueryFunc()
	if err != nil {
		rnr.Log.Warn("Failed preparer config func >%v<", err)
		return err
	}

	rnr.PrepareQuery = pCfg
	if rnr.PrepareQuery == nil {
		rnr.Log.Warn("Repository config is nil, cannot continue")
		return err
	}

	db, err := rnr.Store.GetDb()
	if err != nil {
		rnr.Log.Warn("Failed getting database handle >%v<", err)
		return err
	}

	// Initialise preparer
	if err = rnr.PrepareRepository.Init(db); err != nil {
		rnr.Log.Warn("Failed preparer init >%v<", err)
		return err
	}
	if err = rnr.PrepareQuery.Init(db); err != nil {
		rnr.Log.Warn("Failed preparer config init >%v<", err)
		return err
	}

	// Modeller
	if rnr.ModellerFunc == nil {
		rnr.ModellerFunc = rnr.Modeller
	}

	return nil
}

// Run - Runs the CLI application.
func (rnr *Runner) Run(args map[string]interface{}) (err error) {

	rnr.Log.Debug("** Run **")

	// store init
	tx, err := rnr.Store.GetTx()
	if err != nil {
		rnr.Log.Warn("Failed getting tx >%v<", err)
		return err
	}

	// modeller
	m, err := rnr.ModellerFunc()
	if err != nil {
		rnr.Log.Warn("Failed modeller func >%v<", err)
		return err
	}

	if m == nil {
		rnr.Log.Warn("Modeller is nil, cannot continue")
		return err
	}

	// model init
	err = m.Init(rnr.PrepareRepository, rnr.PrepareQuery, tx)
	if err != nil {
		rnr.Log.Warn("Failed model init >%v<", err)
		return err
	}
	rnr.Model = m

	// run
	err = rnr.App.Run(os.Args)
	if err != nil {
		rnr.Log.Warn("Failed running app >%v<", err)

		// Rollback database transaction on error
		tx.Rollback()
		return err
	}

	// Commit database transaction
	err = tx.Commit()
	if err != nil {
		rnr.Log.Warn("Failed database transaction commit >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) PreparerRepository() (preparer.Repository, error) {

	rnr.Log.Info("** Repository **")

	p, err := prepare.NewRepositoryPreparer(rnr.Log)
	if err != nil {
		rnr.Log.Warn("Failed new prepare repository >%v<", err)
		return nil, err
	}

	return p, nil
}

func (rnr *Runner) PreparerQuery() (preparer.Query, error) {

	rnr.Log.Info("** Query **")

	p, err := prepare.NewQueryPreparer(rnr.Log)
	if err != nil {
		rnr.Log.Warn("Failed new prepare query >%v<", err)
		return nil, err
	}

	return p, nil
}

// Modeller - default ModellerFunc does not provide a modeller
func (rnr *Runner) Modeller() (modeller.Modeller, error) {

	rnr.Log.Info("** Modeller **")

	return nil, nil
}
