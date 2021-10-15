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
	Data       Data
	DataConfig DataConfig
}

// Data -
type Data struct {
	DungeonRecs          []record.Dungeon
	DungeonLocationRecs  []record.DungeonLocation
	DungeonCharacterRecs []record.DungeonCharacter
	DungeonMonsterRecs   []record.DungeonMonster
	DungeonObjectRecs    []record.DungeonObject
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

	for _, dungeonConfig := range t.DataConfig.DungeonConfig {

		dungeonRec, err := t.createDungeonRec(dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed creating dungeon record >%v<", err)
			return err
		}
		data.DungeonRecs = append(t.Data.DungeonRecs, dungeonRec)

		dungeonConfig, err = t.resolveConfigDungeonIdentifiers(dungeonRec, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving dungeon config dungeon identifiers >%v<", err)
			return err
		}

		for _, dungeonLocationConfig := range dungeonConfig.DungeonLocationConfig {
			dungeonLocationRec, err := t.createDungeonLocationRec(dungeonRec, dungeonLocationConfig)
			if err != nil {
				t.Log.Warn("Failed creating game record >%v<", err)
				return err
			}
			data.DungeonLocationRecs = append(data.DungeonLocationRecs, dungeonLocationRec)
		}

		// Resolve all location identifiers on entities where entity
		// configuration references a location by name
		data, err = t.resolveDataLocationIdentifiers(data, dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed resolving config location identifiers >%v<", err)
			return err
		}

		// Update all location records

		// Resolve character config locations

		// Resolve monster config locations

		// Resolve object config locations

	}

	// Assign data once we have successfully set up ll data
	t.Data = data

	return nil
}

func (t *Testing) resolveDataLocationIdentifiers(data Data, dungeonConfig DungeonConfig) (Data, error) {

	findDungeonLocationRec := func(locationName string) *record.DungeonLocation {
		for _, dungeonLocationRec := range data.DungeonLocationRecs {
			if dungeonLocationRec.Name == locationName {
				return &dungeonLocationRec
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

func (t *Testing) resolveConfigDungeonIdentifiers(dungeonRec record.Dungeon, dungeonConfig DungeonConfig) (DungeonConfig, error) {

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

// RemoveData -
func (t *Testing) RemoveData() error {

DUNGEON_LOCATION_RECS:
	for {
		if len(t.Data.DungeonLocationRecs) == 0 {
			break DUNGEON_LOCATION_RECS
		}
		rec := record.DungeonLocation{}
		rec, t.Data.DungeonLocationRecs = t.Data.DungeonLocationRecs[0], t.Data.DungeonLocationRecs[1:]

		err := t.Model.(*model.Model).RemoveDungeonLocationRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon location record >%v<", err)
			return err
		}
	}

DUNGEON_RECS:
	for {
		if len(t.Data.DungeonRecs) == 0 {
			break DUNGEON_RECS
		}
		rec := record.Dungeon{}
		rec, t.Data.DungeonRecs = t.Data.DungeonRecs[0], t.Data.DungeonRecs[1:]

		err := t.Model.(*model.Model).RemoveDungeonRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing dungeon record >%v<", err)
			return err
		}
	}

	return nil
}
