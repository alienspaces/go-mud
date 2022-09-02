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
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/mapper"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

const (
	getDungeons  server.HandlerConfigKey = "get-dungeons"
	getDungeon   server.HandlerConfigKey = "get-dungeon"
	enterDungeon server.HandlerConfigKey = "enter-dungeon"
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
		enterDungeon: {
			Method:      http.MethodPost,
			Path:        "/api/v1/dungeons/:dungeon_id/enter",
			HandlerFunc: rnr.EnterDungeonHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypePublic,
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/dungeon-instance",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/dungeon-instance",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Enter a dungeon.",
			},
		},
	})
}

// GetDungeonHandler -
func (rnr *Runner) GetDungeonHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "GetDungeonHandler")
	l.Info("** Get dungeon handler **")

	var recs []*record.Dungeon
	var err error

	// Path parameters
	id := pp.ByName("dungeon_id")

	if id == "" {
		err := coreerror.NewNotFoundError("dungeon", id)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(id) {
		err := coreerror.NewPathParamInvalidTypeError("dungeon_id", id)
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
		err := coreerror.NewNotFoundError("dungeon", id)
		server.WriteError(l, w, err)
		return err
	}

	recs = append(recs, rec)

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

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// GetDungeonsHandler -
func (rnr *Runner) GetDungeonsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "GetDungeonsHandler")
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

// TODO: (game) Perhaps entering and exiting a dungeon should be an action?

// EnterDungeonHandler -
func (rnr *Runner) EnterDungeonHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "EnterDungeonHandler")
	l.Info("** Enter dungeon handler **")

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	if dungeonID == "" {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(dungeonID) {
		err := coreerror.NewPathParamInvalidTypeError("dungeon_id", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	if characterID == "" {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(characterID) {
		err := coreerror.NewPathParamInvalidTypeError("dungeon_id", characterID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting dungeon record ID >%s<", dungeonID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, false)
	if err != nil {
		l.Warn("failed getting dungeon record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if dungeonRec == nil {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting character record ID >%s<", characterID)

	characterRec, err := m.(*model.Model).GetCharacterRec(characterID, false)
	if err != nil {
		l.Warn("failed getting character record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if characterRec == nil {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Entering dungeon ID with character ID >%s<", dungeonID, characterID)

	characterInstanceRecordSet, err := m.(*model.Model).CharacterEnterDungeon(dungeonID, characterID)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Character instance record set >%#v<", characterInstanceRecordSet)

	dungeonActionRecordSet, err := m.(*model.Model).ProcessCharacterAction(
		characterInstanceRecordSet.CharacterInstanceRec.DungeonInstanceID,
		characterInstanceRecordSet.CharacterInstanceRec.ID,
		"look",
	)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	l.Debug("Resulting action >%#v<", dungeonActionRecordSet.ActionRec)
	l.Debug("Resulting action current location >%#v<", dungeonActionRecordSet.CurrentLocation)
	l.Debug("Resulting action target location >%#v<", dungeonActionRecordSet.TargetLocation)
	l.Debug("Resulting action character >%#v<", dungeonActionRecordSet.ActionCharacterRec)
	l.Debug("Resulting action monster >%#v<", dungeonActionRecordSet.ActionMonsterRec)

	// Response data
	responseData, err := mapper.ActionRecordSetToActionResponse(l, *dungeonActionRecordSet)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	res := schema.ActionResponse{
		Data: []schema.ActionResponseData{
			*responseData,
		},
	}

	l.Info("Response >%#v<", res)
	if len(res.Data) != 0 {
		l.Info("Response character >%#v<", res.Data[0].Character)
		l.Info("Response monster >%#v<", res.Data[0].Monster)
		l.Info("Response target character >%#v<", res.Data[0].TargetCharacter)
		l.Info("Response target monster >%#v<", res.Data[0].TargetMonster)
		l.Info("Response target object >%#v<", res.Data[0].TargetObject)
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// DungeonRequestDataToRecord -
func (rnr *Runner) DungeonRequestDataToRecord(data schema.DungeonData, rec *record.Dungeon) error {

	return nil
}

// RecordToDungeonResponseData -
func (rnr *Runner) RecordToDungeonResponseData(dungeonRec record.Dungeon) (schema.DungeonData, error) {

	data := schema.DungeonData{
		ID:          dungeonRec.ID,
		Name:        dungeonRec.Name,
		Description: dungeonRec.Description,
		CreatedAt:   dungeonRec.CreatedAt,
		UpdatedAt:   dungeonRec.UpdatedAt.Time,
	}

	return data, nil
}
