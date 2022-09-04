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

	l.Info("Creating test data")

	// Objects
	for _, objectConfig := range t.DataConfig.ObjectConfig {
		objectRec, err := t.createObjectRec(objectConfig)
		if err != nil {
			l.Warn("failed creating object record >%v<", err)
			return err
		}
		l.Debug("+ Created object record ID >%s< Name >%s<", objectRec.ID, objectRec.Name)
		data.AddObjectRec(objectRec)
		teardownData.AddObjectRec(objectRec)
	}

	// Monsters
	for _, monsterConfig := range t.DataConfig.MonsterConfig {
		monsterRec, err := t.createMonsterRec(monsterConfig)
		if err != nil {
			l.Warn("failed creating monster record >%v<", err)
			return err
		}
		l.Debug("+ Created monster record ID >%s< Name >%s<", monsterRec.ID, monsterRec.Name)
		data.AddMonsterRec(monsterRec)
		teardownData.AddMonsterRec(monsterRec)

		for _, monsterObjectConfig := range monsterConfig.MonsterObjectConfig {
			monsterObjectRec, err := t.createMonsterObjectRec(data, monsterRec, monsterObjectConfig)
			if err != nil {
				l.Warn("failed creating monster object record >%v<", err)
				return err
			}
			l.Debug("+ Created monster object record ID >%s< monster ID >%s< object ID", monsterObjectRec.ID, monsterObjectRec.MonsterID, monsterObjectRec.ObjectID)
			data.AddMonsterObjectRec(monsterObjectRec)
			teardownData.AddMonsterObjectRec(monsterObjectRec)
		}
	}

	// Characters
	for _, characterConfig := range t.DataConfig.CharacterConfig {

		characterRec, err := t.createCharacterRec(characterConfig)
		if err != nil {
			l.Warn("failed creating character record >%v<", err)
			return err
		}

		l.Debug("+ Created character record ID >%s< Name >%s<", characterRec.ID, characterRec.Name)
		data.AddCharacterRec(characterRec)
		teardownData.AddCharacterRec(characterRec)

		for _, characterObjectConfig := range characterConfig.CharacterObjectConfig {
			characterObjectRec, err := t.createCharacterObjectRec(data, characterRec, characterObjectConfig)
			if err != nil {
				l.Warn("failed creating character object record >%v<", err)
				return err
			}

			l.Debug("+ Created character object record ID >%s< character ID >%s< object ID >%s<", characterObjectRec.ID, characterObjectRec.CharacterID, characterObjectRec.ObjectID)
			data.AddCharacterObjectRec(characterObjectRec)
			teardownData.AddCharacterObjectRec(characterObjectRec)
		}
	}

	// Dungeons
	for _, dungeonConfig := range t.DataConfig.DungeonConfig {

		// Create the dungeon record
		dungeonRec, err := t.createDungeonRec(dungeonConfig)
		if err != nil {
			l.Warn("failed creating dungeon record >%v<", err)
			return err
		}
		l.Debug("+ Created dungeon record ID >%s< Name >%s<", dungeonRec.ID, dungeonRec.Name)
		data.AddDungeonRec(dungeonRec)
		teardownData.AddDungeonRec(dungeonRec)

		// Create the location records
		for _, locationConfig := range dungeonConfig.LocationConfig {
			locationRec, err := t.createLocationRec(dungeonRec, locationConfig)
			if err != nil {
				l.Warn("failed creating location record >%v<", err)
				return err
			}

			l.Debug("+ Created location record ID >%s< Name >%s<", locationRec.ID, locationRec.Name)
			data.AddLocationRec(locationRec)
			teardownData.AddLocationRec(locationRec)

			// Create location objects
			for _, locationObjectConfig := range locationConfig.LocationObjectConfig {
				locationObjectRec, err := t.createLocationObjectRec(data, locationRec, locationObjectConfig)
				if err != nil {
					l.Warn("failed creating location object record >%v<", err)
					return err
				}

				l.Debug("+ Created location object record ID >%s< location ID >%s< object ID >%s<", locationObjectRec.ID, locationObjectRec.LocationID, locationObjectRec.ObjectID)
				data.AddLocationObjectRec(locationObjectRec)
				teardownData.AddLocationObjectRec(locationObjectRec)
			}

			// Create location monster
			for _, locationMonsterConfig := range locationConfig.LocationMonsterConfig {
				locationMonsterRec, err := t.createLocationMonsterRec(data, locationRec, locationMonsterConfig)
				if err != nil {
					l.Warn("failed creating location monster record >%v<", err)
					return err
				}

				l.Debug("+ Created location monster record ID >%s< location ID >%s< monster ID >%s<", locationMonsterRec.ID, locationMonsterRec.LocationID, locationMonsterRec.MonsterID)
				data.AddLocationMonsterRec(locationMonsterRec)
				teardownData.AddLocationMonsterRec(locationMonsterRec)
			}
		}

		// Resolve all location direction identifiers on all dungeon locations
		data, err = t.resolveDataLocationDirectionIdentifiers(data, dungeonConfig)
		if err != nil {
			l.Warn("failed resolving config location identifiers >%v<", err)
			return err
		}

		// Update all previously created location records as they now have all their
		// reference location identifiers set correctly.
		for _, locationRec := range data.LocationRecs {
			err := t.updateLocationRec(locationRec)
			if err != nil {
				l.Warn("failed updating location record >%v<", err)
				return err
			}
		}

		// Dungeon Instances
		for _, dungeonInstanceConfig := range dungeonConfig.DungeonInstanceConfig {
			dungeonInstanceRecordSet, err := t.createDungeonInstance(dungeonRec.ID)
			if err != nil {
				l.Warn("failed creating dungeon instance >%v<", err)
				return err
			}

			data.AddDungeonInstanceRecordSet(dungeonInstanceRecordSet)
			teardownData.AddDungeonInstanceRecordSet(dungeonInstanceRecordSet)

			// Character Instances
			for _, characterInstanceConfig := range dungeonInstanceConfig.CharacterInstanceConfig {
				characterRec, err := data.GetCharacterRecByName(characterInstanceConfig.Name)
				if err != nil {
					l.Warn("failed getting character record >%v<", err)
					return err
				}

				dungeonInstanceRecordSet, characterInstanceRecordSet, err := t.characterEnterDungeon(dungeonRec.ID, characterRec.ID)
				if err != nil {
					l.Warn("failed character enter dungeon >%v<", err)
					return err
				}

				data.AddDungeonInstanceRecordSet(dungeonInstanceRecordSet)
				teardownData.AddDungeonInstanceRecordSet(dungeonInstanceRecordSet)

				data.AddCharacterInstanceRecordSet(characterInstanceRecordSet)
				teardownData.AddCharacterInstanceRecordSet(characterInstanceRecordSet)
			}

			// Actions
			for _, actionConfig := range dungeonInstanceConfig.ActionConfig {

				// Character action
				if actionConfig.CharacterName != "" {
					ciRec, err := data.GetCharacterInstanceRecByName(actionConfig.CharacterName)
					if err != nil {
						l.Warn("failed getting character instance record by name >%s< >%v<", actionConfig.CharacterName, err)
						return err
					}
					actionRecordSet, err := t.createCharacterActionRec(ciRec.DungeonInstanceID, ciRec.ID, actionConfig.Command)
					if err != nil {
						l.Warn("failed creating character action record >%v<", err)
						return err
					}

					data.AddActionRecordSet(actionRecordSet)
					teardownData.AddActionRecordSet(actionRecordSet)
				}

				// Monster action
				if actionConfig.MonsterName != "" {
					miRec, err := data.GetMonsterInstanceRecByName(actionConfig.MonsterName)
					if err != nil {
						l.Warn("failed getting monster instance record by name >%s< >%v<", actionConfig.MonsterName, err)
						return err
					}
					actionRecordSet, err := t.createMonsterActionRec(miRec.DungeonInstanceID, miRec.ID, actionConfig.Command)
					if err != nil {
						l.Warn("failed creating monster action record >%v<", err)
						return err
					}

					data.AddActionRecordSet(actionRecordSet)
					teardownData.AddActionRecordSet(actionRecordSet)
				}
			}
		}
	}

	// Assign data once we have successfully set up all data
	t.Data = *data
	t.teardownData = teardownData

	l.Info("Created test data")

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
	rec := &record.Character{}
	rec.ID = id
	t.teardownData.AddCharacterRec(rec)
}

