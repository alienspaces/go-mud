package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/prepare"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/runnable"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

const (
	// ConfigKeyValidateSchemaLocation - Directory location of JSON schema's
	ConfigKeyValidateSchemaLocation string = "validateSchemaLocation"
	// ConfigKeyValidateMainSchema - Main schema that can include reference schema's
	ConfigKeyValidateMainSchema string = "validateMainSchema"
	// ConfigKeyValidateReferenceSchemas - Schema referenced from the main schema
	ConfigKeyValidateReferenceSchemas string = "validateReferenceSchemas"

	// AuthTypeJWT -
	AuthTypeJWT string = "jwt"
)

var _ runnable.Runnable = &Runner{}

// HandlerFunc - custom service handle
type HandlerFunc func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params, queryParams map[string]interface{}, l logger.Logger, m modeller.Modeller)

// Runner - implements the runnerer interface
type Runner struct {
	Config  configurer.Configurer
	Log     logger.Logger
	Store   storer.Storer
	Prepare preparer.Preparer

	// configuration for routes, handlers and middleware
	HandlerConfig []HandlerConfig

	// composable functions
	RunHTTPFunc    func(args map[string]interface{}) error
	RunDaemonFunc  func(args map[string]interface{}) error
	RouterFunc     func(router *httprouter.Router) error
	MiddlewareFunc func(h HandlerFunc) (HandlerFunc, error)
	// HandlerFunc    func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params, queryParams map[string]interface{}, l logger.Logger, m modeller.Modeller)
	HandlerFunc  HandlerFunc
	PreparerFunc func(l logger.Logger) (preparer.Preparer, error)
	ModellerFunc func(l logger.Logger) (modeller.Modeller, error)
}

// MiddlewareConfig provides configuration for middlewares.
//
// AuthTypes - Supported authentication types.
//
// AuthRequiredAllRoles - All of these roles must exist within the list of roles found in authenticated claims.
//
// AuthRequiredAnyRole - Any of these roles must exist within the list of roles found in authenticated claims.
//
// AuthRequiredAllIdentities - All of these identity keys must have a defined value within the list of identity keys found in authenticated claims.
//
// AuthRequiredAnyIdentity - Any of these identity keys must have a defined value within the list of identity keys found in authenticated claims.
//
// ValidateSchemaLocation - Location of JSON schemas for this endpoint relative to `APP_SERVER_SCHEMA_PATH`.
//
// ValidateSchemaMain - The name of the main JSON schema document to load for this endpoint.
//
// ValidateSchemaReferences - A list of additional JSON schema reference documents to load for this endpoint.
//
// ValidatePathParams - Rules for validating path parameters
//
// ValidateQueryParams - A whitelist of allowed query parameters.
//
type MiddlewareConfig struct {

	// AuthTypes - What auth types are supported by this endpoint
	AuthTypes []string
	// AuthRequireAllRoles - Require all of these roles to access this endpoint
	AuthRequireAllRoles []string
	// AuthRequireAnyRole - Require any of these roles to access this endpoint
	AuthRequireAnyRole []string
	// AuthRequireAllIdentities - Required all of these identities to be defined to access this endpoint
	AuthRequireAllIdentities []string
	// AuthRequireAnyIdentity - Required any of these identities to be defined to access this endpoint
	AuthRequireAnyIdentity []string

	// Validate Schema - JSON schema validation
	ValidateSchemaLocation   string
	ValidateSchemaMain       string
	ValidateSchemaReferences []string

	// ValidatePathParams - Rules for validating path parameters
	ValidatePathParams map[string]ValidatePathParam

	// ValidateQueryParams - A whitelist of allowed query parameters
	ValidateQueryParams []string
}

// HandlerConfig - configuration for routes, handlers and middleware
type HandlerConfig struct {
	// Method - The HTTP method
	Method string
	// Path - The HTTP request URI including :parameter placeholders
	Path string
	// HandlerFunc - Function to handle requests for this method and path
	HandlerFunc func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller)
	// MiddlewareConfig -
	MiddlewareConfig MiddlewareConfig
	// DocumentationConfig -
	DocumentationConfig DocumentationConfig
}

// ValidatePathParam - Rules for validating a path parameter
type ValidatePathParam struct {
	MatchIdentity bool
}

// DocumentationConfig - Configuration describing how to document a route
type DocumentationConfig struct {
	Document    bool
	Description string
}

// NOTE: Request struct definitions are located in the top level `schema` module. We might
// want to consider moving the Response definitions there as well, especially if we decide
// to validate our own response payloads against schema definitions..

