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

	r.HandlerFunc = r.Handler
	r.ModellerFunc = r.Modeller
	r.RunDaemonFunc = r.RunDaemon
	r.HandlerMiddlewareFuncs = r.middlewareFuncs

	// Handler configuration
	hc := r.TemplateHandlerConfig(nil)

	r.HandlerConfig = hc

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

// loggerWithContext provides a logger with function context
func loggerWithContext(l logger.Logger, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext("template/server").WithFunctionContext(functionName)
}

func mergeHandlerConfigs(hc1 map[string]server.HandlerConfig, hc2 map[string]server.HandlerConfig) map[string]server.HandlerConfig {
	if hc1 == nil {
		hc1 = map[string]server.HandlerConfig{}
	}
	for a, b := range hc2 {
		hc1[a] = b
	}
	return hc1
}
