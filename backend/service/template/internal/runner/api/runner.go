package runner

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/model"
)

// Runner -
type Runner struct {
	server.Runner
	// Service specific configuration
	config Config
}

// Fault -
type Fault struct {
	Error error
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// NewRunner -
func NewRunner(c configurer.Configurer, l logger.Logger) (*Runner, error) {

	cr, err := server.NewRunner(c, l)
	if err != nil {
		err := fmt.Errorf("failed core runner >%v<", err)
		l.Warn(err.Error())
		return nil, err
	}

	cfg, err := NewConfig(c)
	if err != nil {
		return nil, err
	}

	r := Runner{
		Runner: *cr,
		config: *cfg,
	}

	r.ModellerFunc = r.Modeller

	r.HandlerFunc = r.Handler
	r.HandlerMiddlewareFuncs = r.middlewareFuncs

	// Handler configuration
	handlerConfig, err := buildHandlerConfig(r)
	if err != nil {
		err := fmt.Errorf("failed to initialise server.MessageConfig >%v<", err)
		l.Warn(err.Error())
		return nil, err
	}

	r.HandlerConfig = handlerConfig

	// Authentication
	if err = server.ValidateAuthenticationTypes(r.HandlerConfig); err != nil {
		l.Warn(err.Error())
		return nil, err
	}

	// Daemon
	r.RunDaemonFunc = r.RunDaemon

	return &r, nil
}

func (rnr *Runner) middlewareFuncs() []server.MiddlewareFunc {
	return []server.MiddlewareFunc{
		rnr.WaitMiddleware,
		rnr.ParamMiddleware,
		rnr.DataMiddleware,
		rnr.AuthzMiddleware,
		rnr.AuthenMiddleware,
		rnr.TxMiddleware,
		rnr.CorrelationMiddleware,
	}
}

// Modeller -
func (rnr *Runner) Modeller(l logger.Logger) (modeller.Modeller, error) {
	l.Info("Template Modeller")

	m, err := model.NewModel(rnr.Config, l, rnr.Store)
	if err != nil {
		l.Warn("Failed new modeller >%v<", err)
		return nil, err
	}

	return m, nil
}
