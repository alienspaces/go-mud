package harness

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/core/harness"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
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
	DungeonRecs                []*record.Dungeon
	DungeonLocationRecs        []*record.DungeonLocation
	DungeonCharacterRecs       []*record.DungeonCharacter
	DungeonMonsterRecs         []*record.DungeonMonster
	DungeonObjectRecs          []*record.DungeonObject
	DungeonActionRecs          []*record.DungeonAction
	DungeonActionCharacterRecs []*record.DungeonActionCharacter
	DungeonActionMonsterRecs   []*record.DungeonActionMonster
	DungeonActionObjectRecs    []*record.DungeonActionObject
}

// teardownData -
type teardownData struct {
	DungeonRecs                []record.Dungeon
	DungeonLocationRecs        []record.DungeonLocation
	DungeonCharacterRecs       []record.DungeonCharacter
	DungeonMonsterRecs         []record.DungeonMonster
	DungeonObjectRecs          []record.DungeonObject
	DungeonActionRecs          []*record.DungeonAction
	DungeonActionCharacterRecs []*record.DungeonActionCharacter
	DungeonActionMonsterRecs   []*record.DungeonActionMonster
	DungeonActionObjectRecs    []*record.DungeonActionObject
}

// NewTesting -
func NewTesting(config DataConfig) (t *Testing, err error) {

	// harness
	t = &Testing{}

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

	m, err := model.NewModel(t.Config, t.Log, t.Store)
	if err != nil {
		t.Log.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
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

		// Resolve all location identifiers on entities where entity
		// configuration references a location by name
		data, err = t.resolveDataLocationIdentifiers(data, dungeonConfig)
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

		// Create object records
		for _, dungeonObjectConfig := range dungeonConfig.DungeonObjectConfig {
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
			}
			if dungeonActionRecordSet.ActionMonsterRec != nil {
				data.DungeonActionMonsterRecs = append(data.DungeonActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)
				teardownData.DungeonActionMonsterRecs = append(teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)
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
			}
			if dungeonActionRecordSet.TargetActionMonsterRec != nil {
				data.DungeonActionMonsterRecs = append(data.DungeonActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)
				teardownData.DungeonActionMonsterRecs = append(teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)
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

func (t *Testing) resolveDataLocationIdentifiers(data Data, dungeonConfig DungeonConfig) (Data, error) {

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
		t.teardownData.DungeonActionCharacterRecs = append(t.teardownData.DungeonActionCharacterRecs, dungeonActionRecordSet.ActionCharacterRec)
	}
	if dungeonActionRecordSet.ActionMonsterRec != nil {
		t.teardownData.DungeonActionMonsterRecs = append(t.teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.ActionMonsterRec)
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
		t.teardownData.DungeonActionCharacterRecs = append(t.teardownData.DungeonActionCharacterRecs, dungeonActionRecordSet.TargetActionCharacterRec)
	}
	if dungeonActionRecordSet.TargetActionMonsterRec != nil {
		t.teardownData.DungeonActionMonsterRecs = append(t.teardownData.DungeonActionMonsterRecs, dungeonActionRecordSet.TargetActionMonsterRec)
	}
	if dungeonActionRecordSet.EquippedActionObjectRec != nil {
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.EquippedActionObjectRec)
	}
	if dungeonActionRecordSet.StashedActionObjectRec != nil {
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.StashedActionObjectRec)
	}
	if dungeonActionRecordSet.TargetActionObjectRec != nil {
		t.teardownData.DungeonActionObjectRecs = append(t.teardownData.DungeonActionObjectRecs, dungeonActionRecordSet.TargetActionObjectRec)
	}

	if t.CommitData {
		t.RollbackTx()
	}
}

// RemoveData -
func (t *Testing) RemoveData() error {

	seen := map[string]bool{}

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

DUNGEON_CHARACTER_RECS:
	for {
		if len(t.teardownData.DungeonCharacterRecs) == 0 {
			break DUNGEON_CHARACTER_RECS
		}
		// rec := &record.DungeonCharacter{}
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
		// rec := &record.DungeonMonster{}
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

DUNGEON_OBJECT_RECS:
	for {
		if len(t.teardownData.DungeonObjectRecs) == 0 {
			break DUNGEON_OBJECT_RECS
		}
		// rec := &record.DungeonObject{}
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

DUNGEON_LOCATION_RECS:
	for {
		if len(t.teardownData.DungeonLocationRecs) == 0 {
			break DUNGEON_LOCATION_RECS
		}
		// rec := &record.DungeonLocation{}
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
		// rec := &record.Dungeon{}
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

	return nil
}
