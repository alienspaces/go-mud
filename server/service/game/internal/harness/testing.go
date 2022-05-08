package harness

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/core/harness"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data         Data
	DataConfig   DataConfig
	teardownData teardownData
}

// Data -
type Data struct {
	// Character
	CharacterRecs []*record.Character

	// Dungeon
	DungeonRecs  []*record.Dungeon
	LocationRecs []*record.Location
	MonsterRecs  []*record.Monster
	ObjectRecs   []*record.Object

	// Dungeon Instance
	DungeonInstanceRecs   []*record.DungeonInstance
	LocationInstanceRecs  []*record.LocationInstance
	CharacterInstanceRecs []*record.CharacterInstance
	MonsterInstanceRecs   []*record.MonsterInstance
	ObjectInstanceRecs    []*record.ObjectInstance

	// Action
	ActionRecs                []*record.Action
	ActionCharacterRecs       []*record.ActionCharacter
	ActionCharacterObjectRecs []*record.ActionCharacterObject
	ActionMonsterRecs         []*record.ActionMonster
	ActionMonsterObjectRecs   []*record.ActionMonsterObject
	ActionObjectRecs          []*record.ActionObject
}

// teardownData -
type teardownData struct {
	// Character
	CharacterRecs []record.Character

	// Dungeon
	DungeonRecs  []record.Dungeon
	LocationRecs []record.Location
	MonsterRecs  []record.Monster
	ObjectRecs   []record.Object

	// Dungeon Instance
	DungeonInstanceRecs   []*record.DungeonInstance
	LocationInstanceRecs  []*record.LocationInstance
	CharacterInstanceRecs []*record.CharacterInstance
	MonsterInstanceRecs   []*record.MonsterInstance
	ObjectInstanceRecs    []*record.ObjectInstance

	// Action
	ActionRecs                []*record.Action
	ActionCharacterRecs       []*record.ActionCharacter
	ActionCharacterObjectRecs []*record.ActionCharacterObject
	ActionMonsterRecs         []*record.ActionMonster
	ActionMonsterObjectRecs   []*record.ActionMonsterObject
	ActionObjectRecs          []*record.ActionObject
}

