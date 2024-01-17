package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

type HttpMethod string

const (
	HttpMethodGet     HttpMethod = http.MethodGet
	HttpMethodHead    HttpMethod = http.MethodHead
	HttpMethodPost    HttpMethod = http.MethodPost
	HttpMethodPut     HttpMethod = http.MethodPut
	HttpMethodPatch   HttpMethod = http.MethodPatch
	HttpMethodDelete  HttpMethod = http.MethodDelete
	HttpMethodConnect HttpMethod = http.MethodConnect
	HttpMethodOptions HttpMethod = http.MethodOptions
	HttpMethodTrace   HttpMethod = http.MethodTrace
)

type WriteResponseOption = func(http.ResponseWriter) error

// RunHTTP - Starts the HTTP server process. Override to implement a custom HTTP server run function.
// The server process exposes a REST API and is intended for clients to manage resources and
// perform actions.
func (rnr *Runner) RunHTTP(args map[string]interface{}) (*http.Server, error) {

	rnr.Log.Debug("** RunHTTP **")

	// Router
	r := httprouter.New()

	r, err := rnr.registerRoutes(r)
	if err != nil {
		rnr.Log.Warn("failed registering routes >%v<", err)
		return nil, err
	}

	// CORS
	allowedOrigins := rnr.HTTPCORSConfig.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	allowedHeaders := []string{
		"X-ProgramID", "X-ProgramName", "Content-Type",
		"Authorization", "X-Authorization-Token",
		"Origin", "X-Requested-With", "Accept",
		"X-CSRF-Token",
	}

	allowedHeaders = append(allowedHeaders, rnr.HTTPCORSConfig.AllowedHeaders...)

	// Access-Control-Allow-Origin, Access-Control-Allow-Headers and Access-Control-Allow-Methods
	// cannot be wildcard if the CORS request is credentialed.
	c := cors.New(cors.Options{
		Debug:            false,
		AllowedOrigins:   allowedOrigins,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   rnr.HTTPCORSConfig.ExposedHeaders,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
		AllowCredentials: rnr.HTTPCORSConfig.AllowCredentials,
	})
	h := c.Handler(r)

	// Serve
	port := rnr.config.AppServerPort
	if port == "" {
		rnr.Log.Warn("missing APP_SERVER_PORT, cannot start server")
		return nil, fmt.Errorf("missing APP_SERVER_PORT, cannot start server")
	}

	rnr.Log.Info("server running at: http://localhost:%s", port)

	srv := &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return srv, srv.ListenAndServe()
}

// registerRoutes - registers routes as implemented by the assigned router function
func (rnr *Runner) registerRoutes(r *httprouter.Router) (*httprouter.Router, error) {
	return rnr.RouterFunc(r)
}

// defaultHandler is the default HandlerFunc
func (rnr *Runner) defaultHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Handler **")

	fmt.Fprint(w, "Ok!\n")

	return nil
}

func (rnr *Runner) RegisterDefaultHealthzRoute(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "RegisterDefaultHealthzRoute")

	h, err := rnr.ApplyMiddleware(
		HandlerConfig{
			Path: "/healthz",
			MiddlewareConfig: MiddlewareConfig{
				AuthenTypes: []AuthenticationType{AuthenticationTypePublic},
			},
		},
		rnr.HandlerFunc,
	)
	if err != nil {
		l.Warn("failed default middleware >%v<", err)
		return nil, err
	}
	r.GET("/healthz", h)

	l.Info("Registered /healthz")

	return r, nil
}

func (rnr *Runner) RegisterDefaultLivenessRoute(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "RegisterDefaultLivenessRoute")

	// This logger should only be used for the liveness endpoint and is to avoid creating
	// a new logger on every request.
	hl := rnr.Log.NewInstance(nil)
	r.GET("/liveness", func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {
		_ = rnr.HandlerFunc(w, r, pp, nil, hl, nil)
	})

	l.Info("Registered /liveness")

	return r, nil
}

