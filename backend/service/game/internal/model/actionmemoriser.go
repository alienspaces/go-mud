package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// TODO: 14-implement-smarter-monsters

type MemoriserArgs struct {
	ActionRecordSet *record.ActionRecordSet
}

func (m *Model) memoriseAction(args *MemoriserArgs) ([]*record.ActionMemory, error) {
	l := m.Logger("memoriseAction")

	if args == nil || args.ActionRecordSet == nil ||
		args.ActionRecordSet.ActionRec == nil ||
		args.ActionRecordSet.ActionCharacterRec == nil ||
		args.ActionRecordSet.ActionMonsterRec == nil {
		err := fmt.Errorf("missing memoriser arguments >%#v<, cannot memorise action", args)
		l.Warn(err.Error())
		return nil, err
	}

	l.Info("Memorising action ID >%s<", args.ActionRecordSet.ActionRec)

	// Memory mapper functions
	memoryMapperFuncs := map[string]func(args *MemoriserArgs) ([]*record.ActionMemory, error){
		"move":   m.memoriseMoveActionMapper,
		"look":   m.memoriseLookActionMapper,
		"use":    m.memoriseUseActionMapper,
		"stash":  m.memoriseStashActionMapper,
		"equip":  m.memoriseEquipActionMapper,
		"drop":   m.memoriseDropActionMapper,
		"attack": m.memoriseAttackActionMapper,
	}

	command := args.ActionRecordSet.ActionRec.ResolvedCommand
	memoryMapperFunc, ok := memoryMapperFuncs[command]
	if !ok {
		err := fmt.Errorf("missing memory mapper function for command >%s<, cannot memorise action", command)
		l.Warn(err.Error())
		return nil, err
	}

	// Multiple memories may be recorded from an action that are then assessed in future actions
	recs, err := memoryMapperFunc(args)
	if err != nil {
		l.Warn("failed calling memory mapper function for command >%s< >%v<", command, err)
		return nil, err
	}

	for idx := range recs {
		rec := recs[idx]
		err = m.CreateActionMemoryRec(rec)
		if err != nil {
			l.Warn("failed creating monster instance memory record >%v<", command, err)
			return nil, err
		}
	}

	return recs, nil
}

// memoriseMoveActionMapper remembers locations a monster has moved from
func (m *Model) memoriseMoveActionMapper(args *MemoriserArgs) ([]*record.ActionMemory, error) {
	l := m.Logger("memoriseMoveActionMapper")

	recs := []*record.ActionMemory{}

	if args == nil || args.ActionRecordSet == nil ||
		args.ActionRecordSet.ActionMonsterRec == nil ||
		args.ActionRecordSet.CurrentLocation == nil ||
		args.ActionRecordSet.CurrentLocation.LocationInstanceViewRec == nil {
		err := fmt.Errorf("missing required move action mapper arguments >%#v<, cannot map move action memory", args)
		l.Warn(err.Error())
		return nil, err
	}

	amRec := args.ActionRecordSet.ActionMonsterRec
	livRec := args.ActionRecordSet.CurrentLocation.LocationInstanceViewRec

	// Moved from a location
	rec := &record.ActionMemory{
		MonsterInstanceID:        amRec.MonsterInstanceID,
		MemoryCommand:            record.ActionCommandMove,
		MemoryType:               record.ActionMemoryTypeLocation,
		MemoryLocationInstanceID: nullstring.FromString(livRec.ID),
	}
	recs = append(recs, rec)

	return recs, nil
}

// memoriseLookActionMapper memorises things a monster has seen
func (m *Model) memoriseLookActionMapper(args *MemoriserArgs) ([]*record.ActionMemory, error) {
	l := m.Logger("memoriseLookActionMapper")

	recs := []*record.ActionMemory{}

	if args == nil || args.ActionRecordSet == nil ||
		args.ActionRecordSet.ActionMonsterRec == nil {
		err := fmt.Errorf("missing required look action mapper arguments >%#v<, cannot map move action memory", args)
		l.Warn(err.Error())
		return nil, err
	}

	tl := args.ActionRecordSet.TargetLocation

	// Currently only remembering things a monster see's when looking into another location
	if tl == nil {
		l.Info("look action was not targetted at a location, not remembering anything")
		return nil, nil
	}

	amRec := args.ActionRecordSet.ActionMonsterRec
	livRec := tl.LocationInstanceViewRec
	acRecs := tl.ActionCharacterRecs

	// Looked into a location
	rec := &record.ActionMemory{
		MonsterInstanceID:        amRec.MonsterInstanceID,
		MemoryCommand:            record.ActionCommandLook,
		MemoryType:               record.ActionMemoryTypeLocation,
		MemoryLocationInstanceID: nullstring.FromString(livRec.ID),
	}
	recs = append(recs, rec)

	// Currently only remembering characters seen at a location
	for idx := range acRecs {
		rec := &record.ActionMemory{
			MonsterInstanceID:         amRec.MonsterInstanceID,
			MemoryCommand:             record.ActionCommandLook,
			MemoryType:                record.ActionMemoryTypeCharacter,
			MemoryLocationInstanceID:  nullstring.FromString(livRec.ID),
			MemoryCharacterInstanceID: nullstring.FromString(acRecs[idx].CharacterInstanceID),
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

// memoriseUseActionMapper doesn't remember anything yet
func (m *Model) memoriseUseActionMapper(args *MemoriserArgs) ([]*record.ActionMemory, error) {

	recs := []*record.ActionMemory{}

	return recs, nil
}

// memoriseStashActionMapper doesn't remember anything yet
func (m *Model) memoriseStashActionMapper(args *MemoriserArgs) ([]*record.ActionMemory, error) {

	recs := []*record.ActionMemory{}

	return recs, nil
}

// memoriseEquipActionMapper doesn't remember anything yet
func (m *Model) memoriseEquipActionMapper(args *MemoriserArgs) ([]*record.ActionMemory, error) {

	recs := []*record.ActionMemory{}

	return recs, nil
}

// memoriseDropActionMapper doesn't remember anything yet
func (m *Model) memoriseDropActionMapper(args *MemoriserArgs) ([]*record.ActionMemory, error) {

	recs := []*record.ActionMemory{}

	return recs, nil
}

// memoriseAttackActionMapper memorises characters the monster has attacked
func (m *Model) memoriseAttackActionMapper(args *MemoriserArgs) ([]*record.ActionMemory, error) {

	recs := []*record.ActionMemory{}

	return recs, nil
}
