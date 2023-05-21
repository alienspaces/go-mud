package runner

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

const (
	getDocumentationRoot  string = "get-documentation-root"
	getDocumentationAPI   string = "get-documentation-api"
	getDocumentationAPIV1 string = "get-documentation-api-v1"
)

func (rnr *Runner) DocumentationHandlerConfig(hc map[string]server.HandlerConfig) map[string]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[string]server.HandlerConfig{
		getDocumentationRoot: {
			Method:      http.MethodGet,
			Path:        "/",
			HandlerFunc: rnr.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
			},
		},
		getDocumentationAPI: {
			Method:      http.MethodGet,
			Path:        "/api",
			HandlerFunc: rnr.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
			},
		},
		getDocumentationAPIV1: {
			Method:      http.MethodGet,
			Path:        "/api/v1",
			HandlerFunc: rnr.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
			},
		},
	})
}

// GetDocumentationHandler -
func (rnr *Runner) GetDocumentationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {

	l.Info("** get schema documentation handler ** p >%#v< m >%#v<", pp, m)

	docs, err := rnr.GenerateHandlerDocumentation(rnr.GetHandlerConfigs())
	if err != nil {
		msg := fmt.Sprintf("unable to load schema documentation >%v<, cannot init runner", err)
		rnr.Log.Error(msg)
		return fmt.Errorf(msg)
	}

	// content type html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// status
	w.WriteHeader(http.StatusOK)

	// content
	written, err := w.Write(docs)
	if err != nil {
		l.Warn("failed writing documentation >%v<", err)
		return err
	}

	l.Info("wrote >%d< bytes", written)

	return nil
}
