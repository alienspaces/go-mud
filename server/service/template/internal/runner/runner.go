package runner

import (
	"net/http"

	"gitlab.com/alienspaces/go-mud/server/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/model"
)

const (
	getTemplates     server.HandlerConfigKey = "get-templates"
	getTemplate      server.HandlerConfigKey = "get-template"
	postTemplates    server.HandlerConfigKey = "post-templates"
	postTemplate     server.HandlerConfigKey = "post-template"
	putTemplate      server.HandlerConfigKey = "put-template"
	getDocumentation server.HandlerConfigKey = "getDocumentation"
)

// Runner -
type Runner struct {
	server.Runner
}

// Fault -
type Fault struct {
	Error error
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// NewRunner -
func NewRunner() *Runner {

	r := Runner{}

	r.RouterFunc = r.Router
	r.MiddlewareFunc = r.Middleware
	r.HandlerFunc = r.Handler
	r.ModellerFunc = r.Modeller

	r.HandlerConfig = map[server.HandlerConfigKey]server.HandlerConfig{
		getTemplates: {
			Method:      http.MethodGet,
			Path:        "/api/templates",
			HandlerFunc: r.GetTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypeJWT,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query templates.",
			},
		},
		getTemplate: {
			Method:      http.MethodGet,
			Path:        "/api/templates/:template_id",
			HandlerFunc: r.GetTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypeJWT,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a template.",
			},
		},
		postTemplates: {
			Method:      http.MethodPost,
			Path:        "/api/templates",
			HandlerFunc: r.PostTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypeJWT,
				},
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/template",
						Name:     "create.request.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/template",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a template.",
			},
		},
		postTemplate: {
			Method:      http.MethodPost,
			Path:        "/api/templates/:template_id",
			HandlerFunc: r.PostTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypeJWT,
				},
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/template",
						Name:     "create.request.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/template",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a template.",
			},
		},
		putTemplate: {
			Method:      http.MethodPut,
			Path:        "/api/templates/:template_id",
			HandlerFunc: r.PutTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypeJWT,
				},
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/template",
						Name:     "update.request.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/template",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a template.",
			},
		},
		getDocumentation: {
			Method:           http.MethodGet,
			Path:             "/api",
			HandlerFunc:      r.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
		},
	}

	return &r
}

// Modeller -
func (rnr *Runner) Modeller(l logger.Logger) (modeller.Modeller, error) {

	l.Info("** Template Model **")

	m, err := model.NewModel(rnr.Config, l, rnr.Store)
	if err != nil {
		l.Warn("failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
