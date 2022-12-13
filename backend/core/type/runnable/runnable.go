package runnable

import (
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

// Runnable -
type Runnable interface {
	Init(s storer.Storer) error
	Run(args map[string]interface{}) error
}

// StatelessRunnable -
type StatelessRunnable interface {
	Init() error
	Run(args map[string]interface{}) error
}
