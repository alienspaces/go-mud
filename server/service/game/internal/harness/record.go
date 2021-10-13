package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (t *Testing) createDungeonRec(dungeonConfig DungeonConfig) (record.Dungeon, error) {

	rec := dungeonConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating testing record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating testing dungeon record >%v<", err)
		return rec, err
	}
	return rec, nil
}

func (t *Testing) createDungeonLocationRec(dungeonRec record.Dungeon, dungeonLocationConfig DungeonLocationConfig) (record.DungeonLocation, error) {

	rec := dungeonLocationConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating testing record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonLocationRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating testing dungeon location record >%v<", err)
		return rec, err
	}
	return rec, nil
}
