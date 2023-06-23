package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/prepare"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

const (
	// ConfigKeyValidateSchemaLocation - Directory location of JSON schema's
	ConfigKeyValidateSchemaLocation string = "validateSchemaLocation"
	// ConfigKeyValidateMainSchema - Main schema that can include reference schema's
	ConfigKeyValidateMainSchema string = "validateMainSchema"
	// ConfigKeyValidateReferenceSchemas - Schema referenced from the main schema
	ConfigKeyValidateReferenceSchemas string = "validateReferenceSchemas"
)

// Runner - implements the runnerer interface
type Runner struct {
	Config             configurer.Configurer
	Log                logger.Logger
	Store              storer.Storer
	RepositoryPreparer preparer.Repository
	QueryPreparer      preparer.Query

	// General configuration
	config Config

	// HTTPCORSConfig
	HTTPCORSConfig HTTPCORSConfig

	// Handler and message configuration
	HandlerConfig map[string]HandlerConfig

	// Run functions
	RunHTTPFunc   func(args map[string]interface{}) (*http.Server, error)
	RunDaemonFunc func(args map[string]interface{}) error

	// RouterFunc
	RouterFunc func(router *httprouter.Router) (*httprouter.Router, error)

	// HandlerFunc is the default handler function. It is used for liveness and healthz. Therefore, it should execute quickly.
	HandlerFunc Handle

	// HandlerMiddlewareFuncs returns a list of middleware functions to apply to routes
	HandlerMiddlewareFuncs func() []MiddlewareFunc

	// Service feature callbacks
	AuthenticateRequestFunc func(l logger.Logger, m modeller.Modeller, apiKey string) (AuthenticatedRequest, error)

	// Domain layer
	ModellerFunc func(l logger.Logger) (modeller.Modeller, error)

	// Data layer
	RepositoryPreparerFunc func(l logger.Logger) (preparer.Repository, error)
	QueryPreparerFunc      func(l logger.Logger) (preparer.Query, error)
}

type HTTPCORSConfig struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
}

var _ runnable.Runnable = &Runner{}

// Handle - custom service handle
type Handle func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error

// HandlerConfig - configuration for routes, handlers and middleware
type HandlerConfig struct {
	Name string
	// Method - The HTTP method
	Method string
	// Path - The HTTP request URI including :parameter placeholders
	Path string
	// HandlerFunc - Function to handle requests for this method and path
	HandlerFunc Handle
	// MiddlewareConfig -
	MiddlewareConfig MiddlewareConfig
	// DocumentationConfig -
	DocumentationConfig DocumentationConfig
}

type AuthenticatedType string

const (
	AuthenticatedTypeUser   AuthenticatedType = "User"
	AuthenticatedTypeAPIKey AuthenticatedType = "APIKey"
)

