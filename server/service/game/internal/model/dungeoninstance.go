package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/core/nullstring"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

type DungeonInstanceRecordSet struct {
	DungeonInstanceRec   *record.DungeonInstance
	LocationInstanceRecs []*record.LocationInstance
	ObjectInstanceRecs   []*record.ObjectInstance
	MonsterInstanceRecs  []*record.MonsterInstance
}

// GetAvailableDungeonInstance returns an available dungeon instance
func (m *Model) GetAvailableDungeonInstance(dungeonID string) (*DungeonInstanceRecordSet, error) {
	l := m.Logger("GetAvailableDungeonInstance")

	l.Info("Finding available dungeon instance for dungeon ID >%s<", dungeonID)

	// Character capacity is equal to the number of locations a dungeon has.

	return nil, nil
}

// CreateDungeonInstance creates a dungeon, locations, monsters and objects instances
func (m *Model) CreateDungeonInstance(dungeonID string) (*DungeonInstanceRecordSet, error) {
	l := m.Logger("CreateDungeonInstance")

	l.Debug("Creating dungeon instance from dungeon ID >%s<", dungeonID)

	r := m.DungeonInstanceRepository()

	locationInstanceRecs := []*record.LocationInstance{}
	monsterInstanceRecs := []*record.MonsterInstance{}
	objectInstanceRecs := []*record.ObjectInstance{}

	dungeonInstanceRec := &record.DungeonInstance{
		DungeonID: dungeonID,
	}

	err := r.CreateOne(dungeonInstanceRec)
	if err != nil {
		l.Warn("Failed creating dungeon instance record >%v<", err)
		return nil, err
	}

	// Create location instance records
	locationRecs, err := m.GetLocationRecs(
		map[string]interface{}{
			"dungeon_id": dungeonID,
		},
		nil, false,
	)
	if err != nil {
		l.Warn("Failed getting locations records >%v<", err)
		return nil, err
	}

	for _, locationRec := range locationRecs {
		locationInstanceRec := &record.LocationInstance{
			DungeonInstanceID: dungeonInstanceRec.ID,
			LocationID:        locationRec.ID,
		}
		err := m.CreateLocationInstanceRec(locationInstanceRec)
		if err != nil {
			l.Warn("Failed creating location instance record >%v<", err)
			return nil, err
		}
		locationInstanceRecs = append(locationInstanceRecs, locationInstanceRec)
	}

	locationMap := makeLocationMap(locationRecs, locationInstanceRecs)

	// Resolve location instance direction IDs
	locationInstanceRecs, err = m.resolveLocationInstanceDirectionIdentifiers(locationMap, locationInstanceRecs)
	if err != nil {
		l.Warn("Failed resolving location instance direction identifiers >%v<", err)
		return nil, err
	}

	// Update location instance records
	for _, locationInstanceRec := range locationInstanceRecs {
		err := m.UpdateLocationInstanceRec(locationInstanceRec)
		if err != nil {
			l.Warn("Failed updating location instance record >%v<", err)
			return nil, err
		}

		// Create location object instance records
		locationObjectRecs, err := m.GetLocationObjectRecs(
			map[string]interface{}{
				"location_id": locationInstanceRec.LocationID,
			},
			nil, false,
		)
		if err != nil {
			l.Warn("Failed getting location monster records >%v<", err)
			return nil, err
		}

		objectInstanceMap := map[string]*record.ObjectInstance{}
		for _, locationObjectRec := range locationObjectRecs {
			objectRec, err := m.GetObjectRec(locationObjectRec.ObjectID, false)
			if err != nil {
				l.Warn("Failed getting object record >%v<", err)
				return nil, err
			}

			objectInstanceRec := &record.ObjectInstance{
				ObjectID:           locationObjectRec.ObjectID,
				DungeonInstanceID:  dungeonInstanceRec.ID,
				LocationInstanceID: nullstring.FromString(locationInstanceRec.ID),
			}

			err = m.CreateObjectInstanceRec(objectInstanceRec)
			if err != nil {
				l.Warn("Failed creating location object instance record >%v<", err)
				return nil, err
			}

			objectInstanceMap[objectRec.ID] = objectInstanceRec
			objectInstanceRecs = append(objectInstanceRecs, objectInstanceRec)
		}

		// Create location monster instance records
		locationMonsterRecs, err := m.GetLocationMonsterRecs(
			map[string]interface{}{
				"location_id": locationInstanceRec.LocationID,
			},
			nil, false,
		)
		if err != nil {
			l.Warn("Failed getting location monster records >%v<", err)
			return nil, err
		}

		monsterInstanceMap := map[string]*record.MonsterInstance{}
		for _, monsterLocationRec := range locationMonsterRecs {
			monsterRec, err := m.GetMonsterRec(monsterLocationRec.MonsterID, false)
			if err != nil {
				l.Warn("Failed getting monster record >%v<", err)
				return nil, err
			}

			monsterInstanceRec := &record.MonsterInstance{
				MonsterID:          monsterRec.ID,
				DungeonInstanceID:  dungeonInstanceRec.ID,
				LocationInstanceID: locationInstanceRec.ID,
				Strength:           monsterRec.Strength,
				Dexterity:          monsterRec.Dexterity,
				Intelligence:       monsterRec.Intelligence,
				Health:             monsterRec.Health,
				Fatigue:            monsterRec.Fatigue,
				Coins:              monsterRec.Coins,
				ExperiencePoints:   monsterRec.ExperiencePoints,
				AttributePoints:    monsterRec.AttributePoints,
			}

			err = m.CreateMonsterInstanceRec(monsterInstanceRec)
			if err != nil {
				l.Warn("Failed creating monster instance record >%v<", err)
				return nil, err
			}

			monsterInstanceMap[monsterRec.ID] = monsterInstanceRec
			monsterInstanceRecs = append(monsterInstanceRecs, monsterInstanceRec)

			monsterObjectRecs, err := m.GetMonsterObjectRecs(
				map[string]interface{}{
					"monster_id": monsterRec.ID,
				}, nil, false,
			)
			if err != nil {
				l.Warn("Failed getting monster object records >%v<", err)
				return nil, err
			}

			for _, monsterObjectRec := range monsterObjectRecs {

				objectInstanceRec := &record.ObjectInstance{
					ObjectID:          monsterObjectRec.ObjectID,
					DungeonInstanceID: dungeonInstanceRec.ID,
					MonsterInstanceID: nullstring.FromString(monsterInstanceRec.ID),
					IsEquipped:        monsterObjectRec.IsEquipped,
					IsStashed:         monsterObjectRec.IsStashed,
				}

				err := m.CreateObjectInstanceRec(objectInstanceRec)
				if err != nil {
					l.Warn("Failed creating monster object instance record >%v<", err)
					return nil, err
				}

				objectInstanceMap[monsterObjectRec.ObjectID] = objectInstanceRec
				objectInstanceRecs = append(objectInstanceRecs, objectInstanceRec)
			}
		}
	}

	dungeonInstanceRecordSet := DungeonInstanceRecordSet{
		DungeonInstanceRec:   dungeonInstanceRec,
		LocationInstanceRecs: locationInstanceRecs,
		MonsterInstanceRecs:  monsterInstanceRecs,
		ObjectInstanceRecs:   objectInstanceRecs,
	}

	return &dungeonInstanceRecordSet, nil
}

