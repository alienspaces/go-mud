package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"syscall"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
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

	// Assignable functions
	RunHTTPFunc   func(args map[string]interface{}) (*http.Server, error)
	RunDaemonFunc func(args map[string]interface{}) error
	RouterFunc    func(router *httprouter.Router) (*httprouter.Router, error)

	// HandlerFunc is the default handler function. It is used for liveness and healthz. Therefore, it should execute quickly.
	HandlerFunc Handle

	// HandlerMiddlewareFuncs returns a list of middleware functions to apply to routes
	HandlerMiddlewareFuncs func() []MiddlewareFunc
	RepositoryPreparerFunc func(l logger.Logger) (preparer.Repository, error)
	QueryPreparerFunc      func(l logger.Logger) (preparer.Query, error)
	ModellerFunc           func(l logger.Logger) (modeller.Modeller, error)

	AuthenticateRequestFunc func(l logger.Logger, m modeller.Modeller, apiKey string) (AuthenticatedRequest, error)
}

type AuditRequest struct {
	RequestID string

	RequesterType  string
	RequesterID    string
	RequesterName  string
	RequesterEmail string
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

type Tag string
type TagGroup string

// TagGroupEndpoint is used to group endpoints related to the same resource
type TagGroupEndpoint struct {
	ResourceName TagGroup `json:"name"`
	Description  string   `json:"description"`
	Tags         []Tag    `json:"tags"`
}

// NOTE: AuthenticatedRequest modelled from the following for possible familliarity.
// https://gitlab.com/msts-enterprise/rock/caas-customer/-/blob/develop/server/src/core/authentication/schemas/x-authenticated-request.schema.json
type AuthenticatedRequest struct {
	Type        AuthenticatedType      `json:"type"`
	User        AuthenticatedUser      `json:"user"`
	Permissions []AuthorizedPermission `json:"permissions"`
	RLSType     RLSType                `json:"-"`
}

type RLSType string

const (
	RLSTypeOpen       RLSType = "open"
	RLSTypeRestricted RLSType = "restricted"
)

type AuthenticatedType string

const (
	AuthenticatedTypeUser   AuthenticatedType = "User"
	AuthenticatedTypeAPIKey AuthenticatedType = "APIKey"
)

type AuthenticatedUser struct {
	ID                  any    `json:"id"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	ServiceCloudProfile string `json:"service_cloud_profile,omitempty"` // Profile is used to support QA testing by switching the Active SC Profile in the UI. It should not be populated in Staging or Production.
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
	ValidateRequestSchema  jsonschema.SchemaWithReferences
	ValidateResponseSchema jsonschema.SchemaWithReferences
	ValidateQueryParams    jsonschema.SchemaWithReferences
}

// DocumentationConfig - Configuration describing how to document a route
type DocumentationConfig struct {
	Document      bool
	Summary       string // used for API doc endpoint title
	Description   string // used for API doc endpoint description
	ErrorRegistry coreerror.Registry
	TagGroup      TagGroupEndpoint
}

type MessageAttribute struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type MessageConfig struct {
	Summary           string
	Name              string
	Source            string
	Topic             string
	Subject           string
	Event             string
	attributesMapping map[string]MessageAttribute
	Attributes        []MessageAttribute
	ValidateSchema    jsonschema.SchemaWithReferences
	TagGroup          TagGroupSchemaModel
}

func (m MessageConfig) AttributesMap() map[string]MessageAttribute {
	if m.attributesMapping != nil {
		return m.attributesMapping
	}

	attributes := map[string]MessageAttribute{}

	for _, a := range m.Attributes {
		attributes[a.Name] = a
	}

	m.attributesMapping = attributes

	return m.attributesMapping
}

// TagGroupSchemaModel is used to group schema models related to the same resource
type TagGroupSchemaModel struct {
	ResourceName TagGroup `json:"name"`
	Description  string   `json:"description"`
	Tag          Tag      `json:"tag"`
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

	cfg, err := NewConfig(c)
	if err != nil {
		return nil, err
	}

	r := Runner{
		Config: c,
		Log:    l,
		config: *cfg,
	}

	return &r, nil
}

// Init - override to perform custom initialization
func (rnr *Runner) Init(s storer.Storer) error {

	rnr.Log.Debug("** Initialise **")

	rnr.Store = s
	if rnr.Store != nil {
		if rnr.RepositoryPreparerFunc == nil {
			rnr.RepositoryPreparerFunc = rnr.defaultRepositoryPreparerFunc
		}

		repoPreparer, err := rnr.RepositoryPreparerFunc(rnr.Log)
		if err != nil {
			rnr.Log.Warn("Failed preparer func >%v<", err)
			return err
		}

		rnr.RepositoryPreparer = repoPreparer
		if rnr.RepositoryPreparer == nil {
			rnr.Log.Warn("RepositoryPreparer is nil, cannot continue")
			return err
		}

		if rnr.QueryPreparerFunc == nil {
			rnr.QueryPreparerFunc = rnr.defaultQueryPreparerFunc
		}

		queryPreparer, err := rnr.QueryPreparerFunc(rnr.Log)
		if err != nil {
			rnr.Log.Warn("Failed query preparer func >%v<", err)
			return err
		}

		rnr.QueryPreparer = queryPreparer
		if rnr.QueryPreparer == nil {
			rnr.Log.Warn("QueryPreparer is nil, cannot continue")
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

	// NOTE: We have a good generic preparer.Repository so we'll provide that here

	l.Debug("** Repository **")

	// Return the existing preparer if we already have one
	if rnr.RepositoryPreparer != nil {
		l.Debug("Returning existing preparer")
		return rnr.RepositoryPreparer, nil
	}

	l.Debug("Creating new preparer")

	p, err := prepare.NewRepositoryPreparer(l)
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

	rnr.RepositoryPreparer = p

	return p, nil
}

// defaultQueryPreparerFunc - returns a default initialised query preparer, set the
// property QueryPreparerFunc to provide your own custom query preparer.
func (rnr *Runner) defaultQueryPreparerFunc(l logger.Logger) (preparer.Query, error) {

	// NOTE: We have a good generic preparer.Query so we'll provide that here

	l.Debug("** Query **")

	if rnr.QueryPreparer != nil {
		l.Debug("Returning existing preparer query")
		return rnr.QueryPreparer, nil
	}

	l.Debug("Creating new preparer query")

	p, err := prepare.NewQueryPreparer(l)
	if err != nil {
		l.Warn("failed new prepare query >%v<", err)
		return nil, err
	}

	db, err := rnr.Store.GetDb()
	if err != nil {
		l.Warn("failed getting database handle >%v<", err)
		return nil, err
	}

	err = p.Init(db)
	if err != nil {
		l.Warn("failed preparer query init >%v<", err)
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

func (rnr *Runner) GetHandlerConfigs() []HandlerConfig {
	var cfgs []HandlerConfig
	for _, cfg := range rnr.HandlerConfig {
		cfgs = append(cfgs, cfg)
	}

	return sortHandlerConfigs(cfgs)
}

func sortHandlerConfigs(handlerConfigs []HandlerConfig) []HandlerConfig {
	var sorted []HandlerConfig
	sorted = append(sorted, handlerConfigs...)

	sort.Slice(sorted, func(i, j int) bool {
		x := sorted[i]
		y := sorted[j]

		if x.Path != y.Path {
			return x.Path < y.Path
		}

		return x.Method < y.Method
	})

	return sorted
}

func ResolveHandlerSchemaLocation(handlerConfig map[string]HandlerConfig, location string) map[string]HandlerConfig {
	for handler, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.ValidateQueryParams.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateQueryParams = jsonschema.ResolveSchemaLocation(cfg.MiddlewareConfig.ValidateQueryParams, location)
		}

		if len(cfg.MiddlewareConfig.ValidateRequestSchema.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateRequestSchema = jsonschema.ResolveSchemaLocation(cfg.MiddlewareConfig.ValidateRequestSchema, location)
		}

		if len(cfg.MiddlewareConfig.ValidateResponseSchema.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateResponseSchema = jsonschema.ResolveSchemaLocation(cfg.MiddlewareConfig.ValidateResponseSchema, location)
		}

		handlerConfig[handler] = cfg
	}

	return handlerConfig
}

func ResolveHandlerSchemaLocationRoot(handlerConfig map[string]HandlerConfig, root string) (map[string]HandlerConfig, error) {
	for handler, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.ValidateQueryParams.Main.Name) > 0 || len(cfg.MiddlewareConfig.ValidateQueryParams.Main.Location) > 0 {
			cfg.MiddlewareConfig.ValidateQueryParams = jsonschema.ResolveSchemaLocationRoot(cfg.MiddlewareConfig.ValidateQueryParams, root)
		}

		if len(cfg.MiddlewareConfig.ValidateRequestSchema.Main.Name) > 0 || len(cfg.MiddlewareConfig.ValidateRequestSchema.Main.Location) > 0 {
			cfg.MiddlewareConfig.ValidateRequestSchema = jsonschema.ResolveSchemaLocationRoot(cfg.MiddlewareConfig.ValidateRequestSchema, root)
		}

		if len(cfg.MiddlewareConfig.ValidateResponseSchema.Main.Name) > 0 || len(cfg.MiddlewareConfig.ValidateResponseSchema.Main.Location) > 0 {
			cfg.MiddlewareConfig.ValidateResponseSchema = jsonschema.ResolveSchemaLocationRoot(cfg.MiddlewareConfig.ValidateResponseSchema, root)
		}

		handlerConfig[handler] = cfg
	}

	return handlerConfig, nil
}

func ResolveDocumentationSummary(handlerConfig map[string]HandlerConfig) map[string]HandlerConfig {
	for name, cfg := range handlerConfig {
		if cfg.DocumentationConfig.Summary == "" {
			cfg.DocumentationConfig.Summary = cfg.DocumentationConfig.Description
		}
		handlerConfig[name] = cfg
	}

	return handlerConfig
}

func ValidateAuthenticationTypes(handlerConfig map[string]HandlerConfig) error {
	for _, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.AuthenTypes) == 0 {
			return fmt.Errorf("handler >%s< with undefined authentication type", cfg.Name)
		}
	}

	return nil
}
