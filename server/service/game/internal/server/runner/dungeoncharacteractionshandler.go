package runner

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/mapper"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
)

// PostDungeonCharacterActionsHandler -
func (rnr *Runner) PostDungeonCharacterActionsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Post characters handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	req := schema.ActionRequest{}

	err := rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	req.Data.Sentence = strings.ToLower(req.Data.Sentence)

	l.Info("Creating dungeon character action >%s<", req.Data.Sentence)

	dungeonActionRecordSet, err := m.(*model.Model).ProcessDungeonCharacterAction(dungeonID, characterID, req.Data.Sentence)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	l.Debug("Resulting action >%#v<", dungeonActionRecordSet.ActionRec)
	l.Debug("Resulting action current location >%#v<", dungeonActionRecordSet.CurrentLocation)
	l.Debug("Resulting action target location >%#v<", dungeonActionRecordSet.TargetLocation)
	l.Debug("Resulting action character >%#v<", dungeonActionRecordSet.ActionCharacterRec)
	l.Debug("Resulting action monster >%#v<", dungeonActionRecordSet.ActionMonsterRec)

	// Response data
	responseData, err := mapper.ActionRecordSetToActionResponse(l, *dungeonActionRecordSet)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
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

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}
