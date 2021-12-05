package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

type DungeonActionRecordSet struct {
	DungeonActionRec    *record.DungeonAction
	DungeonCharacterRec *record.DungeonCharacter
	DungeonMonsterRec   *record.DungeonMonster
	CurrentLocation     *DungeonActionLocationRecordSet
	EquippedObjectRec   *record.DungeonObject
	StashedObjectRec    *record.DungeonObject
	TargetObjectRec     *record.DungeonObject
	TargetCharacterRec  *record.DungeonCharacter
	TargetMonsterRec    *record.DungeonMonster
	TargetLocation      *DungeonActionLocationRecordSet
}

type DungeonActionLocationRecordSet struct {
	DungeonLocationRec         *record.DungeonLocation
	DungeonActionCharacterRecs []record.DungeonActionCharacter
	DungeonCharacterRecs       []record.DungeonCharacter
	DungeonActionMonsterRecs   []record.DungeonActionMonster
	DungeonMonsterRecs         []record.DungeonMonster
	DungeonActionObjectRecs    []record.DungeonActionObject
	DungeonObjectRecs          []record.DungeonObject
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

	m.Log.Info("(testing) Processing dungeon character ID >%s< action command >%s<", dungeonCharacterID, sentence)

	// Verify the character performing the action exists within the specified dungeon
	dungeonCharacterRec, err := m.GetDungeonCharacterRec(dungeonCharacterID, true)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon character record before performing action >%v<", err)
		return nil, err
	}

	if dungeonCharacterRec.DungeonID != dungeonID {
		msg := fmt.Sprintf("Failed performing dungeon character action, character ID >%s< does not exist in dungeon ID >%s<", dungeonCharacterID, dungeonID)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Get the current dungeon location set of related records
	dungeonLocationRecordSet, err := m.GetDungeonLocationRecordSet(dungeonCharacterRec.DungeonLocationID, true)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}

	// Resolve the submitted character action
	dungeonActionRec, err := m.resolveAction(sentence, dungeonCharacterRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("(testing) Failed resolving dungeon character action >%v<", err)
		return nil, err
	}

	// TODO: Debug dungeon action record here..
	m.Log.Warn("(*******) Dungeon action record command >%s<", dungeonActionRec.ResolvedCommand)
	m.Log.Warn("(*******) Dungeon action record location >%s<", dungeonActionRec.DungeonLocationID)

	// Perform the submitted character action
	dungeonActionRec, err = m.performDungeonCharacterAction(dungeonCharacterRec, dungeonActionRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("(testing) Failed performing dungeon character action >%v<", err)
		return nil, err
	}

	m.Log.Warn("(*******) Dungeon action record command >%s<", dungeonActionRec.ResolvedCommand)
	m.Log.Warn("(*******) Dungeon action record location >%s<", dungeonActionRec.DungeonLocationID)

	// Create the resulting dungeon action event record
	err = m.CreateDungeonActionRec(dungeonActionRec)
	if err != nil {
		m.Log.Warn("(testing) Failed creating dungeon action record >%v<", err)
		return nil, err
	}

	m.Log.Info("(testing) Created dungeon action record ID >%s<", dungeonActionRec.ID)

	// TODO: Maybe don't need to do this... Get the updated character record
	dungeonCharacterRec, err = m.GetDungeonCharacterRec(dungeonCharacterID, false)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon character record after performing action >%v<", err)
		return nil, err
	}

	// TODO: And here..

	dungeonActionRecordSet := DungeonActionRecordSet{
		DungeonActionRec:    dungeonActionRec,
		DungeonCharacterRec: dungeonCharacterRec,
	}

	// Get the updated current dungeon location record set
	dungeonLocationRecordSet, err = m.GetDungeonLocationRecordSet(dungeonActionRec.DungeonLocationID, true)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon location record set after performing action >%v<", err)
		return nil, err
	}

	m.Log.Info("(****** testing) Dungeon location record set location name >%s<", dungeonLocationRecordSet.LocationRec.Name)
	m.Log.Info("(****** testing) Dungeon location record set characters >%d<", len(dungeonLocationRecordSet.CharacterRecs))
	m.Log.Info("(****** testing) Dungeon location record set monsters >%d<", len(dungeonLocationRecordSet.MonsterRecs))
	m.Log.Info("(****** testing) Dungeon location record set objects >%d<", len(dungeonLocationRecordSet.ObjectRecs))

	// Current location
	dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.DungeonLocationID, false)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon location record after performing action >%v<", err)
		return nil, err
	}

	currentLocationRecordSet := DungeonActionLocationRecordSet{
		DungeonLocationRec:         dungeonLocationRec,
		DungeonActionCharacterRecs: []record.DungeonActionCharacter{},
		DungeonCharacterRecs:       []record.DungeonCharacter{},
		DungeonActionMonsterRecs:   []record.DungeonActionMonster{},
		DungeonMonsterRecs:         []record.DungeonMonster{},
		DungeonActionObjectRecs:    []record.DungeonActionObject{},
		DungeonObjectRecs:          []record.DungeonObject{},
	}

	// Create the dungeon action character record for each character now at the current location
	if len(dungeonLocationRecordSet.CharacterRecs) > 0 {
		for _, characterRec := range dungeonLocationRecordSet.CharacterRecs {
			dungeonActionCharacterRec := record.DungeonActionCharacter{
				DungeonActionID:    dungeonActionRec.ID,
				DungeonLocationID:  dungeonLocationRec.ID,
				DungeonCharacterID: characterRec.ID,
			}

			err := m.CreateDungeonActionCharacterRec(&dungeonActionCharacterRec)
			if err != nil {
				m.Log.Warn("(testing) Failed creating dungeon action character record >%v<", err)
				return nil, err
			}

			m.Log.Info("(testing) Created dungeon action character record ID >%s<", dungeonActionCharacterRec.ID)
			currentLocationRecordSet.DungeonActionCharacterRecs = append(currentLocationRecordSet.DungeonActionCharacterRecs, dungeonActionCharacterRec)

			_dungeonCharacterRec, err := m.GetDungeonCharacterRec(dungeonActionCharacterRec.DungeonCharacterID, false)
			if err != nil {
				m.Log.Warn("(testing) Failed getting dungeon character record >%v<", err)
				return nil, err
			}

			currentLocationRecordSet.DungeonCharacterRecs = append(currentLocationRecordSet.DungeonCharacterRecs, *_dungeonCharacterRec)
		}
	}

	// Create the dungeon action monster record for each monster now at the current location
	if len(dungeonLocationRecordSet.MonsterRecs) > 0 {
		for _, monsterRec := range dungeonLocationRecordSet.MonsterRecs {
			dungeonActionMonsterRec := record.DungeonActionMonster{
				DungeonActionID:   dungeonActionRec.ID,
				DungeonLocationID: dungeonLocationRec.ID,
				DungeonMonsterID:  monsterRec.ID,
			}
			err := m.CreateDungeonActionMonsterRec(&dungeonActionMonsterRec)
			if err != nil {
				m.Log.Warn("(testing) Failed creating dungeon action monster record >%v<", err)
				return nil, err
			}

			m.Log.Info("(testing) Created dungeon action monster record ID >%s<", dungeonActionMonsterRec.ID)
			currentLocationRecordSet.DungeonActionMonsterRecs = append(currentLocationRecordSet.DungeonActionMonsterRecs, dungeonActionMonsterRec)

			_dungeonMonsterRec, err := m.GetDungeonMonsterRec(dungeonActionMonsterRec.DungeonMonsterID, false)
			if err != nil {
				m.Log.Warn("(testing) Failed getting dungeon monster record >%v<", err)
				return nil, err
			}

			currentLocationRecordSet.DungeonMonsterRecs = append(currentLocationRecordSet.DungeonMonsterRecs, *_dungeonMonsterRec)
		}
	}

	// Create the dungeon action object record for each object now at the current location
	if len(dungeonLocationRecordSet.ObjectRecs) > 0 {
		for _, objectRec := range dungeonLocationRecordSet.ObjectRecs {
			dungeonActionObjectRec := record.DungeonActionObject{
				DungeonActionID:   dungeonActionRec.ID,
				DungeonLocationID: dungeonLocationRec.ID,
				DungeonObjectID:   objectRec.ID,
			}
			err := m.CreateDungeonActionObjectRec(&dungeonActionObjectRec)
			if err != nil {
				m.Log.Warn("(testing) Failed creating dungeon action object record >%v<", err)
				return nil, err
			}

			m.Log.Info("(testing) Created dungeon action object record ID >%s<", dungeonActionObjectRec.ID)
			currentLocationRecordSet.DungeonActionObjectRecs = append(currentLocationRecordSet.DungeonActionObjectRecs, dungeonActionObjectRec)

			_dungeonObjectRec, err := m.GetDungeonObjectRec(dungeonActionObjectRec.DungeonObjectID, false)
			if err != nil {
				m.Log.Warn("(testing) Failed getting dungeon object record >%v<", err)
				return nil, err
			}

			currentLocationRecordSet.DungeonObjectRecs = append(currentLocationRecordSet.DungeonObjectRecs, *_dungeonObjectRec)
		}
	}

	dungeonActionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target dungeon location record set when set
	if dungeonActionRec.ResolvedTargetDungeonLocationID.Valid {

		dungeonLocationRecordSet, err := m.GetDungeonLocationRecordSet(dungeonActionRec.ResolvedTargetDungeonLocationID.String, true)
		if err != nil {
			m.Log.Warn("(testing) Failed getting target dungeon location record set after performing action >%v<", err)
			return nil, err
		}

		// Target location
		dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.ResolvedTargetDungeonLocationID.String, false)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon location record after performing action >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := DungeonActionLocationRecordSet{
			DungeonLocationRec:         dungeonLocationRec,
			DungeonActionCharacterRecs: []record.DungeonActionCharacter{},
			DungeonCharacterRecs:       []record.DungeonCharacter{},
			DungeonActionMonsterRecs:   []record.DungeonActionMonster{},
			DungeonMonsterRecs:         []record.DungeonMonster{},
			DungeonActionObjectRecs:    []record.DungeonActionObject{},
			DungeonObjectRecs:          []record.DungeonObject{},
		}

		// Create the dungeon action character record for each character at the target location
		if len(dungeonLocationRecordSet.CharacterRecs) > 0 {
			for _, characterRec := range dungeonLocationRecordSet.CharacterRecs {
				dungeonActionCharacterRec := record.DungeonActionCharacter{
					DungeonActionID:    dungeonActionRec.ID,
					DungeonLocationID:  dungeonLocationRec.ID,
					DungeonCharacterID: characterRec.ID,
				}

				err := m.CreateDungeonActionCharacterRec(&dungeonActionCharacterRec)
				if err != nil {
					m.Log.Warn("(testing) Failed creating dungeon action character record >%v<", err)
					return nil, err
				}

				m.Log.Info("(testing) Created dungeon action character record ID >%s<", dungeonActionCharacterRec.ID)
				targetLocationRecordSet.DungeonActionCharacterRecs = append(targetLocationRecordSet.DungeonActionCharacterRecs, dungeonActionCharacterRec)

				_dungeonCharacterRec, err := m.GetDungeonCharacterRec(dungeonActionCharacterRec.DungeonCharacterID, false)
				if err != nil {
					m.Log.Warn("(testing) Failed getting dungeon character record >%v<", err)
					return nil, err
				}

				targetLocationRecordSet.DungeonCharacterRecs = append(targetLocationRecordSet.DungeonCharacterRecs, *_dungeonCharacterRec)
			}
		}

		// Create the dungeon action monster record for each monster at the target location
		if len(dungeonLocationRecordSet.MonsterRecs) > 0 {
			for _, monsterRec := range dungeonLocationRecordSet.MonsterRecs {
				dungeonActionMonsterRec := record.DungeonActionMonster{
					DungeonActionID:   dungeonActionRec.ID,
					DungeonLocationID: dungeonLocationRec.ID,
					DungeonMonsterID:  monsterRec.ID,
				}
				err := m.CreateDungeonActionMonsterRec(&dungeonActionMonsterRec)
				if err != nil {
					m.Log.Warn("(testing) Failed creating dungeon action monster record >%v<", err)
					return nil, err
				}

				m.Log.Info("(testing) Created dungeon action monster record ID >%s<", dungeonActionMonsterRec.ID)
				targetLocationRecordSet.DungeonActionMonsterRecs = append(targetLocationRecordSet.DungeonActionMonsterRecs, dungeonActionMonsterRec)

				_dungeonMonsterRec, err := m.GetDungeonMonsterRec(dungeonActionMonsterRec.DungeonMonsterID, false)
				if err != nil {
					m.Log.Warn("(testing) Failed getting dungeon monster record >%v<", err)
					return nil, err
				}

				targetLocationRecordSet.DungeonMonsterRecs = append(targetLocationRecordSet.DungeonMonsterRecs, *_dungeonMonsterRec)
			}
		}

		// Create the dungeon action object record for each object at the target location
		if len(dungeonLocationRecordSet.ObjectRecs) > 0 {
			for _, objectRec := range dungeonLocationRecordSet.ObjectRecs {
				dungeonActionObjectRec := record.DungeonActionObject{
					DungeonActionID:   dungeonActionRec.ID,
					DungeonLocationID: dungeonLocationRec.ID,
					DungeonObjectID:   objectRec.ID,
				}
				err := m.CreateDungeonActionObjectRec(&dungeonActionObjectRec)
				if err != nil {
					m.Log.Warn("(testing) Failed creating dungeon action object record >%v<", err)
					return nil, err
				}

				m.Log.Info("(testing) Created dungeon action object record ID >%s<", dungeonActionObjectRec.ID)
				targetLocationRecordSet.DungeonActionObjectRecs = append(targetLocationRecordSet.DungeonActionObjectRecs, dungeonActionObjectRec)

				_dungeonObjectRec, err := m.GetDungeonObjectRec(dungeonActionObjectRec.DungeonObjectID, false)
				if err != nil {
					m.Log.Warn("(testing) Failed getting dungeon object record >%v<", err)
					return nil, err
				}

				targetLocationRecordSet.DungeonObjectRecs = append(targetLocationRecordSet.DungeonObjectRecs, *_dungeonObjectRec)
			}
		}

		dungeonActionRecordSet.TargetLocation = &targetLocationRecordSet
	}

	return &dungeonActionRecordSet, nil
}

