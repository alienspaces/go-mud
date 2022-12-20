package server

import (
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

// Handle - custom service handle
type Handle func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error

// Runner - implements the runnerer interface
type Runner struct {
	Config            configurer.Configurer
	Log               logger.Logger
	Store             storer.Storer
	PrepareRepository preparer.Repository
	PrepareQuery      preparer.Query

	// configuration for routes, handlers and middleware
	HandlerConfig map[HandlerConfigKey]HandlerConfig
	MessageConfig map[Message]MessageConfig

	// composable functions
	RunHTTPFunc            func(args map[string]interface{}) error
	RunDaemonFunc          func(args map[string]interface{}) error
	RouterFunc             func(router *httprouter.Router) error
	MiddlewareFunc         func(h Handle) (Handle, error)
	HandlerFunc            Handle
	PreparerRepositoryFunc func(l logger.Logger) (preparer.Repository, error)
	PreparerQueryFunc      func(l logger.Logger) (preparer.Query, error)
	ModellerFunc           func(l logger.Logger) (modeller.Modeller, error)

	AuthenticateByAPIKeyFunc            func(m modeller.Modeller, l logger.Logger, apiKey string) (Authentication, error)
	GetAuthorizationsByHashedAPIKeyFunc func(modeller.Modeller, logger.Logger, Authentication) (Authentication, error)
	SetAuditConfigFunc                  func(m modeller.Modeller, l logger.Logger, requesterType string, requesterID string, requestID string) error
}

type Authentication struct {
	IsValid      bool
	HashedAPIKey string
	Roles        map[string]struct{}
	Permissions  map[string]struct{}
}

var _ runnable.Runnable = &Runner{}

type RequestPath string
type RequestMethod string
type HandlerConfigKey string

// HandlerConfig - configuration for routes, handlers and middleware
type HandlerConfig struct {
	Name string
	// Method - The HTTP method
	Method string
	// Path - The HTTP request URI including :parameter placeholders
	Path             string
	DocumentTagGroup DocumentTagGroupEndpoint
	// HandlerFunc - Function to handle requests for this method and path
	HandlerFunc Handle
	// MiddlewareConfig -
	MiddlewareConfig MiddlewareConfig
	// DocumentationConfig -
	DocumentationConfig DocumentationConfig
}

// ServerHandlerConfig provides API endpoint configuration
type ServerHandlerConfig map[RequestPath]map[RequestMethod]HandlerConfig

type DocumentTag string
type DocumentTagGroup string

// DocumentTagGroupEndpoint is used to group endpoints related to the same resource
type DocumentTagGroupEndpoint struct {
	ResourceName DocumentTagGroup `json:"name"`
	Description  string           `json:"description"`
	DocumentTags []DocumentTag    `json:"tags"`
}

type AuthenticationType string

const (
	AuthenTypePublic AuthenticationType = "Public"
	AuthenTypeAPIKey AuthenticationType = "Key"
	AuthenTypeJWT    AuthenticationType = "JWT"
)

type AuthorizationPermission string

// MiddlewareConfig - configuration for global default middleware
type MiddlewareConfig struct {
	AuthenTypes            []AuthenticationType
	AuthzPermissions       []AuthorizationPermission
	ValidateRequestSchema  jsonschema.SchemaWithReferences
	ValidateResponseSchema jsonschema.SchemaWithReferences
	// ValidateQueryParams - A whitelist of allowed query parameters
	ValidateQueryParams jsonschema.SchemaWithReferences
}

type QueryParams struct {
	Keys   []string
	Schema jsonschema.SchemaWithReferences
}

// DocumentationConfig - Configuration describing how to document a route
type DocumentationConfig struct {
	Document      bool
	Summary       string // used for API doc endpoint title
	Description   string // used for API doc endpoint description
	ErrorRegistry coreerror.Registry
}

type Message string
type MessageSource string
type MessageTopic string
type MessageSubject string
type MessageEvent string

type MessageConfig struct {
	Summary          string
	Name             Message
	Source           MessageSource
	Topic            MessageTopic
	Subject          MessageSubject
	Event            MessageEvent
	ValidateSchema   jsonschema.SchemaWithReferences
	DocumentTagGroup DocumentTagGroupSchemaModel
}

// DocumentTagGroupSchemaModel is used to group schema models related to the same resource
type DocumentTagGroupSchemaModel struct {
	ResourceName DocumentTagGroup `json:"name"`
	Description  string           `json:"description"`
	DocumentTag  DocumentTag      `json:"tag"`
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
func (rnr *Runner) Init(s storer.Storer) error {
	l := Logger(rnr.Log, "Init")

	if rnr.Log == nil {
		return fmt.Errorf("logger is nil, cannot initialise server runner")
	}

	l.Debug("** Initialise **")

	rnr.Store = s
	if rnr.Store == nil {
		msg := "storer is nil, cannot init runner"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	// Initialise storer
	err := rnr.Store.Init()
	if err != nil {
		l.Warn("failed store init >%v<", err)
		return err
	}

	// Repository
	if rnr.PreparerRepositoryFunc == nil {
		rnr.PreparerRepositoryFunc = rnr.PreparerRepository
	}

	pRepo, err := rnr.PreparerRepositoryFunc(l)
	if err != nil {
		l.Warn("failed preparer func >%v<", err)
		return err
	}

	rnr.PrepareRepository = pRepo
	if rnr.PrepareRepository == nil {
		l.Warn("PreparerRepository is nil, cannot continue")
		return err
	}

	if rnr.PreparerQueryFunc == nil {
		rnr.PreparerQueryFunc = rnr.PreparerQuery
	}

	pQ, err := rnr.PreparerQueryFunc(l)
	if err != nil {
		l.Warn("failed preparer config func >%v<", err)
		return err
	}

	rnr.PrepareQuery = pQ
	if rnr.PrepareQuery == nil {
		l.Warn("PreparerRepository Config is nil, cannot continue")
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

	// model
	if rnr.ModellerFunc == nil {
		rnr.ModellerFunc = rnr.Modeller
	}

	// HTTP server - router
	if rnr.RouterFunc == nil {
		rnr.RouterFunc = rnr.Router
	}

	// HTTP server - middleware
	if rnr.MiddlewareFunc == nil {
		rnr.MiddlewareFunc = rnr.DefaultMiddlewareFunc
	}

	// HTTP server - handler
	if rnr.HandlerFunc == nil {
		rnr.HandlerFunc = rnr.DefaultHandlerFunc
	}

	// Initialise configured routes
	root := rnr.Config.Get("APP_SERVER_HOME")
	for k, v := range rnr.HandlerConfig {
		l.Debug("Resolving schema location root >%s< key >%s<", root, k)
		v.MiddlewareConfig.ValidateQueryParams = jsonschema.ResolveSchemaLocationRoot(root, v.MiddlewareConfig.ValidateQueryParams)
		v.MiddlewareConfig.ValidateRequestSchema = jsonschema.ResolveSchemaLocationRoot(root, v.MiddlewareConfig.ValidateRequestSchema)
		v.MiddlewareConfig.ValidateResponseSchema = jsonschema.ResolveSchemaLocationRoot(root, v.MiddlewareConfig.ValidateResponseSchema)
		rnr.HandlerConfig[k] = v
	}

	return nil
}

// InitModeller initialises a new database transaction returning a prepared modeller
func (rnr *Runner) InitModeller(l logger.Logger) (modeller.Modeller, error) {

	// preparer
	if rnr.PreparerRepositoryFunc == nil {
		msg := "preparer function is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	p, err := rnr.PreparerRepositoryFunc(l)
	if err != nil {
		l.Warn("failed PreparerRepositoryFunc >%v<", err)
		return nil, err
	}

	if p == nil {
		msg := "preparer is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if rnr.PreparerQueryFunc == nil {
		msg := "preparer config function is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	pCfg, err := rnr.PreparerQueryFunc(l)
	if err != nil {
		l.Warn("failed PreparerQueryFunc >%v<", err)
		return nil, err
	}

	if pCfg == nil {
		msg := "preparer config is nil, cannot continue, cannot initialise database transaction"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// NOTE: The modeller is created and initialised with every request instead of
	// creating and assigning to a runner struct "Model" property at start up.
	// This prevents directly accessing a shared property from with the handler
	// function which is running in a goroutine. Otherwise accessing the "Model"
	// property would require locking and block simultaneous requests.

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

	err = m.Init(p, pCfg, tx)
	if err != nil {
		l.Warn("failed init modeller >%v<", err)
		return m, err
	}

	return m, nil
}

// Run starts the HTTP server and daemon processes. Override to implement a custom run function.
func (rnr *Runner) Run(args map[string]interface{}) error {

	rnr.Log.Info("** Run **")

	// signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// run HTTP server
	go func() {
		rnr.Log.Info("** Running HTTP server process **")
		if err := rnr.RunHTTPFunc(args); err != nil {
			rnr.Log.Error("Failed run server >%v<", err)
			sigChan <- syscall.SIGTERM
		}
		rnr.Log.Info("** HTTP server process ended **")
	}()

	// run daemon server
	go func() {
		rnr.Log.Info("** Running daemon process **")
		if err := rnr.RunDaemonFunc(args); err != nil {
			rnr.Log.Error("Failed run daemon >%v<", err)
			sigChan <- syscall.SIGTERM
		}
		rnr.Log.Info("** Daemon process ended **")
	}()

	// wait
	sig := <-sigChan

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	err := fmt.Errorf("received SIG >%v< Mem Alloc >%d MiB< TotalAlloc >%d MiB< Sys >%d MiB< NumGC >%d<",
		sig,
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC,
	)

	rnr.Log.Warn(">%v<", err)

	return err
}

// PreparerRepository - default PreparerRepositoryFunc, override this function for custom preparer.Repository
func (rnr *Runner) PreparerRepository(l logger.Logger) (preparer.Repository, error) {

	// NOTE: We have a good generic preparer.Repository so we'll provide that here

	l.Debug("** Repository **")

	// Return the existing preparer if we already have one
	if rnr.PrepareRepository != nil {
		l.Debug("Returning existing preparer")
		return rnr.PrepareRepository, nil
	}

	l.Debug("Creating new preparer")

	p, err := prepare.NewRepositoryPreparer(l)
	if err != nil {
		l.Warn("failed new prepare >%v<", err)
		return nil, err
	}

	db, err := rnr.Store.GetDb()
	if err != nil {
		l.Warn("failed getting database handle >%v<", err)
		return nil, err
	}

	err = p.Init(db)
	if err != nil {
		l.Warn("failed preparer init >%v<", err)
		return nil, err
	}

	rnr.PrepareRepository = p

	return p, nil
}

// PreparerQuery - default PreparerQueryFunc, override this function for custom preparer.Query
func (rnr *Runner) PreparerQuery(l logger.Logger) (preparer.Query, error) {

	// NOTE: We have a good generic preparer.Query so we'll provide that here

	l.Debug("** Query **")

	if rnr.PrepareQuery != nil {
		l.Debug("Returning existing preparer query")
		return rnr.PrepareQuery, nil
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

	rnr.PrepareQuery = p

	return p, nil
}

// Modeller - default ModellerFunc, override this function for custom model
func (rnr *Runner) Modeller(l logger.Logger) (modeller.Modeller, error) {

	// NOTE: A modeller is very service agnostic so there is no default generalised modeller we can provide here

	l.Debug("** Modeller **")

	return nil, nil
}

func (rnr *Runner) GetMessageConfigs() []MessageConfig {
	var cfgs []MessageConfig
	for _, cfg := range rnr.MessageConfig {
		cfgs = append(cfgs, cfg)
	}

	return cfgs
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

func ResolveHandlerSchemaLocation(handlerConfig map[HandlerConfigKey]HandlerConfig, location string) map[HandlerConfigKey]HandlerConfig {
	for handler, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.ValidateQueryParams.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateQueryParams = jsonschema.ResolveSchemaLocation(location, cfg.MiddlewareConfig.ValidateQueryParams)
		}

		if len(cfg.MiddlewareConfig.ValidateRequestSchema.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateRequestSchema = jsonschema.ResolveSchemaLocation(location, cfg.MiddlewareConfig.ValidateRequestSchema)
		}

		if len(cfg.MiddlewareConfig.ValidateResponseSchema.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateResponseSchema = jsonschema.ResolveSchemaLocation(location, cfg.MiddlewareConfig.ValidateResponseSchema)
		}

		handlerConfig[handler] = cfg
	}

	return handlerConfig
}

func ResolveHandlerSchemaLocationRoot(handlerConfig map[HandlerConfigKey]HandlerConfig, root string) (map[HandlerConfigKey]HandlerConfig, error) {
	for handler, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.ValidateQueryParams.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateQueryParams = jsonschema.ResolveSchemaLocationRoot(root, cfg.MiddlewareConfig.ValidateQueryParams)
		}

		if len(cfg.MiddlewareConfig.ValidateRequestSchema.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateRequestSchema = jsonschema.ResolveSchemaLocationRoot(root, cfg.MiddlewareConfig.ValidateRequestSchema)
		}

		if len(cfg.MiddlewareConfig.ValidateResponseSchema.Main.Name) > 0 {
			cfg.MiddlewareConfig.ValidateResponseSchema = jsonschema.ResolveSchemaLocationRoot(root, cfg.MiddlewareConfig.ValidateResponseSchema)
		}

		handlerConfig[handler] = cfg
	}

	return handlerConfig, nil
}

func ResolveMessageSchemaLocation(messageConfig map[Message]MessageConfig, location string) map[Message]MessageConfig {
	for message, cfg := range messageConfig {
		cfg.ValidateSchema = jsonschema.ResolveSchemaLocation(location, cfg.ValidateSchema)

		messageConfig[message] = cfg
	}

	return messageConfig
}

func ResolveMessageSchemaLocationRoot(messageConfig map[Message]MessageConfig, root string) (map[Message]MessageConfig, error) {
	for messsage, cfg := range messageConfig {
		cfg.ValidateSchema = jsonschema.ResolveSchemaLocationRoot(root, cfg.ValidateSchema)

		messageConfig[messsage] = cfg
	}

	return messageConfig, nil
}

func ResolveDocumentationSummary(handlerConfig map[HandlerConfigKey]HandlerConfig) map[HandlerConfigKey]HandlerConfig {
	for name, cfg := range handlerConfig {
		if cfg.DocumentationConfig.Summary == "" {
			cfg.DocumentationConfig.Summary = cfg.DocumentationConfig.Description
		}

		handlerConfig[name] = cfg
	}

	return handlerConfig
}

func ValidateAuthenticationTypes(handlerConfig map[HandlerConfigKey]HandlerConfig) error {
	for _, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.AuthenTypes) == 0 {
			return fmt.Errorf("handler >%s< with undefined authentication type", cfg.Name)
		}
	}
	return nil
}

func ToAuthorizationPermissionsSet(permissions ...AuthorizationPermission) map[AuthorizationPermission]struct{} {
	set := map[AuthorizationPermission]struct{}{}

	for _, p := range permissions {
		set[p] = struct{}{}
	}

	return set
}

func ToAuthenticationSet(authen ...AuthenticationType) map[AuthenticationType]struct{} {
	set := map[AuthenticationType]struct{}{}

	for _, p := range authen {
		set[p] = struct{}{}
	}

	return set
}

func InvalidAuthentication() Authentication {
	return Authentication{
		IsValid: false,
	}
}
