package harness

import (
	"github.com/brianvoe/gofakeit"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (t *Testing) createObjectRec(objectConfig ObjectConfig) (*record.Object, error) {

	rec := objectConfig.Record

	t.Log.Debug("(testing) Creating object record >%#v<", rec)

	err := t.Model.(*model.Model).CreateObjectRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating object record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createMonsterRec(monsterConfig MonsterConfig) (*record.Monster, error) {

	rec := monsterConfig.Record

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

	t.Log.Debug("(testing) Creating monster record >%#v<", rec)

	err := t.Model.(*model.Model).CreateMonsterRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating monster record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createMonsterObjectRec(data *Data, monsterRec *record.Monster, monsterObjectConfig MonsterObjectConfig) (*record.MonsterObject, error) {

	objectRec, err := data.GetObjectRecByName(monsterObjectConfig.ObjectName)
	if err != nil {
		t.Log.Warn("(testing) Failed getting object record >%v<", err)
		return nil, err
	}

	rec := monsterObjectConfig.Record
	rec.MonsterID = monsterRec.ID
	rec.ObjectID = objectRec.ID

	t.Log.Debug("(testing) Creating monster object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateMonsterObjectRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating monster object record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createCharacterRec(characterConfig CharacterConfig) (*record.Character, error) {

	rec := characterConfig.Record

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

	err := t.Model.(*model.Model).CreateCharacterRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon character record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createCharacterObjectRec(data *Data, characterRec *record.Character, characterObjectConfig CharacterObjectConfig) (*record.CharacterObject, error) {

	objectRec, err := data.GetObjectRecByName(characterObjectConfig.ObjectName)
	if err != nil {
		t.Log.Warn("(testing) Failed getting object record >%v<", err)
		return nil, err
	}

	rec := characterObjectConfig.Record
	rec.CharacterID = characterRec.ID
	rec.ObjectID = objectRec.ID

	t.Log.Debug("(testing) Creating character object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateCharacterObjectRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating character object record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) characterEnterDungeon(dungeonID, characterID string) (*model.DungeonInstanceRecordSet, error) {

	t.Log.Info("(testing) Character ID >%s< entering dungeon ID >%s<", characterID, dungeonID)

	dungeonInstanceRecordSet, err := t.Model.(*model.Model).CharacterEnterDungeon(dungeonID, characterID)
	if err != nil {
		t.Log.Warn("(testing) Failed character entering dungeon >%v<", err)
		return nil, err
	}

	return dungeonInstanceRecordSet, nil
}

func (t *Testing) createDungeonRec(dungeonConfig DungeonConfig) (*record.Dungeon, error) {

	rec := dungeonConfig.Record

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
	rec.DungeonID = dungeonRec.ID

	t.Log.Debug("(testing) Creating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).CreateLocationRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon location record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createLocationObjectRec(data *Data, locationRec *record.Location, locationObjectConfig LocationObjectConfig) (*record.LocationObject, error) {

	objectRec, err := data.GetObjectRecByName(locationObjectConfig.ObjectName)
	if err != nil {
		t.Log.Warn("(testing) Failed getting object record >%v<", err)
		return nil, err
	}

	rec := locationObjectConfig.Record
	rec.LocationID = locationRec.ID
	rec.ObjectID = objectRec.ID

	if rec.SpawnPercentChance == 0 {
		rec.SpawnPercentChance = 100
	}

	t.Log.Debug("(testing) Creating location object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateLocationObjectRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating location object record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createLocationMonsterRec(data *Data, locationRec *record.Location, locationMonsterConfig LocationMonsterConfig) (*record.LocationMonster, error) {

	monsterRec, err := data.GetMonsterRecByName(locationMonsterConfig.MonsterName)
	if err != nil {
		t.Log.Warn("(testing) Failed getting monster record >%v<", err)
		return nil, err
	}

	rec := locationMonsterConfig.Record
	rec.LocationID = locationRec.ID
	rec.MonsterID = monsterRec.ID

	if rec.SpawnPercentChance == 0 {
		rec.SpawnPercentChance = 100
	}

	t.Log.Debug("(testing) Creating location monster record >%#v<", rec)

	err = t.Model.(*model.Model).CreateLocationMonsterRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating location monster record >%v<", err)
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

// TODO: Character actions are tied to dungeon instances
// func (t *Testing) createCharacterActionRec(dungeonID, dungeonCharacterID, sentence string) (*record.ActionRecordSet, error) {

// 	t.Log.Debug("(testing) Creating dungeon action for character ID >%s< sentence >%s<", dungeonCharacterID, sentence)

// 	dungeonActionRecordSet, err := t.Model.(*model.Model).ProcessCharacterAction(dungeonID, dungeonCharacterID, sentence)
// 	if err != nil {
// 		t.Log.Warn("(testing) Failed creating dungeon character action record >%v<", err)
// 		return nil, err
// 	}

// 	return dungeonActionRecordSet, nil
// }
