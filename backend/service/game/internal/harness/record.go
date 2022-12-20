package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"

	"gitlab.com/alienspaces/go-mud/backend/core/repository"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// UniqueName appends a UUID4 to the end of the name to make it unique
// for parallel test execution.
func UniqueName(name string) string {
	if name == "" {
		gofakeit.Name()
	}
	return fmt.Sprintf("%s (%s)", name, repository.NewRecordID())
}

// removes the unique UUID4 from the end of the name to make it normal for
// test harness functions that return a record based on its non unique name.
func NormalName(name string) string {
	return name[:len(name)-39]
}

func (t *Testing) createObjectRec(objectConfig ObjectConfig) (*record.Object, error) {
	l := t.Logger("createObjectRec")

	rec := objectConfig.Record

	rec.Name = UniqueName(rec.Name)

	l.Debug("Creating object record >%#v<", rec)

	err := t.Model.(*model.Model).CreateObjectRec(&rec)
	if err != nil {
		l.Warn("Failed creating object record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createMonsterRec(monsterConfig MonsterConfig) (*record.Monster, error) {
	l := t.Logger("createMonsterRec")

	rec := monsterConfig.Record

	rec.Name = UniqueName(rec.Name)

	// Default values
	if rec.Strength == 0 {
		rec.Strength = 10
	}
	if rec.Dexterity == 0 {
		rec.Dexterity = 10
	}
	if rec.Intelligence == 0 {
		rec.Intelligence = 10
	}

	l.Debug("Creating monster record >%#v<", rec)

	err := t.Model.(*model.Model).CreateMonsterRec(&rec)
	if err != nil {
		l.Warn("Failed creating monster record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createMonsterObjectRec(data *Data, monsterRec *record.Monster, monsterObjectConfig MonsterObjectConfig) (*record.MonsterObject, error) {
	l := t.Logger("createMonsterObjectRec")

	objectRec, err := data.GetObjectRecByName(monsterObjectConfig.ObjectName)
	if err != nil {
		l.Warn("Failed getting object record >%v<", err)
		return nil, err
	}

	rec := monsterObjectConfig.Record
	rec.MonsterID = monsterRec.ID
	rec.ObjectID = objectRec.ID

	l.Debug("Creating monster object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateMonsterObjectRec(&rec)
	if err != nil {
		l.Warn("Failed creating monster object record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createCharacterRec(characterConfig CharacterConfig) (*record.Character, error) {
	l := t.Logger("createCharacterRec")

	rec := characterConfig.Record

	rec.Name = UniqueName(rec.Name)

	// Default values
	if rec.Strength == 0 {
		rec.Strength = 10
	}
	if rec.Dexterity == 0 {
		rec.Dexterity = 10
	}
	if rec.Intelligence == 0 {
		rec.Intelligence = 10
	}

	l.Debug("Creating dungeon character record >%#v<", rec)

	err := t.Model.(*model.Model).CreateCharacterRec(&rec)
	if err != nil {
		l.Warn("Failed creating dungeon character record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createCharacterObjectRec(data *Data, characterRec *record.Character, characterObjectConfig CharacterObjectConfig) (*record.CharacterObject, error) {
	l := t.Logger("createCharacterObjectRec")

	objectRec, err := data.GetObjectRecByName(characterObjectConfig.ObjectName)
	if err != nil {
		l.Warn("Failed getting object record >%v<", err)
		return nil, err
	}

	rec := characterObjectConfig.Record
	rec.CharacterID = characterRec.ID
	rec.ObjectID = objectRec.ID

	l.Debug("Creating character object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateCharacterObjectRec(&rec)
	if err != nil {
		l.Warn("Failed creating character object record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createDungeonRec(dungeonConfig DungeonConfig) (*record.Dungeon, error) {
	l := t.Logger("createDungeonRec")

	rec := dungeonConfig.Record

	rec.Name = UniqueName(rec.Name)

	l.Debug("Creating dungeon record >%#v<", rec)

	err := t.Model.(*model.Model).CreateDungeonRec(&rec)
	if err != nil {
		l.Warn("Failed creating dungeon record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createLocationRec(dungeonRec *record.Dungeon, dungeonLocationConfig LocationConfig) (*record.Location, error) {
	l := t.Logger("createLocationRec")

	rec := dungeonLocationConfig.Record
	rec.DungeonID = dungeonRec.ID
	rec.Name = UniqueName(rec.Name)

	l.Debug("Creating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).CreateLocationRec(&rec)
	if err != nil {
		l.Warn("Failed creating dungeon location record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createLocationObjectRec(data *Data, locationRec *record.Location, locationObjectConfig LocationObjectConfig) (*record.LocationObject, error) {
	l := t.Logger("createLocationObjectRec")

	objectRec, err := data.GetObjectRecByName(locationObjectConfig.ObjectName)
	if err != nil {
		l.Warn("Failed getting object record >%v<", err)
		return nil, err
	}

	rec := locationObjectConfig.Record
	rec.LocationID = locationRec.ID
	rec.ObjectID = objectRec.ID

	if rec.SpawnPercentChance == 0 {
		rec.SpawnPercentChance = 100
	}

	l.Debug("Creating location object record >%#v<", rec)

	err = t.Model.(*model.Model).CreateLocationObjectRec(&rec)
	if err != nil {
		l.Warn("Failed creating location object record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) createLocationMonsterRec(data *Data, locationRec *record.Location, locationMonsterConfig LocationMonsterConfig) (*record.LocationMonster, error) {
	l := t.Logger("createLocationMonsterRec")

	monsterRec, err := data.GetMonsterRecByName(locationMonsterConfig.MonsterName)
	if err != nil {
		l.Warn("Failed getting monster record >%v<", err)
		return nil, err
	}

	rec := locationMonsterConfig.Record
	rec.LocationID = locationRec.ID
	rec.MonsterID = monsterRec.ID

	if rec.SpawnPercentChance == 0 {
		rec.SpawnPercentChance = 100
	}

	l.Debug("Creating location monster record >%#v<", rec)

	err = t.Model.(*model.Model).CreateLocationMonsterRec(&rec)
	if err != nil {
		l.Warn("Failed creating location monster record >%v<", err)
		return nil, err
	}

	return &rec, nil
}

func (t *Testing) updateLocationRec(rec *record.Location) error {
	l := t.Logger("updateLocationRec")

	l.Debug("Updating dungeon location record >%#v<", rec)

	err := t.Model.(*model.Model).UpdateLocationRec(rec)
	if err != nil {
		l.Warn("Failed updating dungeon location record >%v<", err)
		return err
	}
	return nil
}

func (t *Testing) characterEnterDungeon(dungeonID, characterID string) (*model.DungeonInstanceRecordSet, *model.CharacterInstanceRecordSet, error) {
	l := t.Logger("characterEnterDungeon")

	l.Debug("Character ID >%s< entering dungeon ID >%s<", characterID, dungeonID)

	characterInstanceRecordSet, err := t.Model.(*model.Model).CharacterEnterDungeon(dungeonID, characterID)
	if err != nil {
		l.Warn("Failed character entering dungeon >%v<", err)
		return nil, nil, err
	}

	// Fetch the dungeon instance the character was thrown into in the case where
	// a new dungeon instance was created to accomodate the additional character.
	dungeonInstanceRecordSet, err := t.Model.(*model.Model).GetDungeonInstanceRecordSet(characterInstanceRecordSet.CharacterInstanceRec.DungeonInstanceID)
	if err != nil {
		l.Warn("Failed getting dungeon instance record set >%v<", err)
		return nil, nil, err
	}

	return dungeonInstanceRecordSet, characterInstanceRecordSet, nil
}

func (t *Testing) createDungeonInstance(dungeonID string) (*model.DungeonInstanceRecordSet, error) {
	l := t.Logger("createDungeonInstance")

	l.Debug("Creating dungeon ID >%s< instance", dungeonID)

	dungeonInstanceRecordSet, err := t.Model.(*model.Model).CreateDungeonInstance(dungeonID)
	if err != nil {
		l.Warn("Failed creating dungeon instance >%v<", err)
		return nil, err
	}

	return dungeonInstanceRecordSet, nil
}

func (t *Testing) createCharacterActionRec(dungeonInstanceID, characterInstanceID, sentence string) (*record.ActionRecordSet, error) {
	l := t.Logger("createCharacterActionRec")

	l.Debug("Creating action for dungeon instance ID >%s< character instance ID >%s< sentence >%s<", dungeonInstanceID, characterInstanceID, sentence)

	dungeonActionRecordSet, err := t.Model.(*model.Model).ProcessCharacterAction(dungeonInstanceID, characterInstanceID, sentence)
	if err != nil {
		l.Warn("Failed creating character action record >%v<", err)
		return nil, err
	}

	return dungeonActionRecordSet, nil
}

func (t *Testing) createMonsterActionRec(dungeonInstanceID, monsterInstanceID, sentence string) (*record.ActionRecordSet, error) {
	l := t.Logger("createMonsterActionRec")

	l.Debug("Creating action for dungeon instance ID >%s< monster instance ID >%s< sentence >%s<", dungeonInstanceID, monsterInstanceID, sentence)

	dungeonActionRecordSet, err := t.Model.(*model.Model).ProcessMonsterAction(dungeonInstanceID, monsterInstanceID, sentence)
	if err != nil {
		l.Warn("Failed creating monster action record >%v<", err)
		return nil, err
	}

	return dungeonActionRecordSet, nil
}

func (t *Testing) createTurnRec(dungeonInstanceID string, turnConfig TurnConfig) (*record.Turn, error) {
	l := t.Logger("createTurnRec")

	rec := turnConfig.Record

	rec.DungeonInstanceID = dungeonInstanceID

	l.Debug("Creating turn record >%#v<", rec)

	err := t.Model.(*model.Model).CreateTurnRec(&rec)
	if err != nil {
		l.Warn("Failed creating turn record >%v<", err)
		return nil, err
	}
	return &rec, nil
}
