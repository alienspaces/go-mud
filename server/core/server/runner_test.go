package server

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

// TestRunner - allow Init and Run functions to be defined by tests
type TestRunner struct {
	Runner
	InitFunc func(c configurer.Configurer, l logger.Logger, s storer.Storer) error
}

func (rnr *TestRunner) Init(c configurer.Configurer, l logger.Logger, s storer.Storer) error {
	if rnr.InitFunc == nil {
		return rnr.Runner.Init(c, l, s)
	}
	return rnr.InitFunc(c, l, s)
}

func TestRunnerInit(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tr := TestRunner{}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	// test init override with failure
	tr.InitFunc = func(c configurer.Configurer, l logger.Logger, s storer.Storer) error {
		return errors.New("Init failed")
	}
	err = tr.Init(c, l, s)
	require.Error(t, err, "Runner Init returns with error")
}

func TestRunnerServerError(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tr := TestRunner{}
	tr.RunHTTPFunc = func(args map[string]interface{}) error {
		return fmt.Errorf("Run server error")
	}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	err = tr.Run(nil)
	require.Error(t, err, "Runner Run returns with error")
}

func TestRunnerDaemonError(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tr := TestRunner{}
	tr.RunDaemonFunc = func(args map[string]interface{}) error {
		return fmt.Errorf("Run daemon error")
	}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	err = tr.Run(nil)
	require.Error(t, err, "Runner Run returns with error")
}

func TestRunnerRouter(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tr := TestRunner{}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	// test default routes
	router, err := tr.DefaultRouter()
	require.NoError(t, err, "DefaultRouter returns without error")
	require.NotNil(t, router, "DefaultRouter returns a router")

	// test default configured routes
	handle, params, redirect := router.Lookup(http.MethodGet, "/")
	require.NotNil(t, handle, "Handle is not nil")

	t.Logf("Default route /")
	t.Logf("Have handler >%#v<", handle)
	t.Logf("Have params >%v<", params)
	t.Logf("Have redirect >%t<", redirect)

	// test custom routes
	tr.RouterFunc = func(router *httprouter.Router) error {
		h, err := tr.DefaultMiddleware(HandlerConfig{Path: "/custom"}, tr.HandlerFunc)
		if err != nil {
			return err
		}
		router.GET("/custom", h)
		return nil
	}

	router, err = tr.DefaultRouter()
	require.NoError(t, err, "DefaultRouter returns without error")
	require.NotNil(t, router, "DefaultRouter returns a router")

	// test custom configured routes
	handle, params, redirect = router.Lookup(http.MethodGet, "/custom")
	require.NotNil(t, handle, "Handle is not nil")

	t.Logf("Custom route /custom")
	t.Logf("Have handler >%#v<", handle)
	t.Logf("Have params >%v<", params)
	t.Logf("Have redirect >%t<", redirect)

	// test custom routes error
	tr.RouterFunc = func(router *httprouter.Router) error {
		return errors.New("Failed router")
	}

	router, err = tr.DefaultRouter()
	require.Error(t, err, "DefaultRouter returns with error")
	require.Nil(t, router, "DefaultRouter returns nil")
}

func TestRunnerMiddleware(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tr := TestRunner{}

	err = tr.Init(c, l, s)
	require.NoError(t, err, "Runner Init returns without error")

	// test default middleware
	handle, err := tr.DefaultMiddleware(HandlerConfig{Path: "/"}, tr.HandlerFunc)
	require.NoError(t, err, "DefaultMiddleware returns without error")
	require.NotNil(t, handle, "DefaultMiddleware returns a handle")

	// test custom middleware
	tr.MiddlewareFunc = func(h HandlerFunc) (HandlerFunc, error) {
		return nil, errors.New("Failed middleware")
	}

	handle, err = tr.DefaultMiddleware(HandlerConfig{Path: "/"}, tr.HandlerFunc)
	require.Error(t, err, "DefaultMiddleware returns with error")
	require.Nil(t, handle, "DefaultMiddleware returns nil")
}
