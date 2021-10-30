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

	dungeonActionRec := dungeonActionRecordSet.DungeonActionRec

	dungeonLocationRec := dungeonActionRecordSet.DungeonLocationRec
	dungeonLocations := []string{}
	if dungeonLocationRec.NorthDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "north")
	}
	if dungeonLocationRec.NortheastDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "northeast")
	}
	if dungeonLocationRec.EastDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "east")
	}
	if dungeonLocationRec.SoutheastDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "southeast")
	}
	if dungeonLocationRec.SouthDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "south")
	}
	if dungeonLocationRec.SouthwestDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "southwest")
	}
	if dungeonLocationRec.SouthwestDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "southwest")
	}
	if dungeonLocationRec.WestDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "west")
	}
	if dungeonLocationRec.NorthwestDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "northwest")
	}
	if dungeonLocationRec.UpDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "up")
	}
	if dungeonLocationRec.DownDungeonLocationID.Valid {
		dungeonLocations = append(dungeonLocations, "down")
	}

	var characterData schema.CharacterData
	if dungeonActionRecordSet.DungeonCharacterRec != nil {
		characterData = schema.CharacterData{
			Name: dungeonActionRecordSet.DungeonCharacterRec.Name,
		}
	}

	var monsterData schema.MonsterData
	if dungeonActionRecordSet.DungeonMonsterRec != nil {
		monsterData = schema.MonsterData{
			Name: dungeonActionRecordSet.DungeonMonsterRec.Name,
		}
	}

	var dungeonCharacterData []schema.CharacterData
	if len(dungeonActionRecordSet.DungeonCharacterRecs) > 0 {
		for _, dungeonCharacterRec := range dungeonActionRecordSet.DungeonCharacterRecs {
			dungeonCharacterData = append(dungeonCharacterData,
				schema.CharacterData{
					Name: dungeonCharacterRec.Name,
				})
		}
	}

	var dungeonMonsterData []schema.MonsterData
	if len(dungeonActionRecordSet.DungeonMonsterRecs) > 0 {
		for _, dungeonMonsterRec := range dungeonActionRecordSet.DungeonMonsterRecs {
			dungeonMonsterData = append(dungeonMonsterData,
				schema.MonsterData{
					Name: dungeonMonsterRec.Name,
				})
		}
	}

	var dungeonObjectData []schema.ObjectData
	if len(dungeonActionRecordSet.DungeonObjectRecs) > 0 {
		for _, dungeonObjectRec := range dungeonActionRecordSet.DungeonObjectRecs {
			dungeonObjectData = append(dungeonObjectData,
				schema.ObjectData{
					Name: dungeonObjectRec.Name,
				})
		}
	}

	data := schema.DungeonActionResponseData{
		Action: schema.ActionData{
			ID:                             dungeonActionRec.ID,
			Command:                        dungeonActionRec.ResolvedCommand,
			EquippedDungeonObjectName:      dungeonActionRec.ResolvedEquippedDungeonObjectName.String,
			StashedDungeonObjectName:       dungeonActionRec.ResolvedStashedDungeonObjectName.String,
			TargetDungeonObjectName:        dungeonActionRec.ResolvedTargetDungeonObjectName.String,
			TargetDungeonCharacterName:     dungeonActionRec.ResolvedTargetDungeonCharacterName.String,
			TargetDungeonMonsterName:       dungeonActionRec.ResolvedTargetDungeonMonsterName.String,
			TargetDungeonLocationName:      dungeonActionRec.ResolvedTargetDungeonLocationName.String,
			TargetDungeonLocationDirection: dungeonActionRec.ResolvedTargetDungeonLocationDirection.String,
		},
		Location: schema.LocationData{
			Name:        dungeonLocationRec.Name,
			Description: dungeonLocationRec.Description,
			Directions:  dungeonLocations,
		},
		Character:  characterData,
		Monster:    monsterData,
		Characters: dungeonCharacterData,
		Monsters:   dungeonMonsterData,
		Objects:    dungeonObjectData,
	}

	return data, nil
}
