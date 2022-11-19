package harness

import (
	"gitlab.com/alienspaces/go-mud/server/core/repository"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// teardownData -
type teardownData struct {
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

func (d *teardownData) AddDungeonInstanceRecordSet(rs *model.DungeonInstanceRecordSet) {
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

func (d *teardownData) AddCharacterInstanceRecordSet(rs *model.CharacterInstanceRecordSet) {
	d.AddCharacterInstanceRec(rs.CharacterInstanceRec)
	for idx := range rs.ObjectInstanceRecs {
		d.AddObjectInstanceRec(rs.ObjectInstanceRecs[idx])
	}
}

func (d *teardownData) AddObjectRec(rec *record.Object) {
	for _, r := range d.ObjectRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ObjectRecs = append(d.ObjectRecs, &record.Object{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddMonsterRec(rec *record.Monster) {
	for _, r := range d.MonsterRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.MonsterRecs = append(d.MonsterRecs, &record.Monster{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddMonsterObjectRec(rec *record.MonsterObject) {
	for _, r := range d.MonsterObjectRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.MonsterObjectRecs = append(d.MonsterObjectRecs, &record.MonsterObject{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddCharacterRec(rec *record.Character) {
	for _, r := range d.CharacterRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.CharacterRecs = append(d.CharacterRecs, &record.Character{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddCharacterObjectRec(rec *record.CharacterObject) {
	for _, r := range d.CharacterObjectRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.CharacterObjectRecs = append(d.CharacterObjectRecs, &record.CharacterObject{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddDungeonRec(rec *record.Dungeon) {
	for _, r := range d.DungeonRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.DungeonRecs = append(d.DungeonRecs, &record.Dungeon{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddLocationRec(rec *record.Location) {
	for _, r := range d.LocationRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.LocationRecs = append(d.LocationRecs, &record.Location{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddLocationObjectRec(rec *record.LocationObject) {
	for _, r := range d.LocationObjectRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.LocationObjectRecs = append(d.LocationObjectRecs, &record.LocationObject{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddLocationMonsterRec(rec *record.LocationMonster) {
	for _, r := range d.LocationMonsterRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.LocationMonsterRecs = append(d.LocationMonsterRecs, &record.LocationMonster{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddDungeonInstanceRec(rec *record.DungeonInstance) {
	for _, r := range d.DungeonInstanceRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.DungeonInstanceRecs = append(d.DungeonInstanceRecs, &record.DungeonInstance{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddLocationInstanceRec(rec *record.LocationInstance) {
	for _, r := range d.LocationInstanceRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.LocationInstanceRecs = append(d.LocationInstanceRecs, &record.LocationInstance{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddCharacterInstanceRec(rec *record.CharacterInstance) {
	for _, r := range d.CharacterInstanceRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.CharacterInstanceRecs = append(d.CharacterInstanceRecs, &record.CharacterInstance{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddMonsterInstanceRec(rec *record.MonsterInstance) {
	for _, r := range d.MonsterInstanceRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.MonsterInstanceRecs = append(d.MonsterInstanceRecs, &record.MonsterInstance{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddObjectInstanceRec(rec *record.ObjectInstance) {
	for _, r := range d.ObjectInstanceRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ObjectInstanceRecs = append(d.ObjectInstanceRecs, &record.ObjectInstance{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddActionRec(rec *record.Action) {
	for _, r := range d.ActionRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ActionRecs = append(d.ActionRecs, &record.Action{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddActionCharacterRec(rec *record.ActionCharacter) {
	for _, r := range d.ActionCharacterRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ActionCharacterRecs = append(d.ActionCharacterRecs, &record.ActionCharacter{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddActionCharacterObjectRec(rec *record.ActionCharacterObject) {
	for _, r := range d.ActionCharacterObjectRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ActionCharacterObjectRecs = append(d.ActionCharacterObjectRecs, &record.ActionCharacterObject{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddActionMonsterRec(rec *record.ActionMonster) {
	for _, r := range d.ActionMonsterRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ActionMonsterRecs = append(d.ActionMonsterRecs, &record.ActionMonster{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddActionMonsterObjectRec(rec *record.ActionMonsterObject) {
	for _, r := range d.ActionMonsterObjectRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ActionMonsterObjectRecs = append(d.ActionMonsterObjectRecs, &record.ActionMonsterObject{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddActionObjectRec(rec *record.ActionObject) {
	for _, r := range d.ActionObjectRecs {
		if r.ID == rec.ID {
			return
		}
	}
	d.ActionObjectRecs = append(d.ActionObjectRecs, &record.ActionObject{Record: repository.Record{ID: rec.ID}})
}

func (d *teardownData) AddActionLocationRecordSet(alrs *record.ActionLocationRecordSet) {

	for _, rec := range alrs.ActionCharacterRecs {
		d.AddActionCharacterRec(rec)
	}
	for _, rec := range alrs.ActionMonsterRecs {
		d.AddActionMonsterRec(rec)
	}
	for _, rec := range alrs.ActionObjectRecs {
		d.AddActionObjectRec(rec)
	}
}

func (d *teardownData) AddActionRecordSet(rs *record.ActionRecordSet) {

	d.AddActionRec(rs.ActionRec)

	// Source
	if rs.ActionCharacterRec != nil {
		d.AddActionCharacterRec(rs.ActionCharacterRec)
		for _, rec := range rs.ActionCharacterObjectRecs {
			d.AddActionCharacterObjectRec(rec)
		}
	}
	if rs.ActionMonsterRec != nil {
		d.AddActionMonsterRec(rs.ActionMonsterRec)
		for _, rec := range rs.ActionMonsterObjectRecs {
			d.AddActionMonsterObjectRec(rec)
		}
	}

	// Current location
	if rs.CurrentLocation != nil {
		d.AddActionLocationRecordSet(rs.CurrentLocation)
	}

	// Target location
	if rs.TargetLocation != nil {
		d.AddActionLocationRecordSet(rs.TargetLocation)
	}

	// Targets
	if rs.TargetActionCharacterRec != nil {
		d.AddActionCharacterRec(rs.TargetActionCharacterRec)
		for _, rec := range rs.TargetActionCharacterObjectRecs {
			d.AddActionCharacterObjectRec(rec)
		}
	}
	if rs.TargetActionMonsterRec != nil {
		d.AddActionMonsterRec(rs.TargetActionMonsterRec)
		for _, rec := range rs.TargetActionMonsterObjectRecs {
			d.AddActionMonsterObjectRec(rec)
		}
	}
	if rs.EquippedActionObjectRec != nil {
		d.AddActionObjectRec(rs.EquippedActionObjectRec)
	}
	if rs.StashedActionObjectRec != nil {
		d.AddActionObjectRec(rs.StashedActionObjectRec)
	}
	if rs.DroppedActionObjectRec != nil {
		d.AddActionObjectRec(rs.DroppedActionObjectRec)
	}
	if rs.TargetActionObjectRec != nil {
		d.AddActionObjectRec(rs.TargetActionObjectRec)
	}
}
