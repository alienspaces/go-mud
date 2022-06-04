package runner

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/go-mud/server/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
)

const (
	getDungeons           server.HandlerConfigKey = "get-dungeons"
	getDungeon            server.HandlerConfigKey = "get-dungeon"
	getCharacters         server.HandlerConfigKey = "get-characters"
	getCharacter          server.HandlerConfigKey = "get-character"
	postCharacter         server.HandlerConfigKey = "post-character"
	putCharacter          server.HandlerConfigKey = "put-character"
	postAction            server.HandlerConfigKey = "post-action"
	getDocumentationRoot  server.HandlerConfigKey = "get-documentation-root"
	getDocumentationAPI   server.HandlerConfigKey = "get-documentation-api"
	getDocumentationAPIV1 server.HandlerConfigKey = "get-documentation-api-v1"
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
func NewRunner(c configurer.Configurer, l logger.Logger) (*Runner, error) {

	r := Runner{}

	r.Log = l
	if r.Log == nil {
		msg := "logger undefined, cannot init runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	r.Config = c
	if r.Config == nil {
		msg := "configurer undefined, cannot init runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	r.RouterFunc = r.Router
	r.MiddlewareFunc = r.Middleware
	r.HandlerFunc = r.Handler
	r.ModellerFunc = r.Modeller

	r.HandlerConfig = map[server.HandlerConfigKey]server.HandlerConfig{
		getDungeons: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons",
			HandlerFunc: r.GetDungeonsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/dungeon",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query dungeons.",
			},
		},
		getDungeon: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons/:dungeon_id",
			HandlerFunc: r.GetDungeonHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/dungeon",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a dungeon.",
			},
		},
		getCharacters: {
			Method:      http.MethodGet,
			Path:        "/api/v1/characters",
			HandlerFunc: r.GetCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/dungeoncharacter",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get characters.",
			},
		},
		getCharacter: {
			Method:      http.MethodGet,
			Path:        "/api/v1/characters/:character_id",
			HandlerFunc: r.GetCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/dungeoncharacter",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a character.",
			},
		},
		postCharacter: {
			Method:      http.MethodPost,
			Path:        "/api/v1/characters",
			HandlerFunc: r.PostCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/dungeoncharacter",
						Name:     "create.request.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/dungeoncharacter",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a character.",
			},
		},
		// TODO: Complete schema validation and implementation
		putCharacter: {
			Method:           http.MethodPut,
			Path:             "/api/v1/characters",
			HandlerFunc:      r.PutCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a character.",
			},
		},
		postAction: {
			Method:      http.MethodPost,
			Path:        "/api/v1/dungeons/:dungeon_id/characters/:character_id/actions",
			HandlerFunc: r.PostDungeonCharacterActionsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/action",
						Name:     "create.request.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/action",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a dungeon character action.",
			},
		},
		getDocumentationRoot: {
			Method:           http.MethodGet,
			Path:             "/",
			HandlerFunc:      r.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
		},
		getDocumentationAPI: {
			Method:           http.MethodGet,
			Path:             "/api",
			HandlerFunc:      r.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
		},
		getDocumentationAPIV1: {
			Method:           http.MethodGet,
			Path:             "/api/v1",
			HandlerFunc:      r.GetDocumentationHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
		},
	}

	return &r, nil
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
