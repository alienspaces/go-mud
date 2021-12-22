package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

type DungeonActionRecordSet struct {
	ActionRec                *record.DungeonAction
	ActionCharacterRec       *record.DungeonActionCharacter
	ActionMonsterRec         *record.DungeonActionMonster
	CurrentLocation          *DungeonActionLocationRecordSet
	EquippedActionObjectRec  *record.DungeonActionObject
	StashedActionObjectRec   *record.DungeonActionObject
	TargetActionObjectRec    *record.DungeonActionObject
	TargetActionCharacterRec *record.DungeonActionCharacter
	TargetActionMonsterRec   *record.DungeonActionMonster
	TargetLocation           *DungeonActionLocationRecordSet
}

type DungeonActionLocationRecordSet struct {
	LocationRec         *record.DungeonLocation
	ActionCharacterRecs []*record.DungeonActionCharacter
	ActionMonsterRecs   []*record.DungeonActionMonster
	ActionObjectRecs    []*record.DungeonActionObject
}

type DungeonLocationRecordSet struct {
	LocationRec   *record.DungeonLocation
	CharacterRecs []*record.DungeonCharacter
	MonsterRecs   []*record.DungeonMonster
	ObjectRecs    []*record.DungeonObject
	LocationRecs  []*record.DungeonLocation
}

