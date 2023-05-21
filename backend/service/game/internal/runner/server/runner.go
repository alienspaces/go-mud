package runner

import (
	"fmt"

	coreconfig "gitlab.com/alienspaces/go-mud/backend/core/config"
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

	// Handler config
	hc := r.CharacterHandlerConfig(nil)
	hc = r.DungeonHandlerConfig(hc)
	hc = r.DungeonCharacterHandlerConfig(hc)
	hc = r.DungeonLocationHandlerConfig(hc)
	hc = r.ActionHandlerConfig(hc)
	hc = r.DocumentationHandlerConfig(hc)

	appHome := r.Config.Get(coreconfig.AppServerHome)

	hc, err = server.ResolveHandlerSchemaLocationRoot(hc, appHome)
	if err != nil {
		err := fmt.Errorf("failed to resolve template API handler location root >%v<", err)
		r.Log.Warn(err.Error())
		return nil, err
	}

	hc = server.ResolveHandlerSchemaLocation(hc, "schema/game")
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
	m, err := rnr.InitTx(l)
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

func mergeHandlerConfigs(hc1 map[string]server.HandlerConfig, hc2 map[string]server.HandlerConfig) map[string]server.HandlerConfig {
	if hc1 == nil {
		hc1 = map[string]server.HandlerConfig{}
	}
	for a, b := range hc2 {
		hc1[a] = b
	}
	return hc1
}