// Response -
type Response struct {
	Error      *ResponseError      `json:"error,omitempty"`
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

// ResponseError -
type ResponseError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ResponsePagination -
type ResponsePagination struct {
	Number int `json:"page_number"`
	Size   int `json:"page_size"`
	Count  int `json:"page_count"`
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// Init - override to perform custom initialization
func (rnr *Runner) Init(c configurer.Configurer, l logger.Logger, s storer.Storer) error {

	rnr.Log = l
	if rnr.Log == nil {
		msg := "Logger undefined, cannot init runner"
		return fmt.Errorf(msg)
	}

	rnr.Log.Debug("** Initialise **")

	rnr.Config = c
	if rnr.Config == nil {
		msg := "Configurer undefined, cannot init runner"
		rnr.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	rnr.Store = s
	if rnr.Store == nil {
		msg := "Storer undefined, cannot init runner"
		rnr.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	// Initialise storer
	err := rnr.Store.Init()
	if err != nil {
		rnr.Log.Warn("Failed store init >%v<", err)
		return err
	}

	// Preparer
	if rnr.PreparerFunc == nil {
		rnr.PreparerFunc = rnr.Preparer
	}

	p, err := rnr.PreparerFunc(l)
	if err != nil {
		rnr.Log.Warn("Failed preparer func >%v<", err)
		return err
	}

	rnr.Prepare = p
	if rnr.Prepare == nil {
		rnr.Log.Warn("Preparer is nil, cannot continue")
		return err
	}

	// run server
	if rnr.RunHTTPFunc == nil {
		rnr.RunHTTPFunc = rnr.RunHTTP
	}

	// run daemon
	if rnr.RunDaemonFunc == nil {
		rnr.RunDaemonFunc = rnr.RunDaemon
	}

	// prepare
	if rnr.PreparerFunc == nil {
		rnr.PreparerFunc = rnr.Preparer
	}

	// model
	if rnr.ModellerFunc == nil {
		rnr.ModellerFunc = rnr.Modeller
	}

	// http server - router
	if rnr.RouterFunc == nil {
		rnr.RouterFunc = rnr.Router
	}

	// http server - middleware
	if rnr.MiddlewareFunc == nil {
		rnr.MiddlewareFunc = rnr.Middleware
	}

	// http server - handler
	if rnr.HandlerFunc == nil {
		rnr.HandlerFunc = rnr.Handler
	}

	return nil
}

// TODO: Use this function from HTTP middleware Tx maybe..

// InitTx initialises a new database transaction returning a prepared modeller
func (rnr *Runner) InitTx(l logger.Logger) (modeller.Modeller, error) {

	// NOTE: The modeller is created and initialised with every request instead of
	// creating and assigning to a runner struct "Model" property at start up.
	// This prevents directly accessing a shared property from with the handler
	// function which is running in a goroutine. Otherwise accessing the "Model"
	// property would require locking and block simultaneous requests.

	// modeller
	if rnr.ModellerFunc == nil {
		l.Warn("Runner ModellerFunc is nil")
		return nil, fmt.Errorf("ModellerFunc is nil")
	}

	m, err := rnr.ModellerFunc(l)
	if err != nil {
		l.Warn("Failed ModellerFunc >%v<", err)
		return nil, err
	}

	if m == nil {
		l.Warn("Modeller is nil, cannot continue")
		return nil, err
	}

	tx, err := rnr.Store.GetTx()
	if err != nil {
		l.Warn("Failed getting DB connection >%v<", err)
		return m, err
	}

	err = m.Init(rnr.Prepare, tx)
	if err != nil {
		l.Warn("Failed init modeller >%v<", err)
		return m, err
	}

	return m, nil
}

// Run starts the HTTP server and daemon processes. Override to implement a custom run function.
func (rnr *Runner) Run(args map[string]interface{}) (err error) {

	rnr.Log.Debug("** Run **")

	// signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// run HTTP server
	go func() {
		rnr.Log.Debug("** Running HTTP server process **")
		err = rnr.RunHTTPFunc(args)
		if err != nil {
			rnr.Log.Error("Failed run server >%v<", err)
			sigChan <- syscall.SIGTERM
		}
		rnr.Log.Debug("** HTTP server process ended **")
	}()

	// run daemon server
	go func() {
		rnr.Log.Debug("** Running daemon process **")
		err = rnr.RunDaemonFunc(args)
		if err != nil {
			rnr.Log.Error("Failed run daemon >%v<", err)
			sigChan <- syscall.SIGTERM
		}
		rnr.Log.Debug("** Daemon process ended **")
	}()

	// wait
	sig := <-sigChan

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	err = fmt.Errorf("Received SIG >%v< Mem Alloc >%d MiB< TotalAlloc >%d MiB< Sys >%d MiB< NumGC >%d<",
		sig,
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC,
	)

	rnr.Log.Warn(">%v<", err)

	return err
}

// Preparer - default PreparerFunc, override this function for custom prepare
func (rnr *Runner) Preparer(l logger.Logger) (preparer.Preparer, error) {

	// NOTE: We have a good generic preparer so we'll provide that here

	l.Debug("** Preparer **")

	// Return the existing preparer if we already have one
	if rnr.Prepare != nil {
		l.Debug("Returning existing preparer")
		return rnr.Prepare, nil
	}

	l.Debug("Creating new preparer")

	p, err := prepare.NewPrepare(l)
	if err != nil {
		l.Warn("Failed new prepare >%v<", err)
		return nil, err
	}

	db, err := rnr.Store.GetDb()
	if err != nil {
		l.Warn("Failed getting database handle >%v<", err)
		return nil, err
	}

	err = p.Init(db)
	if err != nil {
		l.Warn("Failed preparer init >%v<", err)
		return nil, err
	}

	rnr.Prepare = p

	return p, nil
}

// Modeller - default ModellerFunc, override this function for custom model
func (rnr *Runner) Modeller(l logger.Logger) (modeller.Modeller, error) {

	// NOTE: A modeller is very service agnostic so there is no default generalised modeller we can provide here

	l.Debug("** Modeller **")

	return nil, nil
}
