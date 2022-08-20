package harness

import (
	"fmt"
	"strings"

	"gitlab.com/alienspaces/go-mud/server/core/nullstring"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// Data -
type Data struct {
	// Object
	ObjectRecs []*record.Object

	// Monster
	MonsterRecs       []*record.Monster
	MonsterObjectRecs []*record.MonsterObject

	// Character
	CharacterRecs       []*record.Character
	CharacterObjectRecs []*record.CharacterObject

	// Dungeon
	DungeonRecs         []*record.Dungeon
	LocationRecs        []*record.Location
	LocationObjectRecs  []*record.LocationObject
	LocationMonsterRecs []*record.LocationMonster

	// Dungeon Instance
	DungeonInstanceRecs   []*record.DungeonInstance
	LocationInstanceRecs  []*record.LocationInstance
	CharacterInstanceRecs []*record.CharacterInstance
	MonsterInstanceRecs   []*record.MonsterInstance
	ObjectInstanceRecs    []*record.ObjectInstance

	// Action
	ActionRecs                []*record.Action
	ActionCharacterRecs       []*record.ActionCharacter
	ActionCharacterObjectRecs []*record.ActionCharacterObject
	ActionMonsterRecs         []*record.ActionMonster
	ActionMonsterObjectRecs   []*record.ActionMonsterObject
	ActionObjectRecs          []*record.ActionObject
}

// Object
func (d *Data) AddObjectRec(rec *record.Object) {
	for idx := range d.ObjectRecs {
		if d.ObjectRecs[idx].ID == rec.ID {
			d.ObjectRecs[idx] = rec
			return
		}
	}
	d.ObjectRecs = append(d.ObjectRecs, rec)
}

func (d *Data) GetObjectRecByID(objectID string) (*record.Object, error) {
	for _, rec := range d.ObjectRecs {
		if rec.ID == objectID {
			return rec, nil
		}
	}

	return nil, fmt.Errorf("failed getting object with ID >%s<", objectID)
}

func (d *Data) GetObjectRecByName(objectName string) (*record.Object, error) {

	for idx := range d.ObjectRecs {
		if d.ObjectRecs[idx].Name == objectName {
			return d.ObjectRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("failed getting object with Name >%s<", objectName)
}

// Monster
func (d *Data) AddMonsterRec(rec *record.Monster) {
	for idx := range d.MonsterRecs {
		if d.MonsterRecs[idx].ID == rec.ID {
			d.MonsterRecs[idx] = rec
			return
		}
	}
	d.MonsterRecs = append(d.MonsterRecs, rec)
}

func (d *Data) GetMonsterRecByID(monsterID string) (*record.Monster, error) {

	for idx := range d.MonsterRecs {
		if d.MonsterRecs[idx].ID == monsterID {
			return d.MonsterRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("failed getting monster by ID >%s<", monsterID)
}

func (d *Data) GetMonsterRecByName(monsterName string) (*record.Monster, error) {

	for idx := range d.MonsterRecs {
		if d.MonsterRecs[idx].Name == monsterName {
			return d.MonsterRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("failed getting monster by Name >%s<", monsterName)
}

func (d *Data) AddMonsterObjectRec(rec *record.MonsterObject) {
	for idx := range d.MonsterObjectRecs {
		if d.MonsterObjectRecs[idx].ID == rec.ID {
			d.MonsterObjectRecs[idx] = rec
			return
		}
	}
	d.MonsterObjectRecs = append(d.MonsterObjectRecs, rec)
}

// Character
func (d *Data) AddCharacterRec(rec *record.Character) {
	for idx := range d.CharacterRecs {
		if d.CharacterRecs[idx].ID == rec.ID {
			d.CharacterRecs[idx] = rec
			return
		}
	}
	d.CharacterRecs = append(d.CharacterRecs, rec)
}

func (d *Data) GetCharacterRecByID(characterID string) (*record.Character, error) {

	for idx := range d.CharacterRecs {
		if d.CharacterRecs[idx].ID == characterID {
			return d.CharacterRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("failed getting character by ID >%s<", characterID)
}

func (d *Data) GetCharacterRecByName(characterName string) (*record.Character, error) {

	for idx := range d.CharacterRecs {
		if d.CharacterRecs[idx].Name == characterName {
			return d.CharacterRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("failed getting character by Name >%s<", characterName)
}

func (d *Data) AddCharacterObjectRec(rec *record.CharacterObject) {
	for idx := range d.CharacterObjectRecs {
		if d.CharacterObjectRecs[idx].ID == rec.ID {
			d.CharacterObjectRecs[idx] = rec
			return
		}
	}
	d.CharacterObjectRecs = append(d.CharacterObjectRecs, rec)
}

// Dungeon
func (d *Data) AddDungeonRec(rec *record.Dungeon) {
	for idx := range d.DungeonRecs {
		if d.DungeonRecs[idx].ID == rec.ID {
			d.DungeonRecs[idx] = rec
			return
		}
	}
	d.DungeonRecs = append(d.DungeonRecs, rec)
}

func (d *Data) GetDungeonRecByID(dungeonID string) (*record.Dungeon, error) {
	for idx := range d.DungeonRecs {
		if d.DungeonRecs[idx].ID == dungeonID {
			return d.DungeonRecs[idx], nil
		}
	}
	return nil, fmt.Errorf("unknown dungeon ID >%s<", dungeonID)
}

func (d *Data) GetDungeonRecByName(dungeonName string) (*record.Dungeon, error) {
	for idx := range d.DungeonRecs {
		if d.DungeonRecs[idx].Name == dungeonName {
			return d.DungeonRecs[idx], nil
		}
	}
	return nil, fmt.Errorf("unknown dungeon name >%s<", dungeonName)
}

// Location
func (d *Data) AddLocationRec(rec *record.Location) {
	for idx := range d.LocationRecs {
		if d.LocationRecs[idx].ID == rec.ID {
			d.LocationRecs[idx] = rec
			return
		}
	}
	d.LocationRecs = append(d.LocationRecs, rec)
}

func (d *Data) GetLocationRecByID(locationID string) (*record.Location, error) {
	for idx := range d.LocationRecs {
		if d.LocationRecs[idx].ID == locationID {
			return d.LocationRecs[idx], nil
		}
	}
	return nil, fmt.Errorf("unknown location ID >%s<", locationID)
}

// Location Object
func (d *Data) AddLocationObjectRec(rec *record.LocationObject) {
	for idx := range d.LocationObjectRecs {
		if d.LocationObjectRecs[idx].ID == rec.ID {
			d.LocationObjectRecs[idx] = rec
			return
		}
	}
	d.LocationObjectRecs = append(d.LocationObjectRecs, rec)
}

func (d *Data) AddLocationMonsterRec(rec *record.LocationMonster) {
	for idx := range d.LocationMonsterRecs {
		if d.LocationMonsterRecs[idx].ID == rec.ID {
			d.LocationMonsterRecs[idx] = rec
			return
		}
	}
	d.LocationMonsterRecs = append(d.LocationMonsterRecs, rec)
}

// Object Instance
func (d *Data) AddObjectInstanceRec(rec *record.ObjectInstance) {
	for idx := range d.ObjectInstanceRecs {
		if d.ObjectInstanceRecs[idx].ID == rec.ID {
			d.ObjectInstanceRecs[idx] = rec
			return
		}
	}
	d.ObjectInstanceRecs = append(d.ObjectInstanceRecs, rec)
}

func (d *Data) GetObjectInstanceRecByName(name string) (*record.ObjectInstance, error) {
	for _, rec := range d.ObjectInstanceRecs {
		oRec, err := d.GetObjectRecByID(rec.ObjectID)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(oRec.Name, name) {
			return rec, nil
		}
	}
	return nil, nil
}

func (d *Data) GetObjectInstanceRecsByLocationInstanceID(locationInstanceID string) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.ToString(rec.LocationInstanceID) == locationInstanceID {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetObjectInstanceRecsByCharacterInstanceID(characterInstanceID string) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.ToString(rec.CharacterInstanceID) == characterInstanceID {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetEquippedObjectInstanceRecsByCharacterInstanceID(characterInstanceID string) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.ToString(rec.CharacterInstanceID) == characterInstanceID &&
			rec.IsEquipped {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetStashedObjectInstanceRecsByCharacterInstanceID(characterInstanceID string) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.ToString(rec.CharacterInstanceID) == characterInstanceID &&
			rec.IsStashed {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetObjectInstanceRecsByMonsterInstanceID(monsterInstanceID string) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.ToString(rec.MonsterInstanceID) == monsterInstanceID {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetEquippedObjectInstanceRecsByMonsterInstanceID(monsterInstanceID string) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.ToString(rec.MonsterInstanceID) == monsterInstanceID &&
			rec.IsEquipped {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetStashedObjectInstanceRecsByMonsterInstanceID(monsterInstanceID string) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.ToString(rec.MonsterInstanceID) == monsterInstanceID &&
			rec.IsStashed {
			recs = append(recs, rec)
		}
	}
	return recs
}

// TODO: Experimental, can potentially use generics here?
func (d *Data) GetMatchingObjectInstanceRecs(mrec *record.ObjectInstance) []*record.ObjectInstance {

	recs := []*record.ObjectInstance{}
	matched := false
	for _, rec := range d.ObjectInstanceRecs {
		if nullstring.IsValid(mrec.LocationInstanceID) {
			if nullstring.ToString(mrec.LocationInstanceID) == nullstring.ToString(rec.LocationInstanceID) {
				matched = true
			} else {
				matched = false
			}
		}
		if nullstring.IsValid(mrec.CharacterInstanceID) {
			if nullstring.ToString(mrec.CharacterInstanceID) == nullstring.ToString(rec.CharacterInstanceID) {
				matched = true
			} else {
				matched = false
			}
		}
		if nullstring.IsValid(mrec.MonsterInstanceID) {
			if nullstring.ToString(mrec.MonsterInstanceID) == nullstring.ToString(rec.MonsterInstanceID) {
				matched = true
			} else {
				matched = false
			}
		}
		if mrec.DungeonInstanceID != "" {
			if mrec.DungeonInstanceID == rec.DungeonInstanceID {
				matched = true
			} else {
				matched = false
			}
		}

		if matched {
			recs = append(recs, rec)
		}
	}

	return recs
}

// Monster Instance
func (d *Data) AddMonsterInstanceRec(rec *record.MonsterInstance) {
	for idx := range d.MonsterInstanceRecs {
		if d.MonsterInstanceRecs[idx].ID == rec.ID {
			d.MonsterInstanceRecs[idx] = rec
			return
		}
	}
	d.MonsterInstanceRecs = append(d.MonsterInstanceRecs, rec)
}

// Character Instance
func (d *Data) AddCharacterInstanceRec(rec *record.CharacterInstance) {
	for idx := range d.CharacterInstanceRecs {
		if d.CharacterInstanceRecs[idx].ID == rec.ID {
			d.CharacterInstanceRecs[idx] = rec
			return
		}
	}
	d.CharacterInstanceRecs = append(d.CharacterInstanceRecs, rec)
}

// Dungeon Instance
func (d *Data) AddDungeonInstanceRec(rec *record.DungeonInstance) {
	for idx := range d.DungeonInstanceRecs {
		if d.DungeonInstanceRecs[idx].ID == rec.ID {
			d.DungeonInstanceRecs[idx] = rec
			return
		}
	}
	d.DungeonInstanceRecs = append(d.DungeonInstanceRecs, rec)
}

func (d *Data) AddDungeonInstanceRecordSet(rs *model.DungeonInstanceRecordSet) {

	d.AddDungeonInstanceRec(rs.DungeonInstanceRec)

	for idx := range rs.LocationInstanceRecs {
		d.AddLocationInstanceRec(rs.LocationInstanceRecs[idx])
	}
	for idx := range rs.ObjectInstanceRecs {
		d.AddObjectInstanceRec(rs.ObjectInstanceRecs[idx])
	}
	for idx := range rs.MonsterInstanceRecs {
		d.AddMonsterInstanceRec(rs.MonsterInstanceRecs[idx])
	}
	for idx := range rs.CharacterInstanceRecs {
		d.AddCharacterInstanceRec(rs.CharacterInstanceRecs[idx])
	}
}

func (d *Data) AddCharacterInstanceRecordSet(rs *model.CharacterInstanceRecordSet) {

	d.AddCharacterInstanceRec(rs.CharacterInstanceRec)

	for idx := range rs.ObjectInstanceRecs {
		d.AddObjectInstanceRec(rs.ObjectInstanceRecs[idx])
	}
}

func (d *Data) GetDungeonInstanceRecByName(dungeonInstanceName string) (*record.DungeonInstance, error) {

	for _, rec := range d.DungeonInstanceRecs {
		dungeonRec, err := d.GetDungeonRecByID(rec.DungeonID)
		if err != nil {
			return nil, err
		}
		if dungeonRec.Name == dungeonInstanceName {
			return rec, nil
		}
	}

	return nil, fmt.Errorf("failed getting dungeon instance with name >%s<", dungeonInstanceName)
}

// Location Instance
func (d *Data) AddLocationInstanceRec(rec *record.LocationInstance) {
	for idx := range d.LocationInstanceRecs {
		if d.LocationInstanceRecs[idx].ID == rec.ID {
			d.LocationInstanceRecs[idx] = rec
			return
		}
	}
	d.LocationInstanceRecs = append(d.LocationInstanceRecs, rec)
}

func (d *Data) GetLocationInstanceRecByName(name string) (*record.LocationInstance, error) {
	for _, rec := range d.LocationInstanceRecs {
		lRec, err := d.GetLocationRecByID(rec.LocationID)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(lRec.Name, name) {
			return rec, nil
		}
	}
	return nil, nil
}

// teardownData -
type teardownData struct {
	// Object
	ObjectRecs []record.Object

	// Monster
	MonsterRecs       []record.Monster
	MonsterObjectRecs []record.MonsterObject

	// Character
	CharacterRecs       []record.Character
	CharacterObjectRecs []record.CharacterObject

	// Dungeon
	DungeonRecs         []record.Dungeon
	LocationRecs        []record.Location
	LocationObjectRecs  []record.LocationObject
	LocationMonsterRecs []record.LocationMonster

	// Dungeon Instance
	DungeonInstanceRecs   []*record.DungeonInstance
	LocationInstanceRecs  []*record.LocationInstance
	CharacterInstanceRecs []*record.CharacterInstance
	MonsterInstanceRecs   []*record.MonsterInstance
	ObjectInstanceRecs    []*record.ObjectInstance

	// Action
	ActionRecs                []*record.Action
	ActionCharacterRecs       []*record.ActionCharacter
	ActionCharacterObjectRecs []*record.ActionCharacterObject
	ActionMonsterRecs         []*record.ActionMonster
	ActionMonsterObjectRecs   []*record.ActionMonsterObject
	ActionObjectRecs          []*record.ActionObject
}

func (d *teardownData) AddDungeonInstanceRecordSet(rs *model.DungeonInstanceRecordSet) {
	d.DungeonInstanceRecs = append(d.DungeonInstanceRecs, rs.DungeonInstanceRec)
	d.LocationInstanceRecs = append(d.LocationInstanceRecs, rs.LocationInstanceRecs...)
	d.ObjectInstanceRecs = append(d.ObjectInstanceRecs, rs.ObjectInstanceRecs...)
	d.MonsterInstanceRecs = append(d.MonsterInstanceRecs, rs.MonsterInstanceRecs...)
	d.CharacterInstanceRecs = append(d.CharacterInstanceRecs, rs.CharacterInstanceRecs...)
}

func (d *teardownData) AddCharacterInstanceRecordSet(rs *model.CharacterInstanceRecordSet) {
	d.CharacterInstanceRecs = append(d.CharacterInstanceRecs, rs.CharacterInstanceRec)
	d.ObjectInstanceRecs = append(d.ObjectInstanceRecs, rs.ObjectInstanceRecs...)
}