// NewTesting -
func NewTesting(c configurer.Configurer, l logger.Logger, s storer.Storer, m modeller.Modeller, config DataConfig) (t *Testing, err error) {

	t = &Testing{
		Testing: harness.Testing{
			Config: c,
			Log:    l,
			Store:  s,
			Model:  m,
		},
	}

	// modeller
	t.ModellerFunc = t.Modeller

	// data
	t.CreateDataFunc = t.CreateData
	t.RemoveDataFunc = t.RemoveData

	t.DataConfig = config
	t.Data = Data{}
	t.teardownData = teardownData{}

	err = t.Init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Modeller -
func (t *Testing) Modeller() (modeller.Modeller, error) {
	return t.Model, nil
}

// CreateData - Custom data
func (t *Testing) CreateData() error {

	data := Data{}
	teardownData := teardownData{}

	// Characters
	for _, characterConfig := range t.DataConfig.CharacterConfig {
		characterRec, err := t.createCharacterRec(characterConfig)
		if err != nil {
			t.Log.Warn("Failed creating character record >%v<", err)
			return err
		}
		data.CharacterRecs = append(data.CharacterRecs, characterRec)
		teardownData.CharacterRecs = append(teardownData.CharacterRecs, *characterRec)
	}

	// Dungeons
	for _, dungeonConfig := range t.DataConfig.DungeonConfig {

		// Create the dungeon record
		dungeonRec, err := t.createDungeonRec(dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed creating dungeon record >%v<", err)
			return err
		}
		data.DungeonRecs = append(data.DungeonRecs, dungeonRec)
		teardownData.DungeonRecs = append(teardownData.DungeonRecs, *dungeonRec)

		// Assign the created dungeon record identifier to child configuration records
		dungeonConfig, err = t.resolveConfigDungeonIdentifiers(dungeonRec, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving dungeon config dungeon identifiers >%v<", err)
			return err
		}

		// Create the dungeons location records
		for _, dungeonLocationConfig := range dungeonConfig.LocationConfig {
			dungeonLocationRec, err := t.createLocationRec(dungeonRec, dungeonLocationConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon location record >%v<", err)
				return err
			}
			data.LocationRecs = append(data.LocationRecs, dungeonLocationRec)
			teardownData.LocationRecs = append(teardownData.LocationRecs, *dungeonLocationRec)
		}

		// Resolve all location direction identifiers on all dungeon locations
		data, err = t.resolveDataLocationDirectionIdentifiers(data, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving config location identifiers >%v<", err)
			return err
		}

		// Update all previously created location records as they now have all their
		// reference location identifiers set correctly.
		for _, dungeonLocationRec := range data.LocationRecs {
			err := t.updateLocationRec(dungeonLocationRec)
			if err != nil {
				t.Log.Warn("Failed updating dungeon location record >%v<", err)
				return err
			}
		}

		// Resolve monster and object config locations
		dungeonConfig, err = t.resolveConfigLocationIdentifiers(data, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving config location identifiers >%v<", err)
			return err
		}

		// Create monster records
		for _, dungeonMonsterConfig := range dungeonConfig.MonsterConfig {
			dungeonMonsterRec, err := t.createMonsterRec(dungeonRec, dungeonMonsterConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon monster record >%v<", err)
				return err
			}
			data.MonsterRecs = append(data.MonsterRecs, dungeonMonsterRec)
			teardownData.MonsterRecs = append(teardownData.MonsterRecs, *dungeonMonsterRec)
		}

		// Resolve object config character identifiers
		dungeonConfig, err = t.resolveConfigObjectCharacterIdentifiers(data, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving config object character identifiers >%v<", err)
			return err
		}

		// Resolve object config monster identifiers
		dungeonConfig, err = t.resolveConfigObjectMonsterIdentifiers(data, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving config object monster identifiers >%v<", err)
			return err
		}

		// Create object records
		for _, dungeonObjectConfig := range dungeonConfig.ObjectConfig {
			t.Log.Info("Creating dungeon object >%#v<", dungeonObjectConfig)
			dungeonObjectRec, err := t.createObjectRec(dungeonRec, dungeonObjectConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon object record >%v<", err)
				return err
			}
			data.ObjectRecs = append(data.ObjectRecs, dungeonObjectRec)
			teardownData.ObjectRecs = append(teardownData.ObjectRecs, *dungeonObjectRec)
		}

		// Create dungeon instances
		for _, dungeonInstanceConfig := range dungeonConfig.DungeonInstanceConfig {

			// Create dungeon instance record
			t.Log.Info("Creating dungeon instance >%#v<", dungeonInstanceConfig)
			dungeonInstanceRec, err := t.createDungeonInstanceRec(dungeonRec.ID)
			if err != nil {
				t.Log.Warn("Failed creating dungeon instance record >%v<", err)
				return err
			}
			data.DungeonInstanceRecs = append(data.DungeonInstanceRecs, dungeonInstanceRec)
			teardownData.DungeonInstanceRecs = append(teardownData.DungeonInstanceRecs, dungeonInstanceRec)

			// Get resulting location instance records
			locationInstanceRecs, err := t.getLocationInstanceRecs(dungeonInstanceRec.ID)
			if err != nil {
				t.Log.Warn("Failed getting location instance records >%v<", err)
				return err
			}
			data.LocationInstanceRecs = append(data.LocationInstanceRecs, locationInstanceRecs...)
			teardownData.LocationInstanceRecs = append(teardownData.LocationInstanceRecs, locationInstanceRecs...)

			// Get resulting monster instance records
			monsterInstanceRecs, err := t.getMonsterInstanceRecs(dungeonInstanceRec.ID)
			if err != nil {
				t.Log.Warn("Failed getting monster instance records >%v<", err)
				return err
			}
			data.MonsterInstanceRecs = append(data.MonsterInstanceRecs, monsterInstanceRecs...)
			teardownData.MonsterInstanceRecs = append(teardownData.MonsterInstanceRecs, monsterInstanceRecs...)

			// Get resulting object instance records
			objectInstanceRecs, err := t.getObjectInstanceRecs(dungeonInstanceRec.ID)
			if err != nil {
				t.Log.Warn("Failed getting object instance records >%v<", err)
				return err
			}
			data.ObjectInstanceRecs = append(data.ObjectInstanceRecs, objectInstanceRecs...)
			teardownData.ObjectInstanceRecs = append(teardownData.ObjectInstanceRecs, objectInstanceRecs...)

			// TODO: Create character instance records

			// TODO: Create action records

		}
	}

	// // Create action records
	// for _, dungeonActionConfig := range dungeonConfig.ActionConfig {
	// 	dungeonID := ""
	// 	dungeonCharacterID := ""
	// 	for _, characterRecord := range data.CharacterRecs {
	// 		if characterRecord.Name == dungeonActionConfig.CharacterName {
	// 			dungeonID = characterRecord.DungeonID
	// 			dungeonCharacterID = characterRecord.ID
	// 		}
	// 	}
	// 	if dungeonID == "" || dungeonCharacterID == "" {
	// 		msg := fmt.Sprintf("Failed to find dungeon character record name >%s<", dungeonActionConfig.CharacterName)
	// 		t.Log.Error(msg)
	// 		return fmt.Errorf(msg)
	// 	}

	// 	actionRecordSet, err := t.createDungeonCharacterActionRec(dungeonID, dungeonCharacterID, dungeonActionConfig.Command)
	// 	if err != nil {
	// 		t.Log.Warn("Failed creating dungeon action record >%v<", err)
	// 		return err
	// 	}

	// 	data.ActionRecs = append(data.ActionRecs, actionRecordSet.ActionRec)
	// 	teardownData.ActionRecs = append(teardownData.ActionRecs, actionRecordSet.ActionRec)

	// 	// Source
	// 	if actionRecordSet.ActionCharacterRec != nil {
	// 		data.ActionCharacterRecs = append(data.ActionCharacterRecs, actionRecordSet.ActionCharacterRec)
	// 		teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, actionRecordSet.ActionCharacterRec)

	// 		data.ActionCharacterObjectRecs = append(data.ActionCharacterObjectRecs, actionRecordSet.ActionCharacterObjectRecs...)
	// 		teardownData.ActionCharacterObjectRecs = append(teardownData.ActionCharacterObjectRecs, actionRecordSet.ActionCharacterObjectRecs...)
	// 	}
	// 	if actionRecordSet.ActionMonsterRec != nil {
	// 		data.ActionMonsterRecs = append(data.ActionMonsterRecs, actionRecordSet.ActionMonsterRec)
	// 		teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, actionRecordSet.ActionMonsterRec)

	// 		data.ActionMonsterObjectRecs = append(data.ActionMonsterObjectRecs, actionRecordSet.ActionMonsterObjectRecs...)
	// 		teardownData.ActionMonsterObjectRecs = append(teardownData.ActionMonsterObjectRecs, actionRecordSet.ActionMonsterObjectRecs...)
	// 	}

	// 	// Current location
	// 	t.Log.Info("Dungeon action record set current location >%#v<", actionRecordSet.CurrentLocation)
	// 	if actionRecordSet.CurrentLocation != nil {
	// 		dungeonActionLocationRecordSet := actionRecordSet.CurrentLocation
	// 		data.ActionCharacterRecs = append(data.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
	// 		data.ActionMonsterRecs = append(data.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
	// 		data.ActionObjectRecs = append(data.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)

	// 		teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
	// 		teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
	// 		teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	// 	}

	// 	// Target location
	// 	t.Log.Info("Dungeon action record set target location >%#v<", actionRecordSet.TargetLocation)
	// 	if actionRecordSet.TargetLocation != nil {
	// 		dungeonActionLocationRecordSet := actionRecordSet.TargetLocation
	// 		data.ActionCharacterRecs = append(data.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
	// 		data.ActionMonsterRecs = append(data.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
	// 		data.ActionObjectRecs = append(data.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)

	// 		teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
	// 		teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
	// 		teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	// 	}

	// 	// Targets
	// 	if actionRecordSet.TargetActionCharacterRec != nil {
	// 		data.ActionCharacterRecs = append(data.ActionCharacterRecs, actionRecordSet.TargetActionCharacterRec)
	// 		teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, actionRecordSet.TargetActionCharacterRec)

	// 		data.ActionCharacterObjectRecs = append(data.ActionCharacterObjectRecs, actionRecordSet.TargetActionCharacterObjectRecs...)
	// 		teardownData.ActionCharacterObjectRecs = append(teardownData.ActionCharacterObjectRecs, actionRecordSet.TargetActionCharacterObjectRecs...)
	// 	}
	// 	if actionRecordSet.TargetActionMonsterRec != nil {
	// 		data.ActionMonsterRecs = append(data.ActionMonsterRecs, actionRecordSet.TargetActionMonsterRec)
	// 		teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, actionRecordSet.TargetActionMonsterRec)

	// 		data.ActionMonsterObjectRecs = append(data.ActionMonsterObjectRecs, actionRecordSet.TargetActionMonsterObjectRecs...)
	// 		teardownData.ActionMonsterObjectRecs = append(teardownData.ActionMonsterObjectRecs, actionRecordSet.TargetActionMonsterObjectRecs...)
	// 	}
	// 	if actionRecordSet.EquippedActionObjectRec != nil {
	// 		data.ActionObjectRecs = append(data.ActionObjectRecs, actionRecordSet.EquippedActionObjectRec)
	// 		teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, actionRecordSet.EquippedActionObjectRec)
	// 	}
	// 	if actionRecordSet.StashedActionObjectRec != nil {
	// 		data.ActionObjectRecs = append(data.ActionObjectRecs, actionRecordSet.StashedActionObjectRec)
	// 		teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, actionRecordSet.StashedActionObjectRec)
	// 	}
	// 	if actionRecordSet.TargetActionObjectRec != nil {
	// 		data.ActionObjectRecs = append(data.ActionObjectRecs, actionRecordSet.TargetActionObjectRec)
	// 		teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, actionRecordSet.TargetActionObjectRec)
	// 	}
	// }

	// Assign data once we have successfully set up ll data
	t.Data = data
	t.teardownData = teardownData

	return nil
}

func (t *Testing) resolveDataLocationDirectionIdentifiers(data Data, dungeonConfig DungeonConfig) (Data, error) {

	findLocationRec := func(locationName string) *record.Location {
		for _, dungeonLocationRec := range data.LocationRecs {
			if dungeonLocationRec.Name == locationName {
				return dungeonLocationRec
			}
		}
		return nil
	}

	if dungeonConfig.LocationConfig != nil {
		for _, dungeonLocationConfig := range dungeonConfig.LocationConfig {
			dungeonLocationRec := findLocationRec(dungeonLocationConfig.Record.Name)
			if dungeonLocationRec == nil {
				msg := fmt.Sprintf("Failed to find dungeon location record name >%s<", dungeonLocationConfig.Record.Name)
				t.Log.Error(msg)
				return data, fmt.Errorf(msg)
			}

			if dungeonLocationConfig.NorthLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.NorthLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find north dungeon location record name >%s<", dungeonLocationConfig.NorthLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.NorthLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.NortheastLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.NortheastLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find north east dungeon location record name >%s<", dungeonLocationConfig.NortheastLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.NortheastLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.EastLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.EastLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find east dungeon location record name >%s<", dungeonLocationConfig.EastLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.EastLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.SoutheastLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.SoutheastLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find south east dungeon location record name >%s<", dungeonLocationConfig.SoutheastLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.SoutheastLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.SouthLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.SouthLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find south dungeon location record name >%s<", dungeonLocationConfig.SouthLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.SouthLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.SouthwestLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.SouthwestLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find south west dungeon location record name >%s<", dungeonLocationConfig.SouthwestLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.SouthwestLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.WestLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.WestLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find west dungeon location record name >%s<", dungeonLocationConfig.WestLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.WestLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.NorthwestLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.NorthwestLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find north west dungeon location record name >%s<", dungeonLocationConfig.NorthwestLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.NorthwestLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.UpLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.UpLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find up dungeon location record name >%s<", dungeonLocationConfig.UpLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.UpLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.DownLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.DownLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find down dungeon location record name >%s<", dungeonLocationConfig.DownLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.DownLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}
		}
	}

	return data, nil
}