type LocationMapItem struct {
	LocationRec         *record.Location
	LocationInstanceRec *record.LocationInstance
}

func makeLocationMap(locationRecs []*record.Location, locationInstanceRecs []*record.LocationInstance) map[string]LocationMapItem {

	// Create a map of location records and location instance records
	// indexed by location record ID
	locationMap := map[string]LocationMapItem{}

	for _, locationRec := range locationRecs {
		locationMap[locationRec.ID] = LocationMapItem{
			LocationRec: locationRec,
		}
	}

	for _, locationInstanceRec := range locationInstanceRecs {
		locationMapItem := locationMap[locationInstanceRec.LocationID]
		locationMapItem.LocationInstanceRec = locationInstanceRec
		locationMap[locationInstanceRec.LocationID] = locationMapItem
	}

	return locationMap
}

// resolveLocationInstanceDirectionIdentifiers -
func (m *Model) resolveLocationInstanceDirectionIdentifiers(locationMap map[string]LocationMapItem, locationInstanceRecs []*record.LocationInstance) ([]*record.LocationInstance, error) {

	l := m.Logger("CreateDungeonInstance")

	l.Debug("Resolving location instance direction identifiers")

	for _, locationInstanceRec := range locationInstanceRecs {
		locationRec := locationMap[locationInstanceRec.LocationID].LocationRec
		if locationRec == nil {
			return nil, fmt.Errorf("missing location record with ID >%s<", locationInstanceRec.LocationID)
		}

		if nullstring.IsValid(locationRec.NorthLocationID) {
			locationInstanceRec.NorthLocationInstanceID = nullstring.FromString(locationMap[locationRec.NorthLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.NortheastLocationID) {
			locationInstanceRec.NortheastLocationInstanceID = nullstring.FromString(locationMap[locationRec.NortheastLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.EastLocationID) {
			locationInstanceRec.EastLocationInstanceID = nullstring.FromString(locationMap[locationRec.EastLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.SoutheastLocationID) {
			locationInstanceRec.SoutheastLocationInstanceID = nullstring.FromString(locationMap[locationRec.SoutheastLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.SouthLocationID) {
			locationInstanceRec.SouthLocationInstanceID = nullstring.FromString(locationMap[locationRec.SouthLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.SouthwestLocationID) {
			locationInstanceRec.SouthwestLocationInstanceID = nullstring.FromString(locationMap[locationRec.SouthwestLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.WestLocationID) {
			locationInstanceRec.WestLocationInstanceID = nullstring.FromString(locationMap[locationRec.WestLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.NorthwestLocationID) {
			locationInstanceRec.NorthwestLocationInstanceID = nullstring.FromString(locationMap[locationRec.NorthwestLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.UpLocationID) {
			locationInstanceRec.UpLocationInstanceID = nullstring.FromString(locationMap[locationRec.UpLocationID.String].LocationInstanceRec.ID)
		}
		if nullstring.IsValid(locationRec.DownLocationID) {
			locationInstanceRec.DownLocationInstanceID = nullstring.FromString(locationMap[locationRec.DownLocationID.String].LocationInstanceRec.ID)
		}
	}

	return locationInstanceRecs, nil
}

// // Resolve all location direction identifiers on all dungeon location instances
// data, err = t.resolveLocationInstanceDirectionIdentifiers(data, dungeonConfig)
// if err != nil {
// 	t.Log.Warn("Failed resolving data location instance identifiers >%v<", err)
// 	return err
// }

// // Update all previously created location instance records as they now have all their
// // reference location instance identifiers now set.
// for _, dungeonLocationInstanceRec := range data.LocationInstanceRecs {
// 	err := t.updateLocationInstanceRec(dungeonLocationInstanceRec)
// 	if err != nil {
// 		t.Log.Warn("Failed updating location instance record >%v<", err)
// 		return err
// 	}
// }
