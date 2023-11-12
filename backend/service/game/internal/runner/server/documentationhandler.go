package runner

import (
	"fmt"
	"net/http"
	"sort"

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

	type pathDocument struct {
		method      string
		path        string
		description string
	}

	paths := []string{}
	pathDocuments := map[string]pathDocument{}

	for _, cfg := range rnr.HandlerConfig {
		pathKey := fmt.Sprintf("%s/%s", cfg.Path, cfg.Method)
		paths = append(paths, pathKey)
		pathDocuments[pathKey] = pathDocument{
			method:      cfg.Method,
			path:        cfg.Path,
			description: cfg.DocumentationConfig.Description,
		}
	}

	sort.Strings(paths)

	rowClass := "alt-one"
	content := "<table>\n"
	for _, pathKey := range paths {
		if rowClass == "alt-two" {
			rowClass = "alt-one"
		} else {
			rowClass = "alt-two"
		}
		content = fmt.Sprintf("%s\n<tr class=\"%s\"><td class=\"method\">%s</td><td class=\"path\">%s</td></tr>", content, rowClass, pathDocuments[pathKey].method, pathDocuments[pathKey].path)
		content = fmt.Sprintf("%s\n<tr class=\"%s\"><td></td><td class=\"description\">%s</td></tr>", content, rowClass, pathDocuments[pathKey].description)
	}
	content = fmt.Sprintf("%s</table>", content)

	html := `
<html>
<head>
	<style>
		body {
			font-family: monospace, monospace;			
		}
		table {
			border-spacing: 0px;
		}
		.alt-one {
			background-color: #bcbcbc;
		}
		.alt-two {
			background-color: #dedede;
		}
		.method {
			padding-top: 6px;
			padding-left: 5px;
			padding-right: 5px;
			padding-bottom: 2px;
		}
		.path {
			padding-top: 6px;
			padding-left: 5px;
			padding-right: 5px;
			padding-bottom: 2px;
		}
		.description {
			padding-top: 4px;
			padding-left: 5px;
			padding-right: 5px;
			padding-bottom: 6px;
		}
	</style>
</head>
%s
</html>
	`

	docs := fmt.Sprintf(html, content)

	// content type html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// status
	w.WriteHeader(http.StatusOK)

	// content
	written, err := w.Write([]byte(docs))
	if err != nil {
		l.Warn("failed writing documentation >%v<", err)
		return err
	}

	l.Info("wrote >%d< bytes", written)

	return nil
}
