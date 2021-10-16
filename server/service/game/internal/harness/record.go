package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (t *Testing) createDungeonRec(dungeonConfig DungeonConfig) (*record.Dungeon, error) {

	rec := dungeonConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating dungeon record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonLocationRec(dungeonRec *record.Dungeon, dungeonLocationConfig DungeonLocationConfig) (*record.DungeonLocation, error) {

	rec := dungeonLocationConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonLocationRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon location record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) updateDungeonLocationRec(rec *record.DungeonLocation) error {

	t.Log.Info("Updating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).UpdateDungeonLocationRec(rec)
	if err != nil {
		t.Log.Warn("Failed updating dungeon location record >%v<", err)
		return err
	}
	return nil
}

func (t *Testing) createDungeonCharacterRec(dungeonRec *record.Dungeon, dungeonCharacterConfig DungeonCharacterConfig) (*record.DungeonCharacter, error) {

	rec := dungeonCharacterConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating dungeon character record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonCharacterRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon character record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonMonsterRec(dungeonRec *record.Dungeon, dungeonMonsterConfig DungeonMonsterConfig) (*record.DungeonMonster, error) {

	rec := dungeonMonsterConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating dungeon monster record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonMonsterRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon monster record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonObjectRec(dungeonRec *record.Dungeon, dungeonObjectConfig DungeonObjectConfig) (*record.DungeonObject, error) {

	rec := dungeonObjectConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating dungeon object record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonObjectRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon object record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonCharacterActionRec(dungeonCharacterID string, sentence string) (*record.DungeonAction, error) {

	t.Log.Info("Creating dungeon action for character ID >%s< sentence >%s<", dungeonCharacterID, sentence)

	dungeonActionRec, err := t.Model.(*model.Model).ProcessDungeonCharacterAction(dungeonCharacterID, sentence)
	if err != nil {
		t.Log.Warn("Failed creating dungeon character action record >%v<", err)
		return nil, err
	}
	return dungeonActionRec, nil
}
