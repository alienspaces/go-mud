package runner

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
)

const (
	getDocumentationRoot  server.HandlerConfigKey = "get-documentation-root"
	getDocumentationAPI   server.HandlerConfigKey = "get-documentation-api"
	getDocumentationAPIV1 server.HandlerConfigKey = "get-documentation-api-v1"
)

func (rnr *Runner) DocumentationHandlerConfig(hc map[server.HandlerConfigKey]server.HandlerConfig) map[server.HandlerConfigKey]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[server.HandlerConfigKey]server.HandlerConfig{
		getDocumentationRoot: {
			Method:           http.MethodGet,
			Path:             "/",
			HandlerFunc:      rnr.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
		},
		getDocumentationAPI: {
			Method:           http.MethodGet,
			Path:             "/api",
			HandlerFunc:      rnr.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
		},
		getDocumentationAPIV1: {
			Method:           http.MethodGet,
			Path:             "/api/v1",
			HandlerFunc:      rnr.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
		},
	})
}

// GetDocumentationHandler -
func (rnr *Runner) GetDocumentationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** get schema documentation handler ** p >%#v< m >%#v<", pp, m)

	docs, err := rnr.GenerateHandlerDocumentation(rnr.GetMessageConfigs(), rnr.GetHandlerConfigs())
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
