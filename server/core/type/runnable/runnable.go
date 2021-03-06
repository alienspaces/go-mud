package runnable

import (
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

// Runnable -
type Runnable interface {
	Init(c configurer.Configurer, l logger.Logger, s storer.Storer, m modeller.Modeller) error
	Run(args map[string]interface{}) error
}