func (t *Testing) AddActionTeardownID(id string) {
	l := t.Logger("AddActionTeardownID")

	if t.CommitData {
		t.InitTx(nil)
	}

	actionRecordSet, err := t.Model.(*model.Model).GetActionRecordSet(id)
	if err != nil {
		l.Warn("failed fetch dungeon action record set >%v<", err)
		return
	}

	t.teardownData.AddActionRecordSet(actionRecordSet)

	if t.CommitData {
		t.RollbackTx()
	}
}

// RemoveData -
func (t *Testing) RemoveData() error {
	l := t.Logger("RemoveData")

	l.Info("Removing test data")

	seen := map[string]bool{}

	l.Info("Removing >%d< action character object records", len(t.teardownData.ActionCharacterObjectRecs))

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
			l.Warn("failed removing action character object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< action character records", len(t.teardownData.ActionCharacterRecs))

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
			l.Warn("failed removing action character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< action monster object records", len(t.teardownData.ActionMonsterObjectRecs))

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
			l.Warn("failed removing action monster object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< action monster records", len(t.teardownData.ActionMonsterRecs))

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
			l.Warn("failed removing action monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< action object records", len(t.teardownData.ActionObjectRecs))

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
			l.Warn("failed removing action object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< action records", len(t.teardownData.ActionRecs))

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
			l.Warn("failed removing action record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< object instance records", len(t.teardownData.ObjectInstanceRecs))

OBJECT_INSTANCE_RECS:
	for {
		if len(t.teardownData.ObjectInstanceRecs) == 0 {
			break OBJECT_INSTANCE_RECS
		}
		var rec *record.ObjectInstance
		rec, t.teardownData.ObjectInstanceRecs = t.teardownData.ObjectInstanceRecs[0], t.teardownData.ObjectInstanceRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveObjectInstanceRec(rec.ID)
		if err != nil {
			l.Warn("failed removing object instance record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< monster instance records", len(t.teardownData.MonsterInstanceRecs))

MONSTER_INSTANCE_RECS:
	for {
		if len(t.teardownData.MonsterInstanceRecs) == 0 {
			break MONSTER_INSTANCE_RECS
		}
		var rec *record.MonsterInstance
		rec, t.teardownData.MonsterInstanceRecs = t.teardownData.MonsterInstanceRecs[0], t.teardownData.MonsterInstanceRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveMonsterInstanceRec(rec.ID)
		if err != nil {
			l.Warn("failed removing monster instance record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< character instance records", len(t.teardownData.CharacterInstanceRecs))

CHARACTER_INSTANCE_RECS:
	for {
		if len(t.teardownData.CharacterInstanceRecs) == 0 {
			break CHARACTER_INSTANCE_RECS
		}
		var rec *record.CharacterInstance
		rec, t.teardownData.CharacterInstanceRecs = t.teardownData.CharacterInstanceRecs[0], t.teardownData.CharacterInstanceRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveCharacterInstanceRec(rec.ID)
		if err != nil {
			l.Warn("failed removing character instance record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< location instance records", len(t.teardownData.LocationInstanceRecs))

LOCATION_INSTANCE_RECS:
	for {
		if len(t.teardownData.LocationInstanceRecs) == 0 {
			break LOCATION_INSTANCE_RECS
		}
		var rec *record.LocationInstance
		rec, t.teardownData.LocationInstanceRecs = t.teardownData.LocationInstanceRecs[0], t.teardownData.LocationInstanceRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationInstanceRec(rec.ID)
		if err != nil {
			l.Warn("failed removing location instance record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< dungeon instance records", len(t.teardownData.DungeonInstanceRecs))

DUNGEON_INSTANCE_RECS:
	for {
		if len(t.teardownData.DungeonInstanceRecs) == 0 {
			break DUNGEON_INSTANCE_RECS
		}
		var rec *record.DungeonInstance
		rec, t.teardownData.DungeonInstanceRecs = t.teardownData.DungeonInstanceRecs[0], t.teardownData.DungeonInstanceRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonInstanceRec(rec.ID)
		if err != nil {
			l.Warn("failed removing dungeon instance record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< monster object records", len(t.teardownData.MonsterObjectRecs))

MONSTER_OBJECT_RECS:
	for {
		if len(t.teardownData.MonsterObjectRecs) == 0 {
			break MONSTER_OBJECT_RECS
		}
		var rec *record.MonsterObject
		rec, t.teardownData.MonsterObjectRecs = t.teardownData.MonsterObjectRecs[0], t.teardownData.MonsterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveMonsterObjectRec(rec.ID)
		if err != nil {
			l.Warn("failed removing monster object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< character object records", len(t.teardownData.CharacterObjectRecs))

CHARACTER_OBJECT_RECS:
	for {
		if len(t.teardownData.CharacterObjectRecs) == 0 {
			break CHARACTER_OBJECT_RECS
		}
		var rec *record.CharacterObject
		rec, t.teardownData.CharacterObjectRecs = t.teardownData.CharacterObjectRecs[0], t.teardownData.CharacterObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveCharacterObjectRec(rec.ID)
		if err != nil {
			l.Warn("failed removing character object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< location object records", len(t.teardownData.LocationObjectRecs))

LOCATION_OBJECT_RECS:
	for {
		if len(t.teardownData.LocationObjectRecs) == 0 {
			break LOCATION_OBJECT_RECS
		}
		var rec *record.LocationObject
		rec, t.teardownData.LocationObjectRecs = t.teardownData.LocationObjectRecs[0], t.teardownData.LocationObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationObjectRec(rec.ID)
		if err != nil {
			l.Warn("failed removing location object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< object records", len(t.teardownData.ObjectRecs))

OBJECT_RECS:
	for {
		if len(t.teardownData.ObjectRecs) == 0 {
			break OBJECT_RECS
		}
		var rec *record.Object
		rec, t.teardownData.ObjectRecs = t.teardownData.ObjectRecs[0], t.teardownData.ObjectRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveObjectRec(rec.ID)
		if err != nil {
			l.Warn("failed removing dungeon object record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< character records", len(t.teardownData.CharacterRecs))

CHARACTER_RECS:
	for {
		if len(t.teardownData.CharacterRecs) == 0 {
			break CHARACTER_RECS
		}
		var rec *record.Character
		rec, t.teardownData.CharacterRecs = t.teardownData.CharacterRecs[0], t.teardownData.CharacterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveCharacterRec(rec.ID)
		if err != nil {
			l.Warn("failed removing dungeon character record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< location monster records", len(t.teardownData.LocationMonsterRecs))

LOCATION_MONSTER_RECS:
	for {
		if len(t.teardownData.LocationMonsterRecs) == 0 {
			break LOCATION_MONSTER_RECS
		}
		var rec *record.LocationMonster
		rec, t.teardownData.LocationMonsterRecs = t.teardownData.LocationMonsterRecs[0], t.teardownData.LocationMonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationMonsterRec(rec.ID)
		if err != nil {
			l.Warn("failed removing location monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< monster records", len(t.teardownData.MonsterRecs))

MONSTER_RECS:
	for {
		if len(t.teardownData.MonsterRecs) == 0 {
			break MONSTER_RECS
		}
		var rec *record.Monster
		rec, t.teardownData.MonsterRecs = t.teardownData.MonsterRecs[0], t.teardownData.MonsterRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveMonsterRec(rec.ID)
		if err != nil {
			l.Warn("failed removing dungeon monster record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< location records", len(t.teardownData.LocationRecs))

LOCATION_RECS:
	for {
		if len(t.teardownData.LocationRecs) == 0 {
			break LOCATION_RECS
		}
		var rec *record.Location
		rec, t.teardownData.LocationRecs = t.teardownData.LocationRecs[0], t.teardownData.LocationRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveLocationRec(rec.ID)
		if err != nil {
			l.Warn("failed removing dungeon location record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Info("Removing >%d< dungeon records", len(t.teardownData.DungeonRecs))

DUNGEON_RECS:
	for {
		if len(t.teardownData.DungeonRecs) == 0 {
			break DUNGEON_RECS
		}
		var rec *record.Dungeon
		rec, t.teardownData.DungeonRecs = t.teardownData.DungeonRecs[0], t.teardownData.DungeonRecs[1:]
		if seen[rec.ID] {
			continue
		}

		err := t.Model.(*model.Model).RemoveDungeonRec(rec.ID)
		if err != nil {
			l.Warn("failed removing dungeon record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	t.Data = Data{}

	l.Info("Removed test data")

	return nil
}

// Logger - Returns a logger with package context and provided function context
func (t *Testing) Logger(functionName string) logger.Logger {
	return t.Log.WithPackageContext("harness").WithFunctionContext(functionName)
}
