package runner

import (
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// actionResponseData -
func actionResponseData(l logger.Logger, rs record.ActionRecordSet) (*schema.ActionResponseData, error) {

	actionRec := rs.ActionRec

	var err error
	var characterData *schema.ActionCharacter
	if rs.ActionCharacterRec != nil {
		characterData, err = actionCharacterResponseData(
			l,
			rs.ActionCharacterRec,
			rs.ActionCharacterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	var monsterData *schema.ActionMonster
	if rs.ActionMonsterRec != nil {
		monsterData, err = actionMonsterResponseData(
			l,
			rs.ActionMonsterRec,
			rs.ActionMonsterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	// Current location
	locationData, err := actionLocationResponseData(l, rs.CurrentLocation)
	if err != nil {
		return nil, err
	}

	// Target location
	var targetActionLocation *schema.ActionLocation
	if rs.TargetLocation != nil {
		targetActionLocation, err = actionLocationResponseData(
			l,
			rs.TargetLocation,
		)
		if err != nil {
			return nil, err
		}
	}

	if actionRec.ResolvedTargetLocationDirection.Valid {
		targetActionLocation.Direction = actionRec.ResolvedTargetLocationDirection.String
	}

	// Equipped object
	var equippedActionLocationObject *schema.ActionObject
	if rs.EquippedActionObjectRec != nil {
		equippedActionLocationObject, err = actionObjectResponseData(
			l,
			rs.EquippedActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Stashed object
	var stashedActionLocationObject *schema.ActionObject
	if rs.StashedActionObjectRec != nil {
		stashedActionLocationObject, err = actionObjectResponseData(
			l,
			rs.StashedActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Dropped object
	var droppedActionLocationObject *schema.ActionObject
	if rs.DroppedActionObjectRec != nil {
		droppedActionLocationObject, err = actionObjectResponseData(
			l,
			rs.DroppedActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Target object
	var targetActionLocationObject *schema.ActionObject
	if rs.TargetActionObjectRec != nil {
		targetActionLocationObject, err = actionObjectResponseData(
			l,
			rs.TargetActionObjectRec,
		)
		if err != nil {
			return nil, err
		}
	}

	// Target character
	var targetActionLocationCharacter *schema.ActionCharacter
	if rs.TargetActionCharacterRec != nil {
		targetActionLocationCharacter, err = actionTargetCharacterResponseData(
			l,
			rs.TargetActionCharacterRec,
			rs.TargetActionCharacterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	// Target monster
	var targetActionLocationMonster *schema.ActionMonster
	if rs.TargetActionMonsterRec != nil {
		targetActionLocationMonster, err = actionTargetMonsterResponseData(
			l,
			rs.TargetActionMonsterRec,
			rs.TargetActionMonsterObjectRecs,
		)
		if err != nil {
			return nil, err
		}
	}

	narrative, err := actionNarrativeResponseData(
		l,
		rs,
	)
	if err != nil {
		return nil, err
	}

	data := schema.ActionResponseData{
		ID:              actionRec.ID,
		Command:         actionRec.ResolvedCommand,
		TurnNumber:      actionRec.TurnNumber,
		SerialNumber:    null.NullInt16ToInt16(actionRec.SerialNumber),
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
		CreatedAt:       actionRec.CreatedAt,
	}

	return &data, nil
}

func actionObjectResponseData(l logger.Logger, dungeonObjectRec *record.ActionObject) (*schema.ActionObject, error) {
	return &schema.ActionObject{
		Name:        dungeonObjectRec.Name,
		Description: dungeonObjectRec.Description,
		IsEquipped:  dungeonObjectRec.IsEquipped,
		IsStashed:   dungeonObjectRec.IsStashed,
	}, nil
}

func actionCharacterResponseData(
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

func actionTargetCharacterResponseData(
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

func actionMonsterResponseData(
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
		EquippedObjects:     []schema.ActionObject{},
	}

	// NOTE: We only ever provide a monsters equipped objects and never provide
	// a monsters stashed objects.
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
	}

	return data, nil
}

func actionTargetMonsterResponseData(
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
		for _, actionCharacterRec := range recordSet.ActionCharacterRecs {
			charactersData = append(charactersData,
				schema.ActionLocationCharacter{
					Name:           actionCharacterRec.Name,
					Health:         actionCharacterRec.Health,
					CurrentHealth:  actionCharacterRec.CurrentHealth,
					Fatigue:        actionCharacterRec.Fatigue,
					CurrentFatigue: actionCharacterRec.CurrentFatigue,
				})
		}
	}

	var monstersData []schema.ActionLocationMonster
	if len(recordSet.ActionMonsterRecs) > 0 {
		for _, actionMonsterRec := range recordSet.ActionMonsterRecs {
			monstersData = append(monstersData,
				schema.ActionLocationMonster{
					Name:           actionMonsterRec.Name,
					Health:         actionMonsterRec.Health,
					CurrentHealth:  actionMonsterRec.CurrentHealth,
					Fatigue:        actionMonsterRec.Fatigue,
					CurrentFatigue: actionMonsterRec.CurrentFatigue,
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

// actionNarrativeResponseData -
func actionNarrativeResponseData(l logger.Logger, set record.ActionRecordSet) (string, error) {

	desc := ""
	if set.ActionCharacterRec != nil {
		desc += set.ActionCharacterRec.Name
	} else if set.ActionMonsterRec != nil {
		desc += set.ActionMonsterRec.Name
	}

	switch set.ActionRec.ResolvedCommand {
	case record.ActionCommandAttack:
		desc += " attacks "
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
