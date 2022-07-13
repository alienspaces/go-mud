package harness

import (
	"fmt"

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

func (d *Data) GetDungeonRecByName(objectName string) (*record.Dungeon, error) {

	for idx := range d.DungeonRecs {
		if d.DungeonRecs[idx].Name == objectName {
			return d.DungeonRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("unknown dungeon name >%s<", objectName)
}

func (d *Data) GetObjectRecByName(objectName string) (*record.Object, error) {

	for idx := range d.ObjectRecs {
		if d.ObjectRecs[idx].Name == objectName {
			return d.ObjectRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("unknown object name >%s<", objectName)
}

func (d *Data) GetMonsterRecByName(monsterName string) (*record.Monster, error) {

	for idx := range d.MonsterRecs {
		if d.MonsterRecs[idx].Name == monsterName {
			return d.MonsterRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("unknown monster name >%s<", monsterName)
}

func (d *Data) AddObjectRec(rec *record.Object) {
	for idx := range d.ObjectRecs {
		if d.ObjectRecs[idx].ID == rec.ID {
			d.ObjectRecs[idx] = rec
			return
		}
	}
	d.ObjectRecs = append(d.ObjectRecs, rec)
}

func (d *Data) AddMonsterRec(rec *record.Monster) {
	for idx := range d.MonsterRecs {
		if d.MonsterRecs[idx].ID == rec.ID {
			d.MonsterRecs[idx] = rec
			return
		}
	}
	d.MonsterRecs = append(d.MonsterRecs, rec)
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

func (d *Data) AddCharacterRec(rec *record.Character) {
	for idx := range d.CharacterRecs {
		if d.CharacterRecs[idx].ID == rec.ID {
			d.CharacterRecs[idx] = rec
			return
		}
	}
	d.CharacterRecs = append(d.CharacterRecs, rec)
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

func (d *Data) AddDungeonRec(rec *record.Dungeon) {
	for idx := range d.DungeonRecs {
		if d.DungeonRecs[idx].ID == rec.ID {
			d.DungeonRecs[idx] = rec
			return
		}
	}
	d.DungeonRecs = append(d.DungeonRecs, rec)
}

func (d *Data) AddLocationRec(rec *record.Location) {
	for idx := range d.LocationRecs {
		if d.LocationRecs[idx].ID == rec.ID {
			d.LocationRecs[idx] = rec
			return
		}
	}
	d.LocationRecs = append(d.LocationRecs, rec)
}

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

func (d *Data) AddDungeonInstanceRec(rec *record.DungeonInstance) {
	for idx := range d.DungeonInstanceRecs {
		if d.DungeonInstanceRecs[idx].ID == rec.ID {
			d.DungeonInstanceRecs[idx] = rec
			return
		}
	}
	d.DungeonInstanceRecs = append(d.DungeonInstanceRecs, rec)
}

func (d *Data) AddLocationInstanceRec(rec *record.LocationInstance) {
	for idx := range d.LocationInstanceRecs {
		if d.LocationInstanceRecs[idx].ID == rec.ID {
			d.LocationInstanceRecs[idx] = rec
			return
		}
	}
	d.LocationInstanceRecs = append(d.LocationInstanceRecs, rec)
}

func (d *Data) AddObjectInstanceRec(rec *record.ObjectInstance) {
	for idx := range d.ObjectInstanceRecs {
		if d.ObjectInstanceRecs[idx].ID == rec.ID {
			d.ObjectInstanceRecs[idx] = rec
			return
		}
	}
	d.ObjectInstanceRecs = append(d.ObjectInstanceRecs, rec)
}

func (d *Data) AddMonsterInstanceRec(rec *record.MonsterInstance) {
	for idx := range d.MonsterInstanceRecs {
		if d.MonsterInstanceRecs[idx].ID == rec.ID {
			d.MonsterInstanceRecs[idx] = rec
			return
		}
	}
	d.MonsterInstanceRecs = append(d.MonsterInstanceRecs, rec)
}

func (d *Data) AddCharacterInstanceRec(rec *record.CharacterInstance) {
	for idx := range d.CharacterInstanceRecs {
		if d.CharacterInstanceRecs[idx].ID == rec.ID {
			d.CharacterInstanceRecs[idx] = rec
			return
		}
	}
	d.CharacterInstanceRecs = append(d.CharacterInstanceRecs, rec)
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
