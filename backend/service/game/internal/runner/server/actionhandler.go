package runner

import (
	"net/http"
	"strings"

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
)

const (
	postAction string = "post-action"
)

func (rnr *Runner) ActionHandlerConfig(hc map[string]server.HandlerConfig) map[string]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[string]server.HandlerConfig{
		postAction: {
			Method:      http.MethodPost,
			Path:        "/api/v1/dungeons/:dungeon_id/characters/:character_id/actions",
			HandlerFunc: rnr.PostActionHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/action",
						Name:     "create.request.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/dungeon",
							Name:     "data.schema.json",
						},
					},
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/action",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/action",
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
func (rnr *Runner) PostActionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
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
		err := coreerror.NewPathParamError("dungeon_id", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	if characterID == "" {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(characterID) {
		err := coreerror.NewPathParamError("character_id", characterID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting dungeon record ID >%s<", dungeonID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, nil)
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

	characterRec, err := m.(*model.Model).GetCharacterRec(characterID, nil)
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

	req := &schema.ActionRequest{}

	req, err = server.ReadRequest(l, r, req)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	sentence := strings.ToLower(req.Data.Sentence)

	l.Info("Verifying character instance for dungeon_id >%s< character_id >%s<", dungeonID, characterID)

	characterInstanceRecs, err := m.(*model.Model).GetCharacterInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "character_id",
					Val: characterID,
				},
			},
		},
	)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	if len(characterInstanceRecs) == 0 {
		err := coreerror.NewPathParamError("character_id", characterID)
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

	dungeonInstanceRec, err := m.(*model.Model).GetDungeonInstanceRec(characterInstanceRec.DungeonInstanceID, nil)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	if dungeonInstanceRec == nil {
		err := coreerror.NewPathParamError("dungeon_id", dungeonID)
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

	// Get every action that occured between this characters previous action and this action
	actionRecs, err := m.(*model.Model).GetActionRecsSincePreviousAction(dungeonActionRecordSet.ActionRec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Responding with >%d< action records", len(actionRecs))

	// Response data
	data := []schema.ActionResponseData{}
	for _, rec := range actionRecs {

		rs, err := m.(*model.Model).GetActionRecordSet(rec.ID)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		responseData, err := actionResponseData(l, *rs)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, *responseData)
	}

	// Assign response properties
	res := schema.ActionResponse{
		Data: data,
	}

	for idx := range res.Data {
		if res.Data[idx].ActionCharacter != nil {
			l.Info("Response ActionCharacter >%s<", res.Data[idx].ActionCharacter.CharacterName)
			for oidx := range res.Data[idx].ActionCharacter.CharacterEquippedObjects {
				l.Info("         -   >%#v<", res.Data[idx].ActionCharacter.CharacterEquippedObjects[oidx])
			}
		}
		if res.Data[idx].ActionMonster != nil {
			l.Info("Response ActionMonster >%s<", res.Data[idx].ActionMonster.MonsterName)
		}
		l.Info("Response ActionLocation LocationName >%s<", res.Data[idx].ActionLocation.LocationName)
		l.Info("Response - Action ID >%s<", res.Data[idx].ActionID)
		l.Info("Response - Action Command >%s<", res.Data[idx].ActionCommand)
		l.Info("Response - Action TurnNumber >%d<", res.Data[idx].ActionTurnNumber)
		l.Info("Response - Action SerialNumber >%d<", res.Data[idx].ActionSerialNumber)
	}

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
