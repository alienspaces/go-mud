package runner

import (
	"net/http"

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
		// Dungeons
		{
			Method:           http.MethodGet,
			Path:             "/api/v1/dungeons",
			HandlerFunc:      r.GetDungeonsHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query dungeons.",
			},
		},
		{
			Method:           http.MethodGet,
			Path:             "/api/v1/dungeons/:dungeon_id",
			HandlerFunc:      r.GetDungeonHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a dungeon.",
			},
		},
		// Characters
		{
			Method:           http.MethodGet,
			Path:             "/api/v1/dungeons/:dungeon_id/characters",
			HandlerFunc:      r.GetDungeonCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get characters.",
			},
		},
		{
			Method:           http.MethodGet,
			Path:             "/api/v1/dungeons/:dungeon_id/characters/:character_id",
			HandlerFunc:      r.GetDungeonCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a character.",
			},
		},
		{
			Method:           http.MethodPost,
			Path:             "/api/v1/dungeons/:dungeon_id/characters",
			HandlerFunc:      r.PostDungeonCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a character.",
			},
		},
		{
			Method:           http.MethodPut,
			Path:             "/api/v1/dungeons/:dungeon_id/characters",
			HandlerFunc:      r.PutDungeonCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a character.",
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/dungeons/:dungeon_id/characters/:character_id/actions",
			HandlerFunc: r.PostDungeonCharacterActionsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateSchemaLocation: "game",
				ValidateSchemaMain:     "create.schema.json",
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a dungeon character action.",
			},
		},
		// {
		// 	Method:      http.MethodPost,
		// 	Path:        "/api/v1/dungeons",
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
		// 	Path:        "/api/v1/dungeons/:dungeon_id",
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
		// 	Path:        "/api/v1/dungeons/:dungeon_id",
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
			Path:             "/",
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
