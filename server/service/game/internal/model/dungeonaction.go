package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// TODO: Complete processing character action
func (m *Model) ProcessDungeonCharacterAction(dungeonCharacterID string, sentence string) (*record.DungeonAction, error) {

	m.Log.Info("Processing dungeon character ID >%s< action command >%s<", dungeonCharacterID, sentence)

	// Get current dungeon location record set

	// Resolve character action

	// Perform character action

	// Refetch current dungeon location record set

	// Create dungeon action event records

	return nil, nil
}

type DungeonLocationRecordSet struct {
	CharacterRec  *record.DungeonCharacter
	LocationRec   *record.DungeonLocation
	CharacterRecs []*record.DungeonCharacter
	MonsterRecs   []*record.DungeonMonster
	ObjectRecs    []*record.DungeonObject
	LocationRecs  []*record.DungeonLocation
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
			"dungeon_location_id": locationIDs,
		},
		nil,
		forUpdate,
	)
	if err != nil {
		m.Log.Warn("Failed to get dungeon location direction location records >%v<", err)
		return nil, err
	}
	dungeonLocationRecordSet.LocationRecs = locationRecs

	// [
	// 'north_dungeon_location_id',
	// 'northeast_dungeon_location_id',
	// 'east_dungeon_location_id',
	// 'southeast_dungeon_location_id',
	// 'south_dungeon_location_id',
	// 'southwest_dungeon_location_id',
	// 'west_dungeon_location_id',
	// 'northwest_dungeon_location_id',
	// 'up_dungeon_location_id',
	// 'down_dungeon_location_id',
	// ].forEach((prop) => {
	// if (locationRecord[prop]) {
	// locationIds.push(locationRecord[prop]);
	// }
	// });
	//
	// Location records
	// const locationRecords = await this.dungeonLocationRepository.getMany({
	// parameters: [
	// {
	// column: 'id',
	// value: locationIds,
	// operator: RepositoryOperator.In,
	// },
	// ],
	// });
	// logger.info(`Fetched ${locationRecords.length} dungeon location records`);
	//
	// Resolve action sentence
	// const records: DungeonLocationRecordSet = {
	// character: characterRecord,
	// location: locationRecord,
	// characters: characterRecords,
	// monsters: monsterRecords,
	// objects: objectRecords,
	// locations: locationRecords,
	// };

	return nil, nil
}