func (t *Testing) resolveConfigObjectCharacterIdentifiers(data Data, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	findCharacterRec := func(characterName string) *record.Character {
		for _, dungeonCharacterRec := range data.CharacterRecs {
			if dungeonCharacterRec.Name == characterName {
				return dungeonCharacterRec
			}
		}
		return nil
	}

	for idx := range dungeonConfig.ObjectConfig {
		if dungeonConfig.ObjectConfig[idx].CharacterName == "" {
			continue
		}
		characterRec := findCharacterRec(dungeonConfig.ObjectConfig[idx].CharacterName)
		if characterRec == nil {
			msg := fmt.Sprintf("Failed to find object's related character record with character name >%s<", dungeonConfig.ObjectConfig[idx].CharacterName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.ObjectConfig[idx].Record.CharacterID = sql.NullString{
			String: characterRec.ID,
			Valid:  true,
		}
	}

	return dungeonConfig, nil
}

func (t *Testing) resolveConfigObjectMonsterIdentifiers(data Data, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	findMonsterRec := func(monsterName string) *record.Monster {
		for _, dungeonMonsterRec := range data.MonsterRecs {
			if dungeonMonsterRec.Name == monsterName {
				return dungeonMonsterRec
			}
		}
		return nil
	}

	for idx := range dungeonConfig.ObjectConfig {
		if dungeonConfig.ObjectConfig[idx].MonsterName == "" {
			continue
		}
		dungeonMonsterRec := findMonsterRec(dungeonConfig.ObjectConfig[idx].MonsterName)
		if dungeonMonsterRec == nil {
			msg := fmt.Sprintf("Failed to find object's related monster record with monster name >%s<", dungeonConfig.ObjectConfig[idx].MonsterName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.ObjectConfig[idx].Record.MonsterID = sql.NullString{
			String: dungeonMonsterRec.ID,
			Valid:  true,
		}
	}

	return dungeonConfig, nil
}

func (t *Testing) resolveConfigLocationIdentifiers(data Data, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	findLocationRec := func(locationName string) *record.Location {
		for _, dungeonLocationRec := range data.LocationRecs {
			if dungeonLocationRec.Name == locationName {
				return dungeonLocationRec
			}
		}
		return nil
	}

	for idx := range dungeonConfig.MonsterConfig {
		dungeonLocationRec := findLocationRec(dungeonConfig.MonsterConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find monster dungeon location record name >%s<", dungeonConfig.MonsterConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.MonsterConfig[idx].Record.LocationID = dungeonLocationRec.ID
	}

	for idx := range dungeonConfig.ObjectConfig {
		if dungeonConfig.ObjectConfig[idx].LocationName == "" {
			continue
		}
		dungeonLocationRec := findLocationRec(dungeonConfig.ObjectConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find object dungeon location record name >%s<", dungeonConfig.ObjectConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.ObjectConfig[idx].Record.LocationID = sql.NullString{
			String: dungeonLocationRec.ID,
			Valid:  true,
		}

	}

	return dungeonConfig, nil
}

// resolveConfigDungeonIdentifiers assigned the provided dungeon record ID to configuration records
func (t *Testing) resolveConfigDungeonIdentifiers(dungeonRec *record.Dungeon, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	if dungeonConfig.LocationConfig != nil {
		for idx := range dungeonConfig.LocationConfig {
			dungeonConfig.LocationConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.MonsterConfig != nil {
		for idx := range dungeonConfig.MonsterConfig {
			dungeonConfig.MonsterConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.ObjectConfig != nil {
		for idx := range dungeonConfig.ObjectConfig {
			dungeonConfig.ObjectConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	return dungeonConfig, nil
}

func (t *Testing) AddCharacterTeardownID(id string) {
	rec := record.Character{}
	rec.ID = id

	t.teardownData.CharacterRecs = append(
		t.teardownData.CharacterRecs,
		rec,
	)
}

func (t *Testing) AddDungeonCharacterActionTeardownID(id string) {
	rec := record.Action{}
	rec.ID = id

	t.teardownData.ActionRecs = append(t.teardownData.ActionRecs, &rec)

	if t.CommitData {
		t.InitTx(nil)
	}

	actionRecordSet, err := t.Model.(*model.Model).GetActionRecordSet(id)
	if err != nil {
		t.Log.Warn("Failed fetch dungeon action record set >%v<", err)
		return
	}

	// Source
	if actionRecordSet.ActionCharacterRec != nil {
		t.Log.Info("Adding action character record ID >%s<", actionRecordSet.ActionCharacterRec.ID)
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, actionRecordSet.ActionCharacterRec)
		t.teardownData.ActionCharacterObjectRecs = append(t.teardownData.ActionCharacterObjectRecs, actionRecordSet.ActionCharacterObjectRecs...)
	}
	if actionRecordSet.ActionMonsterRec != nil {
		t.Log.Info("Adding action monster record ID >%s<", actionRecordSet.ActionMonsterRec.ID)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, actionRecordSet.ActionMonsterRec)
		t.teardownData.ActionMonsterObjectRecs = append(t.teardownData.ActionMonsterObjectRecs, actionRecordSet.ActionMonsterObjectRecs...)
	}

	// Current location
	if actionRecordSet.CurrentLocation != nil {
		dungeonActionLocationRecordSet := actionRecordSet.CurrentLocation
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	}

	// Target location
	if actionRecordSet.TargetLocation != nil {
		dungeonActionLocationRecordSet := actionRecordSet.TargetLocation
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	}

	// Targets
	if actionRecordSet.TargetActionCharacterRec != nil {
		t.Log.Info("Adding target action character record character ID >%s<", actionRecordSet.TargetActionCharacterRec.CharacterInstanceID)
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, actionRecordSet.TargetActionCharacterRec)
		t.teardownData.ActionCharacterObjectRecs = append(t.teardownData.ActionCharacterObjectRecs, actionRecordSet.TargetActionCharacterObjectRecs...)
	}
	if actionRecordSet.TargetActionMonsterRec != nil {
		t.Log.Info("Adding target action monster record monster ID >%s<", actionRecordSet.TargetActionMonsterRec.MonsterInstanceID)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, actionRecordSet.TargetActionMonsterRec)
		t.teardownData.ActionMonsterObjectRecs = append(t.teardownData.ActionMonsterObjectRecs, actionRecordSet.TargetActionMonsterObjectRecs...)
	}
	if actionRecordSet.TargetActionObjectRec != nil {
		t.Log.Info("Adding target action object record object ID >%s<", actionRecordSet.TargetActionObjectRec.ObjectInstanceID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, actionRecordSet.TargetActionObjectRec)
	}
	if actionRecordSet.StashedActionObjectRec != nil {
		t.Log.Info("Adding stashed action object record object ID >%s<", actionRecordSet.StashedActionObjectRec.ObjectInstanceID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, actionRecordSet.StashedActionObjectRec)
	}
	if actionRecordSet.EquippedActionObjectRec != nil {
		t.Log.Info("Adding equipped action object record object ID >%s<", actionRecordSet.EquippedActionObjectRec.ObjectInstanceID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, actionRecordSet.EquippedActionObjectRec)
	}
	if actionRecordSet.DroppedActionObjectRec != nil {
		t.Log.Info("Adding dropped action object record object ID >%s<", actionRecordSet.DroppedActionObjectRec.ObjectInstanceID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, actionRecordSet.DroppedActionObjectRec)
	}

	if t.CommitData {
		t.RollbackTx()
	}
}

// RemoveData -
func (t *Testing) RemoveData() error {

	seen := map[string]bool{}

	t.Log.Info("Removing >%d< dungeon action character object records", len(t.teardownData.ActionCharacterObjectRecs))

DUNGEON_ACTION_CHARACTER_OBJECT_RECS:
	for {
		if len(t.teardownData.ActionCharacterObjectRecs) == 0 {
			break DUNGEON_ACTION_CHARACTER_OBJECT_RECS
		}
		var rec *record.ActionCharacterObject
		rec, t.teardownData.ActionCharacterObjectRecs = t.teardownData.ActionCharacterObjectRecs[0], t.teardownData.ActionCharacterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionCharacterObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action character object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action character records", len(t.teardownData.ActionCharacterRecs))

DUNGEON_ACTION_CHARACTER_RECS:
	for {
		if len(t.teardownData.ActionCharacterRecs) == 0 {
			break DUNGEON_ACTION_CHARACTER_RECS
		}
		var rec *record.ActionCharacter
		rec, t.teardownData.ActionCharacterRecs = t.teardownData.ActionCharacterRecs[0], t.teardownData.ActionCharacterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionCharacterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action monster object records", len(t.teardownData.ActionMonsterObjectRecs))

DUNGEON_ACTION_MONSTER_OBJECT_RECS:
	for {
		if len(t.teardownData.ActionMonsterObjectRecs) == 0 {
			break DUNGEON_ACTION_MONSTER_OBJECT_RECS
		}
		var rec *record.ActionMonsterObject
		rec, t.teardownData.ActionMonsterObjectRecs = t.teardownData.ActionMonsterObjectRecs[0], t.teardownData.ActionMonsterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionMonsterObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action monster object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action monster records", len(t.teardownData.ActionMonsterRecs))

DUNGEON_ACTION_MONSTER_RECS:
	for {
		if len(t.teardownData.ActionMonsterRecs) == 0 {
			break DUNGEON_ACTION_MONSTER_RECS
		}
		var rec *record.ActionMonster
		rec, t.teardownData.ActionMonsterRecs = t.teardownData.ActionMonsterRecs[0], t.teardownData.ActionMonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionMonsterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action object records", len(t.teardownData.ActionObjectRecs))

DUNGEON_ACTION_OBJECT_RECS:
	for {
		if len(t.teardownData.ActionObjectRecs) == 0 {
			break DUNGEON_ACTION_OBJECT_RECS
		}
		var rec *record.ActionObject
		rec, t.teardownData.ActionObjectRecs = t.teardownData.ActionObjectRecs[0], t.teardownData.ActionObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_ACTION_RECS:
	for {
		if len(t.teardownData.ActionRecs) == 0 {
			break DUNGEON_ACTION_RECS
		}
		var rec *record.Action
		rec, t.teardownData.ActionRecs = t.teardownData.ActionRecs[0], t.teardownData.ActionRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_OBJECT_RECS:
	for {
		if len(t.teardownData.ObjectRecs) == 0 {
			break DUNGEON_OBJECT_RECS
		}
		var rec record.Object
		rec, t.teardownData.ObjectRecs = t.teardownData.ObjectRecs[0], t.teardownData.ObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_CHARACTER_RECS:
	for {
		if len(t.teardownData.CharacterRecs) == 0 {
			break DUNGEON_CHARACTER_RECS
		}
		var rec record.Character
		rec, t.teardownData.CharacterRecs = t.teardownData.CharacterRecs[0], t.teardownData.CharacterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveCharacterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_MONSTER_RECS:
	for {
		if len(t.teardownData.MonsterRecs) == 0 {
			break DUNGEON_MONSTER_RECS
		}
		var rec record.Monster
		rec, t.teardownData.MonsterRecs = t.teardownData.MonsterRecs[0], t.teardownData.MonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveMonsterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_LOCATION_RECS:
	for {
		if len(t.teardownData.LocationRecs) == 0 {
			break DUNGEON_LOCATION_RECS
		}
		var rec record.Location
		rec, t.teardownData.LocationRecs = t.teardownData.LocationRecs[0], t.teardownData.LocationRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon location record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_RECS:
	for {
		if len(t.teardownData.DungeonRecs) == 0 {
			break DUNGEON_RECS
		}
		var rec record.Dungeon
		rec, t.teardownData.DungeonRecs = t.teardownData.DungeonRecs[0], t.teardownData.DungeonRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Data = Data{}

	return nil
}