// ProcessDungeonCharacterAction - Processes a submitted character action
func (m *Model) ProcessDungeonCharacterAction(dungeonID string, dungeonCharacterID string, sentence string) (*DungeonActionRecordSet, error) {

	m.Log.Info("Processing dungeon character ID >%s< action command >%s<", dungeonCharacterID, sentence)

	// Verify the character performing the action exists within the specified dungeon
	sourceCharacterRec, err := m.GetDungeonCharacterRec(dungeonCharacterID, true)
	if err != nil {
		m.Log.Warn("failed getting dungeon character record before performing action >%v<", err)
		return nil, err
	}
	if sourceCharacterRec == nil {
		msg := fmt.Sprintf("failed getting dungeon character record ID >%s< before performing action", dungeonCharacterID)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if sourceCharacterRec.DungeonID != dungeonID {
		msg := fmt.Sprintf("dungeon character ID >%s< does not exist in dungeon ID >%s<", dungeonCharacterID, dungeonID)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Get the current dungeon location set of related records
	dungeonLocationRecordSet, err := m.GetDungeonLocationRecordSet(sourceCharacterRec.DungeonLocationID, true)
	if err != nil {
		m.Log.Warn("failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}
	if dungeonLocationRecordSet == nil {
		msg := fmt.Sprintf("failed getting dungeon location record ID >%s< set before performing action", sourceCharacterRec.DungeonLocationID)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Resolve the submitted character action
	dungeonActionRec, err := m.resolveAction(sentence, sourceCharacterRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("failed resolving dungeon character action >%v<", err)
		return nil, err
	}

	// TODO: Debug dungeon action record here..
	m.Log.Info("Dungeon action record command >%s<", dungeonActionRec.ResolvedCommand)
	m.Log.Info("Dungeon action record location >%s<", dungeonActionRec.DungeonLocationID)

	// Perform the submitted character action
	dungeonActionRec, err = m.performDungeonCharacterAction(sourceCharacterRec, dungeonActionRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("failed performing dungeon character action >%v<", err)
		return nil, err
	}

	m.Log.Info("Dungeon action record command >%s<", dungeonActionRec.ResolvedCommand)
	m.Log.Info("Dungeon action record location >%s<", dungeonActionRec.DungeonLocationID)

	// Create the resulting dungeon action event record
	err = m.CreateDungeonActionRec(dungeonActionRec)
	if err != nil {
		m.Log.Warn("failed creating dungeon action record >%v<", err)
		return nil, err
	}

	m.Log.Info("Created dungeon action record ID >%s<", dungeonActionRec.ID)

	// TODO: Maybe don't need to do this... Get the updated character record
	sourceCharacterRec, err = m.GetDungeonCharacterRec(dungeonCharacterID, false)
	if err != nil {
		m.Log.Warn("failed getting dungeon character record after performing action >%v<", err)
		return nil, err
	}

	// Create source dungeon action character
	dungeonActionCharacterRec := record.DungeonActionCharacter{
		RecordType:         record.DungeonActionCharacterRecordTypeSource,
		DungeonActionID:    dungeonActionRec.ID,
		DungeonLocationID:  dungeonActionRec.DungeonLocationID,
		DungeonCharacterID: sourceCharacterRec.ID,
		Name:               sourceCharacterRec.Name,
		Strength:           sourceCharacterRec.Strength,
		Dexterity:          sourceCharacterRec.Dexterity,
		Intelligence:       sourceCharacterRec.Intelligence,
		Health:             sourceCharacterRec.Health,
		Fatigue:            sourceCharacterRec.Fatigue,
	}

	err = m.CreateDungeonActionCharacterRec(&dungeonActionCharacterRec)
	if err != nil {
		m.Log.Warn("failed creating source dungeon action character record >%v<", err)
		return nil, err
	}

	dungeonActionRecordSet := DungeonActionRecordSet{
		ActionRec:          dungeonActionRec,
		ActionCharacterRec: &dungeonActionCharacterRec,
	}

	// Get the updated current dungeon location record set
	dungeonLocationRecordSet, err = m.GetDungeonLocationRecordSet(dungeonActionRec.DungeonLocationID, true)
	if err != nil {
		m.Log.Warn("failed getting dungeon location record set after performing action >%v<", err)
		return nil, err
	}

	m.Log.Info("Dungeon location record set location name >%s<", dungeonLocationRecordSet.LocationRec.Name)
	m.Log.Info("Dungeon location record set characters >%d<", len(dungeonLocationRecordSet.CharacterRecs))
	m.Log.Info("Dungeon location record set monsters >%d<", len(dungeonLocationRecordSet.MonsterRecs))
	m.Log.Info("Dungeon location record set objects >%d<", len(dungeonLocationRecordSet.ObjectRecs))

	// Current location
	dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.DungeonLocationID, false)
	if err != nil {
		m.Log.Warn("failed getting dungeon location record after performing action >%v<", err)
		return nil, err
	}

	currentLocationRecordSet := DungeonActionLocationRecordSet{
		LocationRec:         dungeonLocationRec,
		ActionCharacterRecs: []*record.DungeonActionCharacter{},
		ActionMonsterRecs:   []*record.DungeonActionMonster{},
		ActionObjectRecs:    []*record.DungeonActionObject{},
	}

	// Create the dungeon action character record for each character now at the current location
	if len(dungeonLocationRecordSet.CharacterRecs) > 0 {
		for _, characterRec := range dungeonLocationRecordSet.CharacterRecs {
			dungeonActionCharacterRec := record.DungeonActionCharacter{
				RecordType:         record.DungeonActionCharacterRecordTypeOccupant,
				DungeonActionID:    dungeonActionRec.ID,
				DungeonLocationID:  dungeonLocationRec.ID,
				DungeonCharacterID: characterRec.ID,
				Name:               characterRec.Name,
				Strength:           characterRec.Strength,
				Dexterity:          characterRec.Dexterity,
				Intelligence:       characterRec.Intelligence,
				Health:             characterRec.Health,
				Fatigue:            characterRec.Fatigue,
			}

			err := m.CreateDungeonActionCharacterRec(&dungeonActionCharacterRec)
			if err != nil {
				m.Log.Warn("failed creating current location dungeon action character record >%v<", err)
				return nil, err
			}

			m.Log.Info("Created current location dungeon action character record ID >%s<", dungeonActionCharacterRec.ID)
			currentLocationRecordSet.ActionCharacterRecs = append(currentLocationRecordSet.ActionCharacterRecs, &dungeonActionCharacterRec)
		}
	}

	// Create the dungeon action monster record for each monster now at the current location
	if len(dungeonLocationRecordSet.MonsterRecs) > 0 {
		for _, monsterRec := range dungeonLocationRecordSet.MonsterRecs {
			dungeonActionMonsterRec := record.DungeonActionMonster{
				RecordType:        record.DungeonActionMonsterRecordTypeOccupant,
				DungeonActionID:   dungeonActionRec.ID,
				DungeonLocationID: dungeonLocationRec.ID,
				DungeonMonsterID:  monsterRec.ID,
				Name:              monsterRec.Name,
				Strength:          monsterRec.Strength,
				Dexterity:         monsterRec.Dexterity,
				Intelligence:      monsterRec.Intelligence,
				Health:            monsterRec.Health,
				Fatigue:           monsterRec.Fatigue,
			}
			err := m.CreateDungeonActionMonsterRec(&dungeonActionMonsterRec)
			if err != nil {
				m.Log.Warn("failed creating current location dungeon action monster record >%v<", err)
				return nil, err
			}

			m.Log.Info("Created current location dungeon action monster record ID >%s<", dungeonActionMonsterRec.ID)
			currentLocationRecordSet.ActionMonsterRecs = append(currentLocationRecordSet.ActionMonsterRecs, &dungeonActionMonsterRec)
		}
	}

	// Create the dungeon action object record for each object now at the current location
	if len(dungeonLocationRecordSet.ObjectRecs) > 0 {
		for _, objectRec := range dungeonLocationRecordSet.ObjectRecs {
			dungeonActionObjectRec := record.DungeonActionObject{
				RecordType:        record.DungeonActionObjectRecordTypeOccupant,
				DungeonActionID:   dungeonActionRec.ID,
				DungeonLocationID: dungeonLocationRec.ID,
				DungeonObjectID:   objectRec.ID,
				Name:              objectRec.Name,
			}
			err := m.CreateDungeonActionObjectRec(&dungeonActionObjectRec)
			if err != nil {
				m.Log.Warn("failed creating current location dungeon action object record >%v<", err)
				return nil, err
			}

			m.Log.Info("Created current location dungeon action object record ID >%s<", dungeonActionObjectRec.ID)
			currentLocationRecordSet.ActionObjectRecs = append(currentLocationRecordSet.ActionObjectRecs, &dungeonActionObjectRec)
		}
	}

	dungeonActionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target dungeon location record set when set
	if dungeonActionRec.ResolvedTargetDungeonLocationID.Valid {

		dungeonLocationRecordSet, err := m.GetDungeonLocationRecordSet(dungeonActionRec.ResolvedTargetDungeonLocationID.String, true)
		if err != nil {
			m.Log.Warn("failed getting target dungeon location record set after performing action >%v<", err)
			return nil, err
		}

		// Target location
		dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.ResolvedTargetDungeonLocationID.String, false)
		if err != nil {
			m.Log.Warn("failed getting dungeon location record after performing action >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := DungeonActionLocationRecordSet{
			LocationRec:         dungeonLocationRec,
			ActionCharacterRecs: []*record.DungeonActionCharacter{},
			ActionMonsterRecs:   []*record.DungeonActionMonster{},
			ActionObjectRecs:    []*record.DungeonActionObject{},
		}

		// Create the dungeon action character record for each character at the target location
		if len(dungeonLocationRecordSet.CharacterRecs) > 0 {
			for _, characterRec := range dungeonLocationRecordSet.CharacterRecs {
				dungeonActionCharacterRec := record.DungeonActionCharacter{
					RecordType:         record.DungeonActionCharacterRecordTypeOccupant,
					DungeonActionID:    dungeonActionRec.ID,
					DungeonLocationID:  dungeonLocationRec.ID,
					DungeonCharacterID: characterRec.ID,
					Name:               characterRec.Name,
					Strength:           characterRec.Strength,
					Dexterity:          characterRec.Dexterity,
					Intelligence:       characterRec.Intelligence,
					Health:             characterRec.Health,
					Fatigue:            characterRec.Fatigue,
				}

				err := m.CreateDungeonActionCharacterRec(&dungeonActionCharacterRec)
				if err != nil {
					m.Log.Warn("failed creating target location occupant dungeon action character record >%v<", err)
					return nil, err
				}

				m.Log.Info("Created target location occupant dungeon action character record ID >%s<", dungeonActionCharacterRec.ID)
				targetLocationRecordSet.ActionCharacterRecs = append(targetLocationRecordSet.ActionCharacterRecs, &dungeonActionCharacterRec)
			}
		}

		// Create the dungeon action monster record for each monster at the target location
		if len(dungeonLocationRecordSet.MonsterRecs) > 0 {
			for _, monsterRec := range dungeonLocationRecordSet.MonsterRecs {
				dungeonActionMonsterRec := record.DungeonActionMonster{
					RecordType:        record.DungeonActionMonsterRecordTypeOccupant,
					DungeonActionID:   dungeonActionRec.ID,
					DungeonLocationID: dungeonLocationRec.ID,
					DungeonMonsterID:  monsterRec.ID,
					Name:              monsterRec.Name,
					Strength:          monsterRec.Strength,
					Dexterity:         monsterRec.Dexterity,
					Intelligence:      monsterRec.Intelligence,
					Health:            monsterRec.Health,
					Fatigue:           monsterRec.Fatigue,
				}
				err := m.CreateDungeonActionMonsterRec(&dungeonActionMonsterRec)
				if err != nil {
					m.Log.Warn("failed creating target location occupant dungeon action monster record >%v<", err)
					return nil, err
				}

				m.Log.Info("Created target location occupant dungeon action monster record ID >%s<", dungeonActionMonsterRec.ID)
				targetLocationRecordSet.ActionMonsterRecs = append(targetLocationRecordSet.ActionMonsterRecs, &dungeonActionMonsterRec)
			}
		}

		// Create the dungeon action object record for each object at the target location
		if len(dungeonLocationRecordSet.ObjectRecs) > 0 {
			for _, objectRec := range dungeonLocationRecordSet.ObjectRecs {
				dungeonActionObjectRec := record.DungeonActionObject{
					RecordType:        record.DungeonActionObjectRecordTypeOccupant,
					DungeonActionID:   dungeonActionRec.ID,
					DungeonLocationID: dungeonLocationRec.ID,
					DungeonObjectID:   objectRec.ID,
					Name:              objectRec.Name,
				}
				err := m.CreateDungeonActionObjectRec(&dungeonActionObjectRec)
				if err != nil {
					m.Log.Warn("failed creating target location occupant dungeon action object record >%v<", err)
					return nil, err
				}

				m.Log.Info("Created target location occupant dungeon action object record ID >%s<", dungeonActionObjectRec.ID)
				targetLocationRecordSet.ActionObjectRecs = append(targetLocationRecordSet.ActionObjectRecs, &dungeonActionObjectRec)
			}
		}

		dungeonActionRecordSet.TargetLocation = &targetLocationRecordSet
	}

	// Create the target dungeon character action record
	if dungeonActionRec.ResolvedTargetDungeonCharacterID.Valid {
		targetCharacterRec, err := m.GetDungeonCharacterRec(dungeonActionRec.ResolvedTargetDungeonCharacterID.String, false)
		if err != nil {
			m.Log.Warn("failed getting target dungeon character record >%v<", err)
			return nil, err
		}

		rec := &record.DungeonActionCharacter{
			RecordType:         record.DungeonActionCharacterRecordTypeTarget,
			DungeonActionID:    dungeonActionRec.ID,
			DungeonLocationID:  dungeonLocationRec.ID,
			DungeonCharacterID: targetCharacterRec.ID,
			Name:               targetCharacterRec.Name,
			Strength:           targetCharacterRec.Strength,
			Dexterity:          targetCharacterRec.Dexterity,
			Intelligence:       targetCharacterRec.Intelligence,
			Health:             targetCharacterRec.Health,
			Fatigue:            targetCharacterRec.Fatigue,
		}

		err = m.CreateDungeonActionCharacterRec(rec)
		if err != nil {
			m.Log.Warn("failed creating target dungeon action character record >%v<", err)
			return nil, err
		}
		dungeonActionRecordSet.TargetActionCharacterRec = rec
	}

	// Create the target dungeon monster action record
	if dungeonActionRec.ResolvedTargetDungeonMonsterID.Valid {
		targetMonsterRec, err := m.GetDungeonMonsterRec(dungeonActionRec.ResolvedTargetDungeonMonsterID.String, false)
		if err != nil {
			m.Log.Warn("failed getting target dungeon monster record >%v<", err)
			return nil, err
		}

		rec := &record.DungeonActionMonster{
			RecordType:        record.DungeonActionMonsterRecordTypeTarget,
			DungeonActionID:   dungeonActionRec.ID,
			DungeonLocationID: dungeonLocationRec.ID,
			DungeonMonsterID:  targetMonsterRec.ID,
			Name:              targetMonsterRec.Name,
			Strength:          targetMonsterRec.Strength,
			Dexterity:         targetMonsterRec.Dexterity,
			Intelligence:      targetMonsterRec.Intelligence,
			Health:            targetMonsterRec.Health,
			Fatigue:           targetMonsterRec.Fatigue,
		}

		err = m.CreateDungeonActionMonsterRec(rec)
		if err != nil {
			m.Log.Warn("failed creating target dungeon action monster record >%v<", err)
			return nil, err
		}
		dungeonActionRecordSet.TargetActionMonsterRec = rec
	}

	// Create the target dungeon object action record
	if dungeonActionRec.ResolvedTargetDungeonObjectID.Valid {
		targetObjectRec, err := m.GetDungeonObjectRec(dungeonActionRec.ResolvedTargetDungeonObjectID.String, false)
		if err != nil {
			m.Log.Warn("failed getting target dungeon object record >%v<", err)
			return nil, err
		}

		rec := &record.DungeonActionObject{
			RecordType:        record.DungeonActionObjectRecordTypeTarget,
			DungeonActionID:   dungeonActionRec.ID,
			DungeonLocationID: dungeonLocationRec.ID,
			DungeonObjectID:   targetObjectRec.ID,
			Name:              targetObjectRec.Name,
		}

		err = m.CreateDungeonActionObjectRec(rec)
		if err != nil {
			m.Log.Warn("failed creating target dungeon action object record >%v<", err)
			return nil, err
		}
		dungeonActionRecordSet.TargetActionObjectRec = rec
	}

	// TODO:
	// EquippedObjectRec   *record.DungeonObject
	// StashedObjectRec    *record.DungeonObject

	return &dungeonActionRecordSet, nil
}

func (m *Model) GetDungeonActionRecordSet(dungeonActionID string) (*DungeonActionRecordSet, error) {

	dungeonActionRecordSet := DungeonActionRecordSet{}

	dungeonActionRec, err := m.GetDungeonActionRec(dungeonActionID, false)
	if err != nil {
		m.Log.Warn("failed getting dungeon action record >%v<", err)
		return nil, err
	}
	dungeonActionRecordSet.ActionRec = dungeonActionRec

	// Add the source dungeon action character record that performed the dungeon action.
	if dungeonActionRec.DungeonCharacterID.Valid {
		dungeonActionCharacterRecs, err := m.GetDungeonActionCharacterRecs(
			map[string]interface{}{
				"record_type":          record.DungeonActionCharacterRecordTypeSource,
				"dungeon_action_id":    dungeonActionID,
				"dungeon_character_id": dungeonActionRec.DungeonCharacterID.String,
			}, nil, false)
		if err != nil {
			m.Log.Warn("failed getting dungeon action character records >%v<", err)
			return nil, err
		}
		if len(dungeonActionCharacterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of dungeon action character records returned >%d<", len(dungeonActionCharacterRecs))
			m.Log.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		dungeonActionRecordSet.ActionCharacterRec = dungeonActionCharacterRecs[0]
	}

	// Add the source dungeon action monster record that performed the dungeon action.
	if dungeonActionRec.DungeonMonsterID.Valid {
		dungeonActionMonsterRecs, err := m.GetDungeonActionMonsterRecs(
			map[string]interface{}{
				"record_type":        record.DungeonActionMonsterRecordTypeSource,
				"dungeon_action_id":  dungeonActionID,
				"dungeon_monster_id": dungeonActionRec.DungeonMonsterID.String,
			}, nil, false)
		if err != nil {
			m.Log.Warn("failed getting dungeon action monster records >%v<", err)
			return nil, err
		}
		if len(dungeonActionMonsterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of dungeon action monster records returned >%d<", len(dungeonActionMonsterRecs))
			m.Log.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		dungeonActionRecordSet.ActionMonsterRec = dungeonActionMonsterRecs[0]
	}

	// Add the current location record set when the action was performed.
	dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.DungeonLocationID, false)
	if err != nil {
		m.Log.Warn("failed getting dungeon location record >%v<", err)
		return nil, err
	}

	currentLocationRecordSet := DungeonActionLocationRecordSet{
		LocationRec:         dungeonLocationRec,
		ActionCharacterRecs: []*record.DungeonActionCharacter{},
		ActionMonsterRecs:   []*record.DungeonActionMonster{},
		ActionObjectRecs:    []*record.DungeonActionObject{},
	}

	// Add the current location occupant dungeon action character records
	dungeonActionCharacterRecs, err := m.GetDungeonActionCharacterRecs(
		map[string]interface{}{
			"record_type":         record.DungeonActionCharacterRecordTypeOccupant,
			"dungeon_action_id":   dungeonActionID,
			"dungeon_location_id": dungeonLocationRec.ID,
		},
		nil,
		false,
	)
	if err != nil {
		m.Log.Warn("failed getting current location occupant dungeon action character records >%v<", err)
		return nil, err
	}
	currentLocationRecordSet.ActionCharacterRecs = dungeonActionCharacterRecs

	// Add the current location occupant dungeon action monster records
	dungeonActionMonsterRecs, err := m.GetDungeonActionMonsterRecs(
		map[string]interface{}{
			"record_type":         record.DungeonActionMonsterRecordTypeOccupant,
			"dungeon_action_id":   dungeonActionID,
			"dungeon_location_id": dungeonLocationRec.ID,
		},
		nil,
		false,
	)
	if err != nil {
		m.Log.Warn("failed getting current location occupant dungeon action monster records >%v<", err)
		return nil, err
	}
	currentLocationRecordSet.ActionMonsterRecs = dungeonActionMonsterRecs

	// Add the current location occupant dungeon action object records
	dungeonActionObjectRecs, err := m.GetDungeonActionObjectRecs(
		map[string]interface{}{
			"record_type":         record.DungeonActionObjectRecordTypeOccupant,
			"dungeon_action_id":   dungeonActionID,
			"dungeon_location_id": dungeonLocationRec.ID,
		},
		nil,
		false,
	)
	if err != nil {
		m.Log.Warn("failed getting current location occupant dungeon action monster records >%v<", err)
		return nil, err
	}
	currentLocationRecordSet.ActionObjectRecs = dungeonActionObjectRecs

	dungeonActionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target dungeon location record set when set
	if dungeonActionRec.ResolvedTargetDungeonLocationID.Valid {

		// Add the target location record set when the action was performed.
		dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.ResolvedTargetDungeonLocationID.String, false)
		if err != nil {
			m.Log.Warn("failed getting target dungeon location record >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := DungeonActionLocationRecordSet{
			LocationRec:         dungeonLocationRec,
			ActionCharacterRecs: []*record.DungeonActionCharacter{},
			ActionMonsterRecs:   []*record.DungeonActionMonster{},
			ActionObjectRecs:    []*record.DungeonActionObject{},
		}

		// Add the target location occupant dungeon action character records
		dungeonActionCharacterRecs, err := m.GetDungeonActionCharacterRecs(
			map[string]interface{}{
				"record_type":         record.DungeonActionCharacterRecordTypeOccupant,
				"dungeon_action_id":   dungeonActionID,
				"dungeon_location_id": dungeonLocationRec.ID,
			},
			nil,
			false,
		)
		if err != nil {
			m.Log.Warn("failed getting target location occupant dungeon action character records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionCharacterRecs = dungeonActionCharacterRecs

		// Add the target location occupant dungeon action monster records
		dungeonActionMonsterRecs, err := m.GetDungeonActionMonsterRecs(
			map[string]interface{}{
				"record_type":         record.DungeonActionMonsterRecordTypeOccupant,
				"dungeon_action_id":   dungeonActionID,
				"dungeon_location_id": dungeonLocationRec.ID,
			},
			nil,
			false,
		)
		if err != nil {
			m.Log.Warn("failed getting target location occupant dungeon action monster records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionMonsterRecs = dungeonActionMonsterRecs

		// Add the target location occupant dungeon action object records
		dungeonActionObjectRecs, err := m.GetDungeonActionObjectRecs(
			map[string]interface{}{
				"record_type":         record.DungeonActionObjectRecordTypeOccupant,
				"dungeon_action_id":   dungeonActionID,
				"dungeon_location_id": dungeonLocationRec.ID,
			},
			nil,
			false,
		)
		if err != nil {
			m.Log.Warn("failed getting target location occupant dungeon action monster records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionObjectRecs = dungeonActionObjectRecs

		dungeonActionRecordSet.TargetLocation = &targetLocationRecordSet
	}

	// Get the target dungeon character action record
	if dungeonActionRec.ResolvedTargetDungeonCharacterID.Valid {
		dungeonActionCharacterRecs, err := m.GetDungeonActionCharacterRecs(
			map[string]interface{}{
				"record_type":          record.DungeonActionCharacterRecordTypeTarget,
				"dungeon_action_id":    dungeonActionID,
				"dungeon_character_id": dungeonActionRec.ResolvedTargetDungeonCharacterID.String,
			}, nil, false)
		if err != nil {
			m.Log.Warn("failed getting target dungeon action character record >%v<", err)
			return nil, err
		}
		if len(dungeonActionCharacterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of dungeon action character records returned >%d<", len(dungeonActionCharacterRecs))
			m.Log.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		dungeonActionRecordSet.TargetActionCharacterRec = dungeonActionCharacterRecs[0]
	}

	// Get the target dungeon monster action record
	if dungeonActionRec.ResolvedTargetDungeonMonsterID.Valid {
		dungeonActionMonsterRecs, err := m.GetDungeonActionMonsterRecs(
			map[string]interface{}{
				"record_type":        record.DungeonActionMonsterRecordTypeTarget,
				"dungeon_action_id":  dungeonActionID,
				"dungeon_monster_id": dungeonActionRec.ResolvedTargetDungeonMonsterID.String,
			}, nil, false)
		if err != nil {
			m.Log.Warn("failed getting target dungeon action monster record >%v<", err)
			return nil, err
		}
		if len(dungeonActionMonsterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of dungeon action monster records returned >%d<", len(dungeonActionMonsterRecs))
			m.Log.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		dungeonActionRecordSet.TargetActionMonsterRec = dungeonActionMonsterRecs[0]
	}

	// Get the target dungeon object action record
	if dungeonActionRec.ResolvedTargetDungeonObjectID.Valid {
		dungeonActionObjectRecs, err := m.GetDungeonActionObjectRecs(
			map[string]interface{}{
				"record_type":       record.DungeonActionObjectRecordTypeTarget,
				"dungeon_action_id": dungeonActionID,
				"dungeon_object_id": dungeonActionRec.ResolvedTargetDungeonObjectID.String,
			}, nil, false)
		if err != nil {
			m.Log.Warn("failed getting target dungeon action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of dungeon action object records returned >%d<", len(dungeonActionObjectRecs))
			m.Log.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		dungeonActionRecordSet.TargetActionObjectRec = dungeonActionObjectRecs[0]
	}

	return &dungeonActionRecordSet, nil
}

// func (m *Model) GetDungeonLocationRecordSet(dungeonCharacterID string, forUpdate bool) (*DungeonLocationRecordSet, error) {
func (m *Model) GetDungeonLocationRecordSet(dungeonLocationID string, forUpdate bool) (*DungeonLocationRecordSet, error) {

	dungeonLocationRecordSet := &DungeonLocationRecordSet{}

	// Location record
	locationRec, err := m.GetDungeonLocationRec(dungeonLocationID, forUpdate)
	if err != nil {
		m.Log.Warn("failed to get dungeon location record >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.LocationRec = locationRec

	// All characters at location
	characterRecs, err := m.GetDungeonCharacterRecs(
		map[string]interface{}{
			"dungeon_location_id": locationRec.ID,
		},
		nil,
		forUpdate,
	)
	if err != nil {
		m.Log.Warn("failed to get dungeon location character records >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.CharacterRecs = characterRecs

	// All monsters at location
	monsterRecs, err := m.GetDungeonMonsterRecs(
		map[string]interface{}{
			"dungeon_location_id": locationRec.ID,
		},
		nil,
		forUpdate,
	)
	if err != nil {
		m.Log.Warn("failed to get dungeon location monster records >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.MonsterRecs = monsterRecs

	// All objects at location
	objectRecs, err := m.GetDungeonObjectRecs(
		map[string]interface{}{
			"dungeon_location_id": locationRec.ID,
		},
		nil,
		forUpdate,
	)
	if err != nil {
		m.Log.Warn("failed to get dungeon location object records >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.ObjectRecs = objectRecs

	locationIDs := []string{}
	if locationRec.NorthDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.NorthDungeonLocationID.String)
	}
	if locationRec.NortheastDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.NortheastDungeonLocationID.String)
	}
	if locationRec.EastDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.EastDungeonLocationID.String)
	}
	if locationRec.SoutheastDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.SoutheastDungeonLocationID.String)
	}
	if locationRec.SouthDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.SouthDungeonLocationID.String)
	}
	if locationRec.SouthwestDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.SouthwestDungeonLocationID.String)
	}
	if locationRec.WestDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.WestDungeonLocationID.String)
	}
	if locationRec.NorthwestDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.NorthwestDungeonLocationID.String)
	}
	if locationRec.UpDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.UpDungeonLocationID.String)
	}
	if locationRec.DownDungeonLocationID.Valid {
		locationIDs = append(locationIDs, locationRec.DownDungeonLocationID.String)
	}

	locationRecs, err := m.GetDungeonLocationRecs(
		map[string]interface{}{
			"id": locationIDs,
		},
		nil,
		forUpdate,
	)
	if err != nil {
		m.Log.Warn("failed to get dungeon location direction location records >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.LocationRecs = locationRecs

	return dungeonLocationRecordSet, nil
}
