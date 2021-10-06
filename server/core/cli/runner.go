package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/runnable"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

// Runner - implements the runnerer interface
type Runner struct {
	Config  configurer.Configurer
	Log     logger.Logger
	Store   storer.Storer
	Prepare preparer.Preparer
	Model   modeller.Modeller

	// cli configuration - https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	App *cli.App

	// composable functions
	PreparerFunc func() (preparer.Preparer, error)
	ModellerFunc func() (modeller.Modeller, error)
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// Init - override to perform custom initialization
func (rnr *Runner) Init(c configurer.Configurer, l logger.Logger, s storer.Storer) error {

	rnr.Log = l
	if rnr.Log == nil {
		msg := "Logger undefined, cannot init runner"
		return fmt.Errorf(msg)
	}

	rnr.Log.Info("** Initialise **")

	rnr.Config = c
	if rnr.Config == nil {
		msg := "Configurer undefined, cannot init runner"
		rnr.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	// Storer
	rnr.Store = s
	if rnr.Store == nil {
		msg := "Storer undefined, cannot init runner"
		rnr.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	// Initialise storer
	err := rnr.Store.Init()
	if err != nil {
		rnr.Log.Warn("Failed store init >%v<", err)
		return err
	}

	// Preparer
	if rnr.PreparerFunc == nil {
		rnr.PreparerFunc = rnr.Preparer
	}

	p, err := rnr.PreparerFunc()
	if err != nil {
		rnr.Log.Warn("Failed preparer func >%v<", err)
		return err
	}

	rnr.Prepare = p
	if rnr.Prepare == nil {
		rnr.Log.Warn("Preparer is nil, cannot continue")
		return err
	}

	db, err := rnr.Store.GetDb()
	if err != nil {
		rnr.Log.Warn("Failed getting database handle >%v<", err)
		return err
	}

	// Initialise preparer
	err = rnr.Prepare.Init(db)
	if err != nil {
		rnr.Log.Warn("Failed preparer init >%v<", err)
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
	err = m.Init(rnr.Prepare, tx)
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

// Preparer - default PreparerFunc does not provide a modeller
func (rnr *Runner) Preparer() (preparer.Preparer, error) {

	rnr.Log.Info("** Preparer **")

	return nil, nil
}

// Modeller - default ModellerFunc does not provide a modeller
func (rnr *Runner) Modeller() (modeller.Modeller, error) {

	rnr.Log.Info("** Modeller **")

	return nil, nil
}
