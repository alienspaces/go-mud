package runner

import (
	"net/http"
	"strings"

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
	responseData, err := rnr.RecordToDungeonActionCharacterActionResponseData(*dungeonActionRecordSet)
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
func (rnr *Runner) RecordToDungeonActionCharacterActionResponseData(dungeonActionRecordSet model.DungeonActionRecordSet) (*schema.DungeonActionResponseData, error) {

	dungeonActionRec := dungeonActionRecordSet.ActionRec

	var err error
	var characterData *schema.CharacterDetailedData
	if dungeonActionRecordSet.ActionCharacterRec != nil {
		characterData, err = rnr.dungeonCharacterToResponseCharacter(dungeonActionRecordSet.ActionCharacterRec)
		if err != nil {
			return nil, err
		}
	}

	var monsterData *schema.MonsterDetailedData
	if dungeonActionRecordSet.ActionMonsterRec != nil {
		monsterData, err = rnr.dungeonMonsterToResponseMonster(dungeonActionRecordSet.ActionMonsterRec)
		if err != nil {
			return nil, err
		}
	}

	// Current location
	locationData, err := rnr.dungeonActionLocationToResponseLocation(dungeonActionRecordSet.CurrentLocation)
	if err != nil {
		return nil, err
	}

	// Target location
	var targetLocationData *schema.LocationData
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
	var equippedObjectData *schema.ObjectDetailedData
	if dungeonActionRecordSet.EquippedActionObjectRec != nil {
		equippedObjectData, err = rnr.dungeonObjectToResponseObject(dungeonActionRecordSet.EquippedActionObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Stashed object
	var stashedObjectData *schema.ObjectDetailedData
	if dungeonActionRecordSet.StashedActionObjectRec != nil {
		stashedObjectData, err = rnr.dungeonObjectToResponseObject(dungeonActionRecordSet.StashedActionObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Target object
	var targetObjectData *schema.ObjectDetailedData
	if dungeonActionRecordSet.TargetActionObjectRec != nil {
		targetObjectData, err = rnr.dungeonObjectToResponseObject(dungeonActionRecordSet.TargetActionObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Target character
	var targetCharacterData *schema.CharacterDetailedData
	if dungeonActionRecordSet.TargetActionCharacterRec != nil {
		targetCharacterData, err = rnr.dungeonCharacterToResponseCharacter(dungeonActionRecordSet.TargetActionCharacterRec)
		if err != nil {
			return nil, err
		}
	}

	// Target monster
	var targetMonsterData *schema.MonsterDetailedData
	if dungeonActionRecordSet.TargetActionMonsterRec != nil {
		targetMonsterData, err = rnr.dungeonMonsterToResponseMonster(dungeonActionRecordSet.TargetActionMonsterRec)
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
	}

	return &data, nil
}

func (rnr *Runner) dungeonObjectToResponseObject(dungeonObjectRec *record.DungeonActionObject) (*schema.ObjectDetailedData, error) {
	return &schema.ObjectDetailedData{
		Name: dungeonObjectRec.Name,
	}, nil
}

func (rnr *Runner) dungeonCharacterToResponseCharacter(dungeonCharacterRec *record.DungeonActionCharacter) (*schema.CharacterDetailedData, error) {
	return &schema.CharacterDetailedData{
		Name:         dungeonCharacterRec.Name,
		Strength:     dungeonCharacterRec.Strength,
		Dexterity:    dungeonCharacterRec.Dexterity,
		Intelligence: dungeonCharacterRec.Intelligence,
		Health:       dungeonCharacterRec.Health,
		Fatigue:      dungeonCharacterRec.Fatigue,
	}, nil
}

func (rnr *Runner) dungeonMonsterToResponseMonster(dungeonMonsterRec *record.DungeonActionMonster) (*schema.MonsterDetailedData, error) {
	return &schema.MonsterDetailedData{
		Name:         dungeonMonsterRec.Name,
		Strength:     dungeonMonsterRec.Strength,
		Dexterity:    dungeonMonsterRec.Dexterity,
		Intelligence: dungeonMonsterRec.Intelligence,
		Health:       dungeonMonsterRec.Health,
		Fatigue:      dungeonMonsterRec.Fatigue,
	}, nil
}

// actionLocationToReponseLocation -
func (rnr *Runner) dungeonActionLocationToResponseLocation(recordSet *model.DungeonActionLocationRecordSet) (*schema.LocationData, error) {

	dungeonLocationRec := recordSet.LocationRec

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
	if len(recordSet.ActionCharacterRecs) > 0 {
		for _, dungeonCharacterRec := range recordSet.ActionCharacterRecs {
			charactersData = append(charactersData,
				schema.CharacterData{
					Name: dungeonCharacterRec.Name,
				})
		}
	}

	var monstersData []schema.MonsterData
	if len(recordSet.ActionMonsterRecs) > 0 {
		for _, dungeonMonsterRec := range recordSet.ActionMonsterRecs {
			monstersData = append(monstersData,
				schema.MonsterData{
					Name: dungeonMonsterRec.Name,
				})
		}
	}

	var objectsData []schema.ObjectData
	if len(recordSet.ActionObjectRecs) > 0 {
		for _, dungeonObjectRec := range recordSet.ActionObjectRecs {
			objectsData = append(objectsData,
				schema.ObjectData{
					Name: dungeonObjectRec.Name,
				})
		}
	}

	data := &schema.LocationData{
		Name:        dungeonLocationRec.Name,
		Description: dungeonLocationRec.Description,
		Directions:  directions,
		Characters:  charactersData,
		Monsters:    monstersData,
		Objects:     objectsData,
	}

	return data, nil
}
