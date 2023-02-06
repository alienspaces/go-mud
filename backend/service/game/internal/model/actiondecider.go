package model

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// TODO: 9-implement-monster-actions
// NOTES: Monsters should behave based on their statistics and capabilities.
// Examples:
// - Intelligent humanoid monsters with no equipped weapon should attempt to equip any
//   weapons left lying in the current location.
//   -- Requires item types
// - Intelligent monsters that are not defending anything other than themselves should
//   attempt to run away when losing a fight
//   -- Compare own health with health of characters in the room that have previously
//      attacked it
// - Intelligent monsters that are obviously out of their depth in a battle should run
//   away.

type DeciderArgs struct {
	MonsterInstanceViewRec    *record.MonsterInstanceView
	CharacterInstanceViewRec  *record.CharacterInstanceView
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
}

func (m *Model) decideAction(args *DeciderArgs) (string, error) {
	l := m.Logger("decideAction")

	l.Info("deciding action args >%#v<", args)

	if args.MonsterInstanceViewRec != nil {
		l.Info("deciding action for monster name >%s<", args.MonsterInstanceViewRec.Name)
	} else {
		l.Info("deciding action for character name >%s<", args.CharacterInstanceViewRec.Name)
	}

	deciderFuncs := []func(args *DeciderArgs) (string, error){
		m.decideActionMove,
		m.decideActionStash,
		m.decideActionAttack,
	}

	var err error
	var sentence string

	for idx := range deciderFuncs {
		sentence, err = deciderFuncs[idx](args)
		if err != nil {
			return "", err
		}
		if sentence != "" {
			break
		}
	}

	return sentence, nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionMove(args *DeciderArgs) (string, error) {
	return "", nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionStash(args *DeciderArgs) (string, error) {
	return "", nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionAttack(args *DeciderArgs) (string, error) {
	return "", nil
}
