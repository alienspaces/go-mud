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
	getDungeonCharacter       string = "get-dungeon-character"
	postDungeonCharacterEnter string = "post-dungeon-character-enter"
	postDungeonCharacterExit  string = "post-dungeon-character-exit"
)

func (rnr *Runner) DungeonCharacterHandlerConfig(hc map[string]server.HandlerConfig) map[string]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[string]server.HandlerConfig{
		getDungeonCharacter: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons/:dungeon_id/characters/:character_id",
			HandlerFunc: rnr.GetDungeonCharacterHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateResponseSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/dungeoncharacter",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/dungeoncharacter",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a dungeon character.",
			},
		},
		postDungeonCharacterEnter: {
			Method:      http.MethodPost,
			Path:        "/api/v1/dungeons/:dungeon_id/characters/:character_id/enter",
			HandlerFunc: rnr.PostDungeonCharacterEnterHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateResponseSchema: &jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/dungeoncharacter",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/dungeoncharacter",
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
		postDungeonCharacterExit: {
			Method:      http.MethodPost,
			Path:        "/api/v1/dungeons/:dungeon_id/characters/:character_id/exit",
			HandlerFunc: rnr.PostDungeonCharacterExitHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
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
				Description: "Exit a dungeon.",
			},
		},
	})
}

// GetDungeonCharacterHandler -
func (rnr *Runner) GetDungeonCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithFunctionContext(l, "GetDungeonCharacterHandler")
	l.Info("** Get dungeon character handler **")

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	l.Info("Getting dungeon record ID >%s<", dungeonID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, nil)
	if err != nil {
		l.Warn("failed getting dungeon record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if dungeonRec == nil {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting character record ID >%s<", characterID)

	characterRec, err := m.(*model.Model).GetCharacterRec(characterID, nil)
	if err != nil {
		l.Warn("failed getting character record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if characterRec == nil {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	instanceViewRecordSet, err := rnr.getInstanceViewRecordSetByCharacterID(l, m, characterID)
	if err != nil {
		l.Warn("failed getting character instance record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if instanceViewRecordSet == nil {
		l.Warn("instance record set is nil")
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	data, err := dungeonCharacterResponseData(l, instanceViewRecordSet)
	if err != nil {
		l.Warn("failed mapping instance view record set to character response data")
		server.WriteError(l, w, err)
		return err
	}

	res := schema.DungeonCharacterResponse{
		Data: []schema.DungeonCharacterData{
			data,
		},
	}

	l.Info("Responding with >%#v<", res)

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// PostDungeonCharacterEnterHandler -
func (rnr *Runner) PostDungeonCharacterEnterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithFunctionContext(l, "PostDungeonCharacterEnterHandler")
	l.Info("** Dungeon character enter handler **")

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	l.Info("Getting dungeon record ID >%s<", dungeonID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, nil)
	if err != nil {
		l.Warn("failed getting dungeon record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if dungeonRec == nil {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting character record ID >%s<", characterID)

	characterRec, err := m.(*model.Model).GetCharacterRec(characterID, nil)
	if err != nil {
		l.Warn("failed getting character record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if characterRec == nil {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Entering dungeon ID >%s< with character ID >%s<", dungeonID, characterID)

	characterInstanceRecordSet, err := m.(*model.Model).CharacterEnterDungeon(dungeonID, characterID)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Character instance record set >%#v<", characterInstanceRecordSet)

	instanceViewRecordSet, err := rnr.getInstanceViewRecordSetByCharacterID(l, m, characterID)
	if err != nil {
		l.Warn("failed getting character instance record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if instanceViewRecordSet == nil {
		l.Warn("instance record set is nil")
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	data, err := dungeonCharacterResponseData(l, instanceViewRecordSet)
	if err != nil {
		l.Warn("failed mapping instance view record set to character response data")
		server.WriteError(l, w, err)
		return err
	}

	res := schema.DungeonCharacterResponse{
		Data: []schema.DungeonCharacterData{
			data,
		},
	}

	l.Info("Responding with >%#v<", res)

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// PostDungeonCharacterExitHandler -
func (rnr *Runner) PostDungeonCharacterExitHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithFunctionContext(l, "PostDungeonCharacterExitHandler")
	l.Info("** Dungeon character exit handler **")

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	l.Info("Getting dungeon record ID >%s<", dungeonID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, nil)
	if err != nil {
		l.Warn("failed getting dungeon record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if dungeonRec == nil {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting character record ID >%s<", characterID)

	characterRec, err := m.(*model.Model).GetCharacterRec(characterID, nil)
	if err != nil {
		l.Warn("failed getting character record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if characterRec == nil {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	instanceViewRecordSet, err := rnr.getInstanceViewRecordSetByCharacterID(l, m, characterID)
	if err != nil {
		l.Warn("failed getting character instance record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if instanceViewRecordSet == nil {
		l.Warn("instance record set is nil")
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	if instanceViewRecordSet.DungeonInstanceViewRec.DungeonID != dungeonID {
		l.Warn("dungeon ID >%s< does not contain character ID >%s<", dungeonID, characterID)
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Exiting dungeon ID >%s< with character ID >%s<", dungeonID, characterID)

	err = m.(*model.Model).CharacterExitDungeon(characterID)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Get updated character record
	characterRec, err = m.(*model.Model).GetCharacterRec(characterID, nil)
	if err != nil {
		l.Warn("failed getting character record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	if characterRec == nil {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	responseData, err := characterResponseData(l, characterRec, nil)
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

	return nil
}