// defaultRouter - implements default routes based on runner configuration options
func (rnr *Runner) defaultRouter(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "defaultRouter")

	var err error
	r, err = rnr.RegisterDefaultHealthzRoute(r)
	if err != nil {
		return nil, err
	}

	r, err = rnr.RegisterDefaultLivenessRoute(r)
	if err != nil {
		return nil, err
	}

	// register configured routes
	for _, hc := range rnr.HandlerConfig {

		l.Info("Router method >%s< path >%s<", hc.Method, hc.Path)

		h, err := rnr.ApplyMiddleware(hc, hc.HandlerFunc)
		if err != nil {
			l.Warn("failed registering handler >%v<", err)
			return nil, err
		}

		switch hc.Method {
		case http.MethodGet:
			r.GET(hc.Path, h)
		case http.MethodPost:
			r.POST(hc.Path, h)
		case http.MethodPut:
			r.PUT(hc.Path, h)
		case http.MethodPatch:
			r.PATCH(hc.Path, h)
		case http.MethodDelete:
			r.DELETE(hc.Path, h)
		case http.MethodOptions:
			r.OPTIONS(hc.Path, h)
		case http.MethodHead:
			r.HEAD(hc.Path, h)
		default:
			l.Warn("router HTTP method >%s< not supported", hc.Method)
			return nil, fmt.Errorf("router HTTP method >%s< not supported", hc.Method)
		}
	}

	return r, nil
}

// HttpRouterHandlerWrapper wraps a Handle function in an httprouter.Handle function while also
// providing a new logger for every request. Typically this function should be used to wrap the
// final product of applying all middleware to Handle function.
func (rnr *Runner) HttpRouterHandlerWrapper(h Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {
		// Create new logger with its own context (fields map) on every request because each
		// request maintains its own context (fields map). If the same logger is used, when
		// different requests set the logger context, there will be concurrent map read/writes.
		l := rnr.Log.NewInstance(nil)

		// delegate
		_ = h(w, r, pp, nil, l, nil)
	}
}

func (rnr *Runner) ResolveHandlerSchemaLocations() error {
	l := Logger(rnr.Log, "resolveHandlerSchemaLocationRoot")

	appHome := rnr.config.AppServerHome
	if appHome == "" {
		err := fmt.Errorf("missing app home")
		l.Warn(err.Error())
		return err
	}

	handlerConfig := rnr.HandlerConfig

	for name, cfg := range handlerConfig {

		if cfg.MiddlewareConfig.ValidateParamsConfig != nil {

			// Path parameter schemas
			if cfg.MiddlewareConfig.ValidateParamsConfig.PathParamSchema != nil {
				schema := cfg.MiddlewareConfig.ValidateParamsConfig.PathParamSchema
				if schema.Main.Name == "" {
					return fmt.Errorf("handler name >%s< main path params schema missing name", name)
				}
				if schema.Main.Location == "" {
					return fmt.Errorf("handler name >%s< main path params schema missing location", name)
				}
				if schema.Main.LocationRoot == "" {
					schema = jsonschema.ResolveSchemaLocationRoot(appHome, schema)
				}

				l.Info("PathParamSchema Main Name >%s<", schema.Main.Name)
				l.Info("PathParamSchema Main LocationRoot >%s<", schema.Main.LocationRoot)
				l.Info("PathParamSchema Main Location >%s<", schema.Main.Location)

				for i := range schema.References {
					schemaRef := schema.References[i]
					l.Info("PathParamSchema Ref Name >%s<", schemaRef.Name)
					l.Info("PathParamSchema Ref LocationRoot >%s<", schemaRef.LocationRoot)
					l.Info("PathParamSchema Ref Location >%s<", schemaRef.Location)
				}

				cfg.MiddlewareConfig.ValidateParamsConfig.PathParamSchema = schema
			}

			// Query parameter schemas
			if cfg.MiddlewareConfig.ValidateParamsConfig.QueryParamSchema != nil {
				schema := cfg.MiddlewareConfig.ValidateParamsConfig.QueryParamSchema
				if schema.Main.Name == "" {
					return fmt.Errorf("handler name >%s< main query params schema missing name", name)
				}
				if schema.Main.Location == "" {
					return fmt.Errorf("handler name >%s< main query params schema missing location", name)
				}
				if schema.Main.LocationRoot == "" {
					schema = jsonschema.ResolveSchemaLocationRoot(appHome, schema)
				}

				l.Info("QueryParamSchema Main Name >%s<", schema.Main.Name)
				l.Info("QueryParamSchema Main LocationRoot >%s<", schema.Main.LocationRoot)
				l.Info("QueryParamSchema Main Location >%s<", schema.Main.Location)

				for i := range schema.References {
					schemaRef := schema.References[i]
					l.Info("QueryParamSchema Ref Name >%s<", schemaRef.Name)
					l.Info("QueryParamSchema Ref LocationRoot >%s<", schemaRef.LocationRoot)
					l.Info("QueryParamSchema Ref Location >%s<", schemaRef.Location)
				}

				cfg.MiddlewareConfig.ValidateParamsConfig.QueryParamSchema = schema
			}
		}

		if cfg.MiddlewareConfig.ValidateRequestSchema != nil {

			schema := cfg.MiddlewareConfig.ValidateRequestSchema
			if schema.Main.Name == "" {
				return fmt.Errorf("handler name >%s< main request schema missing name", name)
			}
			if schema.Main.Location == "" {
				return fmt.Errorf("handler name >%s< main request schema missing location", name)
			}

			if schema.Main.LocationRoot == "" {
				schema = jsonschema.ResolveSchemaLocationRoot(appHome, schema)
			}

			cfg.MiddlewareConfig.ValidateRequestSchema = schema
		}

		if cfg.MiddlewareConfig.ValidateResponseSchema != nil {

			schema := cfg.MiddlewareConfig.ValidateResponseSchema
			if schema.Main.Name == "" {
				return fmt.Errorf("handler name >%s< main response schema missing name", name)
			}
			if schema.Main.Location == "" {
				return fmt.Errorf("handler name >%s< main response schema missing location", name)
			}

			if schema.Main.LocationRoot == "" {
				schema = jsonschema.ResolveSchemaLocationRoot(appHome, schema)
			}

			cfg.MiddlewareConfig.ValidateResponseSchema = schema
		}

		handlerConfig[name] = cfg
	}

	rnr.HandlerConfig = handlerConfig

	return nil
}

