package runner

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
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
	l = l.WithContext(logger.ContextApplication, "gameserver")

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

	// Handler configuration
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

// loggerWithContext provides a logger with function context
func loggerWithContext(l logger.Logger, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext("runner").WithFunctionContext(functionName)
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

func queryParamsToSQLOptions(qp *queryparam.QueryParams) *coresql.Options {

	if len(qp.SortColumns) == 0 {
		qp.SortColumns = []queryparam.SortColumn{
			{
				Col: "created_at",
			},
		}
	}

	return queryparam.ToSQLOptions(qp)
}
