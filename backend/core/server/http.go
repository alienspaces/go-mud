package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

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
func (rnr *Runner) RunHTTP(args map[string]interface{}) error {

	rnr.Log.Debug("** RunHTTP **")

	// default handler
	router, err := rnr.DefaultRouter()
	if err != nil {
		rnr.Log.Warn("failed default router >%v<", err)
		return err
	}

	port := rnr.Config.Get("APP_SERVER_PORT")
	if port == "" {
		rnr.Log.Warn("missing APP_SERVER_PORT, cannot start server")
		return fmt.Errorf("missing APP_SERVER_PORT, cannot start server")
	}

	// cors
	c := cors.New(cors.Options{
		Debug:          false,
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{
			"X-ProgramID", "X-ProgramName", "Content-Type",
			"X-Authorization", "X-Authorization-Token",
			"Origin", "X-Requested-With", "Accept",
			"Access-Control-Allow-Origin",
			"X-CSRF-Token",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	// serve
	rnr.Log.Info("server running at: http://localhost:%s", port)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}

// DefaultRouterFunc - default RouterFunc, override this function for custom routes
func (rnr *Runner) DefaultRouterFunc(router *httprouter.Router) error {
	l := loggerWithContext(rnr.Log, "Router")

	l.Info("Using default empty router")

	return nil
}

// Middleware - default MiddlewareFunc, override this function for custom middleware
func (rnr *Runner) DefaultMiddlewareFunc(h Handle) (Handle, error) {
	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

// Handler - default HandlerFunc, override this function for custom handler
func (rnr *Runner) DefaultHandlerFunc(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "DefaultHandlerFunc")

	l.Info("Using default handler")

	fmt.Fprint(w, "Ok!\n")

	return nil
}

// DefaultRouter - implements default routes based on runner configuration options
func (rnr *Runner) DefaultRouter() (*httprouter.Router, error) {
	l := loggerWithContext(rnr.Log, "DefaultRouter")

	l.Info("Using default router")

	r := httprouter.New()

	// Default index handler
	h, err := rnr.DefaultMiddleware(HandlerConfig{Path: "/", MiddlewareConfig: MiddlewareConfig{
		AuthenTypes: []AuthenticationType{AuthenTypePublic},
	}}, rnr.HandlerFunc)
	if err != nil {
		rnr.Log.Warn("failed default middleware >%v<", err)
		return nil, err
	}
	r.GET("/healthz", h)
	r.GET("/liveness", func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {
		_ = rnr.HandlerFunc(w, r, pp, nil, rnr.Log, nil)
	})

	// Register configured routes
	for _, hc := range rnr.HandlerConfig {

		l.Info("Routing method >%s< path >%s<", hc.Method, hc.Path)

		h, err := rnr.DefaultMiddleware(hc, hc.HandlerFunc)
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
			err := fmt.Errorf("router HTTP method >%s< not supported", hc.Method)
			l.Warn(err.Error())
			return nil, err
		}
	}

	// server defined routes
	err = rnr.RouterFunc(r)
	if err != nil {
		l.Warn("failed router >%v<", err)
		return nil, err
	}

	return r, nil
}

// DefaultMiddleware - implements middlewares based on runner configuration
func (rnr *Runner) DefaultMiddleware(hc HandlerConfig, h Handle) (httprouter.Handle, error) {
	l := loggerWithContext(rnr.Log, "DefaultMiddleware")

	l.Info("Using default middleware")

	// NOTE: The order matters here. Different middleware will not be able to function
	// unless prior middleware have run successfully.

	// Validate body data
	h, err := rnr.Validate(hc, h)
	if err != nil {
		l.Warn("failed adding validate middleware >%v<", err)
		return nil, err
	}

	// Request body data
	h, err = rnr.Data(h)
	if err != nil {
		l.Warn("failed adding data middleware >%v<", err)
		return nil, err
	}

	h, err = rnr.Audit(hc, h)
	if err != nil {
		l.Warn("failed adding audit middleware >%v<", err)
		return nil, err
	}

	// Authz
	h, err = rnr.Authz(hc, h)
	if err != nil {
		l.Warn("failed adding authz middleware >%v<", err)
		return nil, err
	}

	// Authen
	h, err = rnr.Authen(hc, h)
	if err != nil {
		l.Warn("failed adding authen middleware >%v<", err)
		return nil, err
	}

	// Tx
	h, err = rnr.Tx(h)
	if err != nil {
		l.Warn("failed adding tx middleware >%v<", err)
		return nil, err
	}

	// Correlation
	h, err = rnr.Correlation(h)
	if err != nil {
		l.Warn("failed adding correlation middleware >%v<", err)
		return nil, err
	}

	// Server defined routes
	h, err = rnr.MiddlewareFunc(h)
	if err != nil {
		l.Warn("failed middleware >%v<", err)
		return nil, err
	}

	// Wrap everything in a httprouter Handler
	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {

		// New logger instance per request
		l, err := l.NewInstance()
		if err != nil {
			WriteSystemError(rnr.Log, w, err)
			return
		}

		// Delegate
		h(w, r, pp, nil, l, nil)
	}

	return handle, nil
}

// ReadRequest -
func ReadRequest(l logger.Logger, r *http.Request, s interface{}) error {

	data := r.Context().Value(ctxKeyData)

	if data != nil {
		r := strings.NewReader(data.(string))
		err := json.NewDecoder(r).Decode(s)
		if err != nil {
			// Include data in error response
			return fmt.Errorf("failed decoding request data >%s< >%v<", data.(string), err)
		}
	}

	return nil
}

func ReadXMLRequest(l logger.Logger, r *http.Request, s interface{}) (*string, error) {

	data := r.Context().Value(ctxKeyData)

	d, ok := data.(string)
	if !ok {
		return nil, nil
	}

	l.Info("xml message body >%s<", d)
	reader := strings.NewReader(d)
	if err := xml.NewDecoder(reader).Decode(s); err != nil {
		return &d, fmt.Errorf("failed decoding request data >%s< >%v<", d, err)
	}

	return &d, nil
}

// WriteResponse -
func WriteResponse(l logger.Logger, w http.ResponseWriter, r interface{}, options ...WriteResponseOption) error {
	l = loggerWithContext(l, "WriteResponse")

	status := http.StatusOK
	l.Info("Write response status >%d<", status)

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

func WriteXMLResponse(l logger.Logger, w http.ResponseWriter, s interface{}) error {
	l = loggerWithContext(l, "WriteXMLResponse")

	status := http.StatusOK
	l.Info("Write response status >%d<", status)

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	w.WriteHeader(status)

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}

	return xml.NewEncoder(w).Encode(s)
}