func ValidateAuthenticationTypes(handlerConfig map[string]HandlerConfig) error {
	for _, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.AuthenTypes) == 0 {
			return fmt.Errorf("handler >%s< with undefined authentication type", cfg.Name)
		}
	}

	return nil
}

// RequestData -
func RequestData(l logger.Logger, r *http.Request) *string {
	value := r.Context().Value(ctxKeyData)
	data, ok := value.(string)
	if !ok {
		return nil
	}

	l.Info("Request data >%s<", data)

	return &data
}

func AuthData(l logger.Logger, r *http.Request) *AuthenticatedRequest {
	auth, ok := (r.Context().Value(ctxKeyAuth)).(AuthenticatedRequest)
	if !ok {
		return nil
	}

	l.Info("Auth data Type >%s< Permissions >%v<", auth.Type, auth.Permissions)

	return &auth
}

func GetCorrelationID(l logger.Logger, r *http.Request) (string, error) {
	correlationID, ok := (r.Context().Value(ctxKeyCorrelationID)).(string)
	if !ok {
		return "", fmt.Errorf("missing correlation ID")
	}

	l.Info("Correlation ID >%s<", correlationID)

	return correlationID, nil
}

// ReadRequest -
func ReadRequest[T any](l logger.Logger, r *http.Request, s *T) (*T, error) {

	data := RequestData(l, r)
	if data == nil {
		return nil, nil
	}

	reader := strings.NewReader(*data)
	err := json.NewDecoder(reader).Decode(s)
	if err != nil {
		return nil, fmt.Errorf("failed decoding request data >%s< >%v<", *data, err)
	}

	return s, nil
}

// ReadXMLRequest -
func ReadXMLRequest(l logger.Logger, r *http.Request, s interface{}) (*string, error) {

	data := RequestData(l, r)
	if data == nil {
		return nil, nil
	}

	reader := strings.NewReader(*data)
	if err := xml.NewDecoder(reader).Decode(s); err != nil {
		return data, fmt.Errorf("failed decoding request data >%s< >%v<", *data, err)
	}

	return data, nil
}

// WriteResponse -
func WriteResponse(l logger.Logger, w http.ResponseWriter, status int, r interface{}, options ...WriteResponseOption) error {
	l.Info("write response status >%d<", status)

	// content type json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	for _, o := range options {
		if err := o(w); err != nil {
			return err
		}
	}

	// status
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(r)
}

func WritePaginatedResponse[R any, D any](l logger.Logger, w http.ResponseWriter, recs []R, mapper func(R) (D, error), pageSize int) error {
	res := []D{}

	buildPageSize := pageSize
	for _, rec := range recs {
		if buildPageSize == 0 {
			break
		}

		responseData, err := mapper(rec)
		if err != nil {
			WriteSystemError(l, w, err)
			return err
		}
		res = append(res, responseData)

		buildPageSize--
	}

	err := WriteResponse(l, w, http.StatusOK, res, XPaginationHeader(len(recs), pageSize))
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func WriteXMLResponse(l logger.Logger, w http.ResponseWriter, s interface{}) error {
	status := http.StatusOK
	l.Info("write response status >%d<", status)

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	w.WriteHeader(status)

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}

	return xml.NewEncoder(w).Encode(s)
}