type AuthenticatedUser struct {
	ID    any    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthenticatedRequest struct {
	Type        AuthenticatedType      `json:"type"`
	User        AuthenticatedUser      `json:"user"`
	Permissions []AuthorizedPermission `json:"permissions"`
}

type AuthenticationType string
type AuthorizedPermission string

const (
	AuthenticationTypePublic     AuthenticationType = "Public"
	AuthenticationTypeRestricted AuthenticationType = "Restricted"
	AuthenticationTypeAPIKey     AuthenticationType = "API Key"
	AuthenticationTypeJWT        AuthenticationType = "JWT"
)

// MiddlewareConfig - configuration for global default middleware
type MiddlewareConfig struct {
	AuthenTypes            []AuthenticationType
	AuthzPermissions       []AuthorizedPermission
	ValidateRequestSchema  *jsonschema.SchemaWithReferences
	ValidateResponseSchema *jsonschema.SchemaWithReferences
	ValidateParamsConfig   *ValidateParamsConfig
}

// ValidateParamsConfig defines how route path parameters should be validated
//
// ExcludePathParamsFromQueryParams
// - By default path parameters will be added as query parameters and validated
// as part of query parameter validation. When disabled path parameters will
// need to be validated in handler functions
//
// RemoveParamPrefixes
// - TODO: Removes the specified prefixes from parameters
//
// PathParamSchema
// - Validate path parameters using this JSON schema set
//
// QueryParamSchema
// - Validate query parameters using this JSON schema set
type ValidateParamsConfig struct {
	ExcludePathParamsFromQueryParams bool
	RemoveParamPrefixes              []string
	PathParamSchema                  *jsonschema.SchemaWithReferences
	QueryParamSchema                 *jsonschema.SchemaWithReferences
}

// DocumentationConfig - Configuration describing how to document a route
type DocumentationConfig struct {
	Document    bool
	Summary     string // used for API doc endpoint title
	Description string // used for API doc endpoint description
}

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

func NewRunner(c configurer.Configurer, l logger.Logger) (*Runner, error) {

	r := Runner{
		Config: c,
		Log:    l,
	}

	return &r, nil
}

// Init - override to perform custom initialization
func (rnr *Runner) Init(s storer.Storer) error {
	l := Logger(rnr.Log, "Init")

	if rnr.Config != nil {
		cfg, err := NewConfig(rnr.Config)
		if err != nil {
			return err
		}
		rnr.config = *cfg
	}

	rnr.Store = s
	if rnr.Store != nil {
		if rnr.RepositoryPreparerFunc == nil {
			l.Info("Using default repository preparer func")
			rnr.RepositoryPreparerFunc = rnr.defaultRepositoryPreparerFunc
		}

		repoPreparer, err := rnr.RepositoryPreparerFunc(rnr.Log)
		if err != nil {
			l.Warn("Failed preparer func >%v<", err)
			return err
		}

		rnr.RepositoryPreparer = repoPreparer
		if rnr.RepositoryPreparer == nil {
			l.Warn("RepositoryPreparer is nil, cannot continue")
			return err
		}

		if rnr.QueryPreparerFunc == nil {
			l.Info("Using default query preparer func")
			rnr.QueryPreparerFunc = rnr.defaultQueryPreparerFunc
		}

		queryPreparer, err := rnr.QueryPreparerFunc(rnr.Log)
		if err != nil {
			l.Warn("Failed query preparer func >%v<", err)
			return err
		}

		rnr.QueryPreparer = queryPreparer
		if rnr.QueryPreparer == nil {
			l.Warn("QueryPreparer is nil, cannot continue")
			return err
		}
	}

	// run server
	if rnr.RunHTTPFunc == nil {
		rnr.RunHTTPFunc = rnr.RunHTTP
	}

	// run daemon
	if rnr.RunDaemonFunc == nil {
		rnr.RunDaemonFunc = rnr.RunDaemon
	}

	// model
	if rnr.ModellerFunc == nil {
		rnr.ModellerFunc = rnr.defaultModellerFunc
	}

	// http server - router
	if rnr.RouterFunc == nil {
		rnr.RouterFunc = rnr.defaultRouter
	}

	// http server - middleware
	if rnr.HandlerMiddlewareFuncs == nil {
		rnr.HandlerMiddlewareFuncs = rnr.defaultMiddlewareFuncs
	}

	// http server - handler
	if rnr.HandlerFunc == nil {
		rnr.HandlerFunc = rnr.defaultHandler
	}

	err := rnr.ResolveHandlerSchemaLocations()
	if err != nil {
		l.Warn("Failed resolving handler schema locations >%v<", err)
		return err
	}

	return nil
}

// InitTx initialises a new database transaction returning a prepared modeller
func (rnr *Runner) InitTx(l logger.Logger) (modeller.Modeller, error) {
	l = Logger(l, "InitTx")

	// preparer
	if rnr.RepositoryPreparerFunc == nil {
		msg := "preparer function is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	rp, err := rnr.RepositoryPreparerFunc(l)
	if err != nil {
		l.Warn("failed RepositoryPreparerFunc >%v<", err)
		return nil, err
	}

	if rp == nil {
		msg := "preparer is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if rnr.QueryPreparerFunc == nil {
		msg := "preparer config function is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	qp, err := rnr.QueryPreparerFunc(l)
	if err != nil {
		l.Warn("failed QueryPreparerFunc >%v<", err)
		return nil, err
	}

	if qp == nil {
		msg := "preparer config is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// NOTE: The modeller is created and initialised with every request instead of
	// creating and assigning to a runner struct "Model" property at start up.
	// This prevents directly accessing a shared property from with the handler
	// function which is running in a goroutine. Otherwise accessing the "Model"
	// property would require locking and blocking simultaneous requests.

	// modeller
	if rnr.ModellerFunc == nil {
		l.Warn("runner ModellerFunc is nil")
		return nil, fmt.Errorf("ModellerFunc is nil")
	}

	m, err := rnr.ModellerFunc(l)
	if err != nil {
		l.Warn("failed ModellerFunc >%v<", err)
		return nil, err
	}
	if m == nil {
		l.Warn("modeller is nil, cannot continue")
		return nil, err
	}

	tx, err := rnr.Store.GetTx()
	if err != nil {
		l.Warn("failed getting DB connection >%v<", err)
		return m, err
	}

	err = m.Init(rp, qp, tx)
	if err != nil {
		l.Warn("failed init modeller >%v<", err)
		return m, err
	}

	return m, nil
}

// Run starts the HTTP server and daemon processes. Override to implement a custom run function.
func (rnr *Runner) Run(args map[string]interface{}) error {

	rnr.Log.Debug("** Run **")

	// signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var srv *http.Server
	var err error

	// run HTTP server
	go func() {
		rnr.Log.Debug("** Running HTTP server process **")
		srv, err = rnr.RunHTTPFunc(args)
		if err != nil {
			rnr.Log.Error("Failed run server >%v<", err)
			sigChan <- syscall.SIGTERM
		}
		rnr.Log.Debug("** HTTP server process ended **")
	}()

	// run daemon server
	go func() {
		rnr.Log.Debug("** Running daemon process **")
		if err := rnr.RunDaemonFunc(args); err != nil {
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

	err = fmt.Errorf("received SIG >%v< Mem Alloc >%d MiB< TotalAlloc >%d MiB< Sys >%d MiB< NumGC >%d<",
		sig,
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC,
	)

	rnr.Log.Warn(">%v<", err)

	if srv != nil {
		// By default, a k8s pod is given a termination grace period of 30s:
		// https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase
		// We specify a read (e.g., 5s) and write (e.g., 10s) timeout of less than 30s
		// to ensure requests will be terminated in that time.
		if err := srv.Shutdown(context.Background()); err != nil {
			rnr.Log.Error("failed shutting down server >%#v<", err)
		}
	}

	return err
}

// defaultRepositoryPreparerFunc - returns a default initialised repository preparer, set the
// property RepositoryPreparerFunc to provide your own custom repository preparer.
func (rnr *Runner) defaultRepositoryPreparerFunc(l logger.Logger) (preparer.Repository, error) {
	l = Logger(l, "defaultRepositoryPreparerFunc")

	// Return the existing preparer if we already have one
	if rnr.RepositoryPreparer != nil {
		l.Info("Returning existing repository preparer")
		return rnr.RepositoryPreparer, nil
	}

	l.Info("Creating new repository preparer")

	p, err := prepare.NewRepositoryPreparer(l)
	if err != nil {
		l.Warn("Failed new repository prepare >%v<", err)
		return nil, err
	}

	db, err := rnr.Store.GetDb()
	if err != nil {
		l.Warn("Failed getting database handle >%v<", err)
		return nil, err
	}

	err = p.Init(db)
	if err != nil {
		l.Warn("Failed repository preparer init >%v<", err)
		return nil, err
	}

	rnr.RepositoryPreparer = p

	return p, nil
}

// defaultQueryPreparerFunc - returns a default initialised query preparer, set the
// property QueryPreparerFunc to provide your own custom query preparer.
func (rnr *Runner) defaultQueryPreparerFunc(l logger.Logger) (preparer.Query, error) {
	l = Logger(l, "defaultQueryPreparerFunc")

	if rnr.QueryPreparer != nil {
		l.Info("Returning existing query preparer")
		return rnr.QueryPreparer, nil
	}

	l.Debug("Creating new query preparer")

	p, err := prepare.NewQueryPreparer(l)
	if err != nil {
		l.Warn("failed new query preparer >%v<", err)
		return nil, err
	}

	db, err := rnr.Store.GetDb()
	if err != nil {
		l.Warn("failed getting database handle >%v<", err)
		return nil, err
	}

	err = p.Init(db)
	if err != nil {
		l.Warn("failed query preparer init >%v<", err)
		return nil, err
	}

	rnr.QueryPreparer = p

	return p, nil
}

// defaultModellerFunc does not provide a modeller, set the property ModellerFunc to
// provide your own custom modeller.
func (rnr *Runner) defaultModellerFunc(l logger.Logger) (modeller.Modeller, error) {

	// NOTE: A modeller is service specific so there is no default modeller we can provide here

	l.Debug("** Modeller **")

	return nil, nil
}
