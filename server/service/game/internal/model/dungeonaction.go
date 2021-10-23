package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

type DungeonActionRecordSet struct {
	DungeonActionRec           *record.DungeonAction
	DungeonActionCharacterRecs []*record.DungeonActionCharacter
	DungeonActionMonsterRecs   []*record.DungeonActionMonster
	DungeonActionObjectRecs    []*record.DungeonActionObject
}

type DungeonLocationRecordSet struct {
	CharacterRec  *record.DungeonCharacter
	LocationRec   *record.DungeonLocation
	CharacterRecs []*record.DungeonCharacter
	MonsterRecs   []*record.DungeonMonster
	ObjectRecs    []*record.DungeonObject
	LocationRecs  []*record.DungeonLocation
}

func (m *Model) ProcessDungeonCharacterAction(dungeonCharacterID string, sentence string) (*DungeonActionRecordSet, error) {

	m.Log.Info("Processing dungeon character ID >%s< action command >%s<", dungeonCharacterID, sentence)

	// Get current dungeon location record set
	dungeonLocationRecordSet, err := m.getDungeonLocationRecordSet(dungeonCharacterID, true)
	if err != nil {
		m.Log.Warn("Failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}

	// Resolve character action
	dungeonActionRec, err := m.resolveAction(sentence, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("Failed resolving dungeon character action >%v<", err)
		return nil, err
	}

	// Perform character action
	dungeonActionRec, err = m.performDungeonCharacterAction(dungeonActionRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("Failed performing dungeon character action >%v<", err)
		return nil, err
	}

	// Refetch current dungeon location record set
	dungeonLocationRecordSet, err = m.getDungeonLocationRecordSet(dungeonCharacterID, true)
	if err != nil {
		m.Log.Warn("Failed getting dungeon location record set after performing action >%v<", err)
		return nil, err
	}

	// Create dungeon action event records
	err = m.CreateDungeonActionRec(dungeonActionRec)
	if err != nil {
		m.Log.Warn("Failed creating dungeon action record >%v<", err)
		return nil, err
	}

	m.Log.Info("Created dungeon action record ID >%s<", dungeonActionRec.ID)

	dungeonActionRecordSet := DungeonActionRecordSet{
		DungeonActionRec:           dungeonActionRec,
		DungeonActionCharacterRecs: []*record.DungeonActionCharacter{},
		DungeonActionMonsterRecs:   []*record.DungeonActionMonster{},
		DungeonActionObjectRecs:    []*record.DungeonActionObject{},
	}

	if len(dungeonLocationRecordSet.CharacterRecs) > 0 {
		for _, characterRec := range dungeonLocationRecordSet.CharacterRecs {
			dungeonActionCharacterRec := record.DungeonActionCharacter{
				DungeonActionID:    dungeonActionRec.ID,
				DungeonCharacterID: characterRec.ID,
			}

			err := m.CreateDungeonActionCharacterRec(&dungeonActionCharacterRec)
			if err != nil {
				m.Log.Warn("Failed creating dungeon action character record >%v<", err)
				return nil, err
			}

			m.Log.Info("Created dungeon action character record ID >%s<", dungeonActionCharacterRec.ID)
			dungeonActionRecordSet.DungeonActionCharacterRecs = append(dungeonActionRecordSet.DungeonActionCharacterRecs, &dungeonActionCharacterRec)
		}
	}

	if len(dungeonLocationRecordSet.MonsterRecs) > 0 {
		for _, monsterRec := range dungeonLocationRecordSet.MonsterRecs {
			dungeonActionMonsterRec := record.DungeonActionMonster{
				DungeonActionID:  dungeonActionRec.ID,
				DungeonMonsterID: monsterRec.ID,
			}
			err := m.CreateDungeonActionMonsterRec(&dungeonActionMonsterRec)
			if err != nil {
				m.Log.Warn("Failed creating dungeon action monster record >%v<", err)
				return nil, err
			}
			m.Log.Info("Created dungeon action monster record ID >%s<", dungeonActionMonsterRec.ID)
			dungeonActionRecordSet.DungeonActionMonsterRecs = append(dungeonActionRecordSet.DungeonActionMonsterRecs, &dungeonActionMonsterRec)
		}
	}

	if len(dungeonLocationRecordSet.ObjectRecs) > 0 {
		for _, objectRec := range dungeonLocationRecordSet.ObjectRecs {
			dungeonActionObjectRec := record.DungeonActionObject{
				DungeonActionID: dungeonActionRec.ID,
				DungeonObjectID: objectRec.ID,
			}
			err := m.CreateDungeonActionObjectRec(&dungeonActionObjectRec)
			if err != nil {
				m.Log.Warn("Failed creating dungeon action object record >%v<", err)
				return nil, err
			}
			m.Log.Info("Created dungeon action object record ID >%s<", dungeonActionObjectRec.ID)
			dungeonActionRecordSet.DungeonActionObjectRecs = append(dungeonActionRecordSet.DungeonActionObjectRecs, &dungeonActionObjectRec)
		}
	}

	return &dungeonActionRecordSet, nil
}

func (m *Model) getDungeonLocationRecordSet(dungeonCharacterID string, forUpdate bool) (*DungeonLocationRecordSet, error) {

	dungeonLocationRecordSet := &DungeonLocationRecordSet{}

	// Character record
	characterRec, err := m.GetDungeonCharacterRec(dungeonCharacterID, forUpdate)
	if err != nil {
		m.Log.Warn("Failed to get dungeon character record >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.CharacterRec = characterRec

	// Location record
	locationRec, err := m.GetDungeonLocationRec(characterRec.DungeonLocationID, forUpdate)
	if err != nil {
		m.Log.Warn("Failed to get dungeon location record >%v<", err)
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
		m.Log.Warn("Failed to get dungeon location character records >%v<", err)
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
		m.Log.Warn("Failed to get dungeon location monster records >%v<", err)
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
		m.Log.Warn("Failed to get dungeon location object records >%v<", err)
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
		m.Log.Warn("Failed to get dungeon location direction location records >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.LocationRecs = locationRecs

	return dungeonLocationRecordSet, nil
}