func (m *Model) GetDungeonActionRecordSet(dungeonActionID string) (*DungeonActionRecordSet, error) {

	dungeonActionRecordSet := DungeonActionRecordSet{}

	dungeonActionRec, err := m.GetDungeonActionRec(dungeonActionID, false)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon action record >%v<", err)
		return nil, err
	}
	dungeonActionRecordSet.DungeonActionRec = dungeonActionRec

	// Add the dungeon action character record that performed the dungeon action.
	if dungeonActionRec.DungeonCharacterID.Valid {
		dungeonCharacterRec, err := m.GetDungeonCharacterRec(dungeonActionRec.DungeonCharacterID.String, false)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon character record >%v<", err)
			return nil, err
		}
		dungeonActionRecordSet.DungeonCharacterRec = dungeonCharacterRec
	}

	// Add the dungeon action monster record that performed the dungeon action.
	if dungeonActionRec.DungeonMonsterID.Valid {
		dungeonMonsterRec, err := m.GetDungeonMonsterRec(dungeonActionRec.DungeonMonsterID.String, false)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon monster record >%v<", err)
			return nil, err
		}
		dungeonActionRecordSet.DungeonMonsterRec = dungeonMonsterRec
	}

	// Add the current location record set where the action was performed.
	dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.DungeonLocationID, false)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon location record >%v<", err)
		return nil, err
	}

	currentLocationRecordSet := DungeonActionLocationRecordSet{
		DungeonLocationRec:         dungeonLocationRec,
		DungeonActionCharacterRecs: []record.DungeonActionCharacter{},
		DungeonCharacterRecs:       []record.DungeonCharacter{},
		DungeonActionMonsterRecs:   []record.DungeonActionMonster{},
		DungeonMonsterRecs:         []record.DungeonMonster{},
		DungeonActionObjectRecs:    []record.DungeonActionObject{},
		DungeonObjectRecs:          []record.DungeonObject{},
	}

	// Add the dungeon action character records that existed at the action location
	// at the time the action was performed.
	dungeonActionCharacterRecs, err := m.GetDungeonActionCharacterRecs(
		map[string]interface{}{
			"dungeon_action_id": dungeonActionID,
		},
		nil,
		false,
	)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon action character records >%v<", err)
		return nil, err
	}

	currentLocationRecordSet.DungeonActionCharacterRecs = []record.DungeonActionCharacter{}
	currentLocationRecordSet.DungeonCharacterRecs = []record.DungeonCharacter{}

	for _, dungeonActionCharacterRec := range dungeonActionCharacterRecs {
		currentLocationRecordSet.DungeonActionCharacterRecs = append(currentLocationRecordSet.DungeonActionCharacterRecs, *dungeonActionCharacterRec)
		dungeonCharacterRec, err := m.GetDungeonCharacterRec(dungeonActionCharacterRec.DungeonCharacterID, false)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon character record >%v<", err)
			return nil, err
		}
		currentLocationRecordSet.DungeonCharacterRecs = append(currentLocationRecordSet.DungeonCharacterRecs, *dungeonCharacterRec)
	}

	// Add the dungeon action monster records that existed at the action location
	// at the time the action was performed.
	dungeonActionMonsterRecs, err := m.GetDungeonActionMonsterRecs(
		map[string]interface{}{
			"dungeon_action_id": dungeonActionID,
		},
		nil,
		false,
	)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon action monster records >%v<", err)
		return nil, err
	}

	currentLocationRecordSet.DungeonActionMonsterRecs = []record.DungeonActionMonster{}
	currentLocationRecordSet.DungeonMonsterRecs = []record.DungeonMonster{}

	for _, dungeonActionMonsterRec := range dungeonActionMonsterRecs {
		currentLocationRecordSet.DungeonActionMonsterRecs = append(currentLocationRecordSet.DungeonActionMonsterRecs, *dungeonActionMonsterRec)
		dungeonMonsterRec, err := m.GetDungeonMonsterRec(dungeonActionMonsterRec.DungeonMonsterID, false)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon monster record >%v<", err)
			return nil, err
		}
		currentLocationRecordSet.DungeonMonsterRecs = append(currentLocationRecordSet.DungeonMonsterRecs, *dungeonMonsterRec)
	}

	// Add the dungeon action object records that existed at the action location
	// at the time the action was performed.
	dungeonActionObjectRecs, err := m.GetDungeonActionObjectRecs(
		map[string]interface{}{
			"dungeon_action_id": dungeonActionID,
		},
		nil,
		false,
	)
	if err != nil {
		m.Log.Warn("(testing) Failed getting dungeon action monster records >%v<", err)
		return nil, err
	}

	currentLocationRecordSet.DungeonActionObjectRecs = []record.DungeonActionObject{}
	currentLocationRecordSet.DungeonObjectRecs = []record.DungeonObject{}

	for _, dungeonActionObjectRec := range dungeonActionObjectRecs {
		currentLocationRecordSet.DungeonActionObjectRecs = append(currentLocationRecordSet.DungeonActionObjectRecs, *dungeonActionObjectRec)
		dungeonObjectRec, err := m.GetDungeonObjectRec(dungeonActionObjectRec.DungeonObjectID, false)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon monster record >%v<", err)
			return nil, err
		}
		currentLocationRecordSet.DungeonObjectRecs = append(currentLocationRecordSet.DungeonObjectRecs, *dungeonObjectRec)
	}

	dungeonActionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target dungeon location record set when set
	if dungeonActionRec.ResolvedTargetDungeonLocationID.Valid {

		// Target location
		dungeonLocationRec, err := m.GetDungeonLocationRec(dungeonActionRec.ResolvedTargetDungeonLocationID.String, false)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon location record after performing action >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := DungeonActionLocationRecordSet{
			DungeonLocationRec:         dungeonLocationRec,
			DungeonActionCharacterRecs: []record.DungeonActionCharacter{},
			DungeonCharacterRecs:       []record.DungeonCharacter{},
			DungeonActionMonsterRecs:   []record.DungeonActionMonster{},
			DungeonMonsterRecs:         []record.DungeonMonster{},
			DungeonActionObjectRecs:    []record.DungeonActionObject{},
			DungeonObjectRecs:          []record.DungeonObject{},
		}

		// Add the dungeon action character records that existed at the action location
		// at the time the action was performed.
		dungeonActionCharacterRecs, err := m.GetDungeonActionCharacterRecs(
			map[string]interface{}{
				"dungeon_action_id": dungeonActionID,
			},
			nil,
			false,
		)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon action character records >%v<", err)
			return nil, err
		}

		targetLocationRecordSet.DungeonActionCharacterRecs = []record.DungeonActionCharacter{}
		targetLocationRecordSet.DungeonCharacterRecs = []record.DungeonCharacter{}

		for _, dungeonActionCharacterRec := range dungeonActionCharacterRecs {
			targetLocationRecordSet.DungeonActionCharacterRecs = append(targetLocationRecordSet.DungeonActionCharacterRecs, *dungeonActionCharacterRec)
			dungeonCharacterRec, err := m.GetDungeonCharacterRec(dungeonActionCharacterRec.DungeonCharacterID, false)
			if err != nil {
				m.Log.Warn("(testing) Failed getting dungeon character record >%v<", err)
				return nil, err
			}
			targetLocationRecordSet.DungeonCharacterRecs = append(targetLocationRecordSet.DungeonCharacterRecs, *dungeonCharacterRec)
		}

		// Add the dungeon action monster records that existed at the action location
		// at the time the action was performed.
		dungeonActionMonsterRecs, err := m.GetDungeonActionMonsterRecs(
			map[string]interface{}{
				"dungeon_action_id": dungeonActionID,
			},
			nil,
			false,
		)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon action monster records >%v<", err)
			return nil, err
		}

		targetLocationRecordSet.DungeonActionMonsterRecs = []record.DungeonActionMonster{}
		targetLocationRecordSet.DungeonMonsterRecs = []record.DungeonMonster{}

		for _, dungeonActionMonsterRec := range dungeonActionMonsterRecs {
			targetLocationRecordSet.DungeonActionMonsterRecs = append(targetLocationRecordSet.DungeonActionMonsterRecs, *dungeonActionMonsterRec)
			dungeonMonsterRec, err := m.GetDungeonMonsterRec(dungeonActionMonsterRec.DungeonMonsterID, false)
			if err != nil {
				m.Log.Warn("(testing) Failed getting dungeon monster record >%v<", err)
				return nil, err
			}
			targetLocationRecordSet.DungeonMonsterRecs = append(targetLocationRecordSet.DungeonMonsterRecs, *dungeonMonsterRec)
		}

		// Add the dungeon action object records that existed at the action location
		// at the time the action was performed.
		dungeonActionObjectRecs, err := m.GetDungeonActionObjectRecs(
			map[string]interface{}{
				"dungeon_action_id": dungeonActionID,
			},
			nil,
			false,
		)
		if err != nil {
			m.Log.Warn("(testing) Failed getting dungeon action monster records >%v<", err)
			return nil, err
		}

		targetLocationRecordSet.DungeonActionObjectRecs = []record.DungeonActionObject{}
		targetLocationRecordSet.DungeonObjectRecs = []record.DungeonObject{}

		for _, dungeonActionObjectRec := range dungeonActionObjectRecs {
			targetLocationRecordSet.DungeonActionObjectRecs = append(targetLocationRecordSet.DungeonActionObjectRecs, *dungeonActionObjectRec)
			dungeonObjectRec, err := m.GetDungeonObjectRec(dungeonActionObjectRec.DungeonObjectID, false)
			if err != nil {
				m.Log.Warn("(testing) Failed getting dungeon monster record >%v<", err)
				return nil, err
			}
			targetLocationRecordSet.DungeonObjectRecs = append(targetLocationRecordSet.DungeonObjectRecs, *dungeonObjectRec)
		}

		dungeonActionRecordSet.CurrentLocation = &currentLocationRecordSet
	}

	return &dungeonActionRecordSet, nil
}

