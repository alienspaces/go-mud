package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/schema"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	getDungeons server.HandlerConfigKey = "get-dungeons"
	getDungeon  server.HandlerConfigKey = "get-dungeon"
)

func (rnr *Runner) DungeonHandlerConfig(hc map[server.HandlerConfigKey]server.HandlerConfig) map[server.HandlerConfigKey]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[server.HandlerConfigKey]server.HandlerConfig{
		getDungeons: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons",
			HandlerFunc: rnr.GetDungeonsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypePublic,
				},
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
				Description: "List dungeons.",
			},
		},
		getDungeon: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons/:dungeon_id",
			HandlerFunc: rnr.GetDungeonHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypePublic,
				},
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
	})
}

// GetDungeonHandler -
func (rnr *Runner) GetDungeonHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "GetDungeonHandler")
	l.Info("** Get dungeon handler **")

	// Path parameters
	id := pp.ByName("dungeon_id")

	if id == "" {
		err := coreerror.NewResourceNotFoundError("dungeon", id)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(id) {
		err := coreerror.NewValidationPathParamTypeError("dungeon_id", id)
		server.WriteError(l, w, err)
		return err

	}

	l.Info("Getting dungeon record ID >%s<", id)

	rec, err := m.(*model.Model).GetDungeonRec(id, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewResourceNotFoundError("dungeon", id)
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	data, err := rnr.RecordToDungeonResponseData(*rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	res := schema.DungeonResponse{
		Data: []schema.DungeonData{
			data,
		},
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// GetDungeonsHandler -
func (rnr *Runner) GetDungeonsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "GetDungeonsHandler")
	l.Info("** Get dungeons handler **")

	var recs []*record.Dungeon
	var err error

	params := make(map[string]interface{})
	for paramName, paramValue := range qp {
		l.Info("Querying dungeon records with param name >%s< value >%v<", paramName, paramValue)
		params[paramName] = paramValue
	}

	l.Info("Querying dungeon records with params >%#v<", params)

	recs, err = m.(*model.Model).GetDungeonRecs(params, nil, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	data := []schema.DungeonData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToDungeonResponseData(*rec)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.DungeonResponse{
		Data: data,
	}

	l.Info("Responding with >%#v<", res)

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
