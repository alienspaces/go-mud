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
	DungeonRecs                      []*record.Dungeon
	LocationRecs              []*record.Location
	DungeonCharacterRecs             []*record.DungeonCharacter
	DungeonMonsterRecs               []*record.DungeonMonster
	DungeonObjectRecs                []*record.DungeonObject
	ActionRecs                []*record.Action
	ActionCharacterRecs       []*record.ActionCharacter
	ActionCharacterObjectRecs []*record.ActionCharacterObject
	ActionMonsterRecs         []*record.ActionMonster
	ActionMonsterObjectRecs   []*record.ActionMonsterObject
	ActionObjectRecs          []*record.ActionObject
}

// teardownData -
type teardownData struct {
	DungeonRecs                      []record.Dungeon
	LocationRecs              []record.Location
	DungeonCharacterRecs             []record.DungeonCharacter
	DungeonMonsterRecs               []record.DungeonMonster
	DungeonObjectRecs                []record.DungeonObject
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

	for _, dungeonConfig := range t.DataConfig.DungeonConfig {

		dungeonRec, err := t.createDungeonRec(dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed creating dungeon record >%v<", err)
			return err
		}
		data.DungeonRecs = append(data.DungeonRecs, dungeonRec)
		teardownData.DungeonRecs = append(teardownData.DungeonRecs, *dungeonRec)

		dungeonConfig, err = t.resolveConfigDungeonIdentifiers(dungeonRec, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving dungeon config dungeon identifiers >%v<", err)
			return err
		}

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

		// Update all location records
		for _, dungeonLocationRec := range data.LocationRecs {
			err := t.updateLocationRec(dungeonLocationRec)
			if err != nil {
				t.Log.Warn("Failed updating dungeon location record >%v<", err)
				return err
			}
		}

		// Resolve character, monster and object config locations
		dungeonConfig, err = t.resolveConfigLocationIdentifiers(data, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving config location identifiers >%v<", err)
			return err
		}

		// Create character records
		for _, dungeonCharacterConfig := range dungeonConfig.DungeonCharacterConfig {
			dungeonCharacterRec, err := t.createDungeonCharacterRec(dungeonRec, dungeonCharacterConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon character record >%v<", err)
				return err
			}
			data.DungeonCharacterRecs = append(data.DungeonCharacterRecs, dungeonCharacterRec)
			teardownData.DungeonCharacterRecs = append(teardownData.DungeonCharacterRecs, *dungeonCharacterRec)
		}

		// Create monster records
		for _, dungeonMonsterConfig := range dungeonConfig.DungeonMonsterConfig {
			dungeonMonsterRec, err := t.createDungeonMonsterRec(dungeonRec, dungeonMonsterConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon monster record >%v<", err)
				return err
			}
			data.DungeonMonsterRecs = append(data.DungeonMonsterRecs, dungeonMonsterRec)
			teardownData.DungeonMonsterRecs = append(teardownData.DungeonMonsterRecs, *dungeonMonsterRec)
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
		for _, dungeonObjectConfig := range dungeonConfig.DungeonObjectConfig {
			t.Log.Info("Creating dungeon object >%#v<", dungeonObjectConfig)
			dungeonObjectRec, err := t.createDungeonObjectRec(dungeonRec, dungeonObjectConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon object record >%v<", err)
				return err
			}
			data.DungeonObjectRecs = append(data.DungeonObjectRecs, dungeonObjectRec)
			teardownData.DungeonObjectRecs = append(teardownData.DungeonObjectRecs, *dungeonObjectRec)
		}

		// Create action records
		for _, dungeonActionConfig := range dungeonConfig.ActionConfig {
			dungeonID := ""
			dungeonCharacterID := ""
			for _, characterRecord := range data.DungeonCharacterRecs {
				if characterRecord.Name == dungeonActionConfig.CharacterName {
					dungeonID = characterRecord.DungeonID
					dungeonCharacterID = characterRecord.ID
				}
			}
			if dungeonID == "" || dungeonCharacterID == "" {
				msg := fmt.Sprintf("Failed to find dungeon character record name >%s<", dungeonActionConfig.CharacterName)
				t.Log.Error(msg)
				return fmt.Errorf(msg)
			}

			dungeonActionRecordSet, err := t.createDungeonCharacterActionRec(dungeonID, dungeonCharacterID, dungeonActionConfig.Command)
			if err != nil {
				t.Log.Warn("Failed creating dungeon action record >%v<", err)
				return err
			}

			data.ActionRecs = append(data.ActionRecs, dungeonActionRecordSet.ActionRec)
			teardownData.ActionRecs = append(teardownData.ActionRecs, dungeonActionRecordSet.ActionRec)

			// Source
			if dungeonActionRecordSet.ActionCharacterRec != nil {
				data.ActionCharacterRecs = append(data.ActionCharacterRecs, dungeonActionRecordSet.ActionCharacterRec)
				teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, dungeonActionRecordSet.ActionCharacterRec)

				data.ActionCharacterObjectRecs = append(data.ActionCharacterObjectRecs, dungeonActionRecordSet.ActionCharacterObjectRecs...)
				teardownData.ActionCharacterObjectRecs = append(teardownData.ActionCharacterObjectRecs, dungeonActionRecordSet.ActionCharacterObjectRecs...)
			}
			if dungeonActionRecordSet.ActionMonsterRec != nil {
				data.ActionMonsterRecs = append(data.ActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)
				teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)

				data.ActionMonsterObjectRecs = append(data.ActionMonsterObjectRecs, dungeonActionRecordSet.ActionMonsterObjectRecs...)
				teardownData.ActionMonsterObjectRecs = append(teardownData.ActionMonsterObjectRecs, dungeonActionRecordSet.ActionMonsterObjectRecs...)
			}

			// Current location
			t.Log.Info("Dungeon action record set current location >%#v<", dungeonActionRecordSet.CurrentLocation)
			if dungeonActionRecordSet.CurrentLocation != nil {
				dungeonActionLocationRecordSet := dungeonActionRecordSet.CurrentLocation
				data.ActionCharacterRecs = append(data.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				data.ActionMonsterRecs = append(data.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				data.ActionObjectRecs = append(data.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)

				teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
			}

			// Target location
			t.Log.Info("Dungeon action record set target location >%#v<", dungeonActionRecordSet.TargetLocation)
			if dungeonActionRecordSet.TargetLocation != nil {
				dungeonActionLocationRecordSet := dungeonActionRecordSet.TargetLocation
				data.ActionCharacterRecs = append(data.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				data.ActionMonsterRecs = append(data.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				data.ActionObjectRecs = append(data.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)

				teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
			}

			// Targets
			if dungeonActionRecordSet.TargetActionCharacterRec != nil {
				data.ActionCharacterRecs = append(data.ActionCharacterRecs, dungeonActionRecordSet.TargetActionCharacterRec)
				teardownData.ActionCharacterRecs = append(teardownData.ActionCharacterRecs, dungeonActionRecordSet.TargetActionCharacterRec)

				data.ActionCharacterObjectRecs = append(data.ActionCharacterObjectRecs, dungeonActionRecordSet.TargetActionCharacterObjectRecs...)
				teardownData.ActionCharacterObjectRecs = append(teardownData.ActionCharacterObjectRecs, dungeonActionRecordSet.TargetActionCharacterObjectRecs...)
			}
			if dungeonActionRecordSet.TargetActionMonsterRec != nil {
				data.ActionMonsterRecs = append(data.ActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)
				teardownData.ActionMonsterRecs = append(teardownData.ActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)

				data.ActionMonsterObjectRecs = append(data.ActionMonsterObjectRecs, dungeonActionRecordSet.TargetActionMonsterObjectRecs...)
				teardownData.ActionMonsterObjectRecs = append(teardownData.ActionMonsterObjectRecs, dungeonActionRecordSet.TargetActionMonsterObjectRecs...)
			}
			if dungeonActionRecordSet.EquippedActionObjectRec != nil {
				data.ActionObjectRecs = append(data.ActionObjectRecs, dungeonActionRecordSet.EquippedActionObjectRec)
				teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, dungeonActionRecordSet.EquippedActionObjectRec)
			}
			if dungeonActionRecordSet.StashedActionObjectRec != nil {
				data.ActionObjectRecs = append(data.ActionObjectRecs, dungeonActionRecordSet.StashedActionObjectRec)
				teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, dungeonActionRecordSet.StashedActionObjectRec)
			}
			if dungeonActionRecordSet.TargetActionObjectRec != nil {
				data.ActionObjectRecs = append(data.ActionObjectRecs, dungeonActionRecordSet.TargetActionObjectRec)
				teardownData.ActionObjectRecs = append(teardownData.ActionObjectRecs, dungeonActionRecordSet.TargetActionObjectRec)
			}
		}
	}

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

	findDungeonCharacterRec := func(characterName string) *record.DungeonCharacter {
		for _, dungeonCharacterRec := range data.DungeonCharacterRecs {
			if dungeonCharacterRec.Name == characterName {
				return dungeonCharacterRec
			}
		}
		return nil
	}

	for idx := range dungeonConfig.DungeonObjectConfig {
		if dungeonConfig.DungeonObjectConfig[idx].CharacterName == "" {
			continue
		}
		dungeonCharacterRec := findDungeonCharacterRec(dungeonConfig.DungeonObjectConfig[idx].CharacterName)
		if dungeonCharacterRec == nil {
			msg := fmt.Sprintf("Failed to find object dungeon character record name >%s<", dungeonConfig.DungeonCharacterConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonObjectConfig[idx].Record.DungeonCharacterID = sql.NullString{
			String: dungeonCharacterRec.ID,
			Valid:  true,
		}
	}

	return dungeonConfig, nil
}

func (t *Testing) resolveConfigObjectMonsterIdentifiers(data Data, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	findDungeonMonsterRec := func(monsterName string) *record.DungeonMonster {
		for _, dungeonMonsterRec := range data.DungeonMonsterRecs {
			if dungeonMonsterRec.Name == monsterName {
				return dungeonMonsterRec
			}
		}
		return nil
	}

	for idx := range dungeonConfig.DungeonObjectConfig {
		if dungeonConfig.DungeonObjectConfig[idx].MonsterName == "" {
			continue
		}
		dungeonMonsterRec := findDungeonMonsterRec(dungeonConfig.DungeonObjectConfig[idx].MonsterName)
		if dungeonMonsterRec == nil {
			msg := fmt.Sprintf("Failed to find object dungeon monster record name >%s<", dungeonConfig.DungeonMonsterConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonObjectConfig[idx].Record.DungeonMonsterID = sql.NullString{
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

	for idx := range dungeonConfig.DungeonCharacterConfig {
		dungeonLocationRec := findLocationRec(dungeonConfig.DungeonCharacterConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find character dungeon location record name >%s<", dungeonConfig.DungeonCharacterConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonCharacterConfig[idx].Record.LocationID = dungeonLocationRec.ID
	}

	for idx := range dungeonConfig.DungeonMonsterConfig {
		dungeonLocationRec := findLocationRec(dungeonConfig.DungeonMonsterConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find monster dungeon location record name >%s<", dungeonConfig.DungeonMonsterConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonMonsterConfig[idx].Record.LocationID = dungeonLocationRec.ID
	}

	for idx := range dungeonConfig.DungeonObjectConfig {
		if dungeonConfig.DungeonObjectConfig[idx].LocationName == "" {
			continue
		}
		dungeonLocationRec := findLocationRec(dungeonConfig.DungeonObjectConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find object dungeon location record name >%s<", dungeonConfig.DungeonObjectConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonObjectConfig[idx].Record.LocationID = sql.NullString{
			String: dungeonLocationRec.ID,
			Valid:  true,
		}

	}

	return dungeonConfig, nil
}

func (t *Testing) resolveConfigDungeonIdentifiers(dungeonRec *record.Dungeon, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	if dungeonConfig.LocationConfig != nil {
		for idx := range dungeonConfig.LocationConfig {
			dungeonConfig.LocationConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.DungeonCharacterConfig != nil {
		for idx := range dungeonConfig.DungeonCharacterConfig {
			dungeonConfig.DungeonCharacterConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.DungeonMonsterConfig != nil {
		for idx := range dungeonConfig.DungeonMonsterConfig {
			dungeonConfig.DungeonMonsterConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.DungeonObjectConfig != nil {
		for idx := range dungeonConfig.DungeonObjectConfig {
			dungeonConfig.DungeonObjectConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	return dungeonConfig, nil
}

func (t *Testing) AddDungeonCharacterTeardownID(id string) {
	rec := record.DungeonCharacter{}
	rec.ID = id

	t.teardownData.DungeonCharacterRecs = append(
		t.teardownData.DungeonCharacterRecs,
		rec,
	)

	// TODO: Get all related character child records and add those for teardown
}

func (t *Testing) AddDungeonCharacterActionTeardownID(id string) {
	rec := record.Action{}
	rec.ID = id

	t.teardownData.ActionRecs = append(t.teardownData.ActionRecs, &rec)

	if t.CommitData {
		t.InitTx(nil)
	}

	dungeonActionRecordSet, err := t.Model.(*model.Model).GetActionRecordSet(id)
	if err != nil {
		t.Log.Warn("Failed fetch dungeon action record set >%v<", err)
		return
	}

	// Source
	if dungeonActionRecordSet.ActionCharacterRec != nil {
		t.Log.Info("Adding action character record ID >%s<", dungeonActionRecordSet.ActionCharacterRec.ID)
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, dungeonActionRecordSet.ActionCharacterRec)
		t.teardownData.ActionCharacterObjectRecs = append(t.teardownData.ActionCharacterObjectRecs, dungeonActionRecordSet.ActionCharacterObjectRecs...)
	}
	if dungeonActionRecordSet.ActionMonsterRec != nil {
		t.Log.Info("Adding action monster record ID >%s<", dungeonActionRecordSet.ActionMonsterRec.ID)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)
		t.teardownData.ActionMonsterObjectRecs = append(t.teardownData.ActionMonsterObjectRecs, dungeonActionRecordSet.ActionMonsterObjectRecs...)
	}

	// Current location
	if dungeonActionRecordSet.CurrentLocation != nil {
		dungeonActionLocationRecordSet := dungeonActionRecordSet.CurrentLocation
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	}

	// Target location
	if dungeonActionRecordSet.TargetLocation != nil {
		dungeonActionLocationRecordSet := dungeonActionRecordSet.TargetLocation
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	}

	// Targets
	if dungeonActionRecordSet.TargetActionCharacterRec != nil {
		t.Log.Info("Adding target action character record character ID >%s<", dungeonActionRecordSet.TargetActionCharacterRec.DungeonCharacterID)
		t.teardownData.ActionCharacterRecs = append(t.teardownData.ActionCharacterRecs, dungeonActionRecordSet.TargetActionCharacterRec)
		t.teardownData.ActionCharacterObjectRecs = append(t.teardownData.ActionCharacterObjectRecs, dungeonActionRecordSet.TargetActionCharacterObjectRecs...)
	}
	if dungeonActionRecordSet.TargetActionMonsterRec != nil {
		t.Log.Info("Adding target action monster record monster ID >%s<", dungeonActionRecordSet.TargetActionMonsterRec.DungeonMonsterID)
		t.teardownData.ActionMonsterRecs = append(t.teardownData.ActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)
		t.teardownData.ActionMonsterObjectRecs = append(t.teardownData.ActionMonsterObjectRecs, dungeonActionRecordSet.TargetActionMonsterObjectRecs...)
	}
	if dungeonActionRecordSet.TargetActionObjectRec != nil {
		t.Log.Info("Adding target action object record object ID >%s<", dungeonActionRecordSet.TargetActionObjectRec.DungeonObjectID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionRecordSet.TargetActionObjectRec)
	}
	if dungeonActionRecordSet.StashedActionObjectRec != nil {
		t.Log.Info("Adding stashed action object record object ID >%s<", dungeonActionRecordSet.StashedActionObjectRec.DungeonObjectID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionRecordSet.StashedActionObjectRec)
	}
	if dungeonActionRecordSet.EquippedActionObjectRec != nil {
		t.Log.Info("Adding equipped action object record object ID >%s<", dungeonActionRecordSet.EquippedActionObjectRec.DungeonObjectID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionRecordSet.EquippedActionObjectRec)
	}
	if dungeonActionRecordSet.DroppedActionObjectRec != nil {
		t.Log.Info("Adding dropped action object record object ID >%s<", dungeonActionRecordSet.DroppedActionObjectRec.DungeonObjectID)
		t.teardownData.ActionObjectRecs = append(t.teardownData.ActionObjectRecs, dungeonActionRecordSet.DroppedActionObjectRec)
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
		if len(t.teardownData.DungeonObjectRecs) == 0 {
			break DUNGEON_OBJECT_RECS
		}
		var rec record.DungeonObject
		rec, t.teardownData.DungeonObjectRecs = t.teardownData.DungeonObjectRecs[0], t.teardownData.DungeonObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_CHARACTER_RECS:
	for {
		if len(t.teardownData.DungeonCharacterRecs) == 0 {
			break DUNGEON_CHARACTER_RECS
		}
		var rec record.DungeonCharacter
		rec, t.teardownData.DungeonCharacterRecs = t.teardownData.DungeonCharacterRecs[0], t.teardownData.DungeonCharacterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonCharacterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_MONSTER_RECS:
	for {
		if len(t.teardownData.DungeonMonsterRecs) == 0 {
			break DUNGEON_MONSTER_RECS
		}
		var rec record.DungeonMonster
		rec, t.teardownData.DungeonMonsterRecs = t.teardownData.DungeonMonsterRecs[0], t.teardownData.DungeonMonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonMonsterRec(rec.ID)
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