// func (m *Model) GetDungeonLocationRecordSet(dungeonCharacterID string, forUpdate bool) (*DungeonLocationRecordSet, error) {
func (m *Model) GetDungeonLocationRecordSet(dungeonLocationID string, forUpdate bool) (*DungeonLocationRecordSet, error) {

	dungeonLocationRecordSet := &DungeonLocationRecordSet{}

	// // Character record
	// characterRec, err := m.GetDungeonCharacterRec(dungeonCharacterID, forUpdate)
	// if err != nil {
	// 	m.Log.Warn("(testing) Failed to get dungeon character record >%v<", err)
	// 	return nil, err
	// }
	// dungeonLocationRecordSet.CharacterRec = characterRec

	// Location record
	locationRec, err := m.GetDungeonLocationRec(dungeonLocationID, forUpdate)
	if err != nil {
		m.Log.Warn("(testing) Failed to get dungeon location record >%v<", err)
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
		m.Log.Warn("(testing) Failed to get dungeon location character records >%v<", err)
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
		m.Log.Warn("(testing) Failed to get dungeon location monster records >%v<", err)
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
		m.Log.Warn("(testing) Failed to get dungeon location object records >%v<", err)
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
		m.Log.Warn("(testing) Failed to get dungeon location direction location records >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.LocationRecs = locationRecs

	return dungeonLocationRecordSet, nil
}
