package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
)

// PostDungeonCharacterActionsHandler -
func (rnr *Runner) PostDungeonCharacterActionsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Post characters handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	req := schema.DungeonActionRequest{}

	err := rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	l.Info("Creating dungeon character action >%s<", req.Data.Sentence)

	dungeonActionRecordSet, err := m.(*model.Model).ProcessDungeonCharacterAction(dungeonID, characterID, req.Data.Sentence)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Response data
	responseData, err := rnr.RecordToDungeonCharacterActionResponseData(dungeonActionRecordSet)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.DungeonActionResponse{
		Data: []schema.DungeonActionResponseData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// RecordToCharacterResponseData -
func (rnr *Runner) RecordToDungeonCharacterActionResponseData(dungeonActionRecordSet *model.DungeonActionRecordSet) (schema.DungeonActionResponseData, error) {

	data := schema.DungeonActionResponseData{
		Action: schema.ActionData{
			ID: dungeonActionRecordSet.DungeonActionRec.ID,
		},
	}

	return data, nil
}
