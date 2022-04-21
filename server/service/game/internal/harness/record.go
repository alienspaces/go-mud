package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (t *Testing) createDungeonRec(dungeonConfig DungeonConfig) (*record.Dungeon, error) {

	rec := dungeonConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Debug("(testing) Creating dungeon record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createLocationRec(dungeonRec *record.Dungeon, dungeonLocationConfig LocationConfig) (*record.Location, error) {

	rec := dungeonLocationConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Debug("(testing) Creating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).CreateLocationRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon location record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) updateLocationRec(rec *record.Location) error {

	t.Log.Debug("(testing) Updating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).UpdateLocationRec(rec)
	if err != nil {
		t.Log.Warn("(testing) Failed updating dungeon location record >%v<", err)
		return err
	}
	return nil
}

func (t *Testing) createDungeonCharacterRec(dungeonRec *record.Dungeon, dungeonCharacterConfig DungeonCharacterConfig) (*record.DungeonCharacter, error) {

	rec := dungeonCharacterConfig.Record

	// Default values
	if rec.Name == "" {
		rec.Name = gofakeit.Name() + " " + gofakeit.Name()
	}
	if rec.Strength == 0 {
		rec.Strength = 10
	}
	if rec.Dexterity == 0 {
		rec.Dexterity = 10
	}
	if rec.Intelligence == 0 {
		rec.Intelligence = 10
	}

	t.Log.Debug("(testing) Creating dungeon character record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonCharacterRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon character record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonMonsterRec(dungeonRec *record.Dungeon, dungeonMonsterConfig DungeonMonsterConfig) (*record.DungeonMonster, error) {

	rec := dungeonMonsterConfig.Record

	// Default values
	if rec.Name == "" {
		rec.Name = gofakeit.Name() + " " + gofakeit.Name()
	}
	if rec.Strength == 0 {
		rec.Strength = 10
	}
	if rec.Dexterity == 0 {
		rec.Dexterity = 10
	}
	if rec.Intelligence == 0 {
		rec.Intelligence = 10
	}

	t.Log.Debug("(testing) Creating dungeon monster record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonMonsterRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon monster record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonObjectRec(dungeonRec *record.Dungeon, dungeonObjectConfig DungeonObjectConfig) (*record.DungeonObject, error) {

	rec := dungeonObjectConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Debug("(testing) Creating dungeon object record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonObjectRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon object record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonCharacterActionRec(dungeonID, dungeonCharacterID, sentence string) (*model.ActionRecordSet, error) {

	t.Log.Debug("(testing) Creating dungeon action for character ID >%s< sentence >%s<", dungeonCharacterID, sentence)

	dungeonActionRecordSet, err := t.Model.(*model.Model).ProcessDungeonCharacterAction(dungeonID, dungeonCharacterID, sentence)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon character action record >%v<", err)
		return nil, err
	}

	return dungeonActionRecordSet, nil
}
