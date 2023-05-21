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

	svr := Server{
		Config: c,
		Log:    l,
		Store:  s,
		Runner: r,
	}

	err := svr.Init()
	if err != nil {
		return nil, err
	}

	return &svr, nil
}

// Init -
func (svr *Server) Init() error {

	// TODO: alerting, retries
	return svr.Runner.Init(svr.Store)
}

// Run -
func (svr *Server) Run(args map[string]interface{}) error {

	// TODO:
	// - alerting on errors
	// - retries on start up
	// - reload  on config changes
	return svr.Runner.Run(args)
}
