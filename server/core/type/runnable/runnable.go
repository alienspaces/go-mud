package runnable

import (
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

// Runnable -
type Runnable interface {
	Init(c configurer.Configurer, l logger.Logger, s storer.Storer) error
	Run(args map[string]interface{}) error
}
