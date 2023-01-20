package runner

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/schema"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
)

const (
	postAction server.HandlerConfigKey = "post-action"
)

func (rnr *Runner) ActionHandlerConfig(hc map[server.HandlerConfigKey]server.HandlerConfig) map[server.HandlerConfigKey]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[server.HandlerConfigKey]server.HandlerConfig{
		postAction: {
			Method:      http.MethodPost,
			Path:        "/api/v1/dungeons/:dungeon_id/characters/:character_id/actions",
			HandlerFunc: rnr.PostActionHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypePublic,
				},
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
							Location: "schema/docs/action",
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
	})
}

// PostActionHandler -
func (rnr *Runner) PostActionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "PostActionHandler")
	l.Info("** Post action handler **")

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
		err := coreerror.NewPathParamInvalidTypeError("character_id", characterID)
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

	req := schema.ActionRequest{}

	err = server.ReadRequest(l, r, &req)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	sentence := strings.ToLower(req.Data.Sentence)

	l.Info("Verifying character instance for dungeon_id >%s< character_id >%s<", dungeonID, characterID)

	characterInstanceRecs, err := m.(*model.Model).GetCharacterInstanceRecs(
		map[string]interface{}{
			"character_id": characterID,
		}, nil, false,
	)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	if len(characterInstanceRecs) == 0 {
		err := coreerror.NewPathParamInvalidError("character_id", characterID, "character has not entered a dungeon")
		server.WriteError(l, w, err)
		return err
	}

	if len(characterInstanceRecs) > 1 {
		l.Warn("Unexpected number of character instance records returned >%d<", len(characterInstanceRecs))
		err := coreerror.NewInternalError()
		server.WriteError(l, w, err)
		return err
	}

	characterInstanceRec := characterInstanceRecs[0]

	dungeonInstanceRec, err := m.(*model.Model).GetDungeonInstanceRec(characterInstanceRec.DungeonInstanceID, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	if dungeonInstanceRec == nil {
		err := coreerror.NewPathParamInvalidError("dungeon_id", dungeonID, "dungeon does not exists")
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Creating dungeon character action >%s<", sentence)

	dungeonActionRecordSet, err := m.(*model.Model).ProcessCharacterAction(
		dungeonInstanceRec.ID,
		characterInstanceRec.ID,
		sentence,
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
	responseData, err := actionResponseData(l, *dungeonActionRecordSet)
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
		l.Info("Response character >%#v<", res.Data[0].ActionCharacter)
		l.Info("Response monster >%#v<", res.Data[0].ActionMonster)
		l.Info("Response target character >%#v<", res.Data[0].ActionTargetCharacter)
		l.Info("Response target monster >%#v<", res.Data[0].ActionTargetMonster)
		l.Info("Response target object >%#v<", res.Data[0].ActionTargetObject)
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
