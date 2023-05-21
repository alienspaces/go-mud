package server

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

// TestRunner - allow Init and Run functions to be defined by tests
type TestRunner struct {
	Runner
	InitFunc func(c configurer.Configurer, l logger.Logger, s storer.Storer) error
}

func (rnr *TestRunner) Init(c configurer.Configurer, l logger.Logger, s storer.Storer) error {
	rnr.Config = c
	rnr.Log = l

	if rnr.InitFunc == nil {
		return rnr.Runner.Init(s)
	}
	return rnr.InitFunc(c, l, s)
}

func TestRunnerInit(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	tr := TestRunner{}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	// Test init override with failure
	tr.InitFunc = func(c configurer.Configurer, l logger.Logger, s storer.Storer) error {
		return errors.New("Init failed")
	}
	err = tr.Init(c, l, s)
	require.Error(t, err, "Runner Init returns with error")
}

func TestRunnerServerError(t *testing.T) {
	t.Parallel()

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	tr := TestRunner{}
	tr.RunHTTPFunc = func(args map[string]interface{}) (*http.Server, error) {
		return nil, fmt.Errorf("Run server error")
	}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	err = tr.Run(nil)
	require.Error(t, err, "Runner Run returns with error")
}

func TestRunnerDaemonError(t *testing.T) {
	t.Parallel()

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	tr := TestRunner{}
	tr.RunDaemonFunc = func(args map[string]interface{}) error {
		return fmt.Errorf("Run daemon error")
	}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	err = tr.Run(nil)
	require.Error(t, err, "Runner Run returns with error")
}

func Test_registerRoutes(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	tr := TestRunner{}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	r := httprouter.New()
	router, err := tr.RegisterRoutes(r)
	require.NoError(t, err, "Router returns without error")
	require.NotNil(t, router, "Router returns a router")

	// Test default configured routes
	handle, _, _ := router.Lookup(http.MethodGet, "/healthz")
	require.NotNil(t, handle, "Handle for /healthz is not nil")

	// Test custom routes
	tr.RouterFunc = func(r *httprouter.Router) (*httprouter.Router, error) {
		h, err := tr.ApplyMiddleware(HandlerConfig{Path: "/custom"}, tr.HandlerFunc)
		if err != nil {
			return nil, err
		}
		r.GET("/custom", h)
		return r, nil
	}

	r = httprouter.New()
	router, err = tr.RegisterRoutes(r)
	require.NoError(t, err, "Router returns without error")
	require.NotNil(t, router, "Router returns a router")

	// Test custom configured routes
	handle, _, _ = router.Lookup(http.MethodGet, "/custom")
	require.NotNil(t, handle, "Handle for /custom is not nil")

	// Test custom router error
	tr.RouterFunc = func(r *httprouter.Router) (*httprouter.Router, error) {
		return nil, errors.New("Failed router")
	}

	r = httprouter.New()
	router, err = tr.RegisterRoutes(r)
	require.Error(t, err, "Router returns with error")
	require.Nil(t, router, "Router returns nil")
}

func Test_ApplyMiddleware(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	tr := TestRunner{}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	// Test default middleware
	handle, err := tr.ApplyMiddleware(HandlerConfig{Path: "/"}, tr.HandlerFunc)
	require.NoError(t, err, "Middleware returns without error")
	require.NotNil(t, handle, "Middleware returns a handle")

	// Test custom middleware
	tr.HandlerMiddlewareFuncs = func() []MiddlewareFunc {
		return []MiddlewareFunc{
			func(hc HandlerConfig, h Handle) (Handle, error) {
				return nil, errors.New("Failed middleware")
			},
		}
	}

	handle, err = tr.ApplyMiddleware(HandlerConfig{Path: "/"}, tr.HandlerFunc)
	require.Error(t, err, "Middleware returns with error")
	require.Nil(t, handle, "Middleware returns nil")
}
