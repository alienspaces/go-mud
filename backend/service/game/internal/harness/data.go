package harness

import (
	"fmt"
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
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

	// Instance
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

	// Turn
	TurnRecs []*record.Turn
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
		if strings.EqualFold(NormalName(d.ObjectRecs[idx].Name), objectName) {
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
		if strings.EqualFold(NormalName(d.MonsterRecs[idx].Name), monsterName) {
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
		if strings.EqualFold(NormalName(d.CharacterRecs[idx].Name), characterName) {
			return d.CharacterRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("failed getting character by Name >%s<", characterName)
}

// CharacterObject
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
		if strings.EqualFold(NormalName(d.DungeonRecs[idx].Name), dungeonName) {
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

func (d *Data) GetLocationRecByName(name string) (*record.Location, error) {
	for idx := range d.LocationRecs {
		if strings.EqualFold(NormalName(d.LocationRecs[idx].Name), name) {
			return d.LocationRecs[idx], nil
		}
	}
	return nil, fmt.Errorf("unknown location Name >%s<", name)
}

// LocationObject
func (d *Data) AddLocationObjectRec(rec *record.LocationObject) {
	for idx := range d.LocationObjectRecs {
		if d.LocationObjectRecs[idx].ID == rec.ID {
			d.LocationObjectRecs[idx] = rec
			return
		}
	}
	d.LocationObjectRecs = append(d.LocationObjectRecs, rec)
}

// LocationMonster
func (d *Data) AddLocationMonsterRec(rec *record.LocationMonster) {
	for idx := range d.LocationMonsterRecs {
		if d.LocationMonsterRecs[idx].ID == rec.ID {
			d.LocationMonsterRecs[idx] = rec
			return
		}
	}
	d.LocationMonsterRecs = append(d.LocationMonsterRecs, rec)
}

// ObjectInstance
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
		if strings.EqualFold(NormalName(oRec.Name), name) {
			return rec, nil
		}
	}
	return nil, nil
}

func (d *Data) GetObjectInstanceRecsByLocationInstanceID(locationInstanceID string) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringToString(rec.LocationInstanceID) == locationInstanceID {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetObjectInstanceRecsByCharacterInstanceID(characterInstanceID string) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringToString(rec.CharacterInstanceID) == characterInstanceID {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetEquippedObjectInstanceRecsByCharacterInstanceID(characterInstanceID string) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringToString(rec.CharacterInstanceID) == characterInstanceID &&
			rec.IsEquipped {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetStashedObjectInstanceRecsByCharacterInstanceID(characterInstanceID string) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringToString(rec.CharacterInstanceID) == characterInstanceID &&
			rec.IsStashed {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetObjectInstanceRecsByMonsterInstanceID(monsterInstanceID string) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringToString(rec.MonsterInstanceID) == monsterInstanceID {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetEquippedObjectInstanceRecsByMonsterInstanceID(monsterInstanceID string) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringToString(rec.MonsterInstanceID) == monsterInstanceID &&
			rec.IsEquipped {
			recs = append(recs, rec)
		}
	}
	return recs
}

func (d *Data) GetStashedObjectInstanceRecsByMonsterInstanceID(monsterInstanceID string) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringToString(rec.MonsterInstanceID) == monsterInstanceID &&
			rec.IsStashed {
			recs = append(recs, rec)
		}
	}
	return recs
}

// TODO: (game) Experimental, can potentially use generics here?
func (d *Data) GetMatchingObjectInstanceRecs(mrec *record.ObjectInstance) []*record.ObjectInstance {
	recs := []*record.ObjectInstance{}
	matched := false
	for _, rec := range d.ObjectInstanceRecs {
		if null.NullStringIsValid(mrec.LocationInstanceID) {
			if null.NullStringToString(mrec.LocationInstanceID) == null.NullStringToString(rec.LocationInstanceID) {
				matched = true
			} else {
				matched = false
			}
		}
		if null.NullStringIsValid(mrec.CharacterInstanceID) {
			if null.NullStringToString(mrec.CharacterInstanceID) == null.NullStringToString(rec.CharacterInstanceID) {
				matched = true
			} else {
				matched = false
			}
		}
		if null.NullStringIsValid(mrec.MonsterInstanceID) {
			if null.NullStringToString(mrec.MonsterInstanceID) == null.NullStringToString(rec.MonsterInstanceID) {
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

// MonsterInstance
func (d *Data) AddMonsterInstanceRec(rec *record.MonsterInstance) {
	for idx := range d.MonsterInstanceRecs {
		if d.MonsterInstanceRecs[idx].ID == rec.ID {
			d.MonsterInstanceRecs[idx] = rec
			return
		}
	}
	d.MonsterInstanceRecs = append(d.MonsterInstanceRecs, rec)
}

func (d *Data) GetMonsterInstanceRecByName(name string) (*record.MonsterInstance, error) {
	for _, miRec := range d.MonsterInstanceRecs {
		mRec, err := d.GetMonsterRecByID(miRec.MonsterID)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(NormalName(mRec.Name), name) {
			return miRec, nil
		}
	}
	return nil, fmt.Errorf("failed getting character instance with name >%s<", name)
}

// CharacterInstance
func (d *Data) AddCharacterInstanceRec(rec *record.CharacterInstance) {
	for idx := range d.CharacterInstanceRecs {
		if d.CharacterInstanceRecs[idx].ID == rec.ID {
			d.CharacterInstanceRecs[idx] = rec
			return
		}
	}
	d.CharacterInstanceRecs = append(d.CharacterInstanceRecs, rec)
}

func (d *Data) GetCharacterInstanceRecByName(name string) (*record.CharacterInstance, error) {
	for _, ciRec := range d.CharacterInstanceRecs {
		cRec, err := d.GetCharacterRecByID(ciRec.CharacterID)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(NormalName(cRec.Name), name) {
			return ciRec, nil
		}
	}
	return nil, fmt.Errorf("failed getting character instance with name >%s<", name)
}

// DungeonInstance
func (d *Data) AddDungeonInstanceRec(rec *record.DungeonInstance) {
	for idx := range d.DungeonInstanceRecs {
		if d.DungeonInstanceRecs[idx].ID == rec.ID {
			d.DungeonInstanceRecs[idx] = rec
			return
		}
	}
	d.DungeonInstanceRecs = append(d.DungeonInstanceRecs, rec)
}

func (d *Data) GetDungeonInstanceRecByName(name string) (*record.DungeonInstance, error) {
	for _, rec := range d.DungeonInstanceRecs {
		dungeonRec, err := d.GetDungeonRecByID(rec.DungeonID)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(NormalName(dungeonRec.Name), name) {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting dungeon instance with name >%s<", name)
}

// DungeonInstanceRecordSet
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

// CharacterInstanceRecordSet
func (d *Data) AddCharacterInstanceRecordSet(rs *model.CharacterInstanceRecordSet) {

	d.AddCharacterInstanceRec(rs.CharacterInstanceRec)

	for idx := range rs.ObjectInstanceRecs {
		d.AddObjectInstanceRec(rs.ObjectInstanceRecs[idx])
	}
}

// LocationInstance
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
		if strings.EqualFold(NormalName(lRec.Name), name) {
			return rec, nil
		}
	}
	return nil, nil
}

// Action
func (d *Data) AddActionRec(rec *record.Action) {
	for idx := range d.ActionRecs {
		if d.ActionRecs[idx].ID == rec.ID {
			d.ActionRecs[idx] = rec
			return
		}
	}
	d.ActionRecs = append(d.ActionRecs, rec)
}

// ActionCharacter
func (d *Data) AddActionCharacterRec(rec *record.ActionCharacter) {
	for idx := range d.ActionCharacterRecs {
		if d.ActionCharacterRecs[idx].ID == rec.ID {
			d.ActionCharacterRecs[idx] = rec
			return
		}
	}
	d.ActionCharacterRecs = append(d.ActionCharacterRecs, rec)
}

// ActionCharacterObject
func (d *Data) AddActionCharacterObjectRec(rec *record.ActionCharacterObject) {
	for idx := range d.ActionCharacterObjectRecs {
		if d.ActionCharacterObjectRecs[idx].ID == rec.ID {
			d.ActionCharacterObjectRecs[idx] = rec
			return
		}
	}
	d.ActionCharacterObjectRecs = append(d.ActionCharacterObjectRecs, rec)
}

// ActionMonster
func (d *Data) AddActionMonsterRec(rec *record.ActionMonster) {
	for idx := range d.ActionMonsterRecs {
		if d.ActionMonsterRecs[idx].ID == rec.ID {
			d.ActionMonsterRecs[idx] = rec
			return
		}
	}
	d.ActionMonsterRecs = append(d.ActionMonsterRecs, rec)
}

// ActionMonsterObject
func (d *Data) AddActionMonsterObjectRec(rec *record.ActionMonsterObject) {
	for idx := range d.ActionMonsterObjectRecs {
		if d.ActionMonsterObjectRecs[idx].ID == rec.ID {
			d.ActionMonsterObjectRecs[idx] = rec
			return
		}
	}
	d.ActionMonsterObjectRecs = append(d.ActionMonsterObjectRecs, rec)
}

// ActionObject
func (d *Data) AddActionObjectRec(rec *record.ActionObject) {
	for idx := range d.ActionObjectRecs {
		if d.ActionObjectRecs[idx].ID == rec.ID {
			d.ActionObjectRecs[idx] = rec
			return
		}
	}
	d.ActionObjectRecs = append(d.ActionObjectRecs, rec)
}

// Turn
func (d *Data) AddTurnRec(rec *record.Turn) {
	for idx := range d.TurnRecs {
		if d.TurnRecs[idx].ID == rec.ID {
			d.TurnRecs[idx] = rec
			return
		}
	}
	d.TurnRecs = append(d.TurnRecs, rec)
}

func (d *Data) AddActionLocationRecordSet(alrs *record.ActionLocationRecordSet) {

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

func (d *Data) AddActionRecordSet(rs *record.ActionRecordSet) {

	d.AddActionRec(rs.ActionRec)

	// Source
	if rs.ActionCharacterRec != nil {
		d.AddActionCharacterRec(rs.ActionCharacterRec)
		for idx := range rs.ActionCharacterObjectRecs {
			d.AddActionCharacterObjectRec(rs.ActionCharacterObjectRecs[idx])
		}
	}
	if rs.ActionMonsterRec != nil {
		d.AddActionMonsterRec(rs.ActionMonsterRec)
		for idx := range rs.ActionMonsterObjectRecs {
			d.AddActionMonsterObjectRec(rs.ActionMonsterObjectRecs[idx])
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
		for idx := range rs.TargetActionCharacterObjectRecs {
			d.AddActionCharacterObjectRec(rs.TargetActionCharacterObjectRecs[idx])
		}
	}
	if rs.TargetActionMonsterRec != nil {
		d.AddActionMonsterRec(rs.TargetActionMonsterRec)
		for idx := range rs.TargetActionMonsterObjectRecs {
			d.AddActionMonsterObjectRec(rs.TargetActionMonsterObjectRecs[idx])
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
