package harness

import (
	"github.com/brianvoe/gofakeit"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (t *Testing) createObjectRec(objectConfig ObjectConfig) (*record.Object, error) {

	rec := objectConfig.Record

	if rec.Name == "" {
		rec.Name = gofakeit.Name()
	}

	t.Log.Debug("Creating object record >%#v<", rec)

	err := t.Model.(*model.Model).CreateObjectRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating object record >%v<", err)
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

	t.Log.Debug("Creating monster record >%#v<", rec)

	err := t.Model.(*model.Model).CreateMonsterRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating monster record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createMonsterObjectRec(data *Data, monsterRec *record.Monster, monsterObjectConfig MonsterObjectConfig) (*record.MonsterObject, error) {

	objectRec, err := data.GetObjectRecByName(monsterObjectConfig.ObjectName)
	if err != nil {
		t.Log.Warn("Failed getting object record >%v<", err)
		return nil, err
	}

	rec := monsterObjectConfig.Record
	rec.MonsterID = monsterRec.ID
	rec.ObjectID = objectRec.ID

	t.Log.Debug("Creating monster object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateMonsterObjectRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating monster object record >%v<", err)
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

	t.Log.Debug("Creating dungeon character record >%#v<", rec)

	err := t.Model.(*model.Model).CreateCharacterRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon character record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createCharacterObjectRec(data *Data, characterRec *record.Character, characterObjectConfig CharacterObjectConfig) (*record.CharacterObject, error) {

	objectRec, err := data.GetObjectRecByName(characterObjectConfig.ObjectName)
	if err != nil {
		t.Log.Warn("Failed getting object record >%v<", err)
		return nil, err
	}

	rec := characterObjectConfig.Record
	rec.CharacterID = characterRec.ID
	rec.ObjectID = objectRec.ID

	t.Log.Debug("Creating character object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateCharacterObjectRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating character object record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createDungeonRec(dungeonConfig DungeonConfig) (*record.Dungeon, error) {

	rec := dungeonConfig.Record

	if rec.Name == "" {
		rec.Name = gofakeit.Name()
	}

	t.Log.Debug("Creating dungeon record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createLocationRec(dungeonRec *record.Dungeon, dungeonLocationConfig LocationConfig) (*record.Location, error) {

	rec := dungeonLocationConfig.Record
	rec.DungeonID = dungeonRec.ID

	t.Log.Debug("Creating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).CreateLocationRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating dungeon location record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createLocationObjectRec(data *Data, locationRec *record.Location, locationObjectConfig LocationObjectConfig) (*record.LocationObject, error) {

	objectRec, err := data.GetObjectRecByName(locationObjectConfig.ObjectName)
	if err != nil {
		t.Log.Warn("Failed getting object record >%v<", err)
		return nil, err
	}

	rec := locationObjectConfig.Record
	rec.LocationID = locationRec.ID
	rec.ObjectID = objectRec.ID

	if rec.SpawnPercentChance == 0 {
		rec.SpawnPercentChance = 100
	}

	t.Log.Debug("Creating location object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateLocationObjectRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating location object record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createLocationMonsterRec(data *Data, locationRec *record.Location, locationMonsterConfig LocationMonsterConfig) (*record.LocationMonster, error) {

	monsterRec, err := data.GetMonsterRecByName(locationMonsterConfig.MonsterName)
	if err != nil {
		t.Log.Warn("Failed getting monster record >%v<", err)
		return nil, err
	}

	rec := locationMonsterConfig.Record
	rec.LocationID = locationRec.ID
	rec.MonsterID = monsterRec.ID

	if rec.SpawnPercentChance == 0 {
		rec.SpawnPercentChance = 100
	}

	t.Log.Debug("Creating location monster record >%#v<", rec)

	err = t.Model.(*model.Model).CreateLocationMonsterRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating location monster record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) updateLocationRec(rec *record.Location) error {

	t.Log.Debug("Updating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).UpdateLocationRec(rec)
	if err != nil {
		t.Log.Warn("Failed updating dungeon location record >%v<", err)
		return err
	}
	return nil
}

func (t *Testing) characterEnterDungeon(dungeonID, characterID string) (*model.DungeonInstanceRecordSet, *model.CharacterInstanceRecordSet, error) {

	t.Log.Debug("Character ID >%s< entering dungeon ID >%s<", characterID, dungeonID)

	characterInstanceRecordSet, err := t.Model.(*model.Model).CharacterEnterDungeon(dungeonID, characterID)
	if err != nil {
		t.Log.Warn("Failed character entering dungeon >%v<", err)
		return nil, nil, err
	}

	// Fetch the dungeon instance the character was thrown into in the case where
	// a new dungeon instance was created to accomodate the additional character.
	dungeonInstanceRecordSet, err := t.Model.(*model.Model).GetDungeonInstanceRecordSet(characterInstanceRecordSet.CharacterInstanceRec.DungeonInstanceID)
	if err != nil {
		t.Log.Warn("Failed getting dungeon instance record set >%v<", err)
		return nil, nil, err
	}

	return dungeonInstanceRecordSet, characterInstanceRecordSet, nil
}

func (t *Testing) createDungeonInstance(dungeonID string) (*model.DungeonInstanceRecordSet, error) {

	t.Log.Debug("Creating dungeon ID >%s< instance", dungeonID)

	dungeonInstanceRecordSet, err := t.Model.(*model.Model).CreateDungeonInstance(dungeonID)
	if err != nil {
		t.Log.Warn("Failed creating dungeon instance >%v<", err)
		return nil, err
	}

	return dungeonInstanceRecordSet, nil
}

func (t *Testing) createCharacterActionRec(dungeonInstanceID, characterInstanceID, sentence string) (*record.ActionRecordSet, error) {

	t.Log.Debug("Creating action for dungeon instance ID >%s< character instance ID >%s< sentence >%s<", dungeonInstanceID, characterInstanceID, sentence)

	dungeonActionRecordSet, err := t.Model.(*model.Model).ProcessCharacterAction(dungeonInstanceID, characterInstanceID, sentence)
	if err != nil {
		t.Log.Warn("Failed creating character action record >%v<", err)
		return nil, err
	}

	return dungeonActionRecordSet, nil
}

func (t *Testing) createMonsterActionRec(dungeonInstanceID, monsterInstanceID, sentence string) (*record.ActionRecordSet, error) {

	t.Log.Debug("Creating action for dungeon instance ID >%s< monster instance ID >%s< sentence >%s<", dungeonInstanceID, monsterInstanceID, sentence)

	dungeonActionRecordSet, err := t.Model.(*model.Model).ProcessMonsterAction(dungeonInstanceID, monsterInstanceID, sentence)
	if err != nil {
		t.Log.Warn("Failed creating monster action record >%v<", err)
		return nil, err
	}

	return dungeonActionRecordSet, nil
}
