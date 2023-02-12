package model

import (
	"fmt"
	"math/rand"

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

	l.Info("Deciding action args >%#v<", args)

	if args.MonsterInstanceViewRec != nil {
		l.Info("Deciding action for monster name >%s<", args.MonsterInstanceViewRec.Name)
	} else {
		l.Info("Deciding action for character name >%s<", args.CharacterInstanceViewRec.Name)
	}

	// Decider functions are typically prioritised as attack if anything is
	// worth attacking, then grab anything thats worth grabbing, look into
	// other rooms to find someone to attack, and then the move if there's
	// somewhere worth moving to.
	deciderFuncs := []func(args *DeciderArgs) (string, error){
		m.decideActionAttack,
		m.decideActionStash,
		m.decideActionLook,
		m.decideActionMove,
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

func (m *Model) decideActionAttack(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionAttack")

	lirs := args.LocationInstanceRecordSet

	action := ""
	if args.MonsterInstanceViewRec != nil {
		l.Info("Character instance count >%d<", len(lirs.CharacterInstanceViewRecs))
		if len(lirs.CharacterInstanceViewRecs) != 0 {
			v := rand.Intn(len(lirs.CharacterInstanceViewRecs))
			action = fmt.Sprintf("attack %s", lirs.CharacterInstanceViewRecs[v].Name)
		}
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionStash(args *DeciderArgs) (string, error) {
	return "", nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionLook(args *DeciderArgs) (string, error) {

	return "", nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionMove(args *DeciderArgs) (string, error) {
	return "", nil
}
