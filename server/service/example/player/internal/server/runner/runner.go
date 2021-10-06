package runner

import (
	"net/http"

	"gitlab.com/alienspaces/go-boilerplate/server/constant"
	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/server"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/runnable"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
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
		// 0 - Authentication
		{
			Method:      http.MethodPost,
			Path:        "/api/auth",
			HandlerFunc: r.PostAuthHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateSchemaLocation: "auth",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Authenticate OAuth provider token.",
			},
		},
		// 1 - Refresh Authentication
		{
			Method:      http.MethodPost,
			Path:        "/api/auth-refresh",
			HandlerFunc: r.PostAuthRefreshHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateSchemaLocation: "authrefresh",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Refresh authentication token.",
			},
		},

		// TODO: Provide different default and administrator player handlers constraining record access as appropriate

		// 2 - Players - Get many
		{
			Method:      http.MethodGet,
			Path:        "/api/players",
			HandlerFunc: r.GetPlayersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleDefault,
					constant.AuthRoleAdministrator,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query players.",
			},
		},
		// 3 - Players - Get one
		{
			Method:      http.MethodGet,
			Path:        "/api/players/:player_id",
			HandlerFunc: r.GetPlayersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleDefault,
					constant.AuthRoleAdministrator,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get an player.",
			},
		},
		// 4 - Players - Create one without ID
		{
			Method:      http.MethodPost,
			Path:        "/api/players",
			HandlerFunc: r.PostPlayersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleDefault,
					constant.AuthRoleAdministrator,
				},
				ValidateSchemaLocation: "player",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a player.",
			},
		},
		// 5 - Players - Create one with ID
		{
			Method:      http.MethodPost,
			Path:        "/api/players/:player_id",
			HandlerFunc: r.PostPlayersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleDefault,
					constant.AuthRoleAdministrator,
				},
				ValidateSchemaLocation: "player",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a player.",
			},
		},
		// 6 - Players - Update one
		{
			Method:      http.MethodPut,
			Path:        "/api/players/:player_id",
			HandlerFunc: r.PutPlayersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleDefault,
					constant.AuthRoleAdministrator,
				},
				ValidateSchemaLocation: "player",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a player.",
			},
		},
		// 7 - Documentation
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

	l.Info("** Player Model **")

	m, err := model.NewModel(rnr.Config, l, rnr.Store)
	if err != nil {
		l.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
