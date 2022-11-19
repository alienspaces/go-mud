package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/server/core/error"
	"gitlab.com/alienspaces/go-mud/server/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

const (
	getCharacters server.HandlerConfigKey = "get-characters"
	getCharacter  server.HandlerConfigKey = "get-character"
	postCharacter server.HandlerConfigKey = "post-character"
	putCharacter  server.HandlerConfigKey = "put-character"
)

// TODO: Add character instance data to all responses as there may only be a single
// character instance.
func (rnr *Runner) CharacterHandlerConfig(hc map[server.HandlerConfigKey]server.HandlerConfig) map[server.HandlerConfigKey]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[server.HandlerConfigKey]server.HandlerConfig{
		getCharacters: {
			Method:      http.MethodGet,
			Path:        "/api/v1/characters",
			HandlerFunc: rnr.GetCharactersHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypePublic,
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/character",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/character",
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
					server.AuthenTypePublic,
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/character",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/character",
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
					server.AuthenTypePublic,
				},
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/character",
						Name:     "create.request.schema.json",
					},
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/character",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/character",
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
func (rnr *Runner) GetCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "GetCharacterHandler")
	l.Info("** Get character handler **")

	var recs []*record.Character
	var err error

	// Path parameters
	characterID := pp.ByName("character_id")

	if characterID == "" {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(characterID) {
		err := coreerror.NewPathParamInvalidTypeError("character_id", characterID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting character record ID >%s<", characterID)

	rec, err := m.(*model.Model).GetCharacterRec(characterID, false)
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
		responseData, err := rnr.CharacterRecordWithInstanceViewRecordSetToDungeonCharacterResponseData(*rec, instanceViewRecordSet)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.CharacterResponse{
		Data: data,
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// GetCharactersHandler -
func (rnr *Runner) GetCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "GetCharactersHandler")
	l.Info("** Get characters handler **")

	var recs []*record.Character
	var err error

	l.Info("Querying character records")

	// Add query parameters
	params := make(map[string]interface{})
	for paramName, paramValue := range qp {
		l.Info("Querying dungeon records with param name >%s< value >%v<", paramName, paramValue)
		params[paramName] = paramValue
	}

	recs, err = m.(*model.Model).GetCharacterRecs(params, nil, false)
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
		responseData, err := rnr.CharacterRecordWithInstanceViewRecordSetToDungeonCharacterResponseData(*rec, instanceViewRecordSet)
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

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// PostCharacterHandler -
func (rnr *Runner) PostCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "PostCharacterHandler")
	l.Info("** Post character handler **")

	req := schema.CharacterRequest{}
	err := server.ReadRequest(l, r, &req)
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
	responseData, err := rnr.CharacterRecordWithInstanceViewRecordSetToDungeonCharacterResponseData(rec, instanceViewRecordSet)
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

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// PutCharactersHandler -
func (rnr *Runner) PutCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "PutCharacterHandler")
	l.Info("** Put character handler **")

	// Path parameters
	id := pp.ByName("character_id")

	if id == "" {
		err := coreerror.NewNotFoundError("character", id)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(id) {
		err := coreerror.NewPathParamInvalidTypeError("character_id", id)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Updating character ID >%s<", id)

	rec, err := m.(*model.Model).GetCharacterRec(id, false)
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

	req := schema.CharacterRequest{}

	err = server.ReadRequest(l, r, &req)
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
	responseData, err := rnr.CharacterRecordWithInstanceViewRecordSetToDungeonCharacterResponseData(*rec, instanceViewRecordSet)
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

	err = server.WriteResponse(l, w, res)
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
