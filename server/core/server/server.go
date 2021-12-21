package server

import (
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

// Server -
type Server struct {
	Config configurer.Configurer
	Log    logger.Logger
	Store  storer.Storer
	Model  modeller.Modeller
	Runner runnable.Runnable
}

// NewServer -
func NewServer(c configurer.Configurer, l logger.Logger, s storer.Storer, m modeller.Modeller, r runnable.Runnable) (*Server, error) {

	svc := Server{
		Config: c,
		Log:    l,
		Store:  s,
		Model:  m,
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

	return svc.Runner.Init(svc.Config, svc.Log, svc.Store, svc.Model)
}

// Run -
func (svc *Server) Run(args map[string]interface{}) error {

	// TODO:
	// - alerting on errors
	// - retries on start up
	// - reload  on config changes
	return svc.Runner.Run(args)
}
