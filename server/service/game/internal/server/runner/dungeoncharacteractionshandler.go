package runner

import (
	"net/http"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"

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

	l.Debug("Resulting action >%#v<", dungeonActionRecordSet.DungeonActionRec)
	l.Debug("Resulting action current location >%#v<", dungeonActionRecordSet.CurrentLocation)
	l.Debug("Resulting action target location >%#v<", dungeonActionRecordSet.TargetLocation)
	l.Debug("Resulting action character >%#v<", dungeonActionRecordSet.DungeonCharacterRec)
	l.Debug("Resulting action monster >%#v<", dungeonActionRecordSet.DungeonMonsterRec)

	// Response data
	responseData, err := rnr.RecordToDungeonCharacterActionResponseData(*dungeonActionRecordSet)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.DungeonActionResponse{
		Data: []schema.DungeonActionResponseData{
			*responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// RecordToCharacterResponseData -
func (rnr *Runner) RecordToDungeonCharacterActionResponseData(dungeonActionRecordSet model.DungeonActionRecordSet) (*schema.DungeonActionResponseData, error) {

	dungeonActionRec := dungeonActionRecordSet.DungeonActionRec

	var characterData *schema.CharacterData
	if dungeonActionRecordSet.DungeonCharacterRec != nil {
		characterData = &schema.CharacterData{
			Name: dungeonActionRecordSet.DungeonCharacterRec.Name,
		}
	}

	var monsterData *schema.MonsterData
	if dungeonActionRecordSet.DungeonMonsterRec != nil {
		monsterData = &schema.MonsterData{
			Name: dungeonActionRecordSet.DungeonMonsterRec.Name,
		}
	}

	// Current location
	locationData, err := rnr.dungeonActionLocationToResponseLocation(dungeonActionRecordSet.CurrentLocation)
	if err != nil {
		return nil, err
	}

	// Target location
	targetLocationData := &schema.LocationData{}
	if dungeonActionRecordSet.TargetLocation != nil {
		targetLocationData, err = rnr.dungeonActionLocationToResponseLocation(dungeonActionRecordSet.TargetLocation)
		if err != nil {
			return nil, err
		}
	}

	if dungeonActionRec.ResolvedTargetDungeonLocationDirection.Valid {
		targetLocationData.Direction = dungeonActionRec.ResolvedTargetDungeonLocationDirection.String
	}

	// Equipped object
	equippedObjectData := &schema.ObjectData{}
	if dungeonActionRecordSet.EquippedObjectRec != nil {
		equippedObjectData, err = rnr.dungeonObjectToResponseObject(dungeonActionRecordSet.EquippedObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Stashed object
	stashedObjectData := &schema.ObjectData{}
	if dungeonActionRecordSet.StashedObjectRec != nil {
		stashedObjectData, err = rnr.dungeonObjectToResponseObject(dungeonActionRecordSet.StashedObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Target object
	targetObjectData := &schema.ObjectData{}
	if dungeonActionRecordSet.TargetObjectRec != nil {
		targetObjectData, err = rnr.dungeonObjectToResponseObject(dungeonActionRecordSet.TargetObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Target character
	targetCharacterData := &schema.CharacterData{}
	if dungeonActionRecordSet.TargetCharacterRec != nil {
		targetCharacterData, err = rnr.dungeonCharacterToResponseCharacter(dungeonActionRecordSet.TargetCharacterRec)
		if err != nil {
			return nil, err
		}
	}

	// Target monster
	targetMonsterData := &schema.MonsterData{}
	if dungeonActionRecordSet.TargetMonsterRec != nil {
		targetMonsterData, err = rnr.dungeonMonsterToResponseMonster(dungeonActionRecordSet.TargetMonsterRec)
		if err != nil {
			return nil, err
		}
	}

	data := schema.DungeonActionResponseData{
		ID:              dungeonActionRec.ID,
		Command:         dungeonActionRec.ResolvedCommand,
		Location:        *locationData,
		Character:       characterData,
		Monster:         monsterData,
		EquippedObject:  equippedObjectData,
		StashedObject:   stashedObjectData,
		TargetObject:    targetObjectData,
		TargetCharacter: targetCharacterData,
		TargetMonster:   targetMonsterData,
		TargetLocation:  targetLocationData,
		CreatedAt:       dungeonActionRec.CreatedAt,
		UpdatedAt:       dungeonActionRec.UpdatedAt.Time,
	}

	return &data, nil
}

func (rnr *Runner) dungeonObjectToResponseObject(dungeonObjectRec *record.DungeonObject) (*schema.ObjectData, error) {
	return &schema.ObjectData{
		Name:        dungeonObjectRec.Name,
		Description: dungeonObjectRec.Description,
	}, nil
}

func (rnr *Runner) dungeonCharacterToResponseCharacter(dungeonCharacterRec *record.DungeonCharacter) (*schema.CharacterData, error) {
	return &schema.CharacterData{
		Name: dungeonCharacterRec.Name,
	}, nil
}

func (rnr *Runner) dungeonMonsterToResponseMonster(dungeonMonsterRec *record.DungeonMonster) (*schema.MonsterData, error) {
	return &schema.MonsterData{
		Name: dungeonMonsterRec.Name,
	}, nil
}

// actionLocationToReponseLocation -
func (rnr *Runner) dungeonActionLocationToResponseLocation(recordSet *model.DungeonActionLocationRecordSet) (*schema.LocationData, error) {

	dungeonLocationRec := recordSet.DungeonLocationRec

	directions := []string{}
	if dungeonLocationRec.NorthDungeonLocationID.Valid {
		directions = append(directions, "north")
	}
	if dungeonLocationRec.NortheastDungeonLocationID.Valid {
		directions = append(directions, "northeast")
	}
	if dungeonLocationRec.EastDungeonLocationID.Valid {
		directions = append(directions, "east")
	}
	if dungeonLocationRec.SoutheastDungeonLocationID.Valid {
		directions = append(directions, "southeast")
	}
	if dungeonLocationRec.SouthDungeonLocationID.Valid {
		directions = append(directions, "south")
	}
	if dungeonLocationRec.SouthwestDungeonLocationID.Valid {
		directions = append(directions, "southwest")
	}
	if dungeonLocationRec.SouthwestDungeonLocationID.Valid {
		directions = append(directions, "southwest")
	}
	if dungeonLocationRec.WestDungeonLocationID.Valid {
		directions = append(directions, "west")
	}
	if dungeonLocationRec.NorthwestDungeonLocationID.Valid {
		directions = append(directions, "northwest")
	}
	if dungeonLocationRec.UpDungeonLocationID.Valid {
		directions = append(directions, "up")
	}
	if dungeonLocationRec.DownDungeonLocationID.Valid {
		directions = append(directions, "down")
	}

	var charactersData []schema.CharacterData
	if len(recordSet.DungeonCharacterRecs) > 0 {
		for _, dungeonCharacterRec := range recordSet.DungeonCharacterRecs {
			charactersData = append(charactersData,
				schema.CharacterData{
					Name: dungeonCharacterRec.Name,
				})
		}
	}

	var monstersData []schema.MonsterData
	if len(recordSet.DungeonMonsterRecs) > 0 {
		for _, dungeonMonsterRec := range recordSet.DungeonMonsterRecs {
			monstersData = append(monstersData,
				schema.MonsterData{
					Name: dungeonMonsterRec.Name,
				})
		}
	}

	var objectsData []schema.ObjectData
	if len(recordSet.DungeonObjectRecs) > 0 {
		for _, dungeonObjectRec := range recordSet.DungeonObjectRecs {
			objectsData = append(objectsData,
				schema.ObjectData{
					Name: dungeonObjectRec.Name,
				})
		}
	}

	data := &schema.LocationData{
		Name:        dungeonLocationRec.Name,
		Description: dungeonLocationRec.Description,
		//		Direction:
		Directions: directions,
		Characters: charactersData,
		Monsters:   monstersData,
		Objects:    objectsData,
	}

	return data, nil
}
