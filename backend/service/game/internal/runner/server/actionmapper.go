package runner

import (
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/nullint"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/schema"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// actionResponseData -
func actionResponseData(l logger.Logger, actionRecordSet record.ActionRecordSet) (*schema.ActionResponseData, error) {

	dungeonActionRec := actionRecordSet.ActionRec

	var err error
	var characterData *schema.ActionCharacter
	if actionRecordSet.ActionCharacterRec != nil {
		characterData, err = actionCharacterResponseData(
			l,
			actionRecordSet.ActionCharacterRec,
			actionRecordSet.ActionCharacterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	var monsterData *schema.ActionMonster
	if actionRecordSet.ActionMonsterRec != nil {
		monsterData, err = actionMonsterResponseData(
			l,
			actionRecordSet.ActionMonsterRec,
			actionRecordSet.ActionMonsterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	// Current location
	locationData, err := actionLocationResponseData(l, actionRecordSet.CurrentLocation)
	if err != nil {
		return nil, err
	}

	// Target location
	var targetActionLocation *schema.ActionLocation
	if actionRecordSet.TargetLocation != nil {
		targetActionLocation, err = actionLocationResponseData(
			l,
			actionRecordSet.TargetLocation,
		)
		if err != nil {
			return nil, err
		}
	}

	if dungeonActionRec.ResolvedTargetLocationDirection.Valid {
		targetActionLocation.LocationDirection = dungeonActionRec.ResolvedTargetLocationDirection.String
	}

	// Equipped object
	var equippedActionLocationObject *schema.ActionObject
	if actionRecordSet.EquippedActionObjectRec != nil {
		equippedActionLocationObject, err = actionObjectResponseData(
			l,
			actionRecordSet.EquippedActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Stashed object
	var stashedActionLocationObject *schema.ActionObject
	if actionRecordSet.StashedActionObjectRec != nil {
		stashedActionLocationObject, err = actionObjectResponseData(
			l,
			actionRecordSet.StashedActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Dropped object
	var droppedActionLocationObject *schema.ActionObject
	if actionRecordSet.DroppedActionObjectRec != nil {
		droppedActionLocationObject, err = actionObjectResponseData(
			l,
			actionRecordSet.DroppedActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Target object
	var targetActionLocationObject *schema.ActionObject
	if actionRecordSet.TargetActionObjectRec != nil {
		targetActionLocationObject, err = actionObjectResponseData(
			l,
			actionRecordSet.TargetActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Target character
	var targetActionLocationCharacter *schema.ActionCharacter
	if actionRecordSet.TargetActionCharacterRec != nil {
		targetActionLocationCharacter, err = actionTargetCharacterResponseData(
			l,
			actionRecordSet.TargetActionCharacterRec,
			actionRecordSet.TargetActionCharacterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	// Target monster
	var targetActionLocationMonster *schema.ActionMonster
	if actionRecordSet.TargetActionMonsterRec != nil {
		targetActionLocationMonster, err = actionTargetMonsterResponseData(
			l,
			actionRecordSet.TargetActionMonsterRec,
			actionRecordSet.TargetActionMonsterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	narrative, err := actionNarrativeResponseData(
		l,
		actionRecordSet,
	)
	if err != nil {
		return nil, err
	}

	data := schema.ActionResponseData{
		ActionID:              dungeonActionRec.ID,
		ActionCommand:         dungeonActionRec.ResolvedCommand,
		ActionTurnNumber:      dungeonActionRec.TurnNumber,
		ActionSerialNumber:    nullint.ToInt16(dungeonActionRec.SerialNumber),
		ActionNarrative:       narrative,
		ActionLocation:        *locationData,
		ActionCharacter:       characterData,
		ActionMonster:         monsterData,
		ActionEquippedObject:  equippedActionLocationObject,
		ActionStashedObject:   stashedActionLocationObject,
		ActionDroppedObject:   droppedActionLocationObject,
		ActionTargetObject:    targetActionLocationObject,
		ActionTargetCharacter: targetActionLocationCharacter,
		ActionTargetMonster:   targetActionLocationMonster,
		ActionTargetLocation:  targetActionLocation,
		ActionCreatedAt:       dungeonActionRec.CreatedAt,
	}

	return &data, nil
}

func actionObjectResponseData(l logger.Logger, dungeonObjectRec *record.ActionObject) (*schema.ActionObject, error) {
	return &schema.ActionObject{
		ObjectName:        dungeonObjectRec.Name,
		ObjectDescription: dungeonObjectRec.Description,
		ObjectIsEquipped:  dungeonObjectRec.IsEquipped,
		ObjectIsStashed:   dungeonObjectRec.IsStashed,
	}, nil
}

func actionCharacterResponseData(
	l logger.Logger,
	characterRec *record.ActionCharacter,
	objectRecs []*record.ActionCharacterObject,
) (*schema.ActionCharacter, error) {

	data := &schema.ActionCharacter{
		CharacterName:                characterRec.Name,
		CharacterStrength:            characterRec.Strength,
		CharacterDexterity:           characterRec.Dexterity,
		CharacterIntelligence:        characterRec.Intelligence,
		CharacterCurrentStrength:     characterRec.CurrentStrength,
		CharacterCurrentDexterity:    characterRec.CurrentDexterity,
		CharacterCurrentIntelligence: characterRec.CurrentIntelligence,
		CharacterHealth:              characterRec.Health,
		CharacterFatigue:             characterRec.Fatigue,
		CharacterCurrentHealth:       characterRec.CurrentHealth,
		CharacterCurrentFatigue:      characterRec.CurrentFatigue,
		CharacterStashedObjects:      []schema.ActionObject{},
		CharacterEquippedObjects:     []schema.ActionObject{},
	}

	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			l.Debug("Adding character equipped object >%#v<", objectRec)
			data.CharacterEquippedObjects = append(data.CharacterEquippedObjects, schema.ActionObject{
				ObjectName:       objectRec.Name,
				ObjectIsEquipped: objectRec.IsEquipped,
				ObjectIsStashed:  objectRec.IsStashed,
			})
			continue
		}
		if objectRec.IsStashed {
			l.Debug("Adding character stashed object >%#v<", objectRec)
			data.CharacterStashedObjects = append(data.CharacterStashedObjects, schema.ActionObject{
				ObjectName:       objectRec.Name,
				ObjectIsEquipped: objectRec.IsEquipped,
				ObjectIsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

func actionTargetCharacterResponseData(
	l logger.Logger,
	characterRec *record.ActionCharacter,
	objectRecs []*record.ActionCharacterObject,
) (*schema.ActionCharacter, error) {
	data := &schema.ActionCharacter{
		CharacterName:                characterRec.Name,
		CharacterStrength:            characterRec.Strength,
		CharacterDexterity:           characterRec.Dexterity,
		CharacterIntelligence:        characterRec.Intelligence,
		CharacterCurrentStrength:     characterRec.CurrentStrength,
		CharacterCurrentDexterity:    characterRec.CurrentDexterity,
		CharacterCurrentIntelligence: characterRec.CurrentIntelligence,
		CharacterHealth:              characterRec.Health,
		CharacterFatigue:             characterRec.Fatigue,
		CharacterCurrentHealth:       characterRec.CurrentHealth,
		CharacterCurrentFatigue:      characterRec.CurrentFatigue,
		CharacterEquippedObjects:     []schema.ActionObject{},
	}

	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			l.Debug("Adding target character equipped object >%#v<", objectRec)
			data.CharacterEquippedObjects = append(data.CharacterEquippedObjects, schema.ActionObject{
				ObjectName:       objectRec.Name,
				ObjectIsEquipped: objectRec.IsEquipped,
				ObjectIsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

func actionMonsterResponseData(
	l logger.Logger,
	dungeonMonsterRec *record.ActionMonster,
	objectRecs []*record.ActionMonsterObject,
) (*schema.ActionMonster, error) {

	data := &schema.ActionMonster{
		MonsterName:                dungeonMonsterRec.Name,
		MonsterStrength:            dungeonMonsterRec.Strength,
		MonsterDexterity:           dungeonMonsterRec.Dexterity,
		MonsterIntelligence:        dungeonMonsterRec.Intelligence,
		MonsterCurrentStrength:     dungeonMonsterRec.CurrentStrength,
		MonsterCurrentDexterity:    dungeonMonsterRec.CurrentDexterity,
		MonsterCurrentIntelligence: dungeonMonsterRec.CurrentIntelligence,
		MonsterHealth:              dungeonMonsterRec.Health,
		MonsterFatigue:             dungeonMonsterRec.Fatigue,
		MonsterCurrentHealth:       dungeonMonsterRec.CurrentHealth,
		MonsterCurrentFatigue:      dungeonMonsterRec.CurrentFatigue,
		MonsterEquippedObjects:     []schema.ActionObject{},
	}

	// NOTE: We only ever provide a monsters equipped objects and never provide
	// a monsters stashed objects.
	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			l.Debug("Adding monster equipped object >%#v<", objectRec)
			data.MonsterEquippedObjects = append(data.MonsterEquippedObjects, schema.ActionObject{
				ObjectName:       objectRec.Name,
				ObjectIsEquipped: objectRec.IsEquipped,
				ObjectIsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

func actionTargetMonsterResponseData(
	l logger.Logger,
	monsterRec *record.ActionMonster,
	objectRecs []*record.ActionMonsterObject,
) (*schema.ActionMonster, error) {
	data := &schema.ActionMonster{
		MonsterName:                monsterRec.Name,
		MonsterStrength:            monsterRec.Strength,
		MonsterDexterity:           monsterRec.Dexterity,
		MonsterIntelligence:        monsterRec.Intelligence,
		MonsterCurrentStrength:     monsterRec.CurrentStrength,
		MonsterCurrentDexterity:    monsterRec.CurrentDexterity,
		MonsterCurrentIntelligence: monsterRec.CurrentIntelligence,
		MonsterHealth:              monsterRec.Health,
		MonsterFatigue:             monsterRec.Fatigue,
		MonsterCurrentHealth:       monsterRec.CurrentHealth,
		MonsterCurrentFatigue:      monsterRec.CurrentFatigue,
		MonsterEquippedObjects:     []schema.ActionObject{},
	}

	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			data.MonsterEquippedObjects = append(data.MonsterEquippedObjects, schema.ActionObject{
				ObjectName:       objectRec.Name,
				ObjectIsEquipped: objectRec.IsEquipped,
				ObjectIsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

// actionLocationResponseData -
func actionLocationResponseData(l logger.Logger, recordSet *record.ActionLocationRecordSet) (*schema.ActionLocation, error) {

	dungeonLocationRec := recordSet.LocationInstanceViewRec

	directions := []string{}
	if dungeonLocationRec.NorthLocationInstanceID.Valid {
		directions = append(directions, "north")
	}
	if dungeonLocationRec.NortheastLocationInstanceID.Valid {
		directions = append(directions, "northeast")
	}
	if dungeonLocationRec.EastLocationInstanceID.Valid {
		directions = append(directions, "east")
	}
	if dungeonLocationRec.SoutheastLocationInstanceID.Valid {
		directions = append(directions, "southeast")
	}
	if dungeonLocationRec.SouthLocationInstanceID.Valid {
		directions = append(directions, "south")
	}
	if dungeonLocationRec.SouthwestLocationInstanceID.Valid {
		directions = append(directions, "southwest")
	}
	if dungeonLocationRec.SouthwestLocationInstanceID.Valid {
		directions = append(directions, "southwest")
	}
	if dungeonLocationRec.WestLocationInstanceID.Valid {
		directions = append(directions, "west")
	}
	if dungeonLocationRec.NorthwestLocationInstanceID.Valid {
		directions = append(directions, "northwest")
	}
	if dungeonLocationRec.UpLocationInstanceID.Valid {
		directions = append(directions, "up")
	}
	if dungeonLocationRec.DownLocationInstanceID.Valid {
		directions = append(directions, "down")
	}

	var charactersData []schema.ActionLocationCharacter
	if len(recordSet.ActionCharacterRecs) > 0 {
		for _, dungeonCharacterRec := range recordSet.ActionCharacterRecs {
			charactersData = append(charactersData,
				schema.ActionLocationCharacter{
					CharacterName: dungeonCharacterRec.Name,
				})
		}
	}

	var monstersData []schema.ActionLocationMonster
	if len(recordSet.ActionMonsterRecs) > 0 {
		for _, dungeonMonsterRec := range recordSet.ActionMonsterRecs {
			monstersData = append(monstersData,
				schema.ActionLocationMonster{
					MonsterName: dungeonMonsterRec.Name,
				})
		}
	}

	var objectsData []schema.ActionLocationObject
	if len(recordSet.ActionObjectRecs) > 0 {
		for _, dungeonObjectRec := range recordSet.ActionObjectRecs {
			objectsData = append(objectsData,
				schema.ActionLocationObject{
					ObjectName: dungeonObjectRec.Name,
				})
		}
	}

	data := &schema.ActionLocation{
		LocationName:        dungeonLocationRec.Name,
		LocationDescription: dungeonLocationRec.Description,
		LocationDirections:  directions,
		LocationCharacters:  charactersData,
		LocationMonsters:    monstersData,
		LocationObjects:     objectsData,
	}

	return data, nil
}

// actionNarrativeResponseData -
func actionNarrativeResponseData(l logger.Logger, set record.ActionRecordSet) (string, error) {

	desc := ""
	if set.ActionCharacterRec != nil {
		desc += set.ActionCharacterRec.Name
	} else if set.ActionMonsterRec != nil {
		desc += set.ActionMonsterRec.Name
	}

	switch set.ActionRec.ResolvedCommand {
	case record.ActionCommandMove:
		desc += " moves "
	case record.ActionCommandLook:
		desc += " looks "
	case record.ActionCommandStash:
		desc += " stashes "
	case record.ActionCommandEquip:
		desc += " equips "
	case record.ActionCommandDrop:
		desc += " drops "
	default:
		// no-op
	}

	if set.TargetActionCharacterRec != nil {
		desc += set.TargetActionCharacterRec.Name
	} else if set.TargetActionMonsterRec != nil {
		desc += set.TargetActionMonsterRec.Name
	} else if set.TargetActionObjectRec != nil {
		desc += set.TargetActionObjectRec.Name
	} else if set.TargetLocation != nil {
		desc += set.ActionRec.ResolvedTargetLocationDirection.String
	}

	desc = strings.TrimRight(desc, " ")

	return desc, nil
}
