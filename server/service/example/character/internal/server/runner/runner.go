package runner

import (
	"net/http"

	"gitlab.com/alienspaces/go-boilerplate/server/constant"

	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/server"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/runnable"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/model"
)

// Runner -
type Runner struct {
	server.Runner
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
		// 0 - Get any, Administrator role or Default role, player ID not required
		{
			Method:      http.MethodGet,
			Path:        "/api/characters",
			HandlerFunc: r.GetCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleAdministrator,
					constant.AuthRoleDefault,
				},
				ValidateQueryParams: []string{},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query characters.",
			},
		},
		// 1 - Get with character ID, Administrator role or Default role, player ID not required
		{
			Method:      http.MethodGet,
			Path:        "/api/characters/:character_id",
			HandlerFunc: r.GetCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleAdministrator,
					constant.AuthRoleDefault,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get an character.",
			},
		},
		// 2 - Get any, Default or Administrator role, playerID required, player ID in path must match idcharacter
		{
			Method:      http.MethodGet,
			Path:        "/api/players/:player_id/characters",
			HandlerFunc: r.GetPlayerCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleAdministrator,
					constant.AuthRoleDefault,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				ValidatePathParams: map[string]server.ValidatePathParam{
					"player_id": server.ValidatePathParam{
						MatchIdentity: true,
					},
				},
				ValidateQueryParams: []string{},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query characters.",
			},
		},
		// 3 - Get with character ID, Default or Administrator role, playerID required, player ID in path must match idcharacter
		{
			Method:      http.MethodGet,
			Path:        "/api/players/:player_id/characters/:character_id",
			HandlerFunc: r.GetPlayerCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleAdministrator,
					constant.AuthRoleDefault,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				ValidatePathParams: map[string]server.ValidatePathParam{
					"player_id": server.ValidatePathParam{
						MatchIdentity: true,
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get an character.",
			},
		},
		// 4 - Post without character ID, Default or Administrator role, playerID required, player ID in path must match idcharacter
		{
			Method:      http.MethodPost,
			Path:        "/api/players/:player_id/characters",
			HandlerFunc: r.PostPlayerCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleAdministrator,
					constant.AuthRoleDefault,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				ValidatePathParams: map[string]server.ValidatePathParam{
					"player_id": server.ValidatePathParam{
						MatchIdentity: true,
					},
				},
				ValidateSchemaLocation: "character",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create an character.",
			},
		},
		// 5 - Post with character ID, Default or Administrator role, playerID required, player ID in path must match idcharacter
		{
			Method:      http.MethodPost,
			Path:        "/api/players/:player_id/characters/:character_id",
			HandlerFunc: r.PostPlayerCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleAdministrator,
					constant.AuthRoleDefault,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				ValidatePathParams: map[string]server.ValidatePathParam{
					"player_id": server.ValidatePathParam{
						MatchIdentity: true,
					},
				},
				ValidateSchemaLocation: "character",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create an character.",
			},
		},
		// 6 - Put with character ID, Default or Administrator role, playerID required, player ID in
		// path must match idcharacter
		{
			Method:      http.MethodPut,
			Path:        "/api/players/:player_id/characters/:character_id",
			HandlerFunc: r.PutPlayerCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthTypes: []string{
					auth.AuthTypeJWT,
				},
				AuthRequireAnyRole: []string{
					constant.AuthRoleAdministrator,
					constant.AuthRoleDefault,
				},
				AuthRequireAllIdentities: []string{
					constant.AuthIdentityPlayerID,
				},
				ValidatePathParams: map[string]server.ValidatePathParam{
					"player_id": server.ValidatePathParam{
						MatchIdentity: true,
					},
				},
				ValidateSchemaLocation: "character",
				ValidateSchemaMain:     "main.schema.json",
				ValidateSchemaReferences: []string{
					"data.schema.json",
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a character.",
			},
		},
		// 6 - Documentation
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

	l.Info("** Character Model **")

	m, err := model.NewModel(rnr.Config, l, rnr.Store)
	if err != nil {
		l.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
