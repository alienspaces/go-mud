package runner

import (
	"net/http"

	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/server"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/runnable"
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/model"
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

	r.HandlerConfig = []server.HandlerConfig{
		{
			Method:      http.MethodGet,
			Path:        "/api/templates",
			HandlerFunc: r.GetTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query templates.",
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/templates/:template_id",
			HandlerFunc: r.GetTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a template.",
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/templates",
			HandlerFunc: r.PostTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				ValidateSchemaLocation: "template",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a template.",
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/templates/:template_id",
			HandlerFunc: r.PostTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				ValidateSchemaLocation: "template",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a template.",
			},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/templates/:template_id",
			HandlerFunc: r.PutTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				ValidateSchemaLocation: "template",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a template.",
			},
		},
		{
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
		l.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
