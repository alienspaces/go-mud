package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/core/nullstring"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// TODO: Dry out/refactor ProcessCharacterAction and ProcessMonsterAction

// ProcessCharacterAction - Processes a submitted character action
func (m *Model) ProcessCharacterAction(dungeonInstanceID string, characterInstanceID string, sentence string) (*record.ActionRecordSet, error) {
	l := m.Logger("ProcessCharacterAction")

	l.Info("Processing character ID >%s< action command >%s<", characterInstanceID, sentence)

	// Verify the character performing the action exists within the specified dungeon
	characterInstanceViewRec, err := m.GetCharacterInstanceViewRec(characterInstanceID)
	if err != nil {
		l.Warn("failed getting character record before performing action >%v<", err)
		return nil, err
	}
	if characterInstanceViewRec == nil {
		msg := fmt.Sprintf("failed getting character record ID >%s< before performing action", characterInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if characterInstanceViewRec.DungeonInstanceID != dungeonInstanceID {
		msg := fmt.Sprintf("character ID >%s< does not exist in dungeon ID >%s<", characterInstanceID, dungeonInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Get the current dungeon location set of related records
	locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(characterInstanceViewRec.LocationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}
	if locationInstanceRecordSet == nil {
		msg := fmt.Sprintf("failed getting dungeon location record ID >%s< set before performing action", characterInstanceViewRec.LocationInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Resolve the submitted character action
	args := &ResolverArgs{
		DungeonInstanceID:  characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID: characterInstanceViewRec.LocationInstanceID,
		EntityType:         EntityTypeCharacter,
		EntityInstanceID:   characterInstanceViewRec.ID,
	}
	actionRec, err := m.resolveAction(sentence, args, locationInstanceRecordSet)
	if err != nil {
		l.Warn("failed resolving character action >%v<", err)
		return nil, err
	}

	l.Info("(before) Dungeon action record command >%s<", actionRec.ResolvedCommand)
	l.Info("(before) Dungeon action record LocationInstanceID >%s<", actionRec.LocationInstanceID)
	l.Info("(before) Dungeon action record CharacterInstanceID >%s<", nullstring.ToString(actionRec.CharacterInstanceID))
	l.Info("(before) Dungeon action record MonsterInstanceID >%s<", nullstring.ToString(actionRec.MonsterInstanceID))

	// Perform the submitted character action
	actionRec, err = m.performAction(characterInstanceViewRec, nil, actionRec, locationInstanceRecordSet)
	if err != nil {
		l.Warn("failed performing character action >%v<", err)
		return nil, err
	}

	l.Info("(after) Dungeon action record command >%s<", actionRec.ResolvedCommand)
	l.Info("(after) Dungeon action record LocationInstanceID >%s<", actionRec.LocationInstanceID)
	l.Info("(after) Dungeon action record CharacterInstanceID >%s<", nullstring.ToString(actionRec.CharacterInstanceID))
	l.Info("(after) Dungeon action record MonsterInstanceID >%s<", nullstring.ToString(actionRec.MonsterInstanceID))

	// Create the resulting action event record
	err = m.CreateActionRec(actionRec)
	if err != nil {
		l.Warn("failed creating action record >%v<", err)
		return nil, err
	}

	l.Info("Created action record ID >%s<", actionRec.ID)

	// TODO: Maybe don't need to do this... Get the updated character record
	characterInstanceViewRec, err = m.GetCharacterInstanceViewRec(characterInstanceID)
	if err != nil {
		l.Warn("failed getting character record after performing action >%v<", err)
		return nil, err
	}

	// Create action character record
	actionCharacterRec := record.ActionCharacter{
		RecordType:          record.ActionCharacterRecordTypeSource,
		ActionID:            actionRec.ID,
		LocationInstanceID:  actionRec.LocationInstanceID,
		CharacterInstanceID: characterInstanceViewRec.ID,
		Name:                characterInstanceViewRec.Name,
		Strength:            characterInstanceViewRec.Strength,
		Dexterity:           characterInstanceViewRec.Dexterity,
		Intelligence:        characterInstanceViewRec.Intelligence,
		CurrentStrength:     characterInstanceViewRec.CurrentStrength,
		CurrentDexterity:    characterInstanceViewRec.CurrentDexterity,
		CurrentIntelligence: characterInstanceViewRec.CurrentIntelligence,
		Health:              characterInstanceViewRec.Health,
		Fatigue:             characterInstanceViewRec.Fatigue,
		CurrentHealth:       characterInstanceViewRec.CurrentHealth,
		CurrentFatigue:      characterInstanceViewRec.CurrentFatigue,
	}

	err = m.CreateActionCharacterRec(&actionCharacterRec)
	if err != nil {
		l.Warn("failed creating source action character record >%v<", err)
		return nil, err
	}

	// Create action character object records
	objectInstanceViewRecs, err := m.GetCharacterInstanceObjectInstanceViewRecs(characterInstanceViewRec.ID)
	if err != nil {
		l.Warn("failed getting source character object instance view records >%v<", err)
		return nil, err
	}

	actionCharacterObjectRecs := []*record.ActionCharacterObject{}
	for _, objectInstanceViewRec := range objectInstanceViewRecs {
		l.Info("Adding character action object record >%#v<", objectInstanceViewRec)
		dungeonCharacterObjectRec := record.ActionCharacterObject{
			ActionID:            actionRec.ID,
			CharacterInstanceID: characterInstanceViewRec.ID,
			ObjectInstanceID:    objectInstanceViewRec.ID,
			Name:                objectInstanceViewRec.Name,
			IsEquipped:          objectInstanceViewRec.IsEquipped,
			IsStashed:           objectInstanceViewRec.IsStashed,
		}
		err := m.CreateActionCharacterObjectRec(&dungeonCharacterObjectRec)
		if err != nil {
			l.Warn("failed creating source action character object record >%v<", err)
			return nil, err
		}
		actionCharacterObjectRecs = append(actionCharacterObjectRecs, &dungeonCharacterObjectRec)
	}

	actionRecordSet := record.ActionRecordSet{
		ActionRec:                 actionRec,
		ActionCharacterRec:        &actionCharacterRec,
		ActionCharacterObjectRecs: actionCharacterObjectRecs,
	}

	// Get the updated current location instance record set
	locationInstanceRecordSet, err = m.GetLocationInstanceViewRecordSet(actionRec.LocationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set after performing action >%v<", err)
		return nil, err
	}

	l.Info("Dungeon location record set location name >%s<", locationInstanceRecordSet.LocationInstanceViewRec.Name)
	l.Info("Dungeon location record set characters >%d<", len(locationInstanceRecordSet.CharacterInstanceViewRecs))
	l.Info("Dungeon location record set monsters >%d<", len(locationInstanceRecordSet.MonsterInstanceViewRecs))
	l.Info("Dungeon location record set objects >%d<", len(locationInstanceRecordSet.ObjectInstanceViewRecs))

	// Current location
	locationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.LocationInstanceID)
	if err != nil {
		l.Warn("failed getting dungeon location record after performing action >%v<", err)
		return nil, err
	}

	// TODO: What is this for, can we use a location instance view record?
	currentLocationRecordSet := record.ActionLocationRecordSet{
		LocationInstanceViewRec: locationInstanceViewRec,
		ActionCharacterRecs:     []*record.ActionCharacter{},
		ActionMonsterRecs:       []*record.ActionMonster{},
		ActionObjectRecs:        []*record.ActionObject{},
	}

	// Character Occupants: Create the action character record for each character now at the current location
	if len(locationInstanceRecordSet.CharacterInstanceViewRecs) > 0 {
		for _, characterInstanceViewRec := range locationInstanceRecordSet.CharacterInstanceViewRecs {

			actionCharacterRec := record.ActionCharacter{
				RecordType:          record.ActionCharacterRecordTypeOccupant,
				ActionID:            actionRec.ID,
				LocationInstanceID:  locationInstanceViewRec.ID,
				CharacterInstanceID: characterInstanceViewRec.ID,
				Name:                characterInstanceViewRec.Name,
				Strength:            characterInstanceViewRec.Strength,
				Dexterity:           characterInstanceViewRec.Dexterity,
				Intelligence:        characterInstanceViewRec.Intelligence,
				CurrentStrength:     characterInstanceViewRec.CurrentStrength,
				CurrentDexterity:    characterInstanceViewRec.CurrentDexterity,
				CurrentIntelligence: characterInstanceViewRec.CurrentIntelligence,
				Health:              characterInstanceViewRec.Health,
				Fatigue:             characterInstanceViewRec.Fatigue,
				CurrentHealth:       characterInstanceViewRec.CurrentHealth,
				CurrentFatigue:      characterInstanceViewRec.CurrentFatigue,
			}

			err := m.CreateActionCharacterRec(&actionCharacterRec)
			if err != nil {
				l.Warn("failed creating current location action character record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action character record ID >%s<", actionCharacterRec.ID)
			currentLocationRecordSet.ActionCharacterRecs = append(currentLocationRecordSet.ActionCharacterRecs, &actionCharacterRec)
		}
	}

	// Monster Occupants: Create the action monster record for each monster now at the current location
	if len(locationInstanceRecordSet.MonsterInstanceViewRecs) > 0 {
		for _, monsterInstanceViewRec := range locationInstanceRecordSet.MonsterInstanceViewRecs {
			actionMonsterRec := record.ActionMonster{
				RecordType:          record.ActionMonsterRecordTypeOccupant,
				ActionID:            actionRec.ID,
				LocationInstanceID:  locationInstanceViewRec.ID,
				MonsterInstanceID:   monsterInstanceViewRec.ID,
				Name:                monsterInstanceViewRec.Name,
				Strength:            monsterInstanceViewRec.Strength,
				Dexterity:           monsterInstanceViewRec.Dexterity,
				Intelligence:        monsterInstanceViewRec.Intelligence,
				CurrentStrength:     monsterInstanceViewRec.CurrentStrength,
				CurrentDexterity:    monsterInstanceViewRec.CurrentDexterity,
				CurrentIntelligence: monsterInstanceViewRec.CurrentIntelligence,
				Health:              monsterInstanceViewRec.Health,
				Fatigue:             monsterInstanceViewRec.Fatigue,
				CurrentHealth:       monsterInstanceViewRec.CurrentHealth,
				CurrentFatigue:      monsterInstanceViewRec.CurrentFatigue,
			}
			err := m.CreateActionMonsterRec(&actionMonsterRec)
			if err != nil {
				l.Warn("failed creating current location action monster record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action monster record ID >%s<", actionMonsterRec.ID)
			currentLocationRecordSet.ActionMonsterRecs = append(currentLocationRecordSet.ActionMonsterRecs, &actionMonsterRec)
		}
	}

	// Object Occupants: Create the action object record for each object now at the current location
	if len(locationInstanceRecordSet.ObjectInstanceViewRecs) > 0 {
		for _, objectInstanceViewRec := range locationInstanceRecordSet.ObjectInstanceViewRecs {
			dungeonActionObjectRec := record.ActionObject{
				RecordType:         record.ActionObjectRecordTypeOccupant,
				ActionID:           actionRec.ID,
				LocationInstanceID: locationInstanceViewRec.ID,
				ObjectInstanceID:   objectInstanceViewRec.ID,
				Name:               objectInstanceViewRec.Name,
				Description:        objectInstanceViewRec.Description,
				IsStashed:          objectInstanceViewRec.IsStashed,
				IsEquipped:         objectInstanceViewRec.IsEquipped,
			}
			err := m.CreateActionObjectRec(&dungeonActionObjectRec)
			if err != nil {
				l.Warn("failed creating current location action object record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action object record ID >%s<", dungeonActionObjectRec.ID)
			currentLocationRecordSet.ActionObjectRecs = append(currentLocationRecordSet.ActionObjectRecs, &dungeonActionObjectRec)
		}
	}

	actionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target location instance record set when set
	if actionRec.ResolvedTargetLocationInstanceID.Valid {

		locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(actionRec.ResolvedTargetLocationInstanceID.String, true)
		if err != nil {
			l.Warn("failed getting target location instance record set after performing action >%v<", err)
			return nil, err
		}

		// Target location
		targetLocationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.ResolvedTargetLocationInstanceID.String)
		if err != nil {
			l.Warn("failed getting target location instance view record after performing action >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := record.ActionLocationRecordSet{
			LocationInstanceViewRec: targetLocationInstanceViewRec,
			ActionCharacterRecs:     []*record.ActionCharacter{},
			ActionMonsterRecs:       []*record.ActionMonster{},
			ActionObjectRecs:        []*record.ActionObject{},
		}

		// Character Occupants: Create the action character record for each character at the target location
		if len(locationInstanceRecordSet.CharacterInstanceViewRecs) > 0 {
			for _, characterInstanceViewRec := range locationInstanceRecordSet.CharacterInstanceViewRecs {
				actionCharacterRec := record.ActionCharacter{
					RecordType:          record.ActionCharacterRecordTypeOccupant,
					ActionID:            actionRec.ID,
					LocationInstanceID:  targetLocationInstanceViewRec.ID,
					CharacterInstanceID: characterInstanceViewRec.ID,
					Name:                characterInstanceViewRec.Name,
					Strength:            characterInstanceViewRec.Strength,
					Dexterity:           characterInstanceViewRec.Dexterity,
					Intelligence:        characterInstanceViewRec.Intelligence,
					CurrentStrength:     characterInstanceViewRec.CurrentStrength,
					CurrentDexterity:    characterInstanceViewRec.CurrentDexterity,
					CurrentIntelligence: characterInstanceViewRec.CurrentIntelligence,
					Health:              characterInstanceViewRec.Health,
					Fatigue:             characterInstanceViewRec.Fatigue,
					CurrentHealth:       characterInstanceViewRec.CurrentHealth,
					CurrentFatigue:      characterInstanceViewRec.CurrentFatigue,
				}

				err := m.CreateActionCharacterRec(&actionCharacterRec)
				if err != nil {
					l.Warn("failed creating target location occupant action character record >%v<", err)
					return nil, err
				}

				l.Info("Created target location occupant action character record ID >%s<", actionCharacterRec.ID)
				targetLocationRecordSet.ActionCharacterRecs = append(targetLocationRecordSet.ActionCharacterRecs, &actionCharacterRec)
			}
		}

		// Monster Occupants: Create the action monster record for each monster at the target location
		if len(locationInstanceRecordSet.MonsterInstanceViewRecs) > 0 {
			for _, monsterInstanceViewRec := range locationInstanceRecordSet.MonsterInstanceViewRecs {
				actionMonsterRec := record.ActionMonster{
					RecordType:          record.ActionMonsterRecordTypeOccupant,
					ActionID:            actionRec.ID,
					LocationInstanceID:  targetLocationInstanceViewRec.ID,
					MonsterInstanceID:   monsterInstanceViewRec.ID,
					Name:                monsterInstanceViewRec.Name,
					Strength:            monsterInstanceViewRec.Strength,
					Dexterity:           monsterInstanceViewRec.Dexterity,
					Intelligence:        monsterInstanceViewRec.Intelligence,
					CurrentStrength:     monsterInstanceViewRec.CurrentStrength,
					CurrentDexterity:    monsterInstanceViewRec.CurrentDexterity,
					CurrentIntelligence: monsterInstanceViewRec.CurrentIntelligence,
					Health:              monsterInstanceViewRec.Health,
					Fatigue:             monsterInstanceViewRec.Fatigue,
					CurrentHealth:       monsterInstanceViewRec.CurrentHealth,
					CurrentFatigue:      monsterInstanceViewRec.CurrentFatigue,
				}
				err := m.CreateActionMonsterRec(&actionMonsterRec)
				if err != nil {
					l.Warn("failed creating target location occupant action monster record >%v<", err)
					return nil, err
				}

				l.Info("Created target location occupant action monster record ID >%s<", actionMonsterRec.ID)
				targetLocationRecordSet.ActionMonsterRecs = append(targetLocationRecordSet.ActionMonsterRecs, &actionMonsterRec)
			}
		}

		// Object Occupants: Create the action object record for each object at the target location
		if len(locationInstanceRecordSet.ObjectInstanceViewRecs) > 0 {
			for _, objectInstanceViewRec := range locationInstanceRecordSet.ObjectInstanceViewRecs {
				dungeonActionObjectRec := record.ActionObject{
					RecordType:         record.ActionObjectRecordTypeOccupant,
					ActionID:           actionRec.ID,
					LocationInstanceID: locationInstanceViewRec.ID,
					ObjectInstanceID:   objectInstanceViewRec.ID,
					Name:               objectInstanceViewRec.Name,
					Description:        objectInstanceViewRec.Description,
					IsStashed:          objectInstanceViewRec.IsStashed,
					IsEquipped:         objectInstanceViewRec.IsEquipped,
				}
				err := m.CreateActionObjectRec(&dungeonActionObjectRec)
				if err != nil {
					l.Warn("failed creating target location occupant action object record >%v<", err)
					return nil, err
				}

				l.Info("Created target location occupant action object record ID >%s<", dungeonActionObjectRec.ID)
				targetLocationRecordSet.ActionObjectRecs = append(targetLocationRecordSet.ActionObjectRecs, &dungeonActionObjectRec)
			}
		}

		actionRecordSet.TargetLocation = &targetLocationRecordSet
	}

	// Create the target character action record
	if actionRec.ResolvedTargetCharacterInstanceID.Valid {

		l.Info("Resolved target character instance ID >%s<", actionRec.ResolvedTargetCharacterInstanceID.String)

		targetCharacterInstanceViewRec, err := m.GetCharacterInstanceViewRec(actionRec.ResolvedTargetCharacterInstanceID.String)
		if err != nil {
			l.Warn("failed getting target character instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionCharacter{
			RecordType:          record.ActionCharacterRecordTypeTarget,
			ActionID:            actionRec.ID,
			LocationInstanceID:  locationInstanceViewRec.ID,
			CharacterInstanceID: targetCharacterInstanceViewRec.ID,
			Name:                targetCharacterInstanceViewRec.Name,
			Strength:            targetCharacterInstanceViewRec.Strength,
			Dexterity:           targetCharacterInstanceViewRec.Dexterity,
			Intelligence:        targetCharacterInstanceViewRec.Intelligence,
			CurrentStrength:     targetCharacterInstanceViewRec.CurrentStrength,
			CurrentDexterity:    targetCharacterInstanceViewRec.CurrentDexterity,
			CurrentIntelligence: targetCharacterInstanceViewRec.CurrentIntelligence,
			Health:              targetCharacterInstanceViewRec.Health,
			Fatigue:             targetCharacterInstanceViewRec.Fatigue,
			CurrentHealth:       targetCharacterInstanceViewRec.CurrentHealth,
			CurrentFatigue:      targetCharacterInstanceViewRec.CurrentFatigue,
		}

		err = m.CreateActionCharacterRec(rec)
		if err != nil {
			l.Warn("failed creating target action character record >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionCharacterRec = rec

		// Create action character object records
		objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(targetCharacterInstanceViewRec.ID)
		if err != nil {
			l.Warn("failed getting target character object records >%v<", err)
			return nil, err
		}

		l.Info("Adding >%d< target character object records", len(objectInstanceViewRecs))

		targetCharacterObjectRecs := []*record.ActionCharacterObject{}
		for _, objectInstanceViewRec := range objectInstanceViewRecs {
			l.Info("Adding target character object record >%v<", objectInstanceViewRecs)
			dungeonCharacterObjectRec := record.ActionCharacterObject{
				ActionID:            actionRec.ID,
				CharacterInstanceID: targetCharacterInstanceViewRec.ID,
				ObjectInstanceID:    objectInstanceViewRec.ID,
				Name:                objectInstanceViewRec.Name,
				IsEquipped:          objectInstanceViewRec.IsEquipped,
				IsStashed:           objectInstanceViewRec.IsStashed,
			}
			err := m.CreateActionCharacterObjectRec(&dungeonCharacterObjectRec)
			if err != nil {
				l.Warn("failed creating source action character object record >%v<", err)
				return nil, err
			}
			targetCharacterObjectRecs = append(targetCharacterObjectRecs, &dungeonCharacterObjectRec)
		}
		actionRecordSet.TargetActionCharacterObjectRecs = targetCharacterObjectRecs
	}

	// Create the target dungeon monster action record
	if actionRec.ResolvedTargetMonsterInstanceID.Valid {

		l.Info("Resolved target monster ID >%s<", actionRec.ResolvedTargetMonsterInstanceID.String)

		targetMonsterInstanceViewRec, err := m.GetMonsterInstanceViewRec(actionRec.ResolvedTargetMonsterInstanceID.String)
		if err != nil {
			l.Warn("failed getting target monster instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionMonster{
			RecordType:          record.ActionMonsterRecordTypeTarget,
			ActionID:            actionRec.ID,
			LocationInstanceID:  locationInstanceViewRec.ID,
			MonsterInstanceID:   targetMonsterInstanceViewRec.ID,
			Name:                targetMonsterInstanceViewRec.Name,
			Strength:            targetMonsterInstanceViewRec.Strength,
			Dexterity:           targetMonsterInstanceViewRec.Dexterity,
			Intelligence:        targetMonsterInstanceViewRec.Intelligence,
			CurrentStrength:     targetMonsterInstanceViewRec.CurrentStrength,
			CurrentDexterity:    targetMonsterInstanceViewRec.CurrentDexterity,
			CurrentIntelligence: targetMonsterInstanceViewRec.CurrentIntelligence,
			Health:              targetMonsterInstanceViewRec.Health,
			Fatigue:             targetMonsterInstanceViewRec.Fatigue,
			CurrentHealth:       targetMonsterInstanceViewRec.CurrentHealth,
			CurrentFatigue:      targetMonsterInstanceViewRec.CurrentFatigue,
		}

		err = m.CreateActionMonsterRec(rec)
		if err != nil {
			l.Warn("failed creating target action monster record >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionMonsterRec = rec

		// Create action monster object records
		objectInstanceViewRecs, err := m.GetMonsterInstanceEquippedObjectInstanceViewRecs(targetMonsterInstanceViewRec.ID)
		if err != nil {
			l.Warn("failed getting target monster object records >%v<", err)
			return nil, err
		}

		l.Info("Adding >%d< target monster object records", len(objectInstanceViewRecs))

		targetMonsterObjectRecs := []*record.ActionMonsterObject{}
		for _, objectInstanceViewRec := range objectInstanceViewRecs {
			l.Info("Adding target monster object record >%v<", objectInstanceViewRec)
			dungeonMonsterObjectRec := record.ActionMonsterObject{
				ActionID:          actionRec.ID,
				MonsterInstanceID: targetMonsterInstanceViewRec.ID,
				ObjectInstanceID:  objectInstanceViewRec.ID,
				Name:              objectInstanceViewRec.Name,
				IsEquipped:        objectInstanceViewRec.IsEquipped,
				IsStashed:         objectInstanceViewRec.IsStashed,
			}
			err := m.CreateActionMonsterObjectRec(&dungeonMonsterObjectRec)
			if err != nil {
				l.Warn("failed creating source action monster object record >%v<", err)
				return nil, err
			}
			targetMonsterObjectRecs = append(targetMonsterObjectRecs, &dungeonMonsterObjectRec)
		}
		actionRecordSet.TargetActionMonsterObjectRecs = targetMonsterObjectRecs
	}

	// Create the target dungeon object action record
	if actionRec.ResolvedTargetObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedTargetObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting target object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeTarget,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating target action object record >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionObjectRec = rec
	}

	// Create the stashed dungeon object action record
	if actionRec.ResolvedStashedObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedStashedObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting stashed object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeStashed,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating stashed action object record >%v<", err)
			return nil, err
		}
		actionRecordSet.StashedActionObjectRec = rec
	}

	// Create the equipped dungeon object action record
	if actionRec.ResolvedEquippedObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedEquippedObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting equipped object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeEquipped,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating equipped action object record >%v<", err)
			return nil, err
		}
		actionRecordSet.EquippedActionObjectRec = rec
	}

	// Create the dropped dungeon object action record
	if actionRec.ResolvedDroppedObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedDroppedObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting dropped object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeDropped,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating dropped action object record >%v<", err)
			return nil, err
		}

		l.Warn("Assigning dropped object >%v< to action record set", rec)

		actionRecordSet.DroppedActionObjectRec = rec
	}

	return &actionRecordSet, nil
}

// ProcessMonsterAction - Processes a submitted character action
func (m *Model) ProcessMonsterAction(dungeonInstanceID string, monsterInstanceID string, sentence string) (*record.ActionRecordSet, error) {
	l := m.Logger("ProcessMonsterAction")

	l.Info("Processing monster ID >%s< action command >%s<", monsterInstanceID, sentence)

	// Verify the monster performing the action exists within the specified dungeon
	monsterInstanceViewRec, err := m.GetMonsterInstanceViewRec(monsterInstanceID)
	if err != nil {
		l.Warn("failed getting monster record before performing action >%v<", err)
		return nil, err
	}
	if monsterInstanceViewRec == nil {
		msg := fmt.Sprintf("failed getting monster record ID >%s< before performing action", monsterInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if monsterInstanceViewRec.DungeonInstanceID != dungeonInstanceID {
		msg := fmt.Sprintf("monster ID >%s< does not exist in dungeon ID >%s<", monsterInstanceID, dungeonInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Get the current dungeon location set of related records
	locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(monsterInstanceViewRec.LocationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}
	if locationInstanceRecordSet == nil {
		msg := fmt.Sprintf("failed getting dungeon location record ID >%s< set before performing action", monsterInstanceViewRec.LocationInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Resolve the submitted monster action
	args := &ResolverArgs{
		DungeonInstanceID:  monsterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID: monsterInstanceViewRec.LocationInstanceID,
		EntityType:         EntityTypeMonster,
		EntityInstanceID:   monsterInstanceViewRec.ID,
	}
	actionRec, err := m.resolveAction(sentence, args, locationInstanceRecordSet)
	if err != nil {
		l.Warn("failed resolving monster action >%v<", err)
		return nil, err
	}

	l.Info("(before) Dungeon action record command >%s<", actionRec.ResolvedCommand)
	l.Info("(before) Dungeon action record LocationInstanceID >%s<", actionRec.LocationInstanceID)
	l.Info("(before) Dungeon action record MonsterInstanceID >%s<", nullstring.ToString(actionRec.MonsterInstanceID))
	l.Info("(before) Dungeon action record MonsterInstanceID >%s<", nullstring.ToString(actionRec.MonsterInstanceID))

	// Perform the submitted monster action
	actionRec, err = m.performAction(nil, monsterInstanceViewRec, actionRec, locationInstanceRecordSet)
	if err != nil {
		l.Warn("failed performing monster action >%v<", err)
		return nil, err
	}

	// BWW

	l.Info("(after) Dungeon action record command >%s<", actionRec.ResolvedCommand)
	l.Info("(after) Dungeon action record LocationInstanceID >%s<", actionRec.LocationInstanceID)
	l.Info("(after) Dungeon action record MonsterInstanceID >%s<", nullstring.ToString(actionRec.MonsterInstanceID))
	l.Info("(after) Dungeon action record CharacterInstanceID >%s<", nullstring.ToString(actionRec.CharacterInstanceID))

	// Create the resulting action event record
	err = m.CreateActionRec(actionRec)
	if err != nil {
		l.Warn("failed creating action record >%v<", err)
		return nil, err
	}

	l.Info("Created action record ID >%s<", actionRec.ID)

	// TODO: Maybe don't need to do this... Get the updated monster record
	monsterInstanceViewRec, err = m.GetMonsterInstanceViewRec(monsterInstanceID)
	if err != nil {
		l.Warn("failed getting monster record after performing action >%v<", err)
		return nil, err
	}

	// Create action monster record
	actionMonsterRec := record.ActionMonster{
		RecordType:          record.ActionMonsterRecordTypeSource,
		ActionID:            actionRec.ID,
		LocationInstanceID:  actionRec.LocationInstanceID,
		MonsterInstanceID:   monsterInstanceViewRec.ID,
		Name:                monsterInstanceViewRec.Name,
		Strength:            monsterInstanceViewRec.Strength,
		Dexterity:           monsterInstanceViewRec.Dexterity,
		Intelligence:        monsterInstanceViewRec.Intelligence,
		CurrentStrength:     monsterInstanceViewRec.CurrentStrength,
		CurrentDexterity:    monsterInstanceViewRec.CurrentDexterity,
		CurrentIntelligence: monsterInstanceViewRec.CurrentIntelligence,
		Health:              monsterInstanceViewRec.Health,
		Fatigue:             monsterInstanceViewRec.Fatigue,
		CurrentHealth:       monsterInstanceViewRec.CurrentHealth,
		CurrentFatigue:      monsterInstanceViewRec.CurrentFatigue,
	}

	err = m.CreateActionMonsterRec(&actionMonsterRec)
	if err != nil {
		l.Warn("failed creating source action monster record >%v<", err)
		return nil, err
	}

	// Create action monster object records
	objectInstanceViewRecs, err := m.GetMonsterInstanceObjectInstanceViewRecs(monsterInstanceViewRec.ID)
	if err != nil {
		l.Warn("failed getting source monster object instance view records >%v<", err)
		return nil, err
	}

	actionMonsterObjectRecs := []*record.ActionMonsterObject{}
	for _, objectInstanceViewRec := range objectInstanceViewRecs {
		l.Info("Adding monster action object record >%#v<", objectInstanceViewRec)
		dungeonMonsterObjectRec := record.ActionMonsterObject{
			ActionID:          actionRec.ID,
			MonsterInstanceID: monsterInstanceViewRec.ID,
			ObjectInstanceID:  objectInstanceViewRec.ID,
			Name:              objectInstanceViewRec.Name,
			IsEquipped:        objectInstanceViewRec.IsEquipped,
			IsStashed:         objectInstanceViewRec.IsStashed,
		}
		err := m.CreateActionMonsterObjectRec(&dungeonMonsterObjectRec)
		if err != nil {
			l.Warn("failed creating source action character object record >%v<", err)
			return nil, err
		}
		actionMonsterObjectRecs = append(actionMonsterObjectRecs, &dungeonMonsterObjectRec)
	}

	actionRecordSet := record.ActionRecordSet{
		ActionRec:               actionRec,
		ActionMonsterRec:        &actionMonsterRec,
		ActionMonsterObjectRecs: actionMonsterObjectRecs,
	}

	// Get the updated current location instance record set
	locationInstanceRecordSet, err = m.GetLocationInstanceViewRecordSet(actionRec.LocationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set after performing action >%v<", err)
		return nil, err
	}

	l.Info("Dungeon location record set location name >%s<", locationInstanceRecordSet.LocationInstanceViewRec.Name)
	l.Info("Dungeon location record set characters >%d<", len(locationInstanceRecordSet.CharacterInstanceViewRecs))
	l.Info("Dungeon location record set monsters >%d<", len(locationInstanceRecordSet.MonsterInstanceViewRecs))
	l.Info("Dungeon location record set objects >%d<", len(locationInstanceRecordSet.ObjectInstanceViewRecs))

	// Current location
	locationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.LocationInstanceID)
	if err != nil {
		l.Warn("failed getting dungeon location record after performing action >%v<", err)
		return nil, err
	}

	// TODO: What is this for, can we use a location instance view record?
	currentLocationRecordSet := record.ActionLocationRecordSet{
		LocationInstanceViewRec: locationInstanceViewRec,
		ActionCharacterRecs:     []*record.ActionCharacter{},
		ActionMonsterRecs:       []*record.ActionMonster{},
		ActionObjectRecs:        []*record.ActionObject{},
	}

	// Character Occupants: Create the action character record for each character now at the current location
	if len(locationInstanceRecordSet.CharacterInstanceViewRecs) > 0 {
		for _, characterInstanceViewRec := range locationInstanceRecordSet.CharacterInstanceViewRecs {

			actionCharacterRec := record.ActionCharacter{
				RecordType:          record.ActionCharacterRecordTypeOccupant,
				ActionID:            actionRec.ID,
				LocationInstanceID:  locationInstanceViewRec.ID,
				CharacterInstanceID: characterInstanceViewRec.ID,
				Name:                characterInstanceViewRec.Name,
				Strength:            characterInstanceViewRec.Strength,
				Dexterity:           characterInstanceViewRec.Dexterity,
				Intelligence:        characterInstanceViewRec.Intelligence,
				CurrentStrength:     characterInstanceViewRec.CurrentStrength,
				CurrentDexterity:    characterInstanceViewRec.CurrentDexterity,
				CurrentIntelligence: characterInstanceViewRec.CurrentIntelligence,
				Health:              characterInstanceViewRec.Health,
				Fatigue:             characterInstanceViewRec.Fatigue,
				CurrentHealth:       characterInstanceViewRec.CurrentHealth,
				CurrentFatigue:      characterInstanceViewRec.CurrentFatigue,
			}

			err := m.CreateActionCharacterRec(&actionCharacterRec)
			if err != nil {
				l.Warn("failed creating current location action character record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action character record ID >%s<", actionCharacterRec.ID)
			currentLocationRecordSet.ActionCharacterRecs = append(currentLocationRecordSet.ActionCharacterRecs, &actionCharacterRec)
		}
	}

	// Monster Occupants: Create the action monster record for each monster now at the current location
	if len(locationInstanceRecordSet.MonsterInstanceViewRecs) > 0 {
		for _, monsterInstanceViewRec := range locationInstanceRecordSet.MonsterInstanceViewRecs {

			actionMonsterRec := record.ActionMonster{
				RecordType:          record.ActionMonsterRecordTypeOccupant,
				ActionID:            actionRec.ID,
				LocationInstanceID:  locationInstanceViewRec.ID,
				MonsterInstanceID:   monsterInstanceViewRec.ID,
				Name:                monsterInstanceViewRec.Name,
				Strength:            monsterInstanceViewRec.Strength,
				Dexterity:           monsterInstanceViewRec.Dexterity,
				Intelligence:        monsterInstanceViewRec.Intelligence,
				CurrentStrength:     monsterInstanceViewRec.CurrentStrength,
				CurrentDexterity:    monsterInstanceViewRec.CurrentDexterity,
				CurrentIntelligence: monsterInstanceViewRec.CurrentIntelligence,
				Health:              monsterInstanceViewRec.Health,
				Fatigue:             monsterInstanceViewRec.Fatigue,
				CurrentHealth:       monsterInstanceViewRec.CurrentHealth,
				CurrentFatigue:      monsterInstanceViewRec.CurrentFatigue,
			}
			err := m.CreateActionMonsterRec(&actionMonsterRec)
			if err != nil {
				l.Warn("failed creating current location action monster record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action monster record ID >%s<", actionMonsterRec.ID)
			currentLocationRecordSet.ActionMonsterRecs = append(currentLocationRecordSet.ActionMonsterRecs, &actionMonsterRec)
		}
	}

	// Object Occupants: Create the action object record for each object now at the current location
	if len(locationInstanceRecordSet.ObjectInstanceViewRecs) > 0 {
		for _, objectInstanceViewRec := range locationInstanceRecordSet.ObjectInstanceViewRecs {
			dungeonActionObjectRec := record.ActionObject{
				RecordType:         record.ActionObjectRecordTypeOccupant,
				ActionID:           actionRec.ID,
				LocationInstanceID: locationInstanceViewRec.ID,
				ObjectInstanceID:   objectInstanceViewRec.ID,
				Name:               objectInstanceViewRec.Name,
				Description:        objectInstanceViewRec.Description,
				IsStashed:          objectInstanceViewRec.IsStashed,
				IsEquipped:         objectInstanceViewRec.IsEquipped,
			}
			err := m.CreateActionObjectRec(&dungeonActionObjectRec)
			if err != nil {
				l.Warn("failed creating current location action object record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action object record ID >%s<", dungeonActionObjectRec.ID)
			currentLocationRecordSet.ActionObjectRecs = append(currentLocationRecordSet.ActionObjectRecs, &dungeonActionObjectRec)
		}
	}

	actionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target location instance record set when set
	if actionRec.ResolvedTargetLocationInstanceID.Valid {

		locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(actionRec.ResolvedTargetLocationInstanceID.String, true)
		if err != nil {
			l.Warn("failed getting target location instance record set after performing action >%v<", err)
			return nil, err
		}

		// Target location
		targetLocationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.ResolvedTargetLocationInstanceID.String)
		if err != nil {
			l.Warn("failed getting target location instance view record after performing action >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := record.ActionLocationRecordSet{
			LocationInstanceViewRec: targetLocationInstanceViewRec,
			ActionCharacterRecs:     []*record.ActionCharacter{},
			ActionMonsterRecs:       []*record.ActionMonster{},
			ActionObjectRecs:        []*record.ActionObject{},
		}

		// Character Occupants: Create the action character record for each character at the target location
		if len(locationInstanceRecordSet.CharacterInstanceViewRecs) > 0 {
			for _, characterInstanceViewRec := range locationInstanceRecordSet.CharacterInstanceViewRecs {
				actionCharacterRec := record.ActionCharacter{
					RecordType:          record.ActionCharacterRecordTypeOccupant,
					ActionID:            actionRec.ID,
					LocationInstanceID:  targetLocationInstanceViewRec.ID,
					CharacterInstanceID: characterInstanceViewRec.ID,
					Name:                characterInstanceViewRec.Name,
					Strength:            characterInstanceViewRec.Strength,
					Dexterity:           characterInstanceViewRec.Dexterity,
					Intelligence:        characterInstanceViewRec.Intelligence,
					CurrentStrength:     characterInstanceViewRec.CurrentStrength,
					CurrentDexterity:    characterInstanceViewRec.CurrentDexterity,
					CurrentIntelligence: characterInstanceViewRec.CurrentIntelligence,
					Health:              characterInstanceViewRec.Health,
					Fatigue:             characterInstanceViewRec.Fatigue,
					CurrentHealth:       characterInstanceViewRec.CurrentHealth,
					CurrentFatigue:      characterInstanceViewRec.CurrentFatigue,
				}

				err := m.CreateActionCharacterRec(&actionCharacterRec)
				if err != nil {
					l.Warn("failed creating target location occupant action character record >%v<", err)
					return nil, err
				}

				l.Info("Created target location occupant action character record ID >%s<", actionCharacterRec.ID)
				targetLocationRecordSet.ActionCharacterRecs = append(targetLocationRecordSet.ActionCharacterRecs, &actionCharacterRec)
			}
		}

		// Monster Occupants: Create the action monster record for each monster at the target location
		if len(locationInstanceRecordSet.MonsterInstanceViewRecs) > 0 {
			for _, monsterInstanceViewRec := range locationInstanceRecordSet.MonsterInstanceViewRecs {
				actionMonsterRec := record.ActionMonster{
					RecordType:          record.ActionMonsterRecordTypeOccupant,
					ActionID:            actionRec.ID,
					LocationInstanceID:  targetLocationInstanceViewRec.ID,
					MonsterInstanceID:   monsterInstanceViewRec.ID,
					Name:                monsterInstanceViewRec.Name,
					Strength:            monsterInstanceViewRec.Strength,
					Dexterity:           monsterInstanceViewRec.Dexterity,
					Intelligence:        monsterInstanceViewRec.Intelligence,
					CurrentStrength:     monsterInstanceViewRec.CurrentStrength,
					CurrentDexterity:    monsterInstanceViewRec.CurrentDexterity,
					CurrentIntelligence: monsterInstanceViewRec.CurrentIntelligence,
					Health:              monsterInstanceViewRec.Health,
					Fatigue:             monsterInstanceViewRec.Fatigue,
					CurrentHealth:       monsterInstanceViewRec.CurrentHealth,
					CurrentFatigue:      monsterInstanceViewRec.CurrentFatigue,
				}
				err := m.CreateActionMonsterRec(&actionMonsterRec)
				if err != nil {
					l.Warn("failed creating target location occupant action monster record >%v<", err)
					return nil, err
				}

				l.Info("Created target location occupant action monster record ID >%s<", actionMonsterRec.ID)
				targetLocationRecordSet.ActionMonsterRecs = append(targetLocationRecordSet.ActionMonsterRecs, &actionMonsterRec)
			}
		}

		// Object Occupants: Create the action object record for each object at the target location
		if len(locationInstanceRecordSet.ObjectInstanceViewRecs) > 0 {
			for _, objectInstanceViewRec := range locationInstanceRecordSet.ObjectInstanceViewRecs {
				dungeonActionObjectRec := record.ActionObject{
					RecordType:         record.ActionObjectRecordTypeOccupant,
					ActionID:           actionRec.ID,
					LocationInstanceID: locationInstanceViewRec.ID,
					ObjectInstanceID:   objectInstanceViewRec.ID,
					Name:               objectInstanceViewRec.Name,
					Description:        objectInstanceViewRec.Description,
					IsStashed:          objectInstanceViewRec.IsStashed,
					IsEquipped:         objectInstanceViewRec.IsEquipped,
				}
				err := m.CreateActionObjectRec(&dungeonActionObjectRec)
				if err != nil {
					l.Warn("failed creating target location occupant action object record >%v<", err)
					return nil, err
				}

				l.Info("Created target location occupant action object record ID >%s<", dungeonActionObjectRec.ID)
				targetLocationRecordSet.ActionObjectRecs = append(targetLocationRecordSet.ActionObjectRecs, &dungeonActionObjectRec)
			}
		}

		actionRecordSet.TargetLocation = &targetLocationRecordSet
	}

	// Create the target character action record
	if actionRec.ResolvedTargetCharacterInstanceID.Valid {

		l.Info("Resolved target character instance ID >%s<", actionRec.ResolvedTargetCharacterInstanceID.String)

		targetCharacterInstanceViewRec, err := m.GetCharacterInstanceViewRec(actionRec.ResolvedTargetCharacterInstanceID.String)
		if err != nil {
			l.Warn("failed getting target character instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionCharacter{
			RecordType:          record.ActionCharacterRecordTypeTarget,
			ActionID:            actionRec.ID,
			LocationInstanceID:  locationInstanceViewRec.ID,
			CharacterInstanceID: targetCharacterInstanceViewRec.ID,
			Name:                targetCharacterInstanceViewRec.Name,
			Strength:            targetCharacterInstanceViewRec.Strength,
			Dexterity:           targetCharacterInstanceViewRec.Dexterity,
			Intelligence:        targetCharacterInstanceViewRec.Intelligence,
			CurrentStrength:     targetCharacterInstanceViewRec.CurrentStrength,
			CurrentDexterity:    targetCharacterInstanceViewRec.CurrentDexterity,
			CurrentIntelligence: targetCharacterInstanceViewRec.CurrentIntelligence,
			Health:              targetCharacterInstanceViewRec.Health,
			Fatigue:             targetCharacterInstanceViewRec.Fatigue,
			CurrentHealth:       targetCharacterInstanceViewRec.CurrentHealth,
			CurrentFatigue:      targetCharacterInstanceViewRec.CurrentFatigue,
		}

		err = m.CreateActionCharacterRec(rec)
		if err != nil {
			l.Warn("failed creating target action character record >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionCharacterRec = rec

		// Create action character object records
		objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(targetCharacterInstanceViewRec.ID)
		if err != nil {
			l.Warn("failed getting target character object records >%v<", err)
			return nil, err
		}

		l.Info("Adding >%d< target character object records", len(objectInstanceViewRecs))

		targetCharacterObjectRecs := []*record.ActionCharacterObject{}
		for _, objectInstanceViewRec := range objectInstanceViewRecs {
			l.Info("Adding target character object record >%v<", objectInstanceViewRecs)
			dungeonCharacterObjectRec := record.ActionCharacterObject{
				ActionID:            actionRec.ID,
				CharacterInstanceID: targetCharacterInstanceViewRec.ID,
				ObjectInstanceID:    objectInstanceViewRec.ID,
				Name:                objectInstanceViewRec.Name,
				IsEquipped:          objectInstanceViewRec.IsEquipped,
				IsStashed:           objectInstanceViewRec.IsStashed,
			}
			err := m.CreateActionCharacterObjectRec(&dungeonCharacterObjectRec)
			if err != nil {
				l.Warn("failed creating source action character object record >%v<", err)
				return nil, err
			}
			targetCharacterObjectRecs = append(targetCharacterObjectRecs, &dungeonCharacterObjectRec)
		}
		actionRecordSet.TargetActionCharacterObjectRecs = targetCharacterObjectRecs
	}

	// Create the target dungeon monster action record
	if actionRec.ResolvedTargetMonsterInstanceID.Valid {

		l.Info("Resolved target monster ID >%s<", actionRec.ResolvedTargetMonsterInstanceID.String)

		targetMonsterInstanceViewRec, err := m.GetMonsterInstanceViewRec(actionRec.ResolvedTargetMonsterInstanceID.String)
		if err != nil {
			l.Warn("failed getting target monster instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionMonster{
			RecordType:          record.ActionMonsterRecordTypeTarget,
			ActionID:            actionRec.ID,
			LocationInstanceID:  locationInstanceViewRec.ID,
			MonsterInstanceID:   targetMonsterInstanceViewRec.ID,
			Name:                targetMonsterInstanceViewRec.Name,
			Strength:            targetMonsterInstanceViewRec.Strength,
			Dexterity:           targetMonsterInstanceViewRec.Dexterity,
			Intelligence:        targetMonsterInstanceViewRec.Intelligence,
			CurrentStrength:     targetMonsterInstanceViewRec.CurrentStrength,
			CurrentDexterity:    targetMonsterInstanceViewRec.CurrentDexterity,
			CurrentIntelligence: targetMonsterInstanceViewRec.CurrentIntelligence,
			Health:              targetMonsterInstanceViewRec.Health,
			Fatigue:             targetMonsterInstanceViewRec.Fatigue,
			CurrentHealth:       targetMonsterInstanceViewRec.CurrentHealth,
			CurrentFatigue:      targetMonsterInstanceViewRec.CurrentFatigue,
		}

		err = m.CreateActionMonsterRec(rec)
		if err != nil {
			l.Warn("failed creating target action monster record >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionMonsterRec = rec

		// Create action monster object records
		objectInstanceViewRecs, err := m.GetMonsterInstanceEquippedObjectInstanceViewRecs(targetMonsterInstanceViewRec.ID)
		if err != nil {
			l.Warn("failed getting target monster object records >%v<", err)
			return nil, err
		}

		l.Info("Adding >%d< target monster object records", len(objectInstanceViewRecs))

		targetMonsterObjectRecs := []*record.ActionMonsterObject{}
		for _, objectInstanceViewRec := range objectInstanceViewRecs {
			l.Info("Adding target monster object record >%v<", objectInstanceViewRec)
			dungeonMonsterObjectRec := record.ActionMonsterObject{
				ActionID:          actionRec.ID,
				MonsterInstanceID: targetMonsterInstanceViewRec.ID,
				ObjectInstanceID:  objectInstanceViewRec.ID,
				Name:              objectInstanceViewRec.Name,
				IsEquipped:        objectInstanceViewRec.IsEquipped,
				IsStashed:         objectInstanceViewRec.IsStashed,
			}
			err := m.CreateActionMonsterObjectRec(&dungeonMonsterObjectRec)
			if err != nil {
				l.Warn("failed creating source action monster object record >%v<", err)
				return nil, err
			}
			targetMonsterObjectRecs = append(targetMonsterObjectRecs, &dungeonMonsterObjectRec)
		}
		actionRecordSet.TargetActionMonsterObjectRecs = targetMonsterObjectRecs
	}

	// Create the target dungeon object action record
	if actionRec.ResolvedTargetObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedTargetObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting target object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeTarget,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating target action object record >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionObjectRec = rec
	}

	// Create the stashed dungeon object action record
	if actionRec.ResolvedStashedObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedStashedObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting stashed object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeStashed,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating stashed action object record >%v<", err)
			return nil, err
		}
		actionRecordSet.StashedActionObjectRec = rec
	}

	// Create the equipped dungeon object action record
	if actionRec.ResolvedEquippedObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedEquippedObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting equipped object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeEquipped,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating equipped action object record >%v<", err)
			return nil, err
		}
		actionRecordSet.EquippedActionObjectRec = rec
	}

	// Create the dropped dungeon object action record
	if actionRec.ResolvedDroppedObjectInstanceID.Valid {
		targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(actionRec.ResolvedDroppedObjectInstanceID.String)
		if err != nil {
			l.Warn("failed getting dropped object instance view record >%v<", err)
			return nil, err
		}

		rec := &record.ActionObject{
			RecordType:         record.ActionObjectRecordTypeDropped,
			ActionID:           actionRec.ID,
			LocationInstanceID: locationInstanceViewRec.ID,
			ObjectInstanceID:   targetObjectInstanceViewRec.ID,
			Name:               targetObjectInstanceViewRec.Name,
			Description:        targetObjectInstanceViewRec.Description,
			IsStashed:          targetObjectInstanceViewRec.IsStashed,
			IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
		}

		err = m.CreateActionObjectRec(rec)
		if err != nil {
			l.Warn("failed creating dropped action object record >%v<", err)
			return nil, err
		}

		l.Warn("Assigning dropped object >%v< to action record set", rec)

		actionRecordSet.DroppedActionObjectRec = rec
	}

	return &actionRecordSet, nil
}

func (m *Model) GetActionRecordSet(actionID string) (*record.ActionRecordSet, error) {

	l := m.Logger("GetActionRecordSet")

	actionRecordSet := record.ActionRecordSet{}

	actionRec, err := m.GetActionRec(actionID, false)
	if err != nil {
		l.Warn("failed getting action record >%v<", err)
		return nil, err
	}
	actionRecordSet.ActionRec = actionRec

	// Add the source action character record that performed the action.
	if actionRec.CharacterInstanceID.Valid {
		actionCharacterRecs, err := m.GetActionCharacterRecs(
			map[string]interface{}{
				"record_type":           record.ActionCharacterRecordTypeSource,
				"action_id":             actionID,
				"character_instance_id": actionRec.CharacterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting action character records >%v<", err)
			return nil, err
		}
		if len(actionCharacterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action character records returned >%d<", len(actionCharacterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.ActionCharacterRec = actionCharacterRecs[0]

		actionCharacterObjectRecs, err := m.GetActionCharacterObjectRecs(
			map[string]interface{}{
				"action_id":             actionID,
				"character_instance_id": actionRec.CharacterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting action character object records >%v<", err)
			return nil, err
		}
		actionRecordSet.ActionCharacterObjectRecs = actionCharacterObjectRecs
	}

	// Add the source action monster record that performed the action.
	if actionRec.MonsterInstanceID.Valid {
		actionMonsterRecs, err := m.GetActionMonsterRecs(
			map[string]interface{}{
				"record_type":         record.ActionMonsterRecordTypeSource,
				"action_id":           actionID,
				"monster_instance_id": actionRec.MonsterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting action monster records >%v<", err)
			return nil, err
		}
		if len(actionMonsterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action monster records returned >%d<", len(actionMonsterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.ActionMonsterRec = actionMonsterRecs[0]

		actionMonsterObjectRecs, err := m.GetActionMonsterObjectRecs(
			map[string]interface{}{
				"action_id":           actionID,
				"monster_instance_id": actionRec.MonsterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting action monster object records >%v<", err)
			return nil, err
		}
		actionRecordSet.ActionMonsterObjectRecs = actionMonsterObjectRecs
	}

	// Add the current location record set where the action was performed.
	locationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.LocationInstanceID)
	if err != nil {
		l.Warn("failed getting location instance view record >%v<", err)
		return nil, err
	}

	currentLocationRecordSet := record.ActionLocationRecordSet{
		LocationInstanceViewRec: locationInstanceViewRec,
		ActionCharacterRecs:     []*record.ActionCharacter{},
		ActionMonsterRecs:       []*record.ActionMonster{},
		ActionObjectRecs:        []*record.ActionObject{},
	}

	// Add the current location occupant action character records
	actionCharacterRecs, err := m.GetActionCharacterRecs(
		map[string]interface{}{
			"record_type":          record.ActionCharacterRecordTypeOccupant,
			"action_id":            actionID,
			"location_instance_id": locationInstanceViewRec.ID,
		},
		nil,
		false,
	)
	if err != nil {
		l.Warn("failed getting current location occupant action character records >%v<", err)
		return nil, err
	}
	currentLocationRecordSet.ActionCharacterRecs = actionCharacterRecs

	// Add the current location occupant action monster records
	actionMonsterRecs, err := m.GetActionMonsterRecs(
		map[string]interface{}{
			"record_type":          record.ActionMonsterRecordTypeOccupant,
			"action_id":            actionID,
			"location_instance_id": locationInstanceViewRec.ID,
		},
		nil,
		false,
	)
	if err != nil {
		l.Warn("failed getting current location occupant action monster records >%v<", err)
		return nil, err
	}
	currentLocationRecordSet.ActionMonsterRecs = actionMonsterRecs

	// Add the current location occupant action object records
	dungeonActionObjectRecs, err := m.GetActionObjectRecs(
		map[string]interface{}{
			"record_type":          record.ActionObjectRecordTypeOccupant,
			"action_id":            actionID,
			"location_instance_id": locationInstanceViewRec.ID,
		},
		nil,
		false,
	)
	if err != nil {
		l.Warn("failed getting current location occupant action monster records >%v<", err)
		return nil, err
	}
	currentLocationRecordSet.ActionObjectRecs = dungeonActionObjectRecs

	actionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target dungeon location record set when set
	if actionRec.ResolvedTargetLocationInstanceID.Valid {

		// Add the target location record set when the action was performed.
		locationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.ResolvedTargetLocationInstanceID.String)
		if err != nil {
			l.Warn("failed getting target location instance view record >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := record.ActionLocationRecordSet{
			LocationInstanceViewRec: locationInstanceViewRec,
			ActionCharacterRecs:     []*record.ActionCharacter{},
			ActionMonsterRecs:       []*record.ActionMonster{},
			ActionObjectRecs:        []*record.ActionObject{},
		}

		// Add the target location occupant action character records
		actionCharacterRecs, err := m.GetActionCharacterRecs(
			map[string]interface{}{
				"record_type":          record.ActionCharacterRecordTypeOccupant,
				"action_id":            actionID,
				"location_instance_id": locationInstanceViewRec.ID,
			},
			nil,
			false,
		)
		if err != nil {
			l.Warn("failed getting target location occupant action character records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionCharacterRecs = actionCharacterRecs

		// Add the target location occupant action monster records
		actionMonsterRecs, err := m.GetActionMonsterRecs(
			map[string]interface{}{
				"record_type":          record.ActionMonsterRecordTypeOccupant,
				"action_id":            actionID,
				"location_instance_id": locationInstanceViewRec.ID,
			},
			nil,
			false,
		)
		if err != nil {
			l.Warn("failed getting target location occupant action monster records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionMonsterRecs = actionMonsterRecs

		// Add the target location occupant action object records
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			map[string]interface{}{
				"record_type":          record.ActionObjectRecordTypeOccupant,
				"action_id":            actionID,
				"location_instance_id": locationInstanceViewRec.ID,
			},
			nil,
			false,
		)
		if err != nil {
			l.Warn("failed getting target location occupant action monster records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionObjectRecs = dungeonActionObjectRecs

		actionRecordSet.TargetLocation = &targetLocationRecordSet
	}

	// Get the target character action record
	if actionRec.ResolvedTargetCharacterInstanceID.Valid {
		actionCharacterRecs, err := m.GetActionCharacterRecs(
			map[string]interface{}{
				"record_type":           record.ActionCharacterRecordTypeTarget,
				"action_id":             actionID,
				"character_instance_id": actionRec.ResolvedTargetCharacterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting target action character record >%v<", err)
			return nil, err
		}
		if len(actionCharacterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action character records returned >%d<", len(actionCharacterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.TargetActionCharacterRec = actionCharacterRecs[0]

		actionCharacterObjectRecs, err := m.GetActionCharacterObjectRecs(
			map[string]interface{}{
				"action_id":             actionID,
				"character_instance_id": actionRec.ResolvedTargetCharacterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting target character object records >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionCharacterObjectRecs = actionCharacterObjectRecs
	}

	// Get the target dungeon monster action record
	if actionRec.ResolvedTargetMonsterInstanceID.Valid {
		actionMonsterRecs, err := m.GetActionMonsterRecs(
			map[string]interface{}{
				"record_type":         record.ActionMonsterRecordTypeTarget,
				"action_id":           actionID,
				"monster_instance_id": actionRec.ResolvedTargetMonsterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting target action monster record >%v<", err)
			return nil, err
		}
		if len(actionMonsterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action monster records returned >%d<", len(actionMonsterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.TargetActionMonsterRec = actionMonsterRecs[0]

		actionMonsterObjectRecs, err := m.GetActionMonsterObjectRecs(
			map[string]interface{}{
				"action_id":           actionID,
				"monster_instance_id": actionRec.ResolvedTargetMonsterInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting target monster object records >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionMonsterObjectRecs = actionMonsterObjectRecs
	}

	// Get the target dungeon object action record
	if actionRec.ResolvedTargetObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			map[string]interface{}{
				"record_type":        record.ActionObjectRecordTypeTarget,
				"action_id":          actionID,
				"object_instance_id": actionRec.ResolvedTargetObjectInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting target action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.TargetActionObjectRec = dungeonActionObjectRecs[0]
	}

	// Get the stashed dungeon object action record
	if actionRec.ResolvedStashedObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			map[string]interface{}{
				"record_type":        record.ActionObjectRecordTypeStashed,
				"action_id":          actionID,
				"object_instance_id": actionRec.ResolvedStashedObjectInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting stashed action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.StashedActionObjectRec = dungeonActionObjectRecs[0]
	}

	// Get the equipped dungeon object action record
	if actionRec.ResolvedEquippedObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			map[string]interface{}{
				"record_type":        record.ActionObjectRecordTypeEquipped,
				"action_id":          actionID,
				"object_instance_id": actionRec.ResolvedEquippedObjectInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting equipped action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.EquippedActionObjectRec = dungeonActionObjectRecs[0]
	}

	// Get the dropped dungeon object action record
	if actionRec.ResolvedDroppedObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			map[string]interface{}{
				"record_type":        record.ActionObjectRecordTypeDropped,
				"action_id":          actionID,
				"object_instance_id": actionRec.ResolvedDroppedObjectInstanceID.String,
			}, nil, false)
		if err != nil {
			l.Warn("failed getting dropped action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.DroppedActionObjectRec = dungeonActionObjectRecs[0]
	}

	return &actionRecordSet, nil
}

func (m *Model) GetLocationInstanceViewRecordSet(locationInstanceID string, forUpdate bool) (*record.LocationInstanceViewRecordSet, error) {

	l := m.Logger("GetLocationInstanceViewRecordSet")

	locationInstanceRecordSet := &record.LocationInstanceViewRecordSet{}

	// Location record
	locationInstanceViewRec, err := m.GetLocationInstanceViewRec(locationInstanceID)
	if err != nil {
		l.Warn("failed to get dungeon location record >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.LocationInstanceViewRec = locationInstanceViewRec

	// All characters at location
	characterInstanceViewRecs, err := m.GetCharacterInstanceViewRecs(
		map[string]interface{}{
			"location_instance_id": locationInstanceViewRec.ID,
		},
		nil,
	)
	if err != nil {
		l.Warn("failed to get dungeon location character records >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.CharacterInstanceViewRecs = characterInstanceViewRecs

	// All monsters at location
	monsterInstanceViewRecs, err := m.GetMonsterInstanceViewRecs(
		map[string]interface{}{
			"location_instance_id": locationInstanceViewRec.ID,
		},
		nil,
	)
	if err != nil {
		l.Warn("failed to get dungeon location monster records >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.MonsterInstanceViewRecs = monsterInstanceViewRecs

	// All objects at location
	objectInstanceViewRecs, err := m.GetObjectInstanceViewRecs(
		map[string]interface{}{
			"location_instance_id": locationInstanceViewRec.ID,
		},
		nil,
	)
	if err != nil {
		l.Warn("failed to get dungeon location object records >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.ObjectInstanceViewRecs = objectInstanceViewRecs

	locationInstanceIDs := []string{}
	if locationInstanceViewRec.NorthLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.NorthLocationInstanceID.String)
	}
	if locationInstanceViewRec.NortheastLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.NortheastLocationInstanceID.String)
	}
	if locationInstanceViewRec.EastLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.EastLocationInstanceID.String)
	}
	if locationInstanceViewRec.SoutheastLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.SoutheastLocationInstanceID.String)
	}
	if locationInstanceViewRec.SouthLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.SouthLocationInstanceID.String)
	}
	if locationInstanceViewRec.SouthwestLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.SouthwestLocationInstanceID.String)
	}
	if locationInstanceViewRec.WestLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.WestLocationInstanceID.String)
	}
	if locationInstanceViewRec.NorthwestLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.NorthwestLocationInstanceID.String)
	}
	if locationInstanceViewRec.UpLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.UpLocationInstanceID.String)
	}
	if locationInstanceViewRec.DownLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.DownLocationInstanceID.String)
	}

	locationInstanceViewRecs, err := m.GetLocationInstanceViewRecs(
		map[string]interface{}{
			"id": locationInstanceIDs,
		},
		nil,
	)
	if err != nil {
		l.Warn("failed to get dungeon location direction location records >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.LocationInstanceViewRecs = locationInstanceViewRecs

	return locationInstanceRecordSet, nil
}
