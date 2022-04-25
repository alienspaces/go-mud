package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

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

func (t *Testing) createMonsterRec(dungeonRec *record.Dungeon, dungeonMonsterConfig MonsterConfig) (*record.Monster, error) {

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

	err := t.Model.(*model.Model).CreateMonsterRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon monster record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

func (t *Testing) createObjectRec(dungeonRec *record.Dungeon, dungeonObjectConfig ObjectConfig) (*record.Object, error) {

	rec := dungeonObjectConfig.Record

	t.Log.Debug("(testing) Creating dungeon object record >%#v<", rec)

	err := t.Model.(*model.Model).CreateObjectRec(&rec)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon object record >%v<", err)
		return nil, err
	}
	return &rec, nil
}

// Dungeon instance
func (t *Testing) createDungeonInstanceRec(dungeonID string) (*record.DungeonInstance, error) {

	t.Log.Debug("(testing) Creating dungeon instance for dungeon ID >%s<", dungeonID)

	rec, err := t.Model.(*model.Model).CreateDungeonInstance(dungeonID)
	if err != nil {
		t.Log.Warn("(testing) Failed creating dungeon instance record >%v<", err)
		return nil, err
	}
	return rec, nil
}

func (t *Testing) getLocationInstanceRecs(dungeonInstanceID string) ([]*record.LocationInstance, error) {

	locationInstanceRecs, err := t.Model.(*model.Model).GetLocationInstanceRecs(
		map[string]interface{}{
			"dungeon_instance_id": dungeonInstanceID,
		},
		nil, false,
	)
	if err != nil {
		t.Log.Warn("(testing) Failed getting location instance records >%v<", err)
		return nil, err
	}
	return locationInstanceRecs, nil
}

func (t *Testing) getMonsterInstanceRecs(dungeonInstanceID string) ([]*record.MonsterInstance, error) {

	monsterInstanceRecs, err := t.Model.(*model.Model).GetMonsterInstanceRecs(
		map[string]interface{}{
			"dungeon_instance_id": dungeonInstanceID,
		},
		nil, false,
	)
	if err != nil {
		t.Log.Warn("(testing) Failed getting monster instance records >%v<", err)
		return nil, err
	}
	return monsterInstanceRecs, nil
}

func (t *Testing) getObjectInstanceRecs(dungeonInstanceID string) ([]*record.ObjectInstance, error) {

	objectInstanceRecs, err := t.Model.(*model.Model).GetObjectInstanceRecs(
		map[string]interface{}{
			"dungeon_instance_id": dungeonInstanceID,
		},
		nil, false,
	)
	if err != nil {
		t.Log.Warn("(testing) Failed getting object instance records >%v<", err)
		return nil, err
	}
	return objectInstanceRecs, nil
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
