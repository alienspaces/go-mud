package runner

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
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

	// Handler config
	hc := r.CharacterHandlerConfig(nil)
	hc = r.DungeonHandlerConfig(hc)
	hc = r.ActionHandlerConfig(hc)
	hc = r.DocumentationHandlerConfig(hc)

	r.HandlerConfig = hc

	return &r, nil
}

// Modeller -
func (rnr *Runner) Modeller(l logger.Logger) (modeller.Modeller, error) {

	l.Info("** Dungeon Model **")

	m, err := model.NewModel(rnr.Config, l, rnr.Store)
	if err != nil {
		l.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
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
