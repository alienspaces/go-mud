package server

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

// Server -
type Server struct {
	Config configurer.Configurer
	Log    logger.Logger
	Store  storer.Storer
	Runner runnable.Runnable
}

// NewServer -
func NewServer(c configurer.Configurer, l logger.Logger, s storer.Storer, r runnable.Runnable) (*Server, error) {

	svc := Server{
		Config: c,
		Log:    l,
		Store:  s,
		Runner: r,
	}

	err := svc.Init()
	if err != nil {
		return nil, err
	}

	return &svc, nil
}

// Init -
func (svc *Server) Init() error {

	err := svc.Store.Init()
	if err != nil {
		return err
	}
	return svc.Runner.Init(svc.Store)
}

// Run -
func (svc *Server) Run(args map[string]interface{}) error {

	// TODO: (core) Alerts, retries, reloads
	// - alerting on errors
	// - retries on start up
	// - reload  on config changes
	return svc.Runner.Run(args)
}
