package runner

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
)

// Runner -
type Runner struct {
	server.Runner
}

// Fault -
type Fault struct {
	Error error
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// NewRunner -
func NewRunner(c configurer.Configurer, l logger.Logger) (*Runner, error) {

	r := Runner{}

	r.Log = l
	if r.Log == nil {
		msg := "logger undefined, cannot init runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	r.Config = c
	if r.Config == nil {
		msg := "configurer undefined, cannot init runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	r.RouterFunc = r.Router
	r.MiddlewareFunc = r.Middleware
	r.HandlerFunc = r.Handler
	r.ModellerFunc = r.Modeller
	r.RunDaemonFunc = r.RunDaemon

	// Handler config
	hc := r.CharacterHandlerConfig(nil)
	hc = r.DungeonHandlerConfig(hc)
	hc = r.DungeonCharacterHandlerConfig(hc)
	hc = r.DungeonLocationHandlerConfig(hc)
	hc = r.ActionHandlerConfig(hc)
	hc = r.DocumentationHandlerConfig(hc)

	r.HandlerConfig = hc

	return &r, nil
}

// Modeller -
func (rnr *Runner) Modeller(l logger.Logger) (modeller.Modeller, error) {
	m, err := model.NewModel(rnr.Config, l, rnr.Store)
	if err != nil {
		l.Warn("failed new model >%v<", err)
		return nil, err
	}
	return m, nil
}

func (rnr *Runner) initModeller(l logger.Logger) (*model.Model, error) {
	m, err := rnr.InitModeller(l)
	if err != nil {
		l.Warn("failed initialising database transaction, cannot authen >%v<", err)
		return nil, err
	}
	return m.(*model.Model), nil
}

// loggerWithContext provides a logger with function context
func loggerWithContext(l logger.Logger, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext("game/server").WithFunctionContext(functionName)
}

func mergeHandlerConfigs(hc1 map[server.HandlerConfigKey]server.HandlerConfig, hc2 map[server.HandlerConfigKey]server.HandlerConfig) map[server.HandlerConfigKey]server.HandlerConfig {
	if hc1 == nil {
		hc1 = map[server.HandlerConfigKey]server.HandlerConfig{}
	}
	for a, b := range hc2 {
		hc1[a] = b
	}
	return hc1
}
