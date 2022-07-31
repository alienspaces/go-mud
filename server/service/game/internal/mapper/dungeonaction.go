package mapper

import (
	"strings"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ActionRecordSetToActionResponse -
func ActionRecordSetToActionResponse(l logger.Logger, actionRecordSet record.ActionRecordSet) (*schema.ActionResponseData, error) {

	dungeonActionRec := actionRecordSet.ActionRec

	var err error
	var characterData *schema.ActionCharacter
	if actionRecordSet.ActionCharacterRec != nil {
		characterData, err = dungeonCharacterToResponseCharacter(
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
		monsterData, err = dungeonMonsterToResponseMonster(
			l,
			actionRecordSet.ActionMonsterRec,
			actionRecordSet.ActionMonsterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	// Current location
	locationData, err := actionLocationToResponseLocation(actionRecordSet.CurrentLocation)
	if err != nil {
		return nil, err
	}

	// Target location
	var targetActionLocation *schema.ActionLocation
	if actionRecordSet.TargetLocation != nil {
		targetActionLocation, err = actionLocationToResponseLocation(actionRecordSet.TargetLocation)
		if err != nil {
			return nil, err
		}
	}

	if dungeonActionRec.ResolvedTargetLocationDirection.Valid {
		targetActionLocation.Direction = dungeonActionRec.ResolvedTargetLocationDirection.String
	}

	// Equipped object
	var equippedActionLocationObject *schema.ActionObject
	if actionRecordSet.EquippedActionObjectRec != nil {
		equippedActionLocationObject, err = dungeonObjectToResponseObject(actionRecordSet.EquippedActionObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Stashed object
	var stashedActionLocationObject *schema.ActionObject
	if actionRecordSet.StashedActionObjectRec != nil {
		stashedActionLocationObject, err = dungeonObjectToResponseObject(actionRecordSet.StashedActionObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Dropped object
	var droppedActionLocationObject *schema.ActionObject
	if actionRecordSet.DroppedActionObjectRec != nil {
		droppedActionLocationObject, err = dungeonObjectToResponseObject(actionRecordSet.DroppedActionObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Target object
	var targetActionLocationObject *schema.ActionObject
	if actionRecordSet.TargetActionObjectRec != nil {
		targetActionLocationObject, err = dungeonObjectToResponseObject(actionRecordSet.TargetActionObjectRec)
		if err != nil {
			return nil, err
		}
	}

	// Target character
	var targetActionLocationCharacter *schema.ActionCharacter
	if actionRecordSet.TargetActionCharacterRec != nil {
		targetActionLocationCharacter, err = dungeonTargetCharacterToResponseTargetCharacter(
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
		targetActionLocationMonster, err = dungeonTargetMonsterToResponseTargetMonster(
			l,
			actionRecordSet.TargetActionMonsterRec,
			actionRecordSet.TargetActionMonsterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	narrative, err := dungeonActionToNarrative(l, actionRecordSet)
	if err != nil {
		return nil, err
	}

	data := schema.ActionResponseData{
		ID:              dungeonActionRec.ID,
		Command:         dungeonActionRec.ResolvedCommand,
		Narrative:       narrative,
		Location:        *locationData,
		Character:       characterData,
		Monster:         monsterData,
		EquippedObject:  equippedActionLocationObject,
		StashedObject:   stashedActionLocationObject,
		DroppedObject:   droppedActionLocationObject,
		TargetObject:    targetActionLocationObject,
		TargetCharacter: targetActionLocationCharacter,
		TargetMonster:   targetActionLocationMonster,
		TargetLocation:  targetActionLocation,
		CreatedAt:       dungeonActionRec.CreatedAt,
	}

	return &data, nil
}

func dungeonObjectToResponseObject(dungeonObjectRec *record.ActionObject) (*schema.ActionObject, error) {
	return &schema.ActionObject{
		Name:        dungeonObjectRec.Name,
		Description: dungeonObjectRec.Description,
		IsEquipped:  dungeonObjectRec.IsEquipped,
		IsStashed:   dungeonObjectRec.IsStashed,
	}, nil
}

func dungeonCharacterToResponseCharacter(
	l logger.Logger,
	characterRec *record.ActionCharacter,
	objectRecs []*record.ActionCharacterObject,
) (*schema.ActionCharacter, error) {

	data := &schema.ActionCharacter{
		Name:                characterRec.Name,
		Strength:            characterRec.Strength,
		Dexterity:           characterRec.Dexterity,
		Intelligence:        characterRec.Intelligence,
		CurrentStrength:     characterRec.CurrentStrength,
		CurrentDexterity:    characterRec.CurrentDexterity,
		CurrentIntelligence: characterRec.CurrentIntelligence,
		Health:              characterRec.Health,
		Fatigue:             characterRec.Fatigue,
		CurrentHealth:       characterRec.CurrentHealth,
		CurrentFatigue:      characterRec.CurrentFatigue,
		StashedObjects:      []schema.ActionObject{},
		EquippedObjects:     []schema.ActionObject{},
	}

	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			l.Debug("Adding character equipped object >%#v<", objectRec)
			data.EquippedObjects = append(data.EquippedObjects, schema.ActionObject{
				Name:       objectRec.Name,
				IsEquipped: objectRec.IsEquipped,
				IsStashed:  objectRec.IsStashed,
			})
			continue
		}
		if objectRec.IsStashed {
			l.Debug("Adding character stashed object >%#v<", objectRec)
			data.StashedObjects = append(data.StashedObjects, schema.ActionObject{
				Name:       objectRec.Name,
				IsEquipped: objectRec.IsEquipped,
				IsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

func dungeonTargetCharacterToResponseTargetCharacter(
	l logger.Logger,
	characterRec *record.ActionCharacter,
	objectRecs []*record.ActionCharacterObject,
) (*schema.ActionCharacter, error) {
	data := &schema.ActionCharacter{
		Name:                characterRec.Name,
		Strength:            characterRec.Strength,
		Dexterity:           characterRec.Dexterity,
		Intelligence:        characterRec.Intelligence,
		CurrentStrength:     characterRec.CurrentStrength,
		CurrentDexterity:    characterRec.CurrentDexterity,
		CurrentIntelligence: characterRec.CurrentIntelligence,
		Health:              characterRec.Health,
		Fatigue:             characterRec.Fatigue,
		CurrentHealth:       characterRec.CurrentHealth,
		CurrentFatigue:      characterRec.CurrentFatigue,
		EquippedObjects:     []schema.ActionObject{},
	}

	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			l.Debug("Adding target character equipped object >%#v<", objectRec)
			data.EquippedObjects = append(data.EquippedObjects, schema.ActionObject{
				Name:       objectRec.Name,
				IsEquipped: objectRec.IsEquipped,
				IsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

func dungeonMonsterToResponseMonster(
	l logger.Logger,
	dungeonMonsterRec *record.ActionMonster,
	objectRecs []*record.ActionMonsterObject,
) (*schema.ActionMonster, error) {

	data := &schema.ActionMonster{
		Name:                dungeonMonsterRec.Name,
		Strength:            dungeonMonsterRec.Strength,
		Dexterity:           dungeonMonsterRec.Dexterity,
		Intelligence:        dungeonMonsterRec.Intelligence,
		CurrentStrength:     dungeonMonsterRec.CurrentStrength,
		CurrentDexterity:    dungeonMonsterRec.CurrentDexterity,
		CurrentIntelligence: dungeonMonsterRec.CurrentIntelligence,
		Health:              dungeonMonsterRec.Health,
		Fatigue:             dungeonMonsterRec.Fatigue,
		CurrentHealth:       dungeonMonsterRec.CurrentHealth,
		CurrentFatigue:      dungeonMonsterRec.CurrentFatigue,
		StashedObjects:      []schema.ActionObject{},
		EquippedObjects:     []schema.ActionObject{},
	}

	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			l.Debug("Adding monster equipped object >%#v<", objectRec)
			data.EquippedObjects = append(data.EquippedObjects, schema.ActionObject{
				Name:       objectRec.Name,
				IsEquipped: objectRec.IsEquipped,
				IsStashed:  objectRec.IsStashed,
			})
			continue
		}
		if objectRec.IsStashed {
			l.Debug("Adding monster stashed object >%#v<", objectRec)
			data.StashedObjects = append(data.StashedObjects, schema.ActionObject{
				Name:       objectRec.Name,
				IsEquipped: objectRec.IsEquipped,
				IsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

func dungeonTargetMonsterToResponseTargetMonster(
	l logger.Logger,
	monsterRec *record.ActionMonster,
	objectRecs []*record.ActionMonsterObject,
) (*schema.ActionMonster, error) {
	data := &schema.ActionMonster{
		Name:                monsterRec.Name,
		Strength:            monsterRec.Strength,
		Dexterity:           monsterRec.Dexterity,
		Intelligence:        monsterRec.Intelligence,
		CurrentStrength:     monsterRec.CurrentStrength,
		CurrentDexterity:    monsterRec.CurrentDexterity,
		CurrentIntelligence: monsterRec.CurrentIntelligence,
		Health:              monsterRec.Health,
		Fatigue:             monsterRec.Fatigue,
		CurrentHealth:       monsterRec.CurrentHealth,
		CurrentFatigue:      monsterRec.CurrentFatigue,
		EquippedObjects:     []schema.ActionObject{},
	}

	for _, objectRec := range objectRecs {
		if objectRec.IsEquipped {
			data.EquippedObjects = append(data.EquippedObjects, schema.ActionObject{
				Name:       objectRec.Name,
				IsEquipped: objectRec.IsEquipped,
				IsStashed:  objectRec.IsStashed,
			})
			continue
		}
	}

	return data, nil
}

// actionLocationToResponseLocation -
func actionLocationToResponseLocation(recordSet *record.ActionLocationRecordSet) (*schema.ActionLocation, error) {

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
					Name: dungeonCharacterRec.Name,
				})
		}
	}

	var monstersData []schema.ActionLocationMonster
	if len(recordSet.ActionMonsterRecs) > 0 {
		for _, dungeonMonsterRec := range recordSet.ActionMonsterRecs {
			monstersData = append(monstersData,
				schema.ActionLocationMonster{
					Name: dungeonMonsterRec.Name,
				})
		}
	}

	var objectsData []schema.ActionLocationObject
	if len(recordSet.ActionObjectRecs) > 0 {
		for _, dungeonObjectRec := range recordSet.ActionObjectRecs {
			objectsData = append(objectsData,
				schema.ActionLocationObject{
					Name: dungeonObjectRec.Name,
				})
		}
	}

	data := &schema.ActionLocation{
		Name:        dungeonLocationRec.Name,
		Description: dungeonLocationRec.Description,
		Directions:  directions,
		Characters:  charactersData,
		Monsters:    monstersData,
		Objects:     objectsData,
	}

	return data, nil
}

// TODO: This might evolve to being its own lifeform so should potentially be
// pulled out of this mapper and into its own package.

// dungeonActionToNarrative -
func dungeonActionToNarrative(l logger.Logger, set record.ActionRecordSet) (string, error) {

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
