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
	DungeonLocationRecs              []*record.DungeonLocation
	DungeonCharacterRecs             []*record.DungeonCharacter
	DungeonMonsterRecs               []*record.DungeonMonster
	DungeonObjectRecs                []*record.DungeonObject
	DungeonActionRecs                []*record.DungeonAction
	DungeonActionCharacterRecs       []*record.DungeonActionCharacter
	DungeonActionCharacterObjectRecs []*record.DungeonActionCharacterObject
	DungeonActionMonsterRecs         []*record.DungeonActionMonster
	DungeonActionMonsterObjectRecs   []*record.DungeonActionMonsterObject
	DungeonActionObjectRecs          []*record.DungeonActionObject
}

// teardownData -
type teardownData struct {
	DungeonRecs                      []record.Dungeon
	DungeonLocationRecs              []record.DungeonLocation
	DungeonCharacterRecs             []record.DungeonCharacter
	DungeonMonsterRecs               []record.DungeonMonster
	DungeonObjectRecs                []record.DungeonObject
	DungeonActionRecs                []*record.DungeonAction
	DungeonActionCharacterRecs       []*record.DungeonActionCharacter
	DungeonActionCharacterObjectRecs []*record.DungeonActionCharacterObject
	DungeonActionMonsterRecs         []*record.DungeonActionMonster
	DungeonActionMonsterObjectRecs   []*record.DungeonActionMonsterObject
	DungeonActionObjectRecs          []*record.DungeonActionObject
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

		for _, dungeonLocationConfig := range dungeonConfig.DungeonLocationConfig {
			dungeonLocationRec, err := t.createDungeonLocationRec(dungeonRec, dungeonLocationConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon location record >%v<", err)
				return err
			}
			data.DungeonLocationRecs = append(data.DungeonLocationRecs, dungeonLocationRec)
			teardownData.DungeonLocationRecs = append(teardownData.DungeonLocationRecs, *dungeonLocationRec)
		}

		// Resolve all location direction identifiers on all dungeon locations
		data, err = t.resolveDataLocationDirectionIdentifiers(data, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving config location identifiers >%v<", err)
			return err
		}

		// Update all location records
		for _, dungeonLocationRec := range data.DungeonLocationRecs {
			err := t.updateDungeonLocationRec(dungeonLocationRec)
			if err != nil {
				t.Log.Warn("Failed updating dungeon location record >%v<", err)
				return err
			}
		}

		// Resolve character, monster and object config locations
		dungeonConfig, err = t.resolveConfigDungeonLocationIdentifiers(data, dungeonConfig)
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
			t.Log.Warn("Creating dungeon object >%#v<", dungeonObjectConfig)
			dungeonObjectRec, err := t.createDungeonObjectRec(dungeonRec, dungeonObjectConfig)
			if err != nil {
				t.Log.Warn("Failed creating dungeon object record >%v<", err)
				return err
			}
			data.DungeonObjectRecs = append(data.DungeonObjectRecs, dungeonObjectRec)
			teardownData.DungeonObjectRecs = append(teardownData.DungeonObjectRecs, *dungeonObjectRec)
		}

		// Create action records
		for _, dungeonActionConfig := range dungeonConfig.DungeonActionConfig {
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

			data.DungeonActionRecs = append(data.DungeonActionRecs, dungeonActionRecordSet.ActionRec)
			teardownData.DungeonActionRecs = append(teardownData.DungeonActionRecs, dungeonActionRecordSet.ActionRec)

			// Source
			if dungeonActionRecordSet.ActionCharacterRec != nil {
				data.DungeonActionCharacterRecs = append(data.DungeonActionCharacterRecs, dungeonActionRecordSet.ActionCharacterRec)
				teardownData.DungeonActionCharacterRecs = append(teardownData.DungeonActionCharacterRecs, dungeonActionRecordSet.ActionCharacterRec)

				data.DungeonActionCharacterObjectRecs = append(data.DungeonActionCharacterObjectRecs, dungeonActionRecordSet.ActionCharacterObjectRecs...)
				teardownData.DungeonActionCharacterObjectRecs = append(teardownData.DungeonActionCharacterObjectRecs, dungeonActionRecordSet.ActionCharacterObjectRecs...)
			}
			if dungeonActionRecordSet.ActionMonsterRec != nil {
				data.DungeonActionMonsterRecs = append(data.DungeonActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)
				teardownData.DungeonActionMonsterRecs = append(teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)

				data.DungeonActionMonsterObjectRecs = append(data.DungeonActionMonsterObjectRecs, dungeonActionRecordSet.ActionMonsterObjectRecs...)
				teardownData.DungeonActionMonsterObjectRecs = append(teardownData.DungeonActionMonsterObjectRecs, dungeonActionRecordSet.ActionMonsterObjectRecs...)
			}

			// Current location
			t.Log.Info("Dungeon action record set current location >%#v<", dungeonActionRecordSet.CurrentLocation)
			if dungeonActionRecordSet.CurrentLocation != nil {
				dungeonActionLocationRecordSet := dungeonActionRecordSet.CurrentLocation
				data.DungeonActionCharacterRecs = append(data.DungeonActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				data.DungeonActionMonsterRecs = append(data.DungeonActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				data.DungeonActionObjectRecs = append(data.DungeonActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)

				teardownData.DungeonActionCharacterRecs = append(teardownData.DungeonActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				teardownData.DungeonActionMonsterRecs = append(teardownData.DungeonActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				teardownData.DungeonActionObjectRecs = append(teardownData.DungeonActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
			}

			// Target location
			t.Log.Info("Dungeon action record set target location >%#v<", dungeonActionRecordSet.TargetLocation)
			if dungeonActionRecordSet.TargetLocation != nil {
				dungeonActionLocationRecordSet := dungeonActionRecordSet.TargetLocation
				data.DungeonActionCharacterRecs = append(data.DungeonActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				data.DungeonActionMonsterRecs = append(data.DungeonActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				data.DungeonActionObjectRecs = append(data.DungeonActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)

				teardownData.DungeonActionCharacterRecs = append(teardownData.DungeonActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
				teardownData.DungeonActionMonsterRecs = append(teardownData.DungeonActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
				teardownData.DungeonActionObjectRecs = append(teardownData.DungeonActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
			}

			// Targets
			if dungeonActionRecordSet.TargetActionCharacterRec != nil {
				data.DungeonActionCharacterRecs = append(data.DungeonActionCharacterRecs, dungeonActionRecordSet.TargetActionCharacterRec)
				teardownData.DungeonActionCharacterRecs = append(teardownData.DungeonActionCharacterRecs, dungeonActionRecordSet.TargetActionCharacterRec)

				data.DungeonActionCharacterObjectRecs = append(data.DungeonActionCharacterObjectRecs, dungeonActionRecordSet.TargetActionCharacterObjectRecs...)
				teardownData.DungeonActionCharacterObjectRecs = append(teardownData.DungeonActionCharacterObjectRecs, dungeonActionRecordSet.TargetActionCharacterObjectRecs...)
			}
			if dungeonActionRecordSet.TargetActionMonsterRec != nil {
				data.DungeonActionMonsterRecs = append(data.DungeonActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)
				teardownData.DungeonActionMonsterRecs = append(teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)

				data.DungeonActionMonsterObjectRecs = append(data.DungeonActionMonsterObjectRecs, dungeonActionRecordSet.TargetActionMonsterObjectRecs...)
				teardownData.DungeonActionMonsterObjectRecs = append(teardownData.DungeonActionMonsterObjectRecs, dungeonActionRecordSet.TargetActionMonsterObjectRecs...)
			}
			if dungeonActionRecordSet.EquippedActionObjectRec != nil {
				data.DungeonActionObjectRecs = append(data.DungeonActionObjectRecs, dungeonActionRecordSet.EquippedActionObjectRec)
				teardownData.DungeonActionObjectRecs = append(teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.EquippedActionObjectRec)
			}
			if dungeonActionRecordSet.StashedActionObjectRec != nil {
				data.DungeonActionObjectRecs = append(data.DungeonActionObjectRecs, dungeonActionRecordSet.StashedActionObjectRec)
				teardownData.DungeonActionObjectRecs = append(teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.StashedActionObjectRec)
			}
			if dungeonActionRecordSet.TargetActionObjectRec != nil {
				data.DungeonActionObjectRecs = append(data.DungeonActionObjectRecs, dungeonActionRecordSet.TargetActionObjectRec)
				teardownData.DungeonActionObjectRecs = append(teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.TargetActionObjectRec)
			}
		}
	}

	// Assign data once we have successfully set up ll data
	t.Data = data
	t.teardownData = teardownData

	return nil
}

func (t *Testing) resolveDataLocationDirectionIdentifiers(data Data, dungeonConfig DungeonConfig) (Data, error) {

	findDungeonLocationRec := func(locationName string) *record.DungeonLocation {
		for _, dungeonLocationRec := range data.DungeonLocationRecs {
			if dungeonLocationRec.Name == locationName {
				return dungeonLocationRec
			}
		}
		return nil
	}

	if dungeonConfig.DungeonLocationConfig != nil {
		for _, dungeonLocationConfig := range dungeonConfig.DungeonLocationConfig {
			dungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.Record.Name)
			if dungeonLocationRec == nil {
				msg := fmt.Sprintf("Failed to find dungeon location record name >%s<", dungeonLocationConfig.Record.Name)
				t.Log.Error(msg)
				return data, fmt.Errorf(msg)
			}

			if dungeonLocationConfig.NorthLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.NorthLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find north dungeon location record name >%s<", dungeonLocationConfig.NorthLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.NorthDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.NortheastLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.NortheastLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find north east dungeon location record name >%s<", dungeonLocationConfig.NortheastLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.NortheastDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.EastLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.EastLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find east dungeon location record name >%s<", dungeonLocationConfig.EastLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.EastDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.SoutheastLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.SoutheastLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find south east dungeon location record name >%s<", dungeonLocationConfig.SoutheastLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.SoutheastDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.SouthLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.SouthLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find south dungeon location record name >%s<", dungeonLocationConfig.SouthLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.SouthDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.SouthwestLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.SouthwestLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find south west dungeon location record name >%s<", dungeonLocationConfig.SouthwestLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.SouthwestDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.WestLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.WestLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find west dungeon location record name >%s<", dungeonLocationConfig.WestLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.WestDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.NorthwestLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.NorthwestLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find north west dungeon location record name >%s<", dungeonLocationConfig.NorthwestLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.NorthwestDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.UpLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.UpLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find up dungeon location record name >%s<", dungeonLocationConfig.UpLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.UpDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
					Valid:  true,
				}
			}

			if dungeonLocationConfig.DownLocationName != "" {
				targetDungeonLocationRec := findDungeonLocationRec(dungeonLocationConfig.DownLocationName)
				if targetDungeonLocationRec == nil {
					msg := fmt.Sprintf("Failed to find down dungeon location record name >%s<", dungeonLocationConfig.DownLocationName)
					t.Log.Error(msg)
					return data, fmt.Errorf(msg)
				}

				dungeonLocationRec.DownDungeonLocationID = sql.NullString{
					String: targetDungeonLocationRec.ID,
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

func (t *Testing) resolveConfigDungeonLocationIdentifiers(data Data, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	findDungeonLocationRec := func(locationName string) *record.DungeonLocation {
		for _, dungeonLocationRec := range data.DungeonLocationRecs {
			if dungeonLocationRec.Name == locationName {
				return dungeonLocationRec
			}
		}
		return nil
	}

	for idx := range dungeonConfig.DungeonCharacterConfig {
		dungeonLocationRec := findDungeonLocationRec(dungeonConfig.DungeonCharacterConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find character dungeon location record name >%s<", dungeonConfig.DungeonCharacterConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonCharacterConfig[idx].Record.DungeonLocationID = dungeonLocationRec.ID
	}

	for idx := range dungeonConfig.DungeonMonsterConfig {
		dungeonLocationRec := findDungeonLocationRec(dungeonConfig.DungeonMonsterConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find monster dungeon location record name >%s<", dungeonConfig.DungeonMonsterConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonMonsterConfig[idx].Record.DungeonLocationID = dungeonLocationRec.ID
	}

	for idx := range dungeonConfig.DungeonObjectConfig {
		if dungeonConfig.DungeonObjectConfig[idx].LocationName == "" {
			continue
		}
		dungeonLocationRec := findDungeonLocationRec(dungeonConfig.DungeonObjectConfig[idx].LocationName)
		if dungeonLocationRec == nil {
			msg := fmt.Sprintf("Failed to find object dungeon location record name >%s<", dungeonConfig.DungeonObjectConfig[idx].LocationName)
			t.Log.Error(msg)
			return dungeonConfig, fmt.Errorf(msg)
		}
		dungeonConfig.DungeonObjectConfig[idx].Record.DungeonLocationID = sql.NullString{
			String: dungeonLocationRec.ID,
			Valid:  true,
		}

	}

	return dungeonConfig, nil
}

func (t *Testing) resolveConfigDungeonIdentifiers(dungeonRec *record.Dungeon, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	if dungeonConfig.DungeonLocationConfig != nil {
		for idx := range dungeonConfig.DungeonLocationConfig {
			dungeonConfig.DungeonLocationConfig[idx].Record.DungeonID = dungeonRec.ID
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
	rec := record.DungeonAction{}
	rec.ID = id

	t.teardownData.DungeonActionRecs = append(t.teardownData.DungeonActionRecs, &rec)

	if t.CommitData {
		t.InitTx(nil)
	}

	dungeonActionRecordSet, err := t.Model.(*model.Model).GetDungeonActionRecordSet(id)
	if err != nil {
		t.Log.Warn("Failed fetch dungeon action record set >%v<", err)
		return
	}

	// Source
	if dungeonActionRecordSet.ActionCharacterRec != nil {
		t.Log.Info("Adding action character record ID >%s<", dungeonActionRecordSet.ActionCharacterRec.ID)
		t.teardownData.DungeonActionCharacterRecs = append(t.teardownData.DungeonActionCharacterRecs, dungeonActionRecordSet.ActionCharacterRec)
		t.teardownData.DungeonActionCharacterObjectRecs = append(t.teardownData.DungeonActionCharacterObjectRecs, dungeonActionRecordSet.ActionCharacterObjectRecs...)
	}
	if dungeonActionRecordSet.ActionMonsterRec != nil {
		t.Log.Info("Adding action monster record ID >%s<", dungeonActionRecordSet.ActionMonsterRec.ID)
		t.teardownData.DungeonActionMonsterRecs = append(t.teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)
		t.teardownData.DungeonActionMonsterObjectRecs = append(t.teardownData.DungeonActionMonsterObjectRecs, dungeonActionRecordSet.ActionMonsterObjectRecs...)
	}

	// Current location
	if dungeonActionRecordSet.CurrentLocation != nil {
		dungeonActionLocationRecordSet := dungeonActionRecordSet.CurrentLocation
		t.teardownData.DungeonActionCharacterRecs = append(t.teardownData.DungeonActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
		t.teardownData.DungeonActionMonsterRecs = append(t.teardownData.DungeonActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	}

	// Target location
	if dungeonActionRecordSet.TargetLocation != nil {
		dungeonActionLocationRecordSet := dungeonActionRecordSet.TargetLocation
		t.teardownData.DungeonActionCharacterRecs = append(t.teardownData.DungeonActionCharacterRecs, dungeonActionLocationRecordSet.ActionCharacterRecs...)
		t.teardownData.DungeonActionMonsterRecs = append(t.teardownData.DungeonActionMonsterRecs, dungeonActionLocationRecordSet.ActionMonsterRecs...)
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionLocationRecordSet.ActionObjectRecs...)
	}

	// Targets
	if dungeonActionRecordSet.TargetActionCharacterRec != nil {
		t.Log.Info("Adding target action character record character ID >%s<", dungeonActionRecordSet.TargetActionCharacterRec.DungeonCharacterID)
		t.teardownData.DungeonActionCharacterRecs = append(t.teardownData.DungeonActionCharacterRecs, dungeonActionRecordSet.TargetActionCharacterRec)
		t.teardownData.DungeonActionCharacterObjectRecs = append(t.teardownData.DungeonActionCharacterObjectRecs, dungeonActionRecordSet.TargetActionCharacterObjectRecs...)
	}
	if dungeonActionRecordSet.TargetActionMonsterRec != nil {
		t.Log.Info("Adding target action monster record monster ID >%s<", dungeonActionRecordSet.TargetActionMonsterRec.DungeonMonsterID)
		t.teardownData.DungeonActionMonsterRecs = append(t.teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)
		t.teardownData.DungeonActionMonsterObjectRecs = append(t.teardownData.DungeonActionMonsterObjectRecs, dungeonActionRecordSet.TargetActionMonsterObjectRecs...)
	}
	if dungeonActionRecordSet.TargetActionObjectRec != nil {
		t.Log.Info("Adding target action object record object ID >%s<", dungeonActionRecordSet.TargetActionObjectRec.DungeonObjectID)
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.TargetActionObjectRec)
	}
	if dungeonActionRecordSet.StashedActionObjectRec != nil {
		t.Log.Info("Adding stashed action object record object ID >%s<", dungeonActionRecordSet.StashedActionObjectRec.DungeonObjectID)
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.StashedActionObjectRec)
	}
	if dungeonActionRecordSet.EquippedActionObjectRec != nil {
		t.Log.Info("Adding equipped action object record object ID >%s<", dungeonActionRecordSet.EquippedActionObjectRec.DungeonObjectID)
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.EquippedActionObjectRec)
	}
	if dungeonActionRecordSet.DroppedActionObjectRec != nil {
		t.Log.Info("Adding dropped action object record object ID >%s<", dungeonActionRecordSet.DroppedActionObjectRec.DungeonObjectID)
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.DroppedActionObjectRec)
	}

	if t.CommitData {
		t.RollbackTx()
	}
}

// RemoveData -
func (t *Testing) RemoveData() error {

	seen := map[string]bool{}

	t.Log.Info("Removing >%d< dungeon action character object records", len(t.teardownData.DungeonActionCharacterObjectRecs))

DUNGEON_ACTION_CHARACTER_OBJECT_RECS:
	for {
		if len(t.teardownData.DungeonActionCharacterObjectRecs) == 0 {
			break DUNGEON_ACTION_CHARACTER_OBJECT_RECS
		}
		var rec *record.DungeonActionCharacterObject
		rec, t.teardownData.DungeonActionCharacterObjectRecs = t.teardownData.DungeonActionCharacterObjectRecs[0], t.teardownData.DungeonActionCharacterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonActionCharacterObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action character object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action character records", len(t.teardownData.DungeonActionCharacterRecs))

DUNGEON_ACTION_CHARACTER_RECS:
	for {
		if len(t.teardownData.DungeonActionCharacterRecs) == 0 {
			break DUNGEON_ACTION_CHARACTER_RECS
		}
		var rec *record.DungeonActionCharacter
		rec, t.teardownData.DungeonActionCharacterRecs = t.teardownData.DungeonActionCharacterRecs[0], t.teardownData.DungeonActionCharacterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonActionCharacterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action monster object records", len(t.teardownData.DungeonActionMonsterObjectRecs))

DUNGEON_ACTION_MONSTER_OBJECT_RECS:
	for {
		if len(t.teardownData.DungeonActionMonsterObjectRecs) == 0 {
			break DUNGEON_ACTION_MONSTER_OBJECT_RECS
		}
		var rec *record.DungeonActionMonsterObject
		rec, t.teardownData.DungeonActionMonsterObjectRecs = t.teardownData.DungeonActionMonsterObjectRecs[0], t.teardownData.DungeonActionMonsterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonActionMonsterObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action monster object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action monster records", len(t.teardownData.DungeonActionMonsterRecs))

DUNGEON_ACTION_MONSTER_RECS:
	for {
		if len(t.teardownData.DungeonActionMonsterRecs) == 0 {
			break DUNGEON_ACTION_MONSTER_RECS
		}
		var rec *record.DungeonActionMonster
		rec, t.teardownData.DungeonActionMonsterRecs = t.teardownData.DungeonActionMonsterRecs[0], t.teardownData.DungeonActionMonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonActionMonsterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon action object records", len(t.teardownData.DungeonActionObjectRecs))

DUNGEON_ACTION_OBJECT_RECS:
	for {
		if len(t.teardownData.DungeonActionObjectRecs) == 0 {
			break DUNGEON_ACTION_OBJECT_RECS
		}
		var rec *record.DungeonActionObject
		rec, t.teardownData.DungeonActionObjectRecs = t.teardownData.DungeonActionObjectRecs[0], t.teardownData.DungeonActionObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonActionObjectRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon action object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

DUNGEON_ACTION_RECS:
	for {
		if len(t.teardownData.DungeonActionRecs) == 0 {
			break DUNGEON_ACTION_RECS
		}
		var rec *record.DungeonAction
		rec, t.teardownData.DungeonActionRecs = t.teardownData.DungeonActionRecs[0], t.teardownData.DungeonActionRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonActionRec(rec.ID)
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
		if len(t.teardownData.DungeonLocationRecs) == 0 {
			break DUNGEON_LOCATION_RECS
		}
		var rec record.DungeonLocation
		rec, t.teardownData.DungeonLocationRecs = t.teardownData.DungeonLocationRecs[0], t.teardownData.DungeonLocationRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonLocationRec(rec.ID)
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
