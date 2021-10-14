package harness

import (
	"gitlab.com/alienspaces/go-mud/server/core/harness"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data       *Data
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
	t.Data = &Data{}

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

	for _, dungeonConfig := range t.DataConfig.DungeonConfig {

		dungeonRec, err := t.createDungeonRec(dungeonConfig)
		if err != nil {
			t.Log.Warn("Failed creating dungeon record >%v<", err)
			return err
		}
		t.Data.DungeonRecs = append(t.Data.DungeonRecs, dungeonRec)

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
			t.Data.DungeonLocationRecs = append(t.Data.DungeonLocationRecs, dungeonLocationRec)
		}
	}

	return nil
}

func (t *Testing) resolveConfigDungeonIdentifiers(dungeonRec record.Dungeon, dungeonConfig DungeonConfig) (DungeonConfig, error) {

	if dungeonConfig.DungeonLocationConfig != nil {
		for idx := range dungeonConfig.DungeonLocationConfig {
			dungeonConfig.DungeonLocationConfig[idx].Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.DungeonCharacterConfig != nil {
		for _, dungeonCharacterConfig := range dungeonConfig.DungeonCharacterConfig {
			dungeonCharacterConfig.Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.DungeonMonsterConfig != nil {
		for _, dungeonMonsterConfig := range dungeonConfig.DungeonMonsterConfig {
			dungeonMonsterConfig.Record.DungeonID = dungeonRec.ID
		}
	}

	if dungeonConfig.DungeonObjectConfig != nil {
		for _, dungeonObjectConfig := range dungeonConfig.DungeonObjectConfig {
			dungeonObjectConfig.Record.DungeonID = dungeonRec.ID
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
