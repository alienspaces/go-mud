package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	getCharacters string = "get-characters"
	getCharacter  string = "get-character"
	postCharacter string = "post-character"
	putCharacter  string = "put-character"
)

func (rnr *Runner) CharacterHandlerConfig(hc map[string]server.HandlerConfig) map[string]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[string]server.HandlerConfig{
		getCharacters: {
			Method:      http.MethodGet,
			Path:        "/api/v1/characters",
			HandlerFunc: rnr.GetCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateParamsConfig: &server.ValidateParamsConfig{
					QueryParamSchema: &jsonschema.SchemaWithReferences{
						Main: jsonschema.Schema{
							Location: "schema/game/character",
							Name:     "query.schema.json",
						},
					},
				},
				ValidateResponseSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/character",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/character",
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
			HandlerFunc: rnr.GetCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateParamsConfig: &server.ValidateParamsConfig{
					PathParamSchema: &jsonschema.SchemaWithReferences{
						Main: jsonschema.Schema{
							Location: "schema/game/character",
							Name:     "path.schema.json",
						},
					},
				},
				ValidateResponseSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/character",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/character",
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
			HandlerFunc: rnr.PostCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateRequestSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/character",
						Name:     "create.request.schema.json",
					},
				},
				ValidateResponseSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/character",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/character",
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
		// TODO: (game) Complete schema validation and implementation
		putCharacter: {
			Method:           http.MethodPut,
			Path:             "/api/v1/characters",
			HandlerFunc:      rnr.PutCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a character.",
			},
		},
	})
}

// GetCharacterHandler -
func (rnr *Runner) GetCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "GetCharacterHandler")
	l.Info("** Get character handler **")

	var recs []*record.Character
	var err error

	// Path parameters
	characterID := pp.ByName("character_id")

	l.Info("Getting character record ID >%s<", characterID)

	rec, err := m.(*model.Model).GetCharacterRec(characterID, nil)
	if err != nil {
		l.Warn("failed getting character record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	recs = append(recs, rec)

	// Assign response properties
	data := []schema.DungeonCharacterData{}
	for _, rec := range recs {

		instanceViewRecordSet, err := rnr.getInstanceViewRecordSetByCharacterID(l, m, characterID)
		if err != nil {
			l.Warn("failed getting character instance record >%v<", err)
			server.WriteError(l, w, err)
			return err
		}

		// Response data
		responseData, err := characterResponseData(l, rec, instanceViewRecordSet)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.CharacterResponse{
		Data: data,
	}

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func QueryParamsToSQLOptions(qp *queryparam.QueryParams) *coresql.Options {

	if len(qp.SortColumns) == 0 {
		qp.SortColumns = []queryparam.SortColumn{
			{
				Col: "created_at",
			},
		}
	}

	return queryparam.ToSQLOptions(qp)
}

// GetCharactersHandler -
func (rnr *Runner) GetCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "GetCharactersHandler")
	l.Info("** Get characters handler **")

	opts := QueryParamsToSQLOptions(qp)

	l.Info("Querying character records with params >%#v<", qp)

	recs, err := m.(*model.Model).GetCharacterRecs(opts)
	if err != nil {
		l.Warn("failed getting dungeon character records >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	data := []schema.DungeonCharacterData{}
	for _, rec := range recs {

		instanceViewRecordSet, err := rnr.getInstanceViewRecordSetByCharacterID(l, m, rec.ID)
		if err != nil {
			l.Warn("failed getting character instance record >%v<", err)
			server.WriteError(l, w, err)
			return err
		}

		// Response data
		responseData, err := characterResponseData(l, rec, instanceViewRecordSet)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.CharacterResponse{
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

// PostCharacterHandler -
func (rnr *Runner) PostCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "PostCharacterHandler")
	l.Info("** Post character handler **")

	req := &schema.CharacterRequest{}
	req, err := server.ReadRequest(l, r, req)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	rec := record.Character{}

	// Record data
	err = rnr.CharacterRequestDataToRecord(req.Data, &rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Creating character record >%#v<", rec)

	err = m.(*model.Model).CreateCharacterRec(&rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	instanceViewRecordSet, err := rnr.getInstanceViewRecordSetByCharacterID(l, m, rec.ID)
	if err != nil {
		l.Warn("failed getting character instance record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	responseData, err := characterResponseData(l, &rec, instanceViewRecordSet)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	res := schema.CharacterResponse{
		Data: []schema.DungeonCharacterData{
			responseData,
		},
	}

	l.Info("Writing response >%#v<", res)

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// PutCharactersHandler -
func (rnr *Runner) PutCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "PutCharacterHandler")
	l.Info("** Put character handler **")

	// Path parameters
	id := pp.ByName("character_id")

	l.Info("Updating character ID >%s<", id)

	rec, err := m.(*model.Model).GetCharacterRec(id, nil)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewNotFoundError("character", id)
		server.WriteError(l, w, err)
		return err
	}

	req := &schema.CharacterRequest{}
	req, err = server.ReadRequest(l, r, req)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Record data
	err = rnr.CharacterRequestDataToRecord(req.Data, rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	err = m.(*model.Model).UpdateCharacterRec(rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	instanceViewRecordSet, err := rnr.getInstanceViewRecordSetByCharacterID(l, m, rec.ID)
	if err != nil {
		l.Warn("failed getting character instance record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	responseData, err := characterResponseData(l, rec, instanceViewRecordSet)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	res := schema.CharacterResponse{
		Data: []schema.DungeonCharacterData{
			responseData,
		},
	}

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return err
}

// CharacterRequestDataToRecord -
func (rnr *Runner) CharacterRequestDataToRecord(data schema.DungeonCharacterData, rec *record.Character) error {

	rec.Name = data.CharacterName
	rec.Strength = data.CharacterStrength
	rec.Dexterity = data.CharacterDexterity
	rec.Intelligence = data.CharacterIntelligence

	return nil
}
