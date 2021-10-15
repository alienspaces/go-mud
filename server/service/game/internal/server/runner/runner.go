package runner

import (
	"net/http"

	"gitlab.com/alienspaces/go-mud/server/core/auth"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
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
			Path:        "/api/dungeons",
			HandlerFunc: r.GetDungeonsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query dungeons.",
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/dungeons/:dungeon_id",
			HandlerFunc: r.GetDungeonHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a dungeon.",
			},
		},
		// {
		// 	Method:      http.MethodPost,
		// 	Path:        "/api/dungeons",
		// 	HandlerFunc: r.PostDungeonsHandler,
		// 	MiddlewareConfig: server.MiddlewareConfig{
		// 		AuthTypes: []string{
		// 			auth.AuthTypeJWT,
		// 		},
		// 		ValidateSchemaLocation: "dungeon",
		// 		ValidateSchemaMain:     "main.schema.json",
		// 		ValidateSchemaReferences: []string{
		// 			"data.schema.json",
		// 		},
		// 	},
		// 	DocumentationConfig: server.DocumentationConfig{
		// 		Document:    true,
		// 		Description: "Create a dungeon.",
		// 	},
		// },
		// {
		// 	Method:      http.MethodPost,
		// 	Path:        "/api/dungeons/:dungeon_id",
		// 	HandlerFunc: r.PostDungeonsHandler,
		// 	MiddlewareConfig: server.MiddlewareConfig{
		// 		AuthTypes: []string{
		// 			auth.AuthTypeJWT,
		// 		},
		// 		ValidateSchemaLocation: "dungeon",
		// 		ValidateSchemaMain:     "main.schema.json",
		// 		ValidateSchemaReferences: []string{
		// 			"data.schema.json",
		// 		},
		// 	},
		// 	DocumentationConfig: server.DocumentationConfig{
		// 		Document:    true,
		// 		Description: "Create a dungeon.",
		// 	},
		// },
		// {
		// 	Method:      http.MethodPut,
		// 	Path:        "/api/dungeons/:dungeon_id",
		// 	HandlerFunc: r.PutDungeonsHandler,
		// 	MiddlewareConfig: server.MiddlewareConfig{
		// 		AuthTypes: []string{
		// 			auth.AuthTypeJWT,
		// 		},
		// 		ValidateSchemaLocation: "dungeon",
		// 		ValidateSchemaMain:     "main.schema.json",
		// 		ValidateSchemaReferences: []string{
		// 			"data.schema.json",
		// 		},
		// 	},
		// 	DocumentationConfig: server.DocumentationConfig{
		// 		Document:    true,
		// 		Description: "Update a dungeon.",
		// 	},
		// },
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

	l.Info("** Dungeon Model **")

	m, err := model.NewModel(rnr.Config, l, rnr.Store)
	if err != nil {
		l.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
