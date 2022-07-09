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

// NewTesting -
func NewTesting(c configurer.Configurer, l logger.Logger, s storer.Storer, config DataConfig) (t *Testing, err error) {

	t = &Testing{
		Testing: harness.Testing{
			Config: c,
			Log:    l,
			Store:  s,
		},
	}

	// Require service config, logger and store
	if t.Config == nil || t.Log == nil || t.Store == nil {
		return nil, fmt.Errorf("missing configurer >%v<, logger >%v< or storer >%v<, cannot create new test harness", t.Config, t.Log, t.Store)
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
	l := t.Logger("Modeller")

	m, err := model.NewModel(t.Config, t.Log, t.Store)
	if err != nil {
		l.Warn("failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}

// CreateData - Custom data
func (t *Testing) CreateData() error {
	l := t.Logger("CreateData")

	data := &Data{}
	teardownData := teardownData{}

	// Create object records
	for _, objectConfig := range t.DataConfig.ObjectConfig {
		t.Log.Info("Creating object >%#v<", objectConfig)
		objectConfig, err := t.createObjectRec(objectConfig)
		if err != nil {
			l.Warn("Failed creating object record >%v<", err)
			return err
		}
		l.Info("+ Created object record ID >%s< Name >%s<", objectConfig.ID, objectConfig.Name)
		data.ObjectRecs = append(data.ObjectRecs, objectConfig)
		teardownData.ObjectRecs = append(teardownData.ObjectRecs, *objectConfig)
	}

	// Create monster records
	for _, monsterConfig := range t.DataConfig.MonsterConfig {
		monsterRec, err := t.createMonsterRec(monsterConfig)
		if err != nil {
			l.Warn("Failed creating monster record >%v<", err)
			return err
		}
		l.Info("+ Created monster record ID >%s< Name >%s<", monsterRec.ID, monsterRec.Name)
		data.MonsterRecs = append(data.MonsterRecs, monsterRec)
		teardownData.MonsterRecs = append(teardownData.MonsterRecs, *monsterRec)

		for _, monsterObjectConfig := range monsterConfig.MonsterObjectConfig {
			monsterObjectRec, err := t.createMonsterObjectRec(data, monsterRec, monsterObjectConfig)
			if err != nil {
				l.Warn("Failed creating monster object record >%v<", err)
				return err
			}
			l.Info("+ Created monster object record ID >%s< monster ID >%s< object ID", monsterObjectRec.ID, monsterObjectRec.MonsterID, monsterObjectRec.ObjectID)
			data.MonsterObjectRecs = append(data.MonsterObjectRecs, monsterObjectRec)
			teardownData.MonsterObjectRecs = append(teardownData.MonsterObjectRecs, *monsterObjectRec)
		}
	}

	// Dungeons
	for _, dungeonConfig := range t.DataConfig.DungeonConfig {

		// Create the dungeon record
		dungeonRec, err := t.createDungeonRec(dungeonConfig)
		if err != nil {
			l.Warn("Failed creating dungeon record >%v<", err)
			return err
		}
		l.Info("+ Created dungeon record ID >%s< Name >%s<", dungeonRec.ID, dungeonRec.Name)
		data.DungeonRecs = append(data.DungeonRecs, dungeonRec)
		teardownData.DungeonRecs = append(teardownData.DungeonRecs, *dungeonRec)

		// Create the location records
		for _, locationConfig := range dungeonConfig.LocationConfig {
			locationRec, err := t.createLocationRec(dungeonRec, locationConfig)
			if err != nil {
				l.Warn("Failed creating location record >%v<", err)
				return err
			}

			l.Info("+ Created location record ID >%s< Name >%s<", locationRec.ID, locationRec.Name)
			data.LocationRecs = append(data.LocationRecs, locationRec)
			teardownData.LocationRecs = append(teardownData.LocationRecs, *locationRec)

			// Create location objects
			for _, locationObjectConfig := range locationConfig.LocationObjectConfig {
				locationObjectRec, err := t.createLocationObjectRec(data, locationRec, locationObjectConfig)
				if err != nil {
					l.Warn("Failed creating location object record >%v<", err)
					return err
				}

				l.Info("+ Created location object record ID >%s< location ID >%s< object ID >%s<", locationObjectRec.ID, locationObjectRec.LocationID, locationObjectRec.ObjectID)
				data.LocationObjectRecs = append(data.LocationObjectRecs, locationObjectRec)
				teardownData.LocationObjectRecs = append(teardownData.LocationObjectRecs, *locationObjectRec)
			}

			// Create location monster
			for _, locationMonsterConfig := range locationConfig.LocationMonsterConfig {
				locationMonsterRec, err := t.createLocationMonsterRec(data, locationRec, locationMonsterConfig)
				if err != nil {
					l.Warn("Failed creating location monster record >%v<", err)
					return err
				}

				l.Info("+ Created location monster record ID >%s< location ID >%s< monster ID >%s<", locationMonsterRec.ID, locationMonsterRec.LocationID, locationMonsterRec.MonsterID)
				data.LocationMonsterRecs = append(data.LocationMonsterRecs, locationMonsterRec)
				teardownData.LocationMonsterRecs = append(teardownData.LocationMonsterRecs, *locationMonsterRec)
			}
		}

		// Resolve all location direction identifiers on all dungeon locations
		data, err = t.resolveDataLocationDirectionIdentifiers(data, dungeonConfig)
		if err != nil {
			l.Warn("Failed resolving config location identifiers >%v<", err)
			return err
		}

		// Update all previously created location records as they now have all their
		// reference location identifiers set correctly.
		for _, locationRec := range data.LocationRecs {
			err := t.updateLocationRec(locationRec)
			if err != nil {
				l.Warn("Failed updating location record >%v<", err)
				return err
			}
		}

		// Resolve monster and object config locations
		// dungeonConfig, err = t.resolveConfigLocationIdentifiers(data, dungeonConfig)
		// if err != nil {
		// 	l.Warn("Failed resolving config location identifiers >%v<", err)
		// 	return err
		// }

		// Resolve object config character identifiers
		// dungeonConfig, err = t.resolveConfigObjectCharacterIdentifiers(data, dungeonConfig)
		// if err != nil {
		// 	l.Warn("Failed resolving config object character identifiers >%v<", err)
		// 	return err
		// }

		// Resolve object config monster identifiers
		// dungeonConfig, err = t.resolveConfigObjectMonsterIdentifiers(data, dungeonConfig)
		// if err != nil {
		// 	l.Warn("Failed resolving config object monster identifiers >%v<", err)
		// 	return err
		// }

		// Create dungeon instances
		// for _, dungeonInstanceConfig := range dungeonConfig.DungeonInstanceConfig {

		// 	// Create dungeon instance record
		// 	t.Log.Info("Creating dungeon instance >%#v<", dungeonInstanceConfig)
		// 	dungeonInstanceRecordSet, err := t.createDungeonInstance(dungeonRec.ID)
		// 	if err != nil {
		// 		l.Warn("Failed creating dungeon instance record >%v<", err)
		// 		return err
		// 	}

		// 	l.Info("+ Created dungeon instance record set >%#v<", dungeonInstanceRecordSet)
		// 	data.DungeonInstanceRecs = append(data.DungeonInstanceRecs, dungeonInstanceRecordSet.DungeonInstanceRec)
		// 	teardownData.DungeonInstanceRecs = append(teardownData.DungeonInstanceRecs, dungeonInstanceRecordSet.DungeonInstanceRec)

		// 	data.LocationInstanceRecs = append(data.LocationInstanceRecs, dungeonInstanceRecordSet.LocationInstanceRecs...)
		// 	teardownData.LocationInstanceRecs = append(teardownData.LocationInstanceRecs, dungeonInstanceRecordSet.LocationInstanceRecs...)

		// 	data.MonsterInstanceRecs = append(data.MonsterInstanceRecs, dungeonInstanceRecordSet.MonsterInstanceRecs...)
		// 	teardownData.MonsterInstanceRecs = append(teardownData.MonsterInstanceRecs, dungeonInstanceRecordSet.MonsterInstanceRecs...)

		// 	data.ObjectInstanceRecs = append(data.ObjectInstanceRecs, dungeonInstanceRecordSet.ObjectInstanceRecs...)
		// 	teardownData.ObjectInstanceRecs = append(teardownData.ObjectInstanceRecs, dungeonInstanceRecordSet.ObjectInstanceRecs...)

		// 	// TODO: Create character instance records

		// 	// TODO: Create action records

		// }
	}

	// Characters
	for _, characterConfig := range t.DataConfig.CharacterConfig {
		characterRec, err := t.createCharacterRec(characterConfig)
		if err != nil {
			l.Warn("Failed creating character record >%v<", err)
			return err
		}

		l.Info("+ Created character record ID >%s< Name >%s<", characterRec.ID, characterRec.Name)
		data.CharacterRecs = append(data.CharacterRecs, characterRec)
		teardownData.CharacterRecs = append(teardownData.CharacterRecs, *characterRec)

		for _, characterObjectConfig := range characterConfig.CharacterObjectConfig {
			characterObjectRec, err := t.createCharacterObjectRec(data, characterRec, characterObjectConfig)
			if err != nil {
				l.Warn("Failed creating character object record >%v<", err)
				return err
			}

			l.Info("+ Created character object record ID >%s< character ID >%s< object ID", characterObjectRec.ID, characterObjectRec.CharacterID, characterObjectRec.ObjectID)
			data.CharacterObjectRecs = append(data.CharacterObjectRecs, characterObjectRec)
			teardownData.CharacterObjectRecs = append(teardownData.CharacterObjectRecs, *characterObjectRec)
		}

		// Character enter dungeon
		if characterConfig.CharacterDungeonConfig != nil {
			t.Log.Info("Character ID >%s< entering dungeon Name >%s<", characterRec.ID, characterConfig.CharacterDungeonConfig.DungeonName)
			// dungeonInstanceRecordSet, err := t.characterEnterDungeon(characterRec.ID, characterConfig.CharacterDungeonConfig.DungeonName)
			// if err != nil {
			// 	l.Warn("Failed character entering dungeon >%v<", err)
			// 	return err
			// }
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
	// 		l.Warn("Failed creating dungeon action record >%v<", err)
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

	// Assign data once we have successfully set up all data
	t.Data = *data
	t.teardownData = teardownData

	return nil
}

func (t *Testing) resolveDataLocationDirectionIdentifiers(d *Data, dungeonConfig DungeonConfig) (*Data, error) {

	findLocationRec := func(locationName string) *record.Location {
		for _, dungeonLocationRec := range d.LocationRecs {
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
				return d, fmt.Errorf(msg)
			}

			if dungeonLocationConfig.NorthLocationName != "" {
				targetLocationRec := findLocationRec(dungeonLocationConfig.NorthLocationName)
				if targetLocationRec == nil {
					msg := fmt.Sprintf("Failed to find north dungeon location record name >%s<", dungeonLocationConfig.NorthLocationName)
					t.Log.Error(msg)
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
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
					return d, fmt.Errorf(msg)
				}

				dungeonLocationRec.DownLocationID = sql.NullString{
					String: targetLocationRec.ID,
					Valid:  true,
				}
			}
		}
	}

	return d, nil
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
	l := t.Logger("AddDungeonCharacterActionTeardownID")

	rec := record.Action{}
	rec.ID = id

	t.teardownData.ActionRecs = append(t.teardownData.ActionRecs, &rec)

	if t.CommitData {
		t.InitTx(nil)
	}

	actionRecordSet, err := t.Model.(*model.Model).GetActionRecordSet(id)
	if err != nil {
		l.Warn("Failed fetch dungeon action record set >%v<", err)
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
	l := t.Logger("RemoveData")

	seen := map[string]bool{}

	t.Log.Info("Removing >%d< action character object records", len(t.teardownData.ActionCharacterObjectRecs))

ACTION_CHARACTER_OBJECT_RECS:
	for {
		if len(t.teardownData.ActionCharacterObjectRecs) == 0 {
			break ACTION_CHARACTER_OBJECT_RECS
		}
		var rec *record.ActionCharacterObject
		rec, t.teardownData.ActionCharacterObjectRecs = t.teardownData.ActionCharacterObjectRecs[0], t.teardownData.ActionCharacterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionCharacterObjectRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing action character object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< action character records", len(t.teardownData.ActionCharacterRecs))

ACTION_CHARACTER_RECS:
	for {
		if len(t.teardownData.ActionCharacterRecs) == 0 {
			break ACTION_CHARACTER_RECS
		}
		var rec *record.ActionCharacter
		rec, t.teardownData.ActionCharacterRecs = t.teardownData.ActionCharacterRecs[0], t.teardownData.ActionCharacterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionCharacterRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing action character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< action monster object records", len(t.teardownData.ActionMonsterObjectRecs))

ACTION_MONSTER_OBJECT_RECS:
	for {
		if len(t.teardownData.ActionMonsterObjectRecs) == 0 {
			break ACTION_MONSTER_OBJECT_RECS
		}
		var rec *record.ActionMonsterObject
		rec, t.teardownData.ActionMonsterObjectRecs = t.teardownData.ActionMonsterObjectRecs[0], t.teardownData.ActionMonsterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionMonsterObjectRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing action monster object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< action monster records", len(t.teardownData.ActionMonsterRecs))

ACTION_MONSTER_RECS:
	for {
		if len(t.teardownData.ActionMonsterRecs) == 0 {
			break ACTION_MONSTER_RECS
		}
		var rec *record.ActionMonster
		rec, t.teardownData.ActionMonsterRecs = t.teardownData.ActionMonsterRecs[0], t.teardownData.ActionMonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionMonsterRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing action monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< action object records", len(t.teardownData.ActionObjectRecs))

ACTION_OBJECT_RECS:
	for {
		if len(t.teardownData.ActionObjectRecs) == 0 {
			break ACTION_OBJECT_RECS
		}
		var rec *record.ActionObject
		rec, t.teardownData.ActionObjectRecs = t.teardownData.ActionObjectRecs[0], t.teardownData.ActionObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionObjectRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing action object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< action records", len(t.teardownData.ActionRecs))

ACTION_RECS:
	for {
		if len(t.teardownData.ActionRecs) == 0 {
			break ACTION_RECS
		}
		var rec *record.Action
		rec, t.teardownData.ActionRecs = t.teardownData.ActionRecs[0], t.teardownData.ActionRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveActionRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing action record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< monster object records", len(t.teardownData.MonsterObjectRecs))

MONSTER_OBJECT_RECS:
	for {
		if len(t.teardownData.MonsterObjectRecs) == 0 {
			break MONSTER_OBJECT_RECS
		}
		var rec record.MonsterObject
		rec, t.teardownData.MonsterObjectRecs = t.teardownData.MonsterObjectRecs[0], t.teardownData.MonsterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveMonsterObjectRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing monster object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< character object records", len(t.teardownData.CharacterObjectRecs))

CHARACTER_OBJECT_RECS:
	for {
		if len(t.teardownData.CharacterObjectRecs) == 0 {
			break CHARACTER_OBJECT_RECS
		}
		var rec record.CharacterObject
		rec, t.teardownData.CharacterObjectRecs = t.teardownData.CharacterObjectRecs[0], t.teardownData.CharacterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveCharacterObjectRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing character object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< location object records", len(t.teardownData.LocationObjectRecs))

LOCATION_OBJECT_RECS:
	for {
		if len(t.teardownData.LocationObjectRecs) == 0 {
			break LOCATION_OBJECT_RECS
		}
		var rec record.LocationObject
		rec, t.teardownData.LocationObjectRecs = t.teardownData.LocationObjectRecs[0], t.teardownData.LocationObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationObjectRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing location object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< object records", len(t.teardownData.ObjectRecs))

OBJECT_RECS:
	for {
		if len(t.teardownData.ObjectRecs) == 0 {
			break OBJECT_RECS
		}
		var rec record.Object
		rec, t.teardownData.ObjectRecs = t.teardownData.ObjectRecs[0], t.teardownData.ObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveObjectRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing dungeon object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< character records", len(t.teardownData.CharacterRecs))

CHARACTER_RECS:
	for {
		if len(t.teardownData.CharacterRecs) == 0 {
			break CHARACTER_RECS
		}
		var rec record.Character
		rec, t.teardownData.CharacterRecs = t.teardownData.CharacterRecs[0], t.teardownData.CharacterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveCharacterRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing dungeon character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< location monster records", len(t.teardownData.LocationMonsterRecs))

LOCATION_MONSTER_RECS:
	for {
		if len(t.teardownData.LocationMonsterRecs) == 0 {
			break LOCATION_MONSTER_RECS
		}
		var rec record.LocationMonster
		rec, t.teardownData.LocationMonsterRecs = t.teardownData.LocationMonsterRecs[0], t.teardownData.LocationMonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationMonsterRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing location monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< monster records", len(t.teardownData.MonsterRecs))

MONSTER_RECS:
	for {
		if len(t.teardownData.MonsterRecs) == 0 {
			break MONSTER_RECS
		}
		var rec record.Monster
		rec, t.teardownData.MonsterRecs = t.teardownData.MonsterRecs[0], t.teardownData.MonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveMonsterRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing dungeon monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< location records", len(t.teardownData.LocationRecs))

LOCATION_RECS:
	for {
		if len(t.teardownData.LocationRecs) == 0 {
			break LOCATION_RECS
		}
		var rec record.Location
		rec, t.teardownData.LocationRecs = t.teardownData.LocationRecs[0], t.teardownData.LocationRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationRec(rec.ID)
		if err != nil {
			l.Warn("Failed removing dungeon location record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Log.Info("Removing >%d< dungeon records", len(t.teardownData.DungeonRecs))

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
			l.Warn("Failed removing dungeon record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Data = Data{}

	return nil
}

// Logger - Returns a logger with package context and provided function context
func (t *Testing) Logger(functionName string) logger.Logger {
	return t.Log.WithPackageContext("harness").WithFunctionContext(functionName)
}
