package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
)

const (
	getDungeons string = "get-dungeons"
	getDungeon  string = "get-dungeon"
)

func (rnr *Runner) DungeonHandlerConfig(hc map[string]server.HandlerConfig) map[string]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[string]server.HandlerConfig{
		getDungeons: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons",
			HandlerFunc: rnr.getDungeonsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateRequestSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/dungeon",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/dungeon",
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
			HandlerFunc: rnr.getDungeonHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateResponseSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/dungeon",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/dungeon",
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

// getDungeonHandler -
func (rnr *Runner) getDungeonHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "getDungeonHandler")

	// Path parameters
	id := pp.ByName("dungeon_id")

	l.Info("Getting dungeon record ID >%s<", id)

	rec, err := m.(*model.Model).GetDungeonRec(id, nil)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewNotFoundError("dungeon", id)
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

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getDungeonsHandler -
func (rnr *Runner) getDungeonsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "getDungeonsHandler")

	opts := queryparam.ToSQLOptions(qp)

	l.Info("Querying dungeon records with opts >%#v<", opts)

	recs, err := m.(*model.Model).GetDungeonRecs(opts)
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

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
